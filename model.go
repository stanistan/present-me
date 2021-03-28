package crap

import (
	"bytes"
	// "encoding/json"
	"io"
	"log"
	"regexp"
	"sort"
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

func BuildReviewModel(c Context, p *ReviewParams, refreshData bool) (*ReviewModel, error) {
	if !refreshData {
		cached, err := cache.Read(p)
		if err != nil {
			return nil, err
		}
		if cached != nil {
			return cached, nil
		}
	}

	pull, err := p.GetPullRequest(c)
	if err != nil {
		return nil, err
	}

	review, err := p.GetReview(c)
	if err != nil {
		return nil, err
	}

	comments, err := p.ListReviewComments(c)
	if err != nil {
		return nil, err
	}

	sort.Slice(comments, func(i, j int) bool {
		c1, c1Ok := orderOf(*comments[i].Body)
		c2, c2Ok := orderOf(*comments[j].Body)
		if !c1Ok && !c2Ok {
			return false
		} else if !c1Ok {
			return false
		} else if !c2Ok {
			return true
		}
		return c1 < c2
	})

	files, err := p.ListFiles(c)
	if err != nil {
		return nil, err
	}

	return cache.Write(p, &ReviewModel{
		Params:   p,
		PR:       pull,
		Review:   review,
		Comments: comments,
		Files:    files,
	})
}

type AsMarkdownOptions struct {
	AsHTML bool `help:"if true will render the mardkown to html"`
	InBody bool `help:"if true will place the rendered HTML into a body/template"`
}

func (r *ReviewModel) AsMarkdown(w io.Writer, opts AsMarkdownOptions) error {
	var (
		buf bytes.Buffer
		_   = func(ss ...string) {
			for _, s := range ss {
				buf.Write([]byte(s))
			}
			buf.Write([]byte("\n\n"))
		}
	)

	log.Printf("rendering %+v", *r.Params)
	if err := reviewBody(&buf, r); err != nil {
		return err
	}

	/*
		write("# (#", strconv.Itoa(r.Params.Number), ") ", *r.PR.Title)
		write(*r.Review.Body)

		for _, c := range r.Comments {
			write("### " + *c.Path)
			write("```diff\n" + *c.DiffHunk + "\n```")
			write(stripLeadingNumber(*c.Body))
		}

		buf.Write([]byte("```json\n"))
		e := json.NewEncoder(&buf)
		e.SetIndent("", "  ")
		_ = e.Encode(r)
		write("```")
	*/

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
