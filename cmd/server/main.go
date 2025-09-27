package main

import (
	"context"

	"go-echo-template/internal/alarm"
	"go-echo-template/internal/cache"
	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
	"go-echo-template/internal/modules/user"
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create the Background Context
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Initiate Custom Logger
	logger, err := log.NewCustomLogger(cfg.Server)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Initiate Alarmer
	alarmer := alarm.NewAlarmer(cfg.Alarmer.Telegram, logger)

	// Create Echo instance
	e := echo.New()

	// Use the custom validator
	e.Validator = response.NewValidator()

	// Set custom error handler
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	// Use global middlewares
	e.Use(middleware.RequestID())
	e.Use(log.RequestIDContextMiddleware())
	e.Use(log.LoggerMiddleware(logger))
	e.Use(i18n.LocaleMiddleware)
	e.Use(middleware.ContextTimeout(cfg.Server.RequestTimeout))
	// e.Use(middleware.Recover())

	// Connect to the PostgreSQL DB
	postgreSQL, err := db.NewPostgreSQL(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer postgreSQL.Close()

	// Connect to the Redis Cache
	redis := cache.NewRedisCache(ctx, *cfg.Redis)
	defer redis.Close()

	// Register user routes
	userCache := user.NewUserCache(redis)
	userRepo := user.NewUserRepository(logger, postgreSQL, userCache)
	userService := user.NewUserService(logger, userRepo)
	user.NewUserHandler(logger, alarmer, userService).RegisterRoutes(e)

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
