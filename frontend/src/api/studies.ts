import { apiFetch } from "./client";

export type Subject = { id: string; name: string; note_slug?: string };
export type StudySession = {
  id: string;
  subject_id: string;
  session_date: string;
  duration_minutes: number;
  topic: string;
};
export type SubjectOverview = {
  subject_id: string;
  name: string;
  session_count: number;
  total_minutes: number;
  current_streak: number;
};

export function listSubjects() {
  return apiFetch<Subject[]>("/api/studies/subjects");
}

export function createSubject(name: string, noteSlug?: string) {
  return apiFetch<Subject>("/api/studies/subjects", { method: "POST", body: { name, note_slug: noteSlug } });
}

export function deleteSubject(id: string) {
  return apiFetch<void>(`/api/studies/subjects/${id}`, { method: "DELETE" });
}

export function listSessions(subjectId?: string) {
  return apiFetch<StudySession[]>("/api/studies/sessions", { query: { subjectId } });
}

export function logSession(input: { subject_id: string; session_date: string; duration_minutes: number; topic?: string }) {
  return apiFetch<StudySession>("/api/studies/sessions", { method: "POST", body: input });
}

export function getOverview() {
  return apiFetch<SubjectOverview[]>("/api/studies/overview");
}
