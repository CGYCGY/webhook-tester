-- name: CreateRequest :exec
INSERT INTO requests (id, webhook_id, method, path, query_params, headers, body, content_type, source_ip, content_length, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRequestByID :one
SELECT id, webhook_id, method, path, query_params, headers, body, content_type, source_ip, content_length, created_at
FROM requests
WHERE id = ?
LIMIT 1;

-- name: ListRequestsByWebhookID :many
SELECT id, webhook_id, method, path, query_params, headers, body, content_type, source_ip, content_length, created_at
FROM requests
WHERE webhook_id = ?
ORDER BY created_at DESC
LIMIT ?;

-- name: DeleteOldRequests :exec
DELETE FROM requests
WHERE id IN (
  SELECT r.id FROM requests r
  WHERE r.webhook_id = ?
  ORDER BY r.created_at DESC
  LIMIT -1 OFFSET 100
);
