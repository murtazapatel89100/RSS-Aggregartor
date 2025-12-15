-- name: CreateFeedFollow :one
INSERT INTO
    feeds_follow (id, creates_at, updated_at, user_id, feed_id)
VALUES
    ($1, NOW (), NOW (), $2, $3)
RETURNING *;

-- name: GetFeedFollow :many
SELECT * FROM
    feeds_follow
WHERE
    user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM
    feeds_follow
WHERE
    user_id = $1 AND feed_id = $2;
