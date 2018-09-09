package middleware

import (
	"fmt"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type graphqlService struct {
	logger logging.Logger
}

func CreateGraphqlService(logger logging.Logger) Service {
	return graphqlService{
		logger: logger,
	}.serve
}

func (g graphqlService) serve(ctx graphql.Context, request interface{}) future.Future {
	gqlRequest := request.(GQLRequest)
	doc, err := parser.ParseQuery(&ast.Source{Input: gqlRequest.Query})
	results := make(map[string]future.Future)
	for _, op := range doc.Operations {
		g.logger.Debug("Operation", "name:", op.Name, "Operation", op.Operation)
		switch op.Operation {
		case "query":
			results = resolveSelections(ctx, op.SelectionSet, g.resolveQueryField)
		default:
			panic("unsupport operation")
		}
	}
	g.logger.Info(fmt.Sprintf("%#v", doc.Operations[0].SelectionSet[0]), err)
	return future.MakeValue(GQLResponse{
		Data:  results,
		Error: nil,
	}, nil)
}

func resolveSelections(
	ctx graphql.Context,
	sels ast.SelectionSet,
	resolver func(graphql.Context, *ast.Field) future.Future,
) map[string]future.Future {
	results := make(map[string]future.Future)
	for _, selection := range sels {
		switch selection.(type) {
		case *ast.Field:
			field := selection.(*ast.Field)
			result := resolver(ctx, field)
			results[field.Alias] = result
		default:
			panic("selection type not supported yet")
		}
	}
	return results
}

var query Query = Query{
	GetUser: func(ctx Context, id ID) (*User, error) {
		return &User{
			Id:      id,
			Name:    graphql.NewStringOption("Name of " + id),
			Friends: UserFriends(id),
		}, nil
	},
}

func UserFriends(userId ID) func(Context, int, int) ([]User, error) {
	return func(ctx Context, start int, pageSize int) ([]User, error) {
		var user1 User = User{
			Id:      "1",
			Name:    graphql.NewStringOption("bbbb"),
			Friends: UserFriends("1"),
		}

		var user2 User = User{
			Id:      "2",
			Name:    graphql.NewStringOption("aaaa"),
			Friends: UserFriends("2"),
		}

		return []User{
			user1,
			user2,
		}, nil
	}
}

func (g graphqlService) resolveQueryField(ctx graphql.Context, field *ast.Field) future.Future {
	switch field.Name {
	case "getUser":
		fu := future.MakeFuture(func() (interface{}, error) {
			idValue, _ := field.Arguments.ForName("id").Value.Value(nil)
			return query.GetUser(ctx, idValue.(string))
		})
		return fu.Then(func(data interface{}, err error) (interface{}, error) {
			user := *data.(*User)
			userResolver := UserResolver{user}
			return resolveSelections(ctx, field.SelectionSet, userResolver.resolverQueryField), nil
		})
	default:
		panic("unsupported field")
	}
}

type UserResolver struct {
	data User
}

func (r UserResolver) resolverQueryField(
	ctx graphql.Context,
	field *ast.Field,
) future.Future {
	switch field.Name {
	case "id":
		return future.MakeValue(r.data.Id, nil)
	case "name":
		return future.MakeValue(r.data.Name, nil)
	case "friends":
		fu := future.MakeFuture(func() (interface{}, error) {
			//idValue, _ := field.Arguments.ForName("id").Value.Value(nil)
			return r.data.Friends(ctx, 0, 1)
		})
		return fu.Then(func(data interface{}, err error) (interface{}, error) {
			users := data.([]User)
			results := make([]map[string]future.Future, len(users))
			for i, user := range users {
				userResolver := UserResolver{user}
				results[i] = resolveSelections(
					ctx,
					field.SelectionSet,
					userResolver.resolverQueryField,
				)
			}
			return future.MakeValue(results, nil), nil
		})

	default:
		panic("unsopported field")
	}
}

type ID = string
type StringOption = graphql.StringOption

type Context = graphql.Context

type User struct {
	Id ID

	Name StringOption

	Friends func(ctx Context, start int, pageSize int) ([]User, error)
}

type Query struct {
	GetUser func(ctx Context, id ID) (*User, error)
}
