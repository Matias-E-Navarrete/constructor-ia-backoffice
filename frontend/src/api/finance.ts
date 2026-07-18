import { apiFetch, downloadFile } from "./client";

export type Transaction = {
  id: string;
  type: "income" | "expense";
  amount: number;
  category: string;
  description: string;
  tx_date: string;
  note_slug?: string;
  group_id?: string;
  account_id?: string;
  currency: string;
  exchange_rate?: number;
  method: string;
  installments_total?: number;
  installment_number?: number;
  recurring: boolean;
  recurrence_interval?: string;
};
export type Summary = {
  balance: number;
  income: number;
  expenses: number;
  by_category: { category: string; total: number }[];
};
export type Account = {
  id: string;
  name: string;
  type: string;
  currency: string;
  initial_balance: number;
  active: boolean;
};
export type UpcomingBill = { transaction: Transaction; next_date: string };

export function listTransactions(groupId?: string) {
  return apiFetch<Transaction[]>("/api/finance/transactions", { query: { groupId } });
}

export function recordTransaction(input: {
  type: "income" | "expense";
  amount: number;
  category: string;
  description?: string;
  tx_date: string;
  note_slug?: string;
  group_id?: string;
  account_id?: string;
  currency?: string;
  exchange_rate?: number;
  method?: string;
  installments_total?: number;
  installment_number?: number;
  recurring?: boolean;
  recurrence_interval?: string;
}) {
  return apiFetch<Transaction>("/api/finance/transactions", { method: "POST", body: input });
}

export function deleteTransaction(id: string) {
  return apiFetch<void>(`/api/finance/transactions/${id}`, { method: "DELETE" });
}

export function getSummary() {
  return apiFetch<Summary>("/api/finance/summary");
}

export function exportTransactions() {
  return downloadFile("/api/finance/export", "transactions.csv");
}

export function listAccounts() {
  return apiFetch<Account[]>("/api/finance/accounts");
}

export function createAccount(input: { name: string; type?: string; currency?: string; initial_balance?: number }) {
  return apiFetch<Account>("/api/finance/accounts", { method: "POST", body: input });
}

export function convertCurrency(amount: number, from: string, to: string) {
  return apiFetch<{ converted_amount: number; rate: number }>("/api/finance/convert", {
    query: { amount: String(amount), from, to },
  });
}

export function getUpcomingBills(from: string, to: string) {
  return apiFetch<UpcomingBill[]>("/api/finance/upcoming", { query: { from, to } });
}

export function exportPDF(filters: {
  type?: string;
  from?: string;
  to?: string;
  categories?: string[];
  account_ids?: string[];
}) {
  return downloadFile("/api/finance/export/pdf", "informe.pdf", filters);
}
