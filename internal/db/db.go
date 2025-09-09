package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	config "go-echo-template/internal/config"
)

func New(DBConfig *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBConfig.GetConnectionString())
	if err != nil {
		return nil, err
	}

	// update the database config
	db.SetMaxOpenConns(DBConfig.MaxOpenConns)
	db.SetMaxIdleConns(DBConfig.MaxIdleConns)
	db.SetConnMaxIdleTime(DBConfig.MaxIdleTime)

	// if it takes more than 5 seconds to ping the
	// database, cancel the context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New("transaction error: " + err.Error() + ", rollback error: " + rbErr.Error())
		}
		return err
	}

	return tx.Commit()
}
