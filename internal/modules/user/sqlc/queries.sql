-- name: GetUserById :one
SELECT id, name, email, phone, role, password, created_at, updated_at, is_delete FROM users WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (name, email, phone, role, password)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteUser :exec
UPDATE users SET is_delete = TRUE, updated_at = NOW() WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET name = $1, email = $2, phone = $3, role = $4, password = $5, updated_at = NOW() WHERE id = $6;