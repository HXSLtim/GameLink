import { Page } from '@playwright/test';

/**
 * Page Object Model for Activity Management Page
 * Encapsulates all activity management interactions and assertions
 */
export class ActivityPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /活动管理/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics
  readonly totalActivitiesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总活动数/i });
  readonly activeActivitiesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /进行中/i });
  readonly upcomingActivitiesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /未开始/i });

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索活动名称/i);
  readonly typeSelect = this.page.getByText(/活动类型/i).first();
  readonly statusSelect = this.page.getByText(/活动状态/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增活动/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });
  private readonly publishButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /发布/i });
  private readonly unpublishButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /下线/i });
  private readonly rewardsButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /奖励/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to activity page
   */
  async goto() {
    await this.page.goto('/admin/activity');
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
   * Get activity count from table
   */
  async getActivityCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Get statistics
   */
  async getStatistics(): Promise<{ total: number; active: number; upcoming: number }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      total: await getValue(this.totalActivitiesStat),
      active: await getValue(this.activeActivitiesStat),
      upcoming: await getValue(this.upcomingActivitiesStat),
    };
  }

  /**
   * Search activities by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by activity type
   */
  async filterByType(type: 'sign_in' | 'referral' | 'consumption' | 'custom') {
    await this.typeSelect.click();
    const typeMap = {
      sign_in: '签到',
      referral: '推荐',
      consumption: '消费',
      custom: '自定义'
    };
    await this.page.getByRole('option', { name: typeMap[type] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: 'draft' | 'upcoming' | 'active' | 'ended') {
    await this.statusSelect.click();
    const statusMap = {
      draft: '草稿',
      upcoming: '预热',
      active: '进行中',
      ended: '已结束'
    };
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
   * View activity details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
  }

  /**
   * Edit activity
   */
  async editActivity(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * Delete activity
   */
  async deleteActivity(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Publish activity
   */
  async publishActivity(rowIndex: number) {
    await this.publishButton(rowIndex).click();
    // Handle Popconfirm if present
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 3000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Unpublish activity
   */
  async unpublishActivity(rowIndex: number) {
    await this.unpublishButton(rowIndex).click();
    // Handle Popconfirm if present
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 3000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Open rewards management
   */
  async openRewards(rowIndex: number) {
    await this.rewardsButton(rowIndex).click();
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
   * Verify activity exists in table
   */
  async verifyActivityExists(name: string): Promise<boolean> {
    const count = await this.getActivityCount();
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

  /**
   * Get activity status from table
   */
  async getActivityStatus(rowIndex: number): Promise<string> {
    const statusCell = this.tableRows.nth(rowIndex).locator('td').filter({ hasText: /草稿|预热|进行中|已结束/i });
    return await statusCell.textContent() || '';
  }
}
