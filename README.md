# webhook-tester

A self-hosted webhook debugging tool. Inspect incoming webhook requests in real time — headers, body, query params, and more.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-AGPL--3.0-blue?style=flat)

## Features

- Generate unique webhook endpoints instantly
- Inspect requests in real time via Server-Sent Events (no polling)
- View full request details — method, headers, body (JSON/XML syntax highlighted), query params, source IP
- Configure custom responses per webhook (status code, content-type, body)
- Copy endpoint URL, headers, body, or cURL command with one click
- Rate limiting per webhook and per IP
- Single Docker container, zero external dependencies (SQLite)
- Dark / light theme

## Quick Start

```bash
docker run -d \
  -p 8090:8090 \
  -v webhook-data:/data \
  -e JWT_SECRET=change-me \
  -e ADMIN_EMAIL=admin@example.com \
  -e ADMIN_PASSWORD=changeme \
  ghcr.io/cgy/webhook-tester:latest
```

Open `http://localhost:8090` and log in with the credentials above.

## Configuration

All configuration is via environment variables.

| Variable | Required | Default | Description |
|---|---|---|---|
| `JWT_SECRET` | Yes | — | Secret for signing JWT tokens |
| `ADMIN_EMAIL` | No | — | Initial admin email (seeded on first run) |
| `ADMIN_PASSWORD` | No | — | Initial admin password (seeded on first run) |
| `PORT` | No | `8090` | Server listen port |
| `DATA_DIR` | No | `/data` | Directory for the SQLite database |
| `MAX_BODY_SIZE` | No | `1048576` | Max request body size in bytes (default 1 MB) |
| `RATE_LIMIT_PER_WEBHOOK` | No | `60` | Max requests per minute per webhook |
| `RATE_LIMIT_PER_IP` | No | `120` | Max requests per minute per IP |

## API

Machine-readable API documentation is available at `/llms.txt` on any running instance.

### Authentication

```bash
curl -X POST https://your-instance/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "changeme"}'
```

Returns `{"token": "<jwt>"}`. The token expires after 24 hours.

### List Webhooks

```bash
curl https://your-instance/api/webhooks \
  -H "Authorization: Bearer <token>"
```

Returns a JSON array of webhooks with their ID, name, description, hook URL, request count, creation time, and response configuration.

### Send a Request to a Webhook

No authentication required — this is the public capture endpoint:

```bash
curl -X POST https://your-instance/hook/<uuid> \
  -H "Content-Type: application/json" \
  -d '{"event": "test"}'
```

Supports any HTTP method and sub-paths (`/hook/<uuid>/any/path`). The full request (headers, query params, body) is captured and stored.

## Password Reset

If you lose access, reset the admin password directly against the database:

```bash
docker exec -it <container> /app/reset-password \
  --email=admin@example.com \
  --password=newpassword
```

## Development

Prerequisites: Docker, [just](https://github.com/casey/just)

```bash
just up       # start dev container
just dev      # run with hot reload (port 8090)
just generate # regenerate templ, sqlc, Tailwind
just test     # run tests
just build    # build production image
```

Run `just` to see all available commands.

## Tech Stack

- **Go** — Chi router, stdlib SSE, JWT auth
- **templ** — type-safe server-side HTML templates
- **SQLite** — via sqlc for type-safe queries, golang-migrate for migrations
- **HTMX** — dynamic UI without a JS framework
- **Tailwind CSS** — utility-first styling

## License

[AGPL-3.0](LICENSE)
