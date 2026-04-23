import { dirname, join } from "node:path";

const TOOLS_DIR = import.meta.dir;
export const SKILL_DIR = dirname(TOOLS_DIR);
export const CONFIG_PATH = join(SKILL_DIR, "config.json");
export const CONFIG_EXAMPLE_PATH = join(SKILL_DIR, "config.example.json");
export const TOKEN_PATH = join(SKILL_DIR, ".token.json");

export interface Config {
  base_url: string;
  email: string;
  password: string;
}

interface TokenCache {
  token: string;
  expires_at: string;
}

export function die(code: number, msg: string): never {
  process.stderr.write(msg.endsWith("\n") ? msg : msg + "\n");
  process.exit(code);
}

export function printJSON(obj: unknown): void {
  process.stdout.write(JSON.stringify(obj) + "\n");
}

export function requireArg(args: string[], index: number, name: string): string {
  const v = args[index];
  if (!v) die(3, `missing argument: ${name}`);
  return v;
}

let cachedConfig: Config | undefined;

export async function loadConfig(): Promise<Config> {
  if (cachedConfig) return cachedConfig;
  const file = Bun.file(CONFIG_PATH);
  if (!(await file.exists())) {
    die(
      1,
      `config.json not found at ${CONFIG_PATH}\nCopy config.example.json to config.json and fill in base_url, email, password.`,
    );
  }
  let cfg: Config;
  try {
    cfg = (await file.json()) as Config;
  } catch (e) {
    die(1, `config.json is not valid JSON: ${(e as Error).message}`);
  }
  if (!cfg.base_url || !cfg.email || !cfg.password) {
    die(1, `config.json is missing required fields (base_url, email, password)`);
  }
  cfg.base_url = cfg.base_url.replace(/\/+$/, "");
  cachedConfig = cfg;
  return cfg;
}

async function readTokenCache(): Promise<TokenCache | null> {
  const file = Bun.file(TOKEN_PATH);
  if (!(await file.exists())) return null;
  try {
    return (await file.json()) as TokenCache;
  } catch {
    return null;
  }
}

async function writeTokenCache(tc: TokenCache): Promise<void> {
  await Bun.write(TOKEN_PATH, JSON.stringify(tc, null, 2));
}

async function doLogin(cfg: Config): Promise<TokenCache> {
  const res = await fetch(`${cfg.base_url}/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email: cfg.email, password: cfg.password }),
  });
  if (!res.ok) {
    let msg = `login failed: ${res.status}`;
    try {
      const b = (await res.json()) as { error?: string };
      if (b.error) msg += ` - ${b.error}`;
    } catch {}
    die(1, msg);
  }
  const j = (await res.json()) as { token?: string };
  if (!j.token) die(1, `login response missing token`);
  // Server token TTL is 24h; refresh ~1h early to avoid mid-request expiry.
  const expires_at = new Date(Date.now() + 23 * 3600 * 1000).toISOString();
  const tc: TokenCache = { token: j.token, expires_at };
  await writeTokenCache(tc);
  return tc;
}

export async function forceLogin(): Promise<TokenCache> {
  const cfg = await loadConfig();
  return doLogin(cfg);
}

async function getToken(cfg: Config, force = false): Promise<string> {
  if (!force) {
    const cached = await readTokenCache();
    if (cached && Date.parse(cached.expires_at) > Date.now()) return cached.token;
  }
  return (await doLogin(cfg)).token;
}

export interface RequestOpts {
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  path: string;
  body?: unknown;
  query?: Record<string, string | number>;
  expect?: number[];
}

export async function request<T = unknown>(
  opts: RequestOpts,
): Promise<{ status: number; data: T }> {
  const cfg = await loadConfig();
  let token = await getToken(cfg);
  let url = cfg.base_url + opts.path;
  if (opts.query && Object.keys(opts.query).length > 0) {
    const qs = new URLSearchParams();
    for (const [k, v] of Object.entries(opts.query)) qs.set(k, String(v));
    url += `?${qs.toString()}`;
  }
  let res = await doFetch(url, opts, token);
  if (res.status === 401) {
    token = await getToken(cfg, true);
    res = await doFetch(url, opts, token);
  }
  const ok = opts.expect
    ? opts.expect.includes(res.status)
    : res.status >= 200 && res.status < 300;
  if (!ok) {
    let msg = `${opts.method} ${opts.path} -> ${res.status}`;
    try {
      const b = (await res.json()) as { error?: string };
      if (b.error) msg += `: ${b.error}`;
    } catch {}
    die(1, msg);
  }
  let data: T;
  const text = await res.text();
  if (!text) {
    data = undefined as unknown as T;
  } else {
    try {
      data = JSON.parse(text) as T;
    } catch (e) {
      die(1, `invalid JSON response from ${opts.path}: ${(e as Error).message}`);
    }
  }
  return { status: res.status, data };
}

async function doFetch(
  url: string,
  opts: RequestOpts,
  token: string,
): Promise<Response> {
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
    Accept: "application/json",
  };
  let body: string | undefined;
  if (opts.body !== undefined) {
    headers["Content-Type"] = "application/json";
    body = JSON.stringify(opts.body);
  }
  return fetch(url, { method: opts.method, headers, body });
}
