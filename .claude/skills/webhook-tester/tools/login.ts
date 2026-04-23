#!/usr/bin/env bun
import { forceLogin, printJSON } from "./_common";

const tc = await forceLogin();
printJSON({ ok: true, expires_at: tc.expires_at });
