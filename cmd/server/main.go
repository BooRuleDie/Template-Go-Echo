package main

import (
	"context"

	"go-echo-template/internal/alarm"
	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
	"go-echo-template/internal/modules/user"
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create the Background Context
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Initiate Global Alarmer
	alarm.SetGlobalAlarmer(cfg.Alarmer.Telegram)

	// Create Echo instance
	e := echo.New()

	// Use the custom validator
	e.Validator = response.NewValidator()

	// Use global middlewares
	e.Use(i18n.LocaleMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set custom error handler
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

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
