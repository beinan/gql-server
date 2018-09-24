package middleware

import (
	"errors"
	"fmt"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/graphql"
	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type Resolver interface {
	ResolveQueryField(graphql.Context, *ast.Field) future.Future
}
type graphqlService struct {
	logger            logging.Logger
	rootQueryResolver Resolver
}

func CreateGraphqlService(logger logging.Logger, rootQueryResolver Resolver) Service {
	return graphqlService{
		logger:            logger,
		rootQueryResolver: rootQueryResolver,
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
			results = ResolveSelections(ctx, op.SelectionSet, g.rootQueryResolver)
		default:
			future.MakeValue(nil, errors.New("unsupported opration"))
		}
	}
	g.logger.Info(fmt.Sprintf("%#v", doc.Operations[0].SelectionSet[0]), err)
	return future.MakeValue(GQLResponse{
		Data:  results,
		Error: nil,
	}, nil)
}

func ResolveSelections(
	ctx graphql.Context,
	sels ast.SelectionSet,
	resolver Resolver,
) map[string]future.Future {
	results := make(map[string]future.Future)
	for _, selection := range sels {
		switch selection.(type) {
		case *ast.Field:
			field := selection.(*ast.Field)
			result := resolver.ResolveQueryField(ctx, field)
			results[field.Alias] = result
		default:
			panic("selection type not supported yet")
		}
	}
	return results
}
