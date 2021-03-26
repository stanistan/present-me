package crap

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

func test() {
	// We get our context set up so we can chat with the client, this is
	// something that's got to be pretty simple, maybe we do this with a server
	// at some point, but to start with running this as a command that takes
	// some args would be pretty neato, and it can dump out some markdown.
	//
	// If this runs as a server we will then want to do take this markdown,
	// render it, etc.
	//
	// We can split this out into different `cmd/` s.
	//
	// The general behavior is this:
	//
	// 1. We get the Review
	// 2. We get the ReviewComments
	// 3. We get the Files
	//
	// Determine the order of the markdown to print:
	//
	// 0. The PR Title
	//
	// 1. the Review.Body
	//
	// 2. The ReviewComments.Body in order that they appear, the ordering of the
	//    comments is by 1) if they start with a number!, and then by the rest.
	//    this should be fine since we're not going to have _too many_ comments, this
	//    is a thing a person would do on their own. Of course.
	//
	//    Each one will also have an associated file path. Keep track of this so we can know
	// 	  what the rest of the file paths are, at the end of the PR?
	//
	//	  Will probably want to play with the representation of these (with _changes_, and links
	//    to the original source, etc).
	//
	// 3. Stuff about the rest of the files?
}

var startsWithNumberRegexp = regexp.MustCompile(`^\s*(\d+)\.`)

func orderOf(c string) (int, bool) {
	m := startsWithNumberRegexp.FindStringSubmatch(c)
	if m == nil {
		return 0, false
	}

	n, _ := strconv.Atoi(m[1])
	return n, true
}

func Run(c Context, p *ReviewParams) (*bytes.Buffer, error) {
	var b bytes.Buffer

	// get the PR
	pull, err := p.GetPullRequest(c)
	if err != nil {
		return nil, err
	}

	// maybe write a pre-amble here?

	_, err = fmt.Fprintf(&b, "# (#%d) %s \n\n", *pull.Number, *pull.Title)
	if err != nil {
		return nil, err
	}

	review, err := p.GetReview(c)
	if err != nil {
		return nil, err
	}

	b.Write([]byte(*review.Body))
	b.Write([]byte("\n\n"))

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

	for _, c := range comments {
		fmt.Fprintf(&b, "### %s\n\n", *c.Path)
		b.Write([]byte("```diff\n"))
		b.Write([]byte(*c.DiffHunk))
		b.Write([]byte("\n```"))
		b.Write([]byte("\n\n"))
		b.Write([]byte(*c.Body))
		b.Write([]byte("\n\n"))
	}

	return &b, nil
}
