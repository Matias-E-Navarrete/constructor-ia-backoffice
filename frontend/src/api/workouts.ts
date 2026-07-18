import { apiFetch } from "./client";

export type SetEntry = { exercise_name: string; set_number: number; reps: number; weight_kg: number };
export type Session = {
  id: string;
  session_date: string;
  notes: string;
  note_slug?: string;
  sets: SetEntry[];
};
export type Progress = {
  entries: { session_date: string; set_number: number; reps: number; weight_kg: number }[];
  personal_record_kg: number;
};

export function listSessions(userId?: string) {
  return apiFetch<Session[]>("/api/workouts", { query: { userId } });
}

export function logSession(input: {
  session_date: string;
  notes?: string;
  note_slug?: string;
  sets: SetEntry[];
}) {
  return apiFetch<Session>("/api/workouts", { method: "POST", body: input });
}

export function deleteSession(id: string) {
  return apiFetch<void>(`/api/workouts/${id}`, { method: "DELETE" });
}

export function getProgress(exercise: string) {
  return apiFetch<Progress>("/api/workouts/progress", { query: { exercise } });
}

export type PersonalRecord = { exercise_name: string; personal_record_kg: number };

export function getPersonalRecords() {
  return apiFetch<PersonalRecord[]>("/api/workouts/personal-records");
}

export type Routine = { id: string; name: string; exercises: string[] };

export function listRoutines() {
  return apiFetch<Routine[]>("/api/workouts/routines");
}

export function createRoutine(name: string, exercises: string[]) {
  return apiFetch<Routine>("/api/workouts/routines", { method: "POST", body: { name, exercises } });
}

export function deleteRoutine(id: string) {
  return apiFetch<void>(`/api/workouts/routines/${id}`, { method: "DELETE" });
}
