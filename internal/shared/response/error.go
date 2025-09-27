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
	Code     string
	Status   int
	Messages i18n.Messages
	Args     []any
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

func (ce *CustomErr) translate(locale i18n.Locale) string {
	// Check if locale exists, else fallback
	if msg, ok := ce.Messages[locale]; ok {
		if len(ce.Args) > 0 {
			return fmt.Sprintf(msg, ce.Args...)
		}
		return msg
	}

	// fallback locale
	if msg, ok := ce.Messages[i18n.DefaultLocale]; ok {
		if len(ce.Args) > 0 {
			return fmt.Sprintf(msg, ce.Args...)
		}
		return msg
	}

	// ultimate fallback: return code
	return ce.Code
}

// Echo Error Handler
func CustomHTTPErrorHandler(err error, c echo.Context) {
	locale := i18n.GetLocaleFromContext(c)

	switch err := err.(type) {
	// 1) Handle validation errors
	case validator.ValidationErrors:
		var fieldErrs []CustomFieldErr
		for _, fe := range err {
			fieldKey := fmt.Sprintf("FIELD:%s", strings.ToUpper(fe.Field()))
			translatedField := i18n.Translate(fieldKey, locale)
			userInput := fmt.Sprintf("%v", fe.Value())
			field := strings.ToLower(fe.Field()[:1]) + fe.Field()[1:]

			valKey, args := TagHandler(fe, translatedField)
			msg := i18n.Translate(valKey, locale, args...)
			fieldErrs = append(fieldErrs, CustomFieldErr{
				Input:   userInput,
				Field:   field,
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

		c.JSON(resp.Status, resp)
		return

	// 2) Handle custom errors
	case *CustomErr:
		resp := errResponse{
			IsError: true,
			Code:    err.Code,
			Status:  err.Status,
			Message: err.translate(locale),
		}

		c.JSON(err.Status, resp)
		return

	// 3) Handle Echo's HTTP errors
	case *echo.HTTPError:
		code := fmt.Sprintf("ERR:HTTP_%d", err.Code)
		resp := errResponse{
			IsError: true,
			Code:    code,
			Status:  err.Code,
			Message: i18n.Translate(code, locale),
		}

		c.JSON(resp.Status, resp)
		return

	// 4) Fallback to internal server error
	default:
		resp := errResponse{
			IsError: true,
			Code:    "ERR:HTTP_500",
			Status:  http.StatusInternalServerError,
			Message: i18n.Translate("ERR:HTTP_500", locale),
		}

		c.JSON(resp.Status, resp)
	}
}
