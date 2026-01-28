// Package repository provides data access layer for chat application.
package repository

import (
	"context"
	"time"

	"github.com/Krokozabra213/test_api/internal/domain"
	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
	"gorm.io/gorm"
)

const (
	ctxTimeout = 5 * time.Second
)

// PostgresClient defines database operations interface.
type PostgresClient interface {
	WithContext(ctx context.Context) *gorm.DB
}

// PostgresRepository implements chat data storage using PostgreSQL.
type PostgresRepository struct {
	сlient PostgresClient
}

// NewPostgresRepository creates new repository instance.
func NewPostgresRepository(client PostgresClient) *PostgresRepository {
	return &PostgresRepository{
		сlient: client,
	}
}

// SaveChat persists new chat and returns it with generated ID & CreatedAt field.
func (r *PostgresRepository) SaveChat(ctx context.Context, title string) (*domain.Chat, error) {
	repoCtx, cancel := EnsureCtxTimeout(ctx, ctxTimeout)
	defer cancel()

	chat := domain.NewChat(title)
	err := r.сlient.WithContext(repoCtx).Create(&chat).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return &chat, nil
}

// GetChat retrieves chat by ID. Returns error if not found.
func (r *PostgresRepository) GetChat(ctx context.Context, chatID int64) (*domain.Chat, error) {
	repoCtx, cancel := EnsureCtxTimeout(ctx, ctxTimeout)
	defer cancel()

	var chat domain.Chat
	err := r.сlient.WithContext(repoCtx).First(&chat, chatID).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return &chat, nil
}

// SaveMessage persists new message and returns it with generated ID & CreatedAt field.
// func (r *PostgresRepository) SaveMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
func (r *PostgresRepository) SaveMessage(ctx context.Context, chatID int64, text string) (*domain.Message, error) {
	repoCtx, cancel := EnsureCtxTimeout(ctx, ctxTimeout)
	defer cancel()

	message := domain.NewMessage(chatID, text)
	err := r.сlient.WithContext(repoCtx).Create(&message).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return &message, nil
}

// GetMessages retrieves messages for chat, ordered by creation time (newest first).
func (r *PostgresRepository) GetMessages(ctx context.Context, chatID int64, limit int) ([]domain.Message, error) {
	repoCtx, cancel := EnsureCtxTimeout(ctx, ctxTimeout)
	defer cancel()

	var messages []domain.Message
	err := r.сlient.WithContext(repoCtx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return messages, nil
}

// handleError wraps database errors into domain-specific errors.
func (r *PostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	customErr := postgresclient.ErrorWrapper(err)
	return ErrorFactory(customErr)
}
