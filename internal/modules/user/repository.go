package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Repository interface {
	GetUserById(id int) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(user *User) error
}

type userRepository struct {
	users map[int]*User
}

func NewUserRepository() Repository {
	return &userRepository{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
			3: {ID: 3, Name: "Carol", Email: "carol@example.com"},
		},
	}
}

func (r *userRepository) GetUserById(id int) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return ErrUserAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) DeleteUser(id int) error {
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *userRepository) UpdateUser(user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}
