//Generated by gql-server
//DO NOT EDIT
package gen

import (
	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/graphql"
)

type ID = string
type StringOption = graphql.StringOption

type Future = future.Future

type User struct {
	Id ID

	Name StringOption
}

type Query struct {
}

type Mutation struct {
}

type UserInput struct {
	Name StringOption

	Email StringOption
}
