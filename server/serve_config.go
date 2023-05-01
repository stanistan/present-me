package presentme

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

type ServeConfig struct {
	// Port describes the port this server runs on.
	Port               string        `default:"8080" env:"PORT"`
	Hostname           string        `default:"localhost" env:"HOSTNAME"`
	ServerReadTimeout  time.Duration `default:"5s"`
	ServerWriteTimeout time.Duration `default:"10s"`

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

type staticFS struct {
	s           http.FileSystem
	defaultPath string
}

func spaFS(dir string) http.FileSystem {
	return &staticFS{
		s:           http.Dir(dir),
		defaultPath: "index.html",
	}
}

func (fs *staticFS) Open(name string) (http.File, error) {
	f, err := fs.s.Open(name)
	if os.IsNotExist(err) {
		return fs.s.Open(fs.defaultPath)
	}

	return f, err
}

func (c *ServeConfig) WebsiteHandler() (http.Handler, error) {
	if !c.IsProxy() {
		return http.FileServer(spaFS(c.StaticDir)), nil
	}

	remote, err := url.Parse(c.ProxyAddress)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ProxyAddress")
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	return proxy, nil
}
