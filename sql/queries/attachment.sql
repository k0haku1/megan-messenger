-- name: CreateAttachment :one
INSERT INTO message_attachments (
    id, message_id, conversation_id, uploader_id,
    object_key, thumb_key, mime, kind, size_bytes,
    width, height, duration_ms, original_name, created_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14
)
RETURNING *;

-- name: GetAttachmentByID :one
SELECT * FROM message_attachments WHERE id = $1;

-- name: ListAttachmentsByMessageIDs :many
SELECT * FROM message_attachments
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY created_at ASC, id ASC;

-- name: AttachToMessage :exec
UPDATE message_attachments
SET message_id = @message_id
WHERE id = ANY(@attachment_ids::uuid[])
  AND uploader_id = @uploader_id
  AND conversation_id = @conversation_id
  AND message_id IS NULL;

-- name: CountPendingAttachments :one
SELECT COUNT(*)::bigint
FROM message_attachments
WHERE id = ANY(@attachment_ids::uuid[])
  AND uploader_id = @uploader_id
  AND conversation_id = @conversation_id
  AND message_id IS NULL;

-- name: ListConversationMedia :many
SELECT *
FROM message_attachments
WHERE conversation_id = @conversation_id
  AND message_id IS NOT NULL
  AND (
    @kind::text = ''
    OR kind = @kind::text
  )
  AND (
    @cursor_created_at::timestamptz IS NULL
    OR (created_at, id) < (@cursor_created_at::timestamptz, @cursor_id::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT @limit_;
