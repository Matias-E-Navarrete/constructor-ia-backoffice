import { defineConfig } from "@playwright/test";

const BACKEND_PORT = 8081;
const FRONTEND_PORT = 4173;

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false,
  workers: 1,
  use: {
    baseURL: `http://localhost:${FRONTEND_PORT}`,
    screenshot: "only-on-failure",
    launchOptions: {
      executablePath: "/opt/pw-browsers/chromium",
    },
  },
  webServer: [
    {
      command: "go run ./cmd/api",
      cwd: "../backend",
      url: `http://localhost:${BACKEND_PORT}/health`,
      timeout: 60_000,
      reuseExistingServer: !process.env.CI,
      env: {
        PORT: String(BACKEND_PORT),
        DATABASE_URL: "postgres://rimu:rimu@localhost:5432/rimu_test?sslmode=disable",
        JWT_SECRET: "e2e-test-secret",
        ADMIN_EMAILS: "admin@rimu.test",
      },
    },
    {
      command: `npm run build && npm run preview -- --port ${FRONTEND_PORT} --strictPort`,
      url: `http://localhost:${FRONTEND_PORT}`,
      timeout: 120_000,
      reuseExistingServer: !process.env.CI,
      env: {
        VITE_API_PROXY_TARGET: `http://localhost:${BACKEND_PORT}`,
      },
    },
  ],
});
