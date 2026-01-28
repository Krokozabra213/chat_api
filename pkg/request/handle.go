// Package request provides utilities for handling HTTP requests and responses.
package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrDecode is returned when JSON decoding fails.
var ErrDecode = errors.New("decode error")

// Validator interface for request validation.
type Validator interface {
	Validate() error
}

// DecodeAndValidate decodes JSON request body and validates it.
func DecodeAndValidate[T Validator](r *http.Request) (T, error) {
	var req T
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, fmt.Errorf("%w: %v", ErrDecode, err)
	}

	if err := req.Validate(); err != nil {
		return req, err
	}

	return req, nil
}
