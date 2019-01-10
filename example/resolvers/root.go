package resolvers

import (
	"context"
	"fmt"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/example/dao"
	"github.com/beinan/gql-server/example/gen"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/middleware"
)

type ID = string
type Value = future.Value
type StringOption = graphql.StringOption
type User = gen.User
type Context = context.Context

var db = make(map[string]*User)
var friendDB = make(map[string][]string)

//MkRootResolvers makes an instance of middleware.GQLResolvers as a root entry of all the custom resolvers
func MkRootResolvers(dao *dao.DAO) middleware.GQLResolvers {
	return middleware.GQLResolvers{
		RootQueryResolver: gen.MkGqlQueryResolver(&RootQueryResolver{
			dao: dao,
		}),
		RootMutationResolver: gen.MkGqlMutationResolver(&RootMutationResolver{
			dao: dao,
		}),
	}
}

type RootMutationResolver struct {
	dao *dao.DAO
}
type RootQueryResolver struct {
	dao *dao.DAO
}

func (r *RootMutationResolver) UpdateUserName(ctx Context, id ID, name string) gen.UserResolver {
	userFuture := r.dao.GetUser(ctx, id).Then(func(value future.Value) (Value, error) {
		userValue := value.(User)
		userValue.Name = graphql.StringOption{Value: name}
		return userValue, nil
	})

	resolver := EnhancedUserResolver{
		r.dao,
		id,
		gen.FutureUserResolver{Value: userFuture},
	}
	return resolver
}

func (r *RootMutationResolver) UpdateUser(ctx Context, id ID, userInput gen.UserInput) gen.UserResolver {
	userFuture := r.dao.GetUser(ctx, id).Then(func(value future.Value) (Value, error) {
		userValue := value.(User)
		userValue.Name = userInput.Name
		userValue.Email = userInput.Email
		fmt.Println("user name", userInput, userValue)
		return userValue, nil
	})

	resolver := EnhancedUserResolver{
		r.dao,
		id,
		gen.FutureUserResolver{Value: userFuture},
	}
	return resolver
}

func (r *RootQueryResolver) GetUser(ctx Context, id ID) gen.UserResolver {
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
