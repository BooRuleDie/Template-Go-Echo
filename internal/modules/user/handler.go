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
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errInvalidID
	}
	user, err := h.service.getUser(ctx, id)
	if err != nil {
		return err
	}
	
	getUserResp := &GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     nil,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if user.Phone.Valid {
		getUserResp.Phone = &user.Phone.String
	}

	return response.Success(c, http.StatusOK).WithData(getUserResp).Send()
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	cur := new(CreateUserRequest)
	if err := c.Bind(cur); err != nil {
		return errInvalidRequestPayload
	}
	if err := c.Validate(cur); err != nil {
		return err
	}

	err := h.service.createUser(ctx, cur)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusCreated).WithMessage("SUC:USER_CREATED").Send()
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errInvalidID
	}
	uur := new(UpdateUserRequest)
	if err := c.Bind(&uur); err != nil {
		return errInvalidRequestPayload
	}

	uur.ID = id // Ensure id from URL is used.
	err = h.service.updateUser(ctx, uur)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK).Send()
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errInvalidID
	}
	err = h.service.deleteUser(ctx, id)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusNoContent).Send()
}
