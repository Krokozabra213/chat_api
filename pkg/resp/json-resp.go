package resp

import (
	"encoding/json"
	"net/http"
)

// JsonResp sends JSON response with the given status code.
func JsonResp(w http.ResponseWriter, statusCode int, resp any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(resp)
}

// JsonError sends JSON error response.
func JsonError(w http.ResponseWriter, statusCode int, message string) error {
	return JsonResp(w, statusCode, map[string]string{
		"error": message,
	})
}
