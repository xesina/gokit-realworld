package http

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, response interface{}, code int) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
