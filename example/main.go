package main

import (
	"net/http"

	"github.com/beinan/gql-server/example/dao"
	"github.com/beinan/gql-server/example/resolvers"
	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/middleware"
)

//go:generate sh -c "gql-server gen model > ./gen/model.go"
//go:generate sh -c "gql-server gen resolver > ./gen/resolver.go"
//go:generate sh -c "gql-server gen gqlresolver > ./gen/gql_resolver.go"
func main() {
	var logger = logging.StandardLogger(logging.DEBUG)
	defer logging.InitOpenTracing("gql-service").Close()

	logger.Debug("server starting...")

	dao, batcherAttacher := dao.MakeDAO()

	resolvers := resolvers.MkRootResolvers(dao)

	graphqlHandler := middleware.InitHttpHandler(logger, resolvers)

	http.Handle("/query", batcherAttacher(graphqlHandler))
	logger.Info(http.ListenAndServe(":8888", nil))
}
