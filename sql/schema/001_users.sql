-- +goose Up
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        creates_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
        username TEXT NOT NULL UNIQUE
    );

-- +goose Down
DROP TABLE users;