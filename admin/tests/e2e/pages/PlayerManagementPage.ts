import { Page, expect } from '@playwright/test';

/**
 * Page Object Model for Player Management Page
 * Encapsulates all player management interactions and assertions
 */
export class PlayerManagementPage {
  constructor(private page: Page) {}

  // Element locators - PageContainer uses h1 for title
  private readonly pageTitle = this.page.locator('h1').filter({ hasText: /陪玩师管理/ });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');
  private readonly editButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /编辑|edit/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除|delete/i });
  private readonly viewButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /查看|view/i });

  // Status buttons
  private readonly approveButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /通过|approve|verify/i });
  private readonly rejectButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /拒绝|reject/i });

  // Filter locators
  private readonly statusFilter = this.page.getByText(/状态|status/i).first();
  private readonly searchInput = this.page.getByPlaceholder(/搜索|search/i);

  // Modal locators
  private readonly modal = this.page.locator('.ant-modal');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定|ok|submit/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消|cancel/i });

  /**
   * Navigate to player management page
   */
  async goto() {
    await this.page.goto('/admin/biz/player');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    // Wait for either the page title or the table to be visible
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.table.waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Get player count from table
   */
  async getPlayerCount(): Promise<number> {
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
   * View player details
   */
  async viewPlayerDetails(rowIndex: number) {
    await this.viewButton(rowIndex).click();
    await expect(this.modal).toBeVisible();
  }

  /**
   * Edit player
   */
  async editPlayer(rowIndex: number, updatedData: {
    nickname?: string;
    bio?: string;
    rank?: string;
    hourlyRateCents?: number;
    mainGameId?: number;
  }) {
    await this.editButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    if (updatedData.nickname) {
      const nicknameInput = this.modal.getByLabel(/昵称|nickname/i);
      await nicknameInput.fill(updatedData.nickname);
    }

    if (updatedData.bio) {
      const bioInput = this.modal.getByLabel(/简介|bio/i);
      await bioInput.fill(updatedData.bio);
    }

    if (updatedData.rank) {
      const rankSelect = this.modal.getByLabel(/段位|rank/i);
      await rankSelect.click();
      await this.page.getByRole('option', { name: updatedData.rank }).click();
    }

    if (updatedData.hourlyRateCents) {
      const rateInput = this.modal.getByLabel(/时薪|hourly rate/i);
      await rateInput.fill(updatedData.hourlyRateCents.toString());
    }

    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Delete player
   */
  async deletePlayer(rowIndex: number) {
    await this.deleteButton(rowIndex).click();
    await expect(this.modal).toBeVisible();
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Approve player verification
   */
  async approvePlayer(rowIndex: number, remark?: string) {
    await this.approveButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    if (remark) {
      const remarkInput = this.modal.getByLabel(/备注|remark/i);
      await remarkInput.fill(remark);
    }

    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Reject player verification
   */
  async rejectPlayer(rowIndex: number, reason: string) {
    await this.rejectButton(rowIndex).click();
    await expect(this.modal).toBeVisible();

    const reasonInput = this.modal.getByLabel(/原因|reason/i);
    await reasonInput.fill(reason);

    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter players by status
   */
  async filterByStatus(status: string) {
    await this.statusFilter.click();
    await this.page.getByRole('option', { name: status }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Search for player
   */
  async searchPlayer(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Verify player status
   */
  async verifyPlayerStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    const cellText = await this.getCellText(rowIndex, 3); // Assuming status is in column 3
    return cellText.includes(expectedStatus);
  }

  /**
   * Verify player certification status
   */
  async verifyCertificationStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    const cellText = await this.getCellText(rowIndex, 4); // Assuming certification status is in column 4
    return cellText.includes(expectedStatus);
  }

  /**
   * Wait for player status to change
   */
  async waitForPlayerStatusChange(rowIndex: number, expectedStatus: string, timeout = 10000) {
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      if (await this.verifyPlayerStatus(rowIndex, expectedStatus)) {
        return true;
      }
      await this.page.waitForTimeout(500);
    }

    throw new Error(`Player status did not change to "${expectedStatus}" within ${timeout}ms`);
  }

  /**
   * Batch update player status
   */
  async batchUpdateStatus(playerNicknames: string[], status: string) {
    // Select checkboxes
    for (const nickname of playerNicknames) {
      const checkbox = this.page.getByText(nickname).locator('..').getByRole('checkbox');
      await checkbox.check();
    }

    // Click batch status button
    const batchStatusButton = this.page.getByRole('button', { name: /批量状态|batch status/i });
    await batchStatusButton.click();

    // Select status
    await this.page.getByRole('option', { name: status }).click();

    // Confirm
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Batch delete players
   */
  async batchDeletePlayers(playerNicknames: string[]) {
    // Select checkboxes
    for (const nickname of playerNicknames) {
      const checkbox = this.page.getByText(nickname).locator('..').getByRole('checkbox');
      await checkbox.check();
    }

    // Click batch delete button
    const batchDeleteButton = this.page.getByRole('button', { name: /批量删除|batch delete/i });
    await batchDeleteButton.click();

    // Confirm deletion
    await this.modalOkButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify player exists in table
   */
  async verifyPlayerExists(nickname: string): Promise<boolean> {
    const playerCount = await this.getPlayerCount();
    for (let i = 0; i < playerCount; i++) {
      const cellText = await this.getCellText(i, 1); // Assuming nickname is in column 1
      if (cellText.includes(nickname)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Verify success message
   */
  async verifySuccessMessage() {
    await expect(this.page.locator('.ant-message-success')).toBeVisible();
  }

  /**
   * Verify error message
   */
  async verifyErrorMessage() {
    await expect(this.page.locator('.ant-message-error')).toBeVisible();
  }

  /**
   * Export player list
   */
  async exportPlayers() {
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
