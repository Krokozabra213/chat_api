package repository

import (
	"errors"

	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
)

var (
	ErrCtxCancelled = errors.New("context cancelled error")
	ErrCtxDeadline  = errors.New("context deadline error")
	ErrUnknown      = errors.New("unknown error")

	//postgres errors
	ErrValidation = errors.New("validation error")
	ErrDuplicate  = errors.New("duplicate key error")
	ErrNotFound   = errors.New("not found error")
	ErrInternal   = errors.New("internal error")
)

func ErrorFactory(entity string, err error) error {
	// errors postgresclient package
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
