import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as notesApi from "../../api/notes";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Notes() {
  const [notes, setNotes] = useState<notesApi.Note[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [selected, setSelected] = useState<notesApi.Note | null>(null);
  const [locked, setLocked] = useState(false);

  async function load() {
    try {
      setNotes(await notesApi.listNotes());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreate() {
    if (!title.trim()) return;
    await notesApi.createNote(title.trim(), body);
    setTitle("");
    setBody("");
    load();
  }

  async function handleExport() {
    try {
      await notesApi.exportVault();
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-3xl grid grid-cols-1 md:grid-cols-2 gap-6">
      <div>
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Notas</h1>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Título"
          className="input-field w-full mb-2"
        />
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="Escribí tu nota en markdown. Usá [[Otra nota]] para enlazar."
          rows={5}
          className="input-field w-full mb-2 font-mono text-sm"
        />
        <button onClick={handleCreate} className="btn-primary mb-4">
          Guardar nota
        </button>

        <div className="flex gap-2 mb-4">
          <Link to="/app/notes/graph" className="btn-secondary">
            Ver grafo
          </Link>
          <button onClick={handleExport} className="btn-secondary">
            Exportar vault
          </button>
        </div>
        {locked && <UpgradeWall feature="Exportar el vault de Obsidian" />}
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Todas las notas</h2>
        <ul className="space-y-2">
          {notes.map((n) => (
            <li key={n.slug} className="surface-card-hover p-3">
              <button onClick={() => setSelected(n)} className="font-medium text-neutral-800 dark:text-neutral-200 text-sm hover:text-accent-strong dark:hover:text-accent-soft transition-colors duration-150">
                {n.title}
              </button>
              {n.links.length > 0 && (
                <p className="text-xs text-neutral-400 dark:text-neutral-500 mt-1">Enlaza a: {n.links.join(", ")}</p>
              )}
            </li>
          ))}
          {notes.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no creaste notas.</p>}
        </ul>
        {selected && (
          <div className="mt-4 surface-card p-3">
            <p className="font-medium text-neutral-800 dark:text-neutral-200">{selected.title}</p>
            <pre className="text-sm text-neutral-600 dark:text-neutral-400 whitespace-pre-wrap">{selected.body}</pre>
          </div>
        )}
      </div>
    </div>
  );
}
