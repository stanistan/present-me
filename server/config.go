package presentme

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/present-me/internal/log"
)

// Config is the main configuration wrapper struct for
// the present-me application.
//
// It is basically the DI provider for all of the runtime
// dependencies and how to configure them.
type Config struct {
	ServeConfig
	Debug     bool                 `env:"DEBUG"`
	Log       log.Config           `embed:"" prefix:"log-"`
	DiskCache cache.CacheOpts      `embed:"" prefix:"disk-cache-"`
	Github    github.ClientOptions `embed:"" prefix:"gh-"`
}

func (c *Config) GithubClient(ctx context.Context) (*github.Client, error) {
	g, err := github.New(ctx, c.Github)
	return g, errors.WithStack(err)
}

func (c *Config) Logger(ctx context.Context) (context.Context, zerolog.Logger) {
	l := log.NewLogger(c.Log)
	return l.WithContext(ctx), l
}

func (c *Config) Cache(ctx context.Context) *cache.Cache {
	return cache.NewCache(ctx, c.DiskCache)
}
