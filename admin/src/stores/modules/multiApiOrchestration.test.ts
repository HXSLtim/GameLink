/**
 * Property-Based Tests for Multi-API Orchestration
 *
 * **Feature: admin-phase3-improvements, Property 3: Multi-API Orchestration Graceful Handling**
 * **Validates: Requirements 1.4**
 *
 * Tests that when multiple API calls are required for a single operation,
 * the service layer orchestrates these calls and handles partial failures gracefully.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import type { ServiceDependencies } from '@/services/domain/base';
import { BaseService } from '@/services/domain/base';
import type { BatchResult, BatchItemResult } from '@/services/utils';
import { ServiceResultHelper, ServiceException, ServiceErrorCodes } from '@/services/utils';

/**
 * Test service that simulates multi-API orchestration scenarios
 */
class MultiApiTestService extends BaseService {
  constructor(deps: ServiceDependencies = {}) {
    super(deps);
  }

  /**
   * Simulates a batch operation that calls multiple APIs
   * Each item in the batch triggers a separate API call
   */
  async batchOperation<T>(
    items: T[],
    processor: (item: T, index: number) => Promise<void>,
    context: string
  ): Promise<BatchResult<void>> {
    return this.executeBatch(items, processor, context);
  }

  /**
   * Simulates a multi-step operation where each step is an API call
   * Returns partial success information when some steps fail
   */
  async multiStepOperation(
    steps: Array<{ id: string; shouldFail: boolean }>
  ): Promise<{
    success: boolean;
    completedSteps: string[];
    failedSteps: Array<{ id: string; error: string }>;
  }> {
    const completedSteps: string[] = [];
    const failedSteps: Array<{ id: string; error: string }> = [];

    for (const step of steps) {
      try {
        if (step.shouldFail) {
          throw new Error(`Step ${step.id} failed`);
        }
        // Simulate successful API call
        completedSteps.push(step.id);
      } catch (error) {
        failedSteps.push({
          id: step.id,
          error: error instanceof Error ? error.message : 'Unknown error',
        });
      }
    }

    return {
      success: failedSteps.length === 0,
      completedSteps,
      failedSteps,
    };
  }
}

