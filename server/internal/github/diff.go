package github

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type diffScanner struct {
	countFrom                 int
	countLinesNotStartingWith string

	start, end int
}

type diffChunk struct {
	prefix string
	lines  []string
}

func (d *diffChunk) isUseful(countLinesNotStartingWith string) bool {

	allEmpty := true
	for _, l := range d.lines {
		if strings.TrimSpace(l[1:]) != "" {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		return false
	}

	return d.prefix != countLinesNotStartingWith
}

func (p *diffScanner) filter(lines []string, auto bool) []string {
	var (
		chunk  []string
		chunks []diffChunk

		lineNo     = p.countFrom
		lastPrefix string

		pushChunk = func() {
			if len(chunk) > 0 {
				chunks = append(chunks, diffChunk{prefix: lastPrefix, lines: chunk})
				chunk = nil
			}
		}
	)

	for idx, line := range lines {
		// N.B. we skip the first one since it's where the metadata is
		if idx == 0 {
			continue
		}

		// sometimes in testing lines are fully trimmed out -- we assume
		// in this case that it's an "empty context line"
		if len(line) == 0 {
			line = " "
		}

		prefix := line[0:1]
		if lineNo >= p.start && lineNo <= p.end {
			if prefix != lastPrefix {
				pushChunk()
			}
			chunk = append(chunk, line)
		}

		// track if we're changing prefixes
		lastPrefix = prefix

		// track if we're moving forward to the desired place
		if !strings.HasPrefix(line, p.countLinesNotStartingWith) {
			lineNo++
		}
	}
	pushChunk()

	var (
		out       []string
		numChunks = len(chunks)
		chunkIdx  = numChunks - 1
	)

	for chunkIdx >= 0 {
		chunk := chunks[chunkIdx]
		if auto && len(out) >= 3 && !chunk.isUseful(p.countLinesNotStartingWith) {
			break
		}

		out = append(chunk.lines, out...)
		chunkIdx--
	}

	return out
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

type diffHunkMetadata struct {
	Left, Right int
}

func (m *diffHunkMetadata) countConfig(side string) (int, string, error) {
	// - using the meta, we can see what the first line of the hunk is.
	// - depending on if we're doing LEFT or RIGHT, it means we will
	//   count (or not) specific lines when deciding which ones to include
	//   or not.
	switch side {
	case "LEFT":
		return m.Left, "+", nil
	case "RIGHT":
		return m.Right, "-", nil
	default:
		return 0, "", fmt.Errorf("side should be one of LEFT/RIGHT got %s", side)
	}
}

func diffRange(c *PullRequestComment) (int, int, bool, error) {

	// - endLine is the line that the comment is on or after,
	// - startLine is the beginning line that we'll include in our diff,
	//   and it looks like github defaults to 4 lines included if there is
	//   no `StartLine`.
	var endLine int
	if c.Line != nil {
		endLine = *c.Line
	} else if c.OriginalLine != nil {
		endLine = *c.OriginalLine
	} else {
		return 0, 0, false, fmt.Errorf("invalid nil line")
	}

	var startLine int
	var auto bool
	if c.StartLine != nil {
		startLine = *c.StartLine
	} else if c.OriginalStartLine != nil {
		startLine = *c.OriginalStartLine
	} else {
		startLine = endLine - 3
		auto = true
	}

	return endLine, startLine, auto, nil
}
