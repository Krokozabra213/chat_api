package suite

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Krokozabra213/test_api/internal/config"
	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
)

const (
	ctxTimeout       = 30 * time.Second
	localHTTPAddress = "http://localhost:8180"
)

type APISuite struct {
	*testing.T
	Config     *config.Config
	DB         *postgresclient.PostgresClient
	HTTPClient *Client
}

func New(t *testing.T) (context.Context, *APISuite) {
	t.Helper()
	configFile := filepath.Join("..", "..", "configs", "main.yml")
	envFile := filepath.Join("..", "..", ".env")

	cfg, err := config.Init(configFile, envFile)
	if err != nil {
		t.Fatalf("config init err: %v", err)
	}

	pgConfig := postgresclient.NewPGConfig(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.MaxOpenConns,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.ConnMaxLifetime,
	)

	db, err := postgresclient.New(pgConfig)
	if err != nil {
		t.Fatalf("db connect err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})
	client := NewClient(localHTTPAddress, nil)

	return ctx, &APISuite{
		T:          t,
		Config:     cfg,
		DB:         db,
		HTTPClient: client,
	}
}

func (s *APISuite) CleanupTestData() error {
	err := s.CleanupMessages()
	if err != nil {
		return err
	}

	return s.CleanupChats()
}

func (s *APISuite) CleanupMessages() error {
	result := s.DB.Exec("TRUNCATE TABLE messages CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *APISuite) CleanupChats() error {
	result := s.DB.Exec("TRUNCATE TABLE chats CASCADE")
	if result.Error != nil {
		return result.Error
	}
	return nil
}
