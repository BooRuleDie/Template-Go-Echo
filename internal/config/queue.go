package config

type QueueConfig struct {
	Redis *RedisConfig
}

func newQueueConfig() *QueueConfig {
	return &QueueConfig{
		Redis: newRedisConfig(),
	}
}
