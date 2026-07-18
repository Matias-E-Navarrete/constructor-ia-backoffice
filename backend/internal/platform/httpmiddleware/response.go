package httpmiddleware

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorBody{Error: code, Message: message})
}
