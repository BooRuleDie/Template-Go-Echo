package user

import (
	"database/sql"
	"go-echo-template/internal/shared/models"
)

type userRepository interface {
	getUserById(id int) (*models.User, error)
	createUser(*CreateUserRequest) error
	deleteUser(id int) error
	updateUser(user *models.User) error
}

type repository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) userRepository {
	return &repository{db: db}
}

func (r *repository) getUserById(id int) (*models.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user models.User
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

func (r *repository) createUser(cur *CreateUserRequest) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2);`
	_, err := r.db.Exec(
		query,
		cur.Name, cur.Email,
	)

	return err
}

func (r *repository) deleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)

	return err
}

func (r *repository) updateUser(user *models.User) error {
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
