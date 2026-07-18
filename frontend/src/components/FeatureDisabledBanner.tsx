export default function FeatureDisabledBanner() {
  return (
    <div className="border border-dashed border-slate-300 dark:border-neutral-700 bg-slate-100 dark:bg-neutral-900 rounded-lg p-6 text-center">
      <p className="text-slate-700 dark:text-slate-300 font-medium">Esta función está temporalmente deshabilitada.</p>
      <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Un administrador la reactivará pronto.</p>
    </div>
  );
}
