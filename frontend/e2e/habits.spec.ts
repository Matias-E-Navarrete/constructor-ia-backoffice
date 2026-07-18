import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("create habit, check in, stats locked then unlocked", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("habits"));

  await page.getByPlaceholder("Nuevo hábito").fill("Meditar");
  await page.getByRole("button", { name: "Agregar" }).click();
  await expect(page.getByText("Meditar")).toBeVisible();

  await page.getByRole("button", { name: "Check-in" }).click();
  await page.getByRole("button", { name: "Estadísticas" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/habits");
  await page.getByRole("button", { name: "Estadísticas" }).click();
  await expect(page.getByText(/Racha actual: 1/)).toBeVisible();
});
