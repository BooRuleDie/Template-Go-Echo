package user

import "go-echo-template/internal/shared/models"

type userService interface {
	getUser(id int) (*models.User, error)
	createUser(*CreateUserRequest) error
	deleteUser(id int) error
	updateUser(user *models.User) error
}

type service struct {
	repo userRepository
}

func newUserService(repo userRepository) userService {
	return &service{repo: repo}
}

func (s *service) getUser(id int) (*models.User, error) {
	return s.repo.getUserById(id)
}

func (s *service) createUser(cur *CreateUserRequest) error {
	return s.repo.createUser(cur)
}

func (s *service) deleteUser(id int) error {
	return s.repo.deleteUser(id)
}

func (s *service) updateUser(user *models.User) error {
	return s.repo.updateUser(user)
}
