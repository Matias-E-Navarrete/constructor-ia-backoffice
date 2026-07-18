export default function FeatureDisabledBanner() {
  return (
    <div className="border border-dashed border-slate-300 bg-slate-100 rounded-lg p-6 text-center">
      <p className="text-slate-700 font-medium">Esta función está temporalmente deshabilitada.</p>
      <p className="text-sm text-slate-500 mt-1">Un administrador la reactivará pronto.</p>
    </div>
  );
}
