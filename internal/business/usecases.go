// Package business implements core application logic.
package business

import (
	"context"
	"errors"

	"github.com/Krokozabra213/test_api/internal/domain"
	"github.com/Krokozabra213/test_api/internal/repository/postgres"
)

// CreateChat creates a new chat with the given title.
func (b *Business) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	chat, err := b.chatProvider.SaveChat(ctx, title)
	if err != nil {
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		return nil, ErrInternal
	}

	return chat, nil
}

// DeleteChat removes a chat by its ID.
func (b *Business) DeleteChat(ctx context.Context, chatID int64) error {
	err := b.chatProvider.DeleteChat(ctx, chatID)
	if err != nil {
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return ErrTimeout
		}
		return ErrInternal
	}

	return nil
}

// CreateMessage adds a new message to the specified chat.
func (b *Business) CreateMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error) {
	message, err := b.messageProvider.SaveMessage(ctx, chatID, text)
	if err != nil {
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		if errors.Is(err, postgres.ErrValidation) {
			return nil, ErrChatNotFound
		}
		return nil, ErrInternal
	}

	return message, nil
}

// ReadChatMessages retrieves a chat with its messages up to the given limit.
func (b *Business) ReadChatMessages(ctx context.Context, chatID int64, limit int) (*domain.ChatMessageOutput, error) {
	chat, err := b.chatProvider.GetChat(ctx, chatID)
	if err != nil {
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, ErrInternal
	}

	messages, err := b.messageProvider.GetMessages(ctx, chatID, limit)
	if err != nil {
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		return nil, ErrInternal
	}

	output := domain.NewChatMessageOutput(chat.ID, chat.Title, chat.CreatedAt, messages)
	return output, nil
}
