import { useEffect, useState } from "react";
import * as adminApi from "../../api/admin";

export default function AdminFeatures() {
  const [flags, setFlags] = useState<adminApi.Flag[]>([]);

  async function load() {
    setFlags(await adminApi.listFeatureFlags());
  }

  useEffect(() => {
    load();
  }, []);

  async function handleToggle(key: string, enabled: boolean) {
    await adminApi.toggleFeatureFlag(key, enabled);
    load();
  }

  return (
    <div className="max-w-xl">
      <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-50 mb-4">Feature flags</h1>
      <ul className="space-y-2">
        {flags.map((f) => (
          <li key={f.key} className="bg-white dark:bg-neutral-900 border border-slate-200 dark:border-neutral-800 rounded-lg p-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-slate-800 dark:text-slate-200">{f.key}</p>
              <p className="text-xs text-slate-500 dark:text-slate-400">{f.description}</p>
            </div>
            <button
              onClick={() => handleToggle(f.key, !f.enabled)}
              className={`text-xs font-medium px-3 py-1 rounded-full ${
                f.enabled
                  ? "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400"
                  : "bg-red-100 text-red-700 dark:bg-red-950 dark:text-red-400"
              }`}
            >
              {f.enabled ? "Activo" : "Apagado"}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
