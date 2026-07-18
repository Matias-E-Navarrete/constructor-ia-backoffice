import { apiFetch } from "./client";

export type User = {
  id: string;
  email: string;
  plan: "free" | "pro";
  role: "user" | "admin";
};

export function register(email: string, password: string) {
  return apiFetch<User>("/api/auth/register", { method: "POST", body: { email, password } });
}

export function login(email: string, password: string) {
  return apiFetch<{ token: string; user: User }>("/api/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export function me() {
  return apiFetch<User>("/api/me");
}

export function upgrade() {
  return apiFetch<User>("/api/me/upgrade", { method: "POST" });
}
