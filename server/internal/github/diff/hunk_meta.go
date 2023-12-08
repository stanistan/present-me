package diff

import (
	"fmt"
	"regexp"
	"strconv"
)

var hunkPrefixRegexp = regexp.MustCompile(`^@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)

type HunkRange struct {
	StartingAt, NumLines int
	IgnorePrefix         string
}

type HunkMeta struct {
	Original, New HunkRange
}

func ParseHunkMeta(hunk string) (HunkMeta, error) {
	m := hunkPrefixRegexp.FindStringSubmatch(hunk)
	if m == nil {
		return HunkMeta{}, fmt.Errorf("could not parse hunk prefix")
	}

	l1, _ := strconv.Atoi(m[1])
	l2, _ := strconv.Atoi(m[2])
	r1, _ := strconv.Atoi(m[3])
	r2, _ := strconv.Atoi(m[4])

	return HunkMeta{
		Original: HunkRange{StartingAt: l1, NumLines: l2, IgnorePrefix: "+"},
		New:      HunkRange{StartingAt: r1, NumLines: r2, IgnorePrefix: "-"},
	}, nil
}

func (m *HunkMeta) RangeForSide(side string) (HunkRange, error) {
	// - using the meta, we can see what the first line of the hunk is.
	// - depending on if we're doing LEFT or RIGHT, it means we will
	//   count (or not) specific lines when deciding which ones to include
	//   or not.
	switch side {
	case "LEFT":
		return m.Original, nil
	case "RIGHT":
		return m.New, nil
	default:
		return HunkRange{}, fmt.Errorf("side should be one of LEFT/RIGHT got %s", side)
	}
}
