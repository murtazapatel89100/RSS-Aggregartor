-- +goose Up
CREATE TABLE
    feeds_follow (
        id UUID PRIMARY KEY,
        creates_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
        UNIQUE (user_id, feed_id)
    );

-- +goose Down
DROP TABLE feeds_follow;