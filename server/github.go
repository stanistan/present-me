package presentme

import (
	"context"
	"net/http"
	"sort"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v52/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	dc "github.com/stanistan/present-me/internal/cache"
)

type GHOpts struct {
	AppID          int64        `name:"app_id" env:"GH_APP_ID" required:""`
	InstallationID int64        `name:"installation_id" env:"GH_INSTALLATION_ID" required:""`
	PrivateKey     GHPrivateKey `embed:"" prefix:"pk-" required:""`
}

type GHPrivateKey struct {
	File string `name:"file" env:"GH_PK_FILE"`
}

func (o *GHOpts) HTTPClient() (*http.Client, error) {
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

type GH struct {
	c *github.Client
}

func NewGH(opts GHOpts) (*GH, error) {
	c, err := opts.HTTPClient()
	if err != nil {
		return nil, err
	}

	return &GH{c: github.NewClient(c)}, nil
}

func (g *GH) ListFiles(ctx context.Context, r *ReviewParams) ([]*github.CommitFile, error) {
	var fs []*github.CommitFile
	return fs, cache.Apply(ctx, &fs, dc.Provider{
		Key: []interface{}{r.Owner, r.Repo, r.Number, "files"},
		Fetch: func() (interface{}, error) {
			d, _, err := g.c.PullRequests.ListFiles(ctx, r.Owner, r.Repo, r.Number, nil)
			return d, WrapGithubErr(err, "call to ListFiles failed")
		},
	})
}

func (g *GH) GetPullRequest(ctx context.Context, r *ReviewParams) (*github.PullRequest, error) {
	var pr *github.PullRequest
	return pr, cache.Apply(ctx, &pr, dc.Provider{
		Key: []interface{}{r.Owner, r.Repo, r.Number, "pr"},
		Fetch: func() (interface{}, error) {
			pr, _, err := g.c.PullRequests.Get(ctx, r.Owner, r.Repo, r.Number)
			return pr, WrapGithubErr(err, "call to GetPullRequest failed")
		},
	})
}

func (g *GH) ListReviews(ctx context.Context, r *ReviewParams) ([]*github.PullRequestReview, error) {
	var reviews []*github.PullRequestReview
	return reviews, cache.Apply(ctx, &reviews, dc.Provider{
		Key: []interface{}{r.Owner, r.Repo, r.Number, "reviews"},
		Fetch: func() (interface{}, error) {
			reviews, _, err := g.c.PullRequests.ListReviews(ctx, r.Owner, r.Repo, r.Number, nil)
			return reviews, WrapGithubErr(err, "call to ListReviews failed")
		},
	})
}

func (g *GH) GetReview(ctx context.Context, r *ReviewParams) (*github.PullRequestReview, error) {
	var review *github.PullRequestReview
	return review, cache.Apply(ctx, &review, dc.Provider{
		Key: []interface{}{r, "review"},
		Fetch: func() (interface{}, error) {
			review, _, err := g.c.PullRequests.GetReview(ctx, r.Owner, r.Repo, r.Number, r.ReviewID)
			return review, WrapGithubErr(err, "call to GetReview failed")
		},
	})
}

func (g *GH) ListReviewComments(ctx context.Context, r *ReviewParams) ([]*github.PullRequestComment, error) {
	var cs []*github.PullRequestComment
	return cs, cache.Apply(ctx, &cs, dc.Provider{
		Key: []interface{}{r, "review-comments"},
		Fetch: func() (interface{}, error) {
			cs, _, err := g.c.PullRequests.ListReviewComments(ctx, r.Owner, r.Repo, r.Number, r.ReviewID, nil)
			return cs, WrapGithubErr(err, "call to ListReviewComments failed")
		},
	})
}

func (g *GH) FetchReviewModel(ctx context.Context, r *ReviewParams) (*ReviewModel, error) {
	pull, err := g.GetPullRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	review, err := g.GetReview(ctx, r)
	if err != nil {
		return nil, err
	}

	comments, err := g.ListReviewComments(ctx, r)
	if err != nil {
		return nil, err
	}

	annotatedFiles := map[string]struct{}{}
	sort.Slice(comments, func(i, j int) bool {

		annotatedFiles[*comments[i].Path] = struct{}{}
		annotatedFiles[*comments[j].Path] = struct{}{}

		c1, c1Ok := orderOf(*comments[i].Body)
		c2, c2Ok := orderOf(*comments[j].Body)
		if !c1Ok && !c2Ok {
			return false
		} else if !c1Ok {
			return false
		} else if !c2Ok {
			return true
		}
		return c1 < c2
	})

	files, err := g.ListFiles(ctx, r)
	if err != nil {
		return nil, err
	}

	filesByPath := map[string]ReviewFile{}
	for _, f := range files {
		_, ok := annotatedFiles[*f.Filename]
		filesByPath[*f.Filename] = ReviewFile{
			IsAnnotated: ok,
			File:        f,
		}
	}

	return &ReviewModel{
		Params:   r,
		PR:       pull,
		Review:   review,
		Comments: comments,
		Files:    filesByPath,
	}, nil
}
