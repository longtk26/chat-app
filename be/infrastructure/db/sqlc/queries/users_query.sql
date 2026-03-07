-- name: GetUserByEmail :one
SELECT id, username, email, password
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    id, username, email, password
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, username, email, password;

-- name: ListUsers :many
SELECT id, username, email, password
FROM users;
