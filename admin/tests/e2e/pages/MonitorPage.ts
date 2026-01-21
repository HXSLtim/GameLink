import { Page } from '@playwright/test';

/**
 * Page Object Model for Monitor Page
 * Encapsulates all monitoring dashboard interactions and assertions
 */
export class MonitorPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /实时监控/i });
  readonly subTitle = this.page.locator('.ant-typography').filter({ hasText: /系统运行状态实时数据/i });

  // System status card
  readonly systemStatusCard = this.page.locator('.ant-card').filter({ hasText: /系统状态/i });
  readonly cpuUsage = this.page.locator('.ant-statistic-title').filter({ hasText: /CPU 使用率/i });
  readonly memoryUsage = this.page.locator('.ant-statistic-title').filter({ hasText: /内存使用/i });
  readonly goroutines = this.page.locator('.ant-statistic-title').filter({ hasText: /Goroutines/i });
  readonly dbConnections = this.page.locator('.ant-statistic-title').filter({ hasText: /数据库连接/i });

  // Connection status tag
  readonly wsConnectedTag = this.page.locator('.ant-tag').filter({ hasText: /已连接/i });
  readonly wsDisconnectedTag = this.page.locator('.ant-tag').filter({ hasText: /未连接/i });

  // Online users card
  readonly onlineUsersCard = this.page.locator('.ant-card').filter({ hasText: /在线用户/i });
  readonly onlineUsersTotal = this.page.locator('.ant-statistic-title').filter({ hasText: /当前在线/i });
  readonly onlineUsersPeak = this.page.locator('.ant-statistic-title').filter({ hasText: /峰值/i });

  // Order queue card
  readonly orderQueueCard = this.page.locator('.ant-card').filter({ hasText: /订单队列/i });
  readonly orderPending = this.page.locator('.ant-statistic-title').filter({ hasText: /待处理/i });
  readonly orderProcessing = this.page.locator('.ant-statistic-title').filter({ hasText: /处理中/i });
  readonly orderCompleted = this.page.locator('.ant-statistic-title').filter({ hasText: /已完成/i });

  // Alerts card
  readonly alertsCard = this.page.locator('.ant-card').filter({ hasText: /警告通知/i });
  readonly alertList = this.page.locator('.ant-list-item');
  readonly noAlertsMessage = this.page.locator('.ant-alert').filter({ hasText: /暂无警告/i });
  readonly clearAllAlertsButton = this.page.getByRole('button', { name: /清空全部/i });

  /**
   * Navigate to monitor page
   */
  async goto() {
    await this.page.goto('/admin/monitor');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.systemStatusCard.waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Get WebSocket connection status
   */
  async getWebSocketStatus(): Promise<'connected' | 'disconnected'> {
    const isConnected = await this.wsConnectedTag.isVisible();
    return isConnected ? 'connected' : 'disconnected';
  }

  /**
   * Get system status value
   */
  async getSystemStatus(): Promise<string | null> {
    const statusTag = this.systemStatusCard.locator('.ant-tag').nth(1);
    if (await statusTag.isVisible()) {
      return await statusTag.textContent();
    }
    return null;
  }

  /**
   * Get CPU usage value
   */
  async getCpuUsage(): Promise<number> {
    const content = this.page.locator('.ant-statistic-title').filter({ hasText: /CPU 使用率/i })
      .locator('..')
      .locator('.ant-statistic-content');
    const text = await content.textContent() || '0';
    const match = text.match(/(\d+)/);
    return match ? parseInt(match[1], 10) : 0;
  }

  /**
   * Get online users count
   */
  async getOnlineUsersCount(): Promise<number> {
    const content = this.page.locator('.ant-statistic-title').filter({ hasText: /当前在线/i })
      .locator('..')
      .locator('.ant-statistic-content');
    const text = await content.textContent() || '0';
    const match = text.match(/(\d+)/);
    return match ? parseInt(match[1], 10) : 0;
  }

  /**
   * Get pending orders count
   */
  async getPendingOrdersCount(): Promise<number> {
    const content = this.page.locator('.ant-statistic-title').filter({ hasText: /待处理/i })
      .locator('..')
      .locator('.ant-statistic-content');
    const text = await content.textContent() || '0';
    const match = text.match(/(\d+)/);
    return match ? parseInt(match[1], 10) : 0;
  }

  /**
   * Get alert count
   */
  async getAlertCount(): Promise<number> {
    await this.page.waitForTimeout(1000);
    return await this.alertList.count();
  }

  /**
   * Clear all alerts
   */
  async clearAllAlerts() {
    await this.clearAllAlertsButton.click();
  }

  /**
   * Wait for data to load
   */
  async waitForDataLoad() {
    // Wait for at least one card to show actual data (not loading)
    await this.page.waitForFunction(() => {
      const stats = document.querySelectorAll('.ant-statistic-content-value');
      return stats.length > 0;
    }, { timeout: 10000 });
  }

  /**
   * Check if page displays properly
   */
  async isDisplayingProperly(): Promise<boolean> {
    const titleVisible = await this.pageTitle.isVisible();
    const systemCardVisible = await this.systemStatusCard.isVisible();
    const onlineUsersCardVisible = await this.onlineUsersCard.isVisible();
    const orderQueueCardVisible = await this.orderQueueCard.isVisible();
    const alertsCardVisible = await this.alertsCard.isVisible();

    return titleVisible && systemCardVisible && onlineUsersCardVisible &&
           orderQueueCardVisible && alertsCardVisible;
  }
}
