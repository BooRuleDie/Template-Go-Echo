package config

type Config struct {
	DB     *DBConfig
	JWT    *JWTConfig
	Server *ServerConfig
	Redis  *RedisConfig
	Mail   *MailConfig
	SMS    *SMSConfig
	Push   *PushConfig
	Object *ObjectConfig
	Queue  *QueueConfig
}

func Load() *Config {
	return &Config{
		Server: newServerConfig(),
		DB:     newDBConfig(),
		Redis:  newRedisConfig(),
		JWT:    newJWTConfig(),
		Mail:   newMailConfig(),
		SMS:    newSMSConfig(),
		Push:   newPushConfig(),
		Object: newObjectConfig(),
		Queue:  newQueueConfig(),
	}
}
