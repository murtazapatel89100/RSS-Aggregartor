-- name: CreateFeedFollow :one
INSERT INTO
    feeds_follow (id, creates_at, updated_at, user_id, feed_id)
VALUES
    ($1, NOW (), NOW (), $2, $3)
RETURNING *;
