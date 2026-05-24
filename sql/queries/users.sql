-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUserPasswordEmail :exec
UPDATE users
SET email = $2,
    hashed_password = $3,
    updated_at = $4
WHERE id = $1;

-- name: UpdateChirpyRed :exec
UPDATE users
SET is_chirpy_red = $2
WHERE id = $1; 