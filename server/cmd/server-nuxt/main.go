package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	pm "github.com/stanistan/present-me"
)

func main() {
	var config pm.Config
	_ = kong.Parse(&config)

	gh, err := config.GH()
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure GH client")
	}

	r := mux.NewRouter()

	// 1. Register API routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(githubMiddleware(gh))
	for _, r := range apiRoutes {
		api.
			Handle(r.Prefix, r.Handler).
			Methods(r.Method)
	}

	// 2. Register fallback website handler
	websiteHandler, err := config.WebsiteHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("could not build handler")
	}
	r.PathPrefix("/").Handler(websiteHandler)

	// 3. Init server
	s := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func githubMiddleware(g *pm.GH) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "gh", g)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
