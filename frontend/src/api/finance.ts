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
};
export type Summary = {
  balance: number;
  income: number;
  expenses: number;
  by_category: { category: string; total: number }[];
};

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
