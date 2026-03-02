# Feature List

## Features

| # | Feature | Description | Scope | Priority |
|---|---------|-------------|-------|----------|
| 1 | Admin seed | Seed single admin user on first startup from `ADMIN_EMAIL` + `ADMIN_PASSWORD` env vars. Seed only on first run (skip if user exists). Provide `reset-password` CLI command for recovery. | BE | Must |
| 2 | Login | Email/password login form. Sets JWT in httpOnly cookie on success. | Both | Must |
| 3 | Change password | Admin can change password from settings page (current password + new password + confirm). | Both | Must |
| 4 | JWT auth middleware | Protect all dashboard routes. Validate JWT from cookie. Redirect to login if invalid/expired. | BE | Must |
| 5 | Create webhook | Generate new webhook with user-provided name and description. Auto-assign UUID for public URL. | Both | Must |
| 6 | List webhooks | Dashboard page showing all webhooks with name, description, URL with copy button, and request count. | Both | Must |
| 7 | View webhook detail | Page showing captured requests for a specific webhook with live SSE updates. Displays webhook info, URL with copy button, and request list. | Both | Must |
| 8 | View request detail | Full request inspection: HTTP method, headers table, body (with syntax highlighting for JSON), query params, source IP, timestamp. Copy buttons for each section. | Both | Must |
| 9 | Delete webhook | Delete a webhook and all its captured requests (cascade). Requires confirmation dialog. | Both | Must |
| 10 | Copy to clipboard | One-click copy for webhook URL, request headers (as JSON), request body (raw), and full request as cURL command. Visual "Copied!" feedback. | FE | Must |
| 11 | Receive webhook requests | Public endpoint `/hook/{uuid}` accepts any HTTP method (GET, POST, PUT, PATCH, DELETE, etc.). Captures method, headers, body, query params, source IP, content type. Returns JSON confirmation. | BE | Must |
| 12 | Handle non-existent webhook | Return clean 404 JSON response for requests to unknown webhook UUIDs. Do not leak information. | BE | Must |
| 13 | SSE real-time updates | Push new captured requests to the webhook detail page in real-time via Server-Sent Events. HTMX SSE extension prepends new request rows without page refresh. | Both | Must |
| 14 | Rate limiting (per-webhook) | In-memory token bucket rate limit per webhook URL (e.g., 60 requests/min). Returns 429 when exceeded. | BE | Must |
| 15 | Rate limiting (per-IP) | Global per-IP rate limit to prevent a single source from flooding all webhooks. In-memory, resets on restart. | BE | Must |
| 16 | Max body size | Reject request bodies over 1MB with 413 status code. Configurable via env var. | BE | Must |
| 17 | Request retention limit | Keep last 100 requests per webhook. Auto-prune oldest requests when limit exceeded. | BE | Should |
| 18 | Edit webhook | Update webhook name and description from the webhook detail page. | Both | Should |
| 19 | Webhook request count | Show request count badge on each webhook card in the dashboard list. | Both | Should |
| 20 | Request body syntax highlight | Pretty-print JSON and XML bodies in request detail view. | FE | Should |
| 21 | Dark/light theme toggle | Toggle between dark and light theme. Preference stored in localStorage. Applied via Tailwind `dark` class. | FE | Must |
| 22 | Auto-expiry | Auto-delete webhooks inactive for N days (configurable via env var). | BE | Could |
| 23 | Export requests | Download all captured requests for a webhook as JSON file. | Both | Could |
| 24 | Search/filter requests | Filter requests by HTTP method, date range, or keyword in body/headers. | Both | Could |
| 25 | Multi-user support | Registration, multiple users with separate data. | Both | Won't |
| 26 | Custom webhook slugs | User-chosen URL paths instead of UUID. | Both | Won't |

## Priority Legend

| Priority | Meaning |
|----------|---------|
| Must | Required for MVP, cannot launch without |
| Should | Important, include if possible |
| Could | Nice to have, if time permits |
| Won't | Explicitly out of scope for this project |
