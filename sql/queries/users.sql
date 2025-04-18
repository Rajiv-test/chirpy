-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserbyEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = $1,updated_at = NOW(),email = $2
WHERE id = $3
RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = TRUE WHERE id = $1;