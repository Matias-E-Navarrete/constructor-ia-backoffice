import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import * as inboxApi from "../../api/inbox";
import { FeatureDisabledError } from "../../api/client";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Inbox() {
  const [items, setItems] = useState<inboxApi.InboxItem[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [content, setContent] = useState("");
  const [search, setSearch] = useState("");
  const navigate = useNavigate();

  async function load() {
    try {
      setItems(await inboxApi.listInbox());
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCapture() {
    if (!content.trim()) return;
    await inboxApi.captureItem(content.trim());
    setContent("");
    load();
  }

  async function handleTogglePin(item: inboxApi.InboxItem) {
    await inboxApi.setPinned(item.id, !item.pinned);
    load();
  }

  async function handleDelete(id: string) {
    await inboxApi.deleteItem(id);
    load();
  }

  function handleOrganize(item: inboxApi.InboxItem) {
    navigate("/app/tasks", { state: { prefillTitle: item.content, fromInboxId: item.id } });
  }

  if (disabled) return <FeatureDisabledBanner />;

  const filtered = items.filter((i) => i.content.toLowerCase().includes(search.toLowerCase()));

  return (
    <div className="max-w-2xl">
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-1">Bandeja de Entrada</h1>
      <p className="text-sm text-neutral-500 dark:text-neutral-400 mb-4">Volcado de cerebro</p>

      <div className="flex gap-2 mb-4">
        <input
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleCapture()}
          placeholder="Escribí lo que se te ocurra..."
          className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <button onClick={handleCapture} className="bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-4 py-2 rounded-md text-sm font-medium">
          Guardar
        </button>
      </div>

      <input
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Buscar por contenido..."
        className="w-full mb-4 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2 text-sm"
      />

      <ul className="space-y-2">
        {filtered.map((item) => (
          <li
            key={item.id}
            className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-3 flex items-center justify-between text-sm text-neutral-800 dark:text-neutral-200"
          >
            <span>
              {item.pinned && <span className="mr-1">📌</span>}
              {item.content}
            </span>
            <div className="flex gap-1 text-xs">
              <button onClick={() => handleTogglePin(item)} className="px-2 py-1 rounded bg-neutral-100 dark:bg-neutral-800 text-neutral-600 dark:text-neutral-300">
                {item.pinned ? "Despinnear" : "Pin"}
              </button>
              <button onClick={() => handleOrganize(item)} className="px-2 py-1 rounded bg-neutral-100 dark:bg-neutral-800 text-neutral-600 dark:text-neutral-300">
                Organizar
              </button>
              <button onClick={() => handleDelete(item.id)} className="px-2 py-1 rounded bg-neutral-100 dark:bg-neutral-800 text-red-500 dark:text-red-400">
                Eliminar
              </button>
            </div>
          </li>
        ))}
        {filtered.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Tu bandeja está vacía.</p>}
      </ul>
    </div>
  );
}
