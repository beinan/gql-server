package middleware

import (
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

func InitHttpHandler(logger logging.Logger) http.Handler {
	return &httpHandler{logger, GraphqlMiddleware(logger)(nil)}
}

type httpHandler struct {
	logger      logging.Logger
	GQLEndpoint Endpoint
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var gqlRequest GQLRequest
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
