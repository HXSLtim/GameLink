/**
 * Base Service Class
 * Provides common functionality for all domain services
 *
 * @module services/domain/base
 */

import { adminApi } from '@/api/admin';
import {
  ServiceException,
  ServiceErrorCodes,
  type ServiceResult,
  type BatchResult,
  type BatchItemResult,
  ServiceResultHelper,
} from '../utils';

/**
 * Logger interface for service operations
 */
export interface ServiceLogger {
  debug(message: string, context?: Record<string, unknown>): void;
  info(message: string, context?: Record<string, unknown>): void;
  warn(message: string, context?: Record<string, unknown>): void;
  error(message: string, error?: Error, context?: Record<string, unknown>): void;
}

/**
 * Performance metrics for service operations
 */
export interface PerformanceMetrics {
  methodName: string;
  duration: number;
  success: boolean;
  timestamp: Date;
}

/**
 * Performance monitor interface
 */
export interface PerformanceMonitor {
  startTimer(operationName: string): () => number;
  recordMetric(metric: PerformanceMetrics): void;
  getMetrics(): PerformanceMetrics[];
}

/**
 * Dependencies that can be injected into services
 */
export interface ServiceDependencies {
  /** API client for making HTTP requests */
  api?: typeof adminApi;
  /** Optional logger for service operations */
  logger?: ServiceLogger;
  /** Optional performance monitor */
  perfMonitor?: PerformanceMonitor;
}

/**
 * Default console-based logger implementation
 */
export class DefaultServiceLogger implements ServiceLogger {
  private serviceName: string;
  private isDev: boolean;

  constructor(serviceName: string) {
    this.serviceName = serviceName;
    this.isDev = import.meta.env?.DEV ?? false;
  }

  debug(message: string, context?: Record<string, unknown>): void {
    if (this.isDev) {
      console.debug(`[${this.serviceName}] ${message}`, context ?? '');
    }
  }

  info(message: string, context?: Record<string, unknown>): void {
    console.info(`[${this.serviceName}] ${message}`, context ?? '');
  }

  warn(message: string, context?: Record<string, unknown>): void {
    console.warn(`[${this.serviceName}] ${message}`, context ?? '');
  }

  error(message: string, error?: Error, context?: Record<string, unknown>): void {
    console.error(`[${this.serviceName}] ${message}`, { error, ...context });
  }
}

/**
 * Default performance monitor implementation
 */
export class DefaultPerformanceMonitor implements PerformanceMonitor {
  private metrics: PerformanceMetrics[] = [];
  private readonly SLOW_THRESHOLD_MS = 3000;
  private readonly MAX_METRICS = 1000;
  private logger?: ServiceLogger;

  constructor(logger?: ServiceLogger) {
    this.logger = logger;
  }

  startTimer(operationName: string): () => number {
    const startTime = performance.now();
    return () => {
      const duration = performance.now() - startTime;
      if (duration > this.SLOW_THRESHOLD_MS && this.logger) {
        this.logger.warn(`Slow operation detected: ${operationName}`, { duration });
      }
      return duration;
    };
  }

  recordMetric(metric: PerformanceMetrics): void {
    this.metrics.push(metric);
    // Keep only last MAX_METRICS metrics
    if (this.metrics.length > this.MAX_METRICS) {
      this.metrics = this.metrics.slice(-this.MAX_METRICS);
    }
  }

  getMetrics(): PerformanceMetrics[] {
    return [...this.metrics];
  }
}

/**
 * Sensitive keys that should be redacted in logs
 */
const SENSITIVE_KEYS = ['password', 'token', 'secret', 'key', 'authorization'];

/**
 * Abstract base class for all domain services
 *
 * Provides:
 * - Dependency injection support
 * - Error handling and wrapping
 * - Async operation helpers
 * - Logging with parameter sanitization
 * - Performance monitoring
 *
 * @example
 * ```typescript
 * class UserService extends BaseService {
 *   async getUser(id: number): Promise<ServiceResult<User>> {
 *     return this.wrapAsync(
 *       async () => {
 *         const response = await this.api.getUser(id);
 *         return response.data.data;
 *       },
 *       'getUser'
 *     );
 *   }
 * }
 * ```
 */
export abstract class BaseService {
  protected api: typeof adminApi;
  protected logger: ServiceLogger;
  protected perfMonitor: PerformanceMonitor;

  constructor(deps: ServiceDependencies = {}) {
    this.api = deps.api ?? adminApi;
    this.logger = deps.logger ?? new DefaultServiceLogger(this.constructor.name);
    this.perfMonitor = deps.perfMonitor ?? new DefaultPerformanceMonitor(this.logger);
  }

