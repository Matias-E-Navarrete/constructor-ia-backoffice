import { useState } from "react";
import * as adminApi from "../../api/admin";

export default function AdminTestRunner() {
  const [running, setRunning] = useState<"backend" | "frontend" | null>(null);
  const [result, setResult] = useState<adminApi.TestRunResult | null>(null);

  async function handleRun(suite: "backend" | "frontend") {
    setRunning(suite);
    setResult(null);
    try {
      setResult(await adminApi.runTests(suite));
    } finally {
      setRunning(null);
    }
  }

  return (
    <div className="max-w-2xl">
      <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-50 mb-4">Tests</h1>
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => handleRun("backend")}
          disabled={running !== null}
          className="bg-slate-900 text-white dark:bg-white dark:text-slate-900 px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
        >
          {running === "backend" ? "Corriendo..." : "Correr tests backend"}
        </button>
        <button
          onClick={() => handleRun("frontend")}
          disabled={running !== null}
          className="bg-slate-900 text-white dark:bg-white dark:text-slate-900 px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
        >
          {running === "frontend" ? "Corriendo..." : "Correr tests frontend"}
        </button>
      </div>

      {result && (
        <div
          className={`border rounded-lg p-4 ${
            result.success
              ? "border-emerald-300 dark:border-emerald-800 bg-emerald-50 dark:bg-emerald-950"
              : "border-red-300 dark:border-red-800 bg-red-50 dark:bg-red-950"
          }`}
        >
          <p className="font-medium mb-2 text-slate-900 dark:text-slate-100">{result.success ? "✅ Todo verde" : "❌ Hay fallas"}</p>
          {(result.passed > 0 || result.failed > 0) && (
            <p className="text-sm mb-2 text-slate-700 dark:text-slate-300">
              {result.passed} pasaron, {result.failed} fallaron
            </p>
          )}
          <pre className="text-xs bg-white/60 dark:bg-black/40 text-slate-800 dark:text-slate-200 p-3 rounded max-h-80 overflow-auto whitespace-pre-wrap">{result.output}</pre>
        </div>
      )}
    </div>
  );
}
