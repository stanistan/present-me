package main

import (
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

	r := mux.NewRouter()

	// 1. Register API routes
	api := r.PathPrefix("/api").Subrouter()
	for _, r := range apiRoutes {
		api.Handle(r.Prefix, r.Handler).Methods(r.Method)
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
