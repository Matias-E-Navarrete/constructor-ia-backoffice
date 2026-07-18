import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("record transaction, summary locked then unlocked", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("finance"));
  await page.goto("/app/finance");

  await page.getByPlaceholder("Monto").fill("50");
  await page.getByPlaceholder("Categoría").fill("comida");
  await page.getByRole("button", { name: "Registrar" }).click();
  await expect(page.getByText("comida")).toBeVisible();

  await page.getByRole("button", { name: "Ver resumen" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/finance");
  await page.getByRole("button", { name: "Ver resumen" }).click();
  await expect(page.getByText(/Balance: -50/)).toBeVisible();
});
