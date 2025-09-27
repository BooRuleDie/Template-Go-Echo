package log

import (
	"context"
)

// Logger defines the interface for structured logging
type CustomLogger interface {
	Info(msg string, fields ...field)
	Warn(msg string, fields ...field)
	Error(msg string, fields ...field)

	InfoWithContext(ctx context.Context, msg string, fields ...field)
	WarnWithContext(ctx context.Context, msg string, fields ...field)
	ErrorWithContext(ctx context.Context, msg string, fields ...field)

	With(fields ...field) CustomLogger
	Sync() error

	// Field helper functions
	String(key, value string) field
	Int(key string, value int) field
	Float64(key string, value float64) field
	Bool(key string, value bool) field
	Err(err error) field
	Any(key string, value any) field
}

// Field represents a structured log field
type field struct {
	Key   string
	Value any
}
