package presentme

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/google/go-github/github"
)

type ReviewModel struct {
	Params *ReviewParams

	PR       *github.PullRequest
	Review   *github.PullRequestReview
	Comments []*github.PullRequestComment
	Files    []*github.CommitFile
}

type AsMarkdownOptions struct {
	AsHTML bool `help:"if true will render the mardkown to html"`
	InBody bool `help:"if true will place the rendered HTML into a body/template"`
}

func (r *ReviewModel) AsMarkdown(w io.Writer, opts AsMarkdownOptions) error {
	if r == nil {
		return fmt.Errorf("model is nil!")
	}

	var (
		buf bytes.Buffer
	)

	log.Printf("rendering %+v", *r.Params)
	if err := reviewBody(&buf, r); err != nil {
		return err
	}

	if !opts.AsHTML {
		_, err := w.Write(buf.Bytes())
		return err
	}

	var html bytes.Buffer
	if err := md.Convert(buf.Bytes(), &html); err != nil {
		return err
	}

	if !opts.InBody {
		_, err := w.Write(html.Bytes())
		return err
	}

	return intoTemplate(w, html.Bytes())
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
