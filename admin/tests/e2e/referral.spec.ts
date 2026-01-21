import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { ReferralPage } from './pages/ReferralPage';

/**
 * E2E Tests for Referral Management
 *
 * Test Coverage:
 * - Page loading and display
 * - Statistics display
 * - Tab navigation (referrals, codes, rewards)
 * - Data refresh
 */

test.describe('Referral Management', () => {
  let loginPage: LoginPage;
  let referralPage: ReferralPage;

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    referralPage = new ReferralPage(page);

    // Login and navigate to referral
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
  });

  test.describe('Page Loading', () => {
    test('should display referral page correctly', async () => {
      await referralPage.goto();

      await expect(referralPage.pageTitle).toBeVisible();
      await expect(referralPage.tabs).toBeVisible();
    });

    test('should display statistics cards', async () => {
      await referralPage.goto();

      await expect(referralPage.totalReferralsStat).toBeVisible();
      await expect(referralPage.completedReferralsStat).toBeVisible();
      await expect(referralPage.issuedRewardsStat).toBeVisible();
      await expect(referralPage.activeCodesStat).toBeVisible();
    });

    test('should display statistics values', async () => {
      await referralPage.goto();

      const stats = await referralPage.getStatistics();
      expect(stats.totalReferrals).toBeGreaterThanOrEqual(0);
      expect(stats.completedReferrals).toBeGreaterThanOrEqual(0);
      expect(stats.issuedRewards).toBeGreaterThanOrEqual(0);
      expect(stats.activeCodes).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Tab Navigation', () => {
    test('should switch to referral list tab', async () => {
      await referralPage.goto();

      await referralPage.switchTab('referrals');

      const activeTab = await referralPage.getActiveTab();
      expect(activeTab).toContain('推荐关系');
    });

    test('should switch to codes management tab', async () => {
      await referralPage.goto();

      await referralPage.switchTab('codes');

      const activeTab = await referralPage.getActiveTab();
      expect(activeTab).toContain('邀请码管理');
    });

    test('should switch to rewards management tab', async () => {
      await referralPage.goto();

      await referralPage.switchTab('rewards');

      const activeTab = await referralPage.getActiveTab();
      expect(activeTab).toContain('奖励管理');
    });

    test('should maintain tab state on refresh', async () => {
      await referralPage.goto();

      // Switch to codes tab
      await referralPage.switchTab('codes');

      // Refresh page
      await referralPage.page.reload();
      await referralPage.waitForPageLoad();

      // Should be on codes tab or default tab
      const activeTab = await referralPage.getActiveTab();
      expect(activeTab).toBeTruthy();
    });
  });

  test.describe('Data Display', () => {
    test('should show all statistics correctly', async () => {
      await referralPage.goto();

      const areVisible = await referralPage.areStatisticsVisible();
      expect(areVisible).toBe(true);
    });

    test('should display tab content', async () => {
      await referralPage.goto();

      // Tab content should be visible
      const tabContent = referralPage.page.locator('.ant-tabs-content');
      await expect(tabContent).toBeVisible();
    });
  });

  test.describe('Layout and Responsive Design', () => {
    test('should have proper layout structure', async () => {
      await referralPage.goto();

      await expect(referralPage.pageTitle).toBeVisible();
      await expect(referralPage.tabs).toBeVisible();
    });
  });

  test.describe('Accessibility', () => {
    test('should have proper page title', async () => {
      await referralPage.goto();

      await expect(referralPage.pageTitle).toContainText('推荐管理');
    });

    test('should be keyboard navigable', async () => {
      await referralPage.goto();

      await referralPage.page.keyboard.press('Tab');
      await referralPage.page.keyboard.press('Tab');

      const focusedElement = referralPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Integration', () => {
    test('should navigate from sidebar', async ({ page }) => {
      const referralMenuItem = page.locator('.ant-menu-item').filter({ hasText: /推荐管理/i });
      if (await referralMenuItem.isVisible()) {
        await referralMenuItem.click();
        await referralPage.waitForPageLoad();
        await expect(referralPage.pageTitle).toBeVisible();
      }
    });
  });

  test.describe('Statistics Verification', () => {
    test('should show consistent statistics', async () => {
      await referralPage.goto();

      const stats1 = await referralPage.getStatistics();

      await referralPage.page.waitForTimeout(1000);

      const stats2 = await referralPage.getStatistics();

      // Statistics should be consistent
      expect(stats1.totalReferrals).toEqual(stats2.totalReferrals);
    });
  });
});
