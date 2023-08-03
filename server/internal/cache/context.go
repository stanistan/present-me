package cache

import (
	"context"
)

type ctxKey[T any] struct{}

func ContextWithOptions(ctx context.Context, v *Options) context.Context {
	return context.WithValue(ctx, ctxKey[Options]{}, v)
}

func OptionsFromContext(ctx context.Context) (*Options, bool) {
	v, ok := ctx.Value(ctxKey[Options]{}).(*Options)
	return v, ok
}

// disabledCache is used as a default cache.
var disabledCache *Cache = &Cache{disabled: true}

func ContextWithCache(ctx context.Context, c *Cache) context.Context {
	return context.WithValue(ctx, ctxKey[Cache]{}, c)
}

func Ctx(ctx context.Context) *Cache {
	v, ok := ctx.Value(ctxKey[Cache]{}).(*Cache)
	if !ok {
		return disabledCache
	}

	return v
}
