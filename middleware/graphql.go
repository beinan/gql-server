package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/resolver"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type Context = context.Context

type graphqlService struct {
	logger            logging.Logger
	rootQueryResolver resolver.FieldResolver
}

func CreateGraphqlService(logger logging.Logger, rootQueryResolver resolver.FieldResolver) Service {
	return graphqlService{
		logger:            logger,
		rootQueryResolver: rootQueryResolver,
	}.serve
}

func (g graphqlService) serve(ctx Context, request interface{}) future.Future {
	gqlRequest := request.(GQLRequest)
	doc, err := parser.ParseQuery(&ast.Source{Input: gqlRequest.Query})
	var results resolver.Results
	for _, op := range doc.Operations {
		g.logger.Debug("Operation", "name:", op.Name, "Operation", op.Operation)
		switch op.Operation {
		case "query":
			results = resolver.ResolveSelections(ctx, op.SelectionSet, g.rootQueryResolver)
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
