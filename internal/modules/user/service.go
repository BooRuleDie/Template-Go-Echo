package user

type UserService interface {
	GetUser(id int) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(user *User) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(id int) (*User, error) {
	return s.repo.GetUserById(id)
}

func (s *userService) CreateUser(user *User) error {
	return s.repo.CreateUser(user)
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) UpdateUser(user *User) error {
	return s.repo.UpdateUser(user)
}
