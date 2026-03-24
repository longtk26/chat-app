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

-- name: ListMessagesByConversationCursor :many
SELECT m.*, u.username AS sender_name
FROM messages m
INNER JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = $1
  AND m.deleted_at IS NULL
  AND (
    $2::uuid IS NULL
    OR m.id = $2
    OR (
      m.sent_at < COALESCE(
        (
          SELECT m2.sent_at
          FROM messages m2
          WHERE m2.id = $2
            AND m2.conversation_id = $1
            AND m2.deleted_at IS NULL
        ),
        NOW()
      )
      AND ($3::timestamptz IS NULL OR m.sent_at < $3)
    )
  )
ORDER BY m.sent_at DESC
LIMIT $4;

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
