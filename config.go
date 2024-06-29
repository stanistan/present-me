package presentme

import (
	"context"
	"time"

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
	Debug       bool                 `env:"DEBUG"`
	Environment string               `env:"ENV" default:"dev" enum:"dev,prod"`
	Log         log.Config           `embed:"" prefix:"log-"`
	DiskCache   cache.CacheOpts      `embed:"" prefix:"disk-cache-"`
	Github      github.ClientOptions `embed:"" prefix:"gh-"`
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

type ServeConfig struct {
	// Port describes the port this server runs on.
	Port               string        `default:"8080" env:"PORT"`
	Hostname           string        `default:"localhost" env:"HOSTNAME"`
	ServerReadTimeout  time.Duration `default:"5s"`
	ServerWriteTimeout time.Duration `default:"10s"`
}

func (c *ServeConfig) Address() string {
	return c.Hostname + ":" + c.Port
}
