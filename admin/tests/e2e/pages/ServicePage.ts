import { Page } from '@playwright/test';

/**
 * Page Object Model for Service Item Management Page
 * Encapsulates all service item management interactions and assertions
 */
export class ServicePage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1, h4').filter({ hasText: /服务项目管理/i });

  // Statistics
  readonly totalServicesStat = this.page.locator('.ant-card').filter({ hasText: /服务项目总数/i });
  readonly activeServicesStat = this.page.locator('.ant-card').filter({ hasText: /已启用/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索服务名称/i);
  readonly gameSelect = this.page.getByPlaceholder(/选择游戏/i);
  readonly serviceTypeSelect = this.page.getByPlaceholder(/服务类型/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增服务/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i });

  // Table action buttons
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.page.getByRole('button', { name: /保存/i });
  private readonly modalCancelButton = this.page.getByRole('button', { name: /取消/i });

  // Form inputs
  readonly itemCodeInput = this.page.getByPlaceholder(/如：ESCORT_SOLO_001/i);
  readonly nameInput = this.page.getByPlaceholder(/如：上分陪玩/i);

  /**
   * Navigate to service page
   */
  async goto() {
    await this.page.goto('/admin/service');
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
   * Get service count from table
   */
  async getServiceCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Search services by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by game
   */
  async filterByGame(gameName: string) {
    await this.gameSelect.click();
    await this.page.getByRole('option', { name: gameName }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by service type
   */
  async filterByServiceType(type: 'solo' | 'team' | 'gift') {
    await this.serviceTypeSelect.click();
    const typeMap = { solo: '单人护航', team: '团队护航', gift: '礼物' };
    await this.page.getByRole('option', { name: typeMap[type] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click add button
   */
  async clickAdd() {
    await this.addButton.click();
  }

  /**
   * Edit service
   */
  async editService(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * Delete service
   */
  async deleteService(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Toggle service status (enable/disable)
   */
  async toggleStatus(rowIndex: number) {
    const switchLocator = this.tableRows.nth(rowIndex).locator('.ant-switch');
    await switchLocator.click();
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
   * Verify service exists in table
   */
  async verifyServiceExists(name: string): Promise<boolean> {
    const count = await this.getServiceCount();
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
  async getStatistics(): Promise<{ total: number; active: number }> {
    const totalText = await this.totalServicesStat.textContent() || '0';
    const activeText = await this.activeServicesStat.textContent() || '0';

    const extractNumber = (text: string): number => {
      const match = text.match(/\d+/);
      return match ? parseInt(match[0], 10) : 0;
    };

    return {
      total: extractNumber(totalText),
      active: extractNumber(activeText),
    };
  }
}
