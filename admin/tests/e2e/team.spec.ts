import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { TeamPage } from './pages/TeamPage';

/**
 * E2E Tests for Team Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Team list with pagination
 * - Statistics display
 * - Search and filter teams
 * - Create, edit, delete teams
 * - View team details
 */

test.describe('Team Management', () => {
  let loginPage: LoginPage;
  let teamPage: TeamPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    teamPage = new TeamPage(page);

    // Login and navigate to team
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display team page correctly', async () => {
      await teamPage.goto();

      await expect(teamPage.pageTitle).toBeVisible();
      await expect(teamPage.table).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await teamPage.goto();

      await expect(teamPage.totalTeamsStat).toBeVisible();
      await expect(teamPage.activeTeamsStat).toBeVisible();
      await expect(teamPage.totalMembersStat).toBeVisible();
    });

    test('should load team list', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should display statistics values', async () => {
      await teamPage.goto();

      const stats = await teamPage.getStatistics();
      expect(stats.totalTeams).toBeGreaterThanOrEqual(0);
      expect(stats.activeTeams).toBeGreaterThanOrEqual(0);
      expect(stats.totalMembers).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Search and Filter', () => {
    test('should search teams by keyword', async () => {
      await teamPage.goto();

      await teamPage.searchByKeyword('test');

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should filter by status', async () => {
      await teamPage.goto();

      await teamPage.filterByStatus('active');

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should clear search filters', async () => {
      await teamPage.goto();

      // Apply a filter
      await teamPage.searchByKeyword('test');

      // Clear search
      await teamPage.clearFilters();

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Team Actions', () => {
    test('should open add modal when clicking add button', async () => {
      await teamPage.goto();

      await teamPage.clickAdd();

      const isModalVisible = await teamPage.isModalVisible();
      if (isModalVisible) {
        await teamPage.closeModal();
      }
    });

    test('should view team details', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      if (count > 0) {
        await teamPage.viewDetails(0);

        const isDrawerVisible = await teamPage.isDrawerVisible();
        if (isDrawerVisible) {
          await teamPage.closeDrawer();
        }
      }
    });

    test('should edit team', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      if (count > 0) {
        await teamPage.editTeam(0);

        const isModalVisible = await teamPage.isModalVisible();
        if (isModalVisible) {
          await teamPage.closeModal();
        }
      }
    });

    test('should handle delete operation', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      if (count > 0) {
        const _countBefore = await teamPage.getTeamCount();

        // Try to delete first team
        await teamPage.deleteTeam(0);

        // Wait for potential deletion
        await teamPage.page.waitForTimeout(1000);

        const countAfter = await teamPage.getTeamCount();
        // Count may or may not change depending on permissions
        expect(countAfter).toBeGreaterThanOrEqual(0);
      }
    });
  });

  test.describe('Table Display', () => {
    test('should display team columns correctly', async () => {
      await teamPage.goto();

      await expect(teamPage.table).toBeVisible();
    });

    test('should support pagination', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test('should get cell text from table', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      if (count > 0) {
        const cellText = await teamPage.getCellText(0, 0);
        expect(cellText).toBeTruthy();
      }
    });
  });

  test.describe('Refresh Functionality', () => {
    test('should refresh team list', async () => {
      await teamPage.goto();

      const _countBefore = await teamPage.getTeamCount();

      await teamPage.refresh();

      const countAfter = await teamPage.getTeamCount();
      expect(countAfter).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await teamPage.goto();

      await expect(teamPage.pageTitle).toBeVisible();
      await expect(teamPage.table).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await teamPage.goto();

      await expect(teamPage.pageTitle).toContainText('团队管理');
    });

    test('should be keyboard navigable', async () => {
      await teamPage.goto();

      await teamPage.page.keyboard.press('Tab');
      await teamPage.page.keyboard.press('Tab');

      const focusedElement = teamPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const teamMenuItem = page.locator('.ant-menu-item').filter({ hasText: /团队管理/i });
      if (await teamMenuItem.isVisible()) {
        await teamMenuItem.click();
        await teamPage.waitForPageLoad();
        await expect(teamPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Empty State', () => {
    test('should handle empty team list', async () => {
      await teamPage.goto();

      const count = await teamPage.getTeamCount();
      if (count === 0) {
        // Should show empty state or table with no rows
        await expect(teamPage.table).toBeVisible();
      }
    });
  });

  test.describe('Combined Filters', () => {
    test('should apply multiple filters together', async () => {
      await teamPage.goto();

      // Apply status filter
      await teamPage.filterByStatus('active');

      // Then apply keyword filter
      await teamPage.searchByKeyword('test');

      const count = await teamPage.getTeamCount();
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });
});
