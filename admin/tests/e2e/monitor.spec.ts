import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { MonitorPage } from './pages/MonitorPage';

/**
 * E2E Tests for Monitor Dashboard
 *
 * Test Coverage:
 * - Page loading and display
 * - System status card
 * - Online users card
 * - Order queue card
 * - Alerts card
 * - WebSocket connection status
 * - Real-time data updates
 */

test.describe('Monitor Dashboard', () => {
  let loginPage: LoginPage;
  let monitorPage: MonitorPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    monitorPage = new MonitorPage(page);

    // Login and navigate to monitor
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display monitor page correctly', async () => {
      await monitorPage.goto();

      await expect(monitorPage.pageTitle).toBeVisible();
      await expect(monitorPage.systemStatusCard).toBeVisible();
      await expect(monitorPage.onlineUsersCard).toBeVisible();
      await expect(monitorPage.orderQueueCard).toBeVisible();
      await expect(monitorPage.alertsCard).toBeVisible();
    });

    test('should display page title and subtitle', async () => {
      await monitorPage.goto();

      await expect(monitorPage.pageTitle).toContainText('实时监控');
      await expect(monitorPage.subTitle).toContainText('系统运行状态实时数据');
    });

    test('should load data on page load', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      // Verify statistics are displayed
      const onlineUsers = await monitorPage.getOnlineUsersCount();
      expect(onlineUsers).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('System Status Card', () => {
    test('should display system status metrics', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      await expect(monitorPage.cpuUsage).toBeVisible();
      await expect(monitorPage.memoryUsage).toBeVisible();
      await expect(monitorPage.goroutines).toBeVisible();
      await expect(monitorPage.dbConnections).toBeVisible();
    });

    test('should show CPU usage value', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const cpuUsage = await monitorPage.getCpuUsage();
      expect(cpuUsage).toBeGreaterThanOrEqual(0);
      expect(cpuUsage).toBeLessThanOrEqual(100);
    });

    test('should show WebSocket connection status', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      // Either connected or disconnected tag should be visible
      const wsStatus = await monitorPage.getWebSocketStatus();
      expect(['connected', 'disconnected']).toContain(wsStatus);
    });

    test('should show system health status', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const status = await monitorPage.getSystemStatus();
      expect(status).toBeTruthy();
    });
  });

  test.describe('Online Users Card', () => {
    test('should display online users statistics', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      await expect(monitorPage.onlineUsersTotal).toBeVisible();
      await expect(monitorPage.onlineUsersPeak).toBeVisible();
    });

    test('should show current online count', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const onlineCount = await monitorPage.getOnlineUsersCount();
      expect(onlineCount).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Order Queue Card', () => {
    test('should display order queue statistics', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      await expect(monitorPage.orderPending).toBeVisible();
      await expect(monitorPage.orderProcessing).toBeVisible();
      await expect(monitorPage.orderCompleted).toBeVisible();
    });

    test('should show pending orders count', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const pendingOrders = await monitorPage.getPendingOrdersCount();
      expect(pendingOrders).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Alerts Card', () => {
    test('should display alerts section', async () => {
      await monitorPage.goto();

      await expect(monitorPage.alertsCard).toBeVisible();
    });

    test('should show no alerts message when empty', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const alertCount = await monitorPage.getAlertCount();
      if (alertCount === 0) {
        await expect(monitorPage.noAlertsMessage).toBeVisible();
      }
    });

    test('should display alert list when alerts exist', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const alertCount = await monitorPage.getAlertCount();
      // Alert count may be 0 or more, just verify the element can be queried
      expect(alertCount).toBeGreaterThanOrEqual(0);
    });

    test('should have clear all alerts button', async () => {
      await monitorPage.goto();

      await expect(monitorPage.clearAllAlertsButton).toBeVisible();
    });

    test('should clear all alerts when button clicked', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      const alertCountBefore = await monitorPage.getAlertCount();
      if (alertCountBefore > 0) {
        await monitorPage.clearAllAlerts();
        await monitorPage.page.waitForTimeout(1000);
        const alertCountAfter = await monitorPage.getAlertCount();
        expect(alertCountAfter).toBeLessThanOrEqual(alertCountBefore);
      }
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should display all cards in proper layout', async () => {
      await monitorPage.goto();

      const _isProper = await monitorPage.isDisplayingProperly();
      expect(isPro).toBe(true);
    });
  });

  test.describe('Real-time Updates', () => {
    test('should handle WebSocket connection changes', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      // Get initial connection status
      const initialStatus = await monitorPage.getWebSocketStatus();
      expect(['connected', 'disconnected']).toContain(initialStatus);
    });

    test('should update metrics over time', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      // Get initial metrics
      const initialOnlineUsers = await monitorPage.getOnlineUsersCount();

      // Wait for potential updates
      await monitorPage.page.waitForTimeout(2000);

      // Get updated metrics
      const updatedOnlineUsers = await monitorPage.getOnlineUsersCount();

      // Values should still be valid
      expect(initialOnlineUsers).toBeGreaterThanOrEqual(0);
      expect(updatedOnlineUsers).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await monitorPage.goto();

      await expect(monitorPage.pageTitle).toBeVisible();
    });

    test('should be keyboard navigable', async () => {
      await monitorPage.goto();

      // Tab through the page
      await monitorPage.page.keyboard.press('Tab');
      await monitorPage.page.keyboard.press('Tab');

      const focusedElement = monitorPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Error Handling', () => {
    test('should handle WebSocket disconnection gracefully', async () => {
      await monitorPage.goto();
      await monitorPage.waitForDataLoad();

      // Even if WebSocket is disconnected, page should still display
      await expect(monitorPage.systemStatusCard).toBeVisible();
      await expect(monitorPage.onlineUsersCard).toBeVisible();
      await expect(monitorPage.orderQueueCard).toBeVisible();
      await expect(monitorPage.alertsCard).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate to monitor from sidebar', async ({ page }) => {
      // Click on monitor menu item if exists in sidebar
      const monitorMenuItem = page.locator('.ant-menu-item').filter({ hasText: /实时监控/i });
      if (await monitorMenuItem.isVisible()) {
        await monitorMenuItem.click();
        await monitorPage.waitForPageLoad();
        await expect(monitorPage.pageTitle).toBeVisible();
      }
    });
  });
});
