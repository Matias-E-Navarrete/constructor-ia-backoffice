import { useEffect, useState } from "react";
import * as tasksApi from "../../api/tasks";
import * as authApi from "../../api/auth";
import { useAuth } from "../../context/AuthContext";
import { FeatureDisabledError } from "../../api/client";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import FocusMode from "../../components/FocusMode";
import { categoryColor } from "../../lib/categoryColor";

function todayISO() {
  return new Date().toISOString().slice(0, 10);
}

export default function Planner() {
  const { user, refreshUser } = useAuth();
  const [tasks, setTasks] = useState<tasksApi.Task[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [dragId, setDragId] = useState<string | null>(null);
  const [showFocus, setShowFocus] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [startHour, setStartHour] = useState(user?.planner_start_hour ?? 6);
  const [endHour, setEndHour] = useState(user?.planner_end_hour ?? 23);

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

  useEffect(() => {
    setStartHour(user?.planner_start_hour ?? 6);
    setEndHour(user?.planner_end_hour ?? 23);
  }, [user]);

  async function handleDrop(hour: number) {
    if (!dragId) return;
    const at = new Date();
    at.setHours(hour, 0, 0, 0);
    await tasksApi.scheduleTask(dragId, at.toISOString());
    setDragId(null);
    load();
  }

  async function handleSaveSettings() {
    await authApi.updatePlannerHours(startHour, endHour);
    await refreshUser();
    setShowSettings(false);
  }

  if (disabled) return <FeatureDisabledBanner />;

  const unscheduled = tasks.filter((t) => t.status !== "done" && !t.scheduled_at);
  const hours: number[] = [];
  for (let h = startHour; h <= endHour; h++) hours.push(h);

  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Planificador</h1>
        <div className="flex gap-2">
          <button onClick={() => setShowSettings(true)} className="btn-secondary text-xs px-3 py-1.5">
            ⚙ Horas
          </button>
          <button onClick={() => setShowFocus(true)} className="btn-accent text-xs px-3 py-1.5">
            ▶ Modo Foco
          </button>
        </div>
      </div>
      <p className="text-sm text-neutral-500 dark:text-neutral-400 mb-4">Arrastra tareas para organizar tu día.</p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="space-y-2">
          {unscheduled.map((t) => (
            <div
              key={t.id}
              draggable
              onDragStart={() => setDragId(t.id)}
              className="surface-card-hover p-3 cursor-grab flex items-center gap-2"
            >
              <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: categoryColor(t.category) }} />
              <span className="text-sm text-neutral-800 dark:text-neutral-200">{t.title}</span>
            </div>
          ))}
          {unscheduled.length === 0 && <p className="text-neutral-400 dark:text-neutral-500 text-sm">Todo agendado.</p>}
        </div>

        <div className="md:col-span-2 surface-card divide-y divide-neutral-100 dark:divide-neutral-800">
          {hours.map((h) => {
            const scheduled = tasks.filter((t) => {
              if (!t.scheduled_at) return false;
              const d = new Date(t.scheduled_at);
              return d.getHours() === h && d.toISOString().slice(0, 10) === todayISO();
            });
            return (
              <div
                key={h}
                data-testid={`hour-${h}`}
                onDragOver={(e) => e.preventDefault()}
                onDrop={() => handleDrop(h)}
                className="flex gap-3 p-2 min-h-[48px]"
              >
                <span className="text-xs text-neutral-400 dark:text-neutral-500 w-16 flex-shrink-0">
                  {h % 12 === 0 ? 12 : h % 12} {h < 12 ? "AM" : "PM"}
                </span>
                <div className="flex-1 space-y-1">
                  {scheduled.map((t) => (
                    <div
                      key={t.id}
                      draggable
                      onDragStart={() => setDragId(t.id)}
                      className="text-sm bg-neutral-50 dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-md px-2 py-1 text-neutral-800 dark:text-neutral-200 hover:border-accent/40 transition-colors duration-150"
                    >
                      {t.title}
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {showSettings && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-40" onClick={() => setShowSettings(false)}>
          <div
            onClick={(e) => e.stopPropagation()}
            className="surface-card shadow-xl p-6 w-full max-w-sm"
          >
            <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-1">Horas del planner</h2>
            <p className="text-sm text-neutral-500 dark:text-neutral-400 mb-4">Define el rango visible de la timeline.</p>
            <div className="grid grid-cols-2 gap-2 mb-4">
              <label className="text-xs text-neutral-500 dark:text-neutral-400">
                Empieza a las
                <input
                  type="number"
                  min={0}
                  max={23}
                  value={startHour}
                  onChange={(e) => setStartHour(Number(e.target.value))}
                  className="input-field w-full mt-1 px-2 py-1"
                />
              </label>
              <label className="text-xs text-neutral-500 dark:text-neutral-400">
                Termina a las
                <input
                  type="number"
                  min={1}
                  max={24}
                  value={endHour}
                  onChange={(e) => setEndHour(Number(e.target.value))}
                  className="input-field w-full mt-1 px-2 py-1"
                />
              </label>
            </div>
            <div className="flex justify-end gap-2">
              <button onClick={() => setShowSettings(false)} className="btn-ghost text-sm px-3 py-1.5">
                Cancelar
              </button>
              <button onClick={handleSaveSettings} className="btn-primary text-sm px-3 py-1.5">
                Guardar
              </button>
            </div>
          </div>
        </div>
      )}

      {showFocus && <FocusMode onClose={() => setShowFocus(false)} />}
    </div>
  );
}
