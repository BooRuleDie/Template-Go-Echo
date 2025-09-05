package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	service UserService
}

func NewUserHandler(db *sql.DB) *UserHandler {
	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	user, err := h.service.GetUser(id)
	if err != nil {
		if err == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}
	err := h.service.CreateUser(&user)
	if err != nil {
		if err == ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "user already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}
	user.ID = id // Ensure id from URL is used.
	err = h.service.UpdateUser(&user)
	if err != nil {
		if err == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	err = h.service.DeleteUser(id)
	if err != nil {
		if err == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.NoContent(http.StatusNoContent)
}
