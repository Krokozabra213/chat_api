package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krokozabra213/test_api/internal/business"
	"github.com/Krokozabra213/test_api/internal/config"
	handler "github.com/Krokozabra213/test_api/internal/delivery/http"
	"github.com/Krokozabra213/test_api/internal/repository/postgres"
	"github.com/Krokozabra213/test_api/internal/server"
	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
	"github.com/Krokozabra213/test_api/pkg/logger"
)

const (
	configFile      = "configs/main.yml"
	envFile         = ".env"
	shutdownTimeout = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Config
	cfg, err := config.Init(configFile, envFile)
	if err != nil {
		return err
	}

	if cfg.App.Environment == "local" {
		cfg.Postgres.Host = "localhost"
	}

	// Logger
	log := logger.Init(cfg.App.Environment)
	log.Info("initialized config", "config", cfg.LogValue())
	log.Info("starting application")

	// Database
	pgConfig := postgresclient.NewPGConfig(
		// cfg.Postgres.Host,
		"localhost",
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.MaxOpenConns,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.ConnMaxLifetime,
	)
	// fmt.Println(cfg.Postgres.Host)
	// fmt.Println(cfg.Postgres.Port)

	db, err := postgresclient.New(pgConfig)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Shutdown(shutdownTimeout); err != nil {
			log.Error("database shutdown error", "error", err)
		}
	}()
	log.Info("connected to postgres")

	// Dependencies
	repo := postgres.NewPostgresRepository(db)
	biz := business.New(log, repo, repo)

	// Router
	router := http.NewServeMux()
	handler.New(router, log, biz)

	// Server
	srv := server.NewServer(cfg, router)

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		log.Info("server started", "address", srv.Addr())
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	// Wait for shutdown signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-quit:
		log.Info("received shutdown signal", "signal", sig)
	case err := <-errCh:
		log.Error("server error", "error", err)
		return err
	}

	// Graceful shutdown
	if err := srv.ShutDown(shutdownTimeout); err != nil {
		log.Error("server shutdown error", "error", err)
		return err
	}

	log.Info("application stopped")
	return nil
}
