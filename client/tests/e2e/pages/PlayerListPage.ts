/**
 * Player List Page Object Model
 */
export class PlayerListPage {
  constructor(private page: any) {}

  /** Page elements */
  get pageTitle() {
    return this.page.getByRole('heading', { name: /find players|找陪玩师/i });
  }
  get filterRating() {
    return this.page.getByText(/rating|评分/i);
  }
  get filterOrders() {
    return this.page.getByText(/orders|订单/i);
  }
  get filterPrice() {
    return this.page.getByText(/price|价格/i);
  }
  get playerCards() {
    return this.page.locator('[class*="player"]').or(this.page.locator('.card'));
  }
  get emptyState() {
    return this.page.getByText(/no players|没有找到/i);
  }

  /** Navigate to players page */
  async goto() {
    await this.page.goto('/players');
  }

  /** Click a player card by index */
  async clickPlayerCard(index: number) {
    const cards = await this.playerCards.all();
    await cards.nth(index).click();
  }

  /** Filter by rating */
  async filterByRating() {
    await this.filterRating.click();
  }

  /** Filter by orders */
  async filterByOrders() {
    await this.filterOrders.click();
  }

  /** Verify page is loaded */
  async isLoaded() {
    const { expect } = await import('@playwright/test');
    await expect(this.pageTitle).toBeVisible();
  }
}
