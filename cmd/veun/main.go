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

var cli struct {
	pm.Config
	Version struct{} `cmd:"version"`
	Serve   struct{} `cmd:"serve"`
}

func main() {
	k := kong.Parse(&cli, cliOptions...)
	config := cli.Config

	// init our logger and add it to the context
	ctx, log := config.Logger(context.Background())

	// configure the app
	app, err := App(ctx, log, config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not initialize app")
	}

	s := app.HTTPServer()
	switch k.Command() {
	case "version":
		log.Info().Str("sha", version).Msg("version")
	case "serve":
		log.Info().
			Str("address", config.Address()).
			Bool("debug", config.Debug).
			Str("env", config.Environment).
			Msg("starting server")

		if err := s.ListenAndServe(); err != nil {
			log.Fatal().
				Err(err).
				Msg("server failed")
		}
	case "":
		log.Fatal().Msg("missing command")
	default:
		log.Fatal().Msg(k.Command())
	}

}
