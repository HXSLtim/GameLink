import { Page } from '@playwright/test';

/**
 * Page Object Model for Game Management Page
 * Encapsulates all game management interactions and assertions
 */
export class GamePage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /游戏管理/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索游戏名称/i);
  readonly categorySelect = this.page.getByText(/游戏分类/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增游戏/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });
  readonly importButton = this.page.getByRole('button', { name: /导入/i });

  // Table action buttons
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to game page
   */
  async goto() {
    await this.page.goto('/admin/game');
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
   * Get game count from table
   */
  async getGameCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Search games by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by category
   */
  async filterByCategory(category: 'moba' | 'fps' | 'rpg' | 'card' | 'casual' | 'other') {
    await this.categorySelect.click();
    const categoryMap = { moba: 'MOBA', fps: '射击', rpg: 'RPG', card: '卡牌', casual: '休闲', other: '其他' };
    await this.page.getByRole('option', { name: categoryMap[category] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click add button
   */
  async clickAdd() {
    await this.addButton.click();
  }

  /**
   * Edit game
   */
  async editGame(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * Delete game
   */
  async deleteGame(rowIndex: number) {
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
   * Verify game exists in table
   */
  async verifyGameExists(name: string): Promise<boolean> {
    const count = await this.getGameCount();
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
