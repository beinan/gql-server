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

func (r *MockQueryResolver) GetUser(ctx Context, id ID) Future {
	user := User{
		Id:   id,
		Name: graphql.NewStringOption("beinan"),
	}
	userResolver := &EnhancedUserResolver{
		DefaultUserResolver{Value: future.MakeValue(user, nil)},
	}
	return future.MakeValue(userResolver, nil)
}

var rootQueryResolver = &MockQueryResolver{}

type EnhancedUserResolver struct {
	DefaultUserResolver
}

func (this EnhancedUserResolver) Friends(
	ctx Context, start int64, pageSize int64) Future {
	user1 := User{
		Id:   "friend1",
		Name: graphql.NewStringOption("friend1's name"),
	}
	userResolver1 := &EnhancedUserResolver{
		DefaultUserResolver{Value: future.MakeValue(user1, nil)},
	}
	resolvers := make([]UserResolver, 1)
	resolvers[0] = userResolver1
	return future.MakeValue(resolvers, nil)
}

func TestGraphQLResolver(t *testing.T) {
	doc, err := parser.ParseQuery(&ast.Source{Input: queryString})
	fmt.Println(doc, err)

	ctx := context.Background()

	var results Results
	for _, op := range doc.Operations {
		fmt.Println("Operation", "name:", op.Name, "Operation", op.Operation)
		switch op.Operation {
		case "query":
			results = ResolveSelections(ctx, op.SelectionSet, GqlQueryResolver(rootQueryResolver))
		default:
			future.MakeValue(nil, errors.New("unsupported opration"))
		}
	}
	json, jsonErr := json.Marshal(results)
	fmt.Println("json result:", string(json), jsonErr)

}

//generated test model and resolvers
type Future = future.Future
type Value = interface{}

type IDResolver interface {
	Id() Future
}

type UserResolver interface {
	Id() Future
	Name() Future
	Friends(ctx Context, start int64, pageSize int64) Future
}

type DefaultUserResolver struct {
	Value future.Future // future of User
}

func (this DefaultUserResolver) Id() Future {
	return this.Value.Then(func(value Value) (Value, error) {
		user := value.(User)
		return user.Id, nil
	})
}

func (this DefaultUserResolver) Name() Future {
	return this.Value.Then(func(value Value) (Value, error) {
		user := value.(User)
		return user.Name, nil
	})
}

func (this DefaultUserResolver) Friends(
	ctx Context, start int64, pageSize int64) Future {
	return future.MakeValue(nil, errors.New("Friends not implemented"))
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
	GetUser(ctx Context, id ID) Future
}

func GqlUserResolver(r UserResolver) FieldResolver {
	return func(
		ctx Context,
		field *graphql.Field,
	) future.Future {
		switch field.Name {

		case "id":

			//for immediate value
			return r.Id()

		case "name":

			//for immediate value
			return r.Name()

		case "friends":

			//for future resolver value
			span, ctx := logging.StartSpanFromContext(ctx, "friends")
			defer span.Finish()

			startValue, _ := field.Arguments.ForName("start").Value.Value(nil)

			pageSizeValue, _ := field.Arguments.ForName("pageSize").Value.Value(nil)

			fu := r.Friends(ctx, startValue.(int64), pageSizeValue.(int64))

			//if it's array, resolver each element
			return HandleFutureUserResolverArray(ctx, fu, field.SelectionSet)

		default:
			panic("unsopported field")
		}
	}
}

func GqlQueryResolver(r QueryResolver) FieldResolver {
	return func(
		ctx Context,
		field *graphql.Field,
	) future.Future {
		switch field.Name {

		case "getUser":

			//for future resolver value
			span, ctx := logging.StartSpanFromContext(ctx, "getUser")
			defer span.Finish()

			idValue, _ := field.Arguments.ForName("id").Value.Value(nil)

			fu := r.GetUser(ctx, idValue.(ID))

			//not array
			return HandleFutureUserResolver(ctx, fu, field.SelectionSet)

		default:
			panic("unsopported field")
		}
	}
}

func HandleFutureUserResolver(
	ctx Context,
	futureResolver Future,
	sels ast.SelectionSet,
) Future {
	return futureResolver.Then(func(data Value) (Value, error) {
		resolver := data.(UserResolver)
		result := ResolveSelections(ctx, sels, GqlUserResolver(resolver))
		return result, nil
	})
}

func HandleFutureUserResolverArray(
	ctx Context,
	futureResolverArray Future,
	sels ast.SelectionSet,
) Future {
	return futureResolverArray.Then(func(data Value) (Value, error) {
		resolverArray := data.([]UserResolver)
		results := make([]Results, len(resolverArray))
		for i, resolver := range resolverArray {
			results[i] = ResolveSelections(ctx, sels, GqlUserResolver(resolver))
		}
		return results, nil
	})
}
