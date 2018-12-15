//Generated by gql-server
//DO NOT EDIT
package gen

import "github.com/beinan/gql-server/concurrent/future"
import "github.com/beinan/gql-server/graphql"
import "github.com/beinan/gql-server/logging"
import "github.com/beinan/gql-server/middleware"

type UserResolver struct {
	Data *User
}

func (r UserResolver) ResolveQueryField(
	ctx Context,
	field *graphql.Field,
) future.Future {
	switch field.Name {

	case "id":

		//for immediate value
		return future.MakeValue(r.Data.Id, nil)

	case "name":

		//for immediate value
		return future.MakeValue(r.Data.Name, nil)

	case "friends":

		//for future resolver value
		span, ctx := logging.StartSpanFromContext(ctx, "friends")
		defer span.Finish()

		fu := future.MakeFuture(func() (interface{}, error) {

			startValue, _ := field.Arguments.ForName("start").Value.Value(nil)

			pageSizeValue, _ := field.Arguments.ForName("pageSize").Value.Value(nil)

			return r.Data.Friends(ctx, startValue.(int64), pageSizeValue.(int64))
		})

		//if it's array, resolver each element
		return fu.Then(func(data interface{}, err error) (interface{}, error) {
			values := data.([]*User) //array of elememnt type
			results := make([]map[string]future.Future, len(values))
			for i, value := range values {
				span, ctx := logging.StartSpanFromContext(ctx, "User")
				valueResolver := UserResolver{value}
				results[i] = middleware.ResolveSelections(ctx, field.SelectionSet, valueResolver)
				span.Finish()
			}
			return future.MakeValue(results, nil), nil
		})

	default:
		panic("unsopported field")
	}
}

type QueryResolver struct {
	Data *Query
}

func (r QueryResolver) ResolveQueryField(
	ctx Context,
	field *graphql.Field,
) future.Future {
	switch field.Name {

	case "getUser":

		//for future resolver value
		span, ctx := logging.StartSpanFromContext(ctx, "getUser")
		defer span.Finish()

		fu := future.MakeFuture(func() (interface{}, error) {

			idValue, _ := field.Arguments.ForName("id").Value.Value(nil)

			return r.Data.GetUser(ctx, idValue.(ID))
		})

		//not array
		return fu.Then(func(data interface{}, err error) (interface{}, error) {
			value := data.(*User)
			valueResolver := UserResolver{value}
			result := middleware.ResolveSelections(ctx, field.SelectionSet, valueResolver)
			return result, nil
		})

	case "getUsers":

		//for future resolver value
		span, ctx := logging.StartSpanFromContext(ctx, "getUsers")
		defer span.Finish()

		fu := future.MakeFuture(func() (interface{}, error) {

			startValue, _ := field.Arguments.ForName("start").Value.Value(nil)

			pageSizeValue, _ := field.Arguments.ForName("pageSize").Value.Value(nil)

			return r.Data.GetUsers(ctx, startValue.(int64), pageSizeValue.(int64))
		})

		//if it's array, resolver each element
		return fu.Then(func(data interface{}, err error) (interface{}, error) {
			values := data.([]*User) //array of elememnt type
			results := make([]map[string]future.Future, len(values))
			for i, value := range values {
				span, ctx := logging.StartSpanFromContext(ctx, "User")
				valueResolver := UserResolver{value}
				results[i] = middleware.ResolveSelections(ctx, field.SelectionSet, valueResolver)
				span.Finish()
			}
			return future.MakeValue(results, nil), nil
		})

	default:
		panic("unsopported field")
	}
}
