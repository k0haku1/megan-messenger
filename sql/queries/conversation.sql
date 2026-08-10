-- name: CreateConversation :one
INSERT INTO conversations (id, type, title, slug, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetConversation :one
SELECT *
FROM conversations
WHERE id = $1;

-- name: GetConversationBySlug :one
SELECT *
FROM conversations
WHERE slug = $1
LIMIT 1;

-- name: GetUserConversations :many
SELECT c.id,
       c.type,
       c.title,
       c.slug,
       c.created_at,
       peer.id            AS peer_id,
       peer.username      AS peer_username,
       peer.avatar_url    AS peer_avatar_url,
       lm.id              AS last_message_id,
       lm.content         AS last_message_content,
       lm.created_at      AS last_message_created_at,
       lm.sender_id       AS last_message_sender_id,
       lm_sender.username AS last_message_sender_username,
       COALESCE(att.kind, '') AS last_message_attachment_kind,
       (
           SELECT COUNT(*)::int
           FROM messages m
           WHERE m.conversation_id = c.id
             AND m.deleted_at IS NULL
             AND m.sender_id <> $1
             AND (cm.last_read_at IS NULL OR m.created_at > cm.last_read_at)
             AND NOT EXISTS (
               SELECT 1
               FROM hidden_messages hm
               WHERE hm.message_id = m.id
                 AND hm.user_id = $1
             )
       ) AS unread_count,
       (
           SELECT CASE
                      WHEN COUNT(*) = 0 THEN NULL
                      WHEN COUNT(*) FILTER (WHERE last_read_at IS NULL) > 0 THEN NULL
                      ELSE MIN(last_read_at)
                      END
           FROM conversation_members om
           WHERE om.conversation_id = c.id
             AND om.user_id <> $1
       )::timestamptz AS others_read_at
FROM conversations c
         JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = $1
         LEFT JOIN conversation_members pcm
                   ON pcm.conversation_id = c.id AND pcm.user_id <> $1 AND c.type = 'dm'
         LEFT JOIN users peer ON peer.id = pcm.user_id
         LEFT JOIN messages lm ON lm.id = (
    SELECT m.id
    FROM messages m
    WHERE m.conversation_id = c.id
      AND m.deleted_at IS NULL
      AND NOT EXISTS (
        SELECT 1
        FROM hidden_messages hm
        WHERE hm.message_id = m.id
          AND hm.user_id = $1
    )
    ORDER BY m.created_at DESC, m.id DESC
    LIMIT 1
    )
         LEFT JOIN users lm_sender ON lm_sender.id = lm.sender_id
         LEFT JOIN LATERAL (
    SELECT a.kind
    FROM message_attachments a
    WHERE a.message_id = lm.id
    ORDER BY a.created_at ASC, a.id ASC
    LIMIT 1
    ) att ON TRUE
ORDER BY COALESCE(lm.created_at, c.created_at) DESC, c.id DESC;

-- name: AddConversationMember :exec
INSERT INTO conversation_members (id, conversation_id, user_id, joined_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (conversation_id, user_id) DO NOTHING;

-- name: RemoveConversationMember :exec
DELETE
FROM conversation_members
WHERE conversation_id = $1
  AND user_id = $2;

-- name: IsConversationMember :one
SELECT EXISTS (SELECT 1
               FROM conversation_members
               WHERE conversation_id = $1
                 AND user_id = $2);

-- name: FindDMConversation :one
SELECT conversation_id
FROM dm_pairs
WHERE user_low = $1
  AND user_high = $2;

-- name: CreateDMPair :exec
INSERT INTO dm_pairs (user_low, user_high, conversation_id)
VALUES ($1, $2, $3);

-- name: ConversationExists :one
SELECT EXISTS (SELECT 1
               FROM conversations
               WHERE id = $1);

-- name: AdvanceMemberLastReadAt :one
UPDATE conversation_members
SET last_read_at = @last_read_at
WHERE conversation_id = @conversation_id
  AND user_id = @user_id
  AND (last_read_at IS NULL OR last_read_at < @last_read_at)
RETURNING last_read_at;

-- name: GetMemberLastReadAt :one
SELECT last_read_at
FROM conversation_members
WHERE conversation_id = @conversation_id
  AND user_id = @user_id;

-- name: GetOthersReadWatermark :one
SELECT CASE
           WHEN COUNT(*) = 0 THEN NULL
           WHEN COUNT(*) FILTER (WHERE last_read_at IS NULL) > 0 THEN NULL
           ELSE MIN(last_read_at)
           END::timestamptz AS watermark
FROM conversation_members
WHERE conversation_id = @conversation_id
  AND user_id <> @viewer_id;

-- name: ListConversationMemberIDs :many
SELECT user_id
FROM conversation_members
WHERE conversation_id = $1;
