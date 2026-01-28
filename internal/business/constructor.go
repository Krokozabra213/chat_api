package business

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/test_api/internal/domain"
)

type ChatProvider interface {
	SaveChat(ctx context.Context, title string) (*domain.Chat, error)
	GetChat(ctx context.Context, chatID int64) (*domain.Chat, error)
}

type MessageProvider interface {
	SaveMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID int64, limit int) ([]domain.Message, error)
}

type Business struct {
	log             *slog.Logger
	chatProvider    ChatProvider
	messageProvider MessageProvider
}

func New(slogger *slog.Logger, chatProvider ChatProvider, messageProvider MessageProvider) *Business {
	return &Business{
		log:             slogger,
		chatProvider:    chatProvider,
		messageProvider: messageProvider,
	}
}
