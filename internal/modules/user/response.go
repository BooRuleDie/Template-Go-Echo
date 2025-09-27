package user

import (
	"net/http"

	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/response"
)

// Success Messages
var (
	succUserCreated = &response.SuccessMessage{
		Code: "SUCC:USER_CREATED",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "User created successfully",
			i18n.TR_TR: "Kullanıcı başarıyla oluşturuldu",
		},
	}
)

// Errors
var (
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
	errUserEmailAlreadyExists = &response.CustomErr{
		Status: http.StatusConflict,
		Code:   "ERR:USER_EMAIL_ALREADY_EXISTS",
		Messages: map[i18n.Locale]string{
			i18n.EN_US: "A user with that email already exists",
			i18n.TR_TR: "Bu e-posta adresine sahip bir kullanıcı zaten kayıtlı",
		},
	}
)
