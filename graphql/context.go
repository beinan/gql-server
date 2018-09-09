package graphql

import "context"

func MakeCtx(ctx context.Context) Context {
	return Context{
		ctx: ctx,
	}
}

type Context struct {
	ctx context.Context
}
