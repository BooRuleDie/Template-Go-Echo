package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
	"go-echo-template/internal/modules/user"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Echo instance
	e := echo.New()

	// Connect to the PostgreSQL DB
	DB, err := db.New(cfg.DB)
	if err != nil {
		e.Logger.Fatalf("failed to initialize database: %v", err)
	}

	// Register user routes
	user.NewUserHandler(DB).RegisterRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
