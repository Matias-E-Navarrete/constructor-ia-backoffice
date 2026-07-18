import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("invite a family member, share a transaction, both see it", async ({ browser }) => {
  const ownerCtx = await browser.newContext();
  const memberCtx = await browser.newContext();
  const ownerPage = await ownerCtx.newPage();
  const memberPage = await memberCtx.newPage();

  const ownerEmail = uniqueEmail("owner");
  const memberEmail = uniqueEmail("member");

  await registerAndLogin(ownerPage, ownerEmail);
  await registerAndLogin(memberPage, memberEmail);

  await ownerPage.goto("/app/family");
  await ownerPage.getByPlaceholder("Nombre del grupo familiar").fill("Familia Test");
  await ownerPage.getByRole("button", { name: "Crear" }).click();
  await ownerPage.getByText("Familia Test").click();
  await ownerPage.getByPlaceholder("Email a invitar").fill(memberEmail);
  await ownerPage.getByRole("button", { name: "Invitar" }).click();

  await ownerPage.goto("/app/finance");
  await ownerPage.getByPlaceholder("Monto").fill("30");
  await ownerPage.getByPlaceholder("Categoría").fill("super");
  await ownerPage.getByLabel("Compartir con grupo familiar").selectOption({ label: "Compartir con Familia Test" });
  await ownerPage.getByRole("button", { name: "Registrar" }).click();
  await expect(ownerPage.getByText("super")).toBeVisible();

  await memberPage.goto("/app/finance");
  await memberPage.getByLabel("Vista de movimientos").selectOption({ label: "Compartidos en Familia Test" });
  await expect(memberPage.getByText("super")).toBeVisible();

  await ownerCtx.close();
  await memberCtx.close();
});
