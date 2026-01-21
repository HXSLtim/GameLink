import { Page } from '@playwright/test';

/**
 * Page Object Model for Routing Rule Management Page
 * Encapsulates all routing rule management interactions and assertions
 */
export class RoutingPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /支付路由规则/i });

  // Statistics
  readonly totalRulesStat = this.page.locator('.ant-card').filter({ hasText: /规则总数/i });
  readonly activeRulesStat = this.page.locator('.ant-card').filter({ hasText: /启用规则/i });
  readonly inactiveRulesStat = this.page.locator('.ant-card').filter({ hasText: /禁用规则/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/规则名称/i);
  readonly statusSelect = this.page.getByText(/状态/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增规则/i });
  readonly testButton = this.page.getByRole('button', { name: /测试规则/i });
  readonly exportButton = this.page.getByRole('button', { name: /导出数据/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i });

  // Table action buttons
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i }).or(
    this.tableRows.nth(rowIndex).locator('button').filter({ hasText: /编辑/i })
  );
  private readonly historyButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /历史/i }).or(
    this.tableRows.nth(rowIndex).locator('button').filter({ hasText: /历史/i })
  );
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i }).or(
    this.tableRows.nth(rowIndex).locator('button').filter({ hasText: /删除/i })
  );

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to routing page
   */
  async goto() {
    await this.page.goto('/admin/routing');
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
   * Get rule count from table
   */
  async getRuleCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Search rules by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: 'active' | 'inactive') {
    // Click status select and choose option
    const statusSelect = this.page.locator('.ant-select').filter({ hasText: /状态/ }).first();
    await statusSelect.click();
    const statusMap = { active: '启用', inactive: '禁用' };
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
   * Edit rule
   */
  async editRule(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * View rule history
   */
  async viewHistory(rowIndex: number) {
    await this.historyButton(rowIndex).click();
  }

  /**
   * Delete rule
   */
  async deleteRule(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Toggle rule status by clicking on status tag
   */
  async toggleStatus(rowIndex: number) {
    const statusTag = this.tableRows.nth(rowIndex).locator('.ant-tag');
    await statusTag.click();
  }

  /**
   * Click test rules button
   */
  async clickTestRules() {
    await this.testButton.click();
  }

  /**
   * Click export button
   */
  async clickExport() {
    await this.exportButton.click();
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
   * Verify rule exists in table
   */
  async verifyRuleExists(name: string): Promise<boolean> {
    const count = await this.getRuleCount();
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
   * Get statistics values
   */
  async getStatistics(): Promise<{ total: number; active: number; inactive: number }> {
    const totalText = await this.totalRulesStat.textContent() || '0';
    const activeText = await this.activeRulesStat.textContent() || '0';
    const inactiveText = await this.inactiveRulesStat.textContent() || '0';

    const extractNumber = (text: string): number => {
      const match = text.match(/\d+/);
      return match ? parseInt(match[0], 10) : 0;
    };

    return {
      total: extractNumber(totalText),
      active: extractNumber(activeText),
      inactive: extractNumber(inactiveText),
    };
  }
}
