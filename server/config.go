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
	DiskCache cache.CacheOpts `embed:"" prefix:"disk-cache-"`
	Github    github.GHOpts   `embed:"" prefix:"gh-"`
}

func (c *Config) GH() (*github.GH, error) {
	g, err := github.NewGH(c.Github)
	return g, errors.WithStack(err)
}

func (c *Config) Logger() zerolog.Logger {
	return log.NewLogger()
}

func (c *Config) Cache(ctx context.Context) *cache.Cache {
	return cache.NewCache(ctx, c.DiskCache)
}
