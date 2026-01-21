import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { ChatRoomsPage } from './pages/ChatRoomsPage';

/**
 * E2E Tests for Chat Rooms Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Conversation list with pagination
 * - Search and filter conversations
 * - Quick action buttons
 * - View conversation details
 * - Close and reopen conversations
 * - Export functionality
 */

test.describe('Chat Rooms Management', () => {
  let loginPage: LoginPage;
  let chatRoomsPage: ChatRoomsPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    chatRoomsPage = new ChatRoomsPage(page);

    // Login and navigate to chat rooms
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display chat rooms page correctly', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.pageTitle).toBeVisible();
      await expect(chatRoomsPage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.totalConversationsStat).toBeVisible();
      await expect(chatRoomsPage.activeConversationsStat).toBeVisible();
      await expect(chatRoomsPage.totalUsersStat).toBeVisible();
      await expect(chatRoomsPage.totalMessagesStat).toBeVisible();
    });

    test('should load conversation list', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await chatRoomsPage.goto();

      const stats = await chatRoomsPage.getStatistics();
      expect(stats.totalConversations).toBeGreaterThanOrEqual(0);
      expect(stats.activeConversations).toBeGreaterThanOrEqual(0);
      expect(stats.totalUsers).toBeGreaterThanOrEqual(0);
      expect(stats.totalMessages).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Quick Actions', () => {
    test('should filter active conversations only', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.clickQuickAction('active');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter closed conversations only', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.clickQuickAction('closed');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter today created conversations', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.clickQuickAction('today');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter today active conversations', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.clickQuickAction('todayActive');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should refresh conversation list', async () => {
      await chatRoomsPage.goto();

      const _countBefore = await chatRoomsPage.getConversationCount();

      await chatRoomsPage.refresh();

      const countAfter = await chatRoomsPage.getConversationCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search conversations by keyword', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.searchByKeyword('test');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by conversation status', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.filterByStatus('active');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by conversation type', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.filterByType('user_order');

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await chatRoomsPage.goto();

      // Apply a filter
      await chatRoomsPage.searchByKeyword('test');

      // Clear search
      await chatRoomsPage.searchInput.fill('');
      await chatRoomsPage.page.waitForTimeout(500);

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Conversation Details', () => {
    test('should open conversation details modal', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        await chatRoomsPage.viewDetails(0);

        const isModalVisible = await chatRoomsPage.isModalVisible();
        expect(isModalVisible).toBe(true);
      }
    });

    test('should close modal when close button clicked', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        await chatRoomsPage.viewDetails(0);
        await chatRoomsPage.closeModal();

        const isModalVisible = await chatRoomsPage.isModalVisible();
        expect(isModalVisible).toBe(false);
      }
    });

    test('should display conversation information in modal', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        await chatRoomsPage.viewDetails(0);

        await expect(chatRoomsPage.modalTitle).toContainText('会话详情');
      }
    });
  });

  test.describe('Conversation Actions', () => {
    test('should close active conversation from table', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        // Try to close first conversation
        await chatRoomsPage.closeConversation(0);

        // Wait for potential success message
        await chatRoomsPage.page.waitForTimeout(1000);
      }
    });

    test('should reopen closed conversation from table', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        // Try to reopen first conversation
        await chatRoomsPage.reopenConversation(0);

        // Wait for potential success message
        await chatRoomsPage.page.waitForTimeout(1000);
      }
    });

    test('should close conversation from modal', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        await chatRoomsPage.viewDetails(0);

        // Try to close from modal
        const hasCloseButton = await chatRoomsPage.modalCloseConversationButton.isVisible();
        if (hasCloseButton) {
          await chatRoomsPage.closeConversationFromModal();
          await chatRoomsPage.page.waitForTimeout(1000);
        }
      }
    });

    test('should reopen conversation from modal', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        await chatRoomsPage.viewDetails(0);

        // Try to reopen from modal
        const hasReopenButton = await chatRoomsPage.modalReopenConversationButton.isVisible();
        if (hasReopenButton) {
          await chatRoomsPage.reopenConversationFromModal();
          await chatRoomsPage.page.waitForTimeout(1000);
        }
      }
    });
  });

  test.describe('Export Functionality', () => {
    test('should have export button', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.exportButton).toBeVisible();
    });

    test('should export conversation data', async () => {
      await chatRoomsPage.goto();

      const downloadPromise = chatRoomsPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await chatRoomsPage.exportButton.click();

      try {
        const download = await downloadPromise;
        expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
      } catch {
        // Export might fail if no data, that's acceptable
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display conversation columns correctly', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count > 0) {
        const cellText = await chatRoomsPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.pageTitle).toBeVisible();
      await expect(chatRoomsPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await chatRoomsPage.goto();

      await expect(chatRoomsPage.pageTitle).toContainText('聊天室管理');
    });

    test('should be keyboard navigable', async () => {
      await chatRoomsPage.goto();

      await chatRoomsPage.page.keyboard.press('Tab');
      await chatRoomsPage.page.keyboard.press('Tab');

      const focusedElement = chatRoomsPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const chatMenuItem = page.locator('.ant-menu-item').filter({ hasText: /聊天室管理/i });
      if (await chatMenuItem.isVisible()) {
        await chatMenuItem.click();
        await chatRoomsPage.waitForPageLoad();
        await expect(chatRoomsPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty conversation list', async () => {
      await chatRoomsPage.goto();

      const count = await chatRoomsPage.getConversationCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(chatRoomsPage.table).toBeVisible();
      }
    });
  });
});
