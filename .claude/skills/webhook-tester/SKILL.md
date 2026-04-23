---
name: webhook-tester
description: JSON API client tools for the webhook-tester service. Use when the user wants to create a webhook capture URL, list or inspect captured requests, configure a custom response, or wait for an incoming callback. Triggers on "test a webhook", "capture webhook callback", "wait for webhook", "verify HTTP callback", or any mention of a webhook-tester URL.
allowed-tools: Bash, Read, Write
user-invocable: true
---

# webhook-tester

Tools for interacting with a deployed webhook-tester instance via its JSON API. Use these when you need to capture, inspect, or respond to HTTP requests as part of a webhook test.

## Files

SKILL_DIR: ${CLAUDE_SKILL_DIR}
TOOLS: ${CLAUDE_SKILL_DIR}/tools
CONFIG: ${CLAUDE_SKILL_DIR}/config.json
CONFIG_EXAMPLE: ${CLAUDE_SKILL_DIR}/config.example.json
TOKEN_CACHE: ${CLAUDE_SKILL_DIR}/.token.json

## Tools

See `reference/tools.md` for the full catalog (login, create, list, delete, requests, request, set-response, wait) — invocations, args, sample output.

## Sample report

```
Webhook: https://webhooks.example.com/hook/7f3c1a...
Captured: POST / (event=checkout.session.completed) ✓
Cleanup: deleted
```
