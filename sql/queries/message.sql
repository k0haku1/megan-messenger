-- name: CreateMessage :one
INSERT INTO messages (id, conversation_id, sender_id, content, reply_to_id, forwarded_from_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMessageByID :one
SELECT m.id,
       m.conversation_id,
       m.sender_id,
       m.content,
       m.reply_to_id,
       m.forwarded_from_id,
       m.deleted_at,
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
  AND m.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM hidden_messages hm WHERE hm.message_id = m.id AND hm.user_id = @user_id::uuid)
  AND (
    @cursor_created_at::timestamptz IS NULL
        OR
    (m.created_at, m.id) < (@cursor_created_at::timestamptz, @cursor_id::uuid)
    )
ORDER BY m.created_at DESC, m.id DESC
LIMIT @limit_;

-- name: HideMessageForUser :exec
INSERT INTO hidden_messages (message_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteMessageForEveryone :one
UPDATE messages
SET deleted_at = NOW(), content = ''
WHERE id = $1 AND sender_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: AddMessageReaction :exec
INSERT INTO message_reactions (message_id, user_id, emoji)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveMessageReaction :exec
DELETE FROM message_reactions
WHERE message_id = $1 AND user_id = $2 AND emoji = $3;

-- name: GetMessageReactions :many
SELECT emoji, COUNT(*)::bigint AS count, BOOL_OR(user_id = $2) AS reacted_by_me
FROM message_reactions
WHERE message_id = $1
GROUP BY emoji
ORDER BY emoji;

-- name: GetMessageReactionUsers :many
SELECT mr.emoji, u.id, u.username, u.avatar_url
FROM message_reactions mr
JOIN users u ON u.id = mr.user_id
WHERE mr.message_id = $1
ORDER BY mr.emoji, u.username;
