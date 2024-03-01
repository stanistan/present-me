package main

import (
	"context"

	"github.com/alecthomas/kong"

	pm "github.com/stanistan/present-me"
)

var (
	name       = "present-me"
	version    = "development" // version is dynamically linked at build time.
	cliOptions = []kong.Option{
		kong.Name(name),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	}
)

func main() {
	var config pm.Config
	_ = kong.Parse(&config, cliOptions...)

	// init our logger and add it to the context
	ctx, log := config.Logger(context.Background())

	// configure the app
	app, err := App(ctx, log, config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not initialize app")
	}

	log.Info().
		Str("address", config.Address()).
		Bool("debug", config.Debug).
		Msg("starting server")

	s := app.HTTPServer()
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}
