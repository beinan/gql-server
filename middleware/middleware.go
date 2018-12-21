package middleware

import (
	"context"
	"net/http"

	"github.com/beinan/gql-server/concurrent/future"
)

// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// type Middleware func(Endpoint) Endpoint

type Service func(context.Context, interface{}) future.Future

type Filter func(Service) Service

type HttpFilter func(next http.Handler) http.Handler

func Chain(chain ...HttpFilter) http.Handler {
	if len(chain) == 1 {
		return chain[0](nil)
	}
	rest := chain[1:] //drop the first one: chain[0]
	return chain[0](Chain(rest...))
}
