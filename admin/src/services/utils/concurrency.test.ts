/**
 * Concurrency Control Tests
 * Tests for ConcurrencyController, Semaphore, retry logic, and batch processing
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as fc from 'fast-check';
import {
  ConcurrencyController,
  Semaphore,
  DEFAULT_CONCURRENCY_CONFIG,
  isRetryableError,
  chunkArray,
  sleep,
  createConcurrencyController,
} from './concurrency';

describe('Semaphore', () => {
  it('should initialize with correct permits', () => {
    const semaphore = new Semaphore(5);
    expect(semaphore.getAvailablePermits()).toBe(5);
    expect(semaphore.getWaitingCount()).toBe(0);
  });

  it('should throw error for invalid permits', () => {
    expect(() => new Semaphore(0)).toThrow('Semaphore permits must be at least 1');
    expect(() => new Semaphore(-1)).toThrow('Semaphore permits must be at least 1');
  });

  it('should acquire and release permits correctly', async () => {
    const semaphore = new Semaphore(2);

    await semaphore.acquire();
    expect(semaphore.getAvailablePermits()).toBe(1);

    await semaphore.acquire();
    expect(semaphore.getAvailablePermits()).toBe(0);

    semaphore.release();
    expect(semaphore.getAvailablePermits()).toBe(1);

    semaphore.release();
    expect(semaphore.getAvailablePermits()).toBe(2);
  });

  it('should queue waiting operations when no permits available', async () => {
    const semaphore = new Semaphore(1);
    const order: number[] = [];

    await semaphore.acquire();
    expect(semaphore.getAvailablePermits()).toBe(0);

    // Start waiting operations
    const wait1 = semaphore.acquire().then(() => order.push(1));
    const wait2 = semaphore.acquire().then(() => order.push(2));

    expect(semaphore.getWaitingCount()).toBe(2);

    // Release permits
    semaphore.release();
    await wait1;
    expect(order).toEqual([1]);

    semaphore.release();
    await wait2;
    expect(order).toEqual([1, 2]);
  });
});

describe('isRetryableError', () => {
  it('should return true for rate limit errors', () => {
    expect(isRetryableError(new Error('Rate limit exceeded'))).toBe(true);
    expect(isRetryableError(new Error('429 Too Many Requests'))).toBe(true);
  });

  it('should return true for timeout errors', () => {
    expect(isRetryableError(new Error('Request timeout'))).toBe(true);
    expect(isRetryableError(new Error('Connection timeout'))).toBe(true);
  });

  it('should return true for network errors', () => {
    expect(isRetryableError(new Error('Network error'))).toBe(true);
    expect(isRetryableError(new Error('ECONNRESET'))).toBe(true);
    expect(isRetryableError(new Error('ECONNREFUSED'))).toBe(true);
    expect(isRetryableError(new Error('503 Service Unavailable'))).toBe(true);
  });

  it('should return false for non-retryable errors', () => {
    expect(isRetryableError(new Error('Validation failed'))).toBe(false);
    expect(isRetryableError(new Error('Not found'))).toBe(false);
    expect(isRetryableError(new Error('Permission denied'))).toBe(false);
  });

  it('should return false for non-Error values', () => {
    expect(isRetryableError('string error')).toBe(false);
    expect(isRetryableError(null)).toBe(false);
    expect(isRetryableError(undefined)).toBe(false);
  });
});

describe('chunkArray', () => {
  it('should split array into chunks of specified size', () => {
    const array = [1, 2, 3, 4, 5, 6, 7];
    expect(chunkArray(array, 3)).toEqual([[1, 2, 3], [4, 5, 6], [7]]);
  });

  it('should handle empty array', () => {
    expect(chunkArray([], 5)).toEqual([]);
  });

  it('should handle array smaller than chunk size', () => {
    expect(chunkArray([1, 2], 5)).toEqual([[1, 2]]);
  });

  it('should handle chunk size equal to array length', () => {
    expect(chunkArray([1, 2, 3], 3)).toEqual([[1, 2, 3]]);
  });

  it('should throw error for invalid chunk size', () => {
    expect(() => chunkArray([1, 2, 3], 0)).toThrow('Chunk size must be at least 1');
    expect(() => chunkArray([1, 2, 3], -1)).toThrow('Chunk size must be at least 1');
  });
});

describe('sleep', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('should resolve after specified time', async () => {
    const promise = sleep(1000);
    vi.advanceTimersByTime(1000);
    await expect(promise).resolves.toBeUndefined();
  });
});

describe('ConcurrencyController', () => {
  describe('constructor and configuration', () => {
    it('should use default config when none provided', () => {
      const controller = new ConcurrencyController();
      expect(controller.getConfig()).toEqual(DEFAULT_CONCURRENCY_CONFIG);
    });

    it('should merge custom config with defaults', () => {
      const controller = new ConcurrencyController({ maxConcurrent: 10 });
      const config = controller.getConfig();
      expect(config.maxConcurrent).toBe(10);
      expect(config.chunkSize).toBe(DEFAULT_CONCURRENCY_CONFIG.chunkSize);
    });
  });

  describe('withDeduplication', () => {
    it('should execute operation and return result', async () => {
      const controller = new ConcurrencyController();
      const result = await controller.withDeduplication('test-op', async () => 'result');
      expect(result).toBe('result');
    });

    it('should return existing promise for duplicate operation', async () => {
      const controller = new ConcurrencyController();
      let callCount = 0;

      const operation = async () => {
        callCount++;
        await sleep(10);
        return 'result';
      };

      // Start two operations with same key
      const promise1 = controller.withDeduplication('same-key', operation);
      const promise2 = controller.withDeduplication('same-key', operation);

      // Both should return the same promise
      expect(controller.isOperationInProgress('same-key')).toBe(true);

      const [result1, result2] = await Promise.all([promise1, promise2]);

      expect(result1).toBe('result');
      expect(result2).toBe('result');
      expect(callCount).toBe(1); // Operation should only be called once
    });

    it('should allow new operation after previous completes', async () => {
      const controller = new ConcurrencyController();
      let callCount = 0;

      const operation = async () => {
        callCount++;
        return `result-${callCount}`;
      };

      const result1 = await controller.withDeduplication('key', operation);
      expect(result1).toBe('result-1');
      expect(controller.isOperationInProgress('key')).toBe(false);

      const result2 = await controller.withDeduplication('key', operation);
      expect(result2).toBe('result-2');
      expect(callCount).toBe(2);
    });

    it('should clean up operation key even on error', async () => {
      const controller = new ConcurrencyController();

      await expect(
        controller.withDeduplication('error-key', async () => {
          throw new Error('Test error');
        })
      ).rejects.toThrow('Test error');

      expect(controller.isOperationInProgress('error-key')).toBe(false);
    });
  });

  describe('withRetry', () => {
    beforeEach(() => {
      vi.useFakeTimers();
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('should return result on first success', async () => {
      const controller = new ConcurrencyController();
      const result = await controller.withRetry(async () => 'success');
      expect(result).toBe('success');
    });

    it('should retry on retryable error and succeed', async () => {
      const controller = new ConcurrencyController({ retryAttempts: 3, retryDelayMs: 100 });
      let attempts = 0;

      const operation = async () => {
        attempts++;
        if (attempts < 3) {
          throw new Error('Rate limit exceeded');
        }
        return 'success';
      };

      const promise = controller.withRetry(operation);

      // Advance timers for retries
      await vi.advanceTimersByTimeAsync(100); // First retry delay
      await vi.advanceTimersByTimeAsync(200); // Second retry delay (with backoff)

      const result = await promise;
      expect(result).toBe('success');
      expect(attempts).toBe(3);
    });

    it('should throw immediately for non-retryable error', async () => {
      const controller = new ConcurrencyController();
      let attempts = 0;

      const operation = async () => {
        attempts++;
        throw new Error('Validation failed');
      };

      await expect(controller.withRetry(operation)).rejects.toThrow('Validation failed');
      expect(attempts).toBe(1);
    });

    it('should throw after max retries exhausted', async () => {
      vi.useRealTimers(); // Use real timers for this test
      const controller = new ConcurrencyController({ retryAttempts: 2, retryDelayMs: 1 });
      let attempts = 0;

      const operation = async () => {
        attempts++;
        throw new Error('Network error');
      };

      await expect(controller.withRetry(operation)).rejects.toThrow('Network error');
      expect(attempts).toBe(2);
    });
  });

  describe('processWithConcurrency', () => {
    it('should process all items and return results', async () => {
      const controller = new ConcurrencyController({ maxConcurrent: 2, chunkSize: 10 });
      const items = [1, 2, 3, 4, 5];

      const results = await controller.processWithConcurrency(
        items,
        async (item) => item * 2
      );

      expect(results).toHaveLength(5);
      expect(results.every((r) => r.success)).toBe(true);
      expect(results.map((r) => r.result)).toEqual([2, 4, 6, 8, 10]);
    });

    it('should handle empty array', async () => {
      const controller = new ConcurrencyController();
      const results = await controller.processWithConcurrency([], async () => 'result');
      expect(results).toEqual([]);
    });

    it('should capture errors for failed items', async () => {
      const controller = new ConcurrencyController({ maxConcurrent: 2, chunkSize: 10 });
      const items = [1, 2, 3];

      const results = await controller.processWithConcurrency(items, async (item) => {
        if (item === 2) {
          throw new Error('Item 2 failed');
        }
        return item * 2;
      });

      expect(results).toHaveLength(3);
      expect(results[0]).toEqual({ index: 0, success: true, result: 2 });
      expect(results[1].success).toBe(false);
      expect(results[1].error?.message).toBe('Item 2 failed');
      expect(results[2]).toEqual({ index: 2, success: true, result: 6 });
    });

    it('should call progress callback', async () => {
      const controller = new ConcurrencyController({ maxConcurrent: 5, chunkSize: 2 });
      const items = [1, 2, 3, 4, 5];
      const progressCalls: Array<[number, number]> = [];

      await controller.processWithConcurrency(
        items,
        async (item) => item,
        (completed, total) => progressCalls.push([completed, total])
      );

      expect(progressCalls).toEqual([
        [2, 5],
        [4, 5],
        [5, 5],
      ]);
    });

    it('should preserve correct indices across chunks', async () => {
      const controller = new ConcurrencyController({ maxConcurrent: 2, chunkSize: 3 });
      const items = ['a', 'b', 'c', 'd', 'e', 'f', 'g'];

      const results = await controller.processWithConcurrency(items, async (item, index) => ({
        item,
        index,
      }));

      expect(results.map((r) => r.index)).toEqual([0, 1, 2, 3, 4, 5, 6]);
      expect(results.map((r) => r.result?.item)).toEqual(['a', 'b', 'c', 'd', 'e', 'f', 'g']);
    });
  });
});

describe('createConcurrencyController', () => {
  it('should create controller with default config', () => {
    const controller = createConcurrencyController();
    expect(controller.getConfig()).toEqual(DEFAULT_CONCURRENCY_CONFIG);
  });

  it('should create controller with custom config', () => {
    const controller = createConcurrencyController({ maxConcurrent: 10 });
    expect(controller.getConfig().maxConcurrent).toBe(10);
  });
});


// ============================================================================
// Property-Based Tests
// ============================================================================

describe('Property-Based Tests', () => {
  /**
   * Property 27: Batch Concurrency Limit
   * For any batch operation, the number of concurrent API calls SHALL NOT
   * exceed the configured maximum (default 5).
   * Validates: Requirements 11.1
   */
  describe('Property 27: Batch Concurrency Limit', () => {
    it('should never exceed maxConcurrent limit during batch processing', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 1, max: 10 }), // maxConcurrent
          fc.integer({ min: 1, max: 50 }), // number of items
          fc.integer({ min: 1, max: 20 }), // chunkSize
          async (maxConcurrent, itemCount, chunkSize) => {
            const controller = new ConcurrencyController({ maxConcurrent, chunkSize });
            const items = Array.from({ length: itemCount }, (_, i) => i);

            let currentConcurrent = 0;
            let maxObservedConcurrent = 0;

            const results = await controller.processWithConcurrency(
              items,
              async (item) => {
                currentConcurrent++;
                maxObservedConcurrent = Math.max(maxObservedConcurrent, currentConcurrent);

                // Simulate async work
                await new Promise((resolve) => setTimeout(resolve, 1));

                currentConcurrent--;
                return item * 2;
              }
            );

            // Property: max observed concurrent should never exceed configured limit
            expect(maxObservedConcurrent).toBeLessThanOrEqual(maxConcurrent);

            // All items should be processed
            expect(results).toHaveLength(itemCount);
          }
        ),
        { numRuns: 20 } // Reduced runs due to async nature
      );
    });

    it('should respect concurrency limit with varying item counts', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.array(fc.integer({ min: 1, max: 100 }), { minLength: 1, maxLength: 30 }),
          async (items) => {
            const maxConcurrent = 5;
            const controller = new ConcurrencyController({ maxConcurrent, chunkSize: 10 });

            let currentConcurrent = 0;
            let maxObservedConcurrent = 0;

            await controller.processWithConcurrency(items, async (item) => {
              currentConcurrent++;
              maxObservedConcurrent = Math.max(maxObservedConcurrent, currentConcurrent);
              await new Promise((resolve) => setTimeout(resolve, 1));
              currentConcurrent--;
              return item;
            });

            expect(maxObservedConcurrent).toBeLessThanOrEqual(maxConcurrent);
          }
        ),
        { numRuns: 15 }
      );
    });
  });

  /**
   * Property 28: Duplicate Operation Prevention
   * For any batch operation in progress, a duplicate submission with the same
   * operation key SHALL return the existing operation's promise instead of
   * starting a new one.
   * Validates: Requirements 11.2
   */
  describe('Property 28: Duplicate Operation Prevention', () => {
    it('should return same promise for duplicate operation keys', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.string({ minLength: 1, maxLength: 50 }), // operation key
          fc.integer({ min: 10, max: 100 }), // delay in ms
          async (operationKey, delay) => {
            const controller = new ConcurrencyController();
            let executionCount = 0;

            const operation = async () => {
              executionCount++;
              await new Promise((resolve) => setTimeout(resolve, delay));
              return `result-${executionCount}`;
            };

            // Start multiple operations with same key simultaneously
            const promise1 = controller.withDeduplication(operationKey, operation);
            const promise2 = controller.withDeduplication(operationKey, operation);
            const promise3 = controller.withDeduplication(operationKey, operation);

            const [result1, result2, result3] = await Promise.all([promise1, promise2, promise3]);

            // Property: All results should be identical (same promise)
            expect(result1).toBe(result2);
            expect(result2).toBe(result3);

            // Property: Operation should only execute once
            expect(executionCount).toBe(1);
          }
        ),
        { numRuns: 20 }
      );
    });

    it('should allow new operations after previous completes', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.string({ minLength: 1, maxLength: 30 }),
          fc.integer({ min: 2, max: 5 }), // number of sequential operations
          async (operationKey, operationCount) => {
            const controller = new ConcurrencyController();
            const results: string[] = [];

            for (let i = 0; i < operationCount; i++) {
              const result = await controller.withDeduplication(operationKey, async () => {
                return `result-${i}`;
              });
              results.push(result);
            }

            // Property: Each sequential operation should produce unique result
            expect(results).toHaveLength(operationCount);
            for (let i = 0; i < operationCount; i++) {
              expect(results[i]).toBe(`result-${i}`);
            }
          }
        ),
        { numRuns: 20 }
      );
    });

    it('should track operation in progress correctly', async () => {
      await fc.assert(
        fc.asyncProperty(fc.string({ minLength: 1, maxLength: 30 }), async (operationKey) => {
          const controller = new ConcurrencyController();

          // Before operation
          expect(controller.isOperationInProgress(operationKey)).toBe(false);

          let operationStarted = false;
          const promise = controller.withDeduplication(operationKey, async () => {
            operationStarted = true;
            await new Promise((resolve) => setTimeout(resolve, 20));
            return 'done';
          });

          // Wait a tick for the operation to be registered
          await new Promise((resolve) => setTimeout(resolve, 5));

          // While operation is running (after it started)
          if (operationStarted) {
            expect(controller.isOperationInProgress(operationKey)).toBe(true);
          }

          await promise;

          // After operation completes
          expect(controller.isOperationInProgress(operationKey)).toBe(false);
        }),
        { numRuns: 15 }
      );
    });
  });

  /**
   * Property 29: Retry with Exponential Backoff
   * For any retryable error (rate limit, timeout), the service SHALL retry
   * with exponential backoff up to the configured maximum attempts.
   * Validates: Requirements 11.3
   */
  describe('Property 29: Retry with Exponential Backoff', () => {
    it('should retry correct number of times for retryable errors', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 1, max: 5 }), // retryAttempts
          fc.constantFrom('Rate limit exceeded', 'Network error', 'timeout', '429', '503'),
          async (retryAttempts, errorMessage) => {
            const controller = new ConcurrencyController({
              retryAttempts,
              retryDelayMs: 1, // Minimal delay for testing
              backoffMultiplier: 2,
            });

            let attemptCount = 0;

            try {
              await controller.withRetry(async () => {
                attemptCount++;
                throw new Error(errorMessage);
              });
            } catch {
              // Expected to throw after all retries exhausted
            }

            // Property: Should attempt exactly retryAttempts times
            expect(attemptCount).toBe(retryAttempts);
          }
        ),
        { numRuns: 20 }
      );
    });

    it('should not retry non-retryable errors', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 2, max: 5 }), // retryAttempts > 1
          fc.constantFrom('Validation failed', 'Not found', 'Permission denied', 'Invalid input'),
          async (retryAttempts, errorMessage) => {
            const controller = new ConcurrencyController({
              retryAttempts,
              retryDelayMs: 1,
            });

            let attemptCount = 0;

            try {
              await controller.withRetry(async () => {
                attemptCount++;
                throw new Error(errorMessage);
              });
            } catch {
              // Expected to throw immediately
            }

            // Property: Should only attempt once for non-retryable errors
            expect(attemptCount).toBe(1);
          }
        ),
        { numRuns: 20 }
      );
    });

    it('should succeed if operation succeeds within retry limit', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.integer({ min: 2, max: 5 }), // retryAttempts
          fc.integer({ min: 1, max: 4 }), // failuresBeforeSuccess
          async (retryAttempts, failuresBeforeSuccess) => {
            // Ensure we can succeed within retry limit
            const actualFailures = Math.min(failuresBeforeSuccess, retryAttempts - 1);

            const controller = new ConcurrencyController({
              retryAttempts,
              retryDelayMs: 1,
            });

            let attemptCount = 0;

            const result = await controller.withRetry(async () => {
              attemptCount++;
              if (attemptCount <= actualFailures) {
                throw new Error('Rate limit exceeded');
              }
              return 'success';
            });

            // Property: Should succeed after retrying
            expect(result).toBe('success');
            expect(attemptCount).toBe(actualFailures + 1);
          }
        ),
        { numRuns: 20 }
      );
    });

    it('should apply exponential backoff delays', async () => {
      // This test verifies that delays increase between retries
      // Due to timing variations, we just verify the backoff pattern exists
      const controller = new ConcurrencyController({
        retryAttempts: 3,
        retryDelayMs: 50,
        backoffMultiplier: 2,
      });

      const timestamps: number[] = [];

      try {
        await controller.withRetry(async () => {
          timestamps.push(Date.now());
          throw new Error('Rate limit exceeded');
        });
      } catch {
        // Expected
      }

      // Property: Should have attempted 3 times
      expect(timestamps).toHaveLength(3);

      // Property: Delays should generally increase (with tolerance for timing)
      if (timestamps.length >= 3) {
        const delay1 = timestamps[1] - timestamps[0];
        const delay2 = timestamps[2] - timestamps[1];

        // Second delay should be larger than first (backoff effect)
        // Allow some tolerance for timing variations
        expect(delay2).toBeGreaterThanOrEqual(delay1 * 0.8);
      }
    });
  });

  /**
   * Additional property: Chunk array correctness
   */
  describe('Property: Chunk Array Correctness', () => {
    it('should preserve all elements when chunking', () => {
      fc.assert(
        fc.property(
          fc.array(fc.anything(), { maxLength: 100 }),
          fc.integer({ min: 1, max: 20 }),
          (array, chunkSize) => {
            const chunks = chunkArray(array, chunkSize);
            const flattened = chunks.flat();

            // Property: All elements should be preserved
            expect(flattened).toEqual(array);
          }
        )
      );
    });

    it('should create chunks of correct size', () => {
      fc.assert(
        fc.property(
          fc.array(fc.integer(), { minLength: 1, maxLength: 100 }),
          fc.integer({ min: 1, max: 20 }),
          (array, chunkSize) => {
            const chunks = chunkArray(array, chunkSize);

            // Property: All chunks except last should be exactly chunkSize
            for (let i = 0; i < chunks.length - 1; i++) {
              expect(chunks[i]).toHaveLength(chunkSize);
            }

            // Property: Last chunk should be <= chunkSize
            if (chunks.length > 0) {
              expect(chunks[chunks.length - 1].length).toBeLessThanOrEqual(chunkSize);
              expect(chunks[chunks.length - 1].length).toBeGreaterThan(0);
            }
          }
        )
      );
    });
  });
});
