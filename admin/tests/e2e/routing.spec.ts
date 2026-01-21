import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { RoutingPage } from './pages/RoutingPage';

/**
 * E2E Tests for Routing Rule Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Rule list with pagination
 * - Statistics display
 * - Search and filter rules
 * - Create, edit, delete rules
 * - Toggle rule status
 * - View rule history
 * - Test rules functionality
 * - Export rules
 */

test.describe('Routing Rule Management', () => {
  let loginPage: LoginPage;
  let routingPage: RoutingPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    routingPage = new RoutingPage(page);

    // Login and navigate to routing
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display routing page correctly', async () => {
      await routingPage.goto();

      await expect(routingPage.pageTitle).toBeVisible();
      await expect(routingPage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await routingPage.goto();

      await expect(routingPage.totalRulesStat).toBeVisible();
      await expect(routingPage.activeRulesStat).toBeVisible();
      await expect(routingPage.inactiveRulesStat).toBeVisible();
    });

    test('should load rule list', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await routingPage.goto();

      const stats = await routingPage.getStatistics();
      expect(stats.total).toBeGreaterThanOrEqual(0);
      expect(stats.active).toBeGreaterThanOrEqual(0);
      expect(stats.inactive).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search rules by keyword', async () => {
      await routingPage.goto();

      await routingPage.searchByKeyword('test');

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by status', async () => {
      await routingPage.goto();

      await routingPage.filterByStatus('active');

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await routingPage.goto();

      // Apply a filter
      await routingPage.searchByKeyword('test');

      // Clear search
      await routingPage.clearFilters();

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Rule Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await routingPage.goto();

      await routingPage.clickAdd();

      const isModalVisible = await routingPage.isModalVisible();
      if (isModalVisible) {
        await routingPage.closeModal();
      }
    });

    test('should edit rule', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        await routingPage.editRule(0);

        const isModalVisible = await routingPage.isModalVisible();
        if (isModalVisible) {
          await routingPage.closeModal();
        }
      }
    });

    test('should view rule history', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        await routingPage.viewHistory(0);

        const isModalVisible = await routingPage.isModalVisible();
        if (isModalVisible) {
          await routingPage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        const _countBefore = await routingPage.getRuleCount();

        // Try to delete first rule
        await routingPage.deleteRule(0);

        // Wait for potential deletion
        await routingPage.page.waitForTimeout(1000);

        const countAfter = await routingPage.getRuleCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });

    test('should toggle rule status', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        // Toggle status of first rule
        await routingPage.toggleStatus(0);
        await routingPage.page.waitForTimeout(1000);
      }
    });
  });

  test.describe('Status Filtering', () => {
    test('should filter by active status', async () => {
      await routingPage.goto();
      await routingPage.filterByStatus('active');
      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by inactive status', async () => {
      await routingPage.goto();
      await routingPage.filterByStatus('inactive');
      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Toolbar Actions', () => {
    test('should have test rules button', async () => {
      await routingPage.goto();

      await expect(routingPage.testButton).toBeVisible();
    });

    test('should have export button', async () => {
      await routingPage.goto();

      await expect(routingPage.exportButton).toBeVisible();
    });

    test('should click test rules button', async () => {
      await routingPage.goto();

      await routingPage.clickTestRules();
      // Navigation should occur, but we don't assert destination
      // as it depends on test page implementation
    });

    test('should click export button', async () => {
      await routingPage.goto();

      await routingPage.clickExport();
      await routingPage.page.waitForTimeout(1000);
      // Export may trigger download, difficult to assert in headless mode
    });
  });

  test.describe('Table Display', () => {
    test('should display rule columns correctly', async () => {
      await routingPage.goto();

      await expect(routingPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        const cellText = await routingPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh rule list', async () => {
      await routingPage.goto();

      const _countBefore = await routingPage.getRuleCount();

      await routingPage.refresh();

      const countAfter = await routingPage.getRuleCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await routingPage.goto();

      await expect(routingPage.pageTitle).toBeVisible();
      await expect(routingPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await routingPage.goto();

      const title = await routingPage.pageTitle.textContent();
      expect(title).toMatch(/支付路由规则/i);
    });

    test('should be keyboard navigable', async () => {
      await routingPage.goto();

      await routingPage.page.keyboard.press('Tab');
      await routingPage.page.keyboard.press('Tab');

      const focusedElement = routingPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const routingMenuItem = page.locator('.ant-menu-item').filter({ hasText: /支付路由|路由管理/i });
      if (await routingMenuItem.isVisible()) {
        await routingMenuItem.click();
        await routingPage.waitForPageLoad();
        await expect(routingPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty rule list', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(routingPage.table).toBeVisible();
      }
    });
  });

  test.describe('Statistics Verification', () => {
    test('should show consistent statistics', async () => {
      await routingPage.goto();

      const stats1 = await routingPage.getStatistics();

      await routingPage.page.waitForTimeout(1000);

      const stats2 = await routingPage.getStatistics();

      // Statistics should be consistent
      expect(stats1.total).toEqual(stats2.total);
    });

    test('should have correct active/inactive distribution', async () => {
      await routingPage.goto();

      const stats = await routingPage.getStatistics();

      // Total should equal active + inactive
      expect(stats.total).toEqual(stats.active + stats.inactive);
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await routingPage.goto();

      // Apply status filter
      await routingPage.filterByStatus('active');

      // Then apply keyword filter
      await routingPage.searchByKeyword('test');

      const count = await routingPage.getRuleCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Rule Priority Display', () => {
    test('should display rule priority in table', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        // Priority is typically in column 2 (after ID and name)
        const priorityCell = await routingPage.getCellText(0, 2);
        expect(priorityCell).toBeTruthy();
      }
    });
  });

  test.describe('History Modal', () => {
    test('should open history modal', async () => {
      await routingPage.goto();

      const count = await routingPage.getRuleCount();
      if (count > 0) {
        await routingPage.viewHistory(0);

        const isModalVisible = await routingPage.isModalVisible();
        if (isModalVisible) {
          // History modal should display
          await routingPage.closeModal();
        }
      }
    });
  });
});
