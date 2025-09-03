package config

import (
	"go-echo-template/internal/shared/utils"
	"time"
)

type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
	RefreshTime    time.Duration
}

func newJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:         utils.MustGetStrEnv("JWT_SECRET"),
		ExpirationTime: utils.MustGetDurationEnv("JWT_EXPIRATION"),
		RefreshTime:    utils.MustGetDurationEnv("JWT_REFRESH"),
	}
}
