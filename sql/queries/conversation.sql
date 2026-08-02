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
       peer.id         AS peer_id,
       peer.username   AS peer_username,
       peer.avatar_url AS peer_avatar_url
FROM conversations c
         JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = $1
         LEFT JOIN conversation_members pcm
                   ON pcm.conversation_id = c.id AND pcm.user_id <> $1 AND c.type = 'dm'
         LEFT JOIN users peer ON peer.id = pcm.user_id
ORDER BY c.created_at DESC;

-- name: AddConversationMember :exec
INSERT INTO conversation_members (id, conversation_id, user_id, joined_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (conversation_id, user_id) DO NOTHING;

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
