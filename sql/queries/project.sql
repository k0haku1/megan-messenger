-- name: CreateProject :one
INSERT INTO projects (id, name, description, created_by, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProject :one
SELECT *
FROM projects
WHERE id = $1;

-- name: ListUserProjects :many
SELECT p.id,
       p.name,
       p.description,
       p.created_by,
       p.created_at,
       (SELECT COUNT(*)::INT FROM project_members pm WHERE pm.project_id = p.id) AS member_count,
       (SELECT COUNT(*)::INT FROM project_docs pd WHERE pd.project_id = p.id)    AS doc_count
FROM projects p
         JOIN project_members pm ON pm.project_id = p.id AND pm.user_id = $1
ORDER BY p.created_at DESC;

-- name: IsProjectMember :one
SELECT EXISTS (SELECT 1
               FROM project_members
               WHERE project_id = $1
                 AND user_id = $2) AS is_member;

-- name: GetProjectMemberRole :one
SELECT role
FROM project_members
WHERE project_id = $1
  AND user_id = $2;

-- name: AddProjectMember :exec
INSERT INTO project_members (project_id, user_id, role, joined_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (project_id, user_id) DO NOTHING;

-- name: ListProjectMembers :many
SELECT pm.project_id,
       pm.user_id,
       pm.role,
       pm.joined_at,
       u.username,
       u.avatar_url
FROM project_members pm
         JOIN users u ON u.id = pm.user_id
WHERE pm.project_id = $1
ORDER BY pm.joined_at ASC;

-- name: CreateProjectDoc :one
INSERT INTO project_docs (id, project_id, slug, title, body_md, created_by, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProjectDoc :one
SELECT *
FROM project_docs
WHERE id = $1
  AND project_id = $2;

-- name: GetProjectDocBySlug :one
SELECT *
FROM project_docs
WHERE project_id = $1
  AND slug = $2;

-- name: ListProjectDocs :many
SELECT id, project_id, slug, title, created_by, updated_at
FROM project_docs
WHERE project_id = $1
ORDER BY updated_at DESC;

-- name: UpdateProjectDoc :one
UPDATE project_docs
SET title    = $3,
    body_md  = $4,
    updated_at = $5
WHERE id = $1
  AND project_id = $2
RETURNING *;

-- name: DeleteProjectDoc :exec
DELETE
FROM project_docs
WHERE id = $1
  AND project_id = $2;

-- name: ProjectDocSlugExists :one
SELECT EXISTS (SELECT 1
               FROM project_docs
               WHERE project_id = $1
                 AND slug = $2) AS exists;

-- name: CreateProjectDecision :one
INSERT INTO project_decisions (id, project_id, summary, context, status, message_id, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListProjectDecisions :many
SELECT *
FROM project_decisions
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: GetProjectDecision :one
SELECT *
FROM project_decisions
WHERE id = $1
  AND project_id = $2;

-- name: UpdateProjectDecisionStatus :one
UPDATE project_decisions
SET status = $3
WHERE id = $1
  AND project_id = $2
RETURNING *;

-- name: ListConversationProjects :many
SELECT p.id,
       p.name,
       p.description,
       p.created_by,
       p.created_at
FROM project_conversations pc
         JOIN projects p ON p.id = pc.project_id
         JOIN project_members pm ON pm.project_id = p.id AND pm.user_id = $2
WHERE pc.conversation_id = $1
ORDER BY p.name ASC;

-- name: IsProjectConversationLinked :one
SELECT EXISTS (SELECT 1
               FROM project_conversations
               WHERE project_id = $1
                 AND conversation_id = $2) AS linked;

-- name: LinkProjectConversation :exec
INSERT INTO project_conversations (project_id, conversation_id, linked_at)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, conversation_id) DO NOTHING;

-- name: UnlinkProjectConversation :exec
DELETE
FROM project_conversations
WHERE project_id = $1
  AND conversation_id = $2;

-- name: ListProjectConversations :many
SELECT pc.project_id,
       pc.conversation_id,
       pc.linked_at,
       c.type,
       c.title,
       c.slug
FROM project_conversations pc
         JOIN conversations c ON c.id = pc.conversation_id
WHERE pc.project_id = $1
ORDER BY pc.linked_at DESC;
