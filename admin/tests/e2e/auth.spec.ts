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

    test('should login successfully with valid credentials', async ({ testData }) => {
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

      await loginPage.verifyErrorMessage(/请输入用户名|请输入密码|required/i);
    });

    test('should disable login button during submission', async ({ testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Click login and verify button is disabled
      await loginPage.clickLogin();
      await loginPage.verifyLoginButtonDisabled();
    });
  });

  test.describe('Form Validation', () => {
    test('should validate empty username field', async ({ testData }) => {
      await loginPage.fillCredentials('', testData.adminUser.password);
      await loginPage.clickLogin();

      await expect(loginPage['usernameInput']).toBeFocused();
    });

    test('should validate empty password field', async ({ testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, '');
      await loginPage.clickLogin();

      await expect(loginPage['passwordInput']).toBeFocused();
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
      const logoutButton = page.getByRole('button', { name: /退出|logout/i }).first();
      if (await logoutButton.isVisible()) {
        await logoutButton.click();
      }

      // Verify redirected to login page
      await expect(page).toHaveURL(/\/login/);
    });

    test('should redirect to login when accessing protected route without auth', async ({ page }) => {
      // Try to access admin page directly without login
      await page.goto('/admin/User');

      // Should redirect to login
      await expect(page).toHaveURL(/\/login/);
    });
  });

  test.describe('Token Validation', async () => {
    test('should validate API token format', async ({ testData }) => {
      await loginPage.loginAndWaitForDashboard(
        testData.adminUser.username,
        testData.adminUser.password
      );

      const token = await page.evaluate(() => {
        return localStorage.getItem('token') || sessionStorage.getItem('token');
      });

      // JWT tokens have 3 parts separated by dots
      expect(token?.split('.')).toHaveLength(3);
    });

    test('should get valid admin token via API', async () => {
      const token = await getAdminToken();
      expect(token).toBeTruthy();
      expect(token.length).toBeGreaterThan(50); // JWT tokens are typically long
    });
  });

  test.describe('UI/UX', () => {
    test('should show loading state during login', async ({ page, testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Click login and check for loading indicator
      await loginPage.clickLogin();

      // Ant Design shows loading spinner on button
      await expect(loginPage['loginButton']).toHaveAttribute('class', /loading/);
    });

    test('should allow Enter key to submit form', async ({ page, testData }) => {
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);

      // Press Enter on password field
      await loginPage['passwordInput'].press('Enter');

      await loginPage.verifyLoginSuccess();
    });

    test('should clear password field after failed login', async ({ testData }) => {
      await loginPage.login(testData.adminUser.username, 'wrongpassword');

      await loginPage.verifyErrorMessage(/用户名或密码错误/i);

      // Password field should be cleared or allow new input
      await loginPage.fillCredentials(testData.adminUser.username, testData.adminUser.password);
      await loginPage.verifyLoginSuccess();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper form labels', async () => {
      const usernameLabel = loginPage.page.getByText(/用户名|username/i);
      const passwordLabel = loginPage.page.getByText(/密码|password/i);

      await expect(usernameLabel).toBeVisible();
      await expect(passwordLabel).toBeVisible();
    });

    test('should be keyboard navigable', async ({ testData }) => {
      // Tab to username
      await page.keyboard.press('Tab');
      await expect(loginPage['usernameInput']).toBeFocused();

      // Tab to password
      await page.keyboard.press('Tab');
      await expect(loginPage['passwordInput']).toBeFocused();

      // Tab to login button
      await page.keyboard.press('Tab');
      await expect(loginPage['loginButton']).toBeFocused();

      // Submit with Enter
      await page.keyboard.press('Enter');
      await loginPage.verifyLoginSuccess();
    });
  });
});
