import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("routines locked on free plan; create routine, start it, log a full session; personal records", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("routines"));
  await page.goto("/app/workouts");

  await page.getByPlaceholder("Nombre de la rutina").fill("Empuje");
  await page.getByPlaceholder("Ejercicios separados por coma").fill("Press banca, Fondos");
  await page.getByRole("button", { name: "Crear rutina →" }).click();
  await expect(page.getByText("Rutinas de entrenamiento es")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/workouts");

  await page.getByPlaceholder("Nombre de la rutina").fill("Empuje");
  await page.getByPlaceholder("Ejercicios separados por coma").fill("Press banca, Fondos");
  await page.getByRole("button", { name: "Crear rutina →" }).click();
  await expect(page.getByText("Empuje")).toBeVisible();
  await expect(page.getByText("Press banca, Fondos")).toBeVisible();

  await page.getByRole("button", { name: "Empezar" }).click();
  await page.getByLabel("Reps Press banca").fill("8");
  await page.getByLabel("Kg Press banca").fill("50");
  await page.getByLabel("Reps Fondos").fill("12");
  await page.getByLabel("Kg Fondos").fill("0");
  await page.getByRole("button", { name: "Terminar rutina" }).click();

  await expect(page.getByText("Rutina: Empuje")).toBeVisible();
  await expect(page.getByText("Press banca: 8 reps @ 50kg")).toBeVisible();
  await expect(page.getByText("Fondos: 12 reps @ 0kg")).toBeVisible();

  await page.getByRole("button", { name: "Ver récords" }).click();
  await expect(page.getByTestId("pr-Press banca").getByText("50kg")).toBeVisible();
});
