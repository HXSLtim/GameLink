import { Page } from '@playwright/test';

/**
 * Page Object Model for Chat Rooms Management Page
 * Encapsulates all chat rooms management interactions and assertions
 */
export class ChatRoomsPage {
  constructor(private page: Page) {}

  // Element locators
  private readonly pageTitle = this.page.locator('h1').filter({ hasText: /聊天室管理/i });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics cards
  readonly totalConversationsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总会话数/i });
  readonly activeConversationsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /活跃会话/i });
  readonly totalUsersStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总用户数/i });
  readonly totalMessagesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总消息数/i });

  // Quick action buttons
  readonly activeOnlyButton = this.page.getByRole('button', { name: /仅活跃会话/i });
  readonly closedOnlyButton = this.page.getByRole('button', { name: /已关闭会话/i });
  readonly todayCreatedButton = this.page.getByRole('button', { name: /今天创建的会话/i });
  readonly todayActiveButton = this.page.getByRole('button', { name: /今天的活跃会话/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // Search fields
  private readonly searchInput = this.page.getByPlaceholder(/搜索用户名、陪玩师名、订单号/i);
  private readonly typeSelect = this.page.getByText(/会话类型/i).first();
  private readonly statusSelect = this.page.getByText(/会话状态/i).first();
  private readonly orderIdInput = this.page.getByPlaceholder(/输入订单ID/i);
  private readonly searchButton = this.page.getByRole('button', { name: /查询/i });

  // Toolbar buttons
  readonly exportButton = this.page.getByRole('button', { name: /导出数据/i });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly closeButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /关闭/i });
  private readonly reopenButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /打开/i });

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  private readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalCloseButton = this.modal.locator('.ant-modal-close');

  // Modal action buttons
  private readonly modalCloseConversationButton = this.modal.getByRole('button', { name: /关闭会话/i });
  private readonly modalReopenConversationButton = this.modal.getByRole('button', { name: /重新打开会话/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to chat rooms page
   */
  async goto() {
    await this.page.goto('/admin/chat/rooms');
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
   * Get conversation count from table
   */
  async getConversationCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Get conversation statistics
   */
  async getStatistics(): Promise<{
    totalConversations: number;
    activeConversations: number;
    totalUsers: number;
    totalMessages: number;
  }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      totalConversations: await getValue(this.totalConversationsStat),
      activeConversations: await getValue(this.activeConversationsStat),
      totalUsers: await getValue(this.totalUsersStat),
      totalMessages: await getValue(this.totalMessagesStat),
    };
  }

  /**
   * Search conversations by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by conversation status
   */
  async filterByStatus(status: 'active' | 'closed') {
    await this.statusSelect.click();
    const optionText = status === 'active' ? '活跃' : '已关闭';
    await this.page.getByRole('option', { name: optionText }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by conversation type
   */
  async filterByType(type: 'user_order' | 'player_service') {
    await this.typeSelect.click();
    const optionText = type === 'user_order' ? '订单陪玩' : '陪练服务';
    await this.page.getByRole('option', { name: optionText }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click quick action button
   */
  async clickQuickAction(action: 'active' | 'closed' | 'today' | 'todayActive') {
    switch (action) {
      case 'active':
        await this.activeOnlyButton.click();
        break;
      case 'closed':
        await this.closedOnlyButton.click();
        break;
      case 'today':
        await this.todayCreatedButton.click();
        break;
      case 'todayActive':
        await this.todayActiveButton.click();
        break;
    }
    await this.page.waitForTimeout(1000);
  }

  /**
   * Refresh the list
   */
  async refresh() {
    await this.refreshButton.click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * View conversation details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
    await this.modal.waitFor({ state: 'visible', timeout: 5000 });
  }

  /**
   * Close conversation from table
   */
  async closeConversation(rowIndex: number) {
    await this.tableRows.nth(rowIndex).getByRole('button', { name: /关闭/i }).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible()) {
      await confirmButton.click();
    }
  }

  /**
   * Reopen conversation from table
   */
  async reopenConversation(rowIndex: number) {
    await this.tableRows.nth(rowIndex).getByRole('button', { name: /打开/i }).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible()) {
      await confirmButton.click();
    }
  }

  /**
   * Close conversation from modal
   */
  async closeConversationFromModal() {
    await this.modalCloseConversationButton.click();
  }

  /**
   * Reopen conversation from modal
   */
  async reopenConversationFromModal() {
    await this.modalReopenConversationButton.click();
    // Handle Popconfirm
    const confirmButton = this.page.locator('.ant-modal').getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible()) {
      await confirmButton.click();
    }
  }

  /**
   * Close modal
   */
  async closeModal() {
    await this.modalCloseButton.click();
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
   * Verify conversation exists in table
   */
  async verifyConversationExists(orderNo: string): Promise<boolean> {
    const count = await this.getConversationCount();
    for (let i = 0; i < count; i++) {
      const text = await this.getCellText(i, 2); // Order number is in column 2
      if (text.includes(orderNo)) {
        return true;
      }
    }
    return false;
  }
}
