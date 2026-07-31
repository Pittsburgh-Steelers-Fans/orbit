-- name: GetUser :one
SELECT id, email, display_name, avatar_url, created_at, updated_at, disabled_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, display_name, avatar_url, created_at, updated_at, disabled_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, display_name, avatar_url, created_at, updated_at, disabled_at
FROM users
WHERE disabled_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, display_name, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, email, display_name, avatar_url, created_at, updated_at, disabled_at;

-- name: UpdateUserProfile :one
UPDATE users
SET display_name = $2,
    avatar_url = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, display_name, avatar_url, created_at, updated_at, disabled_at;
