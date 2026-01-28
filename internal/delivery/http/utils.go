// Package handler provides HTTP handlers for API.
package handler

import (
	"encoding/json"
	"net/http"
)

// JSONResp sends JSON response with the given status code.
func JSONResp(w http.ResponseWriter, statusCode int, resp any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(resp)
}

// JSONError sends JSON error response.
func JSONError(w http.ResponseWriter, statusCode int, message string) error {
	return JSONResp(w, statusCode, map[string]string{
		"error": message,
	})
}

// clamp limits value between min and max.
func clamp(val, minVal, maxVal int) int {
	if val < minVal {
		return minVal
	}
	if val > maxVal {
		return maxVal
	}
	return val
}
