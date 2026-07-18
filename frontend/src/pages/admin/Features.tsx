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
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Feature flags</h1>
      <ul className="space-y-2">
        {flags.map((f) => (
          <li key={f.key} className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-neutral-800 dark:text-neutral-200">{f.key}</p>
              <p className="text-xs text-neutral-500 dark:text-neutral-400">{f.description}</p>
            </div>
            <button
              onClick={() => handleToggle(f.key, !f.enabled)}
              className={`text-xs font-medium px-3 py-1 rounded-full border ${
                f.enabled
                  ? "bg-neutral-900 text-white border-neutral-900 dark:bg-white dark:text-neutral-900 dark:border-white"
                  : "text-neutral-400 dark:text-neutral-500 border-neutral-300 dark:border-neutral-700"
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
