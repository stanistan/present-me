package crap

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
)

type ReviewParams struct {
	Owner, Repo string
	Number      int
	ReviewID    int64
}

func ReviewParamsFromMap(m map[string]string) (*ReviewParams, error) {
	owner, ok := m["owner"]
	if !ok || owner == "" {
		return nil, fmt.Errorf("missing owner")
	}

	repo, ok := m["repo"]
	if !ok || repo == "" {
		return nil, fmt.Errorf("missing repo")
	}

	numberVal, ok := m["number"]
	if !ok || numberVal == "" {
		return nil, fmt.Errorf("missing number")
	}
	number, err := strconv.ParseInt(numberVal, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", err)
	}

	reviewIDVal, ok := m["reviewID"]
	if !ok || reviewIDVal == "" {
		return nil, fmt.Errorf("missing reviewID")
	}
	reviewID, err := strconv.ParseInt(reviewIDVal, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid reviewID: %s", err)
	}

	return &ReviewParams{
		Owner:    owner,
		Repo:     repo,
		Number:   int(number),
		ReviewID: reviewID,
	}, nil
}

func (r *ReviewParams) ListFiles(c Context) ([]*github.CommitFile, error) {
	fs, _, err := c.Client.PullRequests.ListFiles(c.Ctx, r.Owner, r.Repo, r.Number, nil)
	return fs, err
}

func (r *ReviewParams) GetPullRequest(c Context) (*github.PullRequest, error) {
	pr, _, err := c.Client.PullRequests.Get(c.Ctx, r.Owner, r.Repo, r.Number)
	return pr, err
}

func (r *ReviewParams) GetReview(c Context) (*github.PullRequestReview, error) {
	review, _, err := c.Client.PullRequests.GetReview(c.Ctx, r.Owner, r.Repo, r.Number, r.ReviewID)
	return review, err
}

func (r *ReviewParams) ListReviewComments(c Context) ([]*github.PullRequestComment, error) {
	cs, _, err := c.Client.PullRequests.ListReviewComments(c.Ctx, r.Owner, r.Repo, r.Number, r.ReviewID, nil)
	return cs, err
}

func (r *ReviewParams) Model(c Context) (*ReviewModel, error) {
	return BuildReviewModel(c, r)
}
