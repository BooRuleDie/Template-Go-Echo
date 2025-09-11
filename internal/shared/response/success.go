package response

import (
	"go-echo-template/internal/shared/i18n"

	"github.com/labstack/echo/v4"
)

// Standard Success Response
type successResponse struct {
	ctx echo.Context

	IsError bool   `json:"isError"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(c echo.Context, status int) *successResponse {
	return &successResponse{
		ctx:     c,
		IsError: false,
		Status:  status,
	}
}

func (s *successResponse) WithData(data any) *successResponse {
	s.Data = data
	return s
}

func (s *successResponse) WithMessage(code string, args ...any) *successResponse {
	locale := i18n.GetLocaleFromContext(s.ctx)
	message := i18n.Translate(code, locale, args...)
	s.Message = message
	return s
}

func (s *successResponse) Send() error {
	return s.ctx.JSON(s.Status, s)
}
