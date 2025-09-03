package config

import "go-echo-template/internal/shared/utils"

type MailConfig struct {
	SMTP     *SMTPConfig
	SendGrid *EmailSendGridConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailSendGridConfig struct {
	APIKey string
}

func newSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     utils.MustGetStrEnv("SMTP_HOST"),
		Port:     utils.MustGetIntEnv("SMTP_PORT"),
		Username: utils.MustGetStrEnv("SMTP_USERNAME"),
		Password: utils.MustGetStrEnv("SMTP_PASSWORD"),
	}
}

func newEmailSendgridConfig() *EmailSendGridConfig {
	return &EmailSendGridConfig{
		APIKey: utils.MustGetStrEnv("SENDGRID_API_KEY"),
	}
}

func newMailConfig() *MailConfig {
	return &MailConfig{
		SMTP:     newSMTPConfig(),
		SendGrid: newEmailSendgridConfig(),
	}
}
