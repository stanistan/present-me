package github

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ReviewParams struct {
	Owner    string `required:"" help:"owner or organization" json:"owner"`
	Repo     string `required:"" help:"repository name" json:"repo"`
	Pull     int    `required:"" help:"pull request number" json:"pull"`
	ReviewID int64  `required:"" help:"reviewID number" json:"review"`
}

// The following formats are supported:
// - https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
// - github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
// - stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
// - stanistan/invoice-proxy/pull/3
//
// If there is no pullrequestreview number provided, we will find the
// first one that belongs to the author of the PR itself.
func ReviewParamsFromURL(i string) (*ReviewParams, error) {
	// trim protocol and domain if they are there, we add them back
	// to normalize and support urls that might not have the protocol,
	// or just look like "stanistan/...."
	i = strings.TrimPrefix(i, "https://")
	i = strings.TrimPrefix(i, "github.com/")
	i = "https://github.com/" + i

	u, err := url.Parse(i)
	if err != nil {
		return nil, err
	}

	pieces := strings.Split(u.Path, "/")
	if len(pieces) != 5 {
		return nil, fmt.Errorf("invalid url path %s", u.Path)
	}

	return ReviewParamsFromMap(ReviewParamsMap{
		Owner:  pieces[1],
		Repo:   pieces[2],
		Pull:   pieces[4],
		Review: strings.TrimPrefix(u.Fragment, "pullrequestreview-"),
	})
}

// ReviewParamsMap is the string representation,
// struct model that corresponds to a map[string]string,
// but named...
type ReviewParamsMap struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Pull   string `json:"pull"`
	Review string `json:"review"`
}

func NewReviewParamsMap(values url.Values) ReviewParamsMap {
	return ReviewParamsMap{
		Owner:  values.Get("org"),
		Repo:   values.Get("repo"),
		Pull:   values.Get("pull"),
		Review: values.Get("review"),
	}
}

func ReviewParamsFromMap(m ReviewParamsMap) (*ReviewParams, error) {
	owner := m.Owner
	if owner == "" {
		return nil, errors.New("missing owner")
	}

	repo := m.Repo
	if repo == "" {
		return nil, errors.New("missing repo")
	}

	numberVal := m.Pull
	if numberVal == "" {
		return nil, errors.New("missing number")
	}
	number, err := strconv.ParseInt(numberVal, 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number")
	}

	var reviewID int64
	reviewIDVal := m.Review
	if reviewIDVal != "" {
		reviewID, err = strconv.ParseInt(reviewIDVal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid reviewID: %s", err)
		}
	}

	return &ReviewParams{
		Owner:    owner,
		Repo:     repo,
		Pull:     int(number),
		ReviewID: reviewID,
	}, nil
}

func (r *ReviewParams) EnsureReviewID(ctx context.Context, g *Client) (bool, error) {
	if r.ReviewID != 0 {
		return false, nil
	}

	pull, err := g.GetPullRequest(ctx, r)
	if err != nil {
		return false, err
	}

	reviews, err := g.ListReviews(ctx, r)
	if err != nil {
		return false, err
	}

	if len(reviews) == 0 {
		return false, fmt.Errorf("PR has no reviews")
	}

	for _, rev := range reviews {
		if *rev.User.Login == *pull.User.Login {
			r.ReviewID = *rev.ID
			return true, nil
		}
	}

	return false, fmt.Errorf("PR has no review from the PR author %s", *pull.User.Login)
}
