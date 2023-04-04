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
	"github.com/rs/zerolog/log"

	pm "github.com/stanistan/present-me"
	cache "github.com/stanistan/present-me/internal/cache"
)

func main() {
	var config pm.Config
	_ = kong.Parse(&config)
	config.Configure()

	g, err := config.GH()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	r := mux.NewRouter()

	// TODO(https://github.com/stanistan/present-me/issues/9):
	//
	// This can be something that we pull in at runtime instead of building
	// the assets into the binary.
	r.PathPrefix("/static").Handler(http.FileServer(http.FS(pm.StaticContent)))

	// our main routes router, sub routes are set up later
	sub := r.PathPrefix("/{owner}/{repo}/pull/{number}/{reviewID}").
		Methods("GET").
		Subrouter()

	// our server context,
	// it holds the subrouter so we can do consistent url redirection
	server := &server{g: g, sub: sub}

	// all the routes!
	sub.HandleFunc("/slides", server.slides()).Name("slides")                 // MD rendered as slides
	sub.HandleFunc("/md", server.rawMD()).Name("md")                          // raw MD
	sub.HandleFunc("/post", server.post()).Name("post")                       // the MD rendered as a post
	r.HandleFunc("/{owner}/{repo}/pull/{number}", server.pr()).Methods("GET") // redirects to one of the above
	r.Handle("/", server.index()).Methods("GET")                              // renders the form
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, func() error {
			return &pm.Error{
				Msg:      "Not Found",
				HttpCode: 404,
			}
		})
	})

	// start the server
	err = listenAndServe(r)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func listenAndServe(router *mux.Router) error {
	// get our PORT from the env
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info().Str("address", s.Addr).Msg("started server")
	return s.ListenAndServe()
}

type server struct {
	g   *pm.GH
	sub *mux.Router
}

func (s *server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, func() error {
			var (
				url         string
				err         error
				renderIndex = func() error {
					var errMsg string
					if err != nil {
						errMsg = err.Error()
					}
					_ = pm.IndexPage(w, url, errMsg)
					return nil
				}
			)

			url = r.URL.Query().Get("url")
			if url == "" {
				return renderIndex()
			}

			var params *pm.ReviewParams
			params, err = pm.ReviewParamsFromURL(url)
			if err != nil {
				return renderIndex()
			}

			_, err = params.EnsureReviewID(cacheContext(r), s.g)
			if err != nil {
				return renderIndex()
			}

			toURL, err := s.urlForParams(params, urlType(r))
			if err != nil {
				return &pm.Error{
					Msg:      "could not construct valid url",
					Cause:    err,
					HttpCode: 500,
				}
			}

			http.Redirect(w, r, toURL, http.StatusTemporaryRedirect)
			return nil
		})
	}
}

func (s *server) pr() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, func() error {
			params, err := reviewParamsFromVars(mux.Vars(r))
			if err != nil {
				return &pm.Error{
					Msg:      "Malformed PR URL",
					Cause:    err,
					HttpCode: 400,
				}
			}

			_, err = params.EnsureReviewID(cacheContext(r), s.g)
			if err != nil {
				return &pm.Error{
					Msg:      "Failed to find associated review ID",
					Cause:    err,
					HttpCode: 404,
				}
			}

			toURL, err := s.urlForParams(params, "post")
			if err != nil {
				return err
			}

			http.Redirect(w, r, toURL, http.StatusTemporaryRedirect)
			return nil
		})
	}
}

var (
	slidesOpts pm.AsMarkdownOptions = pm.AsMarkdownOptions{AsSlides: true}
	rawMDOpts                       = pm.AsMarkdownOptions{}
	postOpts                        = pm.AsMarkdownOptions{AsHTML: true, InBody: true}
)

func (s *server) slides() http.HandlerFunc {
	return s.doMD(slidesOpts)
}

func (s *server) rawMD() http.HandlerFunc {
	return s.doMD(rawMDOpts)
}

func (s *server) post() http.HandlerFunc {
	return s.doMD(postOpts)
}

func (s *server) doMD(opts pm.AsMarkdownOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, func() error {
			params, err := reviewParamsFromVars(mux.Vars(r))
			if err != nil {
				return err
			}

			changed, err := params.EnsureReviewID(cacheContext(r), s.g)
			if err != nil {
				return &pm.Error{
					Msg:      "missing reviewID",
					Cause:    err,
					HttpCode: 404,
				}
			} else if changed {
				toURL, err := s.urlForParams(params, optsToURLType(opts))
				if err != nil {
					return &pm.Error{
						Msg:      "could not construct valid url",
						Cause:    err,
						HttpCode: 500,
					}
				}
				http.Redirect(w, r, toURL, http.StatusTemporaryRedirect)
				return nil
			}

			model, err := params.Model(cacheContext(r), s.g)
			if err != nil {
				return err
			}

			return model.AsMarkdown(w, opts)
		})
	}
}

func handle(w http.ResponseWriter, r *http.Request, f func() error) {
	if err := f(); err != nil {
		e := pm.WrapErr(err)
		if e.HttpCode > 500 {
			log.Err(e).Int("status", e.HttpCode).Msg("")
		} else {
			log.Info().
				Int("status", e.HttpCode).
				Str("path", r.URL.Path).
				Msgf("error: %s", e.Error())
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(e.HttpCode)
		_ = pm.ErrorPage(w, e)
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

func optsToURLType(opts pm.AsMarkdownOptions) string {
	switch opts {
	case slidesOpts:
		return "slides"
	case rawMDOpts:
		return "md"
	default:
		return "post"
	}
}

func (s *server) urlForParams(params *pm.ReviewParams, t string) (string, error) {
	u, err := s.sub.Get(t).URL(
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

func reviewParamsFromVars(m map[string]string) (*pm.ReviewParams, error) {
	return pm.ReviewParamsFromMap(pm.ReviewParamsMap{
		Owner:  m["owner"],
		Repo:   m["repo"],
		Number: m["number"],
		Review: m["reviewID"],
	})
}
