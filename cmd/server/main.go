package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	pm "github.com/stanistan/present-me"
)

func main() {
	r := mux.NewRouter()

	sub := r.PathPrefix("/{owner}/{repo}/pull/{number}/{reviewID}").
		Methods("GET").
		Subrouter()

	sub.HandleFunc("", doMD(pm.AsMarkdownOptions{AsHTML: true, InBody: true}))
	sub.HandleFunc("/md", doMD(pm.AsMarkdownOptions{}))

	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("starting server at port %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func doMD(opts pm.AsMarkdownOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, func() error {
			params, err := pm.ReviewParamsFromMap(mux.Vars(r))
			if err != nil {
				return err
			}

			model, err := params.Model(
				pm.Context{
					Ctx:    r.Context(),
					Client: pm.GithubClient(r.Context()),
				},
				r.URL.Query().Get("refresh") == "1",
			)
			if err != nil {
				return err
			}

			return model.AsMarkdown(w, opts)
		})
	}
}

func handle(w http.ResponseWriter, f func() error) {
	if err := f(); err != nil {
		log.Printf("Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
