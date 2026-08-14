-- +goose Up
CREATE TABLE chat_folders (
    id         UUID PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL CHECK (char_length(name) BETWEEN 1 AND 64),
    icon       TEXT,
    position   INT         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX chat_folders_user_id_position_idx ON chat_folders (user_id, position);

CREATE TABLE chat_folder_items (
    folder_id       UUID        NOT NULL REFERENCES chat_folders (id) ON DELETE CASCADE,
    conversation_id UUID        NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    added_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (folder_id, conversation_id)
);

CREATE INDEX chat_folder_items_conversation_id_idx ON chat_folder_items (conversation_id);

-- +goose Down
DROP TABLE IF EXISTS chat_folder_items;
DROP TABLE IF EXISTS chat_folders;
