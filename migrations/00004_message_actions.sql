-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages
    ADD COLUMN IF NOT EXISTS reply_to_id UUID REFERENCES messages (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS forwarded_from_id UUID REFERENCES messages (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS hidden_messages
(
    message_id UUID NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE TABLE IF NOT EXISTS message_reactions
(
    message_id UUID        NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    emoji      VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id, emoji)
);

CREATE INDEX IF NOT EXISTS idx_message_reactions_message ON message_reactions (message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_reactions;
DROP TABLE IF EXISTS hidden_messages;
ALTER TABLE messages
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS forwarded_from_id,
    DROP COLUMN IF EXISTS reply_to_id;
-- +goose StatementEnd
