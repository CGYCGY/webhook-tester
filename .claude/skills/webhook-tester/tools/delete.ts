#!/usr/bin/env bun
import { request, requireArg, printJSON } from "./_common";

const uuid = requireArg(Bun.argv.slice(2), 0, "<uuid>");

await request({
  method: "DELETE",
  path: `/api/webhooks/${uuid}`,
  expect: [204, 200],
});

printJSON({ deleted: uuid });
