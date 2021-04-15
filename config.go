package presentme

import (
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
	configureLogger()

	log.Info().Msgf("config %+v", c)
	configureCache(c.DiskCache)
}

func configureLogger() {
	// This is sset to be GOOGLE format (ish)
	// - https://cloud.google.com/logging/docs/structured-logging
	// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "times"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "message"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case zerolog.TraceLevel:
			return "DEBUG"
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARNING"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "ALERT"
		case zerolog.PanicLevel:
			return "ALERT"
		case zerolog.NoLevel:
			return ""
		}
		return ""
	}
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
)
