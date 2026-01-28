package business

import (
	"context"
	"errors"

	"github.com/Krokozabra213/test_api/internal/domain"
	"github.com/Krokozabra213/test_api/internal/repository"
)

func (b *Business) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	chat, err := b.chatProvider.SaveChat(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrCtxCancelled) || errors.Is(err, repository.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		return nil, ErrInternal
	}

	return chat, nil
}

func (b *Business) CreateMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error) {
	chat, err := b.messageProvider.SaveMessage(ctx, chatID, text)
	if err != nil {
		if errors.Is(err, repository.ErrCtxCancelled) || errors.Is(err, repository.ErrCtxDeadline) {
			return nil, ErrTimeout
		}
		if errors.Is(err, repository.ErrValidation) {
			return nil, ErrChatNotFound
		}
		return nil, ErrInternal
	}

	return chat, nil
}
