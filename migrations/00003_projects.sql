-- +goose Up
-- +goose StatementBegin
CREATE TYPE project_member_role AS ENUM ('owner', 'member');
CREATE TYPE decision_status AS ENUM ('proposed', 'accepted', 'deprecated');

CREATE TABLE projects
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_by  UUID         NOT NULL REFERENCES users (id),
    created_at  TIMESTAMPTZ  NOT NULL
);

CREATE TABLE project_members
(
    project_id UUID                NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id    UUID                NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       project_member_role NOT NULL,
    joined_at  TIMESTAMPTZ         NOT NULL,
    PRIMARY KEY (project_id, user_id)
);

CREATE TABLE project_docs
(
    id         UUID PRIMARY KEY,
    project_id UUID         NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    slug       VARCHAR(120) NOT NULL,
    title      VARCHAR(200) NOT NULL,
    body_md    TEXT         NOT NULL DEFAULT '',
    created_by UUID         NOT NULL REFERENCES users (id),
    updated_at TIMESTAMPTZ  NOT NULL,
    UNIQUE (project_id, slug)
);

CREATE TABLE project_decisions
(
    id         UUID PRIMARY KEY,
    project_id UUID             NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    summary    VARCHAR(300)     NOT NULL,
    context    TEXT             NOT NULL DEFAULT '',
    status     decision_status  NOT NULL DEFAULT 'proposed',
    message_id UUID             REFERENCES messages (id) ON DELETE SET NULL,
    created_by UUID             NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ      NOT NULL
);

CREATE TABLE project_conversations
(
    project_id      UUID        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    conversation_id UUID        NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    linked_at       TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (project_id, conversation_id)
);

CREATE INDEX idx_project_members_user_id ON project_members (user_id);
CREATE INDEX idx_project_docs_project_id ON project_docs (project_id);
CREATE INDEX idx_project_decisions_project_id ON project_decisions (project_id);
CREATE INDEX idx_project_conversations_conversation_id ON project_conversations (conversation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_conversations;
DROP TABLE IF EXISTS project_decisions;
DROP TABLE IF EXISTS project_docs;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TYPE IF EXISTS decision_status;
DROP TYPE IF EXISTS project_member_role;
-- +goose StatementEnd
