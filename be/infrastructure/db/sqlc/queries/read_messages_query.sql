-- name: MarkMessageAsRead :one
INSERT INTO read_messages (
    user_id, message_id
) VALUES (
    $1, $2
)
ON CONFLICT (user_id, message_id) DO NOTHING
RETURNING *;

-- name: GetUnreadMessageCount :one
SELECT COUNT(m.id) AS unread_count
FROM messages m
WHERE m.conversation_id = $1
  AND m.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM read_messages rm
    WHERE rm.message_id = m.id
      AND rm.user_id = $2
  );

-- name: ListReadMessagesByUserInConversation :many
SELECT rm.*
FROM read_messages rm
INNER JOIN messages m ON m.id = rm.message_id
WHERE rm.user_id = $1
  AND m.conversation_id = $2;
