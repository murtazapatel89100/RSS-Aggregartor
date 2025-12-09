-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users
ADD COLUMN api_key VARCHAR(255) UNIQUE NOT NULL DEFAULT (encode (sha256 (gen_random_bytes (16)), 'hex'));

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;