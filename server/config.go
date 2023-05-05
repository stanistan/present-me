package presentme

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	dc "github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/present-me/internal/log"
)

type Config struct {
	ServeConfig
	DiskCache dc.CacheOpts  `embed:"" prefix:"disk-cache-"`
	Github    github.GHOpts `embed:"" prefix:"gh-"`
}

func (c *Config) GH() (*github.GH, error) {
	g, err := github.NewGH(c.Github)
	return g, errors.WithStack(err)
}

func (c *Config) Logger() zerolog.Logger {
	return log.NewLogger()
}

func (c *Config) Configure() {
	configureCache(c.DiskCache)
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(context.TODO(), opts)
}

var (
	cache *dc.Cache = dc.NewCache(
		context.TODO(),
		dc.CacheOpts{Enabled: false},
	)
)
