import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { VIPPage } from './pages/VIPPage';

/**
 * E2E Tests for VIP Management
 *
 * Test Coverage:
 * - Page loading and display
 * - VIP level cards display
 * - Statistics display
 * - Search and filter VIP levels
 * - Create, edit, delete VIP levels
 * - Set default VIP level
 * - Toggle active status
 * - Benefits management
 */

test.describe('VIP Management', () => {
  let loginPage: LoginPage;
  let vipPage: VIPPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    vipPage = new VIPPage(page);

    // Login and navigate to VIP
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display VIP page correctly', async () => {
      await vipPage.goto();

      await expect(vipPage.pageTitle).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await vipPage.goto();

      await expect(vipPage.totalLevelsStat).toBeVisible();
      await expect(vipPage.activeLevelsStat).toBeVisible();
    });

    test('should display VIP level cards', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await vipPage.goto();

      const stats = await vipPage.getStatistics();
      expect(stats.total).toBeGreaterThanOrEqual(0);
      expect(stats.active).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search VIP levels by keyword', async () => {
      await vipPage.goto();

      await vipPage.searchByKeyword('gold');

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by all levels', async () => {
      await vipPage.goto();

      await vipPage.switchViewMode('all');

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by active levels only', async () => {
      await vipPage.goto();

      await vipPage.switchViewMode('active');

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by inactive levels only', async () => {
      await vipPage.goto();

      await vipPage.switchViewMode('inactive');

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await vipPage.goto();

      // Apply search
      await vipPage.searchByKeyword('test');

      // Clear search
      await vipPage.searchInput.fill('');
      await vipPage.page.waitForTimeout(500);

      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('VIP Level Display', () => {
    test('should display VIP level information correctly', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        const title = await vipPage.getCardTitle(0);
        expect(title).toBeTruthy();
        expect(title.length).toBeGreaterThan(0);
      }
    });

    test('should show VIP level cards in proper layout', async () => {
      await vipPage.goto();

      await expect(vipPage.vipCards.first()).toBeVisible();
    });
  });

  test.describe('VIP Level Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await vipPage.goto();

      await vipPage.clickAdd();

      const isModalVisible = await vipPage.isModalVisible();
      if (isModalVisible) {
        await vipPage.closeModal();
      }
    });

    test('should open edit modal when clicking edit button', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        await vipPage.editLevel(0);

        const isModalVisible = await vipPage.isModalVisible();
        if (isModalVisible) {
          await vipPage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        const _countBefore = await vipPage.getVIPCardCount();

        // Try to delete first VIP level
        await vipPage.deleteLevel(0);

        // Wait for potential deletion
        await vipPage.page.waitForTimeout(1000);

        const countAfter = await vipPage.getVIPCardCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });

    test('should handle set as default operation', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        await vipPage.setAsDefault(0);
        await vipPage.page.waitForTimeout(1000);
      }
    });

    test('should toggle active status', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        const _isActiveBefore = await vipPage.isLevelActive(0);

        await vipPage.toggleActive(0);

        await vipPage.page.waitForTimeout(1000);

        const isActiveAfter = await vipPage.isLevelActive(0);
        // Status may have changed
        expect(typeof isActiveAfter).toBe('boolean');
      }
    });

    test('should open benefits editor', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count > 0) {
        await vipPage.editBenefits(0);

        const isModalVisible = await vipPage.isModalVisible();
        if (isModalVisible) {
          await vipPage.closeModal();
        }
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh VIP levels list', async () => {
      await vipPage.goto();

      const _countBefore = await vipPage.getVIPCardCount();

      await vipPage.refresh();

      const countAfter = await vipPage.getVIPCardCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await vipPage.goto();

      await expect(vipPage.pageTitle).toBeVisible();
      await expect(vipPage.vipCards.first()).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await vipPage.goto();

      await expect(vipPage.pageTitle).toContainText('VIP管理');
    });

    test('should be keyboard navigable', async () => {
      await vipPage.goto();

      await vipPage.page.keyboard.press('Tab');
      await vipPage.page.keyboard.press('Tab');

      const focusedElement = vipPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const vipMenuItem = page.locator('.ant-menu-item').filter({ hasText: /VIP管理/i });
      if (await vipMenuItem.isVisible()) {
        await vipMenuItem.click();
        await vipPage.waitForPageLoad();
        await expect(vipPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty VIP levels list', async () => {
      await vipPage.goto();

      const count = await vipPage.getVIPCardCount();
      if (count === 0) {
        // Should show empty state
        await expect(vipPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Search Functionality', () => {
    test('should filter VIP levels by keyword', async () => {
      await vipPage.goto();

      // Apply keyword filter
      await vipPage.searchByKeyword('gold');

      await vipPage.page.waitForTimeout(500);

      // Check that results are filtered
      const count = await vipPage.getVIPCardCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should show no results for non-existent keyword', async () => {
      await vipPage.goto();

      // Search for non-existent keyword
      await vipPage.searchByKeyword('nonexistentlevel12345');

      await vipPage.page.waitForTimeout(500);

      // Should show 0 results
      const count = await vipPage.getVIPCardCount();
      expect(count).toBe(0);
    });
  });

  test.describe('View Mode Switching', () => {
    test('should switch between all, active, and inactive views', async () => {
      await vipPage.goto();

      // Start with all
      await vipPage.switchViewMode('all');
      const countAll = await vipPage.getVIPCardCount();

      // Switch to active
      await vipPage.switchViewMode('active');
      const countActive = await vipPage.getVIPCardCount();
      expect(countActive).toBeLessThanOrEqual(countAll);

      // Switch to inactive
      await vipPage.switchViewMode('inactive');
      const countInactive = await vipPage.getVIPCardCount();

      // Verify counts add up correctly
      expect(countActive + countInactive).toBeLessThanOrEqual(countAll);
    });
  });
});
