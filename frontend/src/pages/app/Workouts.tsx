import { useEffect, useState } from "react";
import * as workoutsApi from "../../api/workouts";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import { categoryColor } from "../../lib/categoryColor";

function ProgressChart({ entries }: { entries: workoutsApi.Progress["entries"] }) {
  if (entries.length === 0) return null;
  const width = 480;
  const height = 140;
  const pad = 24;
  const weights = entries.map((e) => e.weight_kg);
  const min = Math.min(...weights);
  const max = Math.max(...weights);
  const span = max - min || 1;

  const points = entries.map((e, i) => {
    const x = pad + (i / Math.max(entries.length - 1, 1)) * (width - pad * 2);
    const y = height - pad - ((e.weight_kg - min) / span) * (height - pad * 2);
    return [x, y] as const;
  });

  return (
    <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`} className="max-w-full">
      <line x1={pad} y1={height - pad} x2={width - pad} y2={height - pad} stroke="currentColor" className="text-neutral-200 dark:text-neutral-800" strokeWidth={1} />
      <polyline points={points.map((p) => p.join(",")).join(" ")} fill="none" stroke="#22c55e" strokeWidth={2} />
      {points.map(([x, y], i) => (
        <circle key={i} cx={x} cy={y} r={3} fill="#22c55e" />
      ))}
    </svg>
  );
}

export default function Workouts() {
  const [sessions, setSessions] = useState<workoutsApi.Session[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [exercise, setExercise] = useState("");
  const [reps, setReps] = useState(10);
  const [weight, setWeight] = useState(20);
  const [progressExercise, setProgressExercise] = useState("");
  const [progress, setProgress] = useState<workoutsApi.Progress | null>(null);
  const [locked, setLocked] = useState(false);

  const [routines, setRoutines] = useState<workoutsApi.Routine[]>([]);
  const [routinesLocked, setRoutinesLocked] = useState(false);
  const [newRoutineName, setNewRoutineName] = useState("");
  const [newRoutineExercises, setNewRoutineExercises] = useState("");
  const [activeRoutine, setActiveRoutine] = useState<{ routine: workoutsApi.Routine; sets: workoutsApi.SetEntry[] } | null>(null);

  const [personalRecords, setPersonalRecords] = useState<workoutsApi.PersonalRecord[] | null>(null);
  const [prLocked, setPrLocked] = useState(false);

  async function load() {
    try {
      setSessions(await workoutsApi.listSessions());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
    try {
      setRoutines(await workoutsApi.listRoutines());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setRoutinesLocked(true);
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

  async function handleShowPersonalRecords() {
    setPrLocked(false);
    setPersonalRecords(null);
    try {
      setPersonalRecords(await workoutsApi.getPersonalRecords());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setPrLocked(true);
    }
  }

  async function handleCreateRoutine() {
    const exercises = newRoutineExercises
      .split(",")
      .map((e) => e.trim())
      .filter(Boolean);
    if (!newRoutineName.trim() || exercises.length === 0) return;
    try {
      await workoutsApi.createRoutine(newRoutineName.trim(), exercises);
      setNewRoutineName("");
      setNewRoutineExercises("");
      setRoutinesLocked(false);
      load();
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setRoutinesLocked(true);
    }
  }

  function handleStartRoutine(routine: workoutsApi.Routine) {
    setActiveRoutine({
      routine,
      sets: routine.exercises.map((name, i) => ({ exercise_name: name, set_number: i + 1, reps: 10, weight_kg: 20 })),
    });
  }

  function updateActiveSet(i: number, field: "reps" | "weight_kg", value: number) {
    if (!activeRoutine) return;
    const sets = activeRoutine.sets.map((s, idx) => (idx === i ? { ...s, [field]: value } : s));
    setActiveRoutine({ ...activeRoutine, sets });
  }

  async function handleFinishRoutine() {
    if (!activeRoutine) return;
    await workoutsApi.logSession({
      session_date: new Date().toISOString().slice(0, 10),
      notes: `Rutina: ${activeRoutine.routine.name}`,
      sets: activeRoutine.sets,
    });
    setActiveRoutine(null);
    load();
  }

  async function handleDeleteRoutine(id: string) {
    await workoutsApi.deleteRoutine(id);
    load();
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Entrenamientos</h1>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-2 mb-6">
          <input
            value={exercise}
            onChange={(e) => setExercise(e.target.value)}
            placeholder="Ejercicio"
            className="input-field col-span-1"
          />
          <input
            type="number"
            value={reps}
            onChange={(e) => setReps(Number(e.target.value))}
            placeholder="Reps"
            className="input-field"
          />
          <input
            type="number"
            value={weight}
            onChange={(e) => setWeight(Number(e.target.value))}
            placeholder="Kg"
            className="input-field"
          />
        </div>
        <button onClick={handleLog} className="btn-primary mb-6">
          Registrar sesión
        </button>

        <ul className="space-y-2">
          {sessions.map((s) => (
            <li key={s.id} className="surface-card-hover p-4">
              <p className="font-medium text-neutral-800 dark:text-neutral-200">{s.session_date}</p>
              {s.notes && <p className="text-xs text-neutral-400 dark:text-neutral-500">{s.notes}</p>}
              <ul className="text-sm text-neutral-600 dark:text-neutral-400">
                {s.sets.map((set, i) => (
                  <li key={i}>
                    {set.exercise_name}: {set.reps} reps @ {set.weight_kg}kg
                  </li>
                ))}
              </ul>
            </li>
          ))}
          {sessions.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no registraste entrenamientos.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Rutinas</h2>
        {routinesLocked && <UpgradeWall feature="Rutinas de entrenamiento" />}

        <div className="flex flex-wrap gap-2 mb-3">
          <input
            value={newRoutineName}
            onChange={(e) => setNewRoutineName(e.target.value)}
            placeholder="Nombre de la rutina"
            className="input-field text-sm"
          />
          <input
            value={newRoutineExercises}
            onChange={(e) => setNewRoutineExercises(e.target.value)}
            placeholder="Ejercicios separados por coma"
            className="input-field flex-1 text-sm"
          />
          <button onClick={handleCreateRoutine} className="btn-primary text-sm px-3 py-2">
            Crear rutina →
          </button>
        </div>

        <ul className="space-y-2">
          {routines.map((r) => (
            <li key={r.id} className="surface-card-hover p-3 flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-neutral-800 dark:text-neutral-200">{r.name}</p>
                <p className="text-xs text-neutral-400 dark:text-neutral-500">{r.exercises.join(", ")}</p>
              </div>
              <div className="flex gap-2">
                <button onClick={() => handleStartRoutine(r)} className="btn-accent text-xs px-3 py-1.5">
                  Empezar
                </button>
                <button onClick={() => handleDeleteRoutine(r.id)} className="text-xs text-red-500 dark:text-red-400 px-2 hover:text-red-600 dark:hover:text-red-300 transition-colors duration-150">
                  Eliminar
                </button>
              </div>
            </li>
          ))}
          {routines.length === 0 && !routinesLocked && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no creaste ninguna rutina.</p>}
        </ul>

        {activeRoutine && (
          <div className="mt-4 surface-card p-4">
            <p className="text-sm font-medium text-neutral-800 dark:text-neutral-200 mb-3">{activeRoutine.routine.name}</p>
            <div className="space-y-2 mb-3">
              {activeRoutine.sets.map((set, i) => (
                <div key={i} className="grid grid-cols-3 gap-2 items-center text-sm">
                  <span className="text-neutral-700 dark:text-neutral-300">{set.exercise_name}</span>
                  <input
                    type="number"
                    value={set.reps}
                    onChange={(e) => updateActiveSet(i, "reps", Number(e.target.value))}
                    aria-label={`Reps ${set.exercise_name}`}
                    className="input-field px-2 py-1"
                  />
                  <input
                    type="number"
                    value={set.weight_kg}
                    onChange={(e) => updateActiveSet(i, "weight_kg", Number(e.target.value))}
                    aria-label={`Kg ${set.exercise_name}`}
                    className="input-field px-2 py-1"
                  />
                </div>
              ))}
            </div>
            <button onClick={handleFinishRoutine} className="btn-primary">
              Terminar rutina
            </button>
          </div>
        )}
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Récords personales</h2>
        <button onClick={handleShowPersonalRecords} className="btn-secondary mb-3">
          Ver récords
        </button>
        {prLocked && <UpgradeWall feature="Récords personales" />}
        {personalRecords && (
          <div className="flex flex-wrap gap-2">
            {personalRecords.map((pr) => (
              <div
                key={pr.exercise_name}
                data-testid={`pr-${pr.exercise_name}`}
                className="surface-card-hover flex items-center gap-2 px-3 py-2"
              >
                <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: categoryColor(pr.exercise_name) }} />
                <span className="text-sm text-neutral-800 dark:text-neutral-200">{pr.exercise_name}</span>
                <span className="text-sm font-semibold text-accent-strong dark:text-accent-soft">{pr.personal_record_kg}kg</span>
              </div>
            ))}
            {personalRecords.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no hay récords registrados.</p>}
          </div>
        )}
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Progreso</h2>
        <div className="flex gap-2 mb-4">
          <input
            value={progressExercise}
            onChange={(e) => setProgressExercise(e.target.value)}
            placeholder="Nombre del ejercicio"
            className="input-field flex-1"
          />
          <button onClick={handleProgress} className="btn-secondary">
            Ver progreso
          </button>
        </div>
        {locked && <UpgradeWall feature="Progreso de entrenamientos" />}
        {progress && (
          <div className="text-sm text-neutral-600 dark:text-neutral-400">
            <p className="mb-2">Récord personal: {progress.personal_record_kg}kg</p>
            <ProgressChart entries={progress.entries} />
            <ul className="mt-2">
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
