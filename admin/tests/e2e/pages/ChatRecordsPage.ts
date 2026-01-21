import { Page } from '@playwright/test';

/**
 * Page Object Model for Chat Records Management Page
 * Encapsulates all chat records management interactions and assertions
 */
export class ChatRecordsPage {
  constructor(private page: Page) {}

  // Element locators
  private readonly pageTitle = this.page.locator('h1').filter({ hasText: /聊天记录管理/i });
  private readonly table = this.page.locator('.ant-table');
  private readonly tableRows = this.page.locator('.ant-table-tbody tr');

  // Statistics cards
  readonly totalMessagesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总消息数/i });
  readonly todayMessagesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /今日消息/i });
  readonly activeConversationsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /活跃会话数/i });
  readonly totalUsersStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总用户数/i });

  // Quick filter buttons
  readonly allButton = this.page.getByRole('button', { name: /全部/i });
  readonly todayButton = this.page.getByRole('button', { name: /今天/i });
  readonly weekButton = this.page.getByRole('button', { name: /最近7天/i });
  readonly textMessagesButton = this.page.getByRole('button', { name: /文本消息/i });
  readonly imageMessagesButton = this.page.getByRole('button', { name: /图片消息/i });
  readonly systemMessagesButton = this.page.getByRole('button', { name: /系统消息/i });
  readonly refreshButton = this.page.getByRole('button', { name: /^刷新$/i });

  // Search fields
  private readonly searchInput = this.page.getByPlaceholder(/搜索消息内容、发送者/i);
  private readonly messageTypeSelect = this.page.getByText(/消息类型/i).first();
  private readonly senderTypeSelect = this.page.getByText(/发送者类型/i).first();
  private readonly conversationIdInput = this.page.getByPlaceholder(/输入会话ID/i);

  // Toolbar buttons
  readonly exportButton = this.page.getByRole('button', { name: /导出数据/i });

  // Table action buttons
  private readonly detailButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /详情/i });
  private readonly deleteButton = (rowIndex: number) => this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i });

  // Drawer
  private readonly drawer = this.page.locator('.ant-drawer');
  private readonly drawerTitle = this.drawer.locator('.ant-drawer-title');
  private readonly drawerCloseButton = this.drawer.locator('.ant-drawer-close');

  /**
   * Navigate to chat records page
   */
  async goto() {
    await this.page.goto('/admin/chat/records');
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
   * Get message count from table
   */
  async getMessageCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.tableRows.count();
  }

  /**
   * Get message statistics
   */
  async getStatistics(): Promise<{
    totalMessages: number;
    todayMessages: number;
    activeConversations: number;
    totalUsers: number;
  }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      totalMessages: await getValue(this.totalMessagesStat),
      todayMessages: await getValue(this.todayMessagesStat),
      activeConversations: await getValue(this.activeConversationsStat),
      totalUsers: await getValue(this.totalUsersStat),
    };
  }

  /**
   * Search messages by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Filter by message type
   */
  async filterByMessageType(type: 'text' | 'image' | 'system') {
    await this.messageTypeSelect.click();
    const typeMap = { text: '文本', image: '图片', system: '系统' };
    await this.page.getByRole('option', { name: typeMap[type] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Filter by sender type
   */
  async filterBySenderType(senderType: 'user' | 'player' | 'admin' | 'system') {
    await this.senderTypeSelect.click();
    const typeMap = { user: '用户', player: '陪玩师', admin: '管理员', system: '系统' };
    await this.page.getByRole('option', { name: typeMap[senderType] }).click();
    await this.page.waitForTimeout(1000);
  }

  /**
   * Click quick filter button
   */
  async clickQuickFilter(filter: 'all' | 'today' | 'week' | 'text' | 'image' | 'system') {
    switch (filter) {
      case 'all':
        await this.allButton.click();
        break;
      case 'today':
        await this.todayButton.click();
        break;
      case 'week':
        await this.weekButton.click();
        break;
      case 'text':
        await this.textMessagesButton.click();
        break;
      case 'image':
        await this.imageMessagesButton.click();
        break;
      case 'system':
        await this.systemMessagesButton.click();
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
   * View message details
   */
  async viewDetails(rowIndex: number) {
    await this.detailButton(rowIndex).click();
    await this.drawer.waitFor({ state: 'visible', timeout: 5000 });
  }

  /**
   * Delete message from table
   */
  async deleteMessage(rowIndex: number) {
    await this.tableRows.nth(rowIndex).getByRole('button', { name: /删除/i }).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Close drawer
   */
  async closeDrawer() {
    await this.drawerCloseButton.click();
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
   * Verify message exists in table
   */
  async verifyMessageExists(content: string): Promise<boolean> {
    const count = await this.getMessageCount();
    for (let i = 0; i < count; i++) {
      const text = await this.getCellText(i, 4); // Message content is in column 4
      if (text.includes(content)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Search by conversation ID
   */
  async searchByConversationId(conversationId: string) {
    await this.conversationIdInput.fill(conversationId);
    await this.page.waitForTimeout(500);
  }

  /**
   * Clear search filters
   */
  async clearFilters() {
    await this.searchInput.fill('');
    await this.page.waitForTimeout(500);
  }
}
