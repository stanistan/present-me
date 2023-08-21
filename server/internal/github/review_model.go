package github

import (
	"regexp"
	"strconv"
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
