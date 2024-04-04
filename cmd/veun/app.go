package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	pm "github.com/stanistan/present-me"
	"github.com/stanistan/present-me/internal/api"
	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/present-me/internal/view/home"
	"github.com/stanistan/present-me/internal/view/layout"
	"github.com/stanistan/present-me/internal/view/review"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"
)

func App(
	ctx context.Context,
	log zerolog.Logger,
	config pm.Config,
) (*app, error) {

	gh, err := config.GithubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not cconfigure github client: %w", err)
	}

	return &app{
		config: config,
		log:    log,
		gh:     gh,
		cache:  config.Cache(ctx),
	}, nil
}

type app struct {
	config pm.Config
	log    zerolog.Logger
	gh     *github.Client
	cache  *cache.Cache
}

func (s *app) debugView(r *http.Request) veun.AsView {
	pathValue := func(name string) []string {
		return []string{"r.PathValue(\"" + name + "\")", r.PathValue(name)}
	}

	namedData := [][]string{
		{"r.URL.Path", r.URL.Path},
		{"r.URL.RawQuery", r.URL.RawQuery},
		pathValue("owner"),
		pathValue("repo"),
		pathValue("pull"),
		pathValue("source"),
		pathValue("kind"),
	}

	return el.Div{
		el.Class("p-3", "bg-red-100", "h-full", "font-mono"),
		el.H1{
			el.Class("text-2xl", "p-4", "font-bold", "text-center"),
			el.Text("present-me (debug)"),
		},
		table(
			[]string{"var r *http.Request", "value"},
			namedData,
			el.Class("mx-auto", "w-[75%]"),
			el.Caption{
				el.Class("caption-bottom", "text-xs"),
				el.Text("debug http-request things!"),
			},
		),
	}
}

// debug is a request handler function that prints the current request name.
func (s *app) debug(r *http.Request) (veun.AsView, http.Handler, error) {
	if !s.config.Debug {
		return nil, http.NotFoundHandler(), nil
	}

	return s.layout(s.debugView(r), nil), nil, nil
}

func (s *app) layout(view veun.AsView, d func() veun.AsView) veun.AsView {

	if d != nil && s.config.Debug {
		view = veun.Views{view, d()}
	}

	cssFile := "/static/styles.dev.css"
	if s.config.Environment == "prod" {
		cssFile = "/static/styles.css"
	}

	s.log.Trace().Msg("returning layout")
	return layout.Layout(layout.Params{
		Title:    "present-me",
		CSSFiles: []string{cssFile},
		JSFiles:  []string{"/static/prism.js"},
		Version: layout.Version{
			URL: "https://github.com/stanistan/present-me/" + version,
			SHA: version[0:7],
		},
	}, view)
}

func (s *app) Handler() http.Handler {
	var (
		gHandler = github.Middleware(s.gh)
		cHandler = cache.Middleware(s.cache, func(r *http.Request) *cache.Options {
			return &cache.Options{
				TTL:          10 * time.Minute,
				ForceRefresh: r.URL.Query().Get("refresh") == "1",
			}
		})
		logHandler = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(s.log.WithContext(r.Context())))
			})
		}

		h = func(r request.Handler) http.Handler {
			return logHandler(gHandler(cHandler(vhttp.Handler(r))))
		}

		hf = func(r request.HandlerFunc) http.Handler {
			return logHandler(gHandler(cHandler(vhttp.HandlerFunc(r))))
		}
	)

	mux := http.NewServeMux()
	mux.Handle("GET /static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.Handle("GET /{owner}/{repo}/pull/{pull}/{source}/{kind}", hf(s.review)) // do the source and kind
	mux.Handle("GET /{owner}/{repo}/pull/{pull}/{source}", hf(s.review))        // do the source thing
	mux.Handle("GET /{owner}/{repo}/pull/{pull}/", hf(s.pr))                    // list review sources
	mux.Handle("GET /{owner}/{repo}/pull/", hf(s.debug))                        // list pulls
	mux.Handle("GET /{owner}/{repo}/", hf(s.debug))                             // list pulls
	mux.Handle("GET /{owner}/", hf(s.debug))                                    // list repos _should we drop this_
	mux.Handle("GET /", handler.OnlyRoot(hf(s.home)))
	mux.Handle("GET /version", h(request.Always(veun.Raw(version)))) // what version are we running!?
	mux.Handle("/", handler.OnlyRoot(h(request.Always(home.Home{}))))

	return mux
}

