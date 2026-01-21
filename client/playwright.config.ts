import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright E2E Test Configuration for GameLink Client
 *
 * Environment Variables:
 * - BASE_URL: Client URL (default: http://localhost:5173)
 * - API_URL: Backend API URL (default: http://localhost:5000/api/v1)
 *
 * Usage:
 * - npx playwright test                    # Run all tests in headless mode
 * - npx playwright test --ui               # Run tests with UI mode
 * - npx playwright test --headed           # Run tests in headed mode
 * - npx playwright test --project=chromium # Run tests on Chromium only
 * - npx playwright test --debug            # Debug tests with inspector
 */
export default defineConfig({
  testDir: './tests/e2e',

  /* Run tests in files in parallel */
  fullyParallel: false,

  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,

  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,

  /* Limit workers to avoid race conditions */
  workers: 1,

  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ['html', { open: 'never', outputFolder: 'playwright-report' }],
    ['list'],
    ['json', { outputFile: 'test-results/results.json' }],
  ],

  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.BASE_URL || 'http://localhost:5000',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'retain-on-failure',

    /* Screenshot on failure */
    screenshot: 'only-on-failure',

    /* Video on failure */
    video: 'retain-on-failure',

    /* Test timeout */
    actionTimeout: 10000,
    navigationTimeout: 30000,
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
