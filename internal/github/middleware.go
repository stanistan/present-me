package github

import (
	"net/http"
)

func Middleware(g *Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithContext(r.Context(), g)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
