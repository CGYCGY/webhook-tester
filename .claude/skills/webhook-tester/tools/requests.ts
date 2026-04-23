#!/usr/bin/env bun
import { request, requireArg, printJSON, die } from "./_common";

const args = Bun.argv.slice(2);
const uuid = requireArg(args, 0, "<uuid>");
const limit = args[1] ? parseInt(args[1], 10) : 10;
const offset = args[2] ? parseInt(args[2], 10) : 0;

if (!Number.isFinite(limit) || limit < 1 || limit > 500) {
  die(3, "limit must be between 1 and 500");
}
if (!Number.isFinite(offset) || offset < 0) {
  die(3, "offset must be >= 0");
}

const { data } = await request<{ requests: unknown[] }>({
  method: "GET",
  path: `/api/webhooks/${uuid}/requests`,
  query: { limit, offset },
});

printJSON(data.requests);
