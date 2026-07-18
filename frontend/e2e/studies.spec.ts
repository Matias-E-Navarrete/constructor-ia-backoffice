import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("create subject, log a study session, overview locked then unlocked", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("studies"));
  await page.goto("/app/studies");

  await page.getByPlaceholder("Nueva materia").fill("Cálculo I");
  await page.getByRole("button", { name: "Agregar" }).click();
  await expect(page.getByText("Cálculo I")).toBeVisible();

  await page.getByRole("button", { name: "Registrar sesión" }).click();
  await page.getByLabel("Minutos").fill("45");
  await page.getByPlaceholder("Tema (opcional)").fill("Límites");
  await page.getByRole("button", { name: "Guardar" }).click();

  await expect(page.getByText("Cálculo I · Límites")).toBeVisible();
  await expect(page.getByText(/45min/)).toBeVisible();

  await page.getByRole("button", { name: "Ver panorama" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/studies");
  await page.getByRole("button", { name: "Ver panorama" }).click();

  const row = page.locator("tr", { hasText: "Cálculo I" });
  await expect(row).toBeVisible();
  await expect(row.getByRole("cell").nth(1)).toHaveText("1");
  await expect(row.getByRole("cell").nth(3)).toHaveText("1 días");
});
