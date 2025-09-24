package user

import (
	"context"
	"go-echo-template/internal/modules/user/sqlc"
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
	return s.repo.getUserById(ctx, id)
}

func (s *service) createUser(ctx context.Context, cur *CreateUserRequest) (int64, error) {
	return s.repo.createUser(ctx, cur)
}

func (s *service) updateUser(ctx context.Context, uur *UpdateUserRequest) error {
	return s.repo.updateUser(ctx, uur)
}

func (s *service) deleteUser(ctx context.Context, id int64) error {
	return s.repo.deleteUser(ctx, id)
}
