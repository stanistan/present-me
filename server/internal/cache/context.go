package cache

import (
	"context"
)

type ctxKey int

var ctxKeyVal ctxKey

func ContextWithOptions(ctx context.Context, v *Options) context.Context {
	return context.WithValue(ctx, ctxKeyVal, v)
}

func OptionsFromContext(ctx context.Context) (*Options, bool) {
	v, ok := ctx.Value(ctxKeyVal).(*Options)
	return v, ok
}
