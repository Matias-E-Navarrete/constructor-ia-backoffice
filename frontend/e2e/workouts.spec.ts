import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("log session, history shown, progress locked then unlocked", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("workouts"));
  await page.goto("/app/workouts");

  await page.getByPlaceholder("Ejercicio", { exact: true }).fill("Sentadilla");
  await page.getByPlaceholder("Reps").fill("10");
  await page.getByPlaceholder("Kg").fill("60");
  await page.getByRole("button", { name: "Registrar sesión" }).click();
  await expect(page.getByText("Sentadilla: 10 reps @ 60kg")).toBeVisible();

  await page.getByPlaceholder("Nombre del ejercicio").fill("Sentadilla");
  await page.getByRole("button", { name: "Ver progreso" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/workouts");
  await page.getByPlaceholder("Nombre del ejercicio").fill("Sentadilla");
  await page.getByRole("button", { name: "Ver progreso" }).click();
  await expect(page.getByText(/Récord personal: 60/)).toBeVisible();
});
