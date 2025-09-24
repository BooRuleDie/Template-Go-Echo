package response

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

type CustomFieldErr struct {
	Input   string `json:"input"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

func NewValidator() *CustomValidator {
	// Enable reporting of all validation errors (not just the first one)
	v := validator.New()

	// Add custom "phone" validation for Turkish GSM numbers
	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		val := fl.Field()
		if val.Kind() != reflect.String {
			return false
		}
		s := val.String()
		// Must start with +905 followed by 9 digits (Turkish GSM: +905XXXXXXXXX)
		if len(s) != 13 {
			return false
		}
		if s[:4] != "+905" {
			return false
		}
		for _, c := range s[4:] {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	})

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func TagHandler(fe validator.FieldError, fieldName string) (string, []any) {
	switch fe.Tag() {
	case "required":
		return "VAL:REQUIRED", []any{fieldName}
	case "email":
		return "VAL:EMAIL", []any{fieldName}
	case "phone":
		return "VAL:PHONE", []any{fieldName}
	case "alpha":
		return "VAL:ALPHA", []any{fieldName}
	case "alphanum":
		return "VAL:ALPHANUM", []any{fieldName}
	case "contains":
		return "VAL:CONTAINS", []any{fieldName, fe.Param()}
	case "oneof":
		return "VAL:ONEOF", []any{fieldName, fe.Param()}
	case "gte":
		switch fe.Kind() {
		case reflect.String:
			return "VAL:MIN_STRING", []any{fieldName, fe.Param()}
		case reflect.Slice, reflect.Array:
			return "VAL:MIN_SLICE", []any{fieldName, fe.Param()}
		default:
			return "VAL:GTE_NUMBER", []any{fieldName, fe.Param()}
		}
	case "lte":
		switch fe.Kind() {
		case reflect.String:
			return "VAL:MAX_STRING", []any{fieldName, fe.Param()}
		case reflect.Slice, reflect.Array:
			return "VAL:MAX_SLICE", []any{fieldName, fe.Param()}
		default:
			return "VAL:LTE_NUMBER", []any{fieldName, fe.Param()}
		}
	case "min":
		switch fe.Kind() {
		case reflect.String:
			return "VAL:MIN_STRING", []any{fieldName, fe.Param()}
		case reflect.Slice, reflect.Array:
			return "VAL:MIN_SLICE", []any{fieldName, fe.Param()}
		default:
			return "VAL:MIN_NUMBER", []any{fieldName, fe.Param()}
		}
	case "max":
		switch fe.Kind() {
		case reflect.String:
			return "VAL:MAX_STRING", []any{fieldName, fe.Param()}
		case reflect.Slice, reflect.Array:
			return "VAL:MAX_SLICE", []any{fieldName, fe.Param()}
		default:
			return "VAL:MAX_NUMBER", []any{fieldName, fe.Param()}
		}
	case "len":
		switch fe.Kind() {
		case reflect.String:
			return "VAL:LEN_STRING", []any{fieldName, fe.Param()}
		case reflect.Slice, reflect.Array:
			return "VAL:LEN_SLICE", []any{fieldName, fe.Param()}
		default:
			return "VAL:LEN_NUMBER", []any{fieldName, fe.Param()}
		}
	default:
		return "VAL:UNKNOWN", []any{fieldName, fe.Tag()}
	}
}
