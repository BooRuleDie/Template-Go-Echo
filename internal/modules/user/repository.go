package user

import (
	"context"
	"database/sql"

	"go-echo-template/internal/modules/user/sqlc"
)

type userRepository interface {
	getUserById(ctx context.Context, id int64) (*sqlc.User, error)
	createUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error)
	deleteUser(ctx context.Context, id int64) error
	updateUser(ctx context.Context, params sqlc.UpdateUserParams) error
}

type repository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func newUserRepository(db *sql.DB) *repository {
	return &repository{db: db, queries: sqlc.New(db)}
}

func (r *repository) getUserById(ctx context.Context, id int64) (*sqlc.User, error) {
	userRow, err := r.queries.GetUserById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound.WithArgs(id)
	}
	if err != nil {
		return nil, err
	}
	return &sqlc.User{
		ID:        userRow.ID,
		Name:      userRow.Name,
		Email:     userRow.Email,
		Phone:     userRow.Phone,
		Role:      userRow.Role,
		Password:  userRow.Password,
		CreatedAt: userRow.CreatedAt,
		UpdatedAt: userRow.UpdatedAt,
		IsDeleted: userRow.IsDeleted,
	}, nil
}

func (r *repository) createUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error) {
	userID, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *repository) updateUser(ctx context.Context, params sqlc.UpdateUserParams) error {
	err := r.queries.UpdateUser(ctx, params)
	if err == sql.ErrNoRows {
		return errUserNotFound.WithArgs(params.ID)
	}
	return err
}

func (r *repository) deleteUser(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}
