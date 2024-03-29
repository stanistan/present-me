package github

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCommentTagRegex(t *testing.T) {
	type testData struct {
		name       string
		input      string
		expectedOk bool
		expected   ReviewTag
	}

	for _, data := range []testData{
		{
			"standard", "foo (prme1)", true, ReviewTag{"1", 0},
		},
		{
			"no-prefix", "(prme2)", true, ReviewTag{"2", 0},
		},
		{
			"trailing space", "(prme2)  ", true, ReviewTag{"2", 0},
		},
		{
			"with order", "end of text (prmesomething-1)\n", true, ReviewTag{"something", 1},
		},
		{
			"with order 3", "(prmesomething-3)", true, ReviewTag{"something", 3},
		},
		{
			"doesnt match in the middle", "(prme1) banana", false, ReviewTag{},
		},
		{
			"can have naked prme", "foo (prme)", true, ReviewTag{"", 0},
		},
		{
			"won't parse one without parentheses", "foo prme", false, ReviewTag{},
		},
		{
			"no tag but order", "anything (prme-5)", true, ReviewTag{"", 5},
		},
	} {
		d := data
		t.Run(d.name, func(t *testing.T) {
			tag, ok := parseReviewTag(d.input)
			assert.Equal(t, d.expectedOk, ok)
			assert.Equal(t, d.expected, tag)
		})
	}
}
