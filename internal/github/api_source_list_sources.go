package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
	"golang.org/x/sync/errgroup"
)

type ListSourcesAPISource struct {
	ReviewParamsMap ReviewParamsMap
}

func (s *ListSourcesAPISource) Sources(ctx context.Context) (api.Review, []api.SourceProvider, error) {
	var r api.Review
	l := log.Ctx(ctx)
	l.Info().Msg("sources")

	params, err := ReviewParamsFromMap(s.ReviewParamsMap)
	if err != nil {
		return r, nil, errors.WithStack(err)
	}

	client, ok := Ctx(ctx)
	if !ok || client == nil {
		return r, nil, errors.New("missing github client")
	}

	var (
		pr       *PullRequest
		reviews  []*PullRequestReview
		comments []*PullRequestComment
	)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		pull, err := client.GetPullRequest(ctx, params)
		if err == nil {
			pr = pull
		}
		return err
	})

	group.Go(func() error {
		rs, err := client.ListReviews(ctx, params)
		if err == nil {
			reviews = rs
		}
		return err
	})

	group.Go(func() error {
		cs, err := client.ListComments(ctx, params, nil)
		if err == nil {
			comments = cs
		}
		return err
	})

	if err := group.Wait(); err != nil {
		e := errors.WrapErr(err)
		log.Ctx(ctx).Error().Int("err_code", e.HttpCode).Msg("failed sources")
		switch e.HttpCode {
		case 403, 404:
			return r, nil, nil
		default:
			return r, nil, err
		}
	}

	l.Info().Msg("did the sources")

	r.Title = api.MaybeLinked{
		Text: fmt.Sprintf("%s (#%d)", *pr.Title, *pr.Number),
		HRef: *pr.Links.HTML.HRef,
	}

	r.MetaData = map[string]any{"params": params}

	if pr.Body != nil && *pr.Body != "" {
		r.Body = *pr.Body
	}

	var sources []api.SourceProvider

	tags, _ := extractReviewTags(comments)
	for _, tag := range tags {
		sources = append(sources, &CommentsAPISource{
			ReviewParamsMap: s.ReviewParamsMap.WithTag(tag),
		})
	}

	commentsByReviewID := map[int64]int{}
	for _, c := range comments {
		if c.PullRequestReviewID == nil {
			continue
		}
		id := *c.PullRequestReviewID
		count := commentsByReviewID[id]
		commentsByReviewID[id] = count + 1
	}

	for _, r := range reviews {
		if count, ok := commentsByReviewID[*r.ID]; ok && count > 0 {
			if r.Body == nil || len(*r.Body) == 0 {
				continue
			}
			sources = append(sources, &ReviewAPISource{
				ReviewParamsMap: s.ReviewParamsMap.WithReview(
					strconv.FormatInt(*r.ID, 10),
				),
			})
		}
	}

	return r, sources, nil
}

func extractReviewTags(
	cs []*PullRequestComment,
) ([]string, map[string]int) {
	var (
		counter = map[string]int{}
		tags    []string
	)

	for _, c := range cs {
		if c.Body == nil {
			continue
		}

		t, ok := parseReviewTag(*c.Body)
		if !ok {
			continue
		}

		count, present := counter[t.Review]
		counter[t.Review] = count + 1
		if !present {
			tags = append(tags, t.Review)
		}
	}

	return tags, counter
}
