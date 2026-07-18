import { apiFetch } from "./client";

export type Subtask = { title: string; done: boolean };
export type RepeatRule = "none" | "daily" | "weekly" | "monthly";
export type Priority = "low" | "normal" | "high" | "urgent";
export type Status = "todo" | "in_progress" | "done";
export type Quadrant = "do_now" | "schedule" | "delegate" | "eliminate";

export type Task = {
  id: string;
  title: string;
  description: string;
  category: string;
  subcategory?: string;
  due_date?: string;
  due_time?: string;
  repeat_rule: RepeatRule;
  priority: Priority;
  status: Status;
  quadrant?: Quadrant;
  sort_order: number;
  scheduled_at?: string;
  subtasks: Subtask[];
  note_slug?: string;
  completed_at?: string;
  overdue: boolean;
};

export type ListFilter = { status?: Status; quadrant?: Quadrant; from?: string; to?: string };

export function listTasks(filter: ListFilter = {}) {
  return apiFetch<Task[]>("/api/tasks", {
    query: { status: filter.status, quadrant: filter.quadrant, from: filter.from, to: filter.to },
  });
}

export function createTask(input: {
  title: string;
  description?: string;
  category?: string;
  subcategory?: string;
  due_date?: string;
  due_time?: string;
  repeat_rule?: RepeatRule;
  priority?: Priority;
  subtasks?: Subtask[];
  note_slug?: string;
}) {
  return apiFetch<Task>("/api/tasks", { method: "POST", body: input });
}

export function updateTask(
  id: string,
  input: Partial<{
    title: string;
    description: string;
    category: string;
    subcategory: string;
    due_date: string;
    clear_date: boolean;
    due_time: string;
    repeat_rule: RepeatRule;
    priority: Priority;
  }>,
) {
  return apiFetch<Task>(`/api/tasks/${id}`, { method: "PATCH", body: input });
}

export function deleteTask(id: string) {
  return apiFetch<void>(`/api/tasks/${id}`, { method: "DELETE" });
}

export function completeTask(id: string) {
  return apiFetch<Task>(`/api/tasks/${id}/complete`, { method: "POST" });
}

export function uncompleteTask(id: string) {
  return apiFetch<Task>(`/api/tasks/${id}/uncomplete`, { method: "POST" });
}

export function setTaskStatus(id: string, status: Status) {
  return apiFetch<Task>(`/api/tasks/${id}/status`, { method: "PATCH", body: { status } });
}

export function setTaskQuadrant(id: string, quadrant: Quadrant) {
  return apiFetch<Task>(`/api/tasks/${id}/quadrant`, { method: "PATCH", body: { quadrant } });
}

export function scheduleTask(id: string, scheduledAt: string | null) {
  return apiFetch<Task>(`/api/tasks/${id}/schedule`, { method: "PATCH", body: { scheduled_at: scheduledAt } });
}

export function reorderTasks(orderedIds: string[]) {
  return apiFetch<void>("/api/tasks/reorder", { method: "POST", body: { ordered_ids: orderedIds } });
}
