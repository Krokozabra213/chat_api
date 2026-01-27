package main

import (
	"log"
	"time"

	postgresclient "github.com/Krokozabra213/test_api/pkg/database/postgres-client"
)

func main() {
	pgConfig := postgresclient.NewPGConfig("localhost", "5432", "postgres", "secret", "mydb", "disable",
		25, 5, 5*time.Minute)
	db, err := postgresclient.New(pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	// Database Shutdown
	if err := db.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
