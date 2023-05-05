package github

import (
	"context"
	"net/http"
	"sort"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v52/github"
	"github.com/rs/zerolog/log"

	"golang.org/x/sync/errgroup"

	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/errors"
)

type ClientOptions struct {
	AppID          int64      `name:"app_id" env:"GH_APP_ID" required:""`
	InstallationID int64      `name:"installation_id" env:"GH_INSTALLATION_ID" required:""`
	PrivateKey     PrivateKey `embed:"" prefix:"pk-" required:""`
}

type PrivateKey struct {
	File string `name:"file" env:"GH_PK_FILE"`
}

func (o *ClientOptions) HTTPClient() (*http.Client, error) {
	var (
		itr http.RoundTripper
		err error
	)

	if o.PrivateKey.File != "" {
		log.Info().Msgf("reading pk at path=%s", o.PrivateKey.File)
		itr, err = ghinstallation.NewKeyFromFile(http.DefaultTransport, o.AppID, o.InstallationID, o.PrivateKey.File)
	} else {
		itr = http.DefaultTransport
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Info().Msg("github client initialized")
	return &http.Client{Transport: itr}, nil
}

type Client struct {
	c *github.Client
}

func New(opts ClientOptions) (*Client, error) {
	c, err := opts.HTTPClient()
	if err != nil {
		return nil, err
	}

	return &Client{c: github.NewClient(c)}, nil
}

func (g *Client) ListFiles(ctx context.Context, r *ReviewParams) ([]*CommitFile, error) {
	var fs []*CommitFile
	return fs, cache.Ctx(ctx).Apply(ctx, &fs, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "files",
			Hashing: []any{r.Owner, r.Repo, r.Number},
		},
		Fetch: func() (any, error) {
			d, _, err := g.c.PullRequests.ListFiles(ctx, r.Owner, r.Repo, r.Number, nil)
			return d, errors.WrapGithubErr(err, "call to ListFiles failed")
		},
	})
}

func (g *Client) GetPullRequest(ctx context.Context, r *ReviewParams) (*PullRequest, error) {
	var pr *PullRequest
	return pr, cache.Apply(ctx, &pr, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "pr",
			Hashing: []any{r.Owner, r.Repo, r.Number},
		},
		Fetch: func() (any, error) {
			pr, _, err := g.c.PullRequests.Get(ctx, r.Owner, r.Repo, r.Number)
			return pr, errors.WrapGithubErr(err, "call to GetPullRequest failed")
		},
	})
}

func (g *Client) ListReviews(ctx context.Context, r *ReviewParams) ([]*PullRequestReview, error) {
	var reviews []*PullRequestReview
	return reviews, cache.Apply(ctx, &reviews, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "reviews",
			Hashing: []any{r.Owner, r.Repo, r.Number},
		},
		Fetch: func() (any, error) {
			reviews, _, err := g.c.PullRequests.ListReviews(ctx, r.Owner, r.Repo, r.Number, nil)
			return reviews, errors.WrapGithubErr(err, "call to ListReviews failed")
		},
	})
}

func (g *Client) GetReview(ctx context.Context, r *ReviewParams) (*PullRequestReview, error) {
	var review *PullRequestReview
	return review, cache.Apply(ctx, &review, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "review",
			Hashing: []any{r.Owner, r.Repo, r.Number, r.ReviewID},
		},
		Fetch: func() (any, error) {
			review, _, err := g.c.PullRequests.GetReview(ctx, r.Owner, r.Repo, r.Number, r.ReviewID)
			return review, errors.WrapGithubErr(err, "call to GetReview failed")
		},
	})
}

func (g *Client) ListReviewComments(ctx context.Context, r *ReviewParams) ([]*PullRequestComment, error) {
	var cs []*PullRequestComment
	return cs, cache.Apply(ctx, &cs, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "review-comments",
			Hashing: []any{r.Owner, r.Repo, r.Number, r.ReviewID},
		},
		Fetch: func() (any, error) {
			cs, _, err := g.c.PullRequests.ListReviewComments(ctx, r.Owner, r.Repo, r.Number, r.ReviewID, nil)
			return cs, errors.WrapGithubErr(err, "call to ListReviewComments failed")
		},
	})
}

func (g *Client) ListComments(ctx context.Context, r *ReviewParams) ([]*PullRequestComment, error) {
	var cs []*PullRequestComment
	err := cache.Apply(ctx, &cs, cache.Provider{
		DataKey: cache.DataKey{
			Prefix:  "pull-comments",
			Hashing: []any{r.Owner, r.Repo, r.Number},
		},
		Fetch: func() (any, error) {
			cs, _, err := g.c.PullRequests.ListComments(ctx, r.Owner, r.Repo, r.Number, nil)
			return cs, errors.WrapGithubErr(err, "call to ListComments failed")
		},
	})
	if err != nil {
		return nil, err
	}

	var ret []*PullRequestComment
	for _, c := range cs {
		if c.PullRequestReviewID == nil || *c.PullRequestReviewID != r.ReviewID {
			continue
		}
		ret = append(ret, c)
	}

	return ret, nil
}

func (g *Client) FetchReviewModel(ctx context.Context, r *ReviewParams) (*ReviewModel, error) {
	model := &ReviewModel{Params: r}
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		pull, err := g.GetPullRequest(ctx, r)
		if err == nil {
			model.PR = pull
		}
		return err
	})

	group.Go(func() error {
		review, err := g.GetReview(ctx, r)
		if err == nil {
			model.Review = review
		}
		return err
	})

	group.Go(func() error {
		comments, err := g.ListComments(ctx, r)
		if err == nil {
			model.Comments = comments
		}
		return err
	})

	group.Go(func() error {
		files, err := g.ListFiles(ctx, r)
		if err == nil {
			filesByPath := map[string]ReviewFile{}
			for _, f := range files {
				filesByPath[*f.Filename] = ReviewFile{
					IsAnnotated: false,
					File:        f,
				}
			}

			model.Files = filesByPath
		}
		return err
	})

	err := group.Wait()
	if err != nil {
		return nil, err
	}

	sort.Slice(model.Comments, func(i, j int) bool {
		path := *model.Comments[i].Path
		if file, exists := model.Files[path]; exists {
			model.Files[path] = ReviewFile{
				File:        file.File,
				IsAnnotated: true,
			}
		}

		c1, c1Ok := orderOf(*model.Comments[i].Body)
		c2, c2Ok := orderOf(*model.Comments[j].Body)
		if !c1Ok && !c2Ok {
			return false
		} else if !c1Ok {
			return false
		} else if !c2Ok {
			return true
		}
		return c1 < c2
	})

	return model, nil
}
