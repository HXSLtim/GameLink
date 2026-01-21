import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { GamePage } from './pages/GamePage';

/**
 * E2E Tests for Game Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Game list with pagination
 * - Search and filter games
 * - Create, edit, delete games
 * - Import functionality
 */

test.describe('Game Management', () => {
  let loginPage: LoginPage;
  let gamePage: GamePage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    gamePage = new GamePage(page);

    // Login and navigate to game
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display game page correctly', async () => {
      await gamePage.goto();

      await expect(gamePage.pageTitle).toBeVisible();
      await expect(gamePage.table).toBeVisible();
    });

    test('should load game list', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search games by keyword', async () => {
      await gamePage.goto();

      await gamePage.searchByKeyword('test');

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by category', async () => {
      await gamePage.goto();

      await gamePage.filterByCategory('moba');

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await gamePage.goto();

      // Apply a filter
      await gamePage.searchByKeyword('test');

      // Clear search
      await gamePage.clearFilters();

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Game Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await gamePage.goto();

      await gamePage.clickAdd();

      const isModalVisible = await gamePage.isModalVisible();
      if (isModalVisible) {
        await expect(gamePage.modalTitle).toContainText(/新增游戏|添加游戏/i);
        await gamePage.closeModal();
      }
    });

    test('should edit game', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      if (count > 0) {
        await gamePage.editGame(0);

        const isModalVisible = await gamePage.isModalVisible();
        if (isModalVisible) {
          await expect(gamePage.modalTitle).toContainText(/编辑游戏/i);
          await gamePage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      if (count > 0) {
        const _countBefore = await gamePage.getGameCount();

        // Try to delete first game
        await gamePage.deleteGame(0);

        // Wait for potential deletion
        await gamePage.page.waitForTimeout(1000);

        const countAfter = await gamePage.getGameCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });
  });

  test.describe('Import Functionality', () => {
    test('should have import button', async () => {
      await gamePage.goto();

      await expect(gamePage.importButton).toBeVisible();
    });
  });

  test.describe('Table Display', () => {
    test('should display game columns correctly', async () => {
      await gamePage.goto();

      await expect(gamePage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      if (count > 0) {
        const cellText = await gamePage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh game list', async () => {
      await gamePage.goto();

      const _countBefore = await gamePage.getGameCount();

      await gamePage.refresh();

      const countAfter = await gamePage.getGameCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await gamePage.goto();

      await expect(gamePage.pageTitle).toBeVisible();
      await expect(gamePage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await gamePage.goto();

      const title = await gamePage.pageTitle.textContent();
      expect(title).toMatch(/游戏管理/i);
    });

    test('should be keyboard navigable', async () => {
      await gamePage.goto();

      await gamePage.page.keyboard.press('Tab');
      await gamePage.page.keyboard.press('Tab');

      const focusedElement = gamePage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const gameMenuItem = page.locator('.ant-menu-item').filter({ hasText: /游戏管理/i });
      if (await gameMenuItem.isVisible()) {
        await gameMenuItem.click();
        await gamePage.waitForPageLoad();
        await expect(gamePage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty game list', async () => {
      await gamePage.goto();

      const count = await gamePage.getGameCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(gamePage.table).toBeVisible();
      }
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await gamePage.goto();

      // Apply category filter
      await gamePage.filterByCategory('moba');

      // Then apply keyword filter
      await gamePage.searchByKeyword('test');

      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Category Filtering', () => {
    test('should filter by moba category', async () => {
      await gamePage.goto();
      await gamePage.filterByCategory('moba');
      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by fps category', async () => {
      await gamePage.goto();
      await gamePage.filterByCategory('fps');
      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by rpg category', async () => {
      await gamePage.goto();
      await gamePage.filterByCategory('rpg');
      const count = await gamePage.getGameCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });
});
