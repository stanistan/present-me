package presentme

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/present-me/internal/log"
)

type Config struct {
	ServeConfig
	Log       log.Config           `embed:"" prefix:"log-"`
	DiskCache cache.CacheOpts      `embed:"" prefix:"disk-cache-"`
	Github    github.ClientOptions `embed:"" prefix:"gh-"`
}

func (c *Config) GithubClient(ctx context.Context) (*github.Client, error) {
	g, err := github.New(ctx, c.Github)
	return g, errors.WithStack(err)
}

func (c *Config) Logger() zerolog.Logger {
	return log.NewLogger(c.Log)
}

func (c *Config) Cache(ctx context.Context) *cache.Cache {
	return cache.NewCache(ctx, c.DiskCache)
}
