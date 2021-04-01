package presentme

import (
	"context"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

type GHOpts struct {
	AppID          int64  `yaml:"app_id"`
	InstallationID int64  `yaml:"installation_id"`
	PrivateKeyFile string `yaml:"private_key_file"`
}

func GHOptsFromFile(path string) (GHOpts, error) {
	var d GHOpts
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return d, err
	}

	err = yaml.Unmarshal(data, &d)
	return d, err
}

func (o *GHOpts) HTTPClient() (*http.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.NewKeyFromFile(tr, o.AppID, o.InstallationID, o.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: itr,
	}, nil
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
