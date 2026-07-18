import { useEffect, useState } from "react";
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
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Hábitos</h1>

      <div className="flex gap-2 mb-6">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Nuevo hábito"
          className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <input
          value={noteSlug}
          onChange={(e) => setNoteSlug(e.target.value)}
          placeholder="Vincular a nota (slug, opcional)"
          className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <button onClick={handleCreate} className="bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-4 py-2 rounded-md text-sm font-medium">
          Agregar
        </button>
      </div>

      <ul className="space-y-2">
        {habits.map((h) => (
          <li key={h.id} className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-4">
            <div className="flex items-center justify-between">
              <span className="font-medium text-neutral-800 dark:text-neutral-200">{h.name}</span>
              <div className="flex gap-2">
                <button onClick={() => handleCheckIn(h.id)} className="text-sm bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-3 py-1 rounded-md">
                  Check-in
                </button>
                <button onClick={() => handleShowStats(h.id)} className="text-sm bg-neutral-100 dark:bg-neutral-800 text-neutral-700 dark:text-neutral-300 px-3 py-1 rounded-md">
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