func (s *app) listSources(ctx context.Context, params github.ReviewParamsMap) (veun.AsView, http.Handler, error) {
	source := &github.ListSourcesAPISource{ReviewParamsMap: params}
	model, sources, err := source.Sources(ctx)
	if err != nil {
		return nil, nil, err
	}

	return s.layout(
		review.SourcesList(params, model, sources),
		nil,
	), nil, nil
}

func (s *app) pr(r *http.Request) (veun.AsView, http.Handler, error) {
	params := github.ReviewParamsMap{
		Owner: r.PathValue("owner"),
		Repo:  r.PathValue("repo"),
		Pull:  r.PathValue("pull"),
		Kind:  "cards",
	}
	return s.listSources(r.Context(), params)
}

func (s *app) home(r *http.Request) (veun.AsView, http.Handler, error) {
	query := r.URL.Query()
	h := home.Home{
		Owner: query.Get("owner"),
		Repo:  query.Get("repo"),
		Pull:  query.Get("pull"),
	}

	if h.Owner != "" && h.Repo != "" && h.Pull != "" {
		params := github.ReviewParamsMap{
			Owner: h.Owner, Repo: h.Repo, Pull: h.Pull, Kind: "cards",
		}
		source := &github.ListSourcesAPISource{ReviewParamsMap: params}
		model, sources, err := source.Sources(r.Context())
		if err != nil {
			h.SearchResults = el.Div{
				el.Class("mx-auto", "bg-pink-800", "text-white", "text-center", "text-xs", "p-3"),
				el.Text(err.Error()),
			}
		} else {
			sources := review.SourcesFragment(params, model, sources)
			h.SearchResults = el.Div{
				el.Class(
					"border-t-2 border-t-pink-200 bg-gray-50 pb-3",
					//"border-b border-b-gray-200",
				),
				el.Div{
					el.Class("max-w-2xl mx-auto"),
					sources,
				},
			}
		}
	}

	return s.layout(h, nil), nil, nil
}

func (s *app) review(r *http.Request) (veun.AsView, http.Handler, error) {
	s.log.Debug().Msg("rendering review")

	pathSource := r.PathValue("source")
	_, reviewID, hasReviewID := strings.Cut(pathSource, "review-")
	_, tag, hasTagID := strings.Cut(pathSource, "tag")

	params := github.ReviewParamsMap{
		Owner:  r.PathValue("owner"),
		Repo:   r.PathValue("repo"),
		Pull:   r.PathValue("pull"),
		Review: reviewID,
		Tag:    tag,
		Kind:   r.PathValue("kind"),
	}

	if params.Kind == "" {
		params.Kind = "slides"
	}

	var source api.Source
	if hasReviewID {
		source = &github.ReviewAPISource{ReviewParamsMap: params}
	} else if hasTagID {
		source = &github.CommentsAPISource{ReviewParamsMap: params}
	}

	model, err := source.GetReview(r.Context())
	if err != nil {
		s.log.Trace().Msg("failed fetching reviw")
		return nil, nil, err
	}

	var content veun.AsView
	switch params.Kind {
	case "cards":
		content = review.PageContent(params, model)
	case "slides":
		content = review.SlideContent(params, model)
	}

	s.log.Trace().Msg("rendering review layout")
	return s.layout(
		content,
		func() veun.AsView { return s.debugView(r) },
	), nil, nil
}

func (s *app) HTTPServer() *http.Server {
	return &http.Server{
		Addr:         s.config.Address(),
		ReadTimeout:  s.config.ServerReadTimeout,
		WriteTimeout: s.config.ServerWriteTimeout,
		Handler:      s.Handler(),
	}
}

func table(heading []string, rows [][]string, ps ...el.Param) el.Table {
	var (
		cell = el.Class("border", "border-white", "p-2", "text-xs")
	)
	return el.Table{
		el.Class("table-auto", "text-sm"),
		el.THead{
			el.Tr{
				el.Th{},
				el.MapFragment(heading, func(title string, _ int) el.Th {
					return el.Th{cell, el.Class("bg-green-400"), el.Text(title)}
				}),
			},
		},
		el.TBody{
			el.MapFragment(rows, func(row []string, idx int) el.Tr {
				return el.Tr{
					el.Td{
						cell,
						el.Class("bg-pink-100", "text-right"),
						el.Text(fmt.Sprintf("%d", idx+1)),
					},
					el.MapFragment(row, func(t string, _ int) el.Td {
						return el.Td{
							cell,
							el.Class("bg-pink-200"),
							el.Text(t),
						}
					}),
				}
			}),
		},
		el.Fragment(ps),
	}
}
