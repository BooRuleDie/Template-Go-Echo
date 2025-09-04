package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	GetUserById(id int) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(user *User) error
}

type repository struct {
	users map[int]*User
}

func NewRepository() UserRepository {
	return &repository{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
			3: {ID: 3, Name: "Carol", Email: "carol@example.com"},
		},
	}
}

func (r *repository) GetUserById(id int) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *repository) CreateUser(user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return ErrUserAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *repository) DeleteUser(id int) error {
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *repository) UpdateUser(user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}
