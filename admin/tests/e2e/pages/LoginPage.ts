import { Page, expect } from '@playwright/test';

/**
 * Page Object Model for Login Page
 * Encapsulates all login page interactions and assertions
 */
export class LoginPage {
  constructor(private page: Page) {}

  // Element locators
  private readonly usernameInput = this.page.getByPlaceholder('请输入用户名');
  private readonly passwordInput = this.page.getByPlaceholder('请输入密码');
  private readonly loginButton = this.page.getByRole('button', { name: /登录|login/i });
  private readonly errorMessage = this.page.locator('.ant-message-error, .error-message');
  private readonly successMessage = this.page.locator('.ant-message-success');

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
    await expect(this.page).toHaveURL(/\/(dashboard|admin)/);
  }

  /**
   * Verify error message is displayed
   */
  async verifyErrorMessage(message: string) {
    await expect(this.errorMessage).toBeVisible();
    await expect(this.errorMessage).toContainText(message);
  }

  /**
   * Verify login button is disabled while loading
   */
  async verifyLoginButtonDisabled() {
    await expect(this.loginButton).toBeDisabled();
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
   * Login and wait for dashboard
   */
  async loginAndWaitForDashboard(username: string, password: string) {
    await this.login(username, password);
    await this.page.waitForURL(/\/(dashboard|admin)/, { timeout: 10000 });
    await this.page.waitForLoadState('networkidle');
  }
}
