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
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-2">
          Ya sos <span className="text-accent-strong dark:text-accent-soft">Pro</span> 🎉
        </h1>
        <p className="text-neutral-600 dark:text-neutral-400 text-sm">Tenés acceso a estadísticas, progreso, resúmenes, exportación y coaching.</p>
      </div>
    );
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Mejorá a Pro</h1>
      <div className="surface-card p-6">
        <p className="text-sm text-neutral-600 dark:text-neutral-400 mb-4">Desbloqueá:</p>
        <ul className="text-sm text-neutral-700 dark:text-neutral-300 space-y-1 mb-6 list-disc list-inside">
          <li>Estadísticas y panorama de hábitos</li>
          <li>Rutinas, récords y progreso en entrenamientos</li>
          <li>Panorama de estudio por materia</li>
          <li>Resumen financiero, próximas cuentas y exportación</li>
          <li>Grafo y exportación del vault de notas</li>
          <li>Crear grupos de coaching</li>
        </ul>
        <button onClick={handleUpgrade} className="btn-accent w-full py-2">
          Mejorar ahora
        </button>
        <p className="text-xs text-neutral-400 dark:text-neutral-500 mt-2 text-center">
          Demo: activa el plan Pro al instante (sin pasarela de pago real).
        </p>
      </div>
    </div>
  );
}
