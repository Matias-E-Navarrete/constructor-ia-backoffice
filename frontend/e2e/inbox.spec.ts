import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("capture in inbox, pin, organize into a task, delete", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("inbox"));
  await page.goto("/app/inbox");

  await page.getByPlaceholder("Escribí lo que se te ocurra...").fill("Comprar café");
  await page.getByRole("button", { name: "Guardar" }).click();
  await expect(page.getByText("Comprar café")).toBeVisible();

  await page.getByPlaceholder("Escribí lo que se te ocurra...").fill("Llamar al dentista");
  await page.getByRole("button", { name: "Guardar" }).click();
  await expect(page.getByText("Llamar al dentista")).toBeVisible();

  await page.getByPlaceholder("Buscar por contenido...").fill("dentista");
  await expect(page.getByText("Llamar al dentista")).toBeVisible();
  await expect(page.getByText("Comprar café")).not.toBeVisible();
  await page.getByPlaceholder("Buscar por contenido...").fill("");

  const cafeItem = page.locator("li", { hasText: "Comprar café" });
  await cafeItem.getByRole("button", { name: "Pin" }).click();
  await expect(cafeItem.getByRole("button", { name: "Despinnear" })).toBeVisible();

  const dentistaItem = page.locator("li", { hasText: "Llamar al dentista" });
  await dentistaItem.getByRole("button", { name: "Organizar" }).click();
  await page.waitForURL("**/app/tasks");
  await expect(page.getByPlaceholder("ej: Reunión con el equipo")).toHaveValue("Llamar al dentista");
  await page.getByRole("button", { name: "Añadir Tarea →" }).click();
  await expect(page.getByText("Llamar al dentista")).toBeVisible();

  await page.goto("/app/inbox");
  await expect(page.getByText("Llamar al dentista")).not.toBeVisible();

  await page.locator("li", { hasText: "Comprar café" }).getByRole("button", { name: "Eliminar" }).click();
  await expect(page.getByText("Comprar café")).not.toBeVisible();
});
