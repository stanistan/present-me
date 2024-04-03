package github

import (
	"context"
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
)

type ListSourcesAPISource struct {
	ReviewParamsMap ReviewParamsMap
}

var _ api.Source = &ListSourcesAPISource{}

func (s *ListSourcesAPISource) GetReview(ctx context.Context) (api.Review, error) {
	params, err := ReviewParamsFromMap(s.ReviewParamsMap)
	if err != nil {
		return api.Review{}, errors.WithStack(err)
	}

	model, err := FetchReviewModel(
		ctx,
		params,
		nil,
		func(body string) (int, bool) { return 0, false },
	)
	if err != nil {
		return api.Review{}, errors.WithStack(err)
	}

	var body string
	if model.PR.Body != nil && *model.PR.Body != "" {
		body = *model.PR.Body
	}

	return api.Review{
		Title: api.MaybeLinked{
			Text: fmt.Sprintf("%s (#%d)", *model.PR.Title, *model.PR.Number),
			HRef: *model.PR.Links.HTML.HRef,
		},
		Links: []api.LabelledLink{
			{
				Label: "Author",
				MaybeLinked: api.MaybeLinked{
					Text: *model.PR.User.Login,
					HRef: *model.PR.User.HTMLURL,
				},
			},
			{
				Label: "Pull Request",
				MaybeLinked: api.MaybeLinked{
					Text: fmt.Sprintf(
						"%s/%s/pull/%d",
						params.Owner,
						params.Repo,
						params.Pull,
					),
					HRef: *model.PR.HTMLURL,
				},
			},
		},
		MetaData: map[string]any{
			"params":     s.ReviewParamsMap,
			"reviewTags": extractReviewTags(model.Comments),
		},
		Body:     body,
		Comments: nil,
	}, nil
}

func extractReviewTags(cs []*PullRequestComment) map[string]int {
	tags := map[string]int{}
	for _, c := range cs {
		if c.Body == nil {
			continue
		}

		t, ok := parseReviewTag(*c.Body)
		if !ok {
			continue
		}

		count, _ := tags[t.Review]
		tags[t.Review] = count + 1
	}

	return tags
}
