import { Page, expect, Locator } from '@playwright/test';

/**
 * Page Object Model for Login Page
 * Encapsulates all login page interactions and assertions
 */
export class LoginPage {
  private readonly usernameInput: Locator;
  private readonly passwordInput: Locator;
  private readonly loginButton: Locator;
  private readonly errorMessage: Locator;

  constructor(private page: Page) {
    this.usernameInput = this.page.getByPlaceholder(/管理员账号\/邮箱|请输入用户名/);
    this.passwordInput = this.page.getByPlaceholder(/密码|请输入密码/);
    this.loginButton = this.page.getByRole('button', { name: /登录|login/i });
    this.errorMessage = this.page.locator('.ant-message-error, .error-message');
  }

  /**
   * Navigate to login page
   */
  async goto() {
    await this.page.goto('/login');
    await expect(this.page).toHaveTitle(/GameLink|管理后台/i);
  }

  /**
   * Fill in login form
   */
  async fillCredentials(username: string, password: string) {
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
  }

  /**
   * Click login button
   */
  async clickLogin() {
    await this.loginButton.click();
  }

  /**
   * Perform complete login flow
   */
  async login(username: string, password: string) {
    await this.fillCredentials(username, password);
    await this.clickLogin();
  }

  /**
   * Verify login is successful by checking URL change
   */
  async verifyLoginSuccess() {
    await expect(this.page).toHaveURL(/(\/dashboard|\/admin$|\/admin\/)/);
    await expect(this.page).not.toHaveURL(/login/);
  }

  /**
   * Verify error message is displayed
   */
  async verifyErrorMessage(message: string | RegExp) {
    const errorLocator = this.page.locator('.ant-message-error, .error-message, .ant-form-item-explain-error');
    await expect(errorLocator.first()).toBeVisible();
    await expect(errorLocator.first()).toContainText(message);
  }

  /**
   * Verify login button is disabled while loading
   */
  async verifyLoginButtonDisabled() {
    // AntD button uses loading class
    await expect(this.loginButton).toHaveClass(/ant-btn-loading/);
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await expect(this.usernameInput).toBeVisible();
  }

  /**
   * Get current URL
   */
  getCurrentUrl(): string {
    return this.page.url();
  }

  /**
   * Check if already logged in (redirected away from login page)
   */
  async isLoggedIn(): Promise<boolean> {
    await this.page.waitForLoadState('networkidle');
    const url = this.page.url();
    return !url.includes('/login');
  }

  /**
   * Clear form fields
   */
  async clearForm() {
    await this.usernameInput.clear();
    await this.passwordInput.clear();
  }

  /**
   * Verify form validation (empty fields)
   */
  async verifyFormValidation() {
    await this.loginButton.click();
    await expect(this.errorMessage).toBeVisible();
  }

  /**
   * Login and wait for dashboard with retry logic
   */
  async loginAndWaitForDashboard(username: string, password: string, maxRetries = 3) {
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        await this.loginAttempt(username, password);
        return; // Success
      } catch (error) {
        if (attempt < maxRetries) {

          await this.page.waitForTimeout(attempt * 2000);
          // Go back to login page for retry
          await this.goto();
        } else {
          // Take screenshot on final failure
          await this.page.screenshot({ path: `test-results/login-failed-attempt-${attempt}.png` });
          throw error;
        }
      }
    }
  }

  /**
   * Single login attempt
   */
  private async loginAttempt(username: string, password: string) {
    await this.login(username, password);

    // Wait for URL to change from login page
    try {
      await this.page.waitForURL(/\/(dashboard|admin)/, { timeout: 15000 });
    } catch {
      const currentUrl = this.page.url();
      if (currentUrl.includes('/login')) {
        throw new Error('Login failed - still on login page');
      }
    }

    // Wait for network to settle
    await this.page.waitForLoadState('networkidle', { timeout: 15000 });

    // Wait for the admin layout to be fully loaded
    await this.waitForDashboardLayout();
  }

  /**
   * Wait for dashboard layout to be ready
   */
  private async waitForDashboardLayout() {
    // Wait for loading spinner to disappear
    try {
      const spinner = this.page.locator('.ant-spin-spinning');
      if (await spinner.count() > 0) {
        await spinner.waitFor({ state: 'hidden', timeout: 15000 });
      }
    } catch {
      // Spinner might not exist, continue
    }

    // Wait for the layout header to be visible
    const header = this.page.locator('.ant-layout-header');
    await header.waitFor({ state: 'visible', timeout: 15000 });

    // Check if system needs initialization
    const initButton = this.page.getByRole('button', { name: /初始化系统/i });
    const hasInitButton = await initButton.isVisible().catch(() => false);

    if (hasInitButton) {
      await initButton.click();
      await this.page.waitForTimeout(3000);
      await this.page.reload();
      await this.page.waitForLoadState('networkidle');
      await header.waitFor({ state: 'visible', timeout: 15000 });
    }

    // Wait for sidebar menu to appear (indicates menus are loaded)
    try {
      const menu = this.page.locator('.ant-menu');
      await menu.waitFor({ state: 'visible', timeout: 10000 });
    } catch {
      // Menu might take longer, but header is visible so we can proceed
    }

    // Small stability wait
    await this.page.waitForTimeout(500);
  }
}
