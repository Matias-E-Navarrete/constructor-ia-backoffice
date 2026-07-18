import { useEffect, useState } from "react";
import * as financeApi from "../../api/finance";
import * as groupsApi from "../../api/groups";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Finance() {
  const [transactions, setTransactions] = useState<financeApi.Transaction[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [type, setType] = useState<"income" | "expense">("expense");
  const [amount, setAmount] = useState(0);
  const [category, setCategory] = useState("general");
  const [familyGroups, setFamilyGroups] = useState<groupsApi.Membership[]>([]);
  const [shareGroupId, setShareGroupId] = useState("");
  const [viewGroupId, setViewGroupId] = useState("");
  const [summary, setSummary] = useState<financeApi.Summary | null>(null);
  const [locked, setLocked] = useState(false);

  async function load(currentViewGroupId = viewGroupId) {
    try {
      setTransactions(await financeApi.listTransactions(currentViewGroupId || undefined));
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
    try {
      const memberships = await groupsApi.listMyGroups();
      setFamilyGroups(memberships.filter((m) => m.group.kind === "family"));
    } catch {
      // groups module may be disabled independently; finance still works without it
    }
  }

  useEffect(() => {
    load();
  }, []);

  function handleViewChange(groupId: string) {
    setViewGroupId(groupId);
    load(groupId);
  }

  async function handleCreate() {
    if (!amount) return;
    await financeApi.recordTransaction({
      type,
      amount,
      category,
      tx_date: new Date().toISOString().slice(0, 10),
      group_id: shareGroupId || undefined,
    });
    setAmount(0);
    load();
  }

  async function handleSummary() {
    setLocked(false);
    setSummary(null);
    try {
      setSummary(await financeApi.getSummary());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  async function handleExport() {
    try {
      await financeApi.exportTransactions();
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-xl font-semibold text-slate-900 mb-4">Finanzas</h1>
        <div className="grid grid-cols-2 gap-2 mb-2">
          <select aria-label="Tipo de movimiento" value={type} onChange={(e) => setType(e.target.value as "income" | "expense")} className="rounded-md border border-slate-300 px-3 py-2">
            <option value="expense">Gasto</option>
            <option value="income">Ingreso</option>
          </select>
          <input
            type="number"
            value={amount || ""}
            onChange={(e) => setAmount(Number(e.target.value))}
            placeholder="Monto"
            className="rounded-md border border-slate-300 px-3 py-2"
          />
          <input
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            placeholder="Categoría"
            className="rounded-md border border-slate-300 px-3 py-2"
          />
          {familyGroups.length > 0 && (
            <select aria-label="Compartir con grupo familiar" value={shareGroupId} onChange={(e) => setShareGroupId(e.target.value)} className="rounded-md border border-slate-300 px-3 py-2">
              <option value="">Personal</option>
              {familyGroups.map((m) => (
                <option key={m.group.id} value={m.group.id}>
                  Compartir con {m.group.name}
                </option>
              ))}
            </select>
          )}
        </div>
        <button onClick={handleCreate} className="bg-slate-900 text-white px-4 py-2 rounded-md text-sm font-medium mb-6">
          Registrar
        </button>

        {familyGroups.length > 0 && (
          <select
            aria-label="Vista de movimientos"
            value={viewGroupId}
            onChange={(e) => handleViewChange(e.target.value)}
            className="mb-4 rounded-md border border-slate-300 px-3 py-2 text-sm"
          >
            <option value="">Mis movimientos personales</option>
            {familyGroups.map((m) => (
              <option key={m.group.id} value={m.group.id}>
                Compartidos en {m.group.name}
              </option>
            ))}
          </select>
        )}

        <ul className="space-y-2">
          {transactions.map((t) => (
            <li key={t.id} className="bg-white border border-slate-200 rounded-lg p-3 flex justify-between text-sm">
              <span>
                {t.category} {t.group_id && <span className="text-xs text-slate-400">(familia)</span>}
              </span>
              <span className={t.type === "income" ? "text-emerald-600" : "text-red-600"}>
                {t.type === "income" ? "+" : "-"}
                {t.amount}
              </span>
            </li>
          ))}
          {transactions.length === 0 && <p className="text-slate-500 text-sm">Todavía no registraste movimientos.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-slate-900 mb-2">Resumen</h2>
        <div className="flex gap-2 mb-4">
          <button onClick={handleSummary} className="bg-slate-100 text-slate-700 px-4 py-2 rounded-md text-sm font-medium">
            Ver resumen
          </button>
          <button onClick={handleExport} className="bg-slate-100 text-slate-700 px-4 py-2 rounded-md text-sm font-medium">
            Exportar CSV
          </button>
        </div>
        {locked && <UpgradeWall feature="Resumen y exportación financiera" />}
        {summary && (
          <div className="text-sm text-slate-600">
            <p>Balance: {summary.balance}</p>
            <p>Ingresos: {summary.income} / Gastos: {summary.expenses}</p>
          </div>
        )}
      </div>
    </div>
  );
}
