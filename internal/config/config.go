package config

type Config struct {
	DB      *DBConfig
	Alarmer *AlarmerConfig
	Server  *ServerConfig
	Redis   *RedisConfig
	Mail    *MailConfig
	Object  *ObjectConfig
	Queue   *QueueConfig
}

func Load() *Config {
	return &Config{
		Server:  newServerConfig(),
		DB:      newDBConfig(),
		Redis:   newRedisConfig(),
		Alarmer: newAlarmerConfig(),
		Mail:    newMailConfig(),
		Object:  newObjectConfig(),
		Queue:   newQueueConfig(),
	}
}
