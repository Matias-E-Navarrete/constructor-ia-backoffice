import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as habitsApi from "../../api/habits";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Habits() {
  const [habits, setHabits] = useState<habitsApi.Habit[]>([]);
  const [name, setName] = useState("");
  const [noteSlug, setNoteSlug] = useState("");
  const [disabled, setDisabled] = useState(false);
  const [statsFor, setStatsFor] = useState<string | null>(null);
  const [stats, setStats] = useState<habitsApi.HabitStats | null>(null);
  const [locked, setLocked] = useState(false);

  async function load() {
    try {
      setHabits(await habitsApi.listHabits());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreate() {
    if (!name.trim()) return;
    await habitsApi.createHabit(name.trim(), noteSlug.trim() || undefined);
    setName("");
    setNoteSlug("");
    load();
  }

  async function handleCheckIn(id: string) {
    await habitsApi.checkInHabit(id);
    load();
  }

  async function handleShowStats(id: string) {
    setStatsFor(id);
    setLocked(false);
    setStats(null);
    try {
      setStats(await habitsApi.getHabitStats(id));
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-2xl">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Hábitos</h1>
        <Link
          to="/app/habits/panorama"
          className="text-sm text-accent-strong dark:text-accent-soft font-medium hover:opacity-80 transition-opacity duration-150"
        >
          Ver Panorama →
        </Link>
      </div>

      <div className="flex flex-col sm:flex-row gap-2 mb-6">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Nuevo hábito"
          className="input-field flex-1 min-w-0"
        />
        <input
          value={noteSlug}
          onChange={(e) => setNoteSlug(e.target.value)}
          placeholder="Vincular a nota (slug, opcional)"
          className="input-field flex-1 min-w-0"
        />
        <button onClick={handleCreate} className="btn-primary">
          Agregar
        </button>
      </div>

      <ul className="space-y-2">
        {habits.map((h) => (
          <li key={h.id} className="surface-card-hover p-4">
            <div className="flex items-center justify-between">
              <span className="font-medium text-neutral-800 dark:text-neutral-200">{h.name}</span>
              <div className="flex gap-2">
                <button onClick={() => handleCheckIn(h.id)} className="btn-primary text-sm px-3 py-1">
                  Check-in
                </button>
                <button onClick={() => handleShowStats(h.id)} className="btn-secondary text-sm px-3 py-1">
                  Estadísticas
                </button>
              </div>
            </div>
            {statsFor === h.id && (
              <div className="mt-3">
                {locked && <UpgradeWall feature="Estadísticas de hábitos" />}
                {stats && (
                  <div className="text-sm text-neutral-600 dark:text-neutral-400 grid grid-cols-3 gap-2">
                    <span>Racha actual: {stats.current_streak}</span>
                    <span>Completados: {stats.completed_logs}</span>
                    <span>Tasa: {Math.round(stats.completion_rate * 100)}%</span>
                  </div>
                )}
              </div>
            )}
          </li>
        ))}
        {habits.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no agregaste ningún hábito.</p>}
      </ul>
    </div>
  );
}
