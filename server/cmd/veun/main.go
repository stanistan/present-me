package main

import (
	"context"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	pm "github.com/stanistan/present-me"
)

const name = "present-me"

var version = "development" // dynamically linked

func main() {
	var config pm.Config
	_ = kong.Parse(
		&config, kong.Name(name),
		kong.Vars{"version": version},
		kong.UsageOnError(),
	)

	_, server, err := newServer(context.Background(), config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialize server")
	}

	s := &http.Server{
		Addr:         config.Address(),
		ReadTimeout:  config.ServerReadTimeout,
		WriteTimeout: config.ServerWriteTimeout,
		Handler:      server.Handler(),
	}

	server.log.Info().Str("address", config.Address()).Msg("starting server")
	if err := s.ListenAndServe(); err != nil {
		server.log.Fatal().Err(err).Msg("failed to start server")
	}
}
