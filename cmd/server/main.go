package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pm "github.com/stanistan/present-me"
	cache "github.com/stanistan/present-me/internal/cache"
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

	// serving css & fonts n things
	r.PathPrefix("/static").Handler(http.FileServer(http.FS(pm.StaticContent)))

	// our main routes
	sub := r.PathPrefix("/{owner}/{repo}/pull/{number}/{reviewID}").
		Methods("GET").
		Subrouter()

	sub.HandleFunc("/slides", doMD(g, pm.AsMarkdownOptions{AsSlides: true})).
		Name("slides")
	sub.HandleFunc("/md", doMD(g, pm.AsMarkdownOptions{})).
		Name("md")
	sub.HandleFunc("/post", doMD(g, pm.AsMarkdownOptions{AsHTML: true, InBody: true})).
		Name("post")

	r.HandleFunc("/{owner}/{repo}/pull/{number}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, func() error {
			params, err := pm.ReviewParamsFromMap(mux.Vars(r))
			if err != nil {
				return err
			}

			err = params.EnsureReviewID(cacheContext(r), g)
			if err != nil {
				return err
			}

			toURL, err := urlForParams(params, "post")(sub)
			if err != nil {
				return err
			}
			http.Redirect(w, r, toURL, http.StatusTemporaryRedirect)
			return nil
		})
	}))

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			_ = pm.IndexPage(w, "", "")
			return
		}

		params, err := pm.ReviewParamsFromURL(url)
		if err != nil {
			_ = pm.IndexPage(w, url, err.Error())
			return
		}

		err = params.EnsureReviewID(r.Context(), g)
		if err != nil {
			_ = pm.IndexPage(w, url, err.Error())
			return
		}

		toURL, err := urlForParams(params, urlType(r))(sub)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, toURL, http.StatusTemporaryRedirect)
	}))

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

			model, err := params.Model(cacheContext(r), g)
			if err != nil {
				return err
			}

			return model.AsMarkdown(w, opts)
		})
	}
}

func handle(w http.ResponseWriter, f func() error) {
	if err := f(); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func cacheContext(r *http.Request) context.Context {
	return cache.ContextWithOptions(r.Context(), &cache.Options{
		TTL:          10 * time.Minute,
		ForceRefresh: r.URL.Query().Get("refresh") == "1",
	})
}

func urlType(r *http.Request) string {
	t := r.URL.Query().Get("to")
	switch t {
	case "post", "md", "slides":
		return t
	default:
		return "post"
	}
}

func urlForParams(params *pm.ReviewParams, t string) func(*mux.Router) (string, error) {
	return func(sub *mux.Router) (string, error) {
		u, err := sub.Get(t).URL(
			"owner", params.Owner,
			"repo", params.Repo,
			"number", strconv.Itoa(params.Number),
			"reviewID", strconv.FormatInt(params.ReviewID, 10),
		)
		if err != nil {
			return "", errors.Wrap(err, "could not construct url")
		}
		return u.String(), nil
	}
}
