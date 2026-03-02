-- name: GetUserByEmail :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE id = ?
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, email, password, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateUserPassword :exec
UPDATE users
SET password = ?, updated_at = ?
WHERE id = ?;

-- name: UserExists :one
SELECT COUNT(*) FROM users LIMIT 1;
