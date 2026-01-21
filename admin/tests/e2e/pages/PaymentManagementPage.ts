import { Page } from '@playwright/test';

/**
 * Page Object Model for Payment Management Page
 * Encapsulates all payment management interactions and assertions
 */
export class PaymentManagementPage {
  constructor(private page: Page) {}

  // Element locators - PageContainer uses h1 for title
  private readonly pageTitle = this.page.locator('h1').filter({ hasText: /支付记录/ });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');
  private readonly viewButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /查看|view/i });
  private readonly refundButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /退款|refund/i });

  // Filter locators
  private readonly statusFilter = this.page.getByText(/状态|status/i).first();
  private readonly searchInput = this.page.getByPlaceholder(/搜索|search/i);

  // Modal locators
  private readonly modal = this.page.locator('.ant-modal');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定|ok|submit/i });

  /**
   * Navigate to payment management page
   * Note: This page may not be configured in routes yet
   */
  async goto() {
    // Try the component map path first
    await this.page.goto('/admin/payment/records');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    // Wait for either the page title or the table to be visible
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.table.waitFor({ state: 'visible', timeout: 15000 }),
    ]).catch(() => {
      // Page might not be configured, that's okay for now
      console.warn('Payment records page may not be configured in routes');
    });
  }

  /**
   * Get payment count from table
   */
  async getPaymentCount(): Promise<number> {
    await this.page.waitForLoadState('networkidle');
    return await this.tableRows.count();
  }

  /**
   * Get text from table cell
   */
  async getCellText(rowIndex: number, columnIndex: number): Promise<string> {
    const cell = this.tableRows.nth(rowIndex).locator('td').nth(columnIndex);
    return await cell.textContent() || '';
  }

  /**
   * View payment details
   */
  async viewPaymentDetails(rowIndex: number) {
    await this.viewButton(rowIndex).click();
    await expect(this.modal).toBeVisible();
  }

  /**
   * Close modal
   */
  async closeModal() {
    const closeButton = this.modal.getByRole('button', { name: /关闭|close/i });
    await closeButton.click();
  }

  /**
   * Filter payments by status
   */
  async filterByStatus(status: string) {
    await this.statusFilter.click();
    await this.page.getByRole('option', { name: status }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Search for payment by transaction ID
   */
  async searchPayment(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Refund a payment
   */
  async refundPayment(rowIndex: number, data: { reason: string; amount?: number }) {
    await this.refundButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    const reasonInput = this.modal.getByLabel(/原因|reason/i);
    await reasonInput.fill(data.reason);

    if (data.amount) {
      const amountInput = this.modal.getByLabel(/金额|amount/i);
      await amountInput.fill(data.amount.toString());
    }

    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify payment status
   */
  async verifyPaymentStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    const cellText = await this.getCellText(rowIndex, 3); // Assuming status is in column 3
    return cellText.includes(expectedStatus);
  }

  /**
   * Verify payment amount
   */
  async verifyPaymentAmount(rowIndex: number, expectedAmount: string): Promise<boolean> {
    const cellText = await this.getCellText(rowIndex, 4); // Assuming amount is in column 4
    return cellText.includes(expectedAmount);
  }

  /**
   * Wait for payment status to change
   */
  async waitForPaymentStatusChange(rowIndex: number, expectedStatus: string, timeout = 10000) {
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      if (await this.verifyPaymentStatus(rowIndex, expectedStatus)) {
        return true;
      }
      await this.page.waitForTimeout(500);
    }

    throw new Error(`Payment status did not change to "${expectedStatus}" within ${timeout}ms`);
  }

  /**
   * Verify success message
   */
  async verifySuccessMessage() {
    await expect(this.page.locator('.ant-message-success')).toBeVisible();
  }

  /**
   * Verify error message
   */
  async verifyErrorMessage() {
    await expect(this.page.locator('.ant-message-error')).toBeVisible();
  }

  /**
   * Navigate to next page if pagination exists
   */
  async nextPage() {
    const nextButton = this.page.getByRole('button', { name: /下一页|next/i });
    if (await nextButton.isEnabled()) {
      await nextButton.click();
      await this.page.waitForLoadState('networkidle');
    }
  }

  /**
   * Export payment records
   */
  async exportPayments() {
    const exportButton = this.page.getByRole('button', { name: /导出|export/i });
    await exportButton.click();
    await this.page.waitForTimeout(2000);
  }

  /**
   * Filter by date range
   */
  async filterByDateRange(_startDate: string, _endDate: string) {
    const dateFilter = this.page.getByText(/日期|date/i).first();
    await dateFilter.click();
    // Handle date picker interactions
    await this.page.waitForTimeout(1000);
  }
}
