package response

import (
	"fmt"
	"go-echo-template/internal/shared/i18n"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Standard Error Response
type errResponse struct {
	IsError          bool             `json:"isError"`
	Code             string           `json:"code"`
	Status           int              `json:"status"`
	Message          string           `json:"message"`
	ValidationErrors []CustomFieldErr `json:"validationErrors,omitempty"`
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

	// 1) Handle validation errors
	if valErrs, ok := err.(validator.ValidationErrors); ok {
		var fieldErrs []CustomFieldErr
		for _, fe := range valErrs {
			fieldKey := fmt.Sprintf("FIELD:%s", strings.ToUpper(fe.Field()))
			translatedField := i18n.Translate(fieldKey, locale)
			userInput := fmt.Sprintf("%v", fe.Value())
			jsonField := strings.ToLower(fe.Field())

			valKey, args := TagHandler(fe, translatedField)
			msg := i18n.Translate(valKey, locale, args...)
			fieldErrs = append(fieldErrs, CustomFieldErr{
				Input:   userInput,
				Field:   jsonField,
				Message: msg,
			})
		}

		resp := errResponse{
			IsError:          true,
			Code:             "VAL:VALIDATION_ERR",
			Status:           http.StatusUnprocessableEntity,
			Message:          i18n.Translate("VAL:VALIDATION_ERR", locale),
			ValidationErrors: fieldErrs,
		}
		// TODO: log the error after log implementation
		c.JSON(resp.Status, resp)
		return
	}

	// 2) Handle custom errors
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

	// 3) Fallback to internal server error
	resp := errResponse{
		IsError: true,
		Code:    "ERR:INTERNAL_SERVER_ERROR",
		Status:  http.StatusInternalServerError,
		Message: i18n.Translate("ERR:INTERNAL_SERVER_ERROR", locale),
	}
	// TODO: log the error after log implementation
	c.JSON(resp.Status, resp)
}
