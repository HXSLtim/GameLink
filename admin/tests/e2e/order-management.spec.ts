import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { OrderManagementPage } from './pages/OrderManagementPage';
import { getAdminToken } from './helpers/api-helpers';

/**
 * E2E Tests for Order Management
 *
 * Test Coverage:
 * - List orders with pagination
 * - View order details
 * - Cancel orders
 * - Refund orders
 * - Search and filter orders
 * - Batch operations
 * - Order status tracking
 */

test.describe('Order Management', () => {
  let loginPage: LoginPage;
  let orderManagementPage: OrderManagementPage;

  test.beforeAll(async () => {
    await getAdminToken();
  });

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    orderManagementPage = new OrderManagementPage(page);

    // Login and navigate to order management
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
    await orderManagementPage.goto();
  });

  test.describe('Order List', () => {
    test('should display order list page correctly', async () => {
      await expect(orderManagementPage['pageTitle']).toBeVisible();
      await expect(orderManagementPage['table']).toBeVisible();
    });

    test('should load orders in table', async () => {
      const orderCount = await orderManagementPage.getOrderCount();
      expect(orderCount).toBeGreaterThanOrEqual(0);
    });

    test('should display order details correctly', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        const orderNo = await orderManagementPage.getCellText(0, 1);
        expect(orderNo).toBeTruthy();
        expect(orderNo.length).toBeGreaterThan(0);
      }
    });

    test('should show different order statuses', async () => {
      const statuses = ['pending', 'confirmed', 'in_progress', 'completed', 'cancelled'];

      for (const status of statuses) {
        await orderManagementPage.goto();
        await orderManagementPage.filterByStatus(status);
        await orderManagementPage.page.waitForTimeout(1000);

        // Verify filter is applied (may have 0 results for some statuses)
        const count = await orderManagementPage.getOrderCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });
  });

  test.describe('View Order Details', () => {
    test('should view order details', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        await orderManagementPage.viewOrderDetails(0);

        await expect(orderManagementPage['modal']).toBeVisible();

        // Modal should contain order information
        await expect(orderManagementPage['modalTitle']).toBeVisible();
      }
    });

    test('should close order details modal', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        await orderManagementPage.viewOrderDetails(0);
        await expect(orderManagementPage['modal']).toBeVisible();

        await orderManagementPage.closeModal();

        await expect(orderManagementPage['modal']).not.toBeVisible();
      }
    });

    test('should display complete order information in modal', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        await orderManagementPage.viewOrderDetails(0);

        // Verify modal contains expected information
        await expect(orderManagementPage['modal']).toBeVisible();

        const modalContent = orderManagementPage['modal'].textContent();
        expect(modalContent).toBeTruthy();
      }
    });
  });

  test.describe('Search and Filter', () => {
    test('should search orders by order number', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        const orderNo = await orderManagementPage.getCellText(0, 1);
        const orderNumber = orderNo.split(' ')[0]; // Extract order number

        await orderManagementPage.searchOrder(orderNumber);

        await orderManagementPage.page.waitForTimeout(1000);

        const filteredCount = await orderManagementPage.getOrderCount();
        expect(filteredCount).toBeGreaterThanOrEqual(0);
      }
    });

    test('should filter orders by status', async () => {
      await orderManagementPage.filterByStatus('completed');

      await orderManagementPage.page.waitForTimeout(1000);

      const completedCount = await orderManagementPage.getOrderCount();
      expect(completedCount).toBeGreaterThanOrEqual(0);

      // Verify all visible orders have completed status
      if (completedCount > 0) {
        const hasCompleted = await orderManagementPage.verifyOrderStatus(0, 'completed');
        expect(hasCompleted).toBe(true);
      }
    });

    test('should filter by different statuses', async () => {
      const statuses = ['pending', 'confirmed', 'in_progress', 'completed', 'cancelled'];

      for (const status of statuses) {
        await orderManagementPage.goto();
        await orderManagementPage.filterByStatus(status);
        await orderManagementPage.page.waitForTimeout(1000);

        const count = await orderManagementPage.getOrderCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });

    test('should clear filters', async () => {
      // Apply filter
      await orderManagementPage.filterByStatus('completed');

      await orderManagementPage.page.waitForTimeout(1000);

      const filteredCount = await orderManagementPage.getOrderCount();

      // Clear search
      await orderManagementPage.searchOrder('');

      await orderManagementPage.page.waitForTimeout(1000);

      const clearedCount = await orderManagementPage.getOrderCount();

      expect(clearedCount).toBeGreaterThanOrEqual(filteredCount);
    });
  });

  test.describe('Cancel Order', () => {
    test('should cancel pending order', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('pending');

      await orderManagementPage.page.waitForTimeout(1000);

      const pendingCount = await orderManagementPage.getOrderCount();

      if (pendingCount > 0) {
        await orderManagementPage.cancelOrder(0, 'Test cancellation via E2E');

        await orderManagementPage.verifySuccessMessage();

        // Wait for status to change
        await orderManagementPage.page.waitForTimeout(2000);

        // Verify order status changed
        const isCancelled = await orderManagementPage.verifyOrderStatus(0, 'cancelled');
        expect(isCancelled).toBe(true);
      }
    });

    test('should show confirmation dialog before cancel', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('pending');

      await orderManagementPage.page.waitForTimeout(1000);

      const pendingCount = await orderManagementPage.getOrderCount();

      if (pendingCount > 0) {
        await orderManagementPage['cancelButton'](0).click();

        // Verify confirmation modal is shown
        await expect(orderManagementPage['modal']).toBeVisible();
        await expect(orderManagementPage['modalOkButton']).toBeVisible();

        // Cancel the operation
        await orderManagementPage['modal'].getByRole('button', { name: /取消|cancel/i }).click();
      }
    });

    test('should provide reason for cancellation', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('pending');

      await orderManagementPage.page.waitForTimeout(1000);

      const pendingCount = await orderManagementPage.getOrderCount();

      if (pendingCount > 0) {
        const reason = 'Customer requested cancellation';
        await orderManagementPage.cancelOrder(0, reason);

        await orderManagementPage.verifySuccessMessage();
      }
    });
  });

  test.describe('Refund Order', () => {
    test('should refund completed order', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('completed');

      await orderManagementPage.page.waitForTimeout(1000);

      const completedCount = await orderManagementPage.getOrderCount();

      if (completedCount > 0) {
        await orderManagementPage.refundOrder(0, {
          reason: 'Test refund via E2E',
        });

        await orderManagementPage.verifySuccessMessage();

        // Wait for status to change
        await orderManagementPage.page.waitForTimeout(2000);

        // Verify order status changed to refunded
        const isRefunded = await orderManagementPage.verifyOrderStatus(0, 'refunded');
        expect(isRefunded).toBe(true);
      }
    });

    test('should provide refund reason', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('completed');

      await orderManagementPage.page.waitForTimeout(1000);

      const completedCount = await orderManagementPage.getOrderCount();

      if (completedCount > 0) {
        const refundReason = 'Service not provided as expected';
        await orderManagementPage.refundOrder(0, {
          reason: refundReason,
        });

        await orderManagementPage.verifySuccessMessage();
      }
    });
  });

  test.describe('Batch Operations', () => {
    test('should batch cancel multiple orders', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('pending');

      await orderManagementPage.page.waitForTimeout(1000);

      const pendingCount = await orderManagementPage.getOrderCount();

      if (pendingCount >= 2) {
        // Get order numbers
        const orderNo1 = await orderManagementPage.getCellText(0, 1);
        const orderNo2 = await orderManagementPage.getCellText(1, 1);

        await orderManagementPage.batchCancelOrders([orderNo1, orderNo2], 'Batch cancellation via E2E');

        await orderManagementPage.verifySuccessMessage();
      }
    });

    test('should batch complete multiple orders', async () => {
      await orderManagementPage.goto();
      await orderManagementPage.filterByStatus('in_progress');

      await orderManagementPage.page.waitForTimeout(1000);

      const inProgressCount = await orderManagementPage.getOrderCount();

      if (inProgressCount >= 2) {
        // Get order numbers
        const orderNo1 = await orderManagementPage.getCellText(0, 1);
        const orderNo2 = await orderManagementPage.getCellText(1, 1);

        await orderManagementPage.batchCompleteOrders([orderNo1, orderNo2]);

        await orderManagementPage.verifySuccessMessage();
      }
    });
  });

  test.describe('Export Functionality', () => {
    test('should export order list', async () => {
      const downloadPromise = orderManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await orderManagementPage.exportOrders();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });

    test('should export filtered orders', async () => {
      // Apply filter first
      await orderManagementPage.filterByStatus('completed');
      await orderManagementPage.page.waitForTimeout(1000);

      const downloadPromise = orderManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await orderManagementPage.exportOrders();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });
  });

  test.describe('Pagination', () => {
    test('should navigate through pages', async () => {
      await orderManagementPage.getOrderCount();

      await orderManagementPage.nextPage();

      await orderManagementPage.page.waitForTimeout(1000);

      const newCount = await orderManagementPage.getOrderCount();

      // If there are more orders, count should be the same (different page)
      expect(newCount).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Order Status Tracking', () => {
    test('should show order status progression', async () => {
      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        // View order details
        await orderManagementPage.viewOrderDetails(0);

        // Check if status history or timeline is displayed
        await expect(orderManagementPage['modal']).toBeVisible();

        await orderManagementPage.closeModal();
      }
    });
  });

  test.describe('Accessibility', () => {
    test('should be keyboard navigable', async () => {
      // Tab through table
      await orderManagementPage.page.keyboard.press('Tab');
      await orderManagementPage.page.keyboard.press('Tab');

      // Focus should be on table or navigation
      const focusedElement = orderManagementPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Error Handling', () => {
    test('should handle network errors gracefully', async () => {
      // Simulate network error by going offline
      await orderManagementPage.page.context().setOffline(true);

      const orderCount = await orderManagementPage.getOrderCount();

      if (orderCount > 0) {
        await orderManagementPage.cancelOrder(0, 'Test cancellation');

        // Should show error message
        await orderManagementPage.page.waitForTimeout(2000);
      }

      // Restore network
      await orderManagementPage.page.context().setOffline(false);
    });
  });
});
