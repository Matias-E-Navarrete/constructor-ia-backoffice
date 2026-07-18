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

// downloadFile is used by the CSV/zip export endpoints, which return raw
// bytes rather than JSON.
export async function downloadFile(path: string, filename: string) {
  const token = getToken();
  const headers: Record<string, string> = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(path, { headers });
  if (!res.ok) throw new ApiError(res.status, "error", "download failed");

  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
