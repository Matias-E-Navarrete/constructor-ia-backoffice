import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("multi-currency account, live conversion, cuotas + recurring, upcoming/PDF gating", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("wallet"));
  await page.goto("/app/finance");

  await page.getByText(/Cuentas \(0\)/).click();
  await page.getByPlaceholder("Nombre de la cuenta o tarjeta").fill("Tarjeta Visa");
  await page.getByPlaceholder("Saldo inicial").fill("1000");
  await page.getByRole("button", { name: "Agregar Cuenta →" }).click();
  await expect(page.getByText("Tarjeta Visa · USD 1000.00")).toBeVisible();

  await page.getByPlaceholder("Monto").fill("100");
  await page.getByLabel("Moneda").selectOption("EUR");
  await expect(page.getByText(/≈ USD/)).toBeVisible();

  await page.getByLabel("Cuenta").selectOption({ label: "Tarjeta Visa" });
  await page.getByPlaceholder("Categoría").fill("Electrónica");
  await page.getByPlaceholder("Método (opcional)").fill("Visa Débito");
  await page.locator('label:has-text("Cuotas") input[type="number"]').fill("3");
  await page.locator('label:has-text("Recurrente") input[type="checkbox"]').check();
  await page.getByRole("button", { name: "Registrar" }).click();

  const txItem = page.locator("li", { hasText: "Electrónica" });
  await expect(txItem).toBeVisible();
  await expect(txItem.getByText("1/3")).toBeVisible();
  await expect(txItem.getByText("recurrente")).toBeVisible();

  await page.getByRole("button", { name: "Próximas cuentas" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await page.getByRole("button", { name: "Exportar PDF" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/finance");
  await page.getByRole("button", { name: "Próximas cuentas" }).click();
  await expect(page.getByText("es una función Pro")).not.toBeVisible();
});
