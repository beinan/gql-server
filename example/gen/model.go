package genmodel

import "github.com/beinan/gql-server/graphql"

type ID = string
type StringOption = graphql.StringOption

type Context = graphql.Context

type User struct {
	Id ID

	Name StringOption

	Friends func(ctx Context, arg struct {
		start    int
		pageSize int
	}) ([]User, error)
}

type Query struct {
	GetUser func(ctx Context, arg struct {
		id ID
	}) (*User, error)
}
