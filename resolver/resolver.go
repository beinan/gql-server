package resolver

import (
	"context"

	"github.com/beinan/gql-server/concurrent/future"
	"github.com/vektah/gqlparser/ast"
)

type Context = context.Context

type FieldResolver = func(Context, *ast.Field) future.Future

type Value = interface{}

func ResolveSelections(
	ctx Context,
	sels ast.SelectionSet,
	fieldResolver FieldResolver,
) Results {
	results := make(Results, len(sels))
	for i, selection := range sels {
		switch selection.(type) {
		case *ast.Field:
			field := selection.(*ast.Field)
			futureValue := fieldResolver(ctx, field)
			results[i] = Result{
				Alias:       field.Alias,
				FutureValue: futureValue,
			}
		default:
			panic("selection type not supported yet")
		}
	}
	return results
}
