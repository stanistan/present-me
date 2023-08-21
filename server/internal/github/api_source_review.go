package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
)

type ReviewAPISource struct {
	ReviewParamsMap ReviewParamsMap
}

var _ api.Source = &ReviewAPISource{}

func (s *ReviewAPISource) GetReview(ctx context.Context) (api.Review, error) {
	params, err := ReviewParamsFromMap(s.ReviewParamsMap)
	if err != nil {
		return api.Review{}, errors.WithStack(err)
	}

	if params.ReviewID == 0 {
		return api.Review{}, errors.New("review must have review id")
	}

	model, err := FetchReviewModel(ctx, params, CommentMatchesReview(params.ReviewID), orderOf)
	if err != nil {
		return api.Review{}, errors.WithStack(err)
	}

	var body []string
	if model.PR.Body != nil && *model.PR.Body != "" {
		body = append(body, *model.PR.Body)
	}

	if model.Review.Body != nil && *model.Review.Body != "" {
		body = append(body, *model.Review.Body)
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
			{
				Label: "Review",
				MaybeLinked: api.MaybeLinked{
					Text: fmt.Sprintf("#review-%d", params.ReviewID),
					HRef: *model.Review.HTMLURL,
				},
			},
		},
		MetaData: map[string]any{
			"params": s.ReviewParamsMap,
		},
		Body:     strings.Join(body, "\n\n---\n\n"),
		Comments: transformComments(model.Comments),
	}, nil
}
