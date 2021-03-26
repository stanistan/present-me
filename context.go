package crap

import (
	"context"

	"github.com/google/go-github/github"
)

type Context struct {
	Ctx    context.Context
	Client *github.Client
}
