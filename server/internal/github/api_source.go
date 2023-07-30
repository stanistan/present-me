package github

import (
	"context"
	"fmt"
	"path"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
)

type ReviewSource struct {
	ReviewParamsMap ReviewParamsMap
}

var _ api.Source = &ReviewSource{}

func (s *ReviewSource) GetReview(ctx context.Context) (api.Review, error) {
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

	var body string
	if model.PR.Body != nil {
		body = body + *model.PR.Body
	}

	if model.Review.Body != nil {
		body = body + *model.Review.Body
	}

	return api.Review{
		Body: body,
		Title: api.MaybeLinked{
			Text: *model.PR.Title,
			HRef: *model.PR.Links.HTML.HRef,
		},
		Comments: transformComments(model.Comments),
		MetaData: map[string]any{
			"params": s.ReviewParamsMap,
		},
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
