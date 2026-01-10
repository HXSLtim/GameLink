/**
 * Service Layer Result Utilities
 * Provides standardized result types for service operations
 *
 * @module services/utils/serviceResult
 */

import type { ServiceError } from './serviceError';

/**
 * Result of a single service operation
 *
 * @template T - The type of data returned on success
 *
 * @example
 * ```typescript
 * // Success case
 * const result: ServiceResult<User> = {
 *   success: true,
 *   data: { id: 1, name: 'John' }
 * };
 *
 * // Error case
 * const result: ServiceResult<User> = {
 *   success: false,
 *   error: { code: 'NOT_FOUND', message: 'User not found' }
 * };
 * ```
 */
export interface ServiceResult<T> {
  /** Whether the operation succeeded */
  success: boolean;
  /** The result data (present when success is true) */
  data?: T;
  /** Error information (present when success is false) */
  error?: ServiceError;
}

/**
 * Individual item result in a batch operation
 *
 * @template T - The type of data returned for each item
 */
export interface BatchItemResult<T> {
  /** Index of the item in the original batch */
  index: number;
  /** Whether this item's operation succeeded */
  success: boolean;
  /** The result data for this item (present when success is true) */
  data?: T;
  /** Error information for this item (present when success is false) */
  error?: ServiceError;
}

/**
 * Result of a batch operation
 *
 * @template T - The type of data returned for each successful item
 *
 * @example
 * ```typescript
 * const result: BatchResult<void> = {
 *   success: false, // false if any item failed
 *   total: 10,
 *   succeeded: 8,
 *   failed: 2,
 *   results: [
 *     { index: 0, success: true },
 *     { index: 1, success: false, error: { code: 'NOT_FOUND', message: '...' } },
 *     // ...
 *   ]
 * };
 * ```
 */
export interface BatchResult<T> {
  /** Whether all operations succeeded */
  success: boolean;
  /** Total number of items in the batch */
  total: number;
  /** Number of successful operations */
  succeeded: number;
  /** Number of failed operations */
  failed: number;
  /** Individual results for each item */
  results: BatchItemResult<T>[];
}

/**
 * Helper functions for creating service results
 */
export const ServiceResultHelper = {
  /**
   * Create a successful result
   */
  success<T>(data: T): ServiceResult<T> {
    return { success: true, data };
  },

  /**
   * Create a failed result
   */
  failure<T>(error: ServiceError): ServiceResult<T> {
    return { success: false, error };
  },

  /**
   * Create a failed result from code and message
   */
  error<T>(code: string, message: string, details?: Record<string, unknown>): ServiceResult<T> {
    return {
      success: false,
      error: { code, message, details },
    };
  },

  /**
   * Create an empty batch result
   */
  emptyBatch<T>(total: number): BatchResult<T> {
    return {
      success: true,
      total,
      succeeded: 0,
      failed: 0,
      results: [],
    };
  },

  /**
   * Create a batch result from individual results
   */
  fromResults<T>(results: BatchItemResult<T>[]): BatchResult<T> {
    const succeeded = results.filter((r) => r.success).length;
    const failed = results.length - succeeded;

    return {
      success: failed === 0,
      total: results.length,
      succeeded,
      failed,
      results,
    };
  },
};
