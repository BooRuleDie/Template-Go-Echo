package user

import (
	"database/sql"
)

type userRepository interface {
	getUserById(id int) (*User, error)
	createUser(user *User) error
	deleteUser(id int) error
	updateUser(user *User) error
}

type repository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) userRepository {
	return &repository{db: db}
}

func (r *repository) getUserById(id int) (*User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user User
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound.WithArgs(id)
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) createUser(user *User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2);`
	_, err := r.db.Exec(
		query,
		user.Name, user.Email,
	)

	return err
}

func (r *repository) deleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)

	return err
}

func (r *repository) updateUser(user *User) error {
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
		return errUserNotFound.WithArgs(user.ID)
	}
	return nil
}
