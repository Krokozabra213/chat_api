package business

import "errors"

var (
	ErrTimeout  = errors.New("timeout request")
	ErrInternal = errors.New("internal service error")

	ErrChatNotFound = errors.New("chat not found")
)
