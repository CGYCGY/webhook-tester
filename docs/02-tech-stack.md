# Tech Stack

## Backend

| Concern | Choice | Reason |
|---------|--------|--------|
| Language | Go | Fast, single binary, excellent HTTP stdlib |
| Router | Chi | Lightweight, idiomatic Go router with middleware support |
| Templates | templ | Type-safe Go templates, compiles to Go code |
| DB Access | sqlc | Type-safe SQL — generates Go code from queries, no runtime overhead |
| Migration | golang-migrate | Standard Go migration tool, works well with SQLite |
| Auth | JWT (self-issued) | Stateless tokens, bcrypt password hashing, single admin user |
| Real-time | SSE (stdlib) | Server-Sent Events via standard net/http, no extra dependencies |
| Testing | Go stdlib (testing) | Built-in testing framework |

## Frontend (server-rendered from Go)

| Concern | Choice | Reason |
|---------|--------|--------|
| Rendering | templ (server-side) | HTML rendered on server, no JS build step |
| Interactivity | HTMX | Declarative AJAX, SSE support built-in via extensions |
| Styling | Tailwind CSS (standalone CLI) | No Node.js needed — standalone binary generates CSS |
| Copy/Clipboard | Vanilla JS | A few lines of inline JS for copy buttons |
| Theme Toggle | Vanilla JS + Tailwind dark mode | localStorage preference, Tailwind `dark` class toggle |
| Icons | Heroicons (SVG) | Copy-paste SVGs into templates, no dependencies |

## Database

| Concern | Choice | Reason |
|---------|--------|--------|
| Type | Relational (SQL) | Structured data with clear relationships |
| Database | SQLite (mattn/go-sqlite3) | Zero config, file-based, perfect for self-hosted single-instance |
| Storage | Docker volume mount | `/data/webhook-tester.db` persists across container restarts |
| Seeding | Auto-seed admin on first run | From `ADMIN_EMAIL` + `ADMIN_PASSWORD` env vars |

## Infrastructure

| Concern | Choice | Reason |
|---------|--------|--------|
| Deployment | Docker (single container) | Multi-stage build, final image ~20MB |
| Build | Multi-stage Dockerfile | Build Go + Tailwind CSS, output single small image |
| Environments | Dev (local) + Prod (Docker) | Two environments, keep it simple |

## Dev Tools

| Tool | Purpose |
|------|---------|
| Air | Go hot-reload during development |
| Tailwind CLI (standalone) | CSS generation without Node.js |
| sqlc | Generate Go code from SQL queries |
| golang-migrate | Run database migrations |
| Just | Task runner for common commands (dev, build, migrate, seed) |
| Docker Compose | Local dev environment orchestration |
