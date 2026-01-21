import { Page } from '@playwright/test';

/**
 * Page Object Model for Settlement Management Page
 * Encapsulates all settlement management interactions and assertions
 */
export class SettlementPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /结算管理|陪玩师归属/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics (if present)
  readonly totalPlayersStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总陪玩师/i });

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索/i);
  readonly companySelect = this.page.getByText(/所属公司/i).first();

  // Toolbar buttons
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly assignButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /分配/i });
  private readonly transferButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /转移/i });
  private readonly historyButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /历史/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  // Drawer
  private readonly drawer = this.page.locator('.ant-drawer');
  private readonly drawerCloseButton = this.drawer.locator('.ant-drawer-close');

  /**
   * Navigate to settlement page
   */
  async goto() {
    await this.page.goto('/admin/settlement/players');
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
   * Get player count from table
   */
  async getPlayerCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Search by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by company
   */
  async filterByCompany(company: string) {
    await this.companySelect.click();
    await this.page.getByRole('option', { name: company }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * View player details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
  }

  /**
   * Assign player to company
   */
  async assignPlayer(rowIndex: number) {
    await this.assignButton(rowIndex).click();
  }

  /**
   * Transfer player
   */
  async transferPlayer(rowIndex: number) {
    await this.transferButton(rowIndex).click();
  }

  /**
   * View player history
   */
  async viewHistory(rowIndex: number) {
    await this.historyButton(rowIndex).click();
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
   * Clear search filters
   */
  async clearFilters() {
    await this.searchInput.fill('');
    await this.page.waitForTimeout(500);
  }
}
