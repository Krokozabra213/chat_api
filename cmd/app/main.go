package main

import (
	"fmt"

	"github.com/Krokozabra213/test_api/internal/config"
	"github.com/Krokozabra213/test_api/pkg/logger"
)

func main() {
	cfg, err := config.Init("configs/main.yml", ".env")
	if err != nil {
		panic(err)
	}

	log := logger.Init(cfg.App.Environment)
	log.Debug("config loaded", "config", fmt.Sprintf("%+v", cfg))
	// pgConfig := postgresclient.NewPGConfig("localhost", "5432", "postgres", "secret", "mydb", "disable",
	// 	25, 5, 5*time.Minute)
	// db, err := postgresclient.New(pgConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Connected to PostgreSQL")

	// // Database Shutdown
	// if err := db.Shutdown(); err != nil {
	// 	log.Printf("shutdown error: %v", err)
	// }
}
