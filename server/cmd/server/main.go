package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	pm "github.com/stanistan/present-me"
	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
)

var (
	version = "development" // dynamically linked
)

func main() {
	var config pm.Config
	_ = kong.Parse(
		&config,
		kong.Name("present-me"),
		kong.UsageOnError(),
		kong.Description(fmt.Sprintf("build version: %s", version)),
	)

	// 0. Standard Deps
	// - logger,
	// - ctx withLogger
	// - disk cache
	// - github client
	log := config.Logger()
	ctx := log.WithContext(context.Background())
	diskCache := config.Cache(ctx)
	gh, err := config.GithubClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure GH client")
	}

	r := mux.NewRouter()

	// 1. Register API routes & middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(
		hlog.NewHandler(log),
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
