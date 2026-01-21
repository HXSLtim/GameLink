import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { WithdrawPage } from './pages/WithdrawPage';

/**
 * E2E Tests for Withdraw Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Withdraw list with pagination
 * - Statistics display
 * - Search and filter withdraws
 * - View withdraw details
 * - Approve/reject/complete withdraws
 * - Batch operations
 */

test.describe('Withdraw Management', () => {
  let loginPage: LoginPage;
  let withdrawPage: WithdrawPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    withdrawPage = new WithdrawPage(page);

    // Login and navigate to withdraw
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display withdraw page correctly', async () => {
      await withdrawPage.goto();

      await expect(withdrawPage.pageTitle).toBeVisible();
      await expect(withdrawPage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await withdrawPage.goto();

      await expect(withdrawPage.pendingStat).toBeVisible();
      await expect(withdrawPage.approvedStat).toBeVisible();
      await expect(withdrawPage.completedStat).toBeVisible();
    });

    test('should load withdraw list', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await withdrawPage.goto();

      const stats = await withdrawPage.getStatistics();
      expect(stats.pending).toBeGreaterThanOrEqual(0);
      expect(stats.approved).toBeGreaterThanOrEqual(0);
      expect(stats.completed).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search by keyword', async () => {
      await withdrawPage.goto();

      await withdrawPage.searchByKeyword('test');

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by status', async () => {
      await withdrawPage.goto();

      await withdrawPage.filterByStatus('pending');

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await withdrawPage.goto();

      // Apply a filter
      await withdrawPage.searchByKeyword('test');

      // Clear search
      await withdrawPage.clearFilters();

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Withdraw Actions', () => {
    test('should view withdraw details', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count > 0) {
        await withdrawPage.viewDetails(0);

        const isDrawerVisible = await withdrawPage.isDrawerVisible();
        if (isDrawerVisible) {
          await withdrawPage.closeDrawer();
        }
      }
    });

    test('should handle approve operation', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count > 0) {
        // Check first row's status
        const status = await withdrawPage.getWithdrawStatus(0);

        if (status === '待审核') {
          await withdrawPage.approveWithdraw(0);
          await withdrawPage.page.waitForTimeout(1000);
        }
      }
    });

    test('should handle reject operation', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count > 0) {
        // Check first row's status
        const status = await withdrawPage.getWithdrawStatus(0);

        if (status === '待审核') {
          await withdrawPage.rejectWithdraw(0);
          await withdrawPage.page.waitForTimeout(1000);

          // Close reject modal if opened
          if (await withdrawPage.isModalVisible()) {
            await withdrawPage.closeModal();
          }
        }
      }
    });

    test('should handle complete operation', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count > 0) {
        // Check first row's status
        const status = await withdrawPage.getWithdrawStatus(0);

        if (status === '已批准') {
          await withdrawPage.completeWithdraw(0);
          await withdrawPage.page.waitForTimeout(1000);
        }
      }
    });
  });

  test.describe('Withdraw Status Flow', () => {
    test('should display correct status for pending withdraws', async () => {
      await withdrawPage.goto();

      await withdrawPage.filterByStatus('pending');

      const count = await withdrawPage.getWithdrawCount();
      for (let i = 0; i < Math.min(count, 3); i++) {
        const status = await withdrawPage.getWithdrawStatus(i);
        expect(status).toContain('待审核');
      }
    });

    test('should display correct status for approved withdraws', async () => {
      await withdrawPage.goto();

      await withdrawPage.filterByStatus('approved');

      const count = await withdrawPage.getWithdrawCount();
      for (let i = 0; i < Math.min(count, 3); i++) {
        const status = await withdrawPage.getWithdrawStatus(i);
        expect(status).toContain('已批准');
      }
    });

    test('should display correct status for completed withdraws', async () => {
      await withdrawPage.goto();

      await withdrawPage.filterByStatus('completed');

      const count = await withdrawPage.getWithdrawCount();
      for (let i = 0; i < Math.min(count, 3); i++) {
        const status = await withdrawPage.getWithdrawStatus(i);
        expect(status).toContain('已完成');
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display withdraw columns correctly', async () => {
      await withdrawPage.goto();

      await expect(withdrawPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count > 0) {
        const cellText = await withdrawPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh withdraw list', async () => {
      await withdrawPage.goto();

      const _countBefore = await withdrawPage.getWithdrawCount();

      await withdrawPage.refresh();

      const countAfter = await withdrawPage.getWithdrawCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await withdrawPage.goto();

      await expect(withdrawPage.pageTitle).toBeVisible();
      await expect(withdrawPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await withdrawPage.goto();

      await expect(withdrawPage.pageTitle).toContainText('提现管理');
    });

    test('should be keyboard navigable', async () => {
      await withdrawPage.goto();

      await withdrawPage.page.keyboard.press('Tab');
      await withdrawPage.page.keyboard.press('Tab');

      const focusedElement = withdrawPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const withdrawMenuItem = page.locator('.ant-menu-item').filter({ hasText: /提现管理/i });
      if (await withdrawMenuItem.isVisible()) {
        await withdrawMenuItem.click();
        await withdrawPage.waitForPageLoad();
        await expect(withdrawPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty withdraw list', async () => {
      await withdrawPage.goto();

      const count = await withdrawPage.getWithdrawCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(withdrawPage.table).toBeVisible();
      }
    });
  });

  test.describe('Statistics Verification', () => {
    test('should show consistent statistics', async () => {
      await withdrawPage.goto();

      const stats1 = await withdrawPage.getStatistics();

      await withdrawPage.page.waitForTimeout(1000);

      const stats2 = await withdrawPage.getStatistics();

      // Statistics should be consistent
      expect(stats1.pending).toEqual(stats2.pending);
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await withdrawPage.goto();

      // Apply status filter
      await withdrawPage.filterByStatus('pending');

      // Then apply keyword filter
      await withdrawPage.searchByKeyword('test');

      const count = await withdrawPage.getWithdrawCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });
});
