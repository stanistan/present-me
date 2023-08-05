package github

import (
	"regexp"
	"strconv"
)

type Pred[T any] func(T) bool

type CommentPredicate = Pred[*PullRequestComment]

func CommentMatchesReview(reviewID int64) CommentPredicate {
	return func(c *PullRequestComment) bool {
		return c.PullRequestReviewID != nil && *c.PullRequestReviewID == reviewID
	}
}

var prmeTagRegexp = regexp.MustCompile(`\(prme([^-\s]+)?(-(\d+))?\)\s*$`)

type reviewTag struct {
	Review string
	Order  int
}

func parseReviewTag(s string) (reviewTag, bool) {
	m := prmeTagRegexp.FindStringSubmatch(s)
	if m == nil {
		return reviewTag{}, false
	}

	n, _ := strconv.Atoi(m[3])
	return reviewTag{Review: m[1], Order: n}, true
}

func CommentMatchesTag(tag string) CommentPredicate {
	return func(c *PullRequestComment) bool {
		if c.Body == nil {
			return false
		}

		rtag, ok := parseReviewTag(*c.Body)
		return ok && rtag.Review == tag
	}
}
