import { Link } from "react-router-dom";

export default function UpgradeWall({ feature }: { feature: string }) {
  return (
    <div className="border border-dashed border-emerald-300 bg-emerald-50 rounded-lg p-6 text-center">
      <p className="text-emerald-900 font-medium mb-2">{feature} es una función Pro</p>
      <p className="text-sm text-emerald-700 mb-4">Mejorá tu plan para desbloquearla.</p>
      <Link to="/app/upgrade" className="inline-block bg-emerald-600 text-white text-sm font-medium px-4 py-2 rounded-md">
        Ver planes
      </Link>
    </div>
  );
}
