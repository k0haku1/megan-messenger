-- name: CreateMessage :one
INSERT INTO messages (id, conversation_id, sender_id, content, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMessageByID :one
SELECT m.id,
       m.conversation_id,
       m.sender_id,
       m.content,
       m.created_at,
       u.username,
       u.avatar_url
FROM messages m
         JOIN users u ON u.id = m.sender_id
WHERE m.id = $1;

-- name: GetMessagesPaging :many
SELECT m.*, u.*
FROM messages m
         JOIN users u ON u.id = m.sender_id
WHERE m.conversation_id = @conversation_id::uuid
  AND (
    @cursor_created_at::timestamptz IS NULL
        OR
    (m.created_at, m.id) < (@cursor_created_at::timestamptz, @cursor_id::uuid)
    )
ORDER BY m.created_at DESC, m.id DESC
LIMIT @limit_;
