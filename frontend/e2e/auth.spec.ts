import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("register, session persists, logout", async ({ page }) => {
  const email = uniqueEmail("auth");
  await registerAndLogin(page, email);

  await expect(page.getByText(email)).toBeVisible();

  await page.getByRole("button", { name: "Cerrar sesión" }).click();
  await page.waitForURL("**/login");

  await expect(page.getByRole("button", { name: "Ingresar" })).toBeVisible();
});
