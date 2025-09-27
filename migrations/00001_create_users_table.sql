-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NULL,
    role TEXT NOT NULL CHECK (role IN ('user', 'admin', 'subadmin')),
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create partial unique index - only enforces uniqueness for non-deleted users
CREATE UNIQUE INDEX users_email_unique_active 
ON users (email) 
WHERE is_deleted = FALSE;

-- +goose Down
DROP INDEX IF EXISTS users_email_unique_active;
DROP TABLE users;