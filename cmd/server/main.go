package main

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go-echo-template/internal/alarm"
	"go-echo-template/internal/cache"
	"go-echo-template/internal/config"
	"go-echo-template/internal/db"
	"go-echo-template/internal/modules/auth"
	"go-echo-template/internal/modules/user"
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"
	"go-echo-template/internal/storage"
	storageAuth "go-echo-template/internal/storage/auth"
	storageUser "go-echo-template/internal/storage/user"
	"go-echo-template/web"

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
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: cfg.Server.RequestTimeout,
		Skipper: func(c echo.Context) bool {
			// Skip timeout middleware for WebSocket upgrade requests
			// This is implemented to prevent local HMR websocket errors for web
			return strings.ToLower(c.Request().Header.Get("Connection")) == "upgrade" &&
				strings.ToLower(c.Request().Header.Get("Upgrade")) == "websocket"
		},
	}))
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

	// API grouping
	api := e.Group("/api")

	// New Storage Dependencies
	authRepo := storageAuth.NewAuthRepository(logger, postgreSQL)
	userCache := storageUser.NewUserCache(redis)
	userRepo := storageUser.NewUserRepository(logger, postgreSQL, userCache)

	// New Storage
	newStorage := storage.NewStorage(userRepo, authRepo)

	// Auth
	authService := auth.NewSessionCookieService(cfg.Server, logger, redis, newStorage)
	auth.NewAuthHandler(logger, alarmer, authService).RegisterRoutes(api)

	// User
	userService := user.NewUserService(logger, newStorage, authService)
	user.NewUserHandler(logger, alarmer, userService, authService).RegisterRoutes(api)

	// Register web route
	if cfg.Server.IsLocal() {
		target, _ := url.Parse(cfg.Server.LocalWebURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Catch-all route for frontend (but not /api)
		e.Any("/*", func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			proxy.ServeHTTP(res, req)
			return nil
		})
	} else {
		// Create a sub filesystem for the dist directory
		distFS, err := fs.Sub(web.WebFS, "dist")
		if err != nil {
			panic(err)
		}

		// Serve static files from the embedded filesystem
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "/",
			Index:      "index.html",
			HTML5:      true, // This enables SPA routing - serves index.html for non-existent routes
			Filesystem: http.FS(distFS),
			Skipper: func(c echo.Context) bool {
				// Skip static middleware for API routes
				return strings.HasPrefix(c.Request().URL.Path, "/api")
			},
		}))
	}

	e.Logger.Fatal(e.Start(cfg.Server.Address))
}
