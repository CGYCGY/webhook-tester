-- name: CreateWebhook :exec
INSERT INTO webhooks (id, user_id, name, description, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetWebhookByID :one
SELECT id, user_id, name, description, response_config, created_at, updated_at
FROM webhooks
WHERE id = ?
LIMIT 1;

-- name: ListWebhooksByUserID :many
SELECT id, user_id, name, description, response_config, created_at, updated_at
FROM webhooks
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: UpdateWebhook :exec
UPDATE webhooks
SET name = ?, description = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateResponseConfig :exec
UPDATE webhooks
SET response_config = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteWebhook :exec
DELETE FROM webhooks
WHERE id = ?;

-- name: GetWebhookRequestCount :one
SELECT COUNT(*) FROM requests
WHERE webhook_id = ?;
