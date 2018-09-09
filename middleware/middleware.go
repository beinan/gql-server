package middleware

import (
	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/graphql"
)

// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// type Middleware func(Endpoint) Endpoint

type Service func(graphql.Context, interface{}) future.Future

type Filter func(Service) Service
