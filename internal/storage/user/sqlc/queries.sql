-- name: GetUserById :one
SELECT 
    id, 
    name, 
    email, 
    phone, 
    role, 
    password, 
    created_at, 
    updated_at, 
    is_deleted
FROM users 
WHERE 
    id = $1 AND
    is_deleted = FALSE;

-- name: CreateUser :one
INSERT INTO users (name, email, phone, role, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users SET name = $1, email = $2, phone = $3, updated_at = NOW() WHERE id = $4 AND is_deleted = FALSE;

-- name: DeleteUser :exec
UPDATE users SET is_deleted = TRUE, updated_at = NOW() WHERE id = $1 AND is_deleted = FALSE;