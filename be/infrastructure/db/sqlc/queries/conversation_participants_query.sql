-- name: AddConversationParticipant :one
INSERT INTO conversation_participants (
    conversation_id, user_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: AddConversationParticipants :many
INSERT INTO conversation_participants (
    conversation_id, user_id
) VALUES (
    $1, UNNEST($2::UUID[])
)
RETURNING *;

-- name: RemoveConversationParticipant :exec
DELETE FROM conversation_participants
WHERE conversation_id = $1
  AND user_id = $2;

-- name: ListParticipantsByConversation :many
SELECT *
FROM conversation_participants
WHERE conversation_id = $1;

-- name: ListConversationUsers :many
SELECT u.id, u.username
FROM conversation_participants cp
INNER JOIN users u ON u.id = cp.user_id
WHERE cp.conversation_id = $1
    AND u.deleted_at IS NULL
ORDER BY u.username ASC;

-- name: IsParticipantInConversation :one
SELECT EXISTS (
    SELECT 1
    FROM conversation_participants
    WHERE conversation_id = $1
      AND user_id = $2
) AS is_participant;
