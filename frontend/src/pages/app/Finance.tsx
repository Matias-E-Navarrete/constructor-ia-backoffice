import { useEffect, useState } from "react";
import * as financeApi from "../../api/finance";
import * as groupsApi from "../../api/groups";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

const CURRENCIES = ["USD", "EUR", "GBP", "CLP", "BRL", "MXN", "ARS"];

export default function Finance() {
  const [transactions, setTransactions] = useState<financeApi.Transaction[]>([]);
  const [accounts, setAccounts] = useState<financeApi.Account[]>([]);
  const [disabled, setDisabled] = useState(false);

  const [type, setType] = useState<"income" | "expense">("expense");
  const [amount, setAmount] = useState(0);
  const [category, setCategory] = useState("General");
  const [accountId, setAccountId] = useState("");
  const [currency, setCurrency] = useState("USD");
  const [method, setMethod] = useState("");
  const [installments, setInstallments] = useState(1);
  const [recurring, setRecurring] = useState(false);
  const [conversion, setConversion] = useState<{ converted_amount: number; rate: number } | null>(null);

  const [newAccountName, setNewAccountName] = useState("");
  const [newAccountCurrency, setNewAccountCurrency] = useState("USD");
  const [newAccountBalance, setNewAccountBalance] = useState(0);

  const [familyGroups, setFamilyGroups] = useState<groupsApi.Membership[]>([]);
  const [shareGroupId, setShareGroupId] = useState("");
  const [viewGroupId, setViewGroupId] = useState("");

  const [summary, setSummary] = useState<financeApi.Summary | null>(null);
  const [upcoming, setUpcoming] = useState<financeApi.UpcomingBill[] | null>(null);
  const [locked, setLocked] = useState(false);
  const [pdfType, setPdfType] = useState("all");

  async function load(currentViewGroupId = viewGroupId) {
    try {
      setTransactions(await financeApi.listTransactions(currentViewGroupId || undefined));
    } catch (err) {
      if (err instanceof FeatureDisabledError) setDisabled(true);
    }
    try {
      setAccounts(await financeApi.listAccounts());
    } catch {
      // ignore
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

  useEffect(() => {
    if (currency === "USD" || !amount) {
      setConversion(null);
      return;
    }
    const handle = setTimeout(() => {
      financeApi.convertCurrency(amount, currency, "USD").then(setConversion).catch(() => setConversion(null));
    }, 250);
    return () => clearTimeout(handle);
  }, [amount, currency]);

  function handleViewChange(groupId: string) {
    setViewGroupId(groupId);
    load(groupId);
  }

  async function handleCreateAccount() {
    if (!newAccountName.trim()) return;
    await financeApi.createAccount({ name: newAccountName.trim(), currency: newAccountCurrency, initial_balance: newAccountBalance });
    setNewAccountName("");
    setNewAccountBalance(0);
    load();
  }

  async function handleCreate() {
    if (!amount) return;
    await financeApi.recordTransaction({
      type,
      amount,
      category,
      tx_date: new Date().toISOString().slice(0, 10),
      group_id: shareGroupId || undefined,
      account_id: accountId || undefined,
      currency,
      exchange_rate: conversion?.rate,
      method: method || undefined,
      installments_total: installments > 1 ? installments : undefined,
      installment_number: installments > 1 ? 1 : undefined,
      recurring,
      recurrence_interval: recurring ? "monthly" : undefined,
    });
    setAmount(0);
    setMethod("");
    setInstallments(1);
    setRecurring(false);
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

  async function handleUpcoming() {
    setLocked(false);
    setUpcoming(null);
    const now = new Date();
    const from = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10);
    const to = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10);
    try {
      setUpcoming(await financeApi.getUpcomingBills(from, to));
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

  async function handleExportPDF() {
    try {
      await financeApi.exportPDF({ type: pdfType });
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  if (disabled) return <FeatureDisabledBanner />;

  const balancesByCurrency = new Map<string, number>();
  for (const a of accounts) {
    balancesByCurrency.set(a.currency, (balancesByCurrency.get(a.currency) ?? 0) + a.initial_balance);
  }
  for (const t of transactions) {
    if (!t.account_id) continue;
    const acc = accounts.find((a) => a.id === t.account_id);
    if (!acc) continue;
    const signed = t.type === "income" ? t.amount : -t.amount;
    balancesByCurrency.set(acc.currency, (balancesByCurrency.get(acc.currency) ?? 0) + signed);
  }

  return (
    <div className="max-w-3xl space-y-8">
      <div>
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-3">Wallet</h1>
        <div className="surface-card p-5 flex flex-wrap gap-x-8 gap-y-1 mb-4">
          {[...balancesByCurrency.entries()].map(([cur, total]) => (
            <span key={cur} className="text-3xl font-light text-neutral-900 dark:text-neutral-50 tabular-nums">
              {total.toFixed(2)} <span className="text-sm text-neutral-400 dark:text-neutral-500">{cur}</span>
            </span>
          ))}
          {balancesByCurrency.size === 0 && <span className="text-neutral-400 dark:text-neutral-500 text-sm">Sin cuentas todavía.</span>}
        </div>

        <details className="mb-4">
          <summary className="text-sm text-neutral-500 dark:text-neutral-400 cursor-pointer">Cuentas ({accounts.length}) · agregar cuenta</summary>
          <div className="mt-2 flex flex-wrap gap-2">
            <input
              value={newAccountName}
              onChange={(e) => setNewAccountName(e.target.value)}
              placeholder="Nombre de la cuenta o tarjeta"
              className="input-field px-3 py-2 text-sm"
            />
            <select
              value={newAccountCurrency}
              onChange={(e) => setNewAccountCurrency(e.target.value)}
              className="input-field px-3 py-2 text-sm"
            >
              {CURRENCIES.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
            <input
              type="number"
              value={newAccountBalance || ""}
              onChange={(e) => setNewAccountBalance(Number(e.target.value))}
              placeholder="Saldo inicial"
              className="input-field px-3 py-2 text-sm w-32"
            />
            <button onClick={handleCreateAccount} className="btn-primary text-sm px-3 py-2">
              Agregar Cuenta →
            </button>
          </div>
          <ul className="mt-2 text-xs text-neutral-500 dark:text-neutral-400 space-y-1">
            {accounts.map((a) => (
              <li key={a.id}>
                {a.name} · {a.currency} {a.initial_balance.toFixed(2)}
              </li>
            ))}
          </ul>
        </details>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mb-2">
          <select aria-label="Tipo de movimiento" value={type} onChange={(e) => setType(e.target.value as "income" | "expense")} className="input-field px-3 py-2">
            <option value="expense">Gasto</option>
            <option value="income">Ingreso</option>
          </select>
          <div className="flex gap-1">
            <input
              type="number"
              value={amount || ""}
              onChange={(e) => setAmount(Number(e.target.value))}
              placeholder="Monto"
              className="flex-1 input-field px-3 py-2"
            />
            <select
              aria-label="Moneda"
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              className="input-field px-2 py-2 text-sm"
            >
              {CURRENCIES.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
          </div>
          <input
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            placeholder="Categoría"
            className="input-field px-3 py-2"
          />
          <select
            aria-label="Cuenta"
            value={accountId}
            onChange={(e) => setAccountId(e.target.value)}
            className="input-field px-3 py-2"
          >
            <option value="">Sin cuenta</option>
            {accounts.map((a) => (
              <option key={a.id} value={a.id}>
                {a.name}
              </option>
            ))}
          </select>
          <input
            value={method}
            onChange={(e) => setMethod(e.target.value)}
            placeholder="Método (opcional)"
            className="input-field px-3 py-2"
          />
          <div className="flex items-center gap-2 text-sm text-neutral-600 dark:text-neutral-400">
            <label className="flex items-center gap-1">
              Cuotas
              <input
                type="number"
                min={1}
                value={installments}
                onChange={(e) => setInstallments(Number(e.target.value))}
                className="w-14 input-field px-2 py-1"
              />
            </label>
            <label className="flex items-center gap-1">
              <input type="checkbox" checked={recurring} onChange={(e) => setRecurring(e.target.checked)} />
              Recurrente
            </label>
          </div>
          {familyGroups.length > 0 && (
            <select aria-label="Compartir con grupo familiar" value={shareGroupId} onChange={(e) => setShareGroupId(e.target.value)} className="input-field px-3 py-2">
              <option value="">Personal</option>
              {familyGroups.map((m) => (
                <option key={m.group.id} value={m.group.id}>
                  Compartir con {m.group.name}
                </option>
              ))}
            </select>
          )}
        </div>
        {conversion && (
          <p className="text-xs text-neutral-400 dark:text-neutral-500 mb-2">
            ≈ USD {conversion.converted_amount.toFixed(2)} · tasa {conversion.rate.toFixed(4)}
          </p>
        )}
        <button onClick={handleCreate} className="btn-primary mb-6">
          Registrar
        </button>

        {familyGroups.length > 0 && (
          <select
            aria-label="Vista de movimientos"
            value={viewGroupId}
            onChange={(e) => handleViewChange(e.target.value)}
            className="mb-4 input-field px-3 py-2 text-sm"
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
            <li key={t.id} className="surface-card-hover p-3 flex justify-between text-sm text-neutral-800 dark:text-neutral-200">
              <span>
                {t.category}
                {t.installments_total && t.installments_total > 1 && (
                  <span className="text-xs text-neutral-400 dark:text-neutral-500"> · {t.installment_number}/{t.installments_total}</span>
                )}
                {t.recurring && <span className="text-xs text-neutral-400 dark:text-neutral-500"> · recurrente</span>}
                {t.group_id && <span className="text-xs text-neutral-400 dark:text-neutral-500"> (familia)</span>}
              </span>
              <span className={t.type === "income" ? "font-medium text-neutral-900 dark:text-neutral-50" : "text-neutral-500 dark:text-neutral-400"}>
                {t.type === "income" ? "+" : "−"}
                {t.amount} {t.currency}
              </span>
            </li>
          ))}
          {transactions.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no registraste movimientos.</p>}
        </ul>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-neutral-900 dark:text-neutral-50 mb-2">Resumen y próximas cuentas</h2>
        <div className="flex flex-wrap gap-2 mb-4">
          <button onClick={handleSummary} className="btn-secondary">
            Ver resumen
          </button>
          <button onClick={handleUpcoming} className="btn-secondary">
            Próximas cuentas
          </button>
          <button onClick={handleExport} className="btn-secondary">
            Exportar CSV
          </button>
          <select
            value={pdfType}
            onChange={(e) => setPdfType(e.target.value)}
            className="input-field px-2 py-2 text-sm"
          >
            <option value="all">Todos</option>
            <option value="income">Solo ingresos</option>
            <option value="expense">Solo gastos</option>
          </select>
          <button onClick={handleExportPDF} className="btn-secondary">
            Exportar PDF
          </button>
        </div>
        {locked && <UpgradeWall feature="Resumen, próximas cuentas y exportación" />}
        {summary && (
          <div className="text-sm text-neutral-600 dark:text-neutral-400 mb-4">
            <p>Balance: {summary.balance.toFixed(2)}</p>
            <p>Ingresos: {summary.income.toFixed(2)} / Gastos: {summary.expenses.toFixed(2)}</p>
          </div>
        )}
        {upcoming && (
          <ul className="text-sm text-neutral-600 dark:text-neutral-400 space-y-1">
            {upcoming.map((b, i) => (
              <li key={i} className="flex justify-between">
                <span>
                  {b.transaction.category} · {b.next_date}
                </span>
                <span>
                  −{b.transaction.amount} {b.transaction.currency}
                </span>
              </li>
            ))}
            {upcoming.length === 0 && <li>No hay cuentas próximas en este período.</li>}
          </ul>
        )}
      </div>
    </div>
  );
}
