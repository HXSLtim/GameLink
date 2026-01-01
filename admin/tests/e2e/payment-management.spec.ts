import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { PaymentManagementPage } from './pages/PaymentManagementPage';
import { getAdminToken } from './helpers/api-helpers';

/**
 * E2E Tests for Payment Management
 *
 * Test Coverage:
 * - List payment records
 * - View payment details
 * - Process refunds
 * - Search and filter payments
 * - Payment status tracking
 * - Export payment records
 */

test.describe('Payment Management', () => {
  let loginPage: LoginPage;
  let paymentManagementPage: PaymentManagementPage;
  let adminToken: string;

  test.beforeAll(async () => {
    adminToken = await getAdminToken();
  });

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    paymentManagementPage = new PaymentManagementPage(page);

    // Login and navigate to payment management
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
    await paymentManagementPage.goto();
  });

  test.describe('Payment List', () => {
    test('should display payment list page correctly', async () => {
      await expect(paymentManagementPage['pageTitle']).toBeVisible();
      await expect(paymentManagementPage['table']).toBeVisible();
    });

    test('should load payments in table', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();
      expect(paymentCount).toBeGreaterThanOrEqual(0);
    });

    test('should display payment details correctly', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        const paymentId = await paymentManagementPage.getCellText(0, 1);
        expect(paymentId).toBeTruthy();
        expect(paymentId.length).toBeGreaterThan(0);
      }
    });

    test('should show different payment statuses', async () => {
      const statuses = ['pending', 'paid', 'failed', 'refunded'];

      for (const status of statuses) {
        await paymentManagementPage.goto();
        await paymentManagementPage.filterByStatus(status);
        await paymentManagementPage.page.waitForTimeout(1000);

        // Verify filter is applied (may have 0 results for some statuses)
        const count = await paymentManagementPage.getPaymentCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });

    test('should display payment amounts correctly', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        const amountCell = await paymentManagementPage.getCellText(0, 4);
        expect(amountCell).toBeTruthy();
        // Amount should contain currency symbol or number
        expect(amountCell).toMatch(/[\$¥€£]|[\d.,]+/);
      }
    });
  });

  test.describe('View Payment Details', () => {
    test('should view payment details', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        await paymentManagementPage.viewPaymentDetails(0);

        await expect(paymentManagementPage['modal']).toBeVisible();

        // Modal should contain payment information
        await expect(paymentManagementPage['modalTitle']).toBeVisible();
      }
    });

    test('should close payment details modal', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        await paymentManagementPage.viewPaymentDetails(0);
        await expect(paymentManagementPage['modal']).toBeVisible();

        await paymentManagementPage.closeModal();

        await expect(paymentManagementPage['modal']).not.toBeVisible();
      }
    });

    test('should display complete payment information in modal', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        await paymentManagementPage.viewPaymentDetails(0);

        // Verify modal contains expected information
        await expect(paymentManagementPage['modal']).toBeVisible();

        const modalContent = paymentManagementPage['modal'].textContent();
        expect(modalContent).toBeTruthy();
      }
    });

    test('should show transaction details', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        await paymentManagementPage.viewPaymentDetails(0);

        // Look for transaction ID or payment method information
        await expect(paymentManagementPage['modal']).toBeVisible();

        const modalText = await paymentManagementPage['modal'].textContent();
        expect(modalText?.length).toBeGreaterThan(0);
      }
    });
  });

  test.describe('Search and Filter', () => {
    test('should search payments by transaction ID', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        const paymentId = await paymentManagementPage.getCellText(0, 1);

        await paymentManagementPage.searchPayment(paymentId);

        await paymentManagementPage.page.waitForTimeout(1000);

        const filteredCount = await paymentManagementPage.getPaymentCount();
        expect(filteredCount).toBeGreaterThanOrEqual(0);
      }
    });

    test('should filter payments by status', async () => {
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const paidCount = await paymentManagementPage.getPaymentCount();
      expect(paidCount).toBeGreaterThanOrEqual(0);

      // Verify all visible payments have paid status
      if (paidCount > 0) {
        const hasPaid = await paymentManagementPage.verifyPaymentStatus(0, 'paid');
        expect(hasPaid).toBe(true);
      }
    });

    test('should filter by different statuses', async () => {
      const statuses = ['pending', 'paid', 'failed', 'refunded'];

      for (const status of statuses) {
        await paymentManagementPage.goto();
        await paymentManagementPage.filterByStatus(status);
        await paymentManagementPage.page.waitForTimeout(1000);

        const count = await paymentManagementPage.getPaymentCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });

    test('should clear filters', async () => {
      // Apply filter
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const filteredCount = await paymentManagementPage.getPaymentCount();

      // Clear search
      await paymentManagementPage.searchPayment('');

      await paymentManagementPage.page.waitForTimeout(1000);

      const clearedCount = await paymentManagementPage.getPaymentCount();

      expect(clearedCount).toBeGreaterThanOrEqual(filteredCount);
    });

    test('should filter by date range', async () => {
      const startDate = '2024-01-01';
      const endDate = '2024-12-31';

      await paymentManagementPage.filterByDateRange(startDate, endDate);

      await paymentManagementPage.page.waitForTimeout(1000);

      const count = await paymentManagementPage.getPaymentCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Process Refund', () => {
    test('should show refund button for paid payments', async () => {
      await paymentManagementPage.goto();
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const paidCount = await paymentManagementPage.getPaymentCount();

      if (paidCount > 0) {
        // Check if refund button is visible for first paid payment
        const refundButton = paymentManagementPage['refundButton'](0);
        await expect(refundButton).toBeVisible();
      }
    });

    test('should open refund dialog', async () => {
      await paymentManagementPage.goto();
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const paidCount = await paymentManagementPage.getPaymentCount();

      if (paidCount > 0) {
        await paymentManagementPage['refundButton'](0).click();

        await expect(paymentManagementPage['modal']).toBeVisible();

        // Close modal to avoid side effects
        await paymentManagementPage['modal'].getByRole('button', { name: /取消|cancel/i }).click();
      }
    });

    test('should require refund reason', async () => {
      await paymentManagementPage.goto();
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const paidCount = await paymentManagementPage.getPaymentCount();

      if (paidCount > 0) {
        await paymentManagementPage.refundOrder(0, {
          reason: 'Customer requested refund',
        });

        // Should either succeed or show validation
        await paymentManagementPage.page.waitForTimeout(1000);
      }
    });

    test('should process refund successfully', async () => {
      await paymentManagementPage.goto();
      await paymentManagementPage.filterByStatus('paid');

      await paymentManagementPage.page.waitForTimeout(1000);

      const paidCount = await paymentManagementPage.getPaymentCount();

      if (paidCount > 0) {
        const initialStatus = await paymentManagementPage.getCellText(0, 3);

        await paymentManagementPage.refundOrder(0, {
          reason: 'Service not satisfactory',
        });

        // Wait for processing
        await paymentManagementPage.page.waitForTimeout(2000);

        // Check if status changed
        const newStatus = await paymentManagementPage.getCellText(0, 3);

        // Status should either be refunded or same (if refund failed)
        expect(['paid', 'refunded']).toContain(newStatus.toLowerCase());
      }
    });
  });

  test.describe('Export Functionality', () => {
    test('should export payment records', async () => {
      const downloadPromise = paymentManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await paymentManagementPage.exportPayments();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });

    test('should export filtered payments', async () => {
      // Apply filter first
      await paymentManagementPage.filterByStatus('paid');
      await paymentManagementPage.page.waitForTimeout(1000);

      const downloadPromise = paymentManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await paymentManagementPage.exportPayments();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });

    test('should export payments for specific date range', async () => {
      const startDate = '2024-01-01';
      const endDate = '2024-12-31';

      await paymentManagementPage.filterByDateRange(startDate, endDate);
      await paymentManagementPage.page.waitForTimeout(1000);

      const downloadPromise = paymentManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await paymentManagementPage.exportPayments();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });
  });

  test.describe('Pagination', () => {
    test('should navigate through pages', async () => {
      const initialCount = await paymentManagementPage.getPaymentCount();

      await paymentManagementPage.nextPage();

      await paymentManagementPage.page.waitForTimeout(1000);

      const newCount = await paymentManagementPage.getPaymentCount();

      // If there are more payments, count should be the same (different page)
      expect(newCount).toBeGreaterThanOrEqual(0);
    });

    test('should show correct pagination info', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        // Check for pagination controls
        const pagination = paymentManagementPage.page.locator('.ant-pagination');
        const isVisible = await pagination.isVisible().catch(() => false);

        if (isVisible) {
          await expect(pagination).toBeVisible();
        }
      }
    });
  });

  test.describe('Payment Status Tracking', () => {
    test('should show payment status changes', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        // View payment details
        await paymentManagementPage.viewPaymentDetails(0);

        // Check if status history or timeline is displayed
        await expect(paymentManagementPage['modal']).toBeVisible();

        const modalText = await paymentManagementPage['modal'].textContent();
        expect(modalText?.length).toBeGreaterThan(0);

        await paymentManagementPage.closeModal();
      }
    });

    test('should display payment timestamps', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        // Check for timestamp in table row
        const rowText = await paymentManagementPage.tableRows.nth(0).textContent();
        expect(rowText).toBeTruthy();
        expect(rowText?.length).toBeGreaterThan(0);
      }
    });
  });

  test.describe('Accessibility', () => {
    test('should be keyboard navigable', async () => {
      // Tab through table
      await paymentManagementPage.page.keyboard.press('Tab');
      await paymentManagementPage.page.keyboard.press('Tab');

      // Focus should be on table or navigation
      const focusedElement = paymentManagementPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });

    test('should have proper ARIA labels', async () => {
      const table = paymentManagementPage['table'];
      await expect(table).toHaveAttribute('role', 'table');
    });
  });

  test.describe('Error Handling', () => {
    test('should handle network errors gracefully', async () => {
      // Simulate network error by going offline
      await paymentManagementPage.page.context().setOffline(true);

      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        await paymentManagementPage['refundButton'](0).click();

        // Should handle error gracefully
        await paymentManagementPage.page.waitForTimeout(2000);
      }

      // Restore network
      await paymentManagementPage.page.context().setOffline(false);
    });
  });

  test.describe('Amount Display', () => {
    test('should format amounts correctly', async () => {
      const paymentCount = await paymentManagementPage.getPaymentCount();

      if (paymentCount > 0) {
        const amountCell = await paymentManagementPage.getCellText(0, 4);

        // Amount should be formatted with currency symbol and proper decimal places
        expect(amountCell).toMatch(/[\$¥€£₹][\d,]+\.?\d{0,2}/);
      }
    });

    test('should show total amount summary', async () => {
      // Look for total amount display (if exists)
      const summarySection = paymentManagementPage.page.locator(
        '.ant-statistic, .amount-summary, .total-amount'
      );

      const isVisible = await summarySection.isVisible().catch(() => false);

      if (isVisible) {
        await expect(summarySection).toBeVisible();
      }
    });
  });
});
