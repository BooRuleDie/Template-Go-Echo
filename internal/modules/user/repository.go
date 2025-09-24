package user

import (
	"context"
	"database/sql"
	"go-echo-template/internal/modules/user/sqlc"
	constant "go-echo-template/internal/shared/const"
)

type userRepository interface {
	getUserById(ctx context.Context, id int64) (*sqlc.User, error)
	createUser(ctx context.Context, cur *CreateUserRequest) error
	deleteUser(ctx context.Context, id int64) error
	updateUser(ctx context.Context, uur *UpdateUserRequest) error
}

type repository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

func newUserRepository(db *sql.DB) *repository {
	return &repository{queries: sqlc.New(db), db: db}
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
		IsDelete:  userRow.IsDelete,
	}, nil
}

func (r *repository) createUser(ctx context.Context, cur *CreateUserRequest) error {
	var phone sql.NullString
	if cur.Phone != nil {
		phone.Valid = true
		phone.String = *cur.Phone
	}

	password, err := cur.HashPassword()
	if err != nil {
		return err
	}
	params := sqlc.CreateUserParams{
		Name:     cur.Name,
		Email:    cur.Email,
		Phone:    phone,
		Role:     constant.RoleCustomer,
		Password: password,
	}
	return r.queries.CreateUser(ctx, params)
}

func (r *repository) deleteUser(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *repository) updateUser(ctx context.Context, uur *UpdateUserRequest) error {
	var phone sql.NullString
	if uur.Phone != nil {
		phone.Valid = true
		phone.String = *uur.Phone
	}
	params := sqlc.UpdateUserParams{
		ID:    uur.ID,
		Name:  uur.Name,
		Email: uur.Email,
		Phone: phone,
	}
	err := r.queries.UpdateUser(ctx, params)
	if err == sql.ErrNoRows {
		return errUserNotFound.WithArgs(uur.ID)
	}
	return err
}
