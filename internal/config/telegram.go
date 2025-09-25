package config

import "go-echo-template/internal/shared/utils"

type AlarmerConfig struct {
	Telegram *TelegramConfig
}

type TelegramConfig struct {
	CHAT_ID   string
	BOT_TOKEN string
}

func newTelegramConfig() *TelegramConfig {
	return &TelegramConfig{
		CHAT_ID:   utils.MustGetStrEnv("TELEGRAM_CHAT_ID"),
		BOT_TOKEN: utils.MustGetStrEnv("TELEGRAM_BOT_TOKEN"),
	}
}

func newAlarmerConfig() *AlarmerConfig {
	return &AlarmerConfig{
		Telegram: newTelegramConfig(),
	}
}
