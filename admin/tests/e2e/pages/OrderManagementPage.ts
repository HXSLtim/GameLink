import { Page, expect } from '@playwright/test';

/**
 * Page Object Model for Order Management Page
 * Encapsulates all order management interactions and assertions
 */
export class OrderManagementPage {
  constructor(private page: Page) {}

  // Element locators
  private readonly pageTitle = this.page.getByRole('heading', { name: /订单管理|order management/i });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');
  private readonly viewButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /查看|view/i });
  private readonly cancelButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /取消|cancel/i });
  private readonly refundButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /退款|refund/i });

  // Status filter locators
  private readonly statusFilter = this.page.getByText(/状态|status/i).first();
  private readonly dateFilter = this.page.getByText(/日期|date/i).first();

  // Modal locators
  private readonly modal = this.page.locator('.ant-modal');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定|ok|submit/i });

  /**
   * Navigate to order management page
   */
  async goto() {
    await this.page.goto('/admin/Order');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await expect(this.pageTitle).toBeVisible({ timeout: 10000 });
  }

  /**
   * Get order count from table
   */
  async getOrderCount(): Promise<number> {
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
   * View order details
   */
  async viewOrderDetails(rowIndex: number) {
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
   * Filter orders by status
   */
  async filterByStatus(status: string) {
    await this.statusFilter.click();
    await this.page.getByRole('option', { name: status }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter orders by date range
   */
  async filterByDateRange(_startDate: string, _endDate: string) {
    await this.dateFilter.click();
    // Handle date picker interactions
    // This will depend on the actual date picker implementation
    await this.page.waitForTimeout(1000);
  }

  /**
   * Search for order by order number
   */
  async searchOrder(orderNumber: string) {
    const searchInput = this.page.getByPlaceholder(/搜索|search/i);
    await searchInput.fill(orderNumber);
    await this.page.waitForTimeout(500);
  }

  /**
   * Cancel an order
   */
  async cancelOrder(rowIndex: number, reason?: string) {
    await this.cancelButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    if (reason) {
      const reasonInput = this.modal.getByLabel(/原因|reason/i);
      await reasonInput.fill(reason);
    }

    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Refund an order
   */
  async refundOrder(rowIndex: number, data: { reason: string; amount?: number }) {
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
   * Batch cancel orders
   */
  async batchCancelOrders(orderNumbers: string[], reason?: string) {
    // Select checkboxes
    for (const orderNumber of orderNumbers) {
      const checkbox = this.page.getByText(orderNumber).locator('..').getByRole('checkbox');
      await checkbox.check();
    }

    // Click batch cancel button
    const batchCancelButton = this.page.getByRole('button', { name: /批量取消/i });
    await batchCancelButton.click();

    // Fill reason if needed
    if (reason) {
      const reasonInput = this.modal.getByLabel(/原因|reason/i);
      await reasonInput.fill(reason);
    }

    // Confirm
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Batch complete orders
   */
  async batchCompleteOrders(orderNumbers: string[]) {
    // Select checkboxes
    for (const orderNumber of orderNumbers) {
      const checkbox = this.page.getByText(orderNumber).locator('..').getByRole('checkbox');
      await checkbox.check();
    }

    // Click batch complete button
    const batchCompleteButton = this.page.getByRole('button', { name: /批量完成/i });
    await batchCompleteButton.click();

    // Confirm
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify order status
   */
  async verifyOrderStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    const cellText = await this.getCellText(rowIndex, 4); // Assuming status is in column 4
    return cellText.includes(expectedStatus);
  }

  /**
   * Wait for order status to change
   */
  async waitForOrderStatusChange(rowIndex: number, expectedStatus: string, timeout = 10000) {
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      if (await this.verifyOrderStatus(rowIndex, expectedStatus)) {
        return true;
      }
      await this.page.waitForTimeout(500);
    }

    throw new Error(`Order status did not change to "${expectedStatus}" within ${timeout}ms`);
  }

  /**
   * Verify success message
   */
  async verifySuccessMessage() {
    await expect(this.page.locator('.ant-message-success')).toBeVisible();
  }

  /**
   * Export order list
   */
  async exportOrders() {
    const exportButton = this.page.getByRole('button', { name: /导出|export/i });
    await exportButton.click();
    await this.page.waitForTimeout(2000);
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
}
