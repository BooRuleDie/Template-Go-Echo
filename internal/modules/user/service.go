package user

import (
	"context"
	"database/sql"

	"go-echo-template/internal/modules/auth"
	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/utils"
	"go-echo-template/internal/storage"
	"go-echo-template/internal/storage/user/sqlc"

	"github.com/labstack/echo/v4"
)

type userService interface {
	getUser(ctx context.Context, id int64) (*GetUserResponse, error)
	createUser(ctx context.Context, cur *CreateUserRequest) (int64, error)
	updateUser(c echo.Context, uur *UpdateUserRequest) error
	deleteUser(c echo.Context, id int64) error
}

type service struct {
	logger  log.CustomLogger
	storage *storage.Storage
	auth    auth.AuthService
}

func NewUserService(logger log.CustomLogger, storage *storage.Storage, authService auth.AuthService) userService {
	return &service{storage: storage, logger: logger, auth: authService}
}

func (s *service) getUser(ctx context.Context, id int64) (*GetUserResponse, error) {
	// repo call
	user, err := s.storage.User.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	// build response
	getUserResp := &GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     nil,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(shared.DefaultDateFormat),
		UpdatedAt: user.UpdatedAt.Format(shared.DefaultDateFormat),
	}
	if user.Phone.Valid {
		getUserResp.Phone = &user.Phone.String
	}

	return getUserResp, nil
}

func (s *service) createUser(ctx context.Context, cur *CreateUserRequest) (int64, error) {
	password, err := utils.HashPassword(cur.Password)
	if err != nil {
		return 0, err
	}

	params := sqlc.CreateUserParams{
		Name:     cur.Name,
		Email:    cur.Email,
		Phone:    sql.NullString{},
		Role:     shared.RoleCustomer,
		Password: password,
	}
	if cur.Phone != nil {
		params.Phone.Valid = true
		params.Phone.String = *cur.Phone
	}

	// repo call
	return s.storage.User.CreateUser(ctx, params)
}

func (s *service) updateUser(c echo.Context, uur *UpdateUserRequest) error {
	ctx := c.Request().Context()

	// convert dto to params
	params := sqlc.UpdateUserParams{
		ID:    uur.ID,
		Name:  uur.Name,
		Email: uur.Email,
		Phone: sql.NullString{},
	}
	if uur.Phone != nil {
		params.Phone.Valid = true
		params.Phone.String = *uur.Phone
	}

	// transaction example in service layer using storage
	var newUser *sqlc.User
	if err := s.storage.WithTx(ctx, func(storageTx *storage.Storage) error {
		err := storageTx.User.UpdateUser(ctx, params)
		if err != nil {
			return err
		}

		user, err := storageTx.User.GetUserById(ctx, params.ID)
		if err != nil {
			return err
		}

		newUser = user
		return nil
	}); err != nil {
		return err
	}

	// refresh token data
	sessionUser := &auth.User{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Phone:     newUser.Phone.String,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	if err := s.auth.Refresh(c, sessionUser); err != nil {
		s.logger.Error("delete user session after removal is failed", s.logger.Err(err))
	}

	return nil
}

func (s *service) deleteUser(c echo.Context, id int64) error {
	if err := s.auth.Logout(c); err != nil {
		s.logger.Error("delete user session after removal is failed", s.logger.Err(err))
	}

	// repo call
	return s.storage.User.DeleteUser(c.Request().Context(), id)
}
