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
      <h1 className="text-xl font-semibold text-slate-900 mb-4">Feature flags</h1>
      <ul className="space-y-2">
        {flags.map((f) => (
          <li key={f.key} className="bg-white border border-slate-200 rounded-lg p-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-slate-800">{f.key}</p>
              <p className="text-xs text-slate-500">{f.description}</p>
            </div>
            <button
              onClick={() => handleToggle(f.key, !f.enabled)}
              className={`text-xs font-medium px-3 py-1 rounded-full ${
                f.enabled ? "bg-emerald-100 text-emerald-700" : "bg-red-100 text-red-700"
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
