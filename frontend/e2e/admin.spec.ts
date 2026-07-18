import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin, registerOrLogin } from "./helpers";

test("non-admin is redirected out of /admin; admin kill switch disables a user-facing tab", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("nonadmin"));
  await page.goto("/admin");
  await page.waitForURL("**/app/habits");

  await registerOrLogin(page, "admin@rimu.test");
  await page.goto("/admin/features");
  await expect(page.getByText("habits", { exact: true })).toBeVisible();

  const habitsFlagRow = page.locator("li", { hasText: "habits" }).first();
  await habitsFlagRow.getByRole("button").click();
  await expect(habitsFlagRow.getByText("Apagado")).toBeVisible();

  await page.goto("/app/habits");
  await expect(page.getByText("temporalmente deshabilitada")).toBeVisible();

  await page.goto("/admin/features");
  await page.locator("li", { hasText: "habits" }).first().getByRole("button").click();
  await page.goto("/app/habits");
  await expect(page.getByPlaceholder("Nuevo hábito")).toBeVisible();
});
