package config

import "go-echo-template/internal/shared/utils"

type PushConfig struct {
	APIKey string
}

func newPushConfig() *PushConfig {
	return &PushConfig{
		APIKey: utils.MustGetStrEnv("PUSH_API_KEY"),
	}
}
