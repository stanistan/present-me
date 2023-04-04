package presentme

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	dc "github.com/stanistan/present-me/internal/cache"
)

type ServeConfig struct {
	// Port describes the port this server runs on.
	Port     string `default:"8080" env:"PORT"`
	Hostname string `default:"localhost" env:"HOSTNAME"`

	// Serve desides if we're running in proxy mode (for development)
	// or if we are going to be serving the content from the static directory
	// --- which is what happens when we've built our docker image.
	Serve        string `required:"" enum:"static,proxy" default:"static"`
	StaticDir    string `optional:"" default:"./static"`
	ProxyAddress string `optional:"" default:"http://localhost:3000"`
}

func (c *ServeConfig) IsProxy() bool {
	return c.Serve == "proxy"
}

func (c *ServeConfig) Address() string {
	return c.Hostname + ":" + c.Port
}

type spaFS struct {
	r http.FileSystem
}

func (fs *spaFS) Open(name string) (http.File, error) {
	f, err := fs.r.Open(name)
	if os.IsNotExist(err) {
		return fs.r.Open("index.html")
	}

	return f, err
}

func (c *ServeConfig) WebsiteHandler() (http.Handler, error) {
	if !c.IsProxy() {
		fs := &spaFS{http.Dir(c.StaticDir)}
		server := http.FileServer(fs)
		return server, nil
	}

	remote, err := url.Parse(c.ProxyAddress)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ProxyAddress")
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	return proxy, nil
}

type Config struct {
	ServeConfig
	DiskCache dc.CacheOpts `embed:"" prefix:"disk-cache-"`
	Github    GHOpts       `embed:"" prefix:"gh-"`
}

func (c *Config) GH() (*GH, error) {
	g, err := NewGH(c.Github)
	return g, errors.WithStack(err)
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
			return "CRITICAL"
		case zerolog.PanicLevel:
			return "ALERT"
		case zerolog.NoLevel:
			return "DEFAULT"
		}
		return "DEFAULT"
	}
}

func configureCache(opts dc.CacheOpts) {
	cache = dc.NewCache(opts)
}

var (
	cache *dc.Cache = dc.NewCache(dc.CacheOpts{Enabled: false})
)
