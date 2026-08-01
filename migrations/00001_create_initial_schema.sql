-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id             UUID PRIMARY KEY,
    phone          VARCHAR(16) NOT NULL UNIQUE,
    username       VARCHAR(32) UNIQUE,
    password_hash  VARCHAR(255),
    avatar_url     TEXT,
    created_at     TIMESTAMPTZ NOT NULL,
    CONSTRAINT users_username_format CHECK (
        username IS NULL OR username ~ '^[a-zA-Z0-9]{3,32}$'
    )
);

CREATE TABLE phone_otp_codes
(
    phone      VARCHAR(16) PRIMARY KEY,
    code_hash  VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    attempts   INT          NOT NULL DEFAULT 0,
    sent_at    TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL
);

CREATE TYPE conversation_type AS ENUM ('dm', 'group');

CREATE TABLE conversations
(
    id         UUID PRIMARY KEY,
    type       conversation_type NOT NULL,
    title      VARCHAR(100),
    slug       VARCHAR(16) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT conversations_group_requires_title
        CHECK (type <> 'group' OR title IS NOT NULL),
    CONSTRAINT conversations_dm_no_slug
        CHECK (type <> 'dm' OR slug IS NULL)
);

CREATE TABLE conversation_members
(
    id              UUID PRIMARY KEY,
    conversation_id UUID        NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    user_id         UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    joined_at       TIMESTAMPTZ NOT NULL,
    last_read_at    TIMESTAMPTZ,
    UNIQUE (conversation_id, user_id)
);

CREATE TABLE dm_pairs
(
    user_low        UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user_high       UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    conversation_id UUID NOT NULL UNIQUE REFERENCES conversations (id) ON DELETE CASCADE,
    PRIMARY KEY (user_low, user_high),
    CHECK (user_low < user_high)
);

CREATE TABLE messages
(
    id              UUID PRIMARY KEY,
    conversation_id UUID        NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    sender_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content         TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_messages_conversation_created
    ON messages (conversation_id, created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS dm_pairs;
DROP TABLE IF EXISTS conversation_members;
DROP TABLE IF EXISTS conversations;
DROP TYPE IF EXISTS conversation_type;
DROP TABLE IF EXISTS phone_otp_codes;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
