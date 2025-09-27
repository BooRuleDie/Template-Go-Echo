package user

import (
	"net/http"
	"strconv"

	"go-echo-template/internal/alarm"
	constant "go-echo-template/internal/shared/constant"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	logger  log.CustomLogger
	alarmer alarm.Alarmer

	service userService
}

func NewUserHandler(logger log.CustomLogger, alarmer alarm.Alarmer, service userService) *UserHandler {
	return &UserHandler{logger: logger, alarmer: alarmer, service: service}
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	users := e.Group("/api/v1/users")
	users.GET("/:id", h.GetUser)
	users.POST("/", h.CreateUser)
	users.PATCH("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}

	// service call
	user, err := h.service.getUser(ctx, id)
	if err != nil {
		return err
	}

	// build response
	getUserResp := &GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     nil,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(constant.DefaultDateFormat),
		UpdatedAt: user.UpdatedAt.Format(constant.DefaultDateFormat),
	}
	if user.Phone.Valid {
		getUserResp.Phone = &user.Phone.String
	}
	return response.Success(c, http.StatusOK).WithData(getUserResp).Send()
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	// validate input
	cur := new(CreateUserRequest)
	if err := c.Bind(cur); err != nil {
		return errInvalidRequestPayload
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
	ctx := c.Request().Context()

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}
	uur := new(UpdateUserRequest)
	if err := c.Bind(&uur); err != nil {
		return errInvalidRequestPayload
	}
	if err := c.Validate(uur); err != nil {
		return err
	}

	// service call
	uur.ID = id // Ensure id from URL is used.
	err = h.service.updateUser(ctx, uur)
	if err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusOK).Send()
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	// validate input
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return errInvalidID.WithArgs(param)
	}

	// service call
	err = h.service.deleteUser(ctx, id)
	if err != nil {
		return err
	}

	// build response
	return response.Success(c, http.StatusNoContent).Send()
}
