-- name: CreateFeed :one
INSERT INTO
    feeds (id, creates_at, updated_at, name, url, user_id)
VALUES
    ($1, NOW (), NOW (), $2, $3, $4)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM
    feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM
    feeds
order by
    last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedsAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW (), updated_at = NOW ()
WHERE
    id = $1
RETURNING *;