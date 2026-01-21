import { Page } from '@playwright/test';

/**
 * Page Object Model for Team Management Page
 * Encapsulates all team management interactions and assertions
 */
export class TeamPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /团队管理/i });

  // Table
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics
  readonly totalTeamsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总团队数/i });
  readonly activeTeamsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /活跃团队/i });
  readonly totalMembersStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总成员数/i });

  // Search fields
  readonly searchInput = this.page.getByPlaceholder(/搜索团队名称/i);
  readonly statusSelect = this.page.getByText(/团队状态/i).first();

  // Toolbar buttons
  readonly addButton = this.page.getByRole('button', { name: /新增团队/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  // Drawer
  private readonly drawer = this.page.locator('.ant-drawer');
  private readonly drawerCloseButton = this.drawer.locator('.ant-drawer-close');

  /**
   * Navigate to team page
   */
  async goto() {
    await this.page.goto('/admin/team');
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
   * Get team count from table
   */
  async getTeamCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Get statistics
   */
  async getStatistics(): Promise<{ totalTeams: number; activeTeams: number; totalMembers: number }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      totalTeams: await getValue(this.totalTeamsStat),
      activeTeams: await getValue(this.activeTeamsStat),
      totalMembers: await getValue(this.totalMembersStat),
    };
  }

  /**
   * Search teams by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: 'active' | 'busy' | 'inactive') {
    await this.statusSelect.click();
    const statusMap = { active: '活跃', busy: '接单中', inactive: '不活跃' };
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
   * View team details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
  }

  /**
   * Edit team
   */
  async editTeam(rowIndex: number) {
    await this.editButton(rowIndex).click();
  }

  /**
   * Delete team
   */
  async deleteTeam(rowIndex: number) {
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
   * Verify drawer is visible
   */
  async isDrawerVisible(): Promise<boolean> {
    return await this.drawer.isVisible();
  }

  /**
   * Close modal
   */
  async closeModal() {
    await this.modalCancelButton.click();
  }

  /**
   * Close drawer
   */
  async closeDrawer() {
    await this.drawerCloseButton.click();
  }

  /**
   * Submit modal form
   */
  async submitModal() {
    await this.modalOkButton.click();
  }

  /**
   * Verify team exists in table
   */
  async verifyTeamExists(name: string): Promise<boolean> {
    const count = await this.getTeamCount();
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
