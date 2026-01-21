import { Page } from '@playwright/test';

/**
 * Page Object Model for Referral Management Page
 * Encapsulates all referral management interactions and assertions
 */
export class ReferralPage {
  constructor(private page: Page) {}

  // Element locators
  readonly pageTitle = this.page.locator('h1').filter({ hasText: /推荐管理/i });

  // Statistics
  readonly totalReferralsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /总推荐数/i });
  readonly completedReferralsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /已完成推荐/i });
  readonly issuedRewardsStat = this.page.locator('.ant-statistic-title').filter({ hasText: /已发放奖励/i });
  readonly activeCodesStat = this.page.locator('.ant-statistic-title').filter({ hasText: /活跃邀请码/i });

  // Tabs
  readonly tabs = this.page.locator('.ant-tabs');
  readonly referralListTab = this.page.getByRole('tab').filter({ hasText: /推荐关系/i });
  readonly codesManagementTab = this.page.getByRole('tab').filter({ hasText: /邀请码管理/i });
  readonly rewardsManagementTab = this.page.getByRole('tab').filter({ hasText: /奖励管理/i });

  /**
   * Navigate to referral page
   */
  async goto() {
    await this.page.goto('/admin/referral');
    await this.waitForPageLoad();
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await Promise.race([
      this.pageTitle.waitFor({ state: 'visible', timeout: 15000 }),
      this.tabs.waitFor({ state: 'visible', timeout: 15000 }),
    ]);
  }

  /**
   * Get statistics
   */
  async getStatistics(): Promise<{
    totalReferrals: number;
    completedReferrals: number;
    issuedRewards: number;
    activeCodes: number;
  }> {
    const getValue = async (statLocator: ReturnType<typeof this.page.locator>) => {
      const content = statLocator.locator('..').locator('.ant-statistic-content');
      const text = await content.textContent() || '0';
      const match = text.match(/([\d.]+)/);
      return match ? parseFloat(match[1]) : 0;
    };

    return {
      totalReferrals: await getValue(this.totalReferralsStat),
      completedReferrals: await getValue(this.completedReferralsStat),
      issuedRewards: await getValue(this.issuedRewardsStat),
      activeCodes: await getValue(this.activeCodesStat),
    };
  }

  /**
   * Switch to specific tab
   */
  async switchTab(tab: 'referrals' | 'codes' | 'rewards') {
    if (tab === 'referrals') {
      await this.referralListTab.click();
    } else if (tab === 'codes') {
      await this.codesManagementTab.click();
    } else if (tab === 'rewards') {
      await this.rewardsManagementTab.click();
    }
    await this.page.waitForTimeout(500);
  }

  /**
   * Get current active tab
   */
  async getActiveTab(): Promise<string> {
    const activeTab = this.tabs.locator('.ant-tabs-tab-active');
    return await activeTab.textContent() || '';
  }

  /**
   * Verify statistics are displayed
   */
  async areStatisticsVisible(): Promise<boolean> {
    return await this.totalReferralsStat.isVisible() &&
           await this.completedReferralsStat.isVisible() &&
           await this.issuedRewardsStat.isVisible() &&
           await this.activeCodesStat.isVisible();
  }
}
