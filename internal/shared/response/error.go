package response

import (
	"go-echo-template/internal/shared/i18n"

	"github.com/labstack/echo/v4"
)

// Standard Error Response
type errResponse struct {
	IsError bool   `json:"isError"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Custom Error
type CustomErr struct {
	Code   string
	Status int
	Args   []any
}

// Satisfies the error interface
func (ce *CustomErr) Error() string {
	return ce.Code
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
			Status:  ce.Status,
			Message: i18n.Translate(ce.Code, locale, ce.Args...),
		}
		// TODO: log the error after log implementation
		c.JSON(ce.Status, resp)
		return
	}

	// fallback
	resp := errResponse{
		IsError: true,
		Code:    "ERR:INTERNAL_SERVER_ERROR",
		Status:  500,
		Message: i18n.Translate("ERR:INTERNAL_SERVER_ERROR", locale),
	}
	// TODO: log the error after log implementation
	c.JSON(500, resp)
}
