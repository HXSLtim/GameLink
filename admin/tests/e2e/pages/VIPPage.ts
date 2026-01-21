import { Page } from '@playwright/test';

/**
 * Page Object Model for VIP Management Page
 * Encapsulates all VIP level management interactions and assertions
 */
export class VIPPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /VIP管理/i });

  // Statistics
  readonly totalLevelsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总等级数/i });
  readonly activeLevelsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /启用等级/i });

  // View mode segmented control
  readonly viewModeAll = this.page.getByText(/全部/i);
  readonly viewModeActive = this.page.getByText(/启用/i);
  readonly viewModeInactive = this.page.getByText(/禁用/i);

  // Search
  readonly searchInput = this.page.getByPlaceholder(/搜索等级名称/i);

  // Action buttons
  readonly addButton = this.page.getByRole('button', { name: /新增等级/i });
  readonly refreshButton = this.page.getByRole('button', { name: /刷新/i }).filter({ hasText: /^刷新$/ });

  // VIP cards
  readonly vipCards = this.page.locator('.ant-card').filter({ has: '.ant-card-head-title' });

  // Card actions
  readonly editButton = (cardIndex: number) => this.vipCards.nth(cardIndex).getByRole('button', { name: /编辑/i });
  readonly deleteButton = (cardIndex: number) => this.vipCards.nth(cardIndex).getByRole('button', { name: /删除/i });
  readonly setDefaultButton = (cardIndex: number) => this.vipCards.nth(cardIndex).getByRole('button', { name: /设为默认/i });
  readonly editBenefitsButton = (cardIndex: number) => this.vipCards.nth(cardIndex).getByRole('button', { name: /编辑权益/i });
  readonly activeSwitch = (cardIndex: number) => this.vipCards.nth(cardIndex).locator('.ant-switch');

  // Modal
  private readonly modal = this.page.locator('.ant-modal');
  readonly modalTitle = this.modal.locator('.ant-modal-title');
  private readonly modalOkButton = this.modal.getByRole('button', { name: /确定/i });
  private readonly modalCancelButton = this.modal.getByRole('button', { name: /取消/i });

  /**
   * Navigate to VIP page
   */
  async goto() {
    await this.page.goto('/admin/vip');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.vipCards.first().waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Get VIP card count
   */
  async getVIPCardCount(): Promise<number> {
    await this.page.waitForTimeout(500);
    return await this.vipCards.count();
  }

  /**
   * Get statistics
   */
  async getStatistics(): Promise<{ total: number; active: number }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/(\d+)/);
      return match ? parseInt(match[1], 10) : 0;
    };

    return {
      total: await getValue(this.totalLevelsStat),
      active: await getValue(this.activeLevelsStat),
    };
  }

  /**
   * Switch view mode
   */
  async switchViewMode(mode: 'all' | 'active' | 'inactive') {
    if (mode === 'all') {
      await this.viewModeAll.click();
    } else if (mode === 'active') {
      await this.viewModeActive.click();
    } else {
      await this.viewModeInactive.click();
    }
    await this.page.waitForTimeout(500);
  }

  /**
   * Search VIP levels by keyword
   */
  async searchByKeyword(keyword: string) {
    await this.searchInput.fill(keyword);
    await this.page.waitForTimeout(500);
  }

  /**
   * Click add button
   */
  async clickAdd() {
    await this.addButton.click();
  }

  /**
   * Edit VIP level
   */
  async editLevel(cardIndex: number) {
    await this.editButton(cardIndex).click();
  }

  /**
   * Delete VIP level
   */
  async deleteLevel(cardIndex: number) {
    await this.deleteButton(cardIndex).click();
    // Handle Popconfirm
    const confirmButton = this.page.getByRole('button', { name: /确定/i });
    if (await confirmButton.isVisible({ timeout: 5000 })) {
      await confirmButton.click();
    }
  }

  /**
   * Set VIP level as default
   */
  async setAsDefault(cardIndex: number) {
    await this.setDefaultButton(cardIndex).click();
  }

  /**
   * Edit VIP benefits
   */
  async editBenefits(cardIndex: number) {
    await this.editBenefitsButton(cardIndex).click();
  }

  /**
   * Toggle active status
   */
  async toggleActive(cardIndex: number) {
    await this.activeSwitch(cardIndex).click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Refresh the list
   */
  async refresh() {
    await this.refreshButton.click();
    await this.page.waitForTimeout(1000);
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
   * Get VIP card title
   */
  async getCardTitle(cardIndex: number): Promise<string> {
    const title = this.vipCards.nth(cardIndex).locator('.ant-card-head-title');
    return await title.textContent() || '';
  }

  /**
   * Check if VIP level is active
   */
  async isLevelActive(cardIndex: number): Promise<boolean> {
    const switchElement = this.activeSwitch(cardIndex);
    const className = await switchElement.getAttribute('class') || '';
    return className.includes('ant-switch-checked');
  }
}
