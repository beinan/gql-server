package resolvers

import (
	"context"
	"fmt"

	"github.com/beinan/gql-server/example/dao"
	"github.com/beinan/gql-server/example/gen"
	"github.com/beinan/gql-server/graphql"
)

type ID = string
type StringOption = graphql.StringOption
type User = gen.User
type Context = context.Context

var db = make(map[string]*User)
var friendDB = make(map[string][]string)

type RootQueryResolver struct {
	dao *dao.DAO
}

func (r *RootQueryResolver) GetUser(ctx Context, id ID) gen.UserResolver {
	fmt.Println("Graphql Resolver: getUser:", id)
	userFuture := r.dao.GetUser(ctx, id)

	resolver := EnhancedUserResolver{
		r.dao,
		id,
		gen.FutureUserResolver{Value: userFuture},
	}
	return resolver
}

func (r *RootQueryResolver) GetUsers(ctx Context, start int64, pageSize int64) []gen.UserResolver {

	results := make([]gen.UserResolver, 1)
	userFuture := r.dao.GetUser(ctx, "1")

	results[0] = EnhancedUserResolver{
		r.dao,
		"1",
		gen.FutureUserResolver{Value: userFuture},
	}
	return results
}

type EnhancedUserResolver struct {
	dao    *dao.DAO
	userId ID
	gen.FutureUserResolver
}

func (r EnhancedUserResolver) Friends(
	ctx Context, start int64, pageSize int64) []gen.UserResolver {
	userFutures, _ := r.dao.GetFriends(ctx, r.userId, start, pageSize)
	resolvers := make([]gen.UserResolver, len(userFutures))
	for i, userFuture := range userFutures {
		userResolver := EnhancedUserResolver{
			r.dao,
			r.userId,
			gen.FutureUserResolver{Value: userFuture},
		}

		resolvers[i] = userResolver
	}
	return resolvers
}

func MkRootQueryResolver(dao *dao.DAO) gen.QueryResolver {
	return &RootQueryResolver{
		dao: dao,
	}
}
