package github

import (
	"bufio"
	"fmt"
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
	//
	// 1. we extract the metadata, so we know which side we are going to be starting
	//    on.
	meta, err := diffHunkPrefix(*comment.DiffHunk)
	if err != nil {
		return "", errors.WithStack(err)
	}

	//
	// 2.
	//
	// - using the meta, we can see what the first line of the hunk is.
	// - depending on if we're doing LEFT or RIGHT, it means we will
	//   count (or not) specific lines when deciding which ones to include
	//   or not.
	var countFrom int
	var countLinesNotStartingWith string
	switch *comment.Side {
	case "LEFT":
		countFrom = meta.Left
		countLinesNotStartingWith = "+"
	case "RIGHT":
		countFrom = meta.Right
		countLinesNotStartingWith = "-"
	}

	//
	// 3.
	//
	// - endLine is the line that the comment is on or after,
	// - startLine is the beginning line that we'll include in our diff,
	//   and it looks like github defaults to 4 lines included if there is
	//   no `StartLine`.
	var endLine = *comment.Line
	var startLine int
	if comment.StartLine == nil {
		startLine = endLine - 3
	} else {
		startLine = *comment.StartLine
	}

	//
	// 4.
	//
	// Ok, so we're now ready to start going through the _input_ diff
	// line by line.
	var idx = countFrom
	var lines = bufio.NewScanner(strings.NewReader(*comment.DiffHunk))
	var out strings.Builder

	// drop the first line, since it has the metadata in it (we parsed it above).
	lines.Scan()

	// actually go through line by line, keeping the diff range we want
	// from this hunk.
	for lines.Scan() {
		line := lines.Text()
		if idx >= startLine && idx <= endLine {
			out.WriteString(line + "\n")
		}

		if !strings.HasPrefix(line, countLinesNotStartingWith) {
			idx++
		}
	}

	// and join our string!
	return out.String(), nil
}

type diffHunkMetadata struct {
	Left, Right int
}

var diffHunkPrefixRegexp = regexp.MustCompile(`^@@ -(\d+),\d+ \+(\d+),\d+ @@`)

func diffHunkPrefix(hunk string) (diffHunkMetadata, error) {
	meta := diffHunkMetadata{}
	m := diffHunkPrefixRegexp.FindStringSubmatch(hunk)
	if m == nil {
		return meta, fmt.Errorf("could not parse hunk prefix")
	}

	meta.Left, _ = strconv.Atoi(m[1])
	meta.Right, _ = strconv.Atoi(m[2])
	return meta, nil
}
