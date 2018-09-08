package middleware

import (
	"context"
	"fmt"

	"github.com/beinan/gql-server/logging"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func GraphqlMiddleware(logger logging.Logger) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			gqlRequest := request.(GQLRequest)
			doc, err := parser.ParseQuery(&ast.Source{Input: gqlRequest.Query})
			for _, op := range doc.Operations {
				logger.Debug("Operation", "name:", op.Name, "Operation", op.Operation)
				for _, selection := range op.SelectionSet {
					logger.Debug("selection", selection)
					switch selection.(type) {
					case *ast.Field:
						field := selection.(*ast.Field)
						logger.Debug("field:", field.Name, field.Alias)
						for _, arg := range field.Arguments {
							logger.Debug("arg", arg.Name, arg.Value)
						}
					default:
						panic("selection type not supported yet")
					}
				}
			}
			logger.Info(fmt.Sprintf("%#v", doc.Operations[0].SelectionSet[0]), err)
			return GQLResponse{
				Data:  gqlRequest,
				Error: nil,
			}, nil
		}
	}
}
