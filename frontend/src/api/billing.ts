import { apiFetch } from "./client";

export function getBillingConfig() {
  return apiFetch<{ configured: boolean }>("/api/billing/config");
}

export function createCheckoutSession() {
  return apiFetch<{ checkout_url: string }>("/api/billing/checkout", { method: "POST" });
}
