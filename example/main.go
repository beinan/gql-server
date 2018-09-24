package main

import (
	"fmt"
	"net/http"

	"github.com/beinan/gql-server/example/gen"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/middleware"
)

//go:generate sh -c "gql-server gen model > ./gen/model.go"
//go:generate sh -c "gql-server gen resolver > ./gen/resolver.go"
func main() {
	var logger = logging.StandardLogger(logging.DEBUG)
	logger.Debug("server starting...")

	rootQueryResolver := gen.QueryResolver{Data: query}
	http.Handle("/query", middleware.InitHttpHandler(logger, rootQueryResolver))
	logger.Info(http.ListenAndServe(":8888", nil))
}

type ID = string
type StringOption = graphql.StringOption
type User = gen.User
type Query = gen.Query
type Context = graphql.Context

var query = &gen.Query{
	GetUser: func(ctx Context, id ID) (*User, error) {
		fmt.Println("exec getUser")
		return &User{
			Id:      id,
			Name:    graphql.NewStringOption("Name of " + id),
			Friends: UserFriends(id),
		}, nil
	},
}

func UserFriends(userId ID) func(Context, int64, int64) ([]*User, error) {
	return func(ctx Context, start int64, pageSize int64) ([]*User, error) {
		fmt.Println("exec friends")
		var user1 = &User{
			Id:      "1",
			Name:    graphql.NewStringOption("bbbb"),
			Friends: UserFriends("1"),
		}

		var user2 = &User{
			Id:      "2",
			Name:    graphql.NewStringOption("aaaa"),
			Friends: UserFriends("2"),
		}

		return []*User{
			user1,
			user2,
		}, nil
	}
}
