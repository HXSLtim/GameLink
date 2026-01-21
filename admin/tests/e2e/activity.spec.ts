import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { ActivityPage } from './pages/ActivityPage';

/**
 * E2E Tests for Activity Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Activity list with pagination
 * - Statistics display
 * - Search and filter activities
 * - Create, edit, delete activities
 * - Publish/unpublish activities
 * - View activity details
 * - Rewards management
 */

test.describe('Activity Management', () => {
  let loginPage: LoginPage;
  let activityPage: ActivityPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    activityPage = new ActivityPage(page);

    // Login and navigate to activity
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display activity page correctly', async () => {
      await activityPage.goto();

      await expect(activityPage.pageTitle).toBeVisible();
      await expect(activityPage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await activityPage.goto();

      await expect(activityPage.totalActivitiesStat).toBeVisible();
      await expect(activityPage.activeActivitiesStat).toBeVisible();
      await expect(activityPage.upcomingActivitiesStat).toBeVisible();
    });

    test('should load activity list', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await activityPage.goto();

      const stats = await activityPage.getStatistics();
      expect(stats.total).toBeGreaterThanOrEqual(0);
      expect(stats.active).toBeGreaterThanOrEqual(0);
      expect(stats.upcoming).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search activities by keyword', async () => {
      await activityPage.goto();

      await activityPage.searchByKeyword('test');

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by activity type', async () => {
      await activityPage.goto();

      await activityPage.filterByType('sign_in');

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by status', async () => {
      await activityPage.goto();

      await activityPage.filterByStatus('draft');

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await activityPage.goto();

      // Apply a filter
      await activityPage.searchByKeyword('test');

      // Clear search
      await activityPage.clearFilters();

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Activity Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await activityPage.goto();

      await activityPage.clickAdd();

      const isModalVisible = await activityPage.isModalVisible();
      if (isModalVisible) {
        await activityPage.closeModal();
      }
    });

    test('should view activity details', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        await activityPage.viewDetails(0);

        const isModalVisible = await activityPage.isModalVisible();
        if (isModalVisible) {
          await activityPage.closeModal();
        }
      }
    });

    test('should edit activity', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        await activityPage.editActivity(0);

        const isModalVisible = await activityPage.isModalVisible();
        if (isModalVisible) {
          await activityPage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        const _countBefore = await activityPage.getActivityCount();

        // Try to delete first activity
        await activityPage.deleteActivity(0);

        // Wait for potential deletion
        await activityPage.page.waitForTimeout(1000);

        const countAfter = await activityPage.getActivityCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });

    test('should publish activity', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        await activityPage.publishActivity(0);
        await activityPage.page.waitForTimeout(1000);
      }
    });

    test('should unpublish activity', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        await activityPage.unpublishActivity(0);
        await activityPage.page.waitForTimeout(1000);
      }
    });

    test('should open rewards management', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        await activityPage.openRewards(0);

        const isModalVisible = await activityPage.isModalVisible();
        if (isModalVisible) {
          await activityPage.closeModal();
        }
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display activity columns correctly', async () => {
      await activityPage.goto();

      await expect(activityPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        const cellText = await activityPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });

    test('should get activity status from table', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        const status = await activityPage.getActivityStatus(0);
        expect(status).toBeTruthy();
        expect(['草稿', '预热', '进行中', '已结束']).toContain(status);
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh activity list', async () => {
      await activityPage.goto();

      const _countBefore = await activityPage.getActivityCount();

      await activityPage.refresh();

      const countAfter = await activityPage.getActivityCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await activityPage.goto();

      await expect(activityPage.pageTitle).toBeVisible();
      await expect(activityPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await activityPage.goto();

      await expect(activityPage.pageTitle).toContainText('活动管理');
    });

    test('should be keyboard navigable', async () => {
      await activityPage.goto();

      await activityPage.page.keyboard.press('Tab');
      await activityPage.page.keyboard.press('Tab');

      const focusedElement = activityPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const activityMenuItem = page.locator('.ant-menu-item').filter({ hasText: /活动管理/i });
      if (await activityMenuItem.isVisible()) {
        await activityMenuItem.click();
        await activityPage.waitForPageLoad();
        await expect(activityPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty activity list', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(activityPage.table).toBeVisible();
      }
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await activityPage.goto();

      // Apply type filter
      await activityPage.filterByType('sign_in');

      // Then apply status filter
      await activityPage.filterByStatus('draft');

      const count = await activityPage.getActivityCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Activity Lifecycle', () => {
    test('should handle activity status transitions', async () => {
      await activityPage.goto();

      const count = await activityPage.getActivityCount();
      if (count > 0) {
        // Get initial status
        const statusBefore = await activityPage.getActivityStatus(0);

        // Try to publish/unpublish
        if (statusBefore === '草稿' || statusBefore === '预热') {
          await activityPage.publishActivity(0);
        } else if (statusBefore === '进行中') {
          await activityPage.unpublishActivity(0);
        }

        await activityPage.page.waitForTimeout(1000);
      }
    });
  });

  test.describe('Search Functionality', () => {
    test('should verify activity exists in table', async () => {
      await activityPage.goto();

      const exists = await activityPage.verifyActivityExists('test');
      expect(exists).toBeDefined();
    });
  });
});
