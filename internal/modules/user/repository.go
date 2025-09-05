package user

import (
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Userepository interface {
	GetUserById(id int) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(user *User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Userepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserById(id int) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	user := new(User)
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(
		query,
		user.Name, user.Email,
	)

	return err
}

func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)

	return err
}

func (r *userRepository) UpdateUser(user *User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	result, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
