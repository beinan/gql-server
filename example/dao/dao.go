package dao

import (
	"context"
	"strings"
	"time"

	"github.com/beinan/gql-server/concurrent/batcher"
	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/example/gen"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/middleware"
	"github.com/opentracing/opentracing-go/log"
)

type DAO struct {
	GetUser batcher.GetFunc
}

type ID = string
type StringOption = graphql.StringOption
type User = gen.User
type Context = context.Context
type Result = batcher.Result

var db = make(map[string]User)
var friendDB = make(map[string][]string)

func MakeDAO() (dao *DAO, batcherAttacher middleware.HttpFilter) {
	bm := batcher.MakeBatcherManager()
	dao = &DAO{
		GetUser: bm.Register(getUsers),
	}
	batcherAttacher = middleware.AttachBatcher(bm)

	makeTestUser := func(id string) User {
		return User{
			Id:   id,
			Name: graphql.NewStringOption("User_" + id),
		}
	}
	db["1"] = makeTestUser("1")
	db["2"] = makeTestUser("2")
	db["3"] = makeTestUser("3")
	friendDB["1"] = []string{"2", "3"}
	friendDB["2"] = []string{"1", "3"}
	return
}

func getUsers(ctx Context, ids []ID) []Result {
	span, ctx := logging.StartSpanFromContext(ctx, "read_users_from_db")
	span.LogFields(
		log.String("ids", strings.Join(ids, ",")),
	)
	defer span.Finish()
	time.Sleep(10 * time.Millisecond)
	results := make([]Result, len(ids))
	for i, id := range ids {
		results[i] = Result{
			Value: db[id],
			Err:   nil,
		}
	}
	return results
}

func (dao *DAO) GetFriends(ctx Context, userId ID, start int64, pageSize int64) ([]future.Future, error) {
	ids := friendDB[userId]
	userFutures := make([]future.Future, len(ids))
	for i, id := range ids {
		userFutures[i] = dao.GetUser(ctx, id)
	}
	return userFutures, nil

}
