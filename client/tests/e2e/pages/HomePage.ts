/**
 * Home Page Object Model
 *
 * Based on actual HomePage.tsx implementation:
 * - Hero section with gradient background and CTA buttons
 * - Popular games grid (League of Legends, Valorant, Apex Legends, Overwatch 2)
 * - Featured players horizontal scroll
 * - Trust indicators section
 */
export class HomePage {
  constructor(private page: any) {}

  /** Hero section */
  get heroTitle() {
    // h1 contains hero_title and hero_highlight from i18n
    return this.page.getByRole('heading', { level: 1 }).or(
      this.page.locator('h1')
    );
  }
  get findPlayersButton() {
    // Primary CTA button - uses text content matching
    return this.page.getByRole('button').filter({ hasText: /arrow/i }).or(
      this.page.locator('button').filter({ hasText: /find|players|陪玩|找/i })
    ).first();
  }
  get becomePlayerButton() {
    // Secondary button - outline variant
    return this.page.getByRole('button').filter({ hasText: /become|player|陪玩/i }).or(
      this.page.locator('button[variant="outline"]')
    ).nth(1);
  }

  /** Popular games section */
  get popularGamesSection() {
    return this.page.getByText(/popular games|热门游戏/i).or(
      this.page.getByRole('heading', { name: /popular games|热门游戏/i })
    );
  }
  get gameCards() {
    // Cards in grid with game names
    return this.page.locator('.grid').locator('.card').or(
      this.page.locator('[class*="cursor-pointer"]')
    );
  }
  get leagueOfLegendsCard() {
    return this.page.getByText(/League of Legends/i);
  }
  get valorantCard() {
    return this.page.getByText(/Valorant/i);
  }
  get apexLegendsCard() {
    return this.page.getByText(/Apex Legends/i);
  }
  get overwatchCard() {
    return this.page.getByText(/Overwatch/i);
  }

  /** Featured players section */
  get featuredPlayersSection() {
    return this.page.getByText(/featured.*pros|推荐.*陪玩|top.*players/i);
  }
  get playerCards() {
    // Cards with player info in ScrollArea
    return this.page.locator('.card').filter({ hasText: /rating|orders|¥/i });
  }
  get skeletonLoaders() {
    return this.page.locator('.skeleton');
  }

  /** Trust indicators */
  get trustIndicatorsSection() {
    return this.page.getByText(/why choose|为什么选择/i);
  }
  get securePaymentsCard() {
    return this.page.getByText(/secure.*payments|安全.*支付/i);
  }
  get verifiedProsCard() {
    return this.page.getByText(/verified.*pros|认证.*陪玩/i);
  }
  get support247Card() {
    return this.page.getByText(/24.*7.*support|24.*7.*客服/i);
  }

  /** View all buttons */
  get viewAllGamesButton() {
    return this.page.getByRole('button', { name: /view.*all|查看全部/i });
  }
  get viewAllPlayersButton() {
    // Multiple "view all" buttons - get the one in featured section
    return this.page.locator('section').filter({ hasText: /featured/i }).getByRole('button', { name: /view.*all|查看全部/i });
  }

  /** Navigate to home page */
  async goto() {
    await this.page.goto('/');
  }

  /** Click Find Players button */
  async clickFindPlayers() {
    await this.findPlayersButton.click();
  }

  /** Click Become Player button */
  async clickBecomePlayer() {
    await this.becomePlayerButton.click();
  }

  /** Click a game card by name */
  async clickGame(gameName: string) {
    await this.page.getByText(gameName).click();
  }

  /** Click View All button for players */
  async clickViewAllPlayers() {
    await this.viewAllPlayersButton.click();
  }

  /** Verify page is loaded */
  async isLoaded() {
    const { expect } = await import('@playwright/test');
    // Wait for any major element to be visible
    await expect(
      this.page.locator('h1').or(
        this.page.locator('.card').or(
          this.page.locator('button')
        )
      ).first()
    ).toBeVisible({ timeout: 10000 });
  }
}
