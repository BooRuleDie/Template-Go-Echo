package user

import (
	"go-echo-template/internal/shared/response"
	"net/http"
)

var (
	errUserNotFound = &response.CustomErr{
		Code:   "ERR:USER_NOT_FOUND",
		Status: http.StatusNotFound,
	}
	errUserAlreadyExists = &response.CustomErr{
		Code:   "ERR:USER_ALREADY_EXISTS",
		Status: http.StatusConflict,
	}
	errInvalidID = &response.CustomErr{
		Code:   "ERR:USER_INVALID_ID",
		Status: http.StatusBadRequest,
	}
	errInvalidRequestPayload = &response.CustomErr{
		Code:   "ERR:USER_INVALID_REQUEST_PAYLOAD",
		Status: http.StatusBadRequest,
	}
)
