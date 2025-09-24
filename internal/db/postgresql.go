package db

import (
	"context"
	"database/sql"
	"time"

	config "go-echo-template/internal/config"
)

func NewPostgreSQL(ctx context.Context, DBConfig *config.DBConfig) (*sql.DB, error) {
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
