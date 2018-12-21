package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/beinan/gql-server/logging"
	"github.com/beinan/gql-server/resolver"
)

type GQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func InitHttpHandler(logger logging.Logger, rootQueryResolver resolver.FieldResolver) http.Handler {
	return &httpHandler{logger, CreateGraphqlService(logger, rootQueryResolver)}
}

type httpHandler struct {
	logger     logging.Logger
	gqlService Service
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var gqlRequest GQLRequest
	if err := json.NewDecoder(r.Body).Decode(&gqlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	span, ctx := logging.StartSpanFromContext(ctx, "http_request")
	defer span.Finish()

	response, err := h.gqlService(ctx, gqlRequest).Value()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
