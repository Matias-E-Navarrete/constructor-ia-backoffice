import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("create note, link a habit to it, confirm backlink in graph", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("notes"));
  await page.goto("/app/notes");

  await page.getByPlaceholder("Título").fill("Meta de meditacion");
  await page.getByRole("button", { name: "Guardar nota" }).click();
  await expect(page.getByText("Meta de meditacion")).toBeVisible();

  await page.goto("/app/habits");
  await page.getByPlaceholder("Nuevo hábito").fill("Meditar");
  await page.getByPlaceholder("Vincular a nota (slug, opcional)").fill("meta-de-meditacion");
  await page.getByRole("button", { name: "Agregar" }).click();
  await expect(page.getByText("Meditar")).toBeVisible();

  await upgradeToPro(page);
  await page.goto("/app/notes/graph");
  await expect(page.getByText("Meta de meditacion")).toBeVisible();
  await expect(page.getByText("Meditar")).toBeVisible();
});
