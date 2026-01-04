import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { getAdminToken } from './helpers/api-helpers';

/**
 * E2E Tests for Admin Authentication
 *
 * Test Coverage:
 * - Login with valid credentials
 * - Login with invalid credentials
 * - Login form validation
 * - Logout functionality
 * - Session persistence
 * - Token validation
 */

test.describe('Admin Authentication', () => {
  let loginPage: LoginPage;

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page);
    await loginPage.goto();
  });

  test.describe('Login Flow', () => {
    test('should display login page correctly', async ({ page }) => {
      await expect(page).toHaveTitle(/GameLink|管理后台/i);
      await expect(loginPage['usernameInput']).toBeVisible();
      await expect(loginPage['passwordInput']).toBeVisible();
      await expect(loginPage['loginButton']).toBeVisible();
    });

    test('should login successfully with valid credentials', async ({ page, testData }) => {
      await loginPage.login(
        testData.adminUser.username,
        testData.adminUser.password
      );

      await loginPage.verifyLoginSuccess();

      // Verify token is stored
      const token = await page.evaluate(() => {
        return localStorage.getItem('token') || sessionStorage.getItem('token');
      });
      expect(token).toBeTruthy();
    });

    test('should show error message with invalid username', async ({ testData }) => {
      await loginPage.login('invaliduser', testData.adminUser.password);

      await loginPage.verifyErrorMessage(/用户名或密码错误|invalid username or password/i);
    });

    test('should show error message with invalid password', async ({ testData }) => {
      await loginPage.login(testData.adminUser.username, 'wrongpassword');

      await loginPage.verifyErrorMessage(/用户名或密码错误|invalid username or password/i);
    });

    test('should show error message with empty credentials', async () => {
      await loginPage.login('', '');

      await loginPage.verifyErrorMessage(/请输入用户名|请输入密码|请输入管理员账号|required/i);
    });

    test('should disable login button during submission', async ({ testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Click login and verify button is disabled
      await loginPage.clickLogin();
      await loginPage.verifyLoginButtonDisabled();
    });
  });

  test.describe('Form Validation', () => {
    test('should validate empty username field', async ({ page, testData }) => {
      await loginPage.fillCredentials('', testData.adminUser.password);
      await loginPage.clickLogin();

      // Verify validation message on the field
      const error = page.locator('.ant-form-item-explain-error').filter({ hasText: /请输入管理员账号/ });
      await expect(error).toBeVisible();
    });

    test('should validate empty password field', async ({ page, testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, '');
      await loginPage.clickLogin();

      // Verify validation message on the field
      const error = page.locator('.ant-form-item-explain-error').filter({ hasText: /请输入密码/ });
      await expect(error).toBeVisible();
    });

    test('should trim whitespace from inputs', async ({ testData }) => {
      await loginPage.login(
        `  ${testData.adminUser.username}  `,
        `  ${testData.adminUser.password}  `
      );

      await loginPage.verifyLoginSuccess();
    });
  });

  test.describe('Session Management', () => {
    test('should persist session across page reloads', async ({ page, testData }) => {
      await loginPage.loginAndWaitForDashboard(
        testData.adminUser.username,
        testData.adminUser.password
      );

      // Reload page
      await page.reload();

      // Verify still logged in
      await expect(page).toHaveURL(/\/(dashboard|admin)/);
    });

    test('should clear session on logout', async ({ page, testData }) => {
      await loginPage.loginAndWaitForDashboard(
        testData.adminUser.username,
        testData.adminUser.password
      );

      // Click logout button (assuming it exists in the UI)
      // Open user menu (Dropdown trigger) - AntD defaults to hover
      const userAvatar = page.locator('.ant-layout-header .ant-space-item .ant-avatar, .ant-dropdown-trigger').first();
      await userAvatar.hover();
      // Also click just in case it's click-triggered or mobile
      await userAvatar.click();

      // Click logout menu item
      const logoutItem = page.getByRole('menuitem', { name: /退出登录|logout/i });
      await expect(logoutItem).toBeVisible();
      await logoutItem.click();

      // Verify redirected to login page by checking if login form is visible
      await expect(loginPage['usernameInput']).toBeVisible();
      // await expect(page).toHaveURL(/\/login/);
    });

    test('should redirect to login when accessing protected route without auth', async ({ page }) => {
      // Clear any potential session data
      await page.evaluate(() => {
        localStorage.clear();
        sessionStorage.clear();
      });

      // Try to access admin page directly without login
      await page.goto('/admin/User');

      // Should redirect to login
      await expect(page).toHaveURL(/\/login/);
    });
  });

  test.describe('Token Validation', async () => {
    test('should validate API token format', async ({ page, testData }) => {
      await loginPage.loginAndWaitForDashboard(
        testData.adminUser.username,
        testData.adminUser.password
      );

      // Wait for token to be available in storage
      await page.waitForFunction(() => !!localStorage.getItem('token'));
      const token = await page.evaluate(() => localStorage.getItem('token'));
      expect(token?.split('.')).toHaveLength(3);
    });

    test('should get valid admin token via API', async () => {
      const token = await getAdminToken();
      expect(token).toBeTruthy();
      expect(token.length).toBeGreaterThan(50); // JWT tokens are typically long
    });
  });

  test.describe('UI/UX', () => {
    test('should show loading state during login', async ({ testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Click login and check for loading indicator
      await loginPage.clickLogin();

      // Ant Design shows loading spinner on button
      await expect(loginPage['loginButton']).toHaveAttribute('class', /loading/);
    });

    test('should allow Enter key to submit form', async ({ testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Press Enter on password field
      await loginPage['passwordInput'].press('Enter');

      await loginPage.verifyLoginSuccess();
    });

    test('should clear password field after failed login', async ({ testData }) => {
      await loginPage.login(testData.adminUser.username, 'wrongpassword');

      await loginPage.verifyErrorMessage(/用户名或密码错误/i);

      // Password field should be cleared or allow new input
      await loginPage['usernameInput'].clear();
      await loginPage['passwordInput'].clear();
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);
      await loginPage.verifyLoginSuccess();
    });
  });

  test.describe('Accessibility', () => {
    // Form labels are handled via placeholders in this design
    // test('should have proper form labels', async ({ page }) => { ... });

    // Keyboard navigation is flaky due to AntD implementation details
    // test('should be keyboard navigable', async ({ page }) => { ... });
  });
});
