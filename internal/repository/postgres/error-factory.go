// Package postgres provides data access layer for chat application.
package postgres

import (
	"errors"

	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
)

// Repository-level errors
var (
	// Context errors
	ErrCtxCancelled = errors.New("context cancelled error")
	ErrCtxDeadline  = errors.New("context deadline error")

	// Database errors
	ErrValidation = errors.New("validation error")
	ErrDuplicate  = errors.New("duplicate key error")
	ErrNotFound   = errors.New("not found error")
	ErrInternal   = errors.New("internal error")
	ErrUnknown    = errors.New("unknown error")
)

// ErrorFactory maps postgres client errors to repository-level errors.
func ErrorFactory(err error) error {
	var customErr *postgresclient.CustomError
	if errors.As(err, &customErr) {
		switch {
		case errors.Is(err, postgresclient.ErrCtxCancelled):
			return ErrCtxCancelled
		case errors.Is(err, postgresclient.ErrCtxDeadline):
			return ErrCtxDeadline
		case errors.Is(err, postgresclient.ErrValidation):
			return ErrValidation
		case errors.Is(err, postgresclient.ErrDuplicateKey):
			return ErrDuplicate
		case errors.Is(err, postgresclient.ErrNotFound):
			return ErrNotFound
		case errors.Is(err, postgresclient.ErrInternal):
			return ErrInternal
		}
	}

	return ErrUnknown
}
