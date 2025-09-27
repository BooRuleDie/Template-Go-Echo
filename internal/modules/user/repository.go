package user

import (
	"context"
	"database/sql"
	"errors"

	DB "go-echo-template/internal/db"
	"go-echo-template/internal/modules/user/sqlc"
	"go-echo-template/internal/shared/log"
	"go-echo-template/internal/shared/response"

	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository interface {
	getUserById(ctx context.Context, userID int64) (*sqlc.User, error)
	createUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error)
	updateUser(ctx context.Context, params sqlc.UpdateUserParams) (*sqlc.User, error)
	deleteUser(ctx context.Context, userID int64) error
}

type repository struct {
	logger log.CustomLogger

	db      *sql.DB
	queries *sqlc.Queries

	cache userCache
}

func NewUserRepository(logger log.CustomLogger, db *sql.DB, cache userCache) userRepository {
	return &repository{logger: logger, db: db, queries: sqlc.New(db), cache: cache}
}

func (r *repository) getUserById(ctx context.Context, userID int64) (*sqlc.User, error) {
	userFromCache, err := r.cache.Get(ctx, userID)
	if err != nil {
		r.logger.WarnWithContext(ctx, "failed to get user from cache", r.logger.Err(err), r.logger.Int("userID", int(userID)))
		// Continue without cache, do not return error
	} else if userFromCache != nil {
		return userFromCache, nil
	}

	userRow, err := r.queries.GetUserById(ctx, userID)
	if err == sql.ErrNoRows {
		return nil, response.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user := &sqlc.User{
		ID:        userRow.ID,
		Name:      userRow.Name,
		Email:     userRow.Email,
		Phone:     userRow.Phone,
		Role:      userRow.Role,
		Password:  userRow.Password,
		CreatedAt: userRow.CreatedAt,
		UpdatedAt: userRow.UpdatedAt,
		IsDeleted: userRow.IsDeleted,
	}

	if err = r.cache.Set(ctx, user); err != nil {
		r.logger.WarnWithContext(ctx, "failed to set user in cache", r.logger.Err(err), r.logger.Int("userID", int(user.ID)))
		// Do not return error, continue
	}

	return user, nil
}

func (r *repository) createUser(ctx context.Context, params sqlc.CreateUserParams) (int64, error) {
	userID, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violation on users_email_unique_active
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, errUserEmailAlreadyExists
		}
		return 0, err
	}
	return userID, nil
}

func (r *repository) updateUser(ctx context.Context, params sqlc.UpdateUserParams) (*sqlc.User, error) {
	var user *sqlc.User
	if txErr := DB.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		qtx := r.queries.WithTx(tx)

		if err := qtx.UpdateUser(ctx, params); err != nil {
			if err == sql.ErrNoRows {
				return response.ErrUserNotFound
			}
			return err
		}

		userRow, err := qtx.GetUserById(ctx, params.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return response.ErrUserNotFound
			}
			return err
		}

		user = &sqlc.User{
			ID:        userRow.ID,
			Name:      userRow.Name,
			Email:     userRow.Email,
			Phone:     userRow.Phone,
			Role:      userRow.Role,
			Password:  userRow.Password,
			CreatedAt: userRow.CreatedAt,
			UpdatedAt: userRow.UpdatedAt,
			IsDeleted: userRow.IsDeleted,
		}
		return nil
	}); txErr != nil {
		return nil, txErr
	}

	if err := r.cache.Set(ctx, user); err != nil {
		r.logger.WarnWithContext(
			ctx,
			"failed to set user in cache during update",
			r.logger.Err(err),
			r.logger.Int("userID", int(params.ID)),
		)
		// Do not return error, continue
	}

	return user, nil
}

func (r *repository) deleteUser(ctx context.Context, userID int64) error {
	err := r.queries.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, userID); err != nil {
		r.logger.WarnWithContext(ctx, "failed to delete user from cache during delete", r.logger.Err(err), r.logger.Int("userID", int(userID)))
		// Do not return error, continue
	}
	return nil
}
