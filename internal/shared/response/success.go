package response

import (
	"fmt"
	"go-echo-template/internal/shared/i18n"

	"github.com/labstack/echo/v4"
)

type SuccessMessage struct {
	Code     string
	Messages i18n.Messages
}

func (s *SuccessMessage) translate(locale i18n.Locale, args ...any) string {
	// Check if locale exists, else fallback
	if msg, ok := s.Messages[locale]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}

	// fallback locale
	if msg, ok := s.Messages[i18n.DefaultLocale]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}

	// ultimate fallback: return code
	return s.Code
}

// Standard Success Response
type successResponse struct {
	c echo.Context

	IsError bool   `json:"isError"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(c echo.Context, status int) *successResponse {
	return &successResponse{
		c:       c,
		IsError: false,
		Status:  status,
	}
}

func (s *successResponse) WithData(data any) *successResponse {
	s.Data = data
	return s
}

func (s *successResponse) WithMessage(m *SuccessMessage, args ...any) *successResponse {
	locale := i18n.GetLocaleFromContext(s.c)
	message := m.translate(locale, args...)
	s.Message = message
	return s
}

func (s *successResponse) Send() error {
	return s.c.JSON(s.Status, s)
}
