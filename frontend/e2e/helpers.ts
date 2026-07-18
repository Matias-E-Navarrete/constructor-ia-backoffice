import type { Page } from "@playwright/test";

export function uniqueEmail(prefix: string): string {
  return `${prefix}-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`;
}

export async function registerAndLogin(page: Page, email: string, password = "password123") {
  await page.goto("/register");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Contraseña").fill(password);
  await page.getByRole("button", { name: "Crear cuenta" }).click();
  await page.waitForURL("**/app/tasks");
}

export async function login(page: Page, email: string, password = "password123") {
  await page.goto("/login");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Contraseña").fill(password);
  await page.getByRole("button", { name: "Ingresar" }).click();
  await page.waitForURL("**/app/tasks");
}

// registerOrLogin is for fixed, non-unique emails (like the ADMIN_EMAILS
// bootstrap address) that may already be registered from a previous test run
// against the same database: it tries to register, and falls back to login
// if the account already exists.
export async function registerOrLogin(page: Page, email: string, password = "password123") {
  await page.goto("/register");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Contraseña").fill(password);
  await page.getByRole("button", { name: "Crear cuenta" }).click();
  try {
    await page.waitForURL("**/app/tasks", { timeout: 5_000 });
  } catch {
    await login(page, email, password);
  }
}

export async function upgradeToPro(page: Page) {
  await page.goto("/app/upgrade");
  await page.getByRole("button", { name: "Mejorar ahora" }).click();
  await page.getByText("Ya sos Pro").waitFor();
}
