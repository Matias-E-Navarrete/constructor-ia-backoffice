import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import * as tasksApi from "../../api/tasks";
import * as inboxApi from "../../api/inbox";
import { FeatureDisabledError } from "../../api/client";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import { categoryColor } from "../../lib/categoryColor";
import QuickAddTask from "../../components/QuickAddTask";

type View = "enfoque" | "eisenhower" | "kanban";

const quadrants: { key: tasksApi.Quadrant; label: string; dot: string }[] = [
  { key: "do_now", label: "Hacer ahora", dot: "#ef4444" },
  { key: "schedule", label: "Programar", dot: "#3b82f6" },
  { key: "delegate", label: "Delegar", dot: "#f59e0b" },
  { key: "eliminate", label: "Evaluar", dot: "#737373" },
];

const columns: { key: tasksApi.Status; label: string; dot: string }[] = [
  { key: "todo", label: "Por hacer", dot: "#737373" },
  { key: "in_progress", label: "En progreso", dot: "#3b82f6" },
  { key: "done", label: "Hecho", dot: "#22c55e" },
];

export default function Tasks() {
  const location = useLocation();
  const inboxState = location.state as { prefillTitle?: string; fromInboxId?: string } | null;

  const [view, setView] = useState<View>("enfoque");
  const [tasks, setTasks] = useState<tasksApi.Task[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [dragId, setDragId] = useState<string | null>(null);

  async function handleQuickAddCreated() {
    if (inboxState?.fromInboxId) {
      await inboxApi.deleteItem(inboxState.fromInboxId);
      window.history.replaceState({}, "");
    }
    load();
  }

  async function load() {
    try {
      setTasks(await tasksApi.listTasks());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleComplete(id: string, done: boolean) {
    if (done) await tasksApi.uncompleteTask(id);
    else await tasksApi.completeTask(id);
    load();
  }

  async function handleDropReorder(targetId: string) {
    if (!dragId || dragId === targetId) return;
    const ids = tasks.map((t) => t.id);
    const from = ids.indexOf(dragId);
    const to = ids.indexOf(targetId);
    ids.splice(to, 0, ids.splice(from, 1)[0]);
    setTasks(ids.map((id) => tasks.find((t) => t.id === id)!));
    await tasksApi.reorderTasks(ids);
    setDragId(null);
  }

  async function handleDropQuadrant(quadrant: tasksApi.Quadrant) {
    if (!dragId) return;
    await tasksApi.setTaskQuadrant(dragId, quadrant);
    setDragId(null);
    load();
  }

  async function handleDropStatus(status: tasksApi.Status) {
    if (!dragId) return;
    await tasksApi.setTaskStatus(dragId, status);
    setDragId(null);
    load();
  }

  if (disabled) return <FeatureDisabledBanner />;

  const overdue = tasks.filter((t) => t.overdue);
  const pending = tasks.filter((t) => t.status !== "done");

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Tareas</h1>
        <div className="flex gap-1 bg-neutral-100 dark:bg-neutral-900 rounded-md p-1">
          {(["enfoque", "eisenhower", "kanban"] as View[]).map((v) => (
            <button
              key={v}
              onClick={() => setView(v)}
              className={`px-3 py-1 text-sm rounded ${
                view === v
                  ? "bg-white dark:bg-neutral-700 text-neutral-900 dark:text-neutral-50 shadow-sm"
                  : "text-neutral-500 dark:text-neutral-400"
              }`}
            >
              {v === "enfoque" ? "Enfoque" : v === "eisenhower" ? "Eisenhower" : "Kanban"}
            </button>
          ))}
        </div>
      </div>

      <QuickAddTask onCreated={handleQuickAddCreated} initialTitle={inboxState?.prefillTitle} />

      {overdue.length > 0 && (
        <div className="mb-4 text-xs text-red-500 dark:text-red-400">
          {overdue.length} tarea{overdue.length > 1 ? "s" : ""} vencida{overdue.length > 1 ? "s" : ""}
        </div>
      )}

      {view === "enfoque" && (
        <ul className="space-y-2 mt-4">
          {pending.map((t) => (
            <li
              key={t.id}
              draggable
              onDragStart={() => setDragId(t.id)}
              onDragOver={(e) => e.preventDefault()}
              onDrop={() => handleDropReorder(t.id)}
              className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-3 flex items-center gap-3 cursor-grab"
            >
              <button
                onClick={() => handleComplete(t.id, t.status === "done")}
                className="w-4 h-4 rounded-full border border-neutral-400 dark:border-neutral-600 flex-shrink-0"
                aria-label="Completar"
              />
              <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: categoryColor(t.category) }} />
              <div className="flex-1">
                <p className="text-sm text-neutral-800 dark:text-neutral-200">{t.title}</p>
                <p className="text-xs text-neutral-400 dark:text-neutral-500">
                  {t.category}
                  {t.due_date && ` · ${t.due_date}`}
                  {t.overdue && <span className="text-red-500 dark:text-red-400"> · Vencida</span>}
                </p>
              </div>
            </li>
          ))}
          {pending.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Sin tareas pendientes.</p>}
        </ul>
      )}

      {view === "eisenhower" && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
          {quadrants.map((q) => (
            <div
              key={q.key}
              data-testid={`quadrant-${q.key}`}
              onDragOver={(e) => e.preventDefault()}
              onDrop={() => handleDropQuadrant(q.key)}
              className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-3 min-h-[140px]"
            >
              <p className="text-xs font-medium mb-2 flex items-center gap-2 text-neutral-600 dark:text-neutral-300">
                <span className="w-2 h-2 rounded-full" style={{ background: q.dot }} />
                {q.label.toUpperCase()}
              </p>
              <ul className="space-y-2">
                {pending
                  .filter((t) => t.quadrant === q.key)
                  .map((t) => (
                    <li
                      key={t.id}
                      data-testid={`task-card-${t.id}`}
                      draggable
                      onDragStart={() => setDragId(t.id)}
                      className="bg-neutral-50 dark:bg-neutral-800 border border-neutral-200 dark:border-neutral-700 rounded-md p-2 text-sm text-neutral-800 dark:text-neutral-200 cursor-grab"
                    >
                      {t.title}
                    </li>
                  ))}
              </ul>
            </div>
          ))}
          <div className="md:col-span-2 border border-dashed border-neutral-300 dark:border-neutral-700 rounded-lg p-3">
            <p className="text-xs font-medium mb-2 text-neutral-500 dark:text-neutral-400">SIN CLASIFICAR</p>
            <ul className="flex flex-wrap gap-2">
              {pending
                .filter((t) => !t.quadrant)
                .map((t) => (
                  <li
                    key={t.id}
                    data-testid={`task-card-${t.id}`}
                    draggable
                    onDragStart={() => setDragId(t.id)}
                    className="bg-neutral-50 dark:bg-neutral-800 border border-neutral-200 dark:border-neutral-700 rounded-md px-2 py-1 text-sm text-neutral-800 dark:text-neutral-200 cursor-grab"
                  >
                    {t.title}
                  </li>
                ))}
            </ul>
          </div>
        </div>
      )}

      {view === "kanban" && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
          {columns.map((c) => (
            <div
              key={c.key}
              data-testid={`column-${c.key}`}
              onDragOver={(e) => e.preventDefault()}
              onDrop={() => handleDropStatus(c.key)}
              className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-3 min-h-[200px]"
            >
              <p className="text-xs font-medium mb-2 flex items-center gap-2 text-neutral-600 dark:text-neutral-300">
                <span className="w-2 h-2 rounded-full" style={{ background: c.dot }} />
                {c.label.toUpperCase()}
              </p>
              <ul className="space-y-2">
                {tasks
                  .filter((t) => t.status === c.key)
                  .map((t) => (
                    <li
                      key={t.id}
                      data-testid={`task-card-${t.id}`}
                      draggable
                      onDragStart={() => setDragId(t.id)}
                      className="bg-neutral-50 dark:bg-neutral-800 border border-neutral-200 dark:border-neutral-700 rounded-md p-2 text-sm text-neutral-800 dark:text-neutral-200 cursor-grab"
                    >
                      {t.title}
                    </li>
                  ))}
              </ul>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
