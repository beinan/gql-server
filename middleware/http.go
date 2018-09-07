package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/beinan/gql-server/logging"
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

func GraphqlMiddleware(next Endpoint) Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		gqlRequest := request.(GQLRequest)
		return GQLResponse{
			Data:  gqlRequest,
			Error: nil,
		}, nil
	}
}

func InitHttpHandler(logger logging.Logger) http.Handler {
	return &httpHandler{logger, GraphqlMiddleware(nil)}
}

type httpHandler struct {
	logger      logging.Logger
	GQLEndpoint Endpoint
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var gqlRequest GQLRequest
	h.logger.Debug("Reqeust:", r)
	if err := json.NewDecoder(r.Body).Decode(&gqlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.GQLEndpoint(r.Context(), gqlRequest)
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
