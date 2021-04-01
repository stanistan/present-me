package presentme

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	dc "github.com/stanistan/present-me/internal/cache"
)

type ReviewParams struct {
	Owner    string `required:"" help:"owner or organization"`
	Repo     string `required:"" help:"repository name"`
	Number   int    `required:"" help:"pull request number"`
	ReviewID int64  `required:"" help:"reviewID number"`
}

//
// The following formats are supported:
// https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
// github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
// stanistan/invoice-proxy/pull/3#pullrequestreview-605888708
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

func (r *ReviewParams) Model(ctx context.Context, g *GH, refreshData bool) (*ReviewModel, error) {
	var data *ReviewModel
	err := cache.Apply(&data, dc.Provider{
		Key:          r,
		TTL:          cacheTTL,
		ForceRefresh: refreshData,
		Fetch: func() (interface{}, error) {
			return g.FetchReviewModel(ctx, r)
		},
	})
	return data, err
}

var (
	cache    = dc.NewCache()
	cacheTTL = 10 * time.Minute
)
