package presentme

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
)

type ReviewModel struct {
	Params *ReviewParams

	PR       *github.PullRequest
	Review   *github.PullRequestReview
	Comments []*github.PullRequestComment
	Files    map[string]ReviewFile
}

type ReviewFile struct {
	IsAnnotated bool
	File        *github.CommitFile
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

func (r *ReviewModel) AsMarkdown(w io.Writer, opts AsMarkdownOptions) error {
	var buf bytes.Buffer

	log.Info().Msgf("rendering %+v", *r.Params)
	if err := reviewBody(&buf, r); err != nil {
		return err
	}

	if opts.AsSlides {
		return asSlide(w, r.Title(), buf.Bytes())
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

	return intoTemplate(w, r.Title(), html.Bytes())
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
