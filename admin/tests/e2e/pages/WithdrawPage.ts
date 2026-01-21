import { Page } from '@playwright/test';

/**
 * Page Object Model for Withdraw Management Page
 * Encapsulates all withdraw management interactions and assertions
 */
export class WithdrawPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /提现管理/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics
  readonly pendingStat = this.page.locator('.ant-statistic-title').filter({ hasText: /待审核/i });
  readonly approvedStat = this.page.locator('.ant-statistic-title').filter({ hasText: /已批准/i });
  readonly completedStat = this.page.locator('.ant-statistic-title').filter({ hasText: /已完成/i });

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索用户名/i);
  readonly statusSelect = this.page.getByText(/状态/i).first();

  // Toolbar buttons
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly approveButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /批准/i });
  private readonly rejectButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /拒绝/i });
  private readonly completeButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /完成/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  // Drawer
  private readonly drawer = this.page.locator('.ant-drawer');
  private readonly drawerCloseButton = this.drawer.locator('.ant-drawer-close');

  /**
   * Navigate to withdraw page
   */
  async goto() {
    await this.page.goto('/admin/withdraw');
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
   * Get withdraw count from table
   */
  async getWithdrawCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Get statistics
   */
  async getStatistics(): Promise<{ pending: number; approved: number; completed: number }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      pending: await getValue(this.pendingStat),
      approved: await getValue(this.approvedStat),
      completed: await getValue(this.completedStat),
    };
  }

  /**
   * Search by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: 'pending' | 'approved' | 'rejected' | 'completed') {
    await this.statusSelect.click();
    const statusMap = { pending: '待审核', approved: '已批准', rejected: '已拒绝', completed: '已完成' };
    await this.page.getByRole('option', { name: statusMap[status] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * View withdraw details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
  }

  /**
   * Approve withdraw
   */
  async approveWithdraw(rowIndex: number) {
    await this.approveButton(rowIndex).click();
    // Handle confirmation if present
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 3000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Reject withdraw
   */
  async rejectWithdraw(rowIndex: number) {
    await this.rejectButton(rowIndex).click();
    // Handle confirmation/reason modal
    await this.page.waitForTimeout(500);
  }

  /**
   * Complete withdraw
   */
  async completeWithdraw(rowIndex: number) {
    await this.completeButton(rowIndex).click();
    // Handle confirmation if present
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 3000 })) {
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
   * Verify drawer is visible
   */
  async isDrawerVisible(): Promise<boolean> {
    return await this.drawer.isVisible();
  }

  /**
   * Verify modal is visible
   */
  async isModalVisible(): Promise<boolean> {
    return await this.modal.isVisible();
  }

  /**
   * Close drawer
   */
  async closeDrawer() {
    await this.drawerCloseButton.click();
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
   * Get withdraw status from table
   */
  async getWithdrawStatus(rowIndex: number): Promise<string> {
    const statusCell = this.tableRows.nth(rowIndex).locator('td').filter({ hasText: /待审核|已批准|已拒绝|已完成/i });
    return await statusCell.textContent() || '';
  }

  /**
   * Clear search filters
   */
  async clearFilters() {
    await this.searchInput.fill('');
    await this.page.waitForTimeout(500);
  }
}
