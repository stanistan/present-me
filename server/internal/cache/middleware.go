package cache

import (
	"net/http"
)

type RequestOptions func(r *http.Request) *Options

func Middleware(c *Cache, opts RequestOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithCache(r.Context(), c)
			if opts != nil {
				if cOpts := opts(r); opts != nil {
					ctx = ContextWithOptions(ctx, cOpts)
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
