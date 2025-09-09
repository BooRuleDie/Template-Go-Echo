package user

type userService interface {
	getUser(id int) (*User, error)
	createUser(user *User) error
	deleteUser(id int) error
	updateUser(user *User) error
}

type service struct {
	repo userRepository
}

func newUserService(repo userRepository) userService {
	return &service{repo: repo}
}

func (s *service) getUser(id int) (*User, error) {
	return s.repo.getUserById(id)
}

func (s *service) createUser(user *User) error {
	return s.repo.createUser(user)
}

func (s *service) deleteUser(id int) error {
	return s.repo.deleteUser(id)
}

func (s *service) updateUser(user *User) error {
	return s.repo.updateUser(user)
}
