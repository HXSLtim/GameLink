import { Page, expect, Locator } from '@playwright/test';

/**
 * Page Object Model for Dispute Management Page
 * Encapsulates all dispute management page interactions and assertions
 */
export class DisputeManagementPage {
  readonly page: Page;
  readonly pageTitle: Locator;
  readonly table: Locator;
  readonly drawer: Locator;
  readonly drawerTitle: Locator;
  readonly assignModal: Locator;
  readonly resolveModal: Locator;

  // Statistics cards
  readonly pendingCard: Locator;
  readonly assignedCard: Locator;
  readonly mediatingCard: Locator;
  readonly resolvedCard: Locator;
  readonly rejectedCard: Locator;
  readonly slaBreachedCard: Locator;

  // Search fields
  readonly orderNoInput: Locator;
  readonly statusSelect: Locator;
  readonly initiatorTypeSelect: Locator;

  constructor(page: Page) {
    this.page = page;
    this.pageTitle = page.getByText('纠纷管理');
    this.table = page.locator('.ant-table');
    this.drawer = page.locator('.ant-drawer');
    this.drawerTitle = page.locator('.ant-drawer-title');
    this.assignModal = page.locator('.ant-modal').filter({ hasText: /分配客服|指派/ });
    this.resolveModal = page.locator('.ant-modal').filter({ hasText: /处理纠纷|解决/ });

    // Statistics cards
    this.pendingCard = page.locator('.ant-card').filter({ hasText: '待处理' });
    this.assignedCard = page.locator('.ant-card').filter({ hasText: '已指派' });
    this.mediatingCard = page.locator('.ant-card').filter({ hasText: '调解中' });
    this.resolvedCard = page.locator('.ant-card').filter({ hasText: '已解决' });
    this.rejectedCard = page.locator('.ant-card').filter({ hasText: '已驳回' });
    this.slaBreachedCard = page.locator('.ant-card').filter({ hasText: 'SLA超时' });

    // Search fields
    this.orderNoInput = page.getByPlaceholder('请输入订单号');
    this.statusSelect = page.locator('.ant-select').filter({ hasText: /状态/ });
    this.initiatorTypeSelect = page.locator('.ant-select').filter({ hasText: /发起人类型/ });
  }

  /**
   * Navigate to dispute management page
   */
  async goto() {
    await this.page.goto('/admin/disputes');
    await this.page.waitForLoadState('networkidle');
    await expect(this.pageTitle).toBeVisible({ timeout: 10000 });
  }

