import { apiFetch, downloadFile } from "./client";

export type Note = { slug: string; title: string; body: string; kind: "freeform" | "profile"; links: string[] };
export type GraphNode = { id: string; title: string; kind: "note" | "habit" | "workout" | "finance" };
export type GraphEdge = { from: string; to: string };
export type Graph = { nodes: GraphNode[]; edges: GraphEdge[] };

export function listNotes() {
  return apiFetch<Note[]>("/api/notes");
}

export function createNote(title: string, body: string, kind?: "freeform" | "profile") {
  return apiFetch<Note>("/api/notes", { method: "POST", body: { title, body, kind } });
}

export function updateNote(slug: string, title: string, body: string) {
  return apiFetch<Note>(`/api/notes/${slug}`, { method: "PATCH", body: { title, body } });
}

export function deleteNote(slug: string) {
  return apiFetch<void>(`/api/notes/${slug}`, { method: "DELETE" });
}

export function getGraph() {
  return apiFetch<Graph>("/api/notes/graph");
}

export function exportVault() {
  return downloadFile("/api/notes/export", "vault.zip");
}
