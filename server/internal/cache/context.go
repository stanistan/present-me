package cache

import (
	"context"
)

type ctxOptionsKey struct{}

func ContextWithOptions(ctx context.Context, v *Options) context.Context {
	return context.WithValue(ctx, ctxOptionsKey{}, v)
}

func OptionsFromContext(ctx context.Context) (*Options, bool) {
	v, ok := ctx.Value(ctxOptionsKey{}).(*Options)
	return v, ok
}

type ctxCacheKey struct{}

func ContextWithCache(ctx context.Context, c *Cache) context.Context {
	return context.WithValue(ctx, ctxCacheKey{}, c)
}

var disabledCache *Cache = &Cache{disabled: true}

func Ctx(ctx context.Context) *Cache {
	v, ok := ctx.Value(ctxCacheKey{}).(*Cache)
	if !ok {
		return disabledCache
	}

	return v
}
