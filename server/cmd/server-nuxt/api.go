package main

import (
	"net/http"

	"github.com/pkg/errors"

	pm "github.com/stanistan/present-me"
)

var apiRoutes []Route = []Route{
	{
		"GET", "/ping",
		func(_ *http.Request) (*JSONResponse, error) {
			return OKResponse(map[string]any{"ping": "pong"}), nil
		},
	},
	{
		"GET", "/review",
		func(r *http.Request) (*JSONResponse, error) {
			var (
				ctx    = r.Context()
				values = r.URL.Query()
			)

			gh, ok := ctx.Value("gh").(*pm.GH)
			if !ok || gh == nil {
				return nil, errors.New("missing github context")
			}

			params, err := pm.ReviewParamsFromMap(pm.ReviewParamsMap{
				Owner:  values.Get("org"),
				Repo:   values.Get("repo"),
				Number: values.Get("pull"),
				Review: values.Get("review"),
			})
			if err != nil {
				return nil, errors.Wrap(err, "invalid params")
			}

			model, err := params.Model(ctx, gh)
			if err != nil {
				return nil, errors.Wrap(err, "error fetching model")
			}

			return OKResponse(model), nil
		},
	},
}
