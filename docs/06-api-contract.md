# API Contract

## Overview

| Concern | Value |
|---------|-------|
| Base URL | `/` (server-rendered app, no versioned API prefix) |
| Auth | JWT in httpOnly cookie |
| Content-Type (pages) | `text/html` (templ-rendered) |
| Content-Type (webhook endpoint) | `application/json` (responses to external callers) |
| Content-Type (HTMX partials) | `text/html` (HTML fragments for HTMX swaps) |

## Common Headers

| Header | Value | Required |
|--------|-------|----------|
| Cookie | JWT token (set automatically by browser) | Yes (dashboard routes) |
| Content-Type | application/x-www-form-urlencoded | Yes (form submissions) |

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 302 | Redirect (after login, logout, etc.) |
| 400 | Validation error |
| 401 | Unauthorized (no/invalid token) |
| 404 | Resource not found |
| 413 | Request body too large |
| 429 | Rate limit exceeded |
| 500 | Server error |

## Public Routes (no auth)

### ANY `/hook/{uuid}`

Receives incoming webhook requests from external services. Accepts any HTTP method.

**Response (200) — webhook exists:**

```json
{
  "status": "ok",
  "message": "Request captured"
}
```

**Response (404) — webhook not found:**

```json
{
  "status": "error",
  "message": "Webhook not found"
}
```

**Response (429) — rate limited:**

```json
{
  "status": "error",
  "message": "Rate limit exceeded"
}
```

**Response (413) — body too large:**

```json
{
  "status": "error",
  "message": "Request body too large"
}
```

---

### GET `/login`

Renders the login page (HTML).

---

### POST `/login`

Admin login. Sets JWT httpOnly cookie on success.

**Request (form data):**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| email | string | Yes | Valid email format |
| password | string | Yes | Non-empty |

**Success:** 302 redirect to `/dashboard` (set httpOnly JWT cookie)

**Error:** Re-render login page with error message

---

## Authenticated Routes (JWT cookie required)

All routes below redirect to `/login` if no valid JWT cookie is present.

### GET `/dashboard`

Renders the webhook list page (HTML). Shows all webhooks with name, description, URL, request count.

---

### POST `/webhooks`

Create a new webhook. HTMX form submission.

**Request (form data):**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| name | string | Yes | 1-100 characters |
| description | string | No | 0-500 characters |

**Response:** HTMX redirect to `/dashboard` or HTMX swap to append new webhook card.

---

### GET `/webhooks/{uuid}`

Renders webhook detail page with list of captured requests (HTML). SSE connection is initiated from this page for live updates.

---

### PUT `/webhooks/{uuid}`

Update webhook name and description.

**Request (form data):**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| name | string | Yes | 1-100 characters |
| description | string | No | 0-500 characters |

**Response:** HTMX swap to update webhook info section.

---

### DELETE `/webhooks/{uuid}`

Delete webhook and all its captured requests (cascade).

**Response:** HTMX redirect to `/dashboard` or HTMX swap to remove webhook card.

---

### GET `/webhooks/{uuid}/requests/{uuid}`

Renders full request detail page (HTML) — method, URL, headers, body, query params, source IP, timestamp.

---

### GET `/webhooks/{uuid}/sse`

SSE endpoint for real-time updates. Streams new captured requests as HTMX-compatible HTML fragments.

**Event format:**

```
event: new-request
data: <tr hx-swap-oob="afterbegin:#request-list">...request row HTML...</tr>
```

---

### GET `/settings`

Renders settings page with change password form (HTML).

---

### PUT `/settings/password`

Change admin password.

**Request (form data):**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| current_password | string | Yes | Must match current password |
| new_password | string | Yes | Min 8 characters |
| confirm_password | string | Yes | Must match new_password |

**Response:** Re-render settings page with success or error message.

---

### POST `/logout`

Clear JWT cookie, redirect to `/login`.
