// Package business implements core application logic.
package business

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Krokozabra213/test_api/internal/domain"
	"github.com/Krokozabra213/test_api/internal/repository/postgres"
)

// CreateChat creates a new chat with the given title.
func (b *Business) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	const op = "business.CreateChat"
	log := b.log.With(
		slog.String("op", op),
		slog.Int("title_len", len(title)),
	)
	log.Info("starting CreateChat process")

	chat, err := b.chatProvider.SaveChat(ctx, title)
	if err != nil {
		log.Error("failed to create chat", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		return nil, ErrInternal
	}
	log.Info("createChat success")

	return chat, nil
}

// DeleteChat removes a chat by its ID.
func (b *Business) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "business.DeleteChat"
	log := b.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)
	log.Info("starting DeleteChat process")

	err := b.chatProvider.DeleteChat(ctx, chatID)
	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return ErrTimeout
		}
		return ErrInternal
	}
	log.Info("deleteChat success")

	return nil
}

// CreateMessage adds a new message to the specified chat.
func (b *Business) CreateMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error) {
	const op = "business.CreateMessage"
	log := b.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)
	log.Info("starting CreateMessage process")

	message, err := b.messageProvider.SaveMessage(ctx, chatID, text)
	if err != nil {
		log.Error("failed to create message", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		if errors.Is(err, postgres.ErrValidation) {
			return nil, ErrChatNotFound
		}
		return nil, ErrInternal
	}
	log.Info("createMessage success")

	return message, nil
}

// ReadChatMessages retrieves a chat with its messages up to the given limit.
func (b *Business) ReadChatMessages(ctx context.Context, chatID int64, limit int) (*domain.ChatMessageOutput, error) {
	const op = "business.ReadChatMessages"
	log := b.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
		slog.Int("limit", limit),
	)
	log.Info("starting ReadChatMessages process")

	chat, err := b.chatProvider.GetChat(ctx, chatID)
	if err != nil {
		log.Error("failed to get chat", slog.String("error", err.Error()))
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
		log.Error("failed to get messages", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrCtxCancelled) || errors.Is(err, postgres.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		return nil, ErrInternal
	}
	log.Info("read messages success")

	output := domain.NewChatMessageOutput(chat.ID, chat.Title, chat.CreatedAt, messages)
	return output, nil
}
