package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	pm "github.com/stanistan/present-me"
	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
)

const name = "present-me"

var version = "development" // dynamically linked

func main() {
	var config pm.Config
	_ = kong.Parse(&config,
		kong.Name(name),
		kong.Vars{"version": version},
		kong.UsageOnError())

	// set up our logger
	ctx, log := config.Logger(context.Background())

	// get our cache setup for the server
	diskCache := config.Cache(ctx)

	// configure our github client
	gh, err := config.GithubClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure GH client")
	}

	r := mux.NewRouter()

	// 1. Register API routes & middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(
		hlog.NewHandler(log),
		hlog.URLHandler("url"),
		github.Middleware(gh),
		cache.Middleware(diskCache, func(r *http.Request) *cache.Options {
			return &cache.Options{
				TTL:          10 * time.Minute,
				ForceRefresh: r.URL.Query().Get("refresh") == "1",
			}
		}),
	)

	for _, r := range apiRoutes {
		api.Handle(r.Prefix, r.Handler).Methods(r.Method)
	}

	// 2. Register fallback website handler
	websiteHandler, err := config.WebsiteHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("could not build handler")
	}
	r.PathPrefix("/").
		Handler(websiteHandler)

	// 3. Init server
	s := &http.Server{
		Addr:         config.Address(),
		ReadTimeout:  config.ServerReadTimeout,
		WriteTimeout: config.ServerWriteTimeout,
		Handler:      r,
	}

	log.Info().Str("address", config.Address()).Msg("starting server")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
