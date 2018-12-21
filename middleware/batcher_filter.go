package middleware

import (
	"net/http"

	"github.com/beinan/gql-server/concurrent/batcher"
)

func AttachBatcher(bm *batcher.BatcherManager) HttpFilter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//Attach batcher to context
			ctx := bm.Attach(req.Context())
			//Using the request with the changed context
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
