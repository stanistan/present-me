package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	pm "github.com/stanistan/present-me"
	"github.com/stanistan/present-me/internal/cache"
	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el-exp"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/request"
)

func App(
	ctx context.Context, log zerolog.Logger, config pm.Config,
) (*app, error) {
	cache := config.Cache(ctx)
	gh, err := config.GithubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not cconfigure github client: %w", err)
	}

	return &app{
		config: config,
		log:    log,
		gh:     gh,
		cache:  cache,
	}, nil
}

type app struct {
	config pm.Config
	log    zerolog.Logger
	gh     *github.Client
	cache  *cache.Cache
}

// debug is a request handler function that prints the current request name.
func (s *app) debug(r *http.Request) (veun.AsView, http.Handler, error) {

	if !s.config.Debug {
		return nil, http.NotFoundHandler(), nil
	}

	pathValue := func(name string) [2]string {
		return [2]string{
			"r.PathValue(\"" + name + "\")",
			r.PathValue(name),
		}
	}

	namedData := [][2]string{
		{"r.URL.Path", r.URL.Path},
		{"r.URL.RawQuery", r.URL.RawQuery},
		pathValue("owner"),
		pathValue("repo"),
		pathValue("pull"),
		pathValue("source"),
		pathValue("kind"),
	}

	var code = func(in string) el.Code {
		return el.Code{el.Text(in)}
	}

	var rows []veun.AsView
	for _, v := range namedData {
		rows = append(rows, el.Tr{
			el.Td{code(v[0])},
			el.Td{el.Em{code(v[1])}},
		})
	}

	return el.Div{
		el.H1{el.Text("present-me (debug)")},
		el.Table{
			el.THead{el.Tr{
				el.Th{code("var r *http.Request")},
				el.Th{code("value")},
			}},
			el.TBody{el.Content(rows)},
		},
	}, nil, nil
}

func (s *app) Handler() http.Handler {
	var (
		h  = vhttp.Handler
		hf = vhttp.HandlerFunc
	)
	_ = cache.Middleware(s.cache, func(r *http.Request) *cache.Options {
		return &cache.Options{
			TTL:          10 * time.Minute,
			ForceRefresh: r.URL.Query().Get("refresh") == "1",
		}
	})

	mux := http.NewServeMux()
	mux.Handle("GET /{owner}/{repo}/pull/{pull}/{source}/{kind}", hf(s.debug)) // list review sources
	mux.Handle("GET /{owner}/{repo}/pull/{pr}/", hf(s.debug))                  // list review sources
	mux.Handle("GET /{owner}/{repo}/pull/", hf(s.debug))                       // list pulls
	mux.Handle("GET /{owner}/{repo}/", hf(s.debug))                            // list pulls
	mux.Handle("GET /{owner}/", hf(s.debug))                                   // list repos _should we drop this_
	mux.Handle("GET /", hf(s.debug))                                           // search
	mux.Handle("GET /version", h(request.Always(veun.Raw(version))))           // what version are we running!?

	mux.Handle("/", hf(s.debug)) // 404
	return mux
}

func (s *app) HTTPServer() *http.Server {
	return &http.Server{
		Addr:         s.config.Address(),
		ReadTimeout:  s.config.ServerReadTimeout,
		WriteTimeout: s.config.ServerWriteTimeout,
		Handler:      s.Handler(),
	}
}
