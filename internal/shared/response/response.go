package response

import (
	"go-echo-template/internal/shared/i18n"

	"github.com/labstack/echo/v4"
)

// Standard Success Response
type successResponse struct {
	IsError bool `json:"isError"`
	Data    any  `json:"data,omitempty"`
}

func NewSuccessResponse(c echo.Context, status int, data any) error {
	resp := successResponse{
		IsError: false,
		Data:    data,
	}
	return c.JSON(status, resp)
}

// Standard Error Response
type errResponse struct {
	IsError bool   `json:"isError"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Custom Error
type CustomErr struct {
	Code       string
	HTTPStatus int
	Args       []any
}

// Satisfies the error interface
func (ce *CustomErr) Error() string {
	return ce.Code
}

// Format localized error message
func (ce *CustomErr) GetMessage(locale i18n.Locale) string {
	return i18n.Translate(ce.Code, locale, ce.Args...)
}

// Populate args for message with dynamic values
func (ce *CustomErr) WithArgs(args ...any) *CustomErr {
	ce.Args = args
	return ce
}

// Echo Error Handler
func HTTPErrHandler(err error, c echo.Context) {
	locale := i18n.GetLocaleFromContext(c)

	if ce, ok := err.(*CustomErr); ok {
		resp := errResponse{
			IsError: true,
			Code:    ce.Code,
			Message: ce.GetMessage(locale),
		}
		// TODO: log the error after log implementation
		c.JSON(ce.HTTPStatus, resp)
		return
	}

	// fallback
	resp := errResponse{
		IsError: true,
		Code:    "ERR:INTERNAL_SERVER_ERROR",
		Message: i18n.Translate("ERR:INTERNAL_SERVER_ERROR", locale),
	}
	// TODO: log the error after log implementation
	c.JSON(500, resp)
}
