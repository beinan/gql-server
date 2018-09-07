package main

import (
	"net/http"

	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/middleware"
)

func main() {
	var logger = logging.StandardLogger(logging.DEBUG)
	logger.Debug("server starting...")

	http.Handle("/query", middleware.InitHttpHandler(logger))
	logger.Info(http.ListenAndServe(":8888", nil))
}
