import { test, expect } from "@playwright/test";
import { uniqueEmail, registerAndLogin } from "./helpers";

test("upgrade page falls back to the mock flow when Stripe isn't configured", async ({ page }) => {
  await registerAndLogin(page, uniqueEmail("billing"));
  await page.goto("/app/upgrade");

  await expect(page.getByText(/sin pasarela de pago real/)).toBeVisible();

  await page.getByRole("button", { name: "Mejorar ahora" }).click();
  await expect(page.getByText("Ya sos")).toBeVisible();
  await expect(page.getByText("Pro", { exact: true })).toBeVisible();
});
