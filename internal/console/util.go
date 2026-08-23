package console

import (
	"encoding/json"
	"net/http"
)

// writeJSON encodes a response with the given status code.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// readJSON decodes a JSON request body into value.
func readJSON(r *http.Request, value any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(value)
}
