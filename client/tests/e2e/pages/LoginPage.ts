/**
 * Login Page Object Model
 *
 * Selectors based on actual LoginPage.tsx implementation:
 * - Input id="username" with placeholder from i18n
 * - Input id="password" with placeholder="••••••••"
 * - Button with type="submit"
 */
export class LoginPage {
  usernameInput: any;
  passwordInput: any;
  loginButton: any;

  constructor(private page: any) {
    // Use id selectors since the actual implementation uses id attributes
    this.usernameInput = page.locator('#username');
    this.passwordInput = page.locator('#password');
    // Login button - use text content "Sign In" or "登录" to be specific
    this.loginButton = page.getByRole('button', { name: /sign in|登录|sign·in/i })
      .filter({ hasText: /^(sign in|登录)$/i });
  }

  /** Form error messages - based on actual implementation using custom error div */
  get errorMessage() {
    // Error is shown in div with classes: text-destructive bg-destructive/10
    return this.page.locator('div.text-destructive').or(
      this.page.getByText(/用户名或密码错误|invalid username or password|登录失败|authentication failed/i)
    );
  }
  get usernameError() {
    return this.page.getByText(/用户名|username/i).locator('..').locator('.text-destructive');
  }
  get passwordError() {
    return this.page.getByText(/密码|password/i).locator('..').locator('.text-destructive');
  }

  /** Loading indicator */
  get loadingSpinner() {
    return this.page.locator('.animate-spin').or(
      this.page.getByRole('button', { name: /sign in|登录/i }).locator('svg[class*="spin"]')
    );
  }

  /** Navigate to login page */
  async goto() {
    await this.page.goto('/login');
  }

  /** Fill in login credentials */
  async fillCredentials(username: string, password: string) {
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
  }

  /** Submit login form */
  async clickLogin() {
    await this.loginButton.click();
  }

  /** Perform login */
  async login(username: string, password: string) {
    await this.fillCredentials(username, password);
    await this.clickLogin();
  }

  /** Verify login success - redirected to home or profile */
  async verifyLoginSuccess() {
    // Wait for navigation and check we're not on login page anymore
    await this.page.waitForURL(/\/(home|profile|players)?$/, { timeout: 10000 });
  }

  /** Verify error message is shown */
  async verifyErrorMessage(_pattern: RegExp) {
    const { expect } = await import('@playwright/test');
    await expect(this.errorMessage).toBeVisible();
  }
}
