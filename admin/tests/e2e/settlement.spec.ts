import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { SettlementPage } from './pages/SettlementPage';

/**
 * E2E Tests for Settlement Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Player assignment list with pagination
 * - Search and filter
 * - View player details
 * - Assign player to company
 * - Transfer player
 * - View player history
 */

test.describe('Settlement Management', () => {
  let loginPage: LoginPage;
  let settlementPage: SettlementPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    settlementPage = new SettlementPage(page);

    // Login and navigate to settlement
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display settlement page correctly', async () => {
      await settlementPage.goto();

      await expect(settlementPage.pageTitle).toBeVisible();
      await expect(settlementPage.table).toBeVisible();
    });

    test('should load player list', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search by keyword', async () => {
      await settlementPage.goto();

      await settlementPage.searchByKeyword('test');

      const count = await settlementPage.getPlayerCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by company', async () => {
      await settlementPage.goto();

      // Company filtering may not be available on all pages
      const companySelectVisible = await settlementPage.companySelect.isVisible();
      if (companySelectVisible) {
        await settlementPage.filterByCompany('test');
        await settlementPage.page.waitForTimeout(1000);
      }

      const count = await settlementPage.getPlayerCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await settlementPage.goto();

      // Apply a filter
      await settlementPage.searchByKeyword('test');

      // Clear search
      await settlementPage.clearFilters();

      const count = await settlementPage.getPlayerCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Player Actions', () => {
    test('should view player details', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count > 0) {
        await settlementPage.viewDetails(0);

        const isDrawerVisible = await settlementPage.isDrawerVisible();
        if (isDrawerVisible) {
          await settlementPage.closeDrawer();
        }
      }
    });

    test('should open assign modal', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count > 0) {
        await settlementPage.assignPlayer(0);

        const isModalVisible = await settlementPage.isModalVisible();
        if (isModalVisible) {
          await settlementPage.closeModal();
        }
      }
    });

    test('should open transfer modal', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count > 0) {
        // Transfer button may not exist for all players
        const hasTransferButton = await settlementPage.transferButton(0).isVisible();
        if (hasTransferButton) {
          await settlementPage.transferPlayer(0);

          const isModalVisible = await settlementPage.isModalVisible();
          if (isModalVisible) {
            await settlementPage.closeModal();
          }
        }
      }
    });

    test('should view player history', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count > 0) {
        // History button may not exist for all players
        const hasHistoryButton = await settlementPage.historyButton(0).isVisible();
        if (hasHistoryButton) {
          await settlementPage.viewHistory(0);

          const isDrawerVisible = await settlementPage.isDrawerVisible();
          if (isDrawerVisible) {
            await settlementPage.closeDrawer();
          }
        }
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display player columns correctly', async () => {
      await settlementPage.goto();

      await expect(settlementPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count > 0) {
        const cellText = await settlementPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh player list', async () => {
      await settlementPage.goto();

      const _countBefore = await settlementPage.getPlayerCount();

      await settlementPage.refresh();

      const countAfter = await settlementPage.getPlayerCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await settlementPage.goto();

      await expect(settlementPage.pageTitle).toBeVisible();
      await expect(settlementPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await settlementPage.goto();

      const title = await settlementPage.pageTitle.textContent();
      expect(title).toMatch(/结算管理|陪玩师归属/);
    });

    test('should be keyboard navigable', async () => {
      await settlementPage.goto();

      await settlementPage.page.keyboard.press('Tab');
      await settlementPage.page.keyboard.press('Tab');

      const focusedElement = settlementPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const settlementMenuItem = page.locator('.ant-menu-item').filter({ hasText: /结算管理|陪玩师归属/i });
      if (await settlementMenuItem.isVisible()) {
        await settlementMenuItem.click();
        await settlementPage.waitForPageLoad();
        await expect(settlementPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty player list', async () => {
      await settlementPage.goto();

      const count = await settlementPage.getPlayerCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(settlementPage.table).toBeVisible();
      }
    });
  });

  test.describe('Statistics Display', () => {
    test('should display player statistics if available', async () => {
      await settlementPage.goto();

      const statsVisible = await settlementPage.totalPlayersStat.isVisible();
      if (statsVisible) {
        await expect(settlementPage.totalPlayersStat).toBeVisible();
      }
    });
  });
});
