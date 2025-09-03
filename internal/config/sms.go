package config

import "go-echo-template/internal/shared/utils"

type OneSignalConfig struct {
	APIKey string
}

type SMSSendGridConfig struct {
	APIKey string
}

type SMSConfig struct {
	OneSignal *OneSignalConfig
	SendGrid  *SMSSendGridConfig
}

func newOneSignalConfig() *OneSignalConfig {
	return &OneSignalConfig{
		APIKey: utils.MustGetStrEnv("ONESIGNAL_API_KEY"),
	}
}

func newSMSSendgridConfig() *SMSSendGridConfig {
	return &SMSSendGridConfig{
		APIKey: utils.MustGetStrEnv("SENDGRID_API_KEY"),
	}
}

func newSMSConfig() *SMSConfig {
	return &SMSConfig{
		OneSignal: newOneSignalConfig(),
		SendGrid:  newSMSSendgridConfig(),
	}
}
