package presentme

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	dc "github.com/stanistan/present-me/internal/cache"
)

type Config struct {
	DiskCache dc.CacheOpts `yaml:"diskcache"`
	Github    GHOpts       `yaml:"github"`
}

func MustConfig(path string) Config {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}

	log.Printf("config %+v", c)
	configureCache(c.DiskCache)
	return c
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache    *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
	cacheTTL           = 10 * time.Minute
)
