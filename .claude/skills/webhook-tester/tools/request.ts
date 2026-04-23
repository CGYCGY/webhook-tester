#!/usr/bin/env bun
import { request, requireArg, printJSON } from "./_common";

const args = Bun.argv.slice(2);
const uuid = requireArg(args, 0, "<uuid>");
const reqID = requireArg(args, 1, "<requestID>");

const { data } = await request({
  method: "GET",
  path: `/api/webhooks/${uuid}/requests/${reqID}`,
});

printJSON(data);
