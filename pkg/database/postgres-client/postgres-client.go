package postgresclient

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrFailedConnect = errors.New("failed to connect to database")
	ErrFailedSQLDB   = errors.New("failed to get sql.DB")
	ErrTimeout       = errors.New("timeout occurred")
)

type PostgresClient struct {
	*gorm.DB
}

func New(cfg PGConfig) (*PostgresClient, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbName, cfg.sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedConnect, err)
	}

	// Настройка connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedSQLDB, err)
	}

	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.connMaxLifetime)

	return &PostgresClient{db}, nil
}

// Ping проверка соединения
func (p *PostgresClient) Ping(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Shutdown корректно закрывает соединение с таймаутом
func (p *PostgresClient) Shutdown(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedSQLDB, err)
	}

	// Канал для результата закрытия
	done := make(chan error, 1)
	go func() {
		done <- sqlDB.Close()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrTimeout, ctx.Err())
	case err := <-done:
		return err
	}
}
