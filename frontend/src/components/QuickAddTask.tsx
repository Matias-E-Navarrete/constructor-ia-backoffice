import { useState } from "react";
import * as tasksApi from "../api/tasks";

const CATEGORY_PRESETS = ["General", "Trabajo", "Personal", "Estudios", "Hábitos", "Gimnasio"];

export default function QuickAddTask({
  onCreated,
  initialTitle,
}: {
  onCreated: () => void;
  initialTitle?: string;
}) {
  const [expanded, setExpanded] = useState(!!initialTitle);
  const [title, setTitle] = useState(initialTitle ?? "");
  const [category, setCategory] = useState("General");
  const [subcategory, setSubcategory] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [dueTime, setDueTime] = useState("");
  const [repeatRule, setRepeatRule] = useState<tasksApi.RepeatRule>("none");
  const [priority, setPriority] = useState<tasksApi.Priority>("normal");
  const [subtasks, setSubtasks] = useState<tasksApi.Subtask[]>([]);
  const [subtaskDraft, setSubtaskDraft] = useState("");

  function addSubtask() {
    if (!subtaskDraft.trim()) return;
    setSubtasks((s) => [...s, { title: subtaskDraft.trim(), done: false }]);
    setSubtaskDraft("");
  }

  async function handleSubmit() {
    if (!title.trim()) return;
    await tasksApi.createTask({
      title: title.trim(),
      category,
      subcategory: subcategory || undefined,
      due_date: dueDate || undefined,
      due_time: dueTime || undefined,
      repeat_rule: repeatRule,
      priority,
      subtasks,
    });
    setTitle("");
    setSubcategory("");
    setDueDate("");
    setDueTime("");
    setRepeatRule("none");
    setPriority("normal");
    setSubtasks([]);
    setExpanded(false);
    onCreated();
  }

  return (
    <div className="surface-card p-3 mb-4 focus-within:border-accent/40 focus-within:shadow-glow-sm">
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onFocus={() => setExpanded(true)}
        placeholder="ej: Reunión con el equipo"
        className="w-full bg-transparent text-neutral-900 dark:text-neutral-100 outline-none placeholder:text-neutral-400"
      />
      {expanded && (
        <div className="mt-3 space-y-2">
          <div className="flex flex-wrap gap-2 text-xs">
            <select
              aria-label="Categoría"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            >
              {CATEGORY_PRESETS.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
            <input
              aria-label="Subcategoría"
              value={subcategory}
              onChange={(e) => setSubcategory(e.target.value)}
              placeholder="Subcategoría"
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 w-28 transition-colors duration-150 hover:border-accent/40"
            />
            <input
              aria-label="Fecha"
              type="date"
              value={dueDate}
              onChange={(e) => setDueDate(e.target.value)}
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            />
            <input
              aria-label="Hora"
              type="time"
              value={dueTime}
              onChange={(e) => setDueTime(e.target.value)}
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            />
            <select
              aria-label="Repetición"
              value={repeatRule}
              onChange={(e) => setRepeatRule(e.target.value as tasksApi.RepeatRule)}
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            >
              <option value="none">No repetir</option>
              <option value="daily">Diario</option>
              <option value="weekly">Semanal</option>
              <option value="monthly">Mensual</option>
            </select>
            <select
              aria-label="Prioridad"
              value={priority}
              onChange={(e) => setPriority(e.target.value as tasksApi.Priority)}
              className="rounded-full border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            >
              <option value="low">Baja</option>
              <option value="normal">Normal</option>
              <option value="high">Alta</option>
              <option value="urgent">Urgente</option>
            </select>
          </div>

          <div className="flex gap-2">
            <input
              value={subtaskDraft}
              onChange={(e) => setSubtaskDraft(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && addSubtask()}
              placeholder="Agregar subtarea..."
              className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-neutral-50 dark:bg-neutral-800 px-2 py-1 text-sm text-neutral-700 dark:text-neutral-300 transition-colors duration-150 hover:border-accent/40"
            />
            <button onClick={addSubtask} className="btn-secondary text-xs px-2 py-1">
              + Subtarea
            </button>
          </div>
          {subtasks.length > 0 && (
            <ul className="text-xs text-neutral-500 dark:text-neutral-400 list-disc list-inside">
              {subtasks.map((s, i) => (
                <li key={i}>{s.title}</li>
              ))}
            </ul>
          )}

          <div className="flex justify-between items-center pt-1">
            <span className="text-xs text-neutral-400 dark:text-neutral-500">Presiona Enter para crear</span>
            <button onClick={handleSubmit} className="btn-primary py-1.5">
              Añadir Tarea →
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
