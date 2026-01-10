/**
 * Property-Based Tests for PlayerService
 *
 * Tests player verification workflow, earnings calculation, and statistics computation
 * using property-based testing with fast-check.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import { PlayerService, VERIFICATION_STATUSES } from './playerService';
import type { Player, Order, CommissionRule } from '../../api/admin';

// Mock the admin API
vi.mock('../../api/admin', () => ({
  adminApi: {
    getPlayers: vi.fn(),
    getPlayer: vi.fn(),
    createPlayer: vi.fn(),
    updatePlayer: vi.fn(),
    deletePlayer: vi.fn(),
    updatePlayerVerification: vi.fn(),
    updatePlayerSkillTags: vi.fn(),
    batchUpdatePlayerStatus: vi.fn(),
    batchDeletePlayers: vi.fn(),
  },
}));

/**
 * Helper to create a mock player
 */
function createMockPlayer(overrides: Partial<Player> = {}): Player {
  return {
    id: 1,
    userId: 1,
    nickname: 'TestPlayer',
    bio: 'Test bio',
    rank: 'gold',
    hourlyRateCents: 5000,
    mainGameId: 1,
    verificationStatus: 'pending',
    ratingAverage: 4.5,
    ratingCount: 10,
    skillTags: ['上分', '陪玩'],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    ...overrides,
  };
}

/**
 * Helper to create a mock order
 */
function createMockOrder(overrides: Partial<Order> = {}): Order {
  return {
    id: 1,
    orderNo: 'ORD001',
    userId: 1,
    playerId: 1,
    gameId: 1,
    title: 'Test Order',
    description: 'Test description',
    totalPriceCents: 10000,
    currency: 'CNY',
    status: 'completed',
    scheduledStart: new Date().toISOString(),
    scheduledEnd: new Date().toISOString(),
    completedAt: new Date().toISOString(),
    cancelReason: '',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    ...overrides,
  };
}

