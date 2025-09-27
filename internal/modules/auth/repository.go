package auth

import (
	"context"
	"database/sql"

	"go-echo-template/internal/modules/auth/sqlc"
	"go-echo-template/internal/shared/log"
)

type authRepository interface {
	getUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error)
	getUserById(ctx context.Context, userID int64) (*sqlc.GetUserByIdRow, error)
}

type repository struct {
	logger log.CustomLogger

	db      *sql.DB
	queries *sqlc.Queries
}

func NewAuthRepository(logger log.CustomLogger, db *sql.DB) authRepository {
	return &repository{logger: logger, db: db, queries: sqlc.New(db)}
}

func (r *repository) getUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error) {
	userRow, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errEmailOrPasswordWrong
		}
		return nil, err
	}

	return &userRow, nil
}

func (r *repository) getUserById(ctx context.Context, userID int64) (*sqlc.GetUserByIdRow, error) {
	userRow, err := r.queries.GetUserById(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errEmailOrPasswordWrong
		}
		return nil, err
	}

	return &userRow, nil
}
