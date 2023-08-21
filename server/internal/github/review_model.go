package github

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
