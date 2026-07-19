import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("assistant: locked on free plan, saves the message and shows not-configured state on pro", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("assistant"));
  await page.goto("/app/assistant");

  await expect(page.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/assistant");

  await page.getByPlaceholder("Escribí tu mensaje...").fill("Hola, ¿qué tareas tengo hoy?");
  await page.getByRole("button", { name: "Enviar" }).click();

  await expect(page.getByText("Hola, ¿qué tareas tengo hoy?")).toBeVisible();
  await expect(page.getByText(/todavía no está configurado/)).toBeVisible();

  await page.reload();
  await expect(page.getByText("Hola, ¿qué tareas tengo hoy?")).toBeVisible();

  await page.getByRole("button", { name: "Borrar conversación" }).click();
  await expect(page.getByText("Hola, ¿qué tareas tengo hoy?")).not.toBeVisible();
});
