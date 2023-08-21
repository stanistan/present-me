package diff

import (
	"strings"

	"github.com/stanistan/present-me/internal/errors"
)

const (
	implicitNumberOfLines = 3
)

type Scanner struct {
	r          HunkRange
	start, end int
	auto       bool
}

func NewScanner(r HunkRange, start, end RangeFrom) (*Scanner, error) {
	endLine, ok := end.extract()
	if !ok {
		return nil, errors.New("cannot parse endLne from range")
	}

	auto := false
	startLine, ok := start.extract()
	if !ok {
		startLine = endLine - implicitNumberOfLines
		auto = true
	}

	return &Scanner{
		r:     r,
		start: startLine,
		end:   endLine,
		auto:  auto,
	}, nil
}

func (s *Scanner) shouldCount(line string) bool {
	return !strings.HasPrefix(line, s.r.IgnorePrefix)
}

func (s *Scanner) isChunkUseful(c *chunk) bool {
	return c.useful && c.prefix != s.r.IgnorePrefix
}

type chunk struct {
	prefix string
	lines  []string
	useful bool
}

func (c *chunk) pushLine(l string) {
	lineIsUseful := strings.TrimSpace(l[1:]) != ""
	c.useful = c.useful || lineIsUseful
	c.lines = append(c.lines, l)
}

func (s *Scanner) Filter(ls []string) []string {
	var (
		c         = chunk{}
		cs        []chunk
		pushChunk = func() {
			if len(c.lines) > 0 {
				cs = append(cs, c)
				c = chunk{}
			}
		}
	)

	lineNo := s.r.StartingAt
	for idx, line := range ls {
		if idx == 0 {
			// we skip the first line
			// it is hunk metadata
			continue
		}

		if len(line) == 0 {
			// ensure line is normalized for diffing,
			// this only happens in testing because of trimmed whitespace
			// and we need to absolutely have a prefix for
			// empty context lines
			line = " "
		}

		prefix := line[0:1]

		// we go through the entire diff hunk, making sure
		// to track that we're in the desired range.
		//
		// and while we're in here we parse out chunks
		// based on the kind of context that they are, `+`, `-`, ` `
		// based on the prefix changing
		if lineNo >= s.start && lineNo <= s.end {
			if prefix != c.prefix {
				pushChunk()
			}
			c.pushLine(line)
		}

		// always track the last prefix for this chunk
		c.prefix = prefix

		if s.shouldCount(line) {
			lineNo++
		}
	}

	pushChunk()

	var (
		out       []string
		numChunks = len(cs)
		chunkIdx  = numChunks - 1
	)

	for chunkIdx >= 0 {
		c := cs[chunkIdx]
		if s.auto && len(out) >= implicitNumberOfLines && !s.isChunkUseful(&c) {
			break
		}

		out = append(c.lines, out...)
		chunkIdx--
	}

	return out
}
