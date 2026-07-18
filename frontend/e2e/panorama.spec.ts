import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("panorama: locked on free plan, shows heatmap/streaks/radar on pro", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("panorama"));
  await page.goto("/app/habits");

  await page.getByPlaceholder("Nuevo hábito").fill("Meditar");
  await page.getByRole("button", { name: "Agregar" }).click();
  await expect(page.getByText("Meditar")).toBeVisible();
  await page.getByRole("button", { name: "Check-in" }).click();

  await page.getByRole("link", { name: "Ver Panorama →" }).click();
  await page.waitForURL("**/app/habits/panorama");
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/habits/panorama");

  await expect(page.locator('[data-testid^="streak-ring-"]').first()).toBeVisible();
  await expect(page.locator('[data-testid^="heatmap-"]').first()).toBeVisible();

  const row = page.locator("tr", { hasText: "Meditar" });
  await expect(row).toBeVisible();
  await expect(row.getByRole("cell").nth(1)).toHaveText("1");
});
