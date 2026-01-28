package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	ctxTimeout = 5 * time.Second
)

type PostgresClient interface {
	WithContext(ctx context.Context) *gorm.DB
}

type PostgresRepository struct {
	Client PostgresClient
}

func NewPostgresRepository(client PostgresClient) *PostgresRepository {
	return &PostgresRepository{
		Client: client,
	}
}
