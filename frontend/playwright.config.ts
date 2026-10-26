import { defineConfig, devices } from '@playwright/test'

// Runs against a real API (E2E_API, default :8080, seeded synthetic data) behind `vite preview`
// with the production CSP. CI starts Postgres + API before this (see .github/workflows/ci.yml).
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: { baseURL: 'http://127.0.0.1:4173', trace: 'retain-on-failure', locale: 'en-GB' },
  projects: [
    { name: 'desktop', use: { ...devices['Desktop Chrome'] }, grepInvert: /@mobile/ },
    // Keep the Pixel 7 viewport, touch input and mobile UA while avoiding
    // Playwright's mobile viewport scaling, which misroutes pointer events on macOS.
    { name: 'mobile', use: { ...devices['Pixel 7'], isMobile: false }, grep: /@mobile/ },
  ],
  webServer: { command: 'npm run build && npx vite preview', url: 'http://127.0.0.1:4173', reuseExistingServer: !process.env.CI, timeout: 120_000 },
})
