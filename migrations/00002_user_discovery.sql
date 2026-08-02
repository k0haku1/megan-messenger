-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN username_searchable BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN dm_policy          TEXT        NOT NULL DEFAULT 'everyone';

ALTER TABLE users
    ADD CONSTRAINT users_dm_policy_check
        CHECK (dm_policy IN ('everyone', 'nobody'));

ALTER TABLE users
    DROP CONSTRAINT users_username_format;

ALTER TABLE users
    ADD CONSTRAINT users_username_format CHECK (
        username IS NULL OR username ~ '^[a-z][a-z0-9_]{4,31}$'
        );

CREATE INDEX idx_users_username_lower_prefix
    ON users (lower(username) text_pattern_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_username_lower_prefix;

ALTER TABLE users
    DROP CONSTRAINT users_username_format;

ALTER TABLE users
    ADD CONSTRAINT users_username_format CHECK (
        username IS NULL OR username ~ '^[a-zA-Z0-9]{3,32}$'
        );

ALTER TABLE users
    DROP CONSTRAINT users_dm_policy_check;

ALTER TABLE users
    DROP COLUMN dm_policy,
    DROP COLUMN username_searchable;
-- +goose StatementEnd
