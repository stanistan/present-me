package github

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/stanistan/present-me/internal/errors"
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

func generateDiff(comment *PullRequestComment) (string, error) {
	// 1. we extract the metadata, we know which side we are going to be starting on.
	meta, err := diffHunkPrefix(*comment.DiffHunk)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// 2. how are we counting lines?
	countFrom, countLinesNotStartingWith, err := meta.countConfig(*comment.Side)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// 3. what is our range?
	endLine, startLine, auto, err := diffRange(comment)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// 4. configure out scanner
	scanner := &diffScanner{
		countFrom:                 countFrom,
		countLinesNotStartingWith: countLinesNotStartingWith,
		start:                     startLine,
		end:                       endLine,
	}

	// 5. filter our diff lines to only what's relevant for this comment
	out := scanner.filter(
		strings.Split(*comment.DiffHunk, "\n"),
		auto,
	)

	return strings.Join(out, "\n"), nil
}
