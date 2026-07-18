// Package httpkit holds the tiny JSON request/response helpers shared by
// every context's interfaces/http handlers, so each handler stays a
// decode -> call use-case -> encode one-liner.
package httpkit

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, map[string]string{"error": code, "message": message})
}