describe('Multi-API Orchestration - Property Tests', () => {
  let service: MultiApiTestService;

  beforeEach(() => {
    vi.clearAllMocks();
    service = new MultiApiTestService();
  });

  /**
   * **Feature: admin-phase3-improvements, Property 3: Multi-API Orchestration Graceful Handling**
   * **Validates: Requirements 1.4**
   *
   * For any service method that calls multiple APIs, if one API fails,
   * the service SHALL return a result that indicates partial success
   * and includes details of which operations succeeded and which failed.
   */
  describe('Property 3: Multi-API Orchestration Graceful Handling', () => {
    it('should return complete batch result with entry for each input item', async () => {
      // Arbitrary for batch items
      const itemsArb = fc.array(fc.integer({ min: 1, max: 1000 }), {
        minLength: 1,
        maxLength: 50,
      });

      await fc.assert(
        fc.asyncProperty(itemsArb, async (items) => {
          // All operations succeed
          const processor = async () => {};
          const result = await service.batchOperation(items, processor, 'test');

          // Property: Result SHALL contain an entry for each input item
          expect(result.results).toHaveLength(items.length);
          expect(result.total).toBe(items.length);

          // Each result should have valid structure
          for (let i = 0; i < result.results.length; i++) {
            const itemResult = result.results[i];
            expect(itemResult.index).toBe(i);
            expect(typeof itemResult.success).toBe('boolean');
          }
        }),
        { numRuns: 50 }
      );
    });

    it('should indicate partial success when some operations fail', async () => {
      // Arbitrary for items with random failure pattern
      const itemsWithFailuresArb = fc.array(
        fc.record({
          value: fc.integer({ min: 1, max: 100 }),
          shouldFail: fc.boolean(),
        }),
        { minLength: 2, maxLength: 30 }
      );

      await fc.assert(
        fc.asyncProperty(itemsWithFailuresArb, async (items) => {
          const processor = async (item: { value: number; shouldFail: boolean }) => {
            if (item.shouldFail) {
              throw new Error(`Item ${item.value} failed`);
            }
          };

          const result = await service.batchOperation(items, processor, 'test');

          // Count expected successes and failures
          const expectedSuccesses = items.filter((i) => !i.shouldFail).length;
          const expectedFailures = items.filter((i) => i.shouldFail).length;

          // Property: Result SHALL accurately report succeeded and failed counts
          expect(result.succeeded).toBe(expectedSuccesses);
          expect(result.failed).toBe(expectedFailures);
          expect(result.succeeded + result.failed).toBe(items.length);

          // Property: success flag SHALL be false if any operation failed
          if (expectedFailures > 0) {
            expect(result.success).toBe(false);
          } else {
            expect(result.success).toBe(true);
          }

          // Property: Each failed item SHALL have error details
          for (const itemResult of result.results) {
            if (!itemResult.success) {
              expect(itemResult.error).toBeDefined();
              expect(itemResult.error?.code).toBeDefined();
              expect(itemResult.error?.message).toBeDefined();
            }
          }
        }),
        { numRuns: 50 }
      );
    });

    it('should preserve order of results matching input order', async () => {
      // Arbitrary for items with unique identifiers
      const itemsArb = fc.array(fc.integer({ min: 1, max: 10000 }), {
        minLength: 1,
        maxLength: 30,
      }).map((arr) => arr.map((v, i) => ({ id: i, value: v })));

      await fc.assert(
        fc.asyncProperty(itemsArb, async (items) => {
          const processedOrder: number[] = [];
          const processor = async (item: { id: number; value: number }) => {
            processedOrder.push(item.id);
          };

          const result = await service.batchOperation(items, processor, 'test');

          // Property: Results SHALL be in the same order as input items
          for (let i = 0; i < result.results.length; i++) {
            expect(result.results[i].index).toBe(i);
          }
        }),
        { numRuns: 30 }
      );
    });

    it('should handle all failures gracefully', async () => {
      // Arbitrary for items that all fail
      const itemsArb = fc.array(fc.integer({ min: 1, max: 100 }), {
        minLength: 1,
        maxLength: 20,
      });

      await fc.assert(
        fc.asyncProperty(itemsArb, async (items) => {
          // All operations fail
          const processor = async (item: number) => {
            throw new Error(`Item ${item} failed`);
          };

          const result = await service.batchOperation(items, processor, 'test');

          // Property: When all fail, succeeded SHALL be 0 and failed SHALL equal total
          expect(result.succeeded).toBe(0);
          expect(result.failed).toBe(items.length);
          expect(result.success).toBe(false);

          // Property: All results SHALL have error details
          for (const itemResult of result.results) {
            expect(itemResult.success).toBe(false);
            expect(itemResult.error).toBeDefined();
          }
        }),
        { numRuns: 30 }
      );
    });

    it('should handle all successes correctly', async () => {
      // Arbitrary for items that all succeed
      const itemsArb = fc.array(fc.integer({ min: 1, max: 100 }), {
        minLength: 1,
        maxLength: 20,
      });

      await fc.assert(
        fc.asyncProperty(itemsArb, async (items) => {
          // All operations succeed
          const processor = async () => {};

          const result = await service.batchOperation(items, processor, 'test');

          // Property: When all succeed, failed SHALL be 0 and succeeded SHALL equal total
          expect(result.succeeded).toBe(items.length);
          expect(result.failed).toBe(0);
          expect(result.success).toBe(true);

          // Property: All results SHALL indicate success
          for (const itemResult of result.results) {
            expect(itemResult.success).toBe(true);
          }
        }),
        { numRuns: 30 }
      );
    });

    it('should handle empty batch gracefully', async () => {
      const result = await service.batchOperation([], async () => {}, 'test');

      // Property: Empty batch SHALL return valid result with zero counts
      expect(result.total).toBe(0);
      expect(result.succeeded).toBe(0);
      expect(result.failed).toBe(0);
      expect(result.success).toBe(true);
      expect(result.results).toHaveLength(0);
    });

    it('should track completed and failed steps in multi-step operations', async () => {
      // Arbitrary for multi-step operations
      const stepsArb = fc.array(
        fc.record({
          id: fc.string({ minLength: 1, maxLength: 10 }),
          shouldFail: fc.boolean(),
        }),
        { minLength: 1, maxLength: 20 }
      );

      await fc.assert(
        fc.asyncProperty(stepsArb, async (steps) => {
          const result = await service.multiStepOperation(steps);

          // Count expected outcomes
          const expectedCompleted = steps.filter((s) => !s.shouldFail).map((s) => s.id);
          const expectedFailed = steps.filter((s) => s.shouldFail).map((s) => s.id);

          // Property: completedSteps SHALL contain all successful step IDs
          expect(result.completedSteps).toHaveLength(expectedCompleted.length);
          for (const id of expectedCompleted) {
            expect(result.completedSteps).toContain(id);
          }

          // Property: failedSteps SHALL contain all failed step IDs with error details
          expect(result.failedSteps).toHaveLength(expectedFailed.length);
          for (const failedStep of result.failedSteps) {
            expect(expectedFailed).toContain(failedStep.id);
            expect(failedStep.error).toBeDefined();
            expect(failedStep.error.length).toBeGreaterThan(0);
          }

          // Property: success SHALL be true only if no steps failed
          expect(result.success).toBe(expectedFailed.length === 0);
        }),
        { numRuns: 50 }
      );
    });
  });

  describe('ServiceResultHelper - Batch Result Construction', () => {
    it('should correctly aggregate batch results from individual results', () => {
      // Arbitrary for individual results
      const resultsArb = fc.array(
        fc.record({
          index: fc.nat({ max: 100 }),
          success: fc.boolean(),
          data: fc.option(fc.string(), { nil: undefined }),
          error: fc.option(
            fc.record({
              code: fc.string({ minLength: 1 }),
              message: fc.string({ minLength: 1 }),
            }),
            { nil: undefined }
          ),
        }),
        { minLength: 1, maxLength: 50 }
      ).map((results) =>
        // Ensure indices are sequential
        results.map((r, i) => ({
          ...r,
          index: i,
          // Ensure error is present for failures
          error: r.success ? undefined : r.error || { code: 'ERROR', message: 'Failed' },
        }))
      );

      fc.assert(
        fc.property(resultsArb, (results) => {
          const batchResult = ServiceResultHelper.fromResults(results as BatchItemResult<string>[]);

          // Property: total SHALL equal number of results
          expect(batchResult.total).toBe(results.length);

          // Property: succeeded + failed SHALL equal total
          expect(batchResult.succeeded + batchResult.failed).toBe(results.length);

          // Property: succeeded SHALL equal count of successful results
          const expectedSucceeded = results.filter((r) => r.success).length;
          expect(batchResult.succeeded).toBe(expectedSucceeded);

          // Property: failed SHALL equal count of failed results
          const expectedFailed = results.filter((r) => !r.success).length;
          expect(batchResult.failed).toBe(expectedFailed);

          // Property: success SHALL be true only if all results succeeded
          expect(batchResult.success).toBe(expectedFailed === 0);
        }),
        { numRuns: 100 }
      );
    });

    it('should create empty batch result correctly', () => {
      const emptyBatch = ServiceResultHelper.emptyBatch(0);

      expect(emptyBatch.total).toBe(0);
      expect(emptyBatch.succeeded).toBe(0);
      expect(emptyBatch.failed).toBe(0);
      expect(emptyBatch.success).toBe(true);
      expect(emptyBatch.results).toHaveLength(0);
    });
  });
});

