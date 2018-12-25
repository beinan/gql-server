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
	rootQueryResolver resolver.GqlResolver
}

func CreateGraphqlService(logger logging.Logger, rootQueryResolver resolver.GqlResolver) Service {
	return graphqlService{
		logger:            logger,
		rootQueryResolver: rootQueryResolver,
	}.serve
}

func (g graphqlService) serve(ctx Context, request interface{}) future.Future {
	gqlRequest := request.(GQLRequest)
	doc, err := parser.ParseQuery(&ast.Source{Input: gqlRequest.Query})
	var results resolver.GqlResults
	for _, op := range doc.Operations {
		g.logger.Debug("Operation", "name:", op.Name, "Operation", op.Operation)
		switch op.Operation {
		case "query":
			results = g.rootQueryResolver.Resolve(ctx, op.SelectionSet)
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
