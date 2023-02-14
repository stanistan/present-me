package main

import "net/http"

var apiRoutes []Route = []Route{
	{
		"GET", "/ping",
		func(_ *http.Request) (*JSONResponse, error) {
			return OKResponse(map[string]string{"ping": "pong"}), nil
		},
	},
	{
		"GET", "/review",
		func(r *http.Request) (*JSONResponse, error) {
			return OKResponse(nil), nil
		},
	},
}
