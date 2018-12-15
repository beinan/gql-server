package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/beinan/gql-server/example/gen"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/middleware"
	"github.com/opentracing/opentracing-go/log"
)

//go:generate sh -c "gql-server gen model > ./gen/model.go"
//go:generate sh -c "gql-server gen resolver > ./gen/resolver.go"
func main() {
	var logger = logging.StandardLogger(logging.DEBUG)
	defer logging.InitOpenTracing("gql-service").Close()

	logger.Debug("server starting...")

	rootQueryResolver := gen.QueryResolver{Data: query}
	http.Handle("/query", middleware.InitHttpHandler(logger, rootQueryResolver))
	logger.Info(http.ListenAndServe(":8888", nil))
}

type ID = string
type StringOption = graphql.StringOption
type User = gen.User
type Query = gen.Query
type Context = context.Context

var db = make(map[string]*User)
var friendDB = make(map[string][]string)

func init() {
	makeTestUser := func(id string) *User {
		return &User{
			Id:      id,
			Name:    graphql.NewStringOption("User_" + id),
			Friends: UserFriends(id),
		}
	}
	db["1"] = makeTestUser("1")
	db["2"] = makeTestUser("2")
	db["3"] = makeTestUser("3")
	friendDB["1"] = []string{"2", "3"}
	friendDB["2"] = []string{"1", "3"}
}

func getUser(ctx Context, id string) *User {
	span, ctx := logging.StartSpanFromContext(ctx, "read_user_from_db")
	span.LogFields(
		log.String("id", id),
	)

	defer span.Finish()
	time.Sleep(10 * time.Millisecond)
	return db[id]
}

var query = &gen.Query{
	GetUser: func(ctx Context, id ID) (*User, error) {
		fmt.Println("get user:" + id)
		return getUser(ctx, id), nil
	},
}

func UserFriends(userId ID) func(Context, int64, int64) ([]*User, error) {
	return func(ctx Context, start int64, pageSize int64) ([]*User, error) {
		ids := friendDB[userId]
		users := make([]*User, len(ids))
		for i, id := range ids {
			users[i] = getUser(ctx, id)
		}
		return users, nil
	}
}
