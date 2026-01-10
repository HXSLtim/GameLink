import { Page, expect, Locator } from '@playwright/test';

/**
 * Page Object Model for Order Management Page
 * Encapsulates all order management interactions and assertions
 */

// Status mapping from English to Chinese
const STATUS_MAP: Record<string, string> = {
  pending: '待确认',
  confirmed: '已确认',
  in_progress: '进行中',
  completed: '已完成',
  canceled: '已取消',
  refunded: '已退款',
  disputed: '争议中',
};

export class OrderManagementPage {
  // Element locators
  public readonly pageTitle: Locator;
  public readonly table: Locator;
  public readonly tableRows: Locator;
  public readonly drawer: Locator;
  public readonly drawerContent: Locator;
  public readonly drawerTitle: Locator;
  public readonly drawerCloseButton: Locator;
  public readonly modal: Locator;
  public readonly modalTitle: Locator;
  public readonly modalOkButton: Locator;

  constructor(public page: Page) {
    // Initialize locators in constructor after page is assigned
    this.pageTitle = this.page.locator('h1').filter({ hasText: /订单管理/ });
    this.table = this.page.locator('.ant-table');
    this.tableRows = this.page.locator('.ant-table-tbody tr.ant-table-row');
    
    // Drawer locators (订单详情使用 Drawer 而不是 Modal)
    this.drawer = this.page.locator('.ant-drawer-content-wrapper');
    this.drawerContent = this.page.locator('.ant-drawer-content');
    this.drawerTitle = this.page.locator('.ant-drawer-title');
    this.drawerCloseButton = this.page.locator('.ant-drawer-close').first();
    
    // Modal locators (用于退款、取消等操作)
    this.modal = this.page.locator('.ant-modal');
    this.modalTitle = this.modal.locator('.ant-modal-title');
    this.modalOkButton = this.modal.getByRole('button', { name: /确定|ok|submit|确认/i });
  }

  // Dynamic locators (methods that return locators)
  viewButton(rowIndex: number): Locator {
    return this.tableRows.nth(rowIndex).locator('button').filter({ hasText: '详情' });
  }

  cancelButton(rowIndex: number): Locator {
    return this.tableRows.nth(rowIndex).locator('button').filter({ hasText: '取消' });
  }

  refundButton(rowIndex: number): Locator {
    return this.tableRows.nth(rowIndex).locator('button').filter({ hasText: '退款' });
  }

