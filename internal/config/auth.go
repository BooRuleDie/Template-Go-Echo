package config

import "go-echo-template/internal/shared/utils"

type AuthConfig struct {
	Session *SessionConfig
}

type SessionConfig struct {
	SESSION_SECRET string
}

func newSessionConfig() *SessionConfig {
	return &SessionConfig{
		SESSION_SECRET: utils.MustGetStrEnv("SESSION_SECRET"),
	}
}

func newAuthConfig() *AuthConfig {
	return &AuthConfig{
		Session: newSessionConfig(),
	}
}
