package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	pm "github.com/stanistan/present-me"
)

func main() {
	var config pm.Config
	_ = kong.Parse(&config)
	config.Configure()

	g, err := pm.NewGH(config.Github)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	sub := r.PathPrefix("/{owner}/{repo}/pull/{number}/{reviewID}").
		Methods("GET").
		Subrouter()

	sub.HandleFunc("/slides", doMD(g, pm.AsMarkdownOptions{AsSlides: true})).Name("slides")
	sub.HandleFunc("/md", doMD(g, pm.AsMarkdownOptions{})).Name("md")
	sub.HandleFunc("/post", doMD(g, pm.AsMarkdownOptions{AsHTML: true, InBody: true})).Name("post")

	r.PathPrefix("/static").Handler(http.FileServer(http.FS(pm.StaticContent)))

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			_ = pm.IndexPage(w, "", "")
			return
		}

		urlType := r.URL.Query().Get("to")
		params, err := pm.ReviewParamsFromURL(url)
		if err != nil {
			_ = pm.IndexPage(w, url, err.Error())
		} else {
			toURL, err := sub.Get(urlType).URL(
				"owner", params.Owner,
				"repo", params.Repo,
				"number", strconv.Itoa(params.Number),
				"reviewID", strconv.FormatInt(params.ReviewID, 10),
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Redirect(w, r, toURL.String(), http.StatusTemporaryRedirect)
			}
		}
	}))

	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("starting server at %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func doMD(g *pm.GH, opts pm.AsMarkdownOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, func() error {
			params, err := pm.ReviewParamsFromMap(mux.Vars(r))
			if err != nil {
				return err
			}

			model, err := params.Model(r.Context(), g, r.URL.Query().Get("refresh") == "1")
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
