package resolver

import (
	"context"

	"github.com/vektah/gqlparser/ast"
)

//Context is the alias to context.Context
type Context = context.Context

type GqlResolver interface {
	Resolve(Context, ast.SelectionSet) GqlResults
	//resolveField(Context, *ast.Field) GqlResultValue
}
type GqlResultValue = interface{}

type StringResolver interface {
	Value() (GqlResultValue, error)
}

type IDResolver interface {
	Value() (GqlResultValue, error)
}

//FieldResolver is a function to resolve a single field
type FieldResolver = func(Context, *ast.Field) (GqlResultValue, error)

type Value = interface{}

func GqlResolveValues(ctx Context, resolvers []GqlResolver, sels ast.SelectionSet) GqlResultValue {
	results := make([]GqlResults, len(resolvers))
	for i, resolver := range resolvers {
		results[i] = resolver.Resolve(ctx, sels)
	}
	return results
}

//GqlResolveSelections is a util function to resolve each field in selections
func GqlResolveSelections(
	ctx Context,
	sels ast.SelectionSet,
	fieldResolver FieldResolver,
) GqlResults {
	results := make(GqlResults, len(sels))
	for i, selection := range sels {
		switch selection.(type) {
		case *ast.Field:
			field := selection.(*ast.Field)
			resultValue, _ := fieldResolver(ctx, field) //TODO: handle error
			results[i] = GqlResult{
				Alias: field.Alias,
				Value: resultValue,
			}
		default:
			panic("selection type not supported yet")
		}
	}
	return results
}
