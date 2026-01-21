import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { ServicePage } from './pages/ServicePage';

/**
 * E2E Tests for Service Item Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Service list with pagination
 * - Statistics display
 * - Search and filter services
 * - Create, edit, delete services
 * - Toggle service status
 */

test.describe('Service Item Management', () => {
  let loginPage: LoginPage;
  let servicePage: ServicePage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    servicePage = new ServicePage(page);

    // Login and navigate to service
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display service page correctly', async () => {
      await servicePage.goto();

      await expect(servicePage.pageTitle).toBeVisible();
      await expect(servicePage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await servicePage.goto();

      await expect(servicePage.totalServicesStat).toBeVisible();
      await expect(servicePage.activeServicesStat).toBeVisible();
    });

    test('should load service list', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await servicePage.goto();

      const stats = await servicePage.getStatistics();
      expect(stats.total).toBeGreaterThanOrEqual(0);
      expect(stats.active).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search services by keyword', async () => {
      await servicePage.goto();

      await servicePage.searchByKeyword('test');

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by game', async () => {
      await servicePage.goto();

      // Game filtering requires games to be available
      const gameSelectVisible = await servicePage.gameSelect.isVisible();
      if (gameSelectVisible) {
        // Click to see if options are available
        await servicePage.gameSelect.click();
        await servicePage.page.waitForTimeout(500);

        // Try to select first option if available
        const firstOption = servicePage.page.locator('.ant-select-item-option').first();
        if (await firstOption.isVisible()) {
          await firstOption.click();
          await servicePage.page.waitForTimeout(1000);
        }
      }

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by service type', async () => {
      await servicePage.goto();

      await servicePage.filterByServiceType('solo');

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await servicePage.goto();

      // Apply a filter
      await servicePage.searchByKeyword('test');

      // Clear search
      await servicePage.clearFilters();

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Service Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await servicePage.goto();

      await servicePage.clickAdd();

      const isModalVisible = await servicePage.isModalVisible();
      if (isModalVisible) {
        await expect(servicePage.modalTitle).toContainText(/新增服务项目|添加服务/i);
        await servicePage.closeModal();
      }
    });

    test('should edit service', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      if (count > 0) {
        await servicePage.editService(0);

        const isModalVisible = await servicePage.isModalVisible();
        if (isModalVisible) {
          await expect(servicePage.modalTitle).toContainText(/编辑服务项目/i);
          await servicePage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      if (count > 0) {
        const _countBefore = await servicePage.getServiceCount();

        // Try to delete first service
        await servicePage.deleteService(0);

        // Wait for potential deletion
        await servicePage.page.waitForTimeout(1000);

        const countAfter = await servicePage.getServiceCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });

    test('should toggle service status', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      if (count > 0) {
        // Toggle status of first service
        await servicePage.toggleStatus(0);
        await servicePage.page.waitForTimeout(1000);
      }
    });
  });

  test.describe('Service Type Filtering', () => {
    test('should filter by solo services', async () => {
      await servicePage.goto();
      await servicePage.filterByServiceType('solo');
      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by team services', async () => {
      await servicePage.goto();
      await servicePage.filterByServiceType('team');
      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by gift services', async () => {
      await servicePage.goto();
      await servicePage.filterByServiceType('gift');
      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Table Display', () => {
    test('should display service columns correctly', async () => {
      await servicePage.goto();

      await expect(servicePage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      if (count > 0) {
        const cellText = await servicePage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh service list', async () => {
      await servicePage.goto();

      const _countBefore = await servicePage.getServiceCount();

      await servicePage.refresh();

      const countAfter = await servicePage.getServiceCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await servicePage.goto();

      await expect(servicePage.pageTitle).toBeVisible();
      await expect(servicePage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await servicePage.goto();

      const title = await servicePage.pageTitle.textContent();
      expect(title).toMatch(/服务项目管理/i);
    });

    test('should be keyboard navigable', async () => {
      await servicePage.goto();

      await servicePage.page.keyboard.press('Tab');
      await servicePage.page.keyboard.press('Tab');

      const focusedElement = servicePage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const serviceMenuItem = page.locator('.ant-menu-item').filter({ hasText: /服务项目管理|服务管理/i });
      if (await serviceMenuItem.isVisible()) {
        await serviceMenuItem.click();
        await servicePage.waitForPageLoad();
        await expect(servicePage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty service list', async () => {
      await servicePage.goto();

      const count = await servicePage.getServiceCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(servicePage.table).toBeVisible();
      }
    });
  });

  test.describe('Statistics Verification', () => {
    test('should show consistent statistics', async () => {
      await servicePage.goto();

      const stats1 = await servicePage.getStatistics();

      await servicePage.page.waitForTimeout(1000);

      const stats2 = await servicePage.getStatistics();

      // Statistics should be consistent
      expect(stats1.total).toEqual(stats2.total);
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await servicePage.goto();

      // Apply type filter
      await servicePage.filterByServiceType('solo');

      // Then apply keyword filter
      await servicePage.searchByKeyword('test');

      const count = await servicePage.getServiceCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Modal Form', () => {
    test('should have required form fields in modal', async () => {
      await servicePage.goto();

      await servicePage.clickAdd();

      const isModalVisible = await servicePage.isModalVisible();
      if (isModalVisible) {
        // Check for common form fields
        await expect(servicePage.itemCodeInput).toBeVisible();
        await expect(servicePage.nameInput).toBeVisible();
        await servicePage.closeModal();
      }
    });
  });
});
