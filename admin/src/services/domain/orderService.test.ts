/**
 * Property-Based Tests for OrderService
 *
 * Tests order cancellation, refund calculation, and statistics computation
 * using property-based testing with fast-check.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import { OrderService } from './orderService';
import type { Order } from '../../api/admin';

// Mock the admin API
vi.mock('../../api/admin', () => ({
  adminApi: {
    getOrders: vi.fn(),
    getOrder: vi.fn(),
    cancelOrder: vi.fn(),
    refundOrder: vi.fn(),
    batchCancelOrders: vi.fn(),
    batchCompleteOrders: vi.fn(),
  },
}));

/**
 * Arbitrary for generating valid Order objects
 */
const orderStatusArb = fc.constantFrom(
  'pending',
  'confirmed',
  'in_progress',
  'completed',
  'canceled',
  'refunded'
) as fc.Arbitrary<Order['status']>;

// Generate valid ISO date strings using timestamp range
const minTimestamp = new Date('2024-01-01').getTime();
const maxTimestamp = new Date('2026-12-31').getTime();
const dateStringArb = fc.integer({ min: minTimestamp, max: maxTimestamp })
  .map(ts => new Date(ts).toISOString());

const orderArb = fc.record({
  id: fc.integer({ min: 1, max: 100000 }),
  orderNo: fc.stringMatching(/^ORD\d{10}$/),
  userId: fc.integer({ min: 1, max: 10000 }),
  playerId: fc.integer({ min: 1, max: 10000 }),
  gameId: fc.integer({ min: 1, max: 100 }),
  title: fc.string({ minLength: 1, maxLength: 100 }),
  description: fc.string({ minLength: 0, maxLength: 500 }),
  totalPriceCents: fc.integer({ min: 100, max: 10000000 }), // 1 yuan to 100k yuan
  currency: fc.constant('CNY'),
  status: orderStatusArb,
  scheduledStart: dateStringArb,
  scheduledEnd: dateStringArb,
  completedAt: dateStringArb,
  cancelReason: fc.string({ minLength: 0, maxLength: 200 }),
  createdAt: dateStringArb,
  updatedAt: dateStringArb,
}) as fc.Arbitrary<Order>;

