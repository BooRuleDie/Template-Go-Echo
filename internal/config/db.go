package config

import (
	"fmt"
	"go-echo-template/internal/shared/utils"
	"time"
)

type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func newDBConfig() *DBConfig {
	return &DBConfig{
		Host:         utils.MustGetStrEnv("DB_HOST"),
		Port:         utils.MustGetIntEnv("DB_PORT"),
		User:         utils.MustGetStrEnv("DB_USER"),
		Password:     utils.MustGetStrEnv("DB_PASSWORD"),
		DBName:       utils.MustGetStrEnv("DB_NAME"),
		SSLMode:      utils.MustGetStrEnv("DB_SSL_MODE"),
		MaxOpenConns: utils.MustGetIntEnv("DB_MAX_OPEN_CONNS"),
		MaxIdleConns: utils.MustGetIntEnv("DB_MAX_IDLE_CONNS"),
		MaxIdleTime:  utils.MustGetDurationEnv("DB_MAX_IDLE_TIME"),
	}
}

func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}
