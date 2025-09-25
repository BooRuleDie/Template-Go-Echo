package user

import (
	"net/http"

	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
)

var (
	errUserNotFound = &response.CustomErr{
		Status: http.StatusNotFound,
		Code:   "ERR:USER_NOT_FOUND",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "User with ID %v not found",
			i18n.TR_TR: "%v ID'li kullanıcı bulunamadı",
		},
	}
	errUserAlreadyExists = &response.CustomErr{
		Status: http.StatusConflict,
		Code:   "ERR:USER_ALREADY_EXISTS",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "User already exists",
			i18n.TR_TR: "Kullanıcı kaydı yapılmış",
		},
	}
	errInvalidID = &response.CustomErr{
		Status: http.StatusBadRequest,
		Code:   "ERR:USER_INVALID_ID",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "User with ID %v not found",
			i18n.TR_TR: "ID'si %v olan kullanıcı bulunamadı",
		},
	}
	errInvalidRequestPayload = &response.CustomErr{
		Status: http.StatusBadRequest,
		Code:   "ERR:USER_INVALID_REQUEST_PAYLOAD",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "Invalid request payload",
			i18n.TR_TR: "Geçersiz istek verisi",
		},
	}
)
