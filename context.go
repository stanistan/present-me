package presentme

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Context struct {
	Ctx    context.Context
	Client *github.Client
}

func GithubClient(ctx context.Context) *github.Client {
	accessToken, _ := os.LookupEnv("PRESENT_ME_ACCESS_TOKEN")
	if accessToken != "" {
		tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: accessToken,
			},
		))
		return github.NewClient(tc)
	}

	return github.NewClient(nil)
}
