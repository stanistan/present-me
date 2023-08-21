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

type ReviewTag struct {
	Review string
	Order  int
}

var prmeTagRegexp = regexp.MustCompile(`\(prme([^-\s]+)?(-(\d+))?\)\s*$`)

func parseReviewTag(s string) (ReviewTag, bool) {
	m := prmeTagRegexp.FindStringSubmatch(s)
	if m == nil {
		return ReviewTag{}, false
	}

	n, _ := strconv.Atoi(m[3])
	return ReviewTag{Review: m[1], Order: n}, true
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
