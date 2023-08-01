package github

import (
	"context"
	"fmt"
	"path"
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

	gh, ok := Ctx(ctx)
	if !ok || gh == nil {
		return api.Review{}, fmt.Errorf("need gh context")
	}

	model, err := params.Model(ctx, gh)
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

func transformComments(cs []*PullRequestComment) []api.Comment {
	out := make([]api.Comment, len(cs))
	for idx, c := range cs {
		c := c
		out[idx] = api.Comment{
			Number: idx + 1,
			Title: api.MaybeLinked{
				Text: *c.Path,
				HRef: *c.HTMLURL,
			},
			Description: *c.Body,
			CodeBlock: api.CodeBlock{
				IsDiff:   true,
				Content:  *c.DiffHunk,
				Language: detectLanguage(*c.Path),
			},
		}
	}
	return out
}

func detectLanguage(p string) string {
	var (
		base = path.Base(p)
		ext  = path.Ext(p)
	)

	switch ext {
	case "":
		if base == "Dockerfile" {
			return "docker"
		}
		return "bash"
	case ".rs":
		return "rust"
	case ".vue":
		return "html"
	case ".Dockerfile":
		return "docker"
	}

	return ext[1:]
}
