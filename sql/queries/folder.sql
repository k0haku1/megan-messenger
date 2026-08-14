-- name: CreateChatFolder :one
INSERT INTO chat_folders (id, user_id, name, icon, position, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetChatFolder :one
SELECT *
FROM chat_folders
WHERE id = $1;

-- name: ListChatFolders :many
SELECT *
FROM chat_folders
WHERE user_id = $1
ORDER BY position, created_at;

-- name: UpdateChatFolder :one
UPDATE chat_folders
SET name       = $2,
    icon       = $3,
    updated_at = $4
WHERE id = $1
  AND user_id = $5
RETURNING *;

-- name: DeleteChatFolder :exec
DELETE FROM chat_folders
WHERE id = $1
  AND user_id = $2;

-- name: MaxChatFolderPosition :one
SELECT COALESCE(MAX(position), -1)::int
FROM chat_folders
WHERE user_id = $1;

-- name: UpdateChatFolderPosition :exec
UPDATE chat_folders
SET position   = $2,
    updated_at = NOW()
WHERE id = $1
  AND user_id = $3;

-- name: CountUserChatFolders :one
SELECT COUNT(*)::int
FROM chat_folders
WHERE user_id = $1;

-- name: AddChatFolderItem :exec
INSERT INTO chat_folder_items (folder_id, conversation_id, added_at)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveChatFolderItem :exec
DELETE FROM chat_folder_items
WHERE folder_id = $1
  AND conversation_id = $2;

-- name: CountChatFolderItems :one
SELECT COUNT(*)::int
FROM chat_folder_items
WHERE folder_id = $1;

-- name: ChatFolderItemExists :one
SELECT EXISTS (
    SELECT 1
    FROM chat_folder_items
    WHERE folder_id = $1
      AND conversation_id = $2
);

-- name: GetFolderIDsForConversation :many
SELECT cfi.folder_id
FROM chat_folder_items cfi
         JOIN chat_folders cf ON cf.id = cfi.folder_id
WHERE cfi.conversation_id = $1
  AND cf.user_id = $2;

-- name: GetFolderIDsByConversations :many
SELECT cfi.conversation_id, cfi.folder_id
FROM chat_folder_items cfi
         JOIN chat_folders cf ON cf.id = cfi.folder_id
WHERE cf.user_id = $1;
