package auth

import (
	"net/http"

	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
)

// Session-related errors
var (
	errSessionCookieNotFound = &response.CustomErr{
		Status: http.StatusUnauthorized,
		Code:   "ERR:SESSION_COOKIE_NOT_FOUND",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Session cookie not found",
			i18n.TR_TR: "Oturum çerezi bulunamadı",
		},
	}
	errEmptySessionID = &response.CustomErr{
		Status: http.StatusUnauthorized,
		Code:   "ERR:SESSION_EMPTY_ID",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Empty session ID",
			i18n.TR_TR: "Boş oturum ID'si",
		},
	}
	errSessionNotFound = &response.CustomErr{
		Status: http.StatusUnauthorized,
		Code:   "ERR:SESSION_NOT_FOUND",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Session expired or not found",
			i18n.TR_TR: "Oturum süresi dolmuş veya bulunamadı",
		},
	}
	errSessionGenID = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_GENERATE_ID",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to generate session ID",
			i18n.TR_TR: "Oturum ID'si oluşturulamadı",
		},
	}
	errSessionSerialize = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_SERIALIZATION",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to serialize user data",
			i18n.TR_TR: "Kullanıcı verisi serileştirilemedi",
		},
	}
	errSessionStore = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_STORE",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to store session in cache",
			i18n.TR_TR: "Oturum önbelleğe kaydedilemedi",
		},
	}
	errSessionCheckExist = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_CHECK_EXIST",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to check session existence",
			i18n.TR_TR: "Oturum varlığı kontrol edilemedi",
		},
	}
	errSessionRefresh = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_REFRESH",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to refresh session expiry",
			i18n.TR_TR: "Oturum süresi yenilenemedi",
		},
	}
	errSessionDeserialize = &response.CustomErr{
		Status: http.StatusInternalServerError,
		Code:   "ERR:SESSION_DESERIALIZE",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Failed to deserialize user data",
			i18n.TR_TR: "Kullanıcı verisi çözümlenemedi",
		},
	}
)
