import { apiFetch } from "./client";

export type RoadmapItem = {
  id: string;
  title: string;
  description: string;
  kind: "feature" | "bug" | "task";
  status: "planned" | "in_progress" | "done";
};
export type Flag = { key: string; enabled: boolean; description: string };
export type AdminUser = { id: string; email: string; plan: "free" | "pro"; role: "user" | "admin" };
export type TestRunResult = { success: boolean; passed: number; failed: number; output: string };

export function listRoadmap() {
  return apiFetch<RoadmapItem[]>("/api/admin/roadmap");
}

export function createRoadmapItem(title: string, description: string, kind: RoadmapItem["kind"]) {
  return apiFetch<RoadmapItem>("/api/admin/roadmap", { method: "POST", body: { title, description, kind } });
}

export function updateRoadmapStatus(id: string, status: RoadmapItem["status"]) {
  return apiFetch<RoadmapItem>(`/api/admin/roadmap/${id}`, { method: "PATCH", body: { status } });
}

export function listFeatureFlags() {
  return apiFetch<Flag[]>("/api/admin/feature-flags");
}

export function toggleFeatureFlag(key: string, enabled: boolean) {
  return apiFetch<Flag>(`/api/admin/feature-flags/${key}`, { method: "PATCH", body: { enabled } });
}

export function listUsers() {
  return apiFetch<AdminUser[]>("/api/admin/users");
}

export function setUserPlan(id: string, plan: "free" | "pro") {
  return apiFetch<AdminUser>(`/api/admin/users/${id}/plan`, { method: "PATCH", body: { plan } });
}

export function setUserRole(id: string, role: "user" | "admin") {
  return apiFetch<AdminUser>(`/api/admin/users/${id}/role`, { method: "PATCH", body: { role } });
}

export function runTests(suite: "backend" | "frontend") {
  return apiFetch<TestRunResult>("/api/admin/tests/run", { method: "POST", query: { suite } });
}
