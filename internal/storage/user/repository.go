package user

import (
	"context"
	"database/sql"

	"go-echo-template/internal/shared"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/storage/user/sqlc"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userID int64) (*sqlc.User, error)
	CreateUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error)
	UpdateUser(ctx context.Context, params sqlc.UpdateUserParams) error
	DeleteUser(ctx context.Context, userID int64) error

	// transaction
	WithTx(tx *sql.Tx) UserRepository
}

type repository struct {
	logger  log.CustomLogger
	queries *sqlc.Queries
	db      sqlc.DBTX
	cache   UserCache
}

func NewUserRepository(logger log.CustomLogger, db *sql.DB, cache UserCache) UserRepository {
	return &repository{logger: logger, db: db, queries: sqlc.New(db), cache: cache}
}

func (r *repository) WithTx(tx *sql.Tx) UserRepository {
	return &repository{
		logger:  r.logger,
		queries: sqlc.New(tx),
		db:      tx,
		cache:   r.cache,
	}
}

func (r *repository) GetUserById(ctx context.Context, userID int64) (*sqlc.User, error) {
	userFromCache, err := r.cache.Get(ctx, userID)
	if err != nil && err != redis.Nil {
		r.logger.WarnWithContext(
			ctx,
			"failed to get user from cache",
			r.logger.Err(err),
			r.logger.Int("userID", int(userID)),
		)
		// Continue without cache, do not return error
	}

	// if it's a cache hit return immediately
	if err == nil && userFromCache != nil {
		return userFromCache, nil
	}

	userRow, err := r.queries.GetUserById(ctx, userID)
	if err == sql.ErrNoRows {
		return nil, shared.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if err = r.cache.Set(ctx, &userRow); err != nil {
		r.logger.WarnWithContext(
			ctx,
			"failed to set user in cache",
			r.logger.Err(err),
			r.logger.Int("userID", int(userRow.ID)),
		)
		// Do not return error, continue
	}

	return &userRow, nil
}

func (r *repository) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *repository) UpdateUser(ctx context.Context, params sqlc.UpdateUserParams) error {
	if err := r.queries.UpdateUser(ctx, params); err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, params.ID); err != nil {
		r.logger.WarnWithContext(
			ctx,
			"failed to delete user from cache during delete",
			r.logger.Err(err),
			r.logger.Int("userID", int(params.ID)),
		)
		// Do not return error, continue
	}

	return nil
}

func (r *repository) DeleteUser(ctx context.Context, userID int64) error {
	err := r.queries.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, userID); err != nil {
		r.logger.WarnWithContext(
			ctx,
			"failed to delete user from cache during delete",
			r.logger.Err(err),
			r.logger.Int("userID", int(userID)),
		)
		// Do not return error, continue
	}

	return nil
}
