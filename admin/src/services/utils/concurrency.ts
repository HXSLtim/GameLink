/**
 * Concurrency Control Utilities
 * Provides concurrency limiting, deduplication, retry with backoff, and batch chunking
 * @module services/utils/concurrency
 */

/**
 * Configuration for concurrency control
 */
export interface ConcurrencyConfig {
  /** Maximum concurrent operations (default: 5) */
  maxConcurrent: number;
  /** Chunk size for batch processing (default: 50) */
  chunkSize: number;
  /** Maximum retry attempts (default: 3) */
  retryAttempts: number;
  /** Initial retry delay in milliseconds (default: 1000) */
  retryDelayMs: number;
  /** Backoff multiplier for retry delays (default: 2) */
  backoffMultiplier: number;
}

/**
 * Default concurrency configuration
 */
export const DEFAULT_CONCURRENCY_CONFIG: ConcurrencyConfig = {
  maxConcurrent: 5,
  chunkSize: 50,
  retryAttempts: 3,
  retryDelayMs: 1000,
  backoffMultiplier: 2,
};

/**
 * Result of a processed item in batch operations
 */
export interface ProcessResult<R> {
  index: number;
  success: boolean;
  result?: R;
  error?: Error;
}

/**
 * Simple semaphore implementation for concurrency limiting
 */
export class Semaphore {
  private permits: number;
  private waiting: Array<() => void> = [];

  constructor(permits: number) {
    if (permits < 1) {
      throw new Error('Semaphore permits must be at least 1');
    }
    this.permits = permits;
  }

  /**
   * Acquire a permit, waiting if none available
   */
  async acquire(): Promise<void> {
    if (this.permits > 0) {
      this.permits--;
      return;
    }

    return new Promise((resolve) => {
      this.waiting.push(resolve);
    });
  }

  /**
   * Release a permit
   */
  release(): void {
    const next = this.waiting.shift();
    if (next) {
      next();
    } else {
      this.permits++;
    }
  }

  /**
   * Get current available permits
   */
  getAvailablePermits(): number {
    return this.permits;
  }

  /**
   * Get number of waiting operations
   */
  getWaitingCount(): number {
    return this.waiting.length;
  }
}

/**
 * Check if an error is retryable (rate limit, timeout, network errors)
 */
export function isRetryableError(error: unknown): boolean {
  if (error instanceof Error) {
    const message = error.message.toLowerCase();
    return (
      message.includes('rate limit') ||
      message.includes('timeout') ||
      message.includes('network') ||
      message.includes('econnreset') ||
      message.includes('econnrefused') ||
      message.includes('429') ||
      message.includes('503')
    );
  }
  return false;
}

/**
 * Split an array into chunks of specified size
 */
export function chunkArray<T>(array: T[], size: number): T[][] {
  if (size < 1) {
    throw new Error('Chunk size must be at least 1');
  }
  const chunks: T[][] = [];
  for (let i = 0; i < array.length; i += size) {
    chunks.push(array.slice(i, i + size));
  }
  return chunks;
}

/**
 * Sleep for specified milliseconds
 */
export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * ConcurrencyController provides utilities for managing concurrent operations
 */
export class ConcurrencyController {
  private activeOperations = new Map<string, Promise<unknown>>();
  private config: ConcurrencyConfig;

  constructor(config: Partial<ConcurrencyConfig> = {}) {
    this.config = { ...DEFAULT_CONCURRENCY_CONFIG, ...config };
  }

  /**
   * Get current configuration
   */
  getConfig(): ConcurrencyConfig {
    return { ...this.config };
  }

  /**
   * Get number of active operations
   */
  getActiveOperationCount(): number {
    return this.activeOperations.size;
  }

  /**
   * Check if an operation with the given key is in progress
   */
  isOperationInProgress(operationKey: string): boolean {
    return this.activeOperations.has(operationKey);
  }

  /**
   * Prevent duplicate submissions of the same operation
   * If an operation with the same key is in progress, returns its promise
   */
  async withDeduplication<T>(
    operationKey: string,
    operation: () => Promise<T>
  ): Promise<T> {
    const existing = this.activeOperations.get(operationKey);
    if (existing) {
      return existing as Promise<T>;
    }

    const promise = operation().finally(() => {
      this.activeOperations.delete(operationKey);
    });

    this.activeOperations.set(operationKey, promise);
    return promise;
  }

  /**
   * Execute an operation with retry and exponential backoff
   */
  async withRetry<T>(
    operation: () => Promise<T>,
    config?: Partial<Pick<ConcurrencyConfig, 'retryAttempts' | 'retryDelayMs' | 'backoffMultiplier'>>
  ): Promise<T> {
    const retryAttempts = config?.retryAttempts ?? this.config.retryAttempts;
    const retryDelayMs = config?.retryDelayMs ?? this.config.retryDelayMs;
    const backoffMultiplier = config?.backoffMultiplier ?? this.config.backoffMultiplier;

    let lastError: Error | undefined;
    let delay = retryDelayMs;

    for (let attempt = 0; attempt < retryAttempts; attempt++) {
      try {
        return await operation();
      } catch (error) {
        lastError = error instanceof Error ? error : new Error(String(error));

        // Check if error is retryable
        if (!isRetryableError(error)) {
          throw lastError;
        }

        // If not the last attempt, wait before retrying
        if (attempt < retryAttempts - 1) {
          await sleep(delay);
          delay *= backoffMultiplier;
        }
      }
    }

    throw lastError;
  }

  /**
   * Process items with concurrency limit
   */
  async processWithConcurrency<T, R>(
    items: T[],
    processor: (item: T, index: number) => Promise<R>,
    onProgress?: (completed: number, total: number) => void
  ): Promise<ProcessResult<R>[]> {
    if (items.length === 0) {
      return [];
    }

    const results: ProcessResult<R>[] = [];
    const chunks = chunkArray(items, this.config.chunkSize);
    let completed = 0;
    let globalIndex = 0;

    for (const chunk of chunks) {
      const chunkStartIndex = globalIndex;
      const chunkResults = await this.processChunk(chunk, processor, chunkStartIndex);
      results.push(...chunkResults);

      completed += chunk.length;
      globalIndex += chunk.length;

      if (onProgress) {
        onProgress(completed, items.length);
      }
    }

    return results;
  }

  /**
   * Process a single chunk with concurrency limiting
   */
  private async processChunk<T, R>(
    chunk: T[],
    processor: (item: T, index: number) => Promise<R>,
    startIndex: number
  ): Promise<ProcessResult<R>[]> {
    const semaphore = new Semaphore(this.config.maxConcurrent);

    return Promise.all(
      chunk.map(async (item, i) => {
        const index = startIndex + i;
        await semaphore.acquire();

        try {
          const result = await this.withRetry(() => processor(item, index));
          return { index, success: true, result };
        } catch (error) {
          return {
            index,
            success: false,
            error: error instanceof Error ? error : new Error(String(error)),
          };
        } finally {
          semaphore.release();
        }
      })
    );
  }
}

/**
 * Create a ConcurrencyController with optional configuration
 */
export function createConcurrencyController(
  config?: Partial<ConcurrencyConfig>
): ConcurrencyController {
  return new ConcurrencyController(config);
}
