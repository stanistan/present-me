package github

import (
	"regexp"
	"strconv"

	"github.com/google/go-github/v52/github"
)

type ReviewModel struct {
	Params *ReviewParams `json:"params"`

	PR       *github.PullRequest          `json:"pr"`
	Review   *github.PullRequestReview    `json:"review"`
	Comments []*github.PullRequestComment `json:"comments"`
	Files    map[string]ReviewFile        `json:"files"`
}

type ReviewFile struct {
	IsAnnotated bool               `json:"isAnnotated"`
	File        *github.CommitFile `json:"file"`
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