describe('PlayerService - Property Tests', () => {
  let playerService: PlayerService;

  beforeEach(() => {
    playerService = new PlayerService();
    vi.clearAllMocks();
  });

  /**
   * **Feature: admin-phase3-improvements, Property 10: Player Verification Workflow Enforcement**
   * **Validates: Requirements 4.2**
   *
   * For any player verification status change, the service SHALL enforce valid state
   * transitions (pending → verified/rejected) and SHALL reject invalid transitions.
   */
  describe('Property 10: Player Verification Workflow Enforcement', () => {
    // Arbitrary for verification statuses
    const verificationStatusArb = fc.constantFrom(...VERIFICATION_STATUSES);

    // Valid transitions map
    const validTransitions: Record<string, string[]> = {
      pending: ['verified', 'rejected'],
      verified: ['rejected'],
      rejected: ['pending'],
    };

    // Invalid transitions map
    const invalidTransitions: Record<string, string[]> = {
      pending: ['pending'], // Same status
      verified: ['pending', 'verified'], // Can't go back to pending or stay same
      rejected: ['verified', 'rejected'], // Can't go directly to verified or stay same
    };

    it('should allow valid state transitions', () => {
      fc.assert(
        fc.property(verificationStatusArb, (currentStatus) => {
          const player = createMockPlayer({ verificationStatus: currentStatus });
          const allowedNextStatuses = validTransitions[currentStatus];

          for (const nextStatus of allowedNextStatuses) {
            const result = playerService.canVerify(player, nextStatus);
            expect(result.allowed).toBe(true);
            expect(result.reason).toBeUndefined();
          }
        }),
        { numRuns: 100 }
      );
    });

    it('should reject invalid state transitions', () => {
      fc.assert(
        fc.property(verificationStatusArb, (currentStatus) => {
          const player = createMockPlayer({ verificationStatus: currentStatus });
          const disallowedNextStatuses = invalidTransitions[currentStatus];

          for (const nextStatus of disallowedNextStatuses) {
            const result = playerService.canVerify(player, nextStatus);
            expect(result.allowed).toBe(false);
            expect(result.reason).toBeDefined();
            expect(result.reason!.length).toBeGreaterThan(0);
          }
        }),
        { numRuns: 100 }
      );
    });

    it('should reject invalid verification status values', () => {
      const invalidStatusArb = fc.oneof(
        fc.constant('invalid'),
        fc.constant('approved'),
        fc.constant('active'),
        fc.constant(''),
        fc.stringMatching(/^[a-z]{5,10}$/).filter(
          (s) => !VERIFICATION_STATUSES.includes(s as typeof VERIFICATION_STATUSES[number])
        )
      );

      fc.assert(
        fc.property(verificationStatusArb, invalidStatusArb, (currentStatus, invalidStatus) => {
          const player = createMockPlayer({ verificationStatus: currentStatus });
          const result = playerService.canVerify(player, invalidStatus);

          expect(result.allowed).toBe(false);
          expect(result.reason).toBeDefined();
          expect(result.reason).toContain('Invalid verification status');
        }),
        { numRuns: 100 }
      );
    });

    it('should reject transition to same status', () => {
      fc.assert(
        fc.property(verificationStatusArb, (status) => {
          const player = createMockPlayer({ verificationStatus: status });
          const result = playerService.canVerify(player, status);

          expect(result.allowed).toBe(false);
          expect(result.reason).toBeDefined();
          expect(result.reason).toContain('already');
        }),
        { numRuns: 100 }
      );
    });

    it('pending -> verified should be allowed', () => {
      const player = createMockPlayer({ verificationStatus: 'pending' });
      const result = playerService.canVerify(player, 'verified');
      expect(result.allowed).toBe(true);
    });

    it('pending -> rejected should be allowed', () => {
      const player = createMockPlayer({ verificationStatus: 'pending' });
      const result = playerService.canVerify(player, 'rejected');
      expect(result.allowed).toBe(true);
    });

    it('verified -> rejected should be allowed (revocation)', () => {
      const player = createMockPlayer({ verificationStatus: 'verified' });
      const result = playerService.canVerify(player, 'rejected');
      expect(result.allowed).toBe(true);
    });

    it('verified -> pending should NOT be allowed', () => {
      const player = createMockPlayer({ verificationStatus: 'verified' });
      const result = playerService.canVerify(player, 'pending');
      expect(result.allowed).toBe(false);
    });

    it('rejected -> pending should be allowed (reapply)', () => {
      const player = createMockPlayer({ verificationStatus: 'rejected' });
      const result = playerService.canVerify(player, 'pending');
      expect(result.allowed).toBe(true);
    });

    it('rejected -> verified should NOT be allowed (must reapply first)', () => {
      const player = createMockPlayer({ verificationStatus: 'rejected' });
      const result = playerService.canVerify(player, 'verified');
      expect(result.allowed).toBe(false);
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 11: Player Earnings Calculation Accuracy**
   * **Validates: Requirements 4.3**
   *
   * For any player earnings calculation, the net earnings SHALL equal gross amount
   * minus commission, and commission SHALL be calculated using the applicable commission rate.
   */
  describe('Property 11: Player Earnings Calculation Accuracy', () => {
    // Arbitrary for order amounts (in cents)
    const orderAmountArb = fc.integer({ min: 100, max: 10000000 }); // 1 yuan to 100,000 yuan

    // Arbitrary for commission rates (0-100%)
    const commissionRateArb = fc.integer({ min: 0, max: 100 });

    it('net earnings should equal gross minus commission', () => {
      fc.assert(
        fc.property(orderAmountArb, commissionRateArb, (amount, ratePercent) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const commissionRule: CommissionRule = {
            id: 1,
            name: 'Test Rule',
            ratePercent,
            isDefault: true,
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };

          const result = playerService.calculateEarnings(order, commissionRule);

          // Net = Gross - Commission
          expect(result.netAmount).toBe(result.grossAmount - result.commissionAmount);
        }),
        { numRuns: 100 }
      );
    });

    it('commission should be calculated using the commission rate', () => {
      fc.assert(
        fc.property(orderAmountArb, commissionRateArb, (amount, ratePercent) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const commissionRule: CommissionRule = {
            id: 1,
            name: 'Test Rule',
            ratePercent,
            isDefault: true,
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };

          const result = playerService.calculateEarnings(order, commissionRule);

          // Commission rate should match
          expect(result.commissionRate).toBeCloseTo(ratePercent / 100, 5);

          // Commission amount should be approximately rate * gross (with rounding)
          const expectedCommission = Math.round(amount * (ratePercent / 100));
          expect(result.commissionAmount).toBe(expectedCommission);
        }),
        { numRuns: 100 }
      );
    });

    it('gross amount should equal order total price', () => {
      fc.assert(
        fc.property(orderAmountArb, (amount) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const result = playerService.calculateEarnings(order);

          expect(result.grossAmount).toBe(amount);
        }),
        { numRuns: 100 }
      );
    });

    it('should use default commission rate when no rule provided', () => {
      fc.assert(
        fc.property(orderAmountArb, (amount) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const result = playerService.calculateEarnings(order);

          // Default rate is 20%
          expect(result.commissionRate).toBe(0.20);
          expect(result.commissionAmount).toBe(Math.round(amount * 0.20));
        }),
        { numRuns: 100 }
      );
    });

    it('should use default rate when commission rule is inactive', () => {
      fc.assert(
        fc.property(orderAmountArb, commissionRateArb, (amount, ratePercent) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const commissionRule: CommissionRule = {
            id: 1,
            name: 'Test Rule',
            ratePercent,
            isDefault: true,
            status: 'inactive', // Inactive rule
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };

          const result = playerService.calculateEarnings(order, commissionRule);

          // Should use default rate (20%) when rule is inactive
          expect(result.commissionRate).toBe(0.20);
        }),
        { numRuns: 100 }
      );
    });

    it('all amounts should be non-negative', () => {
      fc.assert(
        fc.property(orderAmountArb, commissionRateArb, (amount, ratePercent) => {
          const order = createMockOrder({ totalPriceCents: amount });
          const commissionRule: CommissionRule = {
            id: 1,
            name: 'Test Rule',
            ratePercent,
            isDefault: true,
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };

          const result = playerService.calculateEarnings(order, commissionRule);

          expect(result.grossAmount).toBeGreaterThanOrEqual(0);
          expect(result.commissionAmount).toBeGreaterThanOrEqual(0);
          expect(result.netAmount).toBeGreaterThanOrEqual(0);
        }),
        { numRuns: 100 }
      );
    });

    it('commission rate should be clamped to 0-1 range', () => {
      // Test with rate > 100%
      const order = createMockOrder({ totalPriceCents: 10000 });
      const highRateRule: CommissionRule = {
        id: 1,
        name: 'High Rate',
        ratePercent: 150, // 150%
        isDefault: true,
        status: 'active',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const result = playerService.calculateEarnings(order, highRateRule);
      expect(result.commissionRate).toBeLessThanOrEqual(1);
      expect(result.netAmount).toBeGreaterThanOrEqual(0);
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 12: Player Statistics Computation Accuracy**
   * **Validates: Requirements 4.4**
   *
   * For any player statistics computation, the total earnings SHALL equal the sum of
   * net earnings from all completed orders, and rating average SHALL be correctly computed.
   */
  describe('Property 12: Player Statistics Computation Accuracy', () => {
    // Arbitrary for order amounts
    const orderAmountArb = fc.integer({ min: 100, max: 1000000 });

    // Arbitrary for order status
    const orderStatusArb = fc.constantFrom(
      'pending',
      'confirmed',
      'in_progress',
      'completed',
      'canceled',
      'refunded'
    ) as fc.Arbitrary<Order['status']>;

    // Arbitrary for player rating
    const ratingArb = fc.record({
      average: fc.float({ min: 0, max: 5, noNaN: true }),
      count: fc.integer({ min: 0, max: 10000 }),
    });

    it('total earnings should equal sum of net earnings from completed orders', () => {
      fc.assert(
        fc.property(
          fc.array(
            fc.record({
              amount: orderAmountArb,
              status: orderStatusArb,
            }),
            { minLength: 0, maxLength: 20 }
          ),
          (orderData) => {
            const playerId = 1;
            const player = createMockPlayer({ id: playerId });
            const orders = orderData.map((data, index) =>
              createMockOrder({
                id: index + 1,
                playerId,
                totalPriceCents: data.amount,
                status: data.status,
              })
            );

            const stats = playerService.computeStatistics(player, orders);

            // Calculate expected total earnings manually
            const completedOrders = orders.filter((o) => o.status === 'completed');
            let expectedTotalEarnings = 0;
            for (const order of completedOrders) {
              const earnings = playerService.calculateEarnings(order);
              expectedTotalEarnings += earnings.netAmount;
            }

            expect(stats.totalEarnings).toBe(expectedTotalEarnings);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('completed orders count should match actual completed orders', () => {
      fc.assert(
        fc.property(
          fc.array(
            fc.record({
              amount: orderAmountArb,
              status: orderStatusArb,
            }),
            { minLength: 0, maxLength: 20 }
          ),
          (orderData) => {
            const playerId = 1;
            const player = createMockPlayer({ id: playerId });
            const orders = orderData.map((data, index) =>
              createMockOrder({
                id: index + 1,
                playerId,
                totalPriceCents: data.amount,
                status: data.status,
              })
            );

            const stats = playerService.computeStatistics(player, orders);

            // Count completed orders manually
            const expectedCompletedCount = orders.filter((o) => o.status === 'completed').length;

            expect(stats.completedOrders).toBe(expectedCompletedCount);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('rating average and count should come from player data', () => {
      fc.assert(
        fc.property(ratingArb, (rating) => {
          const player = createMockPlayer({
            ratingAverage: rating.average,
            ratingCount: rating.count,
          });

          const stats = playerService.computeStatistics(player, []);

          expect(stats.averageRating).toBe(rating.average);
          expect(stats.ratingCount).toBe(rating.count);
        }),
        { numRuns: 100 }
      );
    });

    it('should only count orders belonging to the player', () => {
      fc.assert(
        fc.property(
          fc.array(orderAmountArb, { minLength: 1, maxLength: 10 }),
          fc.array(orderAmountArb, { minLength: 1, maxLength: 10 }),
          (playerOrderAmounts, otherOrderAmounts) => {
            const playerId = 1;
            const otherPlayerId = 2;
            const player = createMockPlayer({ id: playerId });

            // Create orders for the player
            const playerOrders = playerOrderAmounts.map((amount, index) =>
              createMockOrder({
                id: index + 1,
                playerId,
                totalPriceCents: amount,
                status: 'completed',
              })
            );

            // Create orders for another player
            const otherOrders = otherOrderAmounts.map((amount, index) =>
              createMockOrder({
                id: 100 + index,
                playerId: otherPlayerId,
                totalPriceCents: amount,
                status: 'completed',
              })
            );

            // Mix all orders
            const allOrders = [...playerOrders, ...otherOrders];

            const stats = playerService.computeStatistics(player, allOrders);

            // Should only count player's orders
            expect(stats.completedOrders).toBe(playerOrders.length);

            // Calculate expected earnings from player's orders only
            let expectedEarnings = 0;
            for (const order of playerOrders) {
              const earnings = playerService.calculateEarnings(order);
              expectedEarnings += earnings.netAmount;
            }
            expect(stats.totalEarnings).toBe(expectedEarnings);
          }
        ),
        { numRuns: 50 }
      );
    });

    it('empty orders should result in zero earnings and counts', () => {
      const player = createMockPlayer();
      const stats = playerService.computeStatistics(player, []);

      expect(stats.totalEarnings).toBe(0);
      expect(stats.monthlyEarnings).toBe(0);
      expect(stats.completedOrders).toBe(0);
    });

    it('non-completed orders should not contribute to earnings', () => {
      const nonCompletedStatuses: Order['status'][] = [
        'pending',
        'confirmed',
        'in_progress',
        'canceled',
        'refunded',
      ];

      fc.assert(
        fc.property(
          fc.constantFrom(...nonCompletedStatuses),
          orderAmountArb,
          (status, amount) => {
            const playerId = 1;
            const player = createMockPlayer({ id: playerId });
            const order = createMockOrder({
              playerId,
              totalPriceCents: amount,
              status,
            });

            const stats = playerService.computeStatistics(player, [order]);

            expect(stats.totalEarnings).toBe(0);
            expect(stats.completedOrders).toBe(0);
          }
        ),
        { numRuns: 50 }
      );
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 21: Skill Tag Parsing**
   * **Validates: Requirements 7.4**
   *
   * For any skill tags input string, the parser SHALL correctly split comma-separated
   * values and trim whitespace from each tag.
   */
  describe('Property 21: Skill Tag Parsing', () => {
    // Arbitrary for tag strings (without commas)
    const tagArb = fc.stringMatching(/^[^\s,][^,]*[^\s,]$|^[^\s,]$/).filter((s) => s.length > 0);

    it('should split comma-separated tags correctly', () => {
      fc.assert(
        fc.property(fc.array(tagArb, { minLength: 1, maxLength: 10 }), (tags) => {
          const tagsString = tags.join(',');
          const result = playerService.parseSkillTags(tagsString);

          expect(result.length).toBe(tags.length);
          for (let i = 0; i < tags.length; i++) {
            expect(result[i]).toBe(tags[i].trim());
          }
        }),
        { numRuns: 100 }
      );
    });

    it('should trim whitespace from each tag', () => {
      fc.assert(
        fc.property(fc.array(tagArb, { minLength: 1, maxLength: 10 }), (tags) => {
          // Add random whitespace around tags
          const tagsWithWhitespace = tags.map((tag) => `  ${tag}  `);
          const tagsString = tagsWithWhitespace.join(',');
          const result = playerService.parseSkillTags(tagsString);

          // Each result should be trimmed
          for (const tag of result) {
            expect(tag).toBe(tag.trim());
            expect(tag.startsWith(' ')).toBe(false);
            expect(tag.endsWith(' ')).toBe(false);
          }
        }),
        { numRuns: 100 }
      );
    });

    it('should filter out empty tags', () => {
      const result = playerService.parseSkillTags('tag1,,tag2,  ,tag3');
      expect(result).toEqual(['tag1', 'tag2', 'tag3']);
    });

    it('should handle empty string', () => {
      const result = playerService.parseSkillTags('');
      expect(result).toEqual([]);
    });

    it('should handle single tag without comma', () => {
      fc.assert(
        fc.property(tagArb, (tag) => {
          const result = playerService.parseSkillTags(tag);
          expect(result.length).toBe(1);
          expect(result[0]).toBe(tag.trim());
        }),
        { numRuns: 100 }
      );
    });

    it('should handle null/undefined gracefully', () => {
      expect(playerService.parseSkillTags(null as unknown as string)).toEqual([]);
      expect(playerService.parseSkillTags(undefined as unknown as string)).toEqual([]);
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 5: Batch Operation Result Completeness**
   * **Validates: Requirements 4.5**
   *
   * For any batch operation on players, the result SHALL contain an entry for
   * each input item indicating success or failure with details.
   */
  describe('Property 5 (Players): Batch Operation Result Completeness', () => {
    // Arbitrary for player IDs
    const playerIdsArb = fc.array(fc.integer({ min: 1, max: 10000 }), {
      minLength: 1,
      maxLength: 20,
    });

    it('batchDelete result should have entry for each input player', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(playerIdsArb, async (playerIds) => {
          // Mock API to succeed for all calls
          vi.mocked(adminApi.deletePlayer).mockResolvedValue({
            data: { success: true },
          } as never);

          const result = await playerService.batchDelete(playerIds);

          // Result should have entry for each input
          expect(result.total).toBe(playerIds.length);
          expect(result.results.length).toBe(playerIds.length);

          // succeeded + failed should equal total
          expect(result.succeeded + result.failed).toBe(result.total);

          // success flag should reflect whether all succeeded
          expect(result.success).toBe(result.failed === 0);
        }),
        { numRuns: 50 }
      );
    });

    it('batch operations with empty array should return empty result', async () => {
      const result = await playerService.batchDelete([]);

      expect(result.total).toBe(0);
      expect(result.succeeded).toBe(0);
      expect(result.failed).toBe(0);
      expect(result.results).toHaveLength(0);
      expect(result.success).toBe(true);
    });
  });
});
