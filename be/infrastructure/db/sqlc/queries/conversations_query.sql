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
WITH filtered_ids AS (
    SELECT 
        c.*,
        COUNT(*) OVER() AS total_count
    FROM conversations c
    INNER JOIN conversation_participants cp ON cp.conversation_id = c.id
    WHERE cp.user_id = $1 
      AND c.deleted_at IS NULL
    ORDER BY c.updated_at DESC
    LIMIT $2 OFFSET $3
)
SELECT 
    sqlc.embed(c),
    f.total_count,
    -- Aggregate all participants for the conversations found in the CTE
    COALESCE(
        json_agg(
            json_build_object(
                'id', u.id,
                'username', u.username
            )
        ), '[]'
    )::jsonb AS participants
FROM filtered_ids f
INNER JOIN conversations c ON c.id = f.id
INNER JOIN conversation_participants cp_all ON cp_all.conversation_id = c.id
INNER JOIN users u ON u.id = cp_all.user_id
GROUP BY 
    c.id, f.total_count
ORDER BY c.updated_at DESC;

-- name: UpdateConversationLastMessage :exec
UPDATE conversations
SET last_message_id = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: SoftDeleteConversation :exec
UPDATE conversations
SET deleted_at = NOW()
WHERE id = $1;
