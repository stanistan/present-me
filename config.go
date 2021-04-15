package presentme

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	dc "github.com/stanistan/present-me/internal/cache"
)

type Config struct {
	Port string `env:"PORT" default:"8080"`

	DiskCache dc.CacheOpts `embed:"" prefix:"disk-cache-"`
	Github    GHOpts       `embed:"" prefix:"gh-"`
}

func (c *Config) Configure() {
	// This is sset to be GOOGLE format (ish)
	// https://cloud.google.com/logging/docs/structured-logging
	configureLogger()

	log.Info().Msgf("config %+v", c)
	configureCache(c.DiskCache)
}

func configureLogger() {
	zerolog.MessageFieldName = "message"
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "times"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
)
