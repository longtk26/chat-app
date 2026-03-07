-- name: CreateMessage :one
INSERT INTO messages (
    sender_id, recipient_id, conversation_id, content
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetMessageByID :one
SELECT *
FROM messages
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListMessagesByConversation :many
SELECT *
FROM messages
WHERE conversation_id = $1
  AND deleted_at IS NULL
ORDER BY sent_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMessageContent :one
UPDATE messages
SET content = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteMessage :exec
UPDATE messages
SET deleted_at = NOW()
WHERE id = $1;
