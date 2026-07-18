import { useEffect, useState } from "react";
import * as workoutsApi from "../../api/workouts";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Workouts() {
  const [sessions, setSessions] = useState<workoutsApi.Session[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [exercise, setExercise] = useState("");
  const [reps, setReps] = useState(10);
  const [weight, setWeight] = useState(20);
  const [progressExercise, setProgressExercise] = useState("");
  const [progress, setProgress] = useState<workoutsApi.Progress | null>(null);
  const [locked, setLocked] = useState(false);

  async function load() {
    try {
      setSessions(await workoutsApi.listSessions());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleLog() {
    if (!exercise.trim()) return;
    await workoutsApi.logSession({
      session_date: new Date().toISOString().slice(0, 10),
      sets: [{ exercise_name: exercise.trim(), set_number: 1, reps, weight_kg: weight }],
    });
    setExercise("");
    load();
  }

  async function handleProgress() {
    if (!progressExercise.trim()) return;
    setLocked(false);
    setProgress(null);
    try {
      setProgress(await workoutsApi.getProgress(progressExercise.trim()));
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-50 mb-4">Entrenamientos</h1>
        <div className="grid grid-cols-3 gap-2 mb-6">
          <input
            value={exercise}
            onChange={(e) => setExercise(e.target.value)}
            placeholder="Ejercicio"
            className="rounded-md border border-slate-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-slate-900 dark:text-slate-100 px-3 py-2 col-span-1"
          />
          <input
            type="number"
            value={reps}
            onChange={(e) => setReps(Number(e.target.value))}
            placeholder="Reps"
            className="rounded-md border border-slate-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-slate-900 dark:text-slate-100 px-3 py-2"
          />
          <input
            type="number"
            value={weight}
            onChange={(e) => setWeight(Number(e.target.value))}
            placeholder="Kg"
            className="rounded-md border border-slate-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-slate-900 dark:text-slate-100 px-3 py-2"
          />
        </div>
        <button onClick={handleLog} className="bg-slate-900 text-white dark:bg-white dark:text-slate-900 px-4 py-2 rounded-md text-sm font-medium mb-6">
          Registrar sesión
        </button>

        <ul className="space-y-2">
          {sessions.map((s) => (
            <li key={s.id} className="bg-white dark:bg-neutral-900 border border-slate-200 dark:border-neutral-800 rounded-lg p-4">
              <p className="font-medium text-slate-800 dark:text-slate-200">{s.session_date}</p>
              <ul className="text-sm text-slate-600 dark:text-slate-400">
                {s.sets.map((set, i) => (
                  <li key={i}>
                    {set.exercise_name}: {set.reps} reps @ {set.weight_kg}kg
                  </li>
                ))}
              </ul>
            </li>
          ))}
          {sessions.length === 0 && <p className="text-slate-500 dark:text-slate-400 text-sm">Todavía no registraste entrenamientos.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-50 mb-2">Progreso</h2>
        <div className="flex gap-2 mb-4">
          <input
            value={progressExercise}
            onChange={(e) => setProgressExercise(e.target.value)}
            placeholder="Nombre del ejercicio"
            className="flex-1 rounded-md border border-slate-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-slate-900 dark:text-slate-100 px-3 py-2"
          />
          <button onClick={handleProgress} className="bg-slate-100 dark:bg-neutral-800 text-slate-700 dark:text-slate-300 px-4 py-2 rounded-md text-sm font-medium">
            Ver progreso
          </button>
        </div>
        {locked && <UpgradeWall feature="Progreso de entrenamientos" />}
        {progress && (
          <div className="text-sm text-slate-600 dark:text-slate-400">
            <p className="mb-2">Récord personal: {progress.personal_record_kg}kg</p>
            <ul>
              {progress.entries.map((e, i) => (
                <li key={i}>
                  {e.session_date}: {e.reps} reps @ {e.weight_kg}kg
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
