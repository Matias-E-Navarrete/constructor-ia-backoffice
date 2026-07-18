import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("schedule a task on the planner timeline and adjust planner hours", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("planner"));

  await page.goto("/app/tasks");
  await page.getByPlaceholder("ej: Reunión con el equipo").fill("Revisar métricas");
  await page.getByRole("button", { name: "Añadir Tarea →" }).click();
  await expect(page.getByText("Revisar métricas")).toBeVisible();

  await page.goto("/app/planner");
  const card = page.getByText("Revisar métricas");
  const hourSlot = page.getByTestId("hour-10");
  await card.dragTo(hourSlot);
  await expect(hourSlot.getByText("Revisar métricas")).toBeVisible();

  await page.getByRole("button", { name: "⚙ Horas" }).click();
  await page.getByLabel("Empieza a las").fill("8");
  await page.getByLabel("Termina a las").fill("20");
  await page.getByRole("button", { name: "Guardar" }).click();
  await expect(page.getByTestId("hour-8")).toBeVisible();
  await expect(page.getByTestId("hour-6")).toHaveCount(0);
});

test("focus mode: complete a task and capture a quick note", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("focus"));

  await page.goto("/app/tasks");
  await page.getByPlaceholder("ej: Reunión con el equipo").fill("Escribir informe");
  await page.getByRole("button", { name: "Añadir Tarea →" }).click();
  await expect(page.getByText("Escribir informe")).toBeVisible();

  await page.goto("/app/planner");
  await page.getByRole("button", { name: "▶ Modo Foco" }).click();
  const focusMode = page.getByTestId("focus-mode");
  await expect(focusMode.getByText("Escribir informe")).toBeVisible();

  await focusMode.getByPlaceholder("Nota rápida...").fill("Idea suelta");
  await focusMode.getByRole("button", { name: "➤" }).click();

  await focusMode.getByRole("button", { name: "Completar" }).click();
  await expect(focusMode.getByText("Tarea completada: Escribir informe")).toBeVisible();

  await focusMode.getByRole("button", { name: "Terminar Flow" }).click();
  await expect(page.getByTestId("focus-mode")).toHaveCount(0);

  await page.goto("/app/inbox");
  await expect(page.getByText("Idea suelta")).toBeVisible();
});
