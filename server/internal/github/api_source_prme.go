package github

import (
	"context"
	"fmt"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
)

// CommentsAPISource can construct a api.Review from a tag
// on a pull requests' comments.
type CommentsAPISource struct {
	ReviewParamsMap ReviewParamsMap
}

var _ api.Source = &CommentsAPISource{}

func (s *CommentsAPISource) GetReview(ctx context.Context) (api.Review, error) {
	params, err := ReviewParamsFromMap(s.ReviewParamsMap)
	if err != nil {
		return api.Review{}, errors.WithStack(err)
	}

	model, err := FetchReviewModel(
		ctx,
		params,
		CommentMatchesTag(s.ReviewParamsMap.Tag),
		func(body string) (int, bool) {
			t, ok := parseReviewTag(body)
			return t.Order, ok
		},
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
			"params": s.ReviewParamsMap,
		},
		Body:     body,
		Comments: transformComments(model.Comments),
	}, nil
}
