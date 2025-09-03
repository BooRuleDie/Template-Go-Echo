package config

import (
	"go-echo-template/internal/shared/utils"
	"time"
)

type ServerConfig struct {
	Port         int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         utils.MustGetIntEnv("SERVER_PORT"),
		Environment:  utils.MustGetStrEnv("ENVIRONMENT"),
		ReadTimeout:  utils.MustGetDurationEnv("SERVER_READ_TIMEOUT"),
		WriteTimeout: utils.MustGetDurationEnv("SERVER_WRITE_TIMEOUT"),
		IdleTimeout:  utils.MustGetDurationEnv("SERVER_IDLE_TIMEOUT"),
	}
}

func (s *ServerConfig) IsDevelopment() bool {
	return s.Environment == "dev"
}

func (s *ServerConfig) IsProduction() bool {
	return s.Environment == "prod"
}
