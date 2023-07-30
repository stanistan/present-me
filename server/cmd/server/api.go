package main

import (
	"github.com/pkg/errors"

	"github.com/stanistan/present-me/internal/github"
	"github.com/stanistan/present-me/internal/http"
)

var apiRoutes = http.Routes(

	// GET /version returns the embeded version of the server/client
	http.GET("/version", func(_ *http.Request) (*http.JSONResponse, error) {
		return http.OKResponse(map[string]any{"version": version}), nil
	}),

	// GET /ping is just for basic API functionality.
	http.GET("/ping", func(_ *http.Request) (*http.JSONResponse, error) {
		return http.OKResponse(map[string]any{"ping": "pong"}), nil
	}),

	// GET /search finds a review based on input.
	http.GET("/search", func(r *http.Request) (*http.JSONResponse, error) {
		ctx := r.Context()
		gh, ok := github.Ctx(ctx)
		if !ok || gh == nil {
			return nil, errors.New("missing github context")
		}

		params, err := github.ReviewParamsFromURL(r.URL.Query().Get("search"))
		if err != nil {
			return nil, err
		}

		_, err = params.EnsureReviewID(ctx, gh)
		if err != nil {
			return &http.JSONResponse{
				Code: 404,
				Data: map[string]string{
					"msg": "No review id",
				},
			}, nil
		}

		return http.OKResponse(params), nil
	}),

	// GET /review hydrates the the full review.
	http.GET("/review", func(r *http.Request) (*http.JSONResponse, error) {
		var (
			ctx    = r.Context()
			values = r.URL.Query()
		)

		gh, ok := github.Ctx(ctx)
		if !ok || gh == nil {
			return nil, errors.New("missing github context")
		}

		params, err := github.ReviewParamsFromMap(github.ReviewParamsMap{
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

		return http.OKResponse(model), nil
	}),

	http.GET("/review2.json", func(r *http.Request) (*http.JSONResponse, error) {
		var (
			ctx    = r.Context()
			values = r.URL.Query()
		)

		source := github.ReviewSource{
			ReviewParamsMap: github.ReviewParamsMap{
				Owner:  values.Get("org"),
				Repo:   values.Get("repo"),
				Number: values.Get("pull"),
				Review: values.Get("review"),
			},
		}

		review, err := source.GetReview(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching review")
		}

		return http.OKResponse(review), nil
	}),
)
