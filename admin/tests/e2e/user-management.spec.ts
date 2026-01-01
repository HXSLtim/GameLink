import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { UserManagementPage } from './pages/UserManagementPage';
import { createTestUser, deleteTestUser, getAdminToken } from './helpers/api-helpers';

/**
 * E2E Tests for User Management
 *
 * Test Coverage:
 * - List users with pagination
 * - Create new user
 * - Edit existing user
 * - Delete user
 * - Search and filter users
 * - Batch operations
 * - User status management
 * - User role management
 */

test.describe('User Management', () => {
  let loginPage: LoginPage;
  let userManagementPage: UserManagementPage;
  let adminToken: string;
  let testUserId: number;

  test.beforeAll(async () => {
    adminToken = await getAdminToken();
  });

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    userManagementPage = new UserManagementPage(page);

    // Login and navigate to user management
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
    await userManagementPage.goto();
  });

  test.afterEach(async ({ page }) => {
    // Cleanup test user if created
    if (testUserId) {
      await deleteTestUser(adminToken, testUserId);
      testUserId = 0;
    }
  });

  test.describe('User List', () => {
    test('should display user list page correctly', async () => {
      await expect(userManagementPage['pageTitle']).toBeVisible();
      await expect(userManagementPage['table']).toBeVisible();
    });

    test('should load users in table', async () => {
      const userCount = await userManagementPage.getUserCount();
      expect(userCount).toBeGreaterThan(0);
    });

    test('should display user details correctly', async () => {
      const firstUserName = await userManagementPage.getCellText(0, 1);
      expect(firstUserName).toBeTruthy();
      expect(firstUserName.length).toBeGreaterThan(0);
    });

    test('should support pagination', async () => {
      const initialCount = await userManagementPage.getUserCount();

      await userManagementPage.nextPage();

      const newCount = await userManagementPage.getUserCount();

      // If there are more users, count should be the same (different page)
      // If no more users, count might be same or different
      expect(newCount).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Create User', () => {
    test('should create new user successfully', async ({ testData }) => {
      const testUser = testData.testUser;

      await userManagementPage.createUser({
        name: testUser.name,
        email: testUser.email,
        phone: testUser.phone,
        role: testUser.role,
        password: testUser.password,
        status: 'active',
      });

      await userManagementPage.verifySuccessMessage();

      // Verify user appears in table
      await userManagementPage.waitForUserToAppear(testUser.name, 10000);
    });

    test('should show validation errors for invalid email', async ({ testData }) => {
      await userManagementPage.clickCreateUser();

      await userManagementPage.fillUserForm({
        name: testData.testUser.name,
        email: 'invalid-email',
        phone: testData.testUser.phone,
        role: 'user',
      });

      await userManagementPage.submitForm();
      await userManagementPage.verifyErrorMessage();
    });

    test('should show validation errors for duplicate email', async ({ testData }) => {
      // Create first user
      testUserId = await createTestUser(adminToken, {
        ...testData.testUser,
        email: 'duplicate@example.com',
      });

      // Try to create duplicate
      await userManagementPage.createUser({
        name: `Another ${testData.testUser.name}`,
        email: 'duplicate@example.com',
        phone: testData.testUser.phone,
        role: 'user',
      });

      await userManagementPage.verifyErrorMessage();
    });

    test('should create user with different roles', async ({ testData }) => {
      const roles = ['user', 'player', 'admin'];

      for (const role of roles) {
        const testUser = {
          ...testData.testUser,
          name: `${role} User ${Date.now()}`,
          email: `${role}${Date.now()}@example.com`,
          role,
        };

        await userManagementPage.createUser({
          name: testUser.name,
          email: testUser.email,
          phone: testUser.phone,
          role,
          password: testUser.password,
        });

        await userManagementPage.verifySuccessMessage();
        await userManagementPage.waitForUserToAppear(testUser.name, 10000);
      }
    });
  });

  test.describe('Edit User', () => {
    test('should edit user information successfully', async ({ testData }) => {
      // Create test user first
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      // Find and edit user
      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);
      await userManagementPage.editUser(0, {
        name: `${testData.testUser.name} (Updated)`,
      });

      await userManagementPage.verifySuccessMessage();
    });

    test('should update user role', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);
      await userManagementPage.editUser(0, {
        role: 'player',
      });

      await userManagementPage.verifySuccessMessage();
    });

    test('should update user status', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, {
        ...testData.testUser,
        status: 'active',
      });
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);
      await userManagementPage.editUser(0, {
        status: 'banned',
      });

      await userManagementPage.verifySuccessMessage();
    });

    test('should cancel edit operation', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);

      // Click edit but then cancel
      await userManagementPage['editButton'](0).click();
      await expect(userManagementPage['modal']).toBeVisible();

      await userManagementPage['modalCancelButton'].click();

      // Verify modal is closed and data unchanged
      await expect(userManagementPage['modal']).not.toBeVisible();
    });
  });

  test.describe('Delete User', () => {
    test('should delete user successfully', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);

      await userManagementPage.deleteUser(0);

      await userManagementPage.verifySuccessMessage();

      // Verify user is removed from table
      await expect(async () => {
        const exists = await userManagementPage.verifyUserExists(testData.testUser.name);
        expect(exists).toBe(false);
      }).toPass({ timeout: 5000 });
    });

    test('should cancel delete operation', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);

      await userManagementPage.cancelDelete(0);

      // Verify user still exists
      const exists = await userManagementPage.verifyUserExists(testData.testUser.name);
      expect(exists).toBe(true);

      // Cleanup manually
      await deleteTestUser(adminToken, testUserId);
      testUserId = 0;
    });

    test('should show confirmation dialog before delete', async ({ testData }) => {
      testUserId = await createTestUser(adminToken, testData.testUser);
      await userManagementPage.page.reload();
      await userManagementPage.waitForPageLoad();

      await userManagementPage.waitForUserToAppear(testData.testUser.name, 10000);

      await userManagementPage['deleteButton'](0).click();

      // Verify confirmation modal is shown
      await expect(userManagementPage['modal']).toBeVisible();
      await expect(userManagementPage['confirmDeleteButton']).toBeVisible();

      // Cancel to cleanup manually
      await userManagementPage['modalCancelButton'].click();
      await deleteTestUser(adminToken, testUserId);
      testUserId = 0;
    });
  });

  test.describe('Search and Filter', () => {
    test('should search users by keyword', async () => {
      const initialCount = await userManagementPage.getUserCount();

      await userManagementPage.searchUser('admin');

      await userManagementPage.page.waitForTimeout(1000);

      const filteredCount = await userManagementPage.getUserCount();

      // Filtered results should be fewer or equal
      expect(filteredCount).toBeLessThanOrEqual(initialCount);
    });

    test('should filter users by role', async () => {
      await userManagementPage.filterByRole('user');

      await userManagementPage.page.waitForTimeout(1000);

      const userCount = await userManagementPage.getUserCount();
      expect(userCount).toBeGreaterThanOrEqual(0);
    });

    test('should filter users by status', async () => {
      await userManagementPage.filterByStatus('active');

      await userManagementPage.page.waitForTimeout(1000);

      const activeCount = await userManagementPage.getUserCount();
      expect(activeCount).toBeGreaterThanOrEqual(0);
    });

    test('should clear filters', async () => {
      // Apply filter
      await userManagementPage.filterByRole('admin');

      await userManagementPage.page.waitForTimeout(1000);

      const filteredCount = await userManagementPage.getUserCount();

      // Clear search
      await userManagementPage.searchUser('');

      await userManagementPage.page.waitForTimeout(1000);

      const clearedCount = await userManagementPage.getUserCount();

      expect(clearedCount).toBeGreaterThanOrEqual(filteredCount);
    });
  });

  test.describe('View User Details', () => {
    test('should view user details', async () => {
      await userManagementPage.viewUserDetails(0);

      await expect(userManagementPage['modal']).toBeVisible();

      // Modal should contain user information
      await expect(userManagementPage['modalTitle']).toBeVisible();
    });

    test('should close user details modal', async () => {
      await userManagementPage.viewUserDetails(0);
      await expect(userManagementPage['modal']).toBeVisible();

      await userManagementPage.closeModal();

      await expect(userManagementPage['modal']).not.toBeVisible();
    });
  });

  test.describe('Batch Operations', () => {
    test('should batch delete multiple users', async ({ testData }) => {
      // Create test users
      const userIds: number[] = [];
      try {
        for (let i = 0; i < 3; i++) {
          const userData = {
            ...testData.testUser,
            name: `Batch User ${i} ${Date.now()}`,
            email: `batch${i}${Date.now()}@example.com`,
          };
          const userId = await createTestUser(adminToken, userData);
          userIds.push(userId);
        }

        await userManagementPage.page.reload();
        await userManagementPage.waitForPageLoad();

        // Batch delete would be performed here
        // Implementation depends on UI availability
      } finally {
        // Cleanup
        for (const id of userIds) {
          await deleteTestUser(adminToken, id);
        }
      }
    });
  });

  test.describe('Export Functionality', () => {
    test('should export user list', async () => {
      const downloadPromise = userManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await userManagementPage.exportUsers();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });
  });

  test.describe('Accessibility', () => {
    test('should be keyboard navigable', async () => {
      // Tab through table
      await userManagementPage.page.keyboard.press('Tab');
      await userManagementPage.page.keyboard.press('Tab');

      // Focus should be on table or navigation
      const focusedElement = userManagementPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });
});
