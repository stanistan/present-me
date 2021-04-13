package presentme

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ReviewParams struct {
	Owner    string `required:"" help:"owner or organization"`
	Repo     string `required:"" help:"repository name"`
	Number   int    `required:"" help:"pull request number"`
	ReviewID int64  `required:"" help:"reviewID number"`
}

//
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

	return ReviewParamsFromMap(map[string]string{
		"owner":    pieces[1],
		"repo":     pieces[2],
		"number":   pieces[4],
		"reviewID": strings.TrimPrefix(u.Fragment, "pullrequestreview-"),
	})
}

func ReviewParamsFromMap(m map[string]string) (*ReviewParams, error) {
	owner, ok := m["owner"]
	if !ok || owner == "" {
		return nil, errors.New("missing owner")
	}

	repo, ok := m["repo"]
	if !ok || repo == "" {
		return nil, errors.New("missing repo")
	}

	numberVal, ok := m["number"]
	if !ok || numberVal == "" {
		return nil, errors.New("missing number")
	}
	number, err := strconv.ParseInt(numberVal, 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number")
	}

	var reviewID int64
	reviewIDVal, ok := m["reviewID"]
	if ok && reviewIDVal != "" {
		reviewID, err = strconv.ParseInt(reviewIDVal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid reviewID: %s", err)
		}
	}

	return &ReviewParams{
		Owner:    owner,
		Repo:     repo,
		Number:   int(number),
		ReviewID: reviewID,
	}, nil
}

func (r *ReviewParams) Model(ctx context.Context, g *GH) (*ReviewModel, error) {
	return g.FetchReviewModel(ctx, r)
}

func (r *ReviewParams) EnsureReviewID(ctx context.Context, g *GH) error {
	if r.ReviewID != 0 {
		return nil
	}

	pull, err := g.GetPullRequest(ctx, r)
	if err != nil {
		return errors.Wrap(err, "could not fetch PR")
	}

	reviews, err := g.ListReviews(ctx, r)
	if err != nil {
		return errors.Wrap(err, "could not fetch reviews for PR")
	}
	for _, rev := range reviews {
		if *rev.User.Login == *pull.User.Login {
			r.ReviewID = *rev.ID
			return nil
		}
	}

	return errors.New("could not find review by author of PR")
}
