package cache

import (
	"net/http"
)

type RequestOptions func(r *http.Request) *Options

// Middleware propagates a Cache to each request, and
// injects the caching options for the request.
func Middleware(c *Cache, opts RequestOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithCache(r.Context(), c)

			if opts != nil {
				if cOpts := opts(r); cOpts != nil {
					ctx = ContextWithOptions(ctx, cOpts)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
