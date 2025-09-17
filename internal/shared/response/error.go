package response

import (
	"errors"
	"fmt"
	"go-echo-template/internal/shared/i18n"
	"go-echo-template/internal/shared/validation"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Standard Error Response
type errResponse struct {
	IsError          bool                        `json:"isError"`
	Code             string                      `json:"code"`
	Status           int                         `json:"status"`
	Message          string                      `json:"message"`
	ValidationErrors []validation.CustomFieldErr `json:"validationErrors,omitempty"`
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
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		var fieldErrs []validation.CustomFieldErr
		for _, fe := range valErrs {
			jsonField := strings.ToLower(fe.Field())
			fieldKey := fmt.Sprintf("FIELD:%s", strings.ToUpper(fe.Field()))
			translatedField := i18n.Translate(fieldKey, locale)
			userInput := fmt.Sprintf("%v", fe.Value())

			if handler, ok := validation.TagHandlers[fe.Tag()]; ok {
				code, args := handler(fe, translatedField)
				msg := i18n.Translate(code, locale, args...)
				fieldErrs = append(fieldErrs, validation.CustomFieldErr{
					Input:   userInput,
					Field:   jsonField,
					Message: msg,
				})
			} else {
				// Fallback for unhandled validation tags
				fallbackMsg := fmt.Sprintf("i18n message translation failed for this unhandled tag: %s", fe.Tag())
				fieldErrs = append(fieldErrs, validation.CustomFieldErr{
					Input:   userInput,
					Field:   jsonField,
					Message: fallbackMsg,
				})
			}
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
