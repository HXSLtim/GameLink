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
  private readonly successMessage: Locator;

  constructor(private page: Page) {
    this.usernameInput = this.page.getByPlaceholder(/管理员账号\/邮箱|请输入用户名/);
    this.passwordInput = this.page.getByPlaceholder(/密码|请输入密码/);
    this.loginButton = this.page.getByRole('button', { name: /登录|login/i });
    this.errorMessage = this.page.locator('.ant-message-error, .error-message');
    this.successMessage = this.page.locator('.ant-message-success');
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
