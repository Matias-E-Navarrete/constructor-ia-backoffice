export default function FeatureDisabledBanner() {
  return (
    <div className="border border-dashed border-neutral-300 dark:border-neutral-700 bg-neutral-100 dark:bg-neutral-900 rounded-lg p-6 text-center">
      <p className="text-neutral-700 dark:text-neutral-300 font-medium">Esta función está temporalmente deshabilitada.</p>
      <p className="text-sm text-neutral-500 dark:text-neutral-400 mt-1">Un administrador la reactivará pronto.</p>
    </div>
  );
}
