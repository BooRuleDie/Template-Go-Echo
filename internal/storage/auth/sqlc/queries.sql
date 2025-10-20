-- name: GetUserByEmail :one
SELECT 
    id, 
    name, 
    email, 
    phone, 
    role, 
    password,
    created_at, 
    updated_at
FROM users 
WHERE 
    email = $1 AND
    is_deleted = FALSE
LIMIT 1;

-- name: GetUserById :one
SELECT 
    id, 
    name, 
    email, 
    phone, 
    role, 
    password,
    created_at, 
    updated_at
FROM users 
WHERE 
    id = $1 AND
    is_deleted = FALSE
LIMIT 1;