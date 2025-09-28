package log

import (
	"context"
	"time"

	"go-echo-template/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements the Logger interface using Zap
type zapLogger struct {
	logger *zap.Logger
}

func NewCustomLogger(cfg *config.ServerConfig) (CustomLogger, error) {
	var zapCfg zap.Config

	if cfg.IsLocal() {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	// Add initial fields
	zapCfg.InitialFields = map[string]any{
		"appName": cfg.AppName,
		"version": cfg.Version,
		"env":     cfg.Environment,
	}

	zapLog, err := zapCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	logger := &zapLogger{logger: zapLog}
	return logger, nil
}

// zapLogger implementation
func (z *zapLogger) convertFields(fields []field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		switch v := field.Value.(type) {
		case string:
			zapFields[i] = zap.String(field.Key, v)
		case int:
			zapFields[i] = zap.Int(field.Key, v)
		case bool:
			zapFields[i] = zap.Bool(field.Key, v)
		case float64:
			zapFields[i] = zap.Float64(field.Key, v)
		case time.Time:
			zapFields[i] = zap.Time(field.Key, v)
		case error:
			// Using NamedError allows us to respect the key from the Field struct.
			// If you use log.Err(err), the key is "error".
			// If you use log.Any("db_error", err), the key will be "db_error".
			zapFields[i] = zap.NamedError(field.Key, v)
		default:
			// Fallback for any other type
			zapFields[i] = zap.Any(field.Key, v)
		}
	}
	return zapFields
}

func (z *zapLogger) Info(msg string, fields ...field) {
	z.logger.Info(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Warn(msg string, fields ...field) {
	z.logger.Warn(msg, z.convertFields(fields)...)
}

func (z *zapLogger) Error(msg string, fields ...field) {
	z.logger.Error(msg, z.convertFields(fields)...)
}

func (z *zapLogger) InfoWithContext(ctx context.Context, msg string, fields ...field) {
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
			z.Info(msg, append(fields, z.String(RequestIDKeyStr, requestID))...)
			return
		}
	}
	z.Info(msg, fields...)
}

func (z *zapLogger) WarnWithContext(ctx context.Context, msg string, fields ...field) {
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
			z.Warn(msg, append(fields, z.String(RequestIDKeyStr, requestID))...)
			return
		}
	}
	z.Warn(msg, fields...)
}

func (z *zapLogger) ErrorWithContext(ctx context.Context, msg string, fields ...field) {
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
			z.Error(msg, append(fields, z.String(RequestIDKeyStr, requestID))...)
			return
		}
	}
	z.Error(msg, fields...)
}

func (z *zapLogger) With(fields ...field) CustomLogger {
	return &zapLogger{
		logger: z.logger.With(z.convertFields(fields)...),
	}
}

func (z *zapLogger) String(key, value string) field {
	return field{Key: key, Value: value}
}

func (z *zapLogger) Int(key string, value int) field {
	return field{Key: key, Value: value}
}

func (z *zapLogger) Float64(key string, value float64) field {
	return field{Key: key, Value: value}
}

func (z *zapLogger) Bool(key string, value bool) field {
	return field{Key: key, Value: value}
}

func (z *zapLogger) Err(err error) field {
	// By convention, error fields are keyed with "error".
	return field{Key: "error", Value: err}
}

func (z *zapLogger) Any(key string, value any) field {
	return field{Key: key, Value: value}
}

func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}
