package shared

import (
	"net/http"

	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
)

var (
	ErrInvalidRequestPayload = &response.CustomErr{
		Status: http.StatusBadRequest,
		Code:   "ERR:USER_INVALID_REQUEST_PAYLOAD",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Invalid request payload",
			i18n.TR_TR: "Geçersiz istek verisi",
		},
	}

	ErrUserNotFound = &response.CustomErr{
		Status: http.StatusNotFound,
		Code:   "ERR:USER_NOT_FOUND",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "User not found",
			i18n.TR_TR: "Kullanıcı bulunamadı",
		},
	}

	ErrSessionUnauthorized = &response.CustomErr{
		Status: http.StatusUnauthorized,
		Code:   "ERR:SESSION_UNAUTHORIZED",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Unauthorized access",
			i18n.TR_TR: "Yetkisiz erişim",
		},
	}
)
