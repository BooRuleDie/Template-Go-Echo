package main

import (
	"context"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
	"go-echo-template/internal/modules/user"
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
)

func main() {
	// Create the Background Context
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Create Echo instance
	e := echo.New()

	// Use the custom validator
	e.Validator = response.NewValidator()

	// Use i18n locale middleware
	e.Use(i18n.LocaleMiddleware)

	// Set custom error handler
	e.HTTPErrorHandler = response.HTTPErrHandler

	// Connect to the PostgreSQL DB
	DB, err := db.NewPostgreSQL(ctx, cfg.DB)
	if err != nil {
		e.Logger.Fatalf("failed to initialize database: %v", err)
	}
	defer DB.Close()

	// Register user routes
	user.NewUserHandler(DB).RegisterRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