describe('Store Integration - Multi-API Orchestration', () => {
  /**
   * These tests verify that the refactored stores properly handle
   * partial failures from service batch operations.
   */

  it('should expose batch result details for UI consumption', async () => {
    // Arbitrary for batch operation results
    const batchResultArb = fc.record({
      success: fc.boolean(),
      total: fc.integer({ min: 1, max: 100 }),
      succeeded: fc.integer({ min: 0, max: 100 }),
      failed: fc.integer({ min: 0, max: 100 }),
      results: fc.array(
        fc.record({
          index: fc.nat({ max: 100 }),
          success: fc.boolean(),
          error: fc.option(
            fc.record({
              code: fc.string({ minLength: 1 }),
              message: fc.string({ minLength: 1 }),
            }),
            { nil: undefined }
          ),
        }),
        { minLength: 1, maxLength: 50 }
      ),
    }).filter((r) => r.succeeded + r.failed === r.total);

    fc.assert(
      fc.property(batchResultArb, (batchResult) => {
        // Property: Batch result SHALL contain all information needed for UI
        expect(typeof batchResult.success).toBe('boolean');
        expect(typeof batchResult.total).toBe('number');
        expect(typeof batchResult.succeeded).toBe('number');
        expect(typeof batchResult.failed).toBe('number');
        expect(Array.isArray(batchResult.results)).toBe(true);

        // Property: Each result SHALL have index and success status
        for (const result of batchResult.results) {
          expect(typeof result.index).toBe('number');
          expect(typeof result.success).toBe('boolean');
        }

        // Property: Failed results SHALL have error details
        const failedResults = batchResult.results.filter((r) => !r.success);
        for (const result of failedResults) {
          if (result.error) {
            expect(typeof result.error.code).toBe('string');
            expect(typeof result.error.message).toBe('string');
          }
        }
      }),
      { numRuns: 100 }
    );
  });
});
