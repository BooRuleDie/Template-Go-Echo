package user

import (
	"database/sql"
	"go-echo-template/internal/shared/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	service userService
}

func NewUserHandler(db *sql.DB) *UserHandler {
	userRepo := newUserRepository(db)
	userService := newUserService(userRepo)
	return &UserHandler{service: userService}
}

// RegisterRoutes registers user routes to the Echo group.
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/v1/users/:id", h.GetUser)
	e.POST("/api/v1/users", h.CreateUser)
	e.PUT("/api/v1/users/:id", h.UpdateUser)
	e.DELETE("/api/v1/users/:id", h.DeleteUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errInvalidID
	}
	user, err := h.service.getUser(id)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(c, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return errInvalidID
	}
	err := h.service.createUser(&user)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(c, http.StatusCreated, nil)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errInvalidID
	}
	var user User
	if err := c.Bind(&user); err != nil {
		return errInvalidRequestPayload
	}
	user.ID = id // Ensure id from URL is used.
	err = h.service.updateUser(&user)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(c, http.StatusOK, nil)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errInvalidID
	}
	err = h.service.deleteUser(id)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(c, http.StatusNoContent, nil)
}
