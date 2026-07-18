import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("free account hits an upgrade wall, upgrades, and the locked tab unlocks", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("freemium"));

  await expect(page.getByText("FREE", { exact: true })).toBeVisible();

  await page.goto("/app/finance");
  await page.getByRole("button", { name: "Ver resumen" }).click();
  await expect(page.getByText("es una función Pro")).toBeVisible();
  await expect(page.getByText(/Balance:/)).toHaveCount(0);

  await page.goto("/app/upgrade");
  await page.getByRole("button", { name: "Mejorar ahora" }).click();
  await expect(page.getByText("Ya sos Pro")).toBeVisible();
  await expect(page.getByText("PRO", { exact: true })).toBeVisible();

  await page.goto("/app/finance");
  await page.getByRole("button", { name: "Ver resumen" }).click();
  await expect(page.getByText(/Balance:/)).toBeVisible();
});
