package utils

import (
	"os"
	"strconv"
	"time"
)

func GetStrEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func MustGetStrEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	panic("missing required environment variable: " + key)
}

func GetIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic("invalid integer value for environment variable: " + key)
		}

		return intValue
	}

	return defaultValue
}

func MustGetIntEnv(key string) int {
	value := os.Getenv(key)
	if value == "" {
		panic("missing required environment variable: " + key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("invalid integer value for environment variable: " + key)
	}

	return intValue
}

func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err != nil {
			panic("invalid duration value for environment variable: " + key)
		}

		return duration
	}
	return defaultValue
}

func MustGetDurationEnv(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		panic("missing required environment variable: " + key)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic("invalid duration value for environment variable: " + key)
	}

	return duration
}
