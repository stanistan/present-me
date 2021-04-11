package presentme

import (
	"context"
	"net/http"
	"sort"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/stanistan/present-me/internal/secret"
)

type GHOpts struct {
	AppID          int64        `name:"app_id" env:"GH_APP_ID" required:""`
	InstallationID int64        `name:"installation_id" env:"GH_INSTALLATION_ID" required:""`
	PrivateKey     GHPrivateKey `embed:"" prefix:"pk-" required:""`
}

type GHPrivateKey struct {
	File       string `name:"file" env:"GH_PK_FILE"`
	SecretName string `name:"secret-name" env:"GH_PK_SECRET_NAME"`
}

func (o *GHOpts) HTTPClient() (*http.Client, error) {
	var (
		itr http.RoundTripper
		err error

		tr = http.DefaultTransport
	)

	if o.PrivateKey.File != "" {
		log.Infof("attempting to read PK from File")
		itr, err = ghinstallation.NewKeyFromFile(tr, o.AppID, o.InstallationID, o.PrivateKey.File)
	} else if o.PrivateKey.SecretName != "" {
		log.Info("attempting to read PK from secret")
		var pk []byte
		pk, err = secret.Get(context.Background(), o.PrivateKey.SecretName)
		if err == nil {
			itr, err = ghinstallation.New(tr, o.AppID, o.InstallationID, pk)
		}
	} else {
		itr = tr
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not create HTTP Client")
	}

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
	fs, _, err := g.c.PullRequests.ListFiles(ctx, r.Owner, r.Repo, r.Number, nil)
	return fs, err
}

func (g *GH) GetPullRequest(ctx context.Context, r *ReviewParams) (*github.PullRequest, error) {
	pr, _, err := g.c.PullRequests.Get(ctx, r.Owner, r.Repo, r.Number)
	return pr, err
}

func (g *GH) GetReview(ctx context.Context, r *ReviewParams) (*github.PullRequestReview, error) {
	review, _, err := g.c.PullRequests.GetReview(ctx, r.Owner, r.Repo, r.Number, r.ReviewID)
	return review, err
}

func (g *GH) ListReviewComments(ctx context.Context, r *ReviewParams) ([]*github.PullRequestComment, error) {
	cs, _, err := g.c.PullRequests.ListReviewComments(ctx, r.Owner, r.Repo, r.Number, r.ReviewID, nil)
	return cs, err
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

	sort.Slice(comments, func(i, j int) bool {
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

	return &ReviewModel{
		Params:   r,
		PR:       pull,
		Review:   review,
		Comments: comments,
		Files:    files,
	}, nil
}
