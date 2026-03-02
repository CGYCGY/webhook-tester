CREATE TABLE IF NOT EXISTS users (
    id          TEXT PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS webhooks (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES users(id),
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhooks_user_id ON webhooks(user_id);

CREATE TABLE IF NOT EXISTS requests (
    id              TEXT PRIMARY KEY,
    webhook_id      TEXT NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    method          TEXT NOT NULL,
    path            TEXT DEFAULT '',
    query_params    TEXT DEFAULT '{}',
    headers         TEXT DEFAULT '{}',
    body            TEXT DEFAULT '',
    content_type    TEXT DEFAULT '',
    source_ip       TEXT DEFAULT '',
    content_length  INTEGER DEFAULT 0,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_requests_webhook_id ON requests(webhook_id);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(webhook_id, created_at DESC);
