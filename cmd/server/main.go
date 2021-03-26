package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	"github.com/stanistan/crap"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{owner}/{repo}/pull/{number}/{reviewID}", renderReview).
		Methods("GET").
		Name("review")

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func run(w http.ResponseWriter, r *http.Request, params *crap.ReviewParams) error {
	model, err := params.Model(crap.Context{
		Ctx:    r.Context(),
		Client: github.NewClient(nil),
	})
	if err != nil {
		return err
	}

	return model.AsMarkdown(w, crap.AsMarkdownOptions{
		AsHTML: true,
		InBody: true,
	})
}

func renderReview(w http.ResponseWriter, r *http.Request) {
	handle(w, func() error {
		params, err := crap.ReviewParamsFromMap(mux.Vars(r))
		if err != nil {
			return err
		}

		return run(w, r, params)
	})
}

func handle(w http.ResponseWriter, f func() error) {
	if err := f(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
