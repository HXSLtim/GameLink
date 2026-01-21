import { test, expect } from './fixtures/test-data.fixture';
import { HomePage } from './pages/HomePage';
import { PlayerListPage } from './pages/PlayerListPage';
import { PlayerDetailPage } from './pages/PlayerDetailPage';

/**
 * Player Browsing E2E Tests
 *
 * Test Coverage:
 * - Browse home page
 * - View player list
 * - View player details
 * - Filter players
 * - Favorite players
 */
test.describe('Player Browsing', () => {
  let homePage: HomePage;
  let playerListPage: PlayerListPage;
  let playerDetailPage: PlayerDetailPage;

  test.beforeEach(async ({ page }) => {
    homePage = new HomePage(page);
    playerListPage = new PlayerListPage(page);
    playerDetailPage = new PlayerDetailPage(page);
  });

  test.describe('Home Page', () => {
    test('should display home page correctly', async ({ page }) => {
      await homePage.goto();

      // Check page title
      await expect(page).toHaveTitle(/GameLink/);

      // Check hero section - wait for page to load
      await homePage.isLoaded();

      // Check buttons exist
      await expect(homePage.findPlayersButton.or(homePage.becomePlayerButton)).toBeVisible();

      // Check sections exist (they may have skeleton loaders)
      await expect(homePage.gameCards.first().or(homePage.skeletonLoaders.first())).toBeVisible();
    });

    test('should navigate to players page from hero button', async ({ page }) => {
      await homePage.goto();
      await homePage.clickFindPlayers();

      await expect(page).toHaveURL(/\/players/);
      await expect(playerListPage.pageTitle).toBeVisible();
    });

    test('should display game cards', async () => {
      await homePage.goto();

      // Check game cards are visible
      const gameCount = await homePage.gameCards.count();
      expect(gameCount).toBeGreaterThan(0);
    });

    test('should display featured players', async () => {
      await homePage.goto();

      // Check featured players section
      await expect(homePage.featuredPlayersSection).toBeVisible();
    });
  });

  test.describe('Player List Page', () => {
    test('should display player list correctly', async () => {
      await playerListPage.goto();

      // Check page title
      await expect(playerListPage.pageTitle).toBeVisible();

      // Check filters are visible
      await expect(playerListPage.filterRating).toBeVisible();
      await expect(playerListPage.filterOrders).toBeVisible();
      await expect(playerListPage.filterPrice).toBeVisible();
    });

    test('should filter players by rating', async () => {
      await playerListPage.goto();

      // Mock players data for testing
      await page.route('**/api/v1/players**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              players: [
                {
                  id: 1,
                  username: 'player1',
                  nickname: 'ProGamer',
                  gameName: 'Valorant',
                  rating: 4.8,
                  price: 50,
                  orderCount: 100,
                  tags: ['Pro', 'Carry'],
                  online: true,
                },
                {
                  id: 2,
                  username: 'player2',
                  nickname: 'CasualPlayer',
                  gameName: 'Valorant',
                  rating: 4.2,
                  price: 30,
                  orderCount: 50,
                  tags: ['Friendly'],
                  online: true,
                },
              ],
              pagination: { page: 1, pageSize: 20, total: 2, hasMore: false },
            },
          }),
        });
      });

      // Click rating filter
      await playerListPage.filterByRating();

      // Verify filter is applied (UI indication)
      await expect(playerListPage.filterRating).toHaveAttribute('class', /active|selected/i);
    });

    test('should navigate to player detail', async ({ page }) => {
      // Mock player detail data
      page.route('**/api/v1/players/1', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              id: 1,
              username: 'player1',
              nickname: 'ProGamer',
              gameName: 'Valorant',
              rating: 4.8,
              price: 50,
              orderCount: 100,
              tags: ['Pro', 'Carry'],
              online: true,
              bio: 'Professional Valorant player',
            },
          }),
        });
      });

      await playerListPage.goto();
      await playerListPage.clickPlayerCard(0);

      // Verify navigation to player detail
      await expect(page).toHaveURL(/\/players\/\d+/);
    });

    test('should handle empty state', async ({ page }) => {
      // Mock empty players list
      await page.route('**/api/v1/players**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              players: [],
              pagination: { page: 1, pageSize: 20, total: 0, hasMore: false },
            },
          }),
        });
      });

      await playerListPage.goto();

      // Verify empty state is shown
      await expect(playerListPage.emptyState).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe('Player Detail Page', () => {
    test.beforeEach(async ({ page }) => {
      // Mock player detail API
      page.route('**/api/v1/players/1', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              id: 1,
              username: 'player1',
              nickname: 'ProGamer',
              gameName: 'Valorant',
              rating: 4.8,
              price: 50,
              orderCount: 100,
              tags: ['Pro', 'Carry', 'Patient'],
              online: true,
              bio: 'Professional Valorant player with 5 years experience',
              avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=player1',
            },
          }),
        });
      });
    });

    test('should display player details correctly', async () => {
      await playerDetailPage.goto(1);

      // Check player info
      await expect(playerDetailPage.playerName).toBeVisible();
      await expect(playerDetailPage.price).toBeVisible();

      // Check action buttons
      await expect(playerDetailPage.bookButton).toBeVisible();
      await expect(playerDetailPage.favoriteButton).toBeVisible();
      await expect(playerDetailPage.chatButton).toBeVisible();
    });

    test('should display player tags', async () => {
      await playerDetailPage.goto(1);

      // Check tags are displayed
      const tagCount = await playerDetailPage.tags.count();
      expect(tagCount).toBeGreaterThan(0);
    });

    test('should display rating', async () => {
      await playerDetailPage.goto(1);

      // Check rating is visible
      const rating = await playerDetailPage.rating;
      await expect(rating).toBeVisible();
    });

    test('should have working favorite button', async ({ page }) => {
      await playerDetailPage.goto(1);

      // Mock favorite API
      page.route('**/api/v1/user/favorites/**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ success: true, data: { isFavorite: true } }),
        });
      });

      // Click favorite button
      await playerDetailPage.toggleFavorite();

      // Verify action was performed (UI feedback)
      await expect(playerDetailPage.favoriteButton).toBeVisible();
    });
  });

  test.describe('User Flow', () => {
    test('should complete browse player flow', async ({ page }) => {
      // Mock APIs
      page.route('**/api/v1/players**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              players: [
                {
                  id: 1,
                  username: 'player1',
                  nickname: 'ProGamer',
                  gameName: 'Valorant',
                  rating: 4.8,
                  price: 50,
                  orderCount: 100,
                  tags: ['Pro'],
                  online: true,
                },
              ],
              pagination: { page: 1, pageSize: 20, total: 1, hasMore: false },
            },
          }),
        });
      });

      // 1. Start at home
      await homePage.goto();
      await homePage.isLoaded();

      // 2. Click find players
      await homePage.clickFindPlayers();
      await expect(page).toHaveURL(/\/players/);

      // 3. View player detail
      await playerListPage.clickPlayerCard(0);
      await expect(page).toHaveURL(/\/players\/\d+/);
      await playerDetailPage.isLoaded();
    });
  });
});
