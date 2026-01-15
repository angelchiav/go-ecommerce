-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, role
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, role
FROM users
WHERE id = $1;