# System Architecture

## Architecture Diagram

```mermaid
flowchart LR
    subgraph Internet
        ExtService["External Service\n(Stripe, GitHub, etc.)"]
    end

    subgraph Docker["Docker Container"]
        subgraph GoApp["Go Application"]
            PublicRouter["Public Router\n/hook/{uuid}"]
            AuthMiddleware["JWT Auth\nMiddleware"]
            DashboardRouter["Dashboard Router\nPages + Actions"]
            SSEHub["SSE Hub\n(in-memory)"]
            RateLimiter["Rate Limiter\n(in-memory)"]
        end
        SQLite["SQLite\n/data/webhook-tester.db"]
    end

    subgraph Browser
        AdminUI["Admin Dashboard\n(templ + HTMX)"]
    end

    ExtService -->|"ANY /hook/{uuid}"| RateLimiter
    RateLimiter --> PublicRouter
    PublicRouter -->|"Store request"| SQLite
    PublicRouter -->|"Notify"| SSEHub

    AdminUI -->|"Login / Browse"| AuthMiddleware
    AuthMiddleware --> DashboardRouter
    DashboardRouter -->|"Query"| SQLite
    SSEHub -->|"Push new requests"| AdminUI
```

## Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| **Public Router** | Receives incoming webhook requests at `/hook/{uuid}`. Validates webhook exists, stores request in SQLite, notifies SSE hub. Returns JSON response to caller. |
| **Rate Limiter** | In-memory middleware. Per-webhook and per-IP token bucket rate limits. Rejects excess traffic with 429 JSON response. |
| **Auth Middleware** | Validates JWT from httpOnly cookie on all dashboard routes. Redirects to `/login` if token is invalid or expired. |
| **Dashboard Router** | Serves HTML pages via templ and handles form actions (create/edit/delete webhook, change password, logout). |
| **SSE Hub** | Manages active SSE connections per webhook. When a new request is captured, pushes HTMX-compatible HTML fragment to connected clients. |
| **SQLite** | Single file database (`/data/webhook-tester.db`). Stores users, webhooks, and captured requests. Persisted via Docker volume. |

## Request Flow: Webhook Capture

```mermaid
sequenceDiagram
    participant Ext as External Service
    participant RL as Rate Limiter
    participant PR as Public Router
    participant DB as SQLite
    participant SSE as SSE Hub
    participant UI as Admin Browser

    Ext->>RL: POST /hook/{uuid}
    RL->>RL: Check per-webhook + per-IP limits
    alt Rate limit exceeded
        RL-->>Ext: 429 Too Many Requests
    else Within limits
        RL->>PR: Forward request
        PR->>DB: Lookup webhook by UUID
        alt Webhook not found
            PR-->>Ext: 404 Not Found
        else Webhook exists
            PR->>DB: Insert request record
            PR->>DB: Prune oldest if > 100 requests
            PR->>SSE: Notify (webhook_id, request data)
            PR-->>Ext: 200 OK (JSON confirmation)
            SSE->>UI: Push HTML fragment via SSE
        end
    end
```

## Authentication Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Server as Go Server
    participant DB as SQLite

    Browser->>Server: POST /login (email, password)
    Server->>DB: Lookup user by email
    Server->>Server: bcrypt.Compare(password, hash)
    alt Valid credentials
        Server->>Server: Generate JWT (24h expiry)
        Server-->>Browser: Set httpOnly cookie + redirect to /dashboard
    else Invalid credentials
        Server-->>Browser: Re-render login page with error
    end

    Browser->>Server: GET /dashboard (cookie with JWT)
    Server->>Server: Validate JWT from cookie
    alt Valid token
        Server-->>Browser: Render dashboard HTML
    else Invalid or expired token
        Server-->>Browser: Redirect to /login
    end
```

## Environments

| Environment | Purpose | Notes |
|-------------|---------|-------|
| Development | Local dev with Air hot-reload + Tailwind watch | `just dev` starts Air + Tailwind CLI in watch mode |
| Production | Docker container | `just build` creates image, `docker run` with volume for SQLite persistence |

## External Services

None. The application is fully self-contained with no external dependencies.
