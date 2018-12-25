package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

var queryString string = `
query{
  user1:getUser(id:"1"){
    id
    name
    friends(start: 0, pageSize:10){
      id
      name
    }
  }
  user2:getUser(id:"2"){
    id
    name
    friends(start: 0, pageSize:10){
      id
      name
    }
  }
}
`

type MockResolver struct{}

func (r *MockResolver) fieldResolver(ctx Context, field *ast.Field) future.Future {
	fmt.Println("resolving field", field.Name)
	return future.MakeValue("abc", nil)
}

var resolver *MockResolver = &MockResolver{}

type MockQueryResolver struct{}

func (r *MockQueryResolver) GetUser(ctx Context, id ID) UserResolver {
	user := User{
		Id:   id,
		Name: graphql.NewStringOption("beinan"),
	}
	userResolver := &EnhancedUserResolver{
		FutureUserResolver{Value: future.MakeValue(user, nil)},
	}
	return userResolver
}

var rootQueryResolver = &MockQueryResolver{}

type EnhancedUserResolver struct {
	FutureUserResolver
}

func (r EnhancedUserResolver) Friends(
	ctx Context, start int64, pageSize int64) []UserResolver {
	user1 := User{
		Id:   "friend1",
		Name: graphql.NewStringOption("friend1's name"),
	}
	userResolver1 := &EnhancedUserResolver{
		FutureUserResolver{Value: future.MakeValue(user1, nil)},
	}
	resolvers := make([]UserResolver, 1)
	resolvers[0] = userResolver1
	return resolvers
}

func TestGraphQLResolver(t *testing.T) {
	doc, err := parser.ParseQuery(&ast.Source{Input: queryString})
	fmt.Println(doc, err)

	ctx := context.Background()

	var results GqlResults
	for _, op := range doc.Operations {
		fmt.Println("Operation", "name:", op.Name, "Operation", op.Operation)
		switch op.Operation {
		case "query":
			results = MkGqlQueryResolver(rootQueryResolver).resolve(ctx, op.SelectionSet)
		default:
			future.MakeValue(nil, errors.New("unsupported opration"))
		}
	}
	json, jsonErr := json.Marshal(results)
	fmt.Println("json result:", string(json), jsonErr)

}

//generated test model and resolvers
type Future = future.Future

type IDQuery struct {
	id      ID
	resolve func(Context, ID) Future
}
type IDQueryResolver interface {
	Id() ID
}

type UserResolver interface {
	Id() StringResolver
	Name() StringResolver
	Friends(ctx Context, start int64, pageSize int64) []UserResolver
}

type FutureUserResolver struct {
	Value future.Future // future of User
}

func (r FutureUserResolver) Id() StringResolver {
	return r.Value.Then(func(value Value) (Value, error) {
		user := value.(User)
		return user.Id, nil
	})
}

func (r FutureUserResolver) Name() StringResolver {
	return r.Value.Then(func(value Value) (Value, error) {
		user := value.(User)
		return user.Name, nil
	})
}

func (r FutureUserResolver) Friends(
	ctx Context, start int64, pageSize int64) []UserResolver {
	panic("Friends not implemented")
}

type ID = string
type StringOption = graphql.StringOption

type User struct {
	Id ID

	Name StringOption
}

type Query struct {
}

type QueryResolver interface {
	GetUser(ctx Context, id ID) UserResolver
}

type GqlUserResolver struct {
	resolver UserResolver
}

func MkGqlUserResolver(resolver UserResolver) GqlUserResolver {
	return GqlUserResolver{
		resolver: resolver,
	}
}

func MkGqlUserResolvers(resolvers []UserResolver) []GqlResolver {
	gqlResolvers := make([]GqlResolver, len(resolvers))
	for i, resolver := range resolvers {
		gqlResolvers[i] = GqlUserResolver{
			resolver: resolver,
		}
	}
	return gqlResolvers
}

func (r GqlUserResolver) resolve(ctx Context, sel ast.SelectionSet) GqlResults {
	return GqlResolveSelections(ctx, sel, r.resolveField)
}

func (r GqlUserResolver) resolveField(ctx Context, field *ast.Field) (GqlResultValue, error) {
	switch field.Name {
	case "id":

		//for immediate value
		return r.resolver.Id().Value()

	case "name":

		//for immediate value
		return r.resolver.Name().Value()

	case "friends":

		//for future resolver value
		span, ctx := logging.StartSpanFromContext(ctx, "friends")
		defer span.Finish()

		startValue, _ := field.Arguments.ForName("start").Value.Value(nil)

		pageSizeValue, _ := field.Arguments.ForName("pageSize").Value.Value(nil)

		resolvers := r.resolver.Friends(ctx, startValue.(int64), pageSizeValue.(int64))
		gqlResolvers := MkGqlUserResolvers(resolvers)
		//if it's array, resolver each element
		return GqlResolveValues(ctx, gqlResolvers, field.SelectionSet), nil

	default:
		panic("unsopported field")
	}

}

type GqlQueryResolver struct {
	resolver QueryResolver
}

func MkGqlQueryResolver(resolver QueryResolver) GqlQueryResolver {
	return GqlQueryResolver{
		resolver: resolver,
	}
}
func (g GqlQueryResolver) resolve(ctx Context, sels ast.SelectionSet) GqlResults {
	return GqlResolveSelections(ctx, sels, g.resolveQueryField)
}

func (g GqlQueryResolver) resolveQueryField(ctx Context, field *ast.Field) (GqlResultValue, error) {
	switch field.Name {

	case "getUser":

		//for future resolver value
		span, ctx := logging.StartSpanFromContext(ctx, "getUser")
		defer span.Finish()

		idValue, _ := field.Arguments.ForName("id").Value.Value(nil)

		resolver := g.resolver.GetUser(ctx, idValue.(ID))
		gqlResolver := MkGqlUserResolver(resolver)
		//not array
		return gqlResolver.resolve(ctx, field.SelectionSet), nil

	default:
		panic("unsopported field")
	}

}
