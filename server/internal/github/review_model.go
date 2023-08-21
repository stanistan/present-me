package github

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/stanistan/present-me/internal/errors"
	"github.com/stanistan/present-me/internal/github/diff"
)

type ReviewModel struct {
	Params *ReviewParams `json:"params"`

	PR       *PullRequest          `json:"pr"`
	Review   *PullRequestReview    `json:"review"`
	Comments []*PullRequestComment `json:"comments"`
	Files    map[string]ReviewFile `json:"files"`
}

type ReviewFile struct {
	IsAnnotated bool        `json:"isAnnotated"`
	File        *CommitFile `json:"file"`
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

func generateDiff(c *PullRequestComment) (string, error) {
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
