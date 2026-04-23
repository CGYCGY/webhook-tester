# Tool Catalog

All tools live in `${CLAUDE_SKILL_DIR}/tools/` and are Bun+TypeScript with a `#!/usr/bin/env bun` shebang. Invoke with `bun <path>` or directly if executable. They share `_common.ts` for config loading, JWT auth, and HTTP.

| Tool | Purpose | Invocation | Sample Output |
|--------|---------|------------|---------------|
| `_common.ts` | Shared: config + JWT cache + `request()` helper. Not called directly. | — | — |
| `login.ts` | Force a fresh JWT login. Auto-runs when token missing/expired; only call manually to debug. | `bun tools/login.ts` | `{"ok":true,"expires_at":"2026-04-23T12:34:56Z"}` |
| `create.ts` | Create a new webhook capture URL. | `bun tools/create.ts <name> [description]` | `{"id":"7f3c...","url":"https://webhooks.example.com/hook/7f3c..."}` |
| `list.ts` | List all webhooks for the account. | `bun tools/list.ts` | `[{"id":"7f3c...","name":"stripe-test","request_count":3}, ...]` |
| `delete.ts` | Delete a webhook by uuid. | `bun tools/delete.ts <uuid>` | `{"deleted":"7f3c..."}` |
| `requests.ts` | List captured requests (metadata only). | `bun tools/requests.ts <uuid> [limit=10] [offset=0]` | `[{"id":"req_abc","method":"POST","received_at":"..."}, ...]` |
| `request.ts` | Full detail for one request: headers, body, query, response sent. | `bun tools/request.ts <uuid> <requestID>` | `{"id":"req_abc","method":"POST","headers":{...},"body":"...","query":{...}}` |
| `set-response.ts` | Configure what the webhook replies with when hit. | `bun tools/set-response.ts <uuid> <status> <content_type> <body>` | `{"updated":true,"status":200,"content_type":"application/json"}` |
| `wait.ts` | Block until the next matching request arrives. Adaptive polling 500ms -> 2s. Exits non-zero on timeout. | `bun tools/wait.ts <uuid> [--since=<reqID>] [timeout_sec=60]` | Same shape as `request.ts` output |

## Notes

- `--since=<reqID>` on `wait.ts`: matches only requests strictly newer than that id. Capture it **before** triggering the service under test so you don't miss fast callbacks.
- `wait.ts` without `--since`: matches any request newer than the moment polling started (still racy — prefer `--since`).
- Non-zero exit codes: 1 = API/auth error, 2 = timeout (for `wait.ts`), 3 = bad args.
- All tools print JSON to stdout and human-readable errors to stderr.
