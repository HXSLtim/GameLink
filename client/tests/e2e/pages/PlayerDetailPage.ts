/**
 * Player Detail Page Object Model
 */
export class PlayerDetailPage {
  constructor(private page: any) {}

  /** Player info elements */
  get playerName() {
    return this.page.getByRole('heading', { level: 1 });
  }
  get rating() {
    return this.page.getByText(/rating|评分/i).or(this.page.locator('[class*="star"]'));
  }
  get price() {
    return this.page.getByText(/\s*¥\s*\d+/); // Price with yuan symbol
  }
  get tags() {
    return this.page.locator('.badge').or(this.page.locator('[class*="tag"]'));
  }
  get bookButton() {
    return this.page.getByRole('button', { name: /book|预订|下单/i });
  }
  get chatButton() {
    return this.page.getByRole('button', { name: /chat|聊天/i });
  }
  get favoriteButton() {
    return this.page.locator('[class*="favorite"]');
  }

  /** Reviews section */
  get reviewsSection() {
    return this.page.getByText(/reviews|评价/i);
  }
  get reviewCards() {
    return this.page.locator('[class*="review"]');
  }

  /** Navigate to player detail page */
  async goto(playerId: number) {
    await this.page.goto(`/players/${playerId}`);
  }

  /** Click book button */
  async clickBook() {
    await this.bookButton.click();
  }

  /** Click favorite button */
  async toggleFavorite() {
    await this.favoriteButton.click();
  }

  /** Click chat button */
  async clickChat() {
    await this.chatButton.click();
  }

  /** Verify page is loaded */
  async isLoaded() {
    const { expect } = await import('@playwright/test');
    await expect(this.playerName).toBeVisible();
  }

  /** Get player name */
  async getPlayerName(): Promise<string> {
    return await this.playerName.textContent();
  }

  /** Get price value */
  async getPrice(): Promise<string> {
    const priceText = await this.price.textContent();
    return priceText?.match(/¥(\d+)/)?.[1] || '0';
  }
}
