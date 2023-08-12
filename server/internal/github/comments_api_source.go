package github

import (
	"context"
	"fmt"
	"net/url"

	"github.com/stanistan/present-me/internal/api"
)

// CommentsAPISource can construct a api.Review from a tag
// on a pull requests' comments.
type CommentsAPISource struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Pull  string `json:"pull"`
	Tag   string `json:"tag"`
}

func NewCommentsAPISourceFromValues(values url.Values) api.Source {
	return &CommentsAPISource{
		Owner: values.Get("org"),
		Repo:  values.Get("repo"),
		Pull:  values.Get("pull"),
		Tag:   values.Get("tag"),
	}
}

var _ api.Source = &CommentsAPISource{}

func (s *CommentsAPISource) GetReview(ctx context.Context) (api.Review, error) {
	gh, ok := Ctx(ctx)
	if !ok || gh == nil {
		return api.Review{}, fmt.Errorf("missing gh context")
	}

	return api.Review{}, nil
}
