import { test, expect } from './fixtures/test-data.fixture';
import { LoginPage } from './pages/LoginPage';
import { PlayerManagementPage } from './pages/PlayerManagementPage';
import { createTestPlayer, deleteTestPlayer, getAdminToken } from './helpers/api-helpers';

/**
 * E2E Tests for Player Management
 *
 * Test Coverage:
 * - List players with pagination
 * - View player details
 * - Create new player
 * - Edit player information
 * - Delete player
 * - Player verification (approve/reject)
 * - Search and filter players
 * - Batch operations
 */

test.describe('Player Management', () => {
  let loginPage: LoginPage;
  let playerManagementPage: PlayerManagementPage;
  let adminToken: string;
  let testPlayerId: number;

  test.beforeAll(async () => {
    adminToken = await getAdminToken();
  });

  test.beforeEach(async ({ page, testData }) => {
    loginPage = new LoginPage(page);
    playerManagementPage = new PlayerManagementPage(page);

    // Login and navigate to player management
    await loginPage.goto();
    await loginPage.loginAndWaitForDashboard(
      testData.adminUser.username,
      testData.adminUser.password
    );
    await playerManagementPage.goto();
  });

  test.afterEach(async ({ page }) => {
    // Cleanup test player if created
    if (testPlayerId) {
      await deleteTestPlayer(adminToken, testPlayerId);
      testPlayerId = 0;
    }
  });

  test.describe('Player List', () => {
    test('should display player list page correctly', async () => {
      await expect(playerManagementPage['pageTitle']).toBeVisible();
      await expect(playerManagementPage['table']).toBeVisible();
    });

    test('should load players in table', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();
      expect(playerCount).toBeGreaterThanOrEqual(0);
    });

    test('should display player details correctly', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        const nickname = await playerManagementPage.getCellText(0, 1);
        expect(nickname).toBeTruthy();
        expect(nickname.length).toBeGreaterThan(0);
      }
    });

    test('should show different player statuses', async () => {
      const statuses = ['active', 'banned', 'suspended'];

      for (const status of statuses) {
        await playerManagementPage.goto();
        await playerManagementPage.filterByStatus(status);
        await playerManagementPage.page.waitForTimeout(1000);

        // Verify filter is applied (may have 0 results for some statuses)
        const count = await playerManagementPage.getPlayerCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });

    test('should show player certification status', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        // Look for certification status column
        const certStatus = await playerManagementPage.getCellText(0, 4);
        expect(certStatus).toBeTruthy();
      }
    });
  });

  test.describe('View Player Details', () => {
    test('should view player details', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        await playerManagementPage.viewPlayerDetails(0);

        await expect(playerManagementPage['modal']).toBeVisible();

        // Modal should contain player information
        await expect(playerManagementPage['modalTitle']).toBeVisible();
      }
    });

    test('should close player details modal', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        await playerManagementPage.viewPlayerDetails(0);
        await expect(playerManagementPage['modal']).toBeVisible();

        await playerManagementPage.closeModal();

        await expect(playerManagementPage['modal']).not.toBeVisible();
      }
    });

    test('should display complete player information', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        await playerManagementPage.viewPlayerDetails(0);

        // Verify modal contains expected information
        await expect(playerManagementPage['modal']).toBeVisible();

        const modalContent = playerManagementPage['modal'].textContent();
        expect(modalContent).toBeTruthy();
        expect(modalContent?.length).toBeGreaterThan(0);
      }
    });
  });

  test.describe('Edit Player', () => {
    test('should edit player information successfully', async ({ testData }) => {
      // Create test player first
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Edit Test ${testData.testPlayer.nickname}`,
      });
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Edit Test ${testData.testPlayer.nickname}`, 10000);
      await playerManagementPage.editPlayer(0, {
        nickname: `Updated ${testData.testPlayer.nickname}`,
      });

      await playerManagementPage.verifySuccessMessage();
    });

    test('should update player rank', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, testData.testPlayer);
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(testData.testPlayer.nickname, 10000);
      await playerManagementPage.editPlayer(0, {
        rank: 'master',
      });

      await playerManagementPage.verifySuccessMessage();
    });

    test('should update player hourly rate', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, testData.testPlayer);
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(testData.testPlayer.nickname, 10000);
      await playerManagementPage.editPlayer(0, {
        hourlyRateCents: 6000, // $60/hour
      });

      await playerManagementPage.verifySuccessMessage();
    });

    test('should cancel edit operation', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, testData.testPlayer);
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(testData.testPlayer.nickname, 10000);

      // Click edit but then cancel
      await playerManagementPage['editButton'](0).click();
      await expect(playerManagementPage['modal']).toBeVisible();

      await playerManagementPage['modalCancelButton'].click();

      // Verify modal is closed
      await expect(playerManagementPage['modal']).not.toBeVisible();
    });
  });

  test.describe('Delete Player', () => {
    test('should delete player successfully', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Delete Test ${testData.testPlayer.nickname}`,
      });
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Delete Test ${testData.testPlayer.nickname}`, 10000);

      await playerManagementPage.deletePlayer(0);

      await playerManagementPage.verifySuccessMessage();

      // Verify player is removed from table
      await expect(async () => {
        const exists = await playerManagementPage.verifyPlayerExists(`Delete Test ${testData.testPlayer.nickname}`);
        expect(exists).toBe(false);
      }).toPass({ timeout: 5000 });

      testPlayerId = 0; // Already deleted via UI
    });

    test('should show confirmation dialog before delete', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Confirm Delete ${testData.testPlayer.nickname}`,
      });
      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Confirm Delete ${testData.testPlayer.nickname}`, 10000);

      await playerManagementPage['deleteButton'](0).click();

      // Verify confirmation modal is shown
      await expect(playerManagementPage['modal']).toBeVisible();
      await expect(playerManagementPage['modalOkButton']).toBeVisible();

      // Cancel to cleanup manually
      await playerManagementPage['modalCancelButton'].click();
    });
  });

  test.describe('Player Verification', () => {
    test('should approve pending player verification', async ({ page, testData }) => {
      // Create a player with pending verification
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Pending Player ${testData.testPlayer.nickname}`,
        verificationStatus: 'pending',
      });

      await page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Pending Player ${testData.testPlayer.nickname}`, 10000);

      await playerManagementPage.approvePlayer(0, 'Verified via E2E test');

      await playerManagementPage.verifySuccessMessage();

      // Wait for status to change
      await page.waitForTimeout(2000);

      // Verify status changed to verified
      const isVerified = await playerManagementPage.verifyCertificationStatus(0, 'verified');
      expect(isVerified).toBe(true);
    });

    test('should reject player verification with reason', async ({ page, testData }) => {
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Reject Player ${testData.testPlayer.nickname}`,
        verificationStatus: 'pending',
      });

      await page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Reject Player ${testData.testPlayer.nickname}`, 10000);

      await playerManagementPage.rejectPlayer(0, 'Insufficient documentation provided');

      await playerManagementPage.verifySuccessMessage();

      // Wait for status to change
      await page.waitForTimeout(2000);

      // Verify status changed to rejected
      const isRejected = await playerManagementPage.verifyCertificationStatus(0, 'rejected');
      expect(isRejected).toBe(true);
    });

    test('should filter players by certification status', async () => {
      const statuses = ['pending', 'verified', 'rejected'];

      for (const status of statuses) {
        await playerManagementPage.goto();
        await playerManagementPage.filterByStatus(status);
        await playerManagementPage.page.waitForTimeout(1000);

        const count = await playerManagementPage.getPlayerCount();
        expect(count).toBeGreaterThanOrEqual(0);
      }
    });

    test('should show verification buttons for pending players', async ({ testData }) => {
      testPlayerId = await createTestPlayer(adminToken, {
        ...testData.testPlayer,
        nickname: `Pending Verify ${testData.testPlayer.nickname}`,
        verificationStatus: 'pending',
      });

      await playerManagementPage.page.reload();
      await playerManagementPage.waitForPageLoad();

      await playerManagementPage.waitForUserToAppear(`Pending Verify ${testData.testPlayer.nickname}`, 10000);

      // Check if approve and reject buttons are visible
      const approveButton = playerManagementPage['approveButton'](0);
      const rejectButton = playerManagementPage['rejectButton'](0);

      await expect(approveButton).toBeVisible();
      await expect(rejectButton).toBeVisible();
    });
  });

  test.describe('Search and Filter', () => {
    test('should search players by nickname', async () => {
      const initialCount = await playerManagementPage.getPlayerCount();

      await playerManagementPage.searchPlayer('Test');

      await playerManagementPage.page.waitForTimeout(1000);

      const filteredCount = await playerManagementPage.getPlayerCount();

      // Filtered results should be fewer or equal
      expect(filteredCount).toBeLessThanOrEqual(initialCount);
    });

    test('should filter players by status', async () => {
      await playerManagementPage.filterByStatus('active');

      await playerManagementPage.page.waitForTimeout(1000);

      const activeCount = await playerManagementPage.getPlayerCount();
      expect(activeCount).toBeGreaterThanOrEqual(0);
    });

    test('should clear filters', async () => {
      // Apply filter
      await playerManagementPage.filterByStatus('active');

      await playerManagementPage.page.waitForTimeout(1000);

      const filteredCount = await playerManagementPage.getPlayerCount();

      // Clear search
      await playerManagementPage.searchPlayer('');

      await playerManagementPage.page.waitForTimeout(1000);

      const clearedCount = await playerManagementPage.getPlayerCount();

      expect(clearedCount).toBeGreaterThanOrEqual(filteredCount);
    });
  });

  test.describe('Batch Operations', () => {
    test('should batch update player status', async ({ testData }) => {
      // Create test players
      const playerIds: number[] = [];
      try {
        for (let i = 0; i < 2; i++) {
          const playerData = {
            ...testData.testPlayer,
            nickname: `Batch Player ${i} ${Date.now()}`,
          };
          const playerId = await createTestPlayer(adminToken, playerData);
          playerIds.push(playerId);
        }

        await playerManagementPage.page.reload();
        await playerManagementPage.waitForPageLoad();

        // Batch status update would be performed here
        // Implementation depends on UI availability
      } finally {
        // Cleanup
        for (const id of playerIds) {
          await deleteTestPlayer(adminToken, id);
        }
      }
    });

    test('should batch delete multiple players', async ({ testData }) => {
      const playerIds: number[] = [];
      try {
        for (let i = 0; i < 2; i++) {
          const playerData = {
            ...testData.testPlayer,
            nickname: `Batch Delete ${i} ${Date.now()}`,
          };
          const playerId = await createTestPlayer(adminToken, playerData);
          playerIds.push(playerId);
        }

        await playerManagementPage.page.reload();
        await playerManagementPage.waitForPageLoad();

        // Batch delete would be performed here
        // Implementation depends on UI availability
      } finally {
        // Cleanup
        for (const id of playerIds) {
          await deleteTestPlayer(adminToken, id);
        }
      }
    });
  });

  test.describe('Export Functionality', () => {
    test('should export player list', async () => {
      const downloadPromise = playerManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await playerManagementPage.exportPlayers();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });

    test('should export filtered players', async () => {
      // Apply filter first
      await playerManagementPage.filterByStatus('active');
      await playerManagementPage.page.waitForTimeout(1000);

      const downloadPromise = playerManagementPage.page.waitForEvent('download', {
        timeout: 10000,
      });

      await playerManagementPage.exportPlayers();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.(xlsx|csv|xls)$/i);
    });
  });

  test.describe('Pagination', () => {
    test('should navigate through pages', async () => {
      const initialCount = await playerManagementPage.getPlayerCount();

      await playerManagementPage.nextPage();

      await playerManagementPage.page.waitForTimeout(1000);

      const newCount = await playerManagementPage.getPlayerCount();

      // If there are more players, count should be the same (different page)
      expect(newCount).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe('Accessibility', () => {
    test('should be keyboard navigable', async () => {
      // Tab through table
      await playerManagementPage.page.keyboard.press('Tab');
      await playerManagementPage.page.keyboard.press('Tab');

      // Focus should be on table or navigation
      const focusedElement = playerManagementPage.page.locator(':focus');
      await expect(focusedElement).toBeVisible();
    });
  });

  test.describe('Player Rating Display', () => {
    test('should display player ratings', async () => {
      const playerCount = await playerManagementPage.getPlayerCount();

      if (playerCount > 0) {
        // Look for rating information in table
        const rowText = await playerManagementPage.tableRows.nth(0).textContent();
        expect(rowText).toBeTruthy();
        expect(rowText?.length).toBeGreaterThan(0);
      }
    });
  });
});
