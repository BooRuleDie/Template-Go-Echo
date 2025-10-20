package auth

import (
	"context"
	"database/sql"

	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/storage/auth/sqlc"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error)
	GetUserById(ctx context.Context, userID int64) (*sqlc.GetUserByIdRow, error)
}

type repository struct {
	logger log.CustomLogger

	db      *sql.DB
	queries *sqlc.Queries
}

func NewAuthRepository(logger log.CustomLogger, db *sql.DB) AuthRepository {
	return &repository{logger: logger, db: db, queries: sqlc.New(db)}
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error) {
	userRow, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &userRow, nil
}

func (r *repository) GetUserById(ctx context.Context, userID int64) (*sqlc.GetUserByIdRow, error) {
	userRow, err := r.queries.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &userRow, nil
}
