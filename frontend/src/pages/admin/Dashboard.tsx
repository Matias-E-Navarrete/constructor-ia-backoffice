import { useEffect, useState } from "react";
import * as adminApi from "../../api/admin";

const columns: adminApi.RoadmapItem["status"][] = ["planned", "in_progress", "done"];
const columnLabel: Record<string, string> = { planned: "Planeado", in_progress: "En progreso", done: "Hecho" };

export default function AdminDashboard() {
  const [items, setItems] = useState<adminApi.RoadmapItem[]>([]);
  const [title, setTitle] = useState("");
  const [kind, setKind] = useState<adminApi.RoadmapItem["kind"]>("task");

  async function load() {
    setItems(await adminApi.listRoadmap());
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreate() {
    if (!title.trim()) return;
    await adminApi.createRoadmapItem(title.trim(), "", kind);
    setTitle("");
    load();
  }

  async function handleMove(id: string, status: adminApi.RoadmapItem["status"]) {
    await adminApi.updateRoadmapStatus(id, status);
    load();
  }

  return (
    <div>
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Roadmap</h1>

      <div className="flex gap-2 mb-6">
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Nueva tarea/feature/bug"
          className="input-field flex-1"
        />
        <select value={kind} onChange={(e) => setKind(e.target.value as adminApi.RoadmapItem["kind"])} className="input-field">
          <option value="task">Tarea</option>
          <option value="feature">Feature</option>
          <option value="bug">Bug</option>
        </select>
        <button onClick={handleCreate} className="btn-primary">
          Agregar
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {columns.map((status) => (
          <div key={status} className="surface-card p-3">
            <h2 className="font-medium text-neutral-700 dark:text-neutral-300 text-sm mb-2">{columnLabel[status]}</h2>
            <ul className="space-y-2">
              {items
                .filter((it) => it.status === status)
                .map((it) => (
                  <li key={it.id} className="bg-neutral-50 dark:bg-neutral-800 border border-neutral-200 dark:border-neutral-700 rounded-lg p-2">
                    <p className="text-sm text-neutral-800 dark:text-neutral-200">{it.title}</p>
                    <p className="text-xs text-neutral-400 dark:text-neutral-500 mb-1">{it.kind}</p>
                    <div className="flex gap-1">
                      {columns
                        .filter((s) => s !== status)
                        .map((s) => (
                          <button
                            key={s}
                            onClick={() => handleMove(it.id, s)}
                            className="text-xs bg-white dark:bg-neutral-900 text-neutral-700 dark:text-neutral-300 border border-neutral-200 dark:border-neutral-700 px-2 py-0.5 rounded-md hover:border-accent/40 transition-colors duration-150"
                          >
                            → {columnLabel[s]}
                          </button>
                        ))}
                    </div>
                  </li>
                ))}
            </ul>
          </div>
        ))}
      </div>
    </div>
  );
}
