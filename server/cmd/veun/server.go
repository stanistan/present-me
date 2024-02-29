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

func newServer(ctx context.Context, config pm.Config) (context.Context, *server, error) {
	ctx, log := config.Logger(ctx)
	cache := config.Cache(ctx)
	gh, err := config.GithubClient(ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("could not cconfigure github client: %w", err)
	}

	return ctx, &server{
		log:   log,
		gh:    gh,
		cache: cache,
	}, nil
}

type server struct {
	log   zerolog.Logger
	gh    *github.Client
	cache *cache.Cache
}

// debug is a request handler function that prints the current request name.
func (s *server) debug(r *http.Request) (veun.AsView, http.Handler, error) {

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
		el.H1{el.Text("Debug")},
		el.Table{
			//			el.Caption{code("var r *http.Request")},
			el.THead{el.Tr{
				el.Th{code("var r *http.Request")},
				el.Th{code("value")},
			}},
			el.TBody{el.Content(rows)},
		},
	}, nil, nil
}

func (s *server) Handler() http.Handler {
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

var (
	h  = vhttp.Handler
	hf = vhttp.HandlerFunc
)
