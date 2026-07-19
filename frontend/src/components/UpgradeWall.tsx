import { Link } from "react-router-dom";

export default function UpgradeWall({ feature }: { feature: string }) {
  return (
    <div className="border border-dashed border-neutral-300 dark:border-neutral-700 rounded-lg p-6 text-center">
      <p className="text-neutral-900 dark:text-neutral-100 font-medium mb-2">
        {feature} es una función <span className="text-accent-strong dark:text-accent-soft">Pro</span>
      </p>
      <p className="text-sm text-neutral-500 dark:text-neutral-400 mb-4">Mejorá tu plan para desbloquearla.</p>
      <Link to="/app/upgrade" className="btn-accent text-sm">
        Ver planes
      </Link>
    </div>
  );
}
