-- name: CreatePost :one
INSERT INTO
    posts (id, created_at, updated_at, title, description, publisghed_at, url, feed_id)
VALUES
    ($1, NOW (), NOW (), $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostForUser :many
SELECT posts.*
FROM posts
JOIN feeds_follow
    ON posts.feed_id = feeds_follow.feed_id
WHERE
    feeds_follow.user_id = $1
ORDER BY
    posts.publisghed_at DESC
LIMIT $2 OFFSET $3;
