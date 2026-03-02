# Component Architecture — Web

## Platform

Web (single platform, server-rendered with templ + HTMX)

## Pages / Screens

| Page | Purpose |
|------|---------|
| Login | Email/password form for admin login |
| Dashboard | List all webhooks with name, description, URL, copy button, request count |
| Webhook Detail | View captured requests for a webhook with live SSE updates |
| Request Detail | Full inspection of a single captured request |
| Settings | Change password form |

## Routing

| Path | Page | Auth Required |
|------|------|---------------|
| `/login` | Login | No |
| `/dashboard` | Dashboard (webhook list) | Yes |
| `/webhooks/{uuid}` | Webhook Detail | Yes |
| `/webhooks/{uuid}/requests/{uuid}` | Request Detail | Yes |
| `/settings` | Settings | Yes |

## Shared Components (templ)

Components are reusable templ functions composed into pages.

| Component | Purpose |
|-----------|---------|
| `layout` | Base HTML layout — head, nav, footer, dark mode toggle, Tailwind CSS, HTMX script tags |
| `navbar` | Top navigation bar — logo/app name, link to dashboard, settings link, logout button |
| `copyButton` | Reusable copy-to-clipboard button — accepts target text, shows "Copied!" feedback for 2 seconds |
| `webhookCard` | Single webhook row/card — name, description, URL with copy button, request count badge, edit/delete actions |
| `requestRow` | Single request row in webhook detail — method badge, path, timestamp, content length |
| `requestDetail` | Full request view — method + URL, headers table, body with syntax highlighting, query params table |
| `emptyState` | Placeholder for empty lists — "No webhooks yet" or "No requests captured" with contextual CTA |
| `toast` | Flash message component for success/error notifications — auto-dismiss after 3 seconds |
| `confirmModal` | Confirmation dialog for destructive actions (delete webhook) |
| `themeToggle` | Dark/light mode toggle button — sun/moon icon in navbar |
| `methodBadge` | Colored badge for HTTP method — GET=green, POST=blue, PUT=yellow, PATCH=purple, DELETE=red |
| `createWebhookForm` | Inline form or modal for creating a new webhook — name + description fields |

## State Management

### Approach

No client-side state management framework. This is a server-rendered application:

- **Page state** lives on the server — templ renders full pages from database queries
- **Interactivity** handled by HTMX — partial page swaps via HTML-over-the-wire
- **Real-time** handled by SSE → HTMX swap — server pushes HTML fragments to the browser
- **Theme preference** stored in `localStorage` — applied via Tailwind `dark` class with vanilla JS

### State Location

| State | Where | Mechanism |
|-------|-------|-----------|
| Auth session | httpOnly cookie | JWT token, validated server-side on each request |
| Webhook/request data | SQLite → server-rendered HTML | HTMX fetches partial HTML from server endpoints |
| Theme preference | localStorage | Vanilla JS reads on page load, toggles `dark` class on `<html>` |
| Clipboard feedback | DOM | Vanilla JS shows "Copied!" text for 2 seconds, then reverts |
| Flash messages | Server → rendered HTML | Toast component rendered by server after form actions |