  /**
   * Get the count of disputes in the table
   */
  async getDisputeCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    const rows = this.table.locator('tbody tr').filter({ hasNot: this.page.locator('.ant-table-placeholder') });
    return await rows.count();
  }

  /**
   * Get cell text from table
   */
  async getCellText(rowIndex: number, colIndex: number): Promise<string> {
    const cell = this.table.locator(`tbody tr:nth-child(${rowIndex + 1}) td:nth-child(${colIndex + 1})`);
    return await cell.textContent() || '';
  }

  /**
   * View dispute details
   */
  async viewDisputeDetails(rowIndex: number) {
    const detailButton = this.table.locator(`tbody tr:nth-child(${rowIndex + 1})`).getByRole('button', { name: /详情|查看/ });
    await detailButton.click();
    await expect(this.drawer).toBeVisible();
  }

  /**
   * Close drawer
   */
  async closeDrawer() {
    const closeButton = this.drawer.locator('.ant-drawer-close');
    await closeButton.click();
    await expect(this.drawer).not.toBeVisible();
  }

  /**
   * Get assign button for a row
   */
  assignButton(rowIndex: number): Locator {
    return this.table.locator(`tbody tr:nth-child(${rowIndex + 1})`).getByRole('button', { name: /分配|指派/ });
  }

  /**
   * Check if assign button exists for a row
   */
  async hasAssignButton(rowIndex: number): Promise<boolean> {
    return await this.assignButton(rowIndex).count() > 0;
  }

  /**
   * Get resolve button for a row
   */
  resolveButton(rowIndex: number): Locator {
    return this.table.locator(`tbody tr:nth-child(${rowIndex + 1})`).getByRole('button', { name: /处理|解决/ });
  }

  /**
   * Check if resolve button exists for a row
   */
  async hasResolveButton(rowIndex: number): Promise<boolean> {
    return await this.resolveButton(rowIndex).count() > 0;
  }

  /**
   * Get rollback button for a row
   */
  rollbackButton(rowIndex: number): Locator {
    return this.table.locator(`tbody tr:nth-child(${rowIndex + 1})`).getByRole('button', { name: /回滚/ });
  }

  /**
   * Check if rollback button exists for a row
   */
  async hasRollbackButton(rowIndex: number): Promise<boolean> {
    return await this.rollbackButton(rowIndex).count() > 0;
  }

  /**
   * Search by order number
   */
  async searchByOrderNo(orderNo: string) {
    await this.orderNoInput.fill(orderNo);
    await this.page.keyboard.press('Enter');
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by status
   */
  async filterByStatus(status: string) {
    // Find the status select in the search form
    const statusSelect = this.page.locator('.ant-form-item').filter({ hasText: '状态' }).locator('.ant-select');
    await statusSelect.click();
    await this.page.waitForTimeout(300);
    
    // Select the option
    const option = this.page.locator('.ant-select-dropdown:visible').getByText(status, { exact: false });
    await option.click();
    await this.page.waitForTimeout(500);
    
    // Click search button
    const searchButton = this.page.getByRole('button', { name: /搜索|查询/ });
    if (await searchButton.count() > 0) {
      await searchButton.click();
    }
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by initiator type
   */
  async filterByInitiatorType(type: 'user' | 'player') {
    const typeLabel = type === 'user' ? '用户' : '陪玩师';
    
    // Find the initiator type select in the search form
    const typeSelect = this.page.locator('.ant-form-item').filter({ hasText: '发起人类型' }).locator('.ant-select');
    await typeSelect.click();
    await this.page.waitForTimeout(300);
    
    // Select the option
    const option = this.page.locator('.ant-select-dropdown:visible').getByText(typeLabel);
    await option.click();
    await this.page.waitForTimeout(500);
    
    // Click search button
    const searchButton = this.page.getByRole('button', { name: /搜索|查询/ });
    if (await searchButton.count() > 0) {
      await searchButton.click();
    }
    await this.page.waitForTimeout(1000);
  }

  /**
   * Assign dispute to customer service
   */
  async assignDispute(rowIndex: number, assignedServiceId: number, originalServiceId?: number) {
    await this.assignButton(rowIndex).click();
    await expect(this.assignModal).toBeVisible();

    // Fill in assigned service ID
    const assignedInput = this.assignModal.locator('input').first();
    await assignedInput.fill(assignedServiceId.toString());

    // Fill in original service ID if provided
    if (originalServiceId) {
      const originalInput = this.assignModal.locator('input').nth(1);
      await originalInput.fill(originalServiceId.toString());
    }

    // Confirm
    const confirmButton = this.assignModal.getByRole('button', { name: /确定|确认/ });
    await confirmButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Resolve dispute
   */
  async resolveDispute(rowIndex: number, resolution: string, remark: string) {
    await this.resolveButton(rowIndex).click();
    await expect(this.resolveModal).toBeVisible();

    // Select resolution type
    const resolutionSelect = this.resolveModal.locator('.ant-select');
    await resolutionSelect.click();
    await this.page.waitForTimeout(300);
    
    const option = this.page.locator('.ant-select-dropdown:visible').getByText(resolution, { exact: false });
    await option.click();

    // Fill in remark
    const remarkInput = this.resolveModal.locator('textarea');
    await remarkInput.fill(remark);

    // Confirm
    const confirmButton = this.resolveModal.getByRole('button', { name: /确定|确认/ });
    await confirmButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Rollback dispute assignment
   */
  async rollbackAssignment(rowIndex: number) {
    await this.rollbackButton(rowIndex).click();
    
    // Confirm in the modal
    const confirmButton = this.page.locator('.ant-modal-confirm').getByRole('button', { name: /确认|确定/ });
    await confirmButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify success message
   */
  async verifySuccessMessage() {
    const successMessage = this.page.locator('.ant-message-success');
    await expect(successMessage).toBeVisible({ timeout: 5000 });
  }

  /**
   * Verify error message
   */
  async verifyErrorMessage() {
    const errorMessage = this.page.locator('.ant-message-error');
    await expect(errorMessage).toBeVisible({ timeout: 5000 });
  }

  /**
   * Get statistics value
   */
  async getStatValue(statName: string): Promise<string> {
    const card = this.page.locator('.ant-card').filter({ hasText: statName });
    const value = card.locator('.ant-statistic-content-value');
    return await value.textContent() || '0';
  }

  /**
   * Check if SLA warning is displayed
   */
  async hasSlaWarning(): Promise<boolean> {
    const warning = this.page.locator('.ant-alert-warning').filter({ hasText: /SLA|超时/ });
    return await warning.count() > 0;
  }

  /**
   * Verify dispute status in row
   */
  async verifyDisputeStatus(rowIndex: number, expectedStatus: string): Promise<boolean> {
    const statusCell = this.table.locator(`tbody tr:nth-child(${rowIndex + 1})`).locator('.ant-tag');
    const statusText = await statusCell.textContent();
    return statusText?.includes(expectedStatus) || false;
  }

  /**
   * Export disputes data
   */
  async exportData() {
    const exportButton = this.page.getByRole('button', { name: /导出/ });
    await exportButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Refresh data
   */
  async refresh() {
    const refreshButton = this.page.getByRole('button', { name: /刷新/ });
    if (await refreshButton.count() > 0) {
      await refreshButton.click();
      await this.page.waitForTimeout(1000);
    }
  }

  /**
   * Navigate to next page
   */
  async nextPage() {
    const nextButton = this.page.locator('.ant-pagination-next:not(.ant-pagination-disabled)');
    if (await nextButton.count() > 0) {
      await nextButton.click();
      await this.page.waitForTimeout(1000);
    }
  }

  /**
   * Navigate to previous page
   */
  async prevPage() {
    const prevButton = this.page.locator('.ant-pagination-prev:not(.ant-pagination-disabled)');
    if (await prevButton.count() > 0) {
      await prevButton.click();
      await this.page.waitForTimeout(1000);
    }
  }
}
