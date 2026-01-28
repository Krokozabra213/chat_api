// Package repository provides data access layer for chat application.
package repository

import (
	"context"
	"time"
)

// EnsureCtxTimeout adds timeout to context if not already set.
// Returns original context with no-op cancel if deadline exists.
func EnsureCtxTimeout(ctx context.Context, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		return context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, func() {}
}
