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
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	websiteHandler, err := config.WebsiteHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("could not build handler")
	}

	r.PathPrefix("/").Handler(websiteHandler)

	s := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
