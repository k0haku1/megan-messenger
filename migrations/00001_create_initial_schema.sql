-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id             UUID PRIMARY KEY,
    username       VARCHAR(50)  NOT NULL UNIQUE,
    email          VARCHAR(255) NOT NULL UNIQUE,
    email_verified BOOLEAN      NOT NULL,
    password_hash  VARCHAR(255),
    created_at     TIMESTAMPTZ  NOT NULL
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
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_messages_conversation_created;
-- +goose StatementEnd
