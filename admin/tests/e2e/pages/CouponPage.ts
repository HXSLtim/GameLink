import { Page } from '@playwright/test';

/**
 * Page Object Model for Coupon Management Page
 * Encapsulates all coupon management interactions and assertions
 */
export class CouponPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /优惠券管理/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索优惠券名称或代码/i);
  readonly typeSelect = this.page.getByText(/优惠券类型/i).first();
  readonly scopeSelect = this.page.getByText(/适用范围/i).first();
  readonly statusSelect = this.page.getByText(/状态/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增优惠券/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });
  readonly exportButton = this.page.getByRole('button', { name: /导出/i });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to coupon page
   */
  async goto() {
    await this.page.goto('/admin/coupon');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.table.waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Get coupon count from table
   */
  async getCouponCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Search coupons by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by coupon type
   */
  async filterByType(type: 'discount' | 'fixed' | 'gift') {
    await this.typeSelect.click();
    const typeMap = { discount: '折扣券', fixed: '代金券', gift: '礼品券' };
    await this.page.getByRole('option', { name: typeMap[type] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by scope
   */
  async filterByScope(scope: 'all' | 'category' | 'item') {
    await this.scopeSelect.click();
    const scopeMap = { all: '全场通用', category: '指定分类', item: '指定商品' };
    await this.page.getByRole('option', { name: scopeMap[scope] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: 'active' | 'inactive' | 'expired') {
    await this.statusSelect.click();
    const statusMap = { active: '进行中', inactive: '未开始', expired: '已结束' };
    await this.page.getByRole('option', { name: statusMap[status] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click add button
   */
  async clickAdd() {
    await this.addButton.click();
  }

  /**
   * View coupon details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
  }

  /**
   * Edit coupon
   */
  async editCoupon(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * Delete coupon
   */
  async deleteCoupon(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Refresh the list
   */
  async refresh() {
    await this.refreshButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Get cell text from table
   */
  async getCellText(rowIndex: number, cellIndex: number): Promise<string> {
    const cell = this.tableRows.nth(rowIndex).locator('td').nth(cellIndex);
    return await cell.textContent() || '';
  }

  /**
   * Verify modal is visible
   */
  async isModalVisible(): Promise<boolean> {
    return await this.modal.isVisible();
  }

  /**
   * Close modal
   */
  async closeModal() {
    await this.modalCancelButton.click();
  }

  /**
   * Submit modal form
   */
  async submitModal() {
    await this.modalOkButton.click();
  }

  /**
   * Verify coupon exists in table
   */
  async verifyCouponExists(name: string): Promise<boolean> {
    const count = await this.getCouponCount();
    for (let i = 0; i < count; i++) {
      const text = await this.getCellText(i, 1); // Name is typically in column 1
      if (text.includes(name)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Clear search filters
   */
  async clearFilters() {
    await this.searchInput.fill('');
    await this.page.waitForTimeout(500);
  }
}
