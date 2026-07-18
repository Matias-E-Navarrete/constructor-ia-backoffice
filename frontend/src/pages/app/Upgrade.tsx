import { useAuth } from "../../context/AuthContext";
import * as authApi from "../../api/auth";

export default function Upgrade() {
  const { user, refreshUser } = useAuth();

  async function handleUpgrade() {
    await authApi.upgrade();
    await refreshUser();
  }

  if (user?.plan === "pro") {
    return (
      <div className="max-w-md">
        <h1 className="text-xl font-semibold text-slate-900 mb-2">Ya sos Pro 🎉</h1>
        <p className="text-slate-600 text-sm">Tenés acceso a estadísticas, progreso, resúmenes, exportación y coaching.</p>
      </div>
    );
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-semibold text-slate-900 mb-4">Mejorá a Pro</h1>
      <div className="bg-white border border-slate-200 rounded-xl p-6">
        <p className="text-sm text-slate-600 mb-4">Desbloqueá:</p>
        <ul className="text-sm text-slate-700 space-y-1 mb-6 list-disc list-inside">
          <li>Estadísticas de hábitos</li>
          <li>Progreso y récords en entrenamientos</li>
          <li>Resumen financiero y exportación CSV</li>
          <li>Grafo y exportación del vault de notas</li>
          <li>Crear grupos de coaching</li>
        </ul>
        <button onClick={handleUpgrade} className="w-full bg-emerald-600 text-white rounded-md py-2 font-medium">
          Mejorar ahora
        </button>
        <p className="text-xs text-slate-400 mt-2 text-center">
          Demo: activa el plan Pro al instante (sin pasarela de pago real).
        </p>
      </div>
    </div>
  );
}
