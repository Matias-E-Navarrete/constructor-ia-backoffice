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
          className="btn-primary"
        >
          {running === "backend" ? "Corriendo..." : "Correr tests backend"}
        </button>
        <button
          onClick={() => handleRun("frontend")}
          disabled={running !== null}
          className="btn-primary"
        >
          {running === "frontend" ? "Corriendo..." : "Correr tests frontend"}
        </button>
      </div>

      {result && (
        <div
          className={`surface-card p-4 ${result.success ? "" : "border-l-4 border-l-red-500 dark:border-l-red-500"}`}
        >
          <p className="font-medium mb-2 text-neutral-900 dark:text-neutral-100">{result.success ? "✅ Todo verde" : "❌ Hay fallas"}</p>
          {(result.passed > 0 || result.failed > 0) && (
            <p className="text-sm mb-2 text-neutral-700 dark:text-neutral-300">
              {result.passed} pasaron, {result.failed} fallaron
            </p>
          )}
          <pre className="text-xs bg-white/60 dark:bg-black/40 text-neutral-800 dark:text-neutral-200 p-3 rounded-md max-h-80 overflow-auto whitespace-pre-wrap">{result.output}</pre>
        </div>
      )}
    </div>
  );
}