  /**
   * Handle and wrap errors into ServiceException
   *
   * @param error - The error to handle
   * @param context - Context string for the error message
   * @returns ServiceException with proper error code and message
   */
  protected handleError(error: unknown, context: string): ServiceException {
    if (error instanceof ServiceException) {
      return error;
    }

    // Handle Axios errors
    if (this.isAxiosError(error)) {
      const status = error.response?.status;
      const message = error.response?.data?.message || error.message;

      switch (status) {
        case 400:
          return new ServiceException(
            ServiceErrorCodes.VALIDATION_ERROR,
            `${context}: ${message}`,
            { status },
            error
          );
        case 401:
          return new ServiceException(
            ServiceErrorCodes.UNAUTHORIZED,
            `${context}: ${message}`,
            { status },
            error
          );
        case 403:
          return new ServiceException(
            ServiceErrorCodes.FORBIDDEN,
            `${context}: ${message}`,
            { status },
            error
          );
        case 404:
          return new ServiceException(
            ServiceErrorCodes.NOT_FOUND,
            `${context}: ${message}`,
            { status },
            error
          );
        default:
          return new ServiceException(
            ServiceErrorCodes.UNKNOWN_ERROR,
            `${context}: ${message}`,
            { status },
            error
          );
      }
    }

    // Handle network errors
    if (error instanceof Error && error.message.includes('Network')) {
      return new ServiceException(
        ServiceErrorCodes.NETWORK_ERROR,
        `${context}: Network error`,
        undefined,
        error
      );
    }

    // Handle timeout errors
    if (error instanceof Error && error.message.includes('timeout')) {
      return new ServiceException(
        ServiceErrorCodes.TIMEOUT_ERROR,
        `${context}: Request timeout`,
        undefined,
        error
      );
    }

    // Generic error handling
    const message = error instanceof Error ? error.message : String(error);
    return new ServiceException(
      ServiceErrorCodes.UNKNOWN_ERROR,
      `${context}: ${message}`,
      undefined,
      error instanceof Error ? error : undefined
    );
  }

  /**
   * Wrap an async operation with error handling and return ServiceResult
   *
   * @param operation - The async operation to execute
   * @param context - Context string for error messages
   * @returns ServiceResult with data or error
   */
  protected async wrapAsync<T>(
    operation: () => Promise<T>,
    context: string
  ): Promise<ServiceResult<T>> {
    try {
      const data = await operation();
      return ServiceResultHelper.success(data);
    } catch (error) {
      const serviceError = this.handleError(error, context);
      this.logger.error(`${context} failed`, serviceError, {
        code: serviceError.code,
      });
      return ServiceResultHelper.failure(serviceError.toError());
    }
  }

  /**
   * Wrap an async operation with logging and performance monitoring
   *
   * @param methodName - Name of the method being executed
   * @param params - Parameters to log (will be sanitized)
   * @param operation - The async operation to execute
   * @returns The result of the operation
   */
  protected async withLogging<T>(
    methodName: string,
    params: Record<string, unknown>,
    operation: () => Promise<T>
  ): Promise<T> {
    const sanitizedParams = this.sanitizeParams(params);
    this.logger.debug(`${methodName} started`, { params: sanitizedParams });

    const stopTimer = this.perfMonitor.startTimer(methodName);

    try {
      const result = await operation();
      const duration = stopTimer();

      this.logger.info(`${methodName} completed`, { duration, success: true });
      this.perfMonitor.recordMetric({
        methodName,
        duration,
        success: true,
        timestamp: new Date(),
      });

      return result;
    } catch (error) {
      const duration = stopTimer();

      this.logger.error(`${methodName} failed`, error as Error, {
        params: sanitizedParams,
        duration,
      });
      this.perfMonitor.recordMetric({
        methodName,
        duration,
        success: false,
        timestamp: new Date(),
      });

      throw error;
    }
  }

  /**
   * Execute a batch operation with progress tracking
   *
   * @param items - Items to process
   * @param processor - Function to process each item
   * @param context - Context string for error messages
   * @param onProgress - Optional progress callback
   * @returns BatchResult with individual results
   */
  protected async executeBatch<T, R>(
    items: T[],
    processor: (item: T, index: number) => Promise<R>,
    context: string,
    onProgress?: (completed: number, total: number) => void
  ): Promise<BatchResult<R>> {
    const results: BatchItemResult<R>[] = [];
    const total = items.length;
    let lastProgressLog = 0;

    for (let index = 0; index < items.length; index++) {
      const item = items[index];
      try {
        const data = await processor(item, index);
        results.push({ index, success: true, data });
      } catch (error) {
        const serviceError = this.handleError(error, `${context} item ${index}`);
        results.push({ index, success: false, error: serviceError.toError() });
      }

      // Report progress
      const completed = index + 1;
      if (onProgress) {
        onProgress(completed, total);
      }

      // Log progress at 10% intervals
      const progressPercent = Math.floor((completed / total) * 100);
      if (progressPercent >= lastProgressLog + 10) {
        this.logger.info(`${context} progress`, {
          completed,
          total,
          percent: progressPercent,
        });
        lastProgressLog = progressPercent;
      }
    }

    return ServiceResultHelper.fromResults(results);
  }

  /**
   * Sanitize parameters by redacting sensitive values
   *
   * @param params - Parameters to sanitize
   * @returns Sanitized parameters
   */
  protected sanitizeParams(params: Record<string, unknown>): Record<string, unknown> {
    const sanitized: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(params)) {
      if (SENSITIVE_KEYS.some((sk) => key.toLowerCase().includes(sk))) {
        sanitized[key] = '[REDACTED]';
      } else if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
        sanitized[key] = this.sanitizeParams(value as Record<string, unknown>);
      } else {
        sanitized[key] = value;
      }
    }

    return sanitized;
  }

  /**
   * Check if an error is an Axios error
   */
  private isAxiosError(
    error: unknown
  ): error is Error & { response?: { status: number; data?: { message?: string } } } {
    return (
      error instanceof Error &&
      'response' in error &&
      typeof (error as { response?: unknown }).response === 'object'
    );
  }
}
