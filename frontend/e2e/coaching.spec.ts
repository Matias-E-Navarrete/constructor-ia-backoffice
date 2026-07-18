import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, upgradeToPro } from "./helpers";

test("upgrade, create coaching group, invite client, view their habits read-only", async ({ browser }) => {
  const coachCtx = await browser.newContext();
  const clientCtx = await browser.newContext();
  const coachPage = await coachCtx.newPage();
  const clientPage = await clientCtx.newPage();

  const coachEmail = uniqueEmail("coach");
  const clientEmail = uniqueEmail("client");

  await registerAndLogin(coachPage, coachEmail);
  await registerAndLogin(clientPage, clientEmail);

  await clientPage.goto("/app/habits");
  await clientPage.getByPlaceholder("Nuevo hábito").fill("Correr");
  await clientPage.getByRole("button", { name: "Agregar" }).click();
  await expect(clientPage.getByText("Correr")).toBeVisible();

  await coachPage.goto("/app/coaching");
  await coachPage.getByPlaceholder("Nombre del grupo de coaching").fill("Mi coaching");
  await coachPage.getByRole("button", { name: "Crear" }).click();
  await expect(coachPage.getByText("es una función Pro")).toBeVisible();

  await upgradeToPro(coachPage);
  await coachPage.goto("/app/coaching");
  await coachPage.getByPlaceholder("Nombre del grupo de coaching").fill("Mi coaching");
  await coachPage.getByRole("button", { name: "Crear" }).click();
  await expect(coachPage.getByText("Mi coaching")).toBeVisible();

  await coachPage.getByText("Mi coaching").click();
  await coachPage.getByPlaceholder("Email del cliente").fill(clientEmail);
  await coachPage.getByRole("button", { name: "Invitar" }).click();

  await coachPage.getByRole("button", { name: "Ver progreso" }).click();
  await expect(coachPage.getByText("Correr")).toBeVisible();

  await coachCtx.close();
  await clientCtx.close();
});
