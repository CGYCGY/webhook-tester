#!/usr/bin/env bun
import { request, printJSON } from "./_common";

const { data } = await request({ method: "GET", path: "/api/webhooks" });
printJSON(data);
