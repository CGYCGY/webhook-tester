#!/usr/bin/env bun
import { request, requireArg, printJSON, die } from "./_common";

const args = Bun.argv.slice(2);
const uuid = requireArg(args, 0, "<uuid>");

let since: string | undefined;
let timeoutSec = 60;

for (let i = 1; i < args.length; i++) {
  const a = args[i];
  if (a.startsWith("--since=")) {
    since = a.slice("--since=".length);
  } else {
    const n = parseInt(a, 10);
    if (!Number.isFinite(n) || n <= 0) die(3, `invalid timeout: ${a}`);
    timeoutSec = n;
  }
}

// No --since provided → snapshot the current latest id so we only match later events.
if (since === undefined) {
  const { data } = await request<{ requests: { id: string }[] }>({
    method: "GET",
    path: `/api/webhooks/${uuid}/requests`,
    query: { limit: 1 },
  });
  since = data.requests[0]?.id ?? "";
}

const deadline = Date.now() + timeoutSec * 1000;
const minDelay = 500;
const maxDelay = 2000;
let delay = minDelay;

while (Date.now() < deadline) {
  const { data } = await request<{ requests: { id: string }[] }>({
    method: "GET",
    path: `/api/webhooks/${uuid}/requests`,
    query: { limit: 50 },
  });
  const items = data.requests;

  let target: { id: string } | undefined;
  if (items.length > 0) {
    if (since === "") {
      target = items[items.length - 1]; // fresh webhook: oldest captured is "first new"
    } else {
      const idx = items.findIndex((r) => r.id === since);
      if (idx > 0) {
        target = items[idx - 1]; // one newer than anchor
      } else if (idx === -1) {
        target = items[items.length - 1]; // anchor fell off the page: all newer
      }
      // idx === 0 → anchor is still newest, nothing new
    }
  }

  if (target) {
    const detail = await request({
      method: "GET",
      path: `/api/webhooks/${uuid}/requests/${target.id}`,
    });
    printJSON(detail.data);
    process.exit(0);
  }

  const remaining = deadline - Date.now();
  if (remaining <= 0) break;
  await Bun.sleep(Math.min(delay, remaining));
  delay = Math.min(Math.round(delay * 1.5), maxDelay);
}

die(2, `timeout after ${timeoutSec}s waiting for new request on ${uuid}`);
