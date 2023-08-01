// Package api holds the API interface for what present me actually renders.
//
// This should be able to abstract away the *sources* of the data and only
// refer to the data-model that is used present "cards" and "slides".
package api

import "context"

type MaybeLinked struct {
	Text string `json:"text"`
	HRef string `json:"href,omitempty"`
}

type LabelledLink struct {
	MaybeLinked
	Label string `json:"label"`
}

type CodeBlock struct {
	Content  string `json:"content"`
	Language string `json:"lang"`
	IsDiff   bool   `json:"diff,omitempty"`
}

type Comment struct {
	Number      int         `json:"number"`
	Title       MaybeLinked `json:"title"`
	Description string      `json:"description"`
	CodeBlock   CodeBlock   `json:"code"`
}

type Review struct {
	Title    MaybeLinked `json:"title"`
	Body     string      `json:"body"`
	Comments []Comment   `json:"comments"`

	Links []LabelledLink `json:"links,omitempty"`

	// MetaData is :shrug:
	// probably an author association, and other links
	// that should be showing up in a place on the page.
	//
	// This is mostly untyped since we're not sure
	// what it would need to be for non-gh related items
	// at the moment.
	MetaData map[string]any `json:"meta,omitempty"`
}

type Source interface {
	GetReview(context.Context) (Review, error)
}

var _ Source = &Review{}

func (r *Review) GetReview(_ context.Context) (Review, error) {
	return *r, nil
}

// possible implemenations
// 1. github pr review
// 2. github pr and prme-id
// 3. github pr and NO prme-id
// 4. local, straight json API
