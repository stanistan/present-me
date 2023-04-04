package presentme

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/go-github/v50/github"
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

func (r *ReviewModel) CommitFile(filepath string) (*github.CommitFile, error) {
	f, ok := r.Files[filepath]
	if !ok || f.File == nil {
		return nil, fmt.Errorf("Missing file for path %s", filepath)
	}

	return f.File, nil
}

func (r *ReviewModel) Title() string {
	return fmt.Sprintf("%s/%s/pull/%d", r.Params.Owner, r.Params.Repo, r.Params.Number)
}

type AsMarkdownOptions struct {
	AsSlides bool
	AsHTML   bool `help:"if true will render the mardkown to html"`
	InBody   bool `help:"if true will place the rendered HTML into a body/template"`
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

func stripLeadingNumber(s string) string {
	return startsWithNumberRegexp.ReplaceAllString(s, "")
}
