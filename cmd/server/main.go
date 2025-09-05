package main

import (
	"github.com/labstack/echo/v4"

	"go-echo-template/internal/config"
	"go-echo-template/internal/modules/user"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Echo instance
	e := echo.New()

	// Register user routes
	user.NewUserHandler().RegisterRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
