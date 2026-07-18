import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("create task, classify in Eisenhower, move in Kanban", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("tasks"));
  await page.goto("/app/tasks");

  await page.getByPlaceholder("ej: Reunión con el equipo").fill("Finish Product Design Course");
  await page.getByRole("button", { name: "Añadir Tarea →" }).click();
  await expect(page.getByText("Finish Product Design Course")).toBeVisible();

  await page.getByRole("button", { name: "Eisenhower" }).click();
  const card = page.getByText("Finish Product Design Course");
  const doNowQuadrant = page.getByTestId("quadrant-do_now");
  await card.dragTo(doNowQuadrant);
  await expect(doNowQuadrant.getByText("Finish Product Design Course")).toBeVisible();

  await page.getByRole("button", { name: "Kanban" }).click();
  const kanbanCard = page.getByText("Finish Product Design Course");
  const inProgressColumn = page.getByTestId("column-in_progress");
  await kanbanCard.dragTo(inProgressColumn);
  await expect(inProgressColumn.getByText("Finish Product Design Course")).toBeVisible();
});
