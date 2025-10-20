package log

import (
	"context"
	"net/http"

	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const RequestIDKey shared.ContextKey = "request_id"
const RequestIDKeyStr string = string(RequestIDKey)

// logger middleware
func LoggerMiddleware(logger CustomLogger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogProtocol:  true,
		LogURI:       true,
		LogMethod:    true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Determine status code (may override if error is present)
			if v.Error != nil {
				switch error := v.Error.(type) {
				case *response.CustomErr:
					v.Status = error.Status
				case validator.ValidationErrors:
					v.Status = http.StatusUnprocessableEntity
				default:
					v.Status = http.StatusInternalServerError
				}
			}

			fields := []field{
				logger.String("protocol", v.Protocol),
				logger.String("method", v.Method),
				logger.String("uri", v.URI),
				logger.Int("status", v.Status),
				logger.String("latency", v.Latency.String()),
				logger.String("remote_ip", v.RemoteIP),
				logger.String("user_agent", v.UserAgent),
			}

			if v.RequestID != "" {
				fields = append(fields, logger.String(RequestIDKeyStr, v.RequestID))
			}

			if v.Error != nil {
				fields = append(fields, logger.Err(v.Error))
				logger.Error("HTTP request error", fields...)
			} else {
				logger.Info("HTTP request", fields...)
			}

			return nil
		},
	})
}

// requestID middleware, it injects the requestID to the echo.Request.Context automatically
func RequestIDContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the request ID that was set by RequestID middleware
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)

			if requestID != "" {
				// Create new context with request_id and replace the original
				ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)

				// Replace the request with the new context
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}
