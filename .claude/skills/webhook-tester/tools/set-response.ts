#!/usr/bin/env bun
import { request, requireArg, printJSON, die } from "./_common";

const args = Bun.argv.slice(2);
const uuid = requireArg(args, 0, "<uuid>");
const statusStr = requireArg(args, 1, "<status>");
const content_type = requireArg(args, 2, "<content_type>");
const body = args[3] ?? "";

const status = parseInt(statusStr, 10);
if (!Number.isFinite(status) || (status !== 0 && (status < 100 || status > 599))) {
  die(3, "status must be 0 or in [100, 599]");
}

const { data } = await request<{
  status: number;
  content_type: string;
  body: string;
}>({
  method: "PATCH",
  path: `/api/webhooks/${uuid}/response`,
  body: { status, content_type, body },
});

printJSON({
  updated: true,
  status: data.status,
  content_type: data.content_type,
});
