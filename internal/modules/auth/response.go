package auth

import (
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
	"net/http"
)

// Success Messages
var (
	succLogin = &response.SuccessMessage{
		Code: "SUCC:LOGIN_SUCCESS",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Login successful",
			i18n.TR_TR: "Giriş başarılı",
		},
	}
)

// Error Messages
var (
	errEmailOrPasswordWrong = &response.CustomErr{
		Status: http.StatusUnauthorized,
		Code:   "ERR:EMAIL_OR_PASSWORD_WRONG",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Email or password is wrong",
			i18n.TR_TR: "E-posta veya şifre yanlış",
		},
	}
)
