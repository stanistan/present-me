package presentme

import (
	"time"

	log "github.com/sirupsen/logrus"

	dc "github.com/stanistan/present-me/internal/cache"
)

type Config struct {
	Port string `env:"PORT" default:"8080"`

	DiskCache dc.CacheOpts `embed:"" prefix:"disk-cache-"`
	Github    GHOpts       `embed:"" prefix:"gh-"`
}

func (c *Config) Configure() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Infof("config %+v", c)
	configureCache(c.DiskCache)
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache    *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
	cacheTTL           = 10 * time.Minute
)
