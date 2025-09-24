package user

import (
	"context"
	"database/sql"
	"go-echo-template/internal/modules/user/sqlc"
	constant "go-echo-template/internal/shared/const"
	"go-echo-template/internal/shared/utils"
)

type userService interface {
	getUser(ctx context.Context, id int64) (*sqlc.User, error)
	createUser(ctx context.Context, cur *CreateUserRequest) (int64, error)
	deleteUser(ctx context.Context, id int64) error
	updateUser(ctx context.Context, uur *UpdateUserRequest) error
}

type service struct {
	repo userRepository
}

func newUserService(repo userRepository) userService {
	return &service{repo: repo}
}

func (s *service) getUser(ctx context.Context, id int64) (*sqlc.User, error) {
	// repo call
	return s.repo.getUserById(ctx, id)
}

func (s *service) createUser(ctx context.Context, cur *CreateUserRequest) (int64, error) {
	// convert dto to params
	var phone sql.NullString
	if cur.Phone != nil {
		phone.Valid = true
		phone.String = *cur.Phone
	}
	password, err := utils.HashPassword(cur.Password)
	if err != nil {
		return 0, err
	}
	params := sqlc.CreateUserParams{
		Name:     cur.Name,
		Email:    cur.Email,
		Phone:    phone,
		Role:     constant.RoleCustomer,
		Password: password,
	}
	
	// repo call
	return s.repo.createUser(ctx, params)
}

func (s *service) updateUser(ctx context.Context, uur *UpdateUserRequest) error {
	// convert dto to params
	var phone sql.NullString
	if uur.Phone != nil {
		phone.Valid = true
		phone.String = *uur.Phone
	}
	params := sqlc.UpdateUserParams{
		ID:    uur.ID,
		Name:  uur.Name,
		Email: uur.Email,
		Phone: phone,
	}
	
	// repo call
	return s.repo.updateUser(ctx, params)
}

func (s *service) deleteUser(ctx context.Context, id int64) error {
	// repo call
	return s.repo.deleteUser(ctx, id)
}
