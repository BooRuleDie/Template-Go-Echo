package cache

import (
	"context"
	"go-echo-template/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisCache(ctx context.Context, config config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr(),
		Password: config.Password,
		DB:       config.DB,
	})

	// Create a timeout context for the ping to Redis.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}

	return client
}
