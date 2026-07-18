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
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Tests</h1>
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => handleRun("backend")}
          disabled={running !== null}
          className="bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
        >
          {running === "backend" ? "Corriendo..." : "Correr tests backend"}
        </button>
        <button
          onClick={() => handleRun("frontend")}
          disabled={running !== null}
          className="bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
        >
          {running === "frontend" ? "Corriendo..." : "Correr tests frontend"}
        </button>
      </div>

      {result && (
        <div
          className={`border rounded-lg p-4 bg-white dark:bg-neutral-900 border-neutral-200 dark:border-neutral-800 ${
            result.success ? "" : "border-l-4 border-l-red-500 dark:border-l-red-500"
          }`}
        >
          <p className="font-medium mb-2 text-neutral-900 dark:text-neutral-100">{result.success ? "✅ Todo verde" : "❌ Hay fallas"}</p>
          {(result.passed > 0 || result.failed > 0) && (
            <p className="text-sm mb-2 text-neutral-700 dark:text-neutral-300">
              {result.passed} pasaron, {result.failed} fallaron
            </p>
          )}
          <pre className="text-xs bg-white/60 dark:bg-black/40 text-neutral-800 dark:text-neutral-200 p-3 rounded max-h-80 overflow-auto whitespace-pre-wrap">{result.output}</pre>
        </div>
      )}
    </div>
  );
}
