// Package api holds the API interface for what present me actually renders.
//
// This should be able to abstract away the *sources* of the data and only
// refer to the data-model that is used present "cards" and "slides".
package api

import "context"

// Source is anything that can give us a `Review` representation given a context.
type Source interface {
	GetReview(context.Context) (Review, error)
}

// Review represents anything that we can present, it is explicitly agnostic
// of data provider.
//
// This is not yet an exhaustive API description as when we add features, like
// "files" this should have them.
type Review struct {
	// Title is the title of the review, as will be presented.
	Title MaybeLinked `json:"title"`

	// Body is a raw markdown description of the Review.
	Body string `json:"body"`

	// Comments are individual code block descriptions.
	Comments []Comment `json:"comments"`

	// Links are metadata links (with label, text, and href).
	Links []LabelledLink `json:"links,omitempty"`

	// MetaData is anything else that we might want to share.
	MetaData map[string]any `json:"meta,omitempty"`
}

var _ Source = &Review{}

// GetReview for a Review will always return itself, and fullfil
// the Source interface.
func (r *Review) GetReview(_ context.Context) (Review, error) {
	return *r, nil
}

// Comment (basically a standard card), is a named and annotated code block.
type Comment struct {
	Number      int         `json:"number"`
	Title       MaybeLinked `json:"title"`
	Description string      `json:"description"`
	CodeBlock   CodeBlock   `json:"code"`
}

// CodeBlock is the code block and its diff.
type CodeBlock struct {
	Content  string `json:"content"`
	Language string `json:"lang"`
	IsDiff   bool   `json:"diff,omitempty"`
}

// MaybeLinked refers to text that can have a url associated with it.
type MaybeLinked struct {
	Text string `json:"text"`
	HRef string `json:"href,omitempty"`
}

// LabelledLink is a text link with a label.
type LabelledLink struct {
	MaybeLinked
	Label string `json:"label"`
}
