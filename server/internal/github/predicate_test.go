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
		expected   reviewTag
	}

	for _, data := range []testData{
		{
			"standard", "foo (prme1)", true, reviewTag{"1", 0},
		},
		{
			"no-prefix", "(prme2)", true, reviewTag{"2", 0},
		},
		{
			"trailing space", "(prme2)  ", true, reviewTag{"2", 0},
		},
		{
			"with order", "end of text (prmesomething-1)\n", true, reviewTag{"something", 1},
		},
		{
			"with order 3", "(prmesomething-3)", true, reviewTag{"something", 3},
		},
		{
			"doesnt match in the middle", "(prme1) banana", false, reviewTag{},
		},
		{
			"can have naked prme", "foo (prme)", true, reviewTag{"", 0},
		},
		{
			"won't parse one without parentheses", "foo prme", false, reviewTag{},
		},
		{
			"no tag but order", "anything (prme-5)", true, reviewTag{"", 5},
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
