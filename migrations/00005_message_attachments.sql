-- +goose Up
CREATE TABLE message_attachments (
    id              UUID PRIMARY KEY,
    message_id      UUID REFERENCES messages (id) ON DELETE CASCADE,
    conversation_id UUID NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    uploader_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    object_key      TEXT NOT NULL,
    thumb_key       TEXT,
    mime            TEXT NOT NULL,
    kind            TEXT NOT NULL CHECK (kind IN ('image', 'video', 'file')),
    size_bytes      BIGINT NOT NULL DEFAULT 0,
    width           INT,
    height          INT,
    duration_ms     INT,
    original_name   TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX message_attachments_message_id_idx ON message_attachments (message_id);
CREATE INDEX message_attachments_conversation_media_idx
    ON message_attachments (conversation_id, kind, created_at DESC, id DESC)
    WHERE message_id IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS message_attachments;
