const TOKEN_KEY = "rimu_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null) {
  if (token) localStorage.setItem(TOKEN_KEY, token);
  else localStorage.removeItem(TOKEN_KEY);
}

export class ApiError extends Error {
  status: number;
  code: string;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

export class UpgradeRequiredError extends ApiError {}
export class FeatureDisabledError extends ApiError {}

type RequestOptions = {
  method?: string;
  body?: unknown;
  query?: Record<string, string | undefined>;
};

export async function apiFetch<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  let url = path;
  if (options.query) {
    const params = new URLSearchParams();
    for (const [key, value] of Object.entries(options.query)) {
      if (value !== undefined && value !== "") params.set(key, value);
    }
    const qs = params.toString();
    if (qs) url += `?${qs}`;
  }

  const res = await fetch(url, {
    method: options.method ?? "GET",
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  });

  if (res.status === 401) {
    setToken(null);
    throw new ApiError(401, "unauthorized", "session expired");
  }

  const contentType = res.headers.get("content-type") ?? "";
  const payload = contentType.includes("application/json") ? await res.json() : undefined;

  if (!res.ok) {
    const code = payload?.error ?? "error";
    const message = payload?.message ?? res.statusText;
    if (res.status === 403 && code === "upgrade_required") {
      throw new UpgradeRequiredError(403, code, message);
    }
    if (res.status === 503 && code === "feature_disabled") {
      throw new FeatureDisabledError(503, code, message);
    }
    throw new ApiError(res.status, code, message);
  }

  return payload as T;
}

// downloadFile is used by the CSV/PDF/zip export endpoints, which return raw
// bytes rather than JSON. Pass `body` to issue a POST with a JSON payload
// (e.g. PDF export filters) instead of a plain GET.
export async function downloadFile(path: string, filename: string, body?: unknown) {
  const token = getToken();
  const headers: Record<string, string> = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;
  if (body !== undefined) headers["Content-Type"] = "application/json";

  const res = await fetch(path, {
    method: body !== undefined ? "POST" : "GET",
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) {
    if (res.status === 403) throw new UpgradeRequiredError(403, "upgrade_required", "this feature requires the pro plan");
    throw new ApiError(res.status, "error", "download failed");
  }

  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
