package github

import "github.com/stanistan/present-me/internal/api"

func transformComments(cs []*PullRequestComment) []api.Comment {
	out := make([]api.Comment, len(cs))
	for idx, c := range cs {
		c := c
		out[idx] = api.Comment{
			Number: idx + 1,
			Title: api.MaybeLinked{
				Text: *c.Path,
				HRef: *c.HTMLURL,
			},
			Description: commentBody(*c.Body),
			CodeBlock: api.CodeBlock{
				IsDiff:   true,
				Content:  *c.DiffHunk,
				Language: detectLanguage(*c.Path),
			},
		}
	}
	return out
}

func commentBody(s string) string {
	out := startsWithNumberRegexp.ReplaceAllString(s, "")
	return prmeTagRegexp.ReplaceAllString(out, "")
}
