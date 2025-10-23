package auth

import (
	"net/http"

	"go-echo-template/internal/alarm"
	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	logger  log.CustomLogger
	alarmer alarm.Alarmer

	service AuthService
}

func NewAuthHandler(logger log.CustomLogger, alarmer alarm.Alarmer, service AuthService) *AuthHandler {
	return &AuthHandler{logger: logger, alarmer: alarmer, service: service}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Group) {
	users := e.Group("/v1/auth")
	users.POST("/login", h.Login)
	users.GET("/refresh", h.Refresh)
	users.GET("/logout", h.Logout)
}

func (h *AuthHandler) Login(c echo.Context) error {
	// validate input
	lr := new(LoginRequest)
	if err := c.Bind(lr); err != nil {
		return shared.ErrInvalidRequestPayload
	}
	if err := c.Validate(lr); err != nil {
		return err
	}

	// service call
	if err := h.service.apiLogin(c, lr); err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).WithMessage(succLogin).Send()
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	// service call
	if err := h.service.apiRefresh(c); err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).Send()
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// service call
	if err := h.service.Logout(c); err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).Send()
}
