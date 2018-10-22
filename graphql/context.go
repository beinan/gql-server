package graphql

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

func MakeCtx(ctx context.Context) Context {
	return Context{
		ctx: ctx,
	}
}

type Context struct {
	ctx context.Context
}

func (c Context) StartSpanFromContext(name string) (opentracing.Span, Context) {
	span, newCtx := opentracing.StartSpanFromContext(c.ctx, name)
	return span, Context{ctx: newCtx}
}
