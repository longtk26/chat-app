-- name: CreateConversation :one
INSERT INTO conversations (
    title, type
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetConversationByID :one
SELECT *
FROM conversations
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListConversationsByUser :many
SELECT c.*
FROM conversations c
INNER JOIN conversation_participants cp ON cp.conversation_id = c.id
WHERE cp.user_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.updated_at DESC;

-- name: ListConversationsByUserWithPagination :many
SELECT 
    c.*,
    COUNT(*) OVER() AS total_conversations
FROM conversations c
INNER JOIN conversation_participants cp 
    ON cp.conversation_id = c.id
WHERE cp.user_id = $1
  AND c.deleted_at IS NULL
ORDER BY c.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateConversationLastMessage :exec
UPDATE conversations
SET last_message_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: SoftDeleteConversation :exec
UPDATE conversations
SET deleted_at = NOW()
WHERE id = $1;
