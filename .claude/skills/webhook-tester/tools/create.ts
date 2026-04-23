#!/usr/bin/env bun
import { request, requireArg, printJSON } from "./_common";

const args = Bun.argv.slice(2);
const name = requireArg(args, 0, "<name>");
const description = args[1] ?? "";

const { data } = await request<{ id: string; url: string }>({
  method: "POST",
  path: "/api/webhooks",
  body: { name, description },
});

printJSON({ id: data.id, url: data.url });
