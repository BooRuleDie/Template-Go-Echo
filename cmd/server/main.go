package main

import (
	"github.com/labstack/echo/v4"

	"go-echo-template/internal/config"
	"go-echo-template/internal/modules/user"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up user repository, service, and handler
	userRepo := user.NewRepository()
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Create Echo instance
	e := echo.New()

	// Register user routes
	userHandler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
