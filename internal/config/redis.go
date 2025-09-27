package config

import (
	"go-echo-template/internal/shared/utils"
	"strconv"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (r *RedisConfig) Addr() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}

func newRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     utils.MustGetStrEnv("REDIS_HOST"),
		Port:     utils.MustGetIntEnv("REDIS_PORT"),
		Password: utils.MustGetStrEnv("REDIS_PASSWORD"),
		DB:       utils.MustGetIntEnv("REDIS_DB"),
	}
}
