package config

import (
	"time"

	"go-echo-template/internal/shared/utils"
)

type ServerConfig struct {
	AppName        string
	Version        string
	Address        string
	Port           int
	Environment    string
	RequestTimeout time.Duration
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		AppName:        utils.MustGetStrEnv("APP_NAME"),
		Version:        utils.MustGetStrEnv("VERSION"),
		Address:        utils.MustGetStrEnv("SERVER_ADDRESS"),
		Port:           utils.MustGetIntEnv("SERVER_PORT"),
		Environment:    utils.MustGetStrEnv("ENVIRONMENT"),
		RequestTimeout: utils.MustGetDurationEnv("REQUEST_TIMEOUT"),
	}
}

func (s *ServerConfig) IsDevelopment() bool {
	return s.Environment == "dev"
}

func (s *ServerConfig) IsProduction() bool {
	return s.Environment == "prod"
}

func (s *ServerConfig) IsLocal() bool {
	return s.Environment == "local"
}
