import { apiFetch } from "./client";

export type User = {
  id: string;
  email: string;
  plan: "free" | "pro";
  role: "user" | "admin";
  planner_start_hour: number;
  planner_end_hour: number;
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

export function updatePlannerHours(startHour: number, endHour: number) {
  return apiFetch<User>("/api/me/planner-hours", {
    method: "PATCH",
    body: { start_hour: startHour, end_hour: endHour },
  });
}
