package user

import (
	"go-echo-template/internal/shared/response"
	"net/http"
)

var (
	errUserNotFound = &response.CustomErr{
		Code:       "ERR:USER_NOT_FOUND",
		HTTPStatus: http.StatusNotFound,
	}
	errUserAlreadyExists = &response.CustomErr{
		Code:       "ERR:USER_ALREADY_EXISTS",
		HTTPStatus: http.StatusConflict,
	}
	errInvalidID = &response.CustomErr{
		Code:       "ERR:INVALID_ID",
		HTTPStatus: http.StatusBadRequest,
	}
	errInvalidRequestPayload = &response.CustomErr{
		Code:       "ERR:INVALID_REQUEST_PAYLOAD",
		HTTPStatus: http.StatusBadRequest,
	}
)
