# Project Overview

## Problem Statement

Developers integrating with third-party services (Stripe, GitHub, Slack, etc.) need a way to inspect incoming webhook payloads during development and debugging. Existing tools are either SaaS with limitations, require accounts, or are overly complex for simple inspection tasks.

## Target Users

Solo developer (self-hosted tool for personal use). Primarily used during webhook integration development to inspect and debug incoming HTTP requests.

## Goals

- Generate unique webhook endpoints on demand with a name and description
- Capture and store all incoming HTTP requests (any method) to those endpoints
- Display captured requests in a clean, real-time dashboard
- Provide easy one-click copy for webhook URLs, request headers, and request bodies
- Protect against spam and abuse on public-facing webhook endpoints
- Run as a single self-hosted Docker container with zero external dependencies

## Success Criteria

- A webhook URL can be generated and receive requests within seconds of creation
- Captured requests appear in the UI in real-time via SSE (no manual refresh)
- Copy-to-clipboard works reliably for URLs, headers, and bodies
- Single Docker container with SQLite — no external database or services required
- Admin password can be changed via UI, and reset via CLI if forgotten

## Constraints

- Single Go binary + SQLite — no external services or databases
- Self-hosted via Docker
- Single admin user (no multi-user or registration)
- Keep the stack simple: Go + templ + HTMX, no JavaScript build toolchain

## Assumptions

- Low-to-moderate traffic (not designed for high-throughput production webhook routing)
- Single-instance deployment (no horizontal scaling needed)
- Always online (no offline support needed)
- Webhook data is ephemeral/test data — acceptable to lose if DB is reset
