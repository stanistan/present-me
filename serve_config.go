package presentme

import (
	"time"
)

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
