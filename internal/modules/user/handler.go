package user

import (
	"net/http"
	"strconv"

	"go-echo-template/internal/alarm"
	"go-echo-template/internal/modules/auth"
	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	logger  log.CustomLogger
	alarmer alarm.Alarmer

	service userService
	auth    auth.AuthService
}

func NewUserHandler(logger log.CustomLogger, alarmer alarm.Alarmer, service userService, authService auth.AuthService) *UserHandler {
	return &UserHandler{logger: logger, alarmer: alarmer, service: service, auth: authService}
}

func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	users := e.Group("/v1/users")
	// public API
	users.POST("/", h.CreateUser)

	// authenticated APIs
	usersAuth := users.Group("", h.auth.CheckAuth(false, shared.RoleCustomer))
	usersAuth.GET("/:id", h.GetUser)
	usersAuth.PATCH("/:id", h.UpdateUser)
	usersAuth.DELETE("/:id", h.DeleteUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	userFromCtx, ok := auth.GetUserFromContext(c)
	if !ok {
		return shared.ErrUserNotFound
	}

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}

	// Access Control
	if userFromCtx.ID != id {
		return shared.ErrSessionUnauthorized
	}

	// service call
	user, err := h.service.getUser(ctx, id)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK).WithData(user).Send()
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	// validate input
	cur := new(CreateUserRequest)
	if err := c.Bind(cur); err != nil {
		return shared.ErrInvalidRequestPayload
	}
	if err := c.Validate(cur); err != nil {
		return err
	}

	// service call
	newUserID, err := h.service.createUser(ctx, cur)
	if err != nil {
		return err
	}

	// build response
	//
	resData := &CreateUserResponse{
		UserID: newUserID,
	}
	return response.Success(c, http.StatusCreated).WithMessage(succUserCreated).WithData(resData).Send()
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userFromCtx, ok := auth.GetUserFromContext(c)
	if !ok {
		return shared.ErrUserNotFound
	}

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}
	uur := new(UpdateUserRequest)
	if err := c.Bind(&uur); err != nil {
		return shared.ErrInvalidRequestPayload
	}
	if err := c.Validate(uur); err != nil {
		return err
	}

	// Access Control
	if userFromCtx.ID != id {
		return shared.ErrSessionUnauthorized
	}

	// service call
	uur.ID = id // Ensure id from URL is used.
	err = h.service.updateUser(c, uur)
	if err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).Send()
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userFromCtx, ok := auth.GetUserFromContext(c)
	if !ok {
		return shared.ErrUserNotFound
	}

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}

	// Access Control
	if userFromCtx.ID != id {
		return shared.ErrSessionUnauthorized
	}

	// service call
	err = h.service.deleteUser(c, id)
	if err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).Send()
}
