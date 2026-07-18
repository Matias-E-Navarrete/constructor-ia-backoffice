import { useEffect, useState } from "react";
import * as studiesApi from "../../api/studies";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import { categoryColor } from "../../lib/categoryColor";

export default function Studies() {
  const [subjects, setSubjects] = useState<studiesApi.Subject[]>([]);
  const [sessions, setSessions] = useState<studiesApi.StudySession[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [newSubjectName, setNewSubjectName] = useState("");

  const [logFor, setLogFor] = useState<string | null>(null);
  const [duration, setDuration] = useState(30);
  const [topic, setTopic] = useState("");

  const [overview, setOverview] = useState<studiesApi.SubjectOverview[] | null>(null);
  const [locked, setLocked] = useState(false);

  async function load() {
    try {
      setSubjects(await studiesApi.listSubjects());
      setSessions(await studiesApi.listSessions());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreateSubject() {
    if (!newSubjectName.trim()) return;
    await studiesApi.createSubject(newSubjectName.trim());
    setNewSubjectName("");
    load();
  }

  async function handleDeleteSubject(id: string) {
    await studiesApi.deleteSubject(id);
    load();
  }

  async function handleLogSession(subjectId: string) {
    await studiesApi.logSession({
      subject_id: subjectId,
      session_date: new Date().toISOString().slice(0, 10),
      duration_minutes: duration,
      topic: topic || undefined,
    });
    setLogFor(null);
    setDuration(30);
    setTopic("");
    load();
  }

  async function handleShowOverview() {
    setLocked(false);
    setOverview(null);
    try {
      setOverview(await studiesApi.getOverview());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  const subjectName = (id: string) => subjects.find((s) => s.id === id)?.name ?? "—";

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Estudios</h1>

        <div className="flex gap-2 mb-6">
          <input
            value={newSubjectName}
            onChange={(e) => setNewSubjectName(e.target.value)}
            placeholder="Nueva materia"
            className="input-field flex-1"
          />
          <button onClick={handleCreateSubject} className="btn-primary">
            Agregar
          </button>
        </div>

        <ul className="space-y-2">
          {subjects.map((s) => (
            <li key={s.id} className="surface-card-hover p-4">
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-2 font-medium text-neutral-800 dark:text-neutral-200">
                  <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: categoryColor(s.name) }} />
                  {s.name}
                </span>
                <div className="flex gap-2">
                  <button
                    onClick={() => setLogFor(logFor === s.id ? null : s.id)}
                    className="btn-primary text-sm px-3 py-1"
                  >
                    Registrar sesión
                  </button>
                  <button onClick={() => handleDeleteSubject(s.id)} className="text-sm text-red-500 dark:text-red-400 px-2 hover:text-red-600 dark:hover:text-red-300 transition-colors duration-150">
                    Eliminar
                  </button>
                </div>
              </div>
              {logFor === s.id && (
                <div className="mt-3 flex gap-2 items-center">
                  <input
                    type="number"
                    value={duration}
                    onChange={(e) => setDuration(Number(e.target.value))}
                    aria-label="Minutos"
                    className="input-field w-24 px-2 py-1 text-sm"
                  />
                  <input
                    value={topic}
                    onChange={(e) => setTopic(e.target.value)}
                    placeholder="Tema (opcional)"
                    className="input-field flex-1 px-2 py-1 text-sm"
                  />
                  <button onClick={() => handleLogSession(s.id)} className="btn-accent text-sm px-3 py-1.5">
                    Guardar
                  </button>
                </div>
              )}
            </li>
          ))}
          {subjects.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no agregaste ninguna materia.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Sesiones recientes</h2>
        <ul className="space-y-2">
          {sessions.map((sess) => (
            <li key={sess.id} className="surface-card flex justify-between text-sm p-3">
              <span className="text-neutral-800 dark:text-neutral-200">
                {subjectName(sess.subject_id)}
                {sess.topic && <span className="text-neutral-400 dark:text-neutral-500"> · {sess.topic}</span>}
              </span>
              <span className="text-neutral-500 dark:text-neutral-400">
                {sess.session_date} · {sess.duration_minutes}min
              </span>
            </li>
          ))}
          {sessions.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no registraste sesiones de estudio.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Panorama de estudio</h2>
        <button onClick={handleShowOverview} className="btn-secondary mb-3">
          Ver panorama
        </button>
        {locked && <UpgradeWall feature="Panorama de estudio" />}
        {overview && (
          <table className="w-full text-sm text-left surface-card p-4">
            <thead className="text-xs text-neutral-400 dark:text-neutral-500">
              <tr>
                <th className="font-medium pb-2">Materia</th>
                <th className="font-medium pb-2">Sesiones</th>
                <th className="font-medium pb-2">Tiempo total</th>
                <th className="font-medium pb-2">Racha</th>
              </tr>
            </thead>
            <tbody className="text-neutral-800 dark:text-neutral-200">
              {overview.map((o) => (
                <tr key={o.subject_id} className="border-t border-neutral-100 dark:border-neutral-800">
                  <td className="py-2">{o.name}</td>
                  <td className="py-2">{o.session_count}</td>
                  <td className="py-2">{Math.floor(o.total_minutes / 60)}h {o.total_minutes % 60}min</td>
                  <td className="py-2">{o.current_streak} días</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
