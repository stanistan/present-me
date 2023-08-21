package github

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/errors"
	"github.com/stanistan/present-me/internal/github/diff"
)

func transformComments(cs []*PullRequestComment) []api.Comment {
	out := make([]api.Comment, len(cs))
	for idx, c := range cs {
		c := c

		diff, err := commentCodeDiff(c)
		if err != nil {
			diff = *c.DiffHunk
		}

		out[idx] = api.Comment{
			Number: idx + 1,
			Title: api.MaybeLinked{
				Text: *c.Path,
				HRef: *c.HTMLURL,
			},
			Description: commentBody(*c.Body),
			CodeBlock: api.CodeBlock{
				IsDiff:   true,
				Content:  diff,
				Language: detectLanguage(*c.Path),
			},
		}
	}
	return out
}

func commentBody(s string) string {
	out := startsWithNumberRegexp.ReplaceAllString(s, "")
	return prmeTagRegexp.ReplaceAllString(out, "")
}

func commentCodeDiff(c *PullRequestComment) (string, error) {
	// we extract the metadata, we know which side we are going to be starting on.
	meta, err := diff.ParseHunkMeta(*c.DiffHunk)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// how are we counting lines?
	hunkRange, err := meta.RangeForSide(*c.Side)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// - endLine is the line that the comment is on or after,
	// - startLine is the beginning line that we'll include in our diff,
	//   and it looks like github defaults to 4 lines included if there is
	//   no `StartLine`.
	scanner, err := diff.NewScanner(
		hunkRange,
		diff.RangeFrom{c.OriginalStartLine, c.StartLine},
		diff.RangeFrom{c.OriginalLine, c.Line},
	)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// filter our diff lines to only what's relevant for this comment
	out := scanner.Filter(strings.Split(*c.DiffHunk, "\n"))
	return strings.Join(out, "\n"), nil
}

var startsWithNumberRegexp = regexp.MustCompile(`^\s*(\d+)\.\s*`)

func orderOf(c string) (int, bool) {
	m := startsWithNumberRegexp.FindStringSubmatch(c)
	if m == nil {
		return 0, false
	}

	n, _ := strconv.Atoi(m[1])
	return n, true
}
