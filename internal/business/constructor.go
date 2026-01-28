// Package business implements core application logic.
package business

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/test_api/internal/domain"
)

// ChatDBProvider defines methods for chat persistence operations.
type ChatDBProvider interface {
	SaveChat(ctx context.Context, title string) (*domain.Chat, error)
	GetChat(ctx context.Context, chatID int64) (*domain.Chat, error)
	DeleteChat(ctx context.Context, chatID int64) error
}

// MessageDBProvider defines methods for message persistence operations.
type MessageDBProvider interface {
	SaveMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID int64, limit int) ([]domain.Message, error)
}

// Business contains the core business logic and dependencies.
type Business struct {
	log             *slog.Logger
	chatProvider    ChatDBProvider
	messageProvider MessageDBProvider
}

// New creates a new Business instance with the provided dependencies.
func New(slogger *slog.Logger, chatProvider ChatDBProvider, messageProvider MessageDBProvider) *Business {
	return &Business{
		log:             slogger,
		chatProvider:    chatProvider,
		messageProvider: messageProvider,
	}
}
