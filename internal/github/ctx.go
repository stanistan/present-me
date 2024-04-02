package github

import "context"

type ctxKey struct{}

func Ctx(ctx context.Context) (*Client, bool) {
	v, ok := ctx.Value(ctxKey{}).(*Client)
	return v, ok
}

func WithContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, ctxKey{}, client)
}
