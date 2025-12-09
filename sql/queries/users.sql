-- name: CreateUser :one
INSERT INTO
    users (id, creates_at, updated_at, username)
VALUES
    ($1, NOW (), NOW (), $2)
RETURNING *;