describe('OrderService - Property Tests', () => {
  let orderService: OrderService;

  beforeEach(() => {
    orderService = new OrderService();
    vi.clearAllMocks();
  });

  /**
   * **Feature: admin-phase3-improvements, Property 7: Order Cancellation State Validation**
   * **Validates: Requirements 3.2**
   *
   * For any order cancellation request, the service SHALL only allow cancellation
   * if the order status is in a cancellable state (pending, confirmed) and SHALL
   * reject with a clear reason otherwise.
   */
  describe('Property 7: Order Cancellation State Validation', () => {
    // Cancellable statuses
    const cancellableStatusArb = fc.constantFrom('pending', 'confirmed') as fc.Arbitrary<Order['status']>;
    
    // Non-cancellable statuses
    const nonCancellableStatusArb = fc.constantFrom(
      'in_progress',
      'completed',
      'canceled',
      'refunded'
    ) as fc.Arbitrary<Order['status']>;

    it('should allow cancellation for pending and confirmed orders', () => {
      fc.assert(
        fc.property(
          orderArb.chain(order => 
            cancellableStatusArb.map(status => ({ ...order, status }))
          ),
          (order) => {
            const result = orderService.canCancel(order);
            
            expect(result.allowed).toBe(true);
            expect(result.reason).toBeUndefined();
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should reject cancellation for non-cancellable orders with clear reason', () => {
      fc.assert(
        fc.property(
          orderArb.chain(order => 
            nonCancellableStatusArb.map(status => ({ ...order, status }))
          ),
          (order) => {
            const result = orderService.canCancel(order);
            
            expect(result.allowed).toBe(false);
            expect(result.reason).toBeDefined();
            expect(typeof result.reason).toBe('string');
            expect(result.reason!.length).toBeGreaterThan(0);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should provide specific reason for each non-cancellable status', () => {
      const statusReasonMap: Record<string, string[]> = {
        in_progress: ['in progress', 'cannot be canceled'],
        completed: ['completed', 'cannot be canceled'],
        canceled: ['already canceled'],
        refunded: ['refunded', 'cannot be canceled'],
      };

      for (const [status, expectedPhrases] of Object.entries(statusReasonMap)) {
        const order: Order = {
          id: 1,
          orderNo: 'ORD0000000001',
          userId: 1,
          playerId: 1,
          gameId: 1,
          title: 'Test Order',
          description: '',
          totalPriceCents: 10000,
          currency: 'CNY',
          status: status as Order['status'],
          scheduledStart: new Date().toISOString(),
          scheduledEnd: new Date().toISOString(),
          completedAt: new Date().toISOString(),
          cancelReason: '',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };

        const result = orderService.canCancel(order);
        
        expect(result.allowed).toBe(false);
        expect(result.reason).toBeDefined();
        
        // Check that reason contains at least one expected phrase
        const reasonLower = result.reason!.toLowerCase();
        const hasExpectedPhrase = expectedPhrases.some(phrase => 
          reasonLower.includes(phrase.toLowerCase())
        );
        expect(hasExpectedPhrase).toBe(true);
      }
    });

    it('canCancel result should be deterministic for same order', () => {
      fc.assert(
        fc.property(orderArb, (order) => {
          const result1 = orderService.canCancel(order);
          const result2 = orderService.canCancel(order);
          
          expect(result1.allowed).toBe(result2.allowed);
          expect(result1.reason).toBe(result2.reason);
        }),
        { numRuns: 100 }
      );
    });

    it('cancelOrder should fail for non-cancellable orders', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(
          orderArb.chain(order => 
            nonCancellableStatusArb.map(status => ({ ...order, status }))
          ),
          async (order) => {
            // Mock getOrder to return the order
            vi.mocked(adminApi.getOrder).mockResolvedValue({
              data: { success: true, data: order },
            } as never);

            const result = await orderService.cancelOrder(order.id, 'Test cancellation');
            
            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
            expect(result.error?.code).toBe('ORDER_CANNOT_CANCEL');
          }
        ),
        { numRuns: 50 }
      );
    });

    it('cancelOrder should succeed for cancellable orders', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(
          orderArb.chain(order => 
            cancellableStatusArb.map(status => ({ ...order, status }))
          ),
          async (order) => {
            const canceledOrder = { ...order, status: 'canceled' as const };
            
            // Mock getOrder to return the order
            vi.mocked(adminApi.getOrder).mockResolvedValue({
              data: { success: true, data: order },
            } as never);
            
            // Mock cancelOrder to return canceled order
            vi.mocked(adminApi.cancelOrder).mockResolvedValue({
              data: { success: true, data: canceledOrder },
            } as never);

            const result = await orderService.cancelOrder(order.id, 'Test cancellation');
            
            expect(result.success).toBe(true);
            expect(result.data?.status).toBe('canceled');
          }
        ),
        { numRuns: 50 }
      );
    });
  });


  /**
   * **Feature: admin-phase3-improvements, Property 8: Refund Calculation Accuracy**
   * **Validates: Requirements 3.3**
   *
   * For any refund calculation, the refund amount SHALL NOT exceed the original
   * payment amount, and the sum of refund amount, platform fee, and player amount
   * SHALL equal the original amount.
   */
  describe('Property 8: Refund Calculation Accuracy', () => {
    // Arbitrary for positive amounts
    const amountArb = fc.integer({ min: 100, max: 10000000 }); // 1 yuan to 100k yuan

    it('refund amount should never exceed original payment', () => {
      fc.assert(
        fc.property(
          orderArb,
          fc.integer({ min: 0, max: 20000000 }), // Request up to 200k yuan
          (order, requestedAmount) => {
            const result = orderService.calculateRefund(order, requestedAmount);
            
            expect(result.refundAmount).toBeLessThanOrEqual(result.originalAmount);
            expect(result.refundAmount).toBeLessThanOrEqual(order.totalPriceCents);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('refund amount should be non-negative', () => {
      fc.assert(
        fc.property(
          orderArb,
          fc.integer({ min: -1000000, max: 10000000 }), // Include negative requests
          (order, requestedAmount) => {
            const result = orderService.calculateRefund(order, requestedAmount);
            
            expect(result.refundAmount).toBeGreaterThanOrEqual(0);
            expect(result.platformFee).toBeGreaterThanOrEqual(0);
            expect(result.playerAmount).toBeGreaterThanOrEqual(0);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('platform fee plus player amount should equal refund amount', () => {
      fc.assert(
        fc.property(
          orderArb,
          amountArb,
          (order, requestedAmount) => {
            const result = orderService.calculateRefund(order, requestedAmount);
            
            // platformFee + playerAmount = refundAmount
            expect(result.platformFee + result.playerAmount).toBe(result.refundAmount);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('full refund should equal original amount', () => {
      fc.assert(
        fc.property(orderArb, (order) => {
          const result = orderService.calculateRefund(order, order.totalPriceCents);
          
          expect(result.refundAmount).toBe(order.totalPriceCents);
          expect(result.originalAmount).toBe(order.totalPriceCents);
        }),
        { numRuns: 100 }
      );
    });

    it('partial refund should be clamped to requested amount', () => {
      fc.assert(
        fc.property(
          orderArb.filter(o => o.totalPriceCents > 1000), // Ensure order has enough value
          fc.integer({ min: 100, max: 500 }), // Small partial refund
          (order, partialAmount) => {
            // Ensure partial amount is less than order total
            const requestedAmount = Math.min(partialAmount, order.totalPriceCents - 100);
            const result = orderService.calculateRefund(order, requestedAmount);
            
            expect(result.refundAmount).toBe(requestedAmount);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('refundOrder should validate amount against original payment', async () => {
      const { adminApi } = await import('../../api/admin');

      const order: Order = {
        id: 1,
        orderNo: 'ORD0000000001',
        userId: 1,
        playerId: 1,
        gameId: 1,
        title: 'Test Order',
        description: '',
        totalPriceCents: 10000, // 100 yuan
        currency: 'CNY',
        status: 'completed',
        scheduledStart: new Date().toISOString(),
        scheduledEnd: new Date().toISOString(),
        completedAt: new Date().toISOString(),
        cancelReason: '',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      // Mock getOrder
      vi.mocked(adminApi.getOrder).mockResolvedValue({
        data: { success: true, data: order },
      } as never);

      // Test refund exceeding original amount
      const result = await orderService.refundOrder(order.id, {
        reason: 'Test refund',
        amount_cents: 20000, // 200 yuan - exceeds original
      });

      expect(result.success).toBe(false);
      expect(result.error?.code).toBe('ORDER_INVALID_REFUND_AMOUNT');
    });

    it('refundOrder should reject zero or negative amounts', async () => {
      const { adminApi } = await import('../../api/admin');

      const order: Order = {
        id: 1,
        orderNo: 'ORD0000000001',
        userId: 1,
        playerId: 1,
        gameId: 1,
        title: 'Test Order',
        description: '',
        totalPriceCents: 10000,
        currency: 'CNY',
        status: 'completed',
        scheduledStart: new Date().toISOString(),
        scheduledEnd: new Date().toISOString(),
        completedAt: new Date().toISOString(),
        cancelReason: '',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      vi.mocked(adminApi.getOrder).mockResolvedValue({
        data: { success: true, data: order },
      } as never);

      // Test zero amount
      const zeroResult = await orderService.refundOrder(order.id, {
        reason: 'Test refund',
        amount_cents: 0,
      });
      expect(zeroResult.success).toBe(false);
      expect(zeroResult.error?.code).toBe('ORDER_INVALID_REFUND_AMOUNT');

      // Test negative amount
      const negativeResult = await orderService.refundOrder(order.id, {
        reason: 'Test refund',
        amount_cents: -100,
      });
      expect(negativeResult.success).toBe(false);
      expect(negativeResult.error?.code).toBe('ORDER_INVALID_REFUND_AMOUNT');
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 9: Order Statistics Computation Accuracy**
   * **Validates: Requirements 3.5**
   *
   * For any set of orders, the computed statistics SHALL accurately reflect
   * total revenue (sum of all completed order amounts), order counts by status,
   * and completion rate.
   */
  describe('Property 9: Order Statistics Computation Accuracy', () => {
    const ordersArb = fc.array(orderArb, { minLength: 0, maxLength: 100 });

    it('totalOrders should equal input array length', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          expect(stats.totalOrders).toBe(orders.length);
        }),
        { numRuns: 100 }
      );
    });

    it('ordersByStatus counts should sum to totalOrders', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          const statusSum = Object.values(stats.ordersByStatus).reduce((a, b) => a + b, 0);
          expect(statusSum).toBe(stats.totalOrders);
        }),
        { numRuns: 100 }
      );
    });

    it('totalRevenue should equal sum of completed order amounts', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          const expectedRevenue = orders
            .filter(o => o.status === 'completed')
            .reduce((sum, o) => sum + o.totalPriceCents, 0);
          
          expect(stats.totalRevenue).toBe(expectedRevenue);
        }),
        { numRuns: 100 }
      );
    });

    it('completionRate should be between 0 and 1', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          expect(stats.completionRate).toBeGreaterThanOrEqual(0);
          expect(stats.completionRate).toBeLessThanOrEqual(1);
        }),
        { numRuns: 100 }
      );
    });

    it('completionRate should equal completed orders / total orders', () => {
      fc.assert(
        fc.property(
          fc.array(orderArb, { minLength: 1, maxLength: 100 }), // At least 1 order
          (orders) => {
            const stats = orderService.computeStatistics(orders);
            
            const completedCount = orders.filter(o => o.status === 'completed').length;
            const expectedRate = completedCount / orders.length;
            
            expect(stats.completionRate).toBeCloseTo(expectedRate, 10);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('averageOrderValue should be calculated from completed orders only', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          const completedOrders = orders.filter(o => o.status === 'completed');
          
          if (completedOrders.length === 0) {
            expect(stats.averageOrderValue).toBe(0);
          } else {
            const expectedAvg = Math.round(
              completedOrders.reduce((sum, o) => sum + o.totalPriceCents, 0) / completedOrders.length
            );
            expect(stats.averageOrderValue).toBe(expectedAvg);
          }
        }),
        { numRuns: 100 }
      );
    });

    it('empty orders array should return zero statistics', () => {
      const stats = orderService.computeStatistics([]);
      
      expect(stats.totalOrders).toBe(0);
      expect(stats.totalRevenue).toBe(0);
      expect(stats.averageOrderValue).toBe(0);
      expect(stats.completionRate).toBe(0);
      expect(Object.keys(stats.ordersByStatus)).toHaveLength(0);
    });

    it('ordersByStatus should correctly count each status', () => {
      fc.assert(
        fc.property(ordersArb, (orders) => {
          const stats = orderService.computeStatistics(orders);
          
          // Manually count each status
          const expectedCounts: Record<string, number> = {};
          for (const order of orders) {
            expectedCounts[order.status] = (expectedCounts[order.status] || 0) + 1;
          }
          
          // Compare with computed stats
          for (const [status, count] of Object.entries(expectedCounts)) {
            expect(stats.ordersByStatus[status]).toBe(count);
          }
        }),
        { numRuns: 100 }
      );
    });
  });

  /**
   * Additional tests for batch operations
   */
  describe('Batch Operations', () => {
    const orderIdsArb = fc.array(fc.integer({ min: 1, max: 10000 }), {
      minLength: 1,
      maxLength: 20,
    });

    it('batchCancel should return result for each order', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(orderIdsArb, async (orderIds) => {
          // Mock getOrder to return cancellable orders
          vi.mocked(adminApi.getOrder).mockImplementation(async (id) => ({
            data: {
              success: true,
              data: {
                id,
                orderNo: `ORD${String(id).padStart(10, '0')}`,
                userId: 1,
                playerId: 1,
                gameId: 1,
                title: 'Test',
                description: '',
                totalPriceCents: 10000,
                currency: 'CNY',
                status: 'pending' as const,
                scheduledStart: new Date().toISOString(),
                scheduledEnd: new Date().toISOString(),
                completedAt: new Date().toISOString(),
                cancelReason: '',
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
              },
            },
          } as never));

          vi.mocked(adminApi.cancelOrder).mockResolvedValue({
            data: { success: true, data: {} as Order },
          } as never);

          const result = await orderService.batchCancel(orderIds, 'Batch cancel');

          expect(result.total).toBe(orderIds.length);
          expect(result.results.length).toBe(orderIds.length);
          expect(result.succeeded + result.failed).toBe(result.total);
        }),
        { numRuns: 30 }
      );
    });

    it('batchCancel with empty array should return empty result', async () => {
      const result = await orderService.batchCancel([], 'Test');

      expect(result.total).toBe(0);
      expect(result.succeeded).toBe(0);
      expect(result.failed).toBe(0);
      expect(result.results).toHaveLength(0);
      expect(result.success).toBe(true);
    });

    it('batchComplete should return result for each order', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(orderIdsArb, async (orderIds) => {
          vi.mocked(adminApi.batchCompleteOrders).mockResolvedValue({
            data: { success: true },
          } as never);

          const result = await orderService.batchComplete(orderIds);

          expect(result.total).toBe(orderIds.length);
          expect(result.results.length).toBe(orderIds.length);
        }),
        { numRuns: 30 }
      );
    });
  });

  /**
   * Tests for trend computation
   */
  describe('Trend Computation', () => {
    it('computeTrend should return correct number of days', () => {
      fc.assert(
        fc.property(
          fc.array(orderArb, { minLength: 0, maxLength: 50 }),
          fc.integer({ min: 1, max: 90 }),
          (orders, days) => {
            const trend = orderService.computeTrend(orders, days);
            
            expect(trend.length).toBe(days);
          }
        ),
        { numRuns: 50 }
      );
    });

    it('computeTrend should have non-negative values', () => {
      fc.assert(
        fc.property(
          fc.array(orderArb, { minLength: 0, maxLength: 50 }),
          fc.integer({ min: 1, max: 30 }),
          (orders, days) => {
            const trend = orderService.computeTrend(orders, days);
            
            for (const point of trend) {
              expect(point.value).toBeGreaterThanOrEqual(0);
            }
          }
        ),
        { numRuns: 50 }
      );
    });

    it('computeTrend dates should be in YYYY-MM-DD format', () => {
      fc.assert(
        fc.property(
          fc.array(orderArb, { minLength: 0, maxLength: 10 }),
          fc.integer({ min: 1, max: 30 }),
          (orders, days) => {
            const trend = orderService.computeTrend(orders, days);
            const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
            
            for (const point of trend) {
              expect(point.date).toMatch(dateRegex);
            }
          }
        ),
        { numRuns: 50 }
      );
    });
  });
});
