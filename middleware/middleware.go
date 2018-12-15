package middleware

import (
	"context"

	"github.com/beinan/gql-server/concurrent/future"
)

// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// type Middleware func(Endpoint) Endpoint

type Service func(context.Context, interface{}) future.Future

type Filter func(Service) Service
