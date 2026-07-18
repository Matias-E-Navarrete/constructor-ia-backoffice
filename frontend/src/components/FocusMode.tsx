import { useEffect, useRef, useState } from "react";
import * as tasksApi from "../api/tasks";
import * as inboxApi from "../api/inbox";

const DEFAULT_SECONDS = 15 * 60;

export default function FocusMode({ onClose }: { onClose: () => void }) {
  const [tasks, setTasks] = useState<tasksApi.Task[]>([]);
  const [secondsLeft, setSecondsLeft] = useState(DEFAULT_SECONDS);
  const [running, setRunning] = useState(false);
  const [note, setNote] = useState("");
  const [lastCompleted, setLastCompleted] = useState<string | null>(null);
  const intervalRef = useRef<number | null>(null);

  useEffect(() => {
    tasksApi.listTasks({ status: "todo" }).then(setTasks);
  }, []);

  useEffect(() => {
    if (running) {
      intervalRef.current = window.setInterval(() => {
        setSecondsLeft((s) => (s > 0 ? s - 1 : 0));
      }, 1000);
    } else if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [running]);

  function reset() {
    setRunning(false);
    setSecondsLeft(DEFAULT_SECONDS);
  }

  async function handleComplete(task: tasksApi.Task) {
    await tasksApi.completeTask(task.id);
    setTasks((ts) => ts.filter((t) => t.id !== task.id));
    setLastCompleted(task.title);
  }

  async function handleSaveNote() {
    if (!note.trim()) return;
    await inboxApi.captureItem(note.trim());
    setNote("");
  }

  const minutes = String(Math.floor(secondsLeft / 60)).padStart(2, "0");
  const seconds = String(secondsLeft % 60).padStart(2, "0");
  const progress = 1 - secondsLeft / DEFAULT_SECONDS;

  return (
    <div data-testid="focus-mode" className="fixed inset-0 bg-white dark:bg-black z-50 flex flex-col md:flex-row">
      <button onClick={onClose} className="absolute top-4 left-4 text-neutral-500 dark:text-neutral-400 text-xl" aria-label="Cerrar">
        ×
      </button>

      {lastCompleted && (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-neutral-900 dark:bg-white text-white dark:text-neutral-900 text-sm px-4 py-2 rounded-md">
          Tarea completada: {lastCompleted}
        </div>
      )}

      <div className="flex-1 flex flex-col items-center justify-center">
        <span className="text-xs uppercase tracking-widest bg-neutral-100 dark:bg-neutral-900 text-neutral-500 dark:text-neutral-400 px-3 py-1 rounded-full mb-6">
          Foco
        </span>
        <div className="text-7xl font-light text-neutral-900 dark:text-neutral-50 tabular-nums">
          {minutes}:{seconds}
        </div>
        <div className="flex gap-3 mt-6">
          <button
            onClick={() => setRunning((r) => !r)}
            className="w-10 h-10 rounded-md bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 flex items-center justify-center"
          >
            {running ? "❚❚" : "▶"}
          </button>
          <button onClick={reset} className="w-10 h-10 rounded-md border border-neutral-300 dark:border-neutral-700 text-neutral-600 dark:text-neutral-300">
            ↺
          </button>
        </div>
        <div className="w-64 h-0.5 bg-neutral-200 dark:bg-neutral-800 mt-8">
          <div className="h-0.5 bg-accent" style={{ width: `${progress * 100}%` }} />
        </div>
      </div>

      <div className="md:w-96 border-t md:border-t-0 md:border-l border-neutral-200 dark:border-neutral-800 p-6 flex flex-col">
        <h2 className="text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">Flow de hoy</h2>
        <p className="text-xs text-neutral-400 dark:text-neutral-500 mb-4">{tasks.length} pendientes</p>
        <ul className="space-y-2 flex-1 overflow-auto">
          {tasks.map((t) => (
            <li key={t.id} className="flex items-center gap-2 bg-neutral-50 dark:bg-neutral-900 rounded-md p-2">
              <button
                onClick={() => handleComplete(t)}
                className="w-4 h-4 rounded-full border border-neutral-400 dark:border-neutral-600 flex-shrink-0"
                aria-label="Completar"
              />
              <span className="text-sm text-neutral-800 dark:text-neutral-200">{t.title}</span>
            </li>
          ))}
        </ul>
        <div className="flex gap-2 mt-4">
          <input
            value={note}
            onChange={(e) => setNote(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSaveNote()}
            placeholder="Nota rápida..."
            className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-2 py-1 text-sm"
          />
          <button onClick={handleSaveNote} className="text-sm px-2 rounded-md bg-neutral-100 dark:bg-neutral-800 text-neutral-600 dark:text-neutral-300">
            ➤
          </button>
        </div>
        <button onClick={onClose} className="mt-3 text-sm text-neutral-500 dark:text-neutral-400 border-t border-neutral-200 dark:border-neutral-800 pt-3">
          Terminar Flow
        </button>
      </div>
    </div>
  );
}
