package github

import (
	"fmt"
	"strings"

	"github.com/stanistan/present-me/internal/github/diff"
)

type diffScanner struct {
	hunkRange diff.HunkRange

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

		lineNo     = p.hunkRange.StartingAt
		lastPrefix = "BANANA" // sentinel

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
		if !strings.HasPrefix(line, p.hunkRange.IgnorePrefix) {
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
		if auto && len(out) >= 3 && !chunk.isUseful(p.hunkRange.IgnorePrefix) {
			break
		}

		out = append(chunk.lines, out...)
		chunkIdx--
	}

	return out
}

func diffRange(c *PullRequestComment) (int, int, bool, error) {
	// - endLine is the line that the comment is on or after,
	// - startLine is the beginning line that we'll include in our diff,
	//   and it looks like github defaults to 4 lines included if there is
	//   no `StartLine`.
	var endLine int
	if c.OriginalLine != nil {
		endLine = *c.OriginalLine
	} else if c.Line != nil {
		endLine = *c.Line
	} else {
		return 0, 0, false, fmt.Errorf("invalid nil line")
	}

	var startLine int
	var auto bool
	if c.OriginalStartLine != nil {
		startLine = *c.OriginalStartLine
	} else if c.StartLine != nil {
		startLine = *c.StartLine
	} else {
		startLine = endLine - 3
		auto = true
	}

	return endLine, startLine, auto, nil
}
