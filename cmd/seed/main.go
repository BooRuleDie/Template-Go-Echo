package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
)

func main() {
	// Create the Background Context
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Connect to the PostgreSQL DB
	DB, err := db.NewPostgreSQL(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer DB.Close()

	// Start seeding

}
