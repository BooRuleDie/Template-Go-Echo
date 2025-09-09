package main

import (
	"log"

	_ "github.com/lib/pq"

	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to the PostgreSQL DB
	DB, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer DB.Close()
	
	// Start seeding

}
