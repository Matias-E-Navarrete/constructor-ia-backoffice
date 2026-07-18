import { apiFetch } from "./client";

export type InboxItem = { id: string; content: string; pinned: boolean };

export function listInbox() {
  return apiFetch<InboxItem[]>("/api/inbox");
}

export function captureItem(content: string) {
  return apiFetch<InboxItem>("/api/inbox", { method: "POST", body: { content } });
}

export function setPinned(id: string, pinned: boolean) {
  return apiFetch<void>(`/api/inbox/${id}/pin`, { method: "PATCH", body: { pinned } });
}

export function deleteItem(id: string) {
  return apiFetch<void>(`/api/inbox/${id}`, { method: "DELETE" });
}
