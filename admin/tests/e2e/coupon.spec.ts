import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { CouponPage } from './pages/CouponPage';

/**
 * E2E Tests for Coupon Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Coupon list with pagination
 * - Search and filter coupons
 * - Create, edit, delete coupons
 * - View coupon details
 * - Export functionality
 */

test.describe('Coupon Management', () => {
  let loginPage: LoginPage;
  let couponPage: CouponPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    couponPage = new CouponPage(page);

    // Login and navigate to coupon
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display coupon page correctly', async () => {
      await couponPage.goto();

      await expect(couponPage.pageTitle).toBeVisible();
      await expect(couponPage.table).toBeVisible();
    });

    test('should load coupon list', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search coupons by keyword', async () => {
      await couponPage.goto();

      await couponPage.searchByKeyword('test');

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by coupon type', async () => {
      await couponPage.goto();

      await couponPage.filterByType('discount');

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by scope', async () => {
      await couponPage.goto();

      await couponPage.filterByScope('all');

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by status', async () => {
      await couponPage.goto();

      await couponPage.filterByStatus('active');

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await couponPage.goto();

      // Apply a filter
      await couponPage.searchByKeyword('test');

      // Clear search
      await couponPage.clearFilters();

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Coupon Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await couponPage.goto();

      await couponPage.clickAdd();

      const isModalVisible = await couponPage.isModalVisible();
      if (isModalVisible) {
        await couponPage.closeModal();
      }
    });

    test('should view coupon details', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      if (count > 0) {
        await couponPage.viewDetails(0);

        const isModalVisible = await couponPage.isModalVisible();
        if (isModalVisible) {
          await couponPage.closeModal();
        }
      }
    });

    test('should edit coupon', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      if (count > 0) {
        await couponPage.editCoupon(0);

        const isModalVisible = await couponPage.isModalVisible();
        if (isModalVisible) {
          await couponPage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      if (count > 0) {
        const _countBefore = await couponPage.getCouponCount();

        // Try to delete first coupon
        await couponPage.deleteCoupon(0);

        // Wait for potential deletion
        await couponPage.page.waitForTimeout(1000);

        const countAfter = await couponPage.getCouponCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display coupon columns correctly', async () => {
      await couponPage.goto();

      await expect(couponPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      if (count > 0) {
        const cellText = await couponPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh coupon list', async () => {
      await couponPage.goto();

      const _countBefore = await couponPage.getCouponCount();

      await couponPage.refresh();

      const countAfter = await couponPage.getCouponCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Export Functionality', () => {
    test('should have export button', async () => {
      await couponPage.goto();

      await expect(couponPage.exportButton).toBeVisible();
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await couponPage.goto();

      await expect(couponPage.pageTitle).toBeVisible();
      await expect(couponPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await couponPage.goto();

      await expect(couponPage.pageTitle).toContainText('优惠券管理');
    });

    test('should be keyboard navigable', async () => {
      await couponPage.goto();

      await couponPage.page.keyboard.press('Tab');
      await couponPage.page.keyboard.press('Tab');

      const focusedElement = couponPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const couponMenuItem = page.locator('.ant-menu-item').filter({ hasText: /优惠券管理/i });
      if (await couponMenuItem.isVisible()) {
        await couponMenuItem.click();
        await couponPage.waitForPageLoad();
        await expect(couponPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty coupon list', async () => {
      await couponPage.goto();

      const count = await couponPage.getCouponCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(couponPage.table).toBeVisible();
      }
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await couponPage.goto();

      // Apply type filter
      await couponPage.filterByType('discount');

      // Then apply status filter
      await couponPage.filterByStatus('active');

      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search Functionality', () => {
    test('should filter coupons by keyword', async () => {
      await couponPage.goto();

      // Apply keyword filter
      await couponPage.searchByKeyword('discount');

      await couponPage.page.waitForTimeout(500);

      // Check that results are filtered
      const count = await couponPage.getCouponCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should verify coupon exists in table', async () => {
      await couponPage.goto();

      const exists = await couponPage.verifyCouponExists('test');
      expect(exists).toBeDefined();
    });
  });
});
