package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/test_api/internal/config"
	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
	"github.com/Krokozabra213/test_api/pkg/logger"
)

const (
	configFile = "configs/main.yml"
	envFile    = ".env"
)

func main() {
	cfg, err := config.Init(configFile, envFile)
	if err != nil {
		panic(err)
	}

	slogger := logger.Init(cfg.App.Environment)
	slogger.Info("application started", "config", cfg)

	pgConfig := postgresclient.NewPGConfig(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName,
		cfg.Postgres.SSLMode, cfg.Postgres.MaxOpenConns, cfg.Postgres.MaxIdleConns, cfg.Postgres.ConnMaxLifetime)
	db, err := postgresclient.New(pgConfig)
	if err != nil {
		log.Fatal(err)
	}
	slogger.Info("connected to postgres")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	slogger.Info("received shutdown signal", "signal", sig)

	// Database Shutdown
	if err := db.Shutdown(); err != nil {
		slogger.Error("database shutdown error", "error", err)
		log.Printf("database shutdown error: %v", err)
	}
}
