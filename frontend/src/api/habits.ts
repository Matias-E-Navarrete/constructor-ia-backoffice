import { apiFetch } from "./client";

export type Habit = { id: string; name: string; note_slug?: string };
export type HabitStats = {
  total_logs: number;
  completed_logs: number;
  completion_rate: number;
  current_streak: number;
};

export function listHabits(userId?: string) {
  return apiFetch<Habit[]>("/api/habits", { query: { userId } });
}

export function createHabit(name: string, noteSlug?: string) {
  return apiFetch<Habit>("/api/habits", { method: "POST", body: { name, note_slug: noteSlug } });
}

export function checkInHabit(id: string, date?: string) {
  return apiFetch<void>(`/api/habits/${id}/log`, { method: "POST", body: date ? { date } : {} });
}

export function deleteHabit(id: string) {
  return apiFetch<void>(`/api/habits/${id}`, { method: "DELETE" });
}

export function getHabitStats(id: string) {
  return apiFetch<HabitStats>(`/api/habits/${id}/stats`);
}

export type HeatmapDay = { date: string; completed: boolean };
export type PanoramaHabit = {
  habit_id: string;
  name: string;
  current_streak: number;
  best_streak: number;
  heatmap: HeatmapDay[];
  weekly_rates: number[];
};

export function getPanorama() {
  return apiFetch<PanoramaHabit[]>("/api/habits/panorama");
}