  /**
   * Navigate to order management page
   */
  async goto() {
    await this.page.goto('/admin/biz/order');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    // First wait for basic page load
    await this.page.waitForLoadState('domcontentloaded');
    
    // Wait for the layout to be ready (header visible means layout is loaded)
    await this.page.waitForSelector('.ant-layout-header', { state: 'visible', timeout: 15000 });
    
    // Then wait for network to settle
    await this.page.waitForLoadState('networkidle');
    
    // Wait for the table to be visible and have data
    try {
      await this.table.waitFor({ state: 'visible', timeout: 10000 });
      // Wait for table rows to appear (at least one row)
      await this.page.waitForSelector('.ant-table-tbody tr.ant-table-row', { state: 'visible', timeout: 10000 });
    } catch {
      // If table doesn't load, wait a bit more
      await this.page.waitForTimeout(2000);
    }
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
   * View order details (opens Drawer)
   */
  async viewOrderDetails(rowIndex: number) {
    await this.viewButton(rowIndex).click();
    // Wait for drawer to open - use content-wrapper which is the visible part
    await expect(this.drawer).toBeVisible({ timeout: 10000 });
  }

  /**
   * Close drawer by clicking the close button or mask
   */
  async closeDrawer() {
    // Method 1: Click the close button
    const closeBtn = this.page.locator('.ant-drawer-close').first();
    if (await closeBtn.isVisible()) {
      await closeBtn.click();
      await this.page.waitForTimeout(800);
      
      // Check if drawer is closed
      const isVisible = await this.drawer.isVisible();
      if (!isVisible) {
        return;
      }
    }
    
    // Method 2: Click the mask
    const mask = this.page.locator('.ant-drawer-mask');
    if (await mask.isVisible()) {
      await mask.click({ force: true });
      await this.page.waitForTimeout(800);
      
      const isVisible = await this.drawer.isVisible();
      if (!isVisible) {
        return;
      }
    }
    
    // Method 3: Press Escape
    await this.page.keyboard.press('Escape');
    await this.page.waitForTimeout(800);
    
    // Final assertion
    await expect(this.drawer).not.toBeVisible({ timeout: 5000 });
  }

  /**
   * Close modal (for refund, cancel operations)
   */
  async closeModal() {
    const closeButton = this.modal.getByRole('button', { name: /关闭|close|取消|cancel/i });
    await closeButton.click();
  }

  /**
   * Filter orders by status
   */
  async filterByStatus(status: string) {
    // Map English status to Chinese
    const chineseStatus = STATUS_MAP[status] || status;
    
    // The SearchForm has searchFields: orderNo (input), status (select), dateRange (dateRange)
    // Find all selects on the page - the status select should be the first one
    const allSelects = this.page.locator('.ant-select:not(.ant-pagination-options-size-changer)');
    
    // Get the first select (status filter)
    const statusSelect = allSelects.first();
    
    // Click to open dropdown
    await statusSelect.click({ timeout: 10000 });
    await this.page.waitForTimeout(500);
    
    // Select the option from dropdown - Ant Design uses .ant-select-dropdown
    const dropdown = this.page.locator('.ant-select-dropdown:visible');
    await dropdown.waitFor({ state: 'visible', timeout: 5000 });
    
    const option = dropdown.locator('.ant-select-item-option').filter({ hasText: chineseStatus });
    
    if (await option.count() > 0) {
      await option.first().click();
    }
    
    await this.page.waitForTimeout(300);
    
    // Click the search button to apply filter
    const searchButton = this.page.getByRole('button', { name: /搜索/ });
    await searchButton.click();
    
    // Wait for filter to apply and table to reload
    await this.page.waitForTimeout(1000);
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Filter orders by date range
   */
  async filterByDateRange(_startDate: string, _endDate: string) {
    const dateFilter = this.page.getByText(/日期|date/i).first();
    await dateFilter.click();
    // Handle date picker interactions
    // This will depend on the actual date picker implementation
    await this.page.waitForTimeout(1000);
  }

  /**
   * Search for order by order number
   */
  async searchOrder(orderNumber: string) {
    // The search input has placeholder "请输入订单号"
    const searchInput = this.page.getByPlaceholder(/订单号|orderNo/i);
    
    // Check if search input exists
    if (await searchInput.isVisible()) {
      await searchInput.fill(orderNumber);
      await this.page.waitForTimeout(500);
    }
  }

  /**
   * Cancel an order (uses Popconfirm, not Modal)
   */
  async cancelOrder(rowIndex: number, _reason?: string) {
    // Click cancel button - this opens a Popconfirm
    await this.cancelButton(rowIndex).click();
    
    // Wait for Popconfirm to appear
    await this.page.waitForTimeout(300);
    
    // Click the confirm button in Popconfirm
    const confirmBtn = this.page.locator('.ant-popconfirm-buttons button.ant-btn-primary, .ant-popover-buttons button.ant-btn-primary').first();
    if (await confirmBtn.isVisible()) {
      await confirmBtn.click();
    }
    
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
   * Verify order status in a specific row
   * Status column is at index 6 (0-based)
   */
  async verifyOrderStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    // Map English status to Chinese for comparison
    const chineseStatus = STATUS_MAP[expectedStatus] || expectedStatus;
    
    // Get the status cell - it's in column 7 (index 6)
    const row = this.tableRows.nth(rowIndex);
    const statusCell = row.locator('td').nth(6);
    
    // Get the text content of the status tag
    const cellText = await statusCell.textContent() || '';
    
    return cellText.includes(chineseStatus);
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
    // Ant Design pagination uses li.ant-pagination-next for next button
    const nextButton = this.page.locator('.ant-pagination-next:not(.ant-pagination-disabled)');
    
    if (await nextButton.count() > 0 && await nextButton.isEnabled()) {
      await nextButton.click();
      await this.page.waitForLoadState('networkidle');
      await this.page.waitForTimeout(500);
    }
  }

  /**
   * Check if there are orders with specific status
   */
  async hasOrdersWithStatus(status: string): Promise<boolean> {
    await this.filterByStatus(status);
    await this.page.waitForTimeout(1000);
    const count = await this.getOrderCount();
    return count > 0;
  }

  /**
   * Check if cancel button exists for a row
   */
  async hasCancelButton(rowIndex: number): Promise<boolean> {
    const btn = this.cancelButton(rowIndex);
    return await btn.count() > 0 && await btn.isVisible();
  }

  /**
   * Check if refund button exists for a row
   */
  async hasRefundButton(rowIndex: number): Promise<boolean> {
    const btn = this.refundButton(rowIndex);
    return await btn.count() > 0 && await btn.isVisible();
  }
}
