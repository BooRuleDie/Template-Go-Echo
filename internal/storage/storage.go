package storage

import (
	"context"
	"database/sql"
	"go-echo-template/internal/storage/auth"
	"go-echo-template/internal/storage/user"
)

type Storage struct {
	db   *sql.DB
	User user.UserRepository
	Auth auth.AuthRepository
}

func NewStorage(db *sql.DB, user user.UserRepository, auth auth.AuthRepository) *Storage {
	return &Storage{
		db:   db,
		User: user,
		Auth: auth,
	}
}

func (s *Storage) WithTx(ctx context.Context, fn func(*Storage) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txStorage := &Storage{
		db:   s.db,
		User: s.User.WithTx(tx),
		Auth: s.Auth.WithTx(tx),
	}

	if err := fn(txStorage); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
