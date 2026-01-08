/**
 * Performance Monitor Module
 *
 * Provides performance monitoring infrastructure for service layer operations with:
 * - Operation timing with startTimer/stopTimer pattern
 * - Slow operation detection (>3s threshold by default)
 * - Metrics storage for debugging and analysis
 *
 * @module services/utils/performance
 */

import type { ServiceLogger } from './logger';

/**
 * Performance metrics for a single operation
 */
export interface PerformanceMetrics {
  /** Name of the method/operation */
  methodName: string;
  /** Duration in milliseconds */
  duration: number;
  /** Whether the operation succeeded */
  success: boolean;
  /** Timestamp when the metric was recorded */
  timestamp: Date;
  /** Optional additional context */
  context?: Record<string, unknown>;
}

/**
 * Performance monitor configuration
 */
export interface PerformanceMonitorConfig {
  /** Threshold in milliseconds for slow operation warnings (default: 3000ms) */
  slowThresholdMs?: number;
  /** Maximum number of metrics to store (default: 1000) */
  maxMetrics?: number;
  /** Logger for slow operation warnings */
  logger?: ServiceLogger;
  /** Custom handler for metrics (for testing or external reporting) */
  metricsHandler?: (metric: PerformanceMetrics) => void;
}

/**
 * Performance Monitor Interface
 *
 * Provides methods for tracking operation performance and detecting slow operations.
 */
export interface PerformanceMonitor {
  /**
   * Start a timer for an operation
   * @param operationName - Name of the operation being timed
   * @returns A function that stops the timer and returns the duration in milliseconds
   */
  startTimer(operationName: string): () => number;

  /**
   * Record a performance metric
   * @param metric - The metric to record
   */
  recordMetric(metric: PerformanceMetrics): void;

  /**
   * Get all recorded metrics
   * @returns Array of performance metrics
   */
  getMetrics(): PerformanceMetrics[];

  /**
   * Get metrics for a specific method
   * @param methodName - Name of the method to filter by
   * @returns Array of metrics for the specified method
   */
  getMetricsByMethod(methodName: string): PerformanceMetrics[];

  /**
   * Get slow operations (operations exceeding the threshold)
   * @returns Array of metrics for slow operations
   */
  getSlowOperations(): PerformanceMetrics[];

  /**
   * Clear all stored metrics
   */
  clearMetrics(): void;

  /**
   * Get the slow operation threshold in milliseconds
   */
  getSlowThreshold(): number;

  /**
   * Get average duration for a specific method
   * @param methodName - Name of the method
   * @returns Average duration in milliseconds, or undefined if no metrics exist
   */
  getAverageDuration(methodName: string): number | undefined;
}

/**
 * Default Performance Monitor Implementation
 *
 * Provides:
 * - Operation timing with automatic slow operation detection
 * - Metrics storage with configurable limits
 * - Integration with ServiceLogger for warnings
 */
export class DefaultPerformanceMonitor implements PerformanceMonitor {
  private metrics: PerformanceMetrics[] = [];
  private readonly slowThresholdMs: number;
  private readonly maxMetrics: number;
  private readonly logger?: ServiceLogger;
  private readonly metricsHandler?: (metric: PerformanceMetrics) => void;

  constructor(config: PerformanceMonitorConfig = {}) {
    this.slowThresholdMs = config.slowThresholdMs ?? 3000;
    this.maxMetrics = config.maxMetrics ?? 1000;
    this.logger = config.logger;
    this.metricsHandler = config.metricsHandler;
  }

  startTimer(operationName: string): () => number {
    const startTime = performance.now();

    return () => {
      const duration = performance.now() - startTime;

      // Check for slow operation and emit warning
      if (duration > this.slowThresholdMs) {
        if (this.logger) {
          this.logger.warn(`Slow operation detected: ${operationName}`, {
            duration: Math.round(duration),
            threshold: this.slowThresholdMs,
            operationName,
          });
        }
      }

      return duration;
    };
  }

  recordMetric(metric: PerformanceMetrics): void {
    this.metrics.push(metric);

    // Keep only last maxMetrics metrics
    if (this.metrics.length > this.maxMetrics) {
      this.metrics = this.metrics.slice(-this.maxMetrics);
    }

    // Call custom handler if provided
    if (this.metricsHandler) {
      this.metricsHandler(metric);
    }
  }

  getMetrics(): PerformanceMetrics[] {
    return [...this.metrics];
  }

  getMetricsByMethod(methodName: string): PerformanceMetrics[] {
    return this.metrics.filter((m) => m.methodName === methodName);
  }

  getSlowOperations(): PerformanceMetrics[] {
    return this.metrics.filter((m) => m.duration > this.slowThresholdMs);
  }

  clearMetrics(): void {
    this.metrics = [];
  }

  getSlowThreshold(): number {
    return this.slowThresholdMs;
  }

  getAverageDuration(methodName: string): number | undefined {
    const methodMetrics = this.getMetricsByMethod(methodName);
    if (methodMetrics.length === 0) return undefined;

    const total = methodMetrics.reduce((sum, m) => sum + m.duration, 0);
    return total / methodMetrics.length;
  }
}

/**
 * Create a performance monitor instance
 *
 * @param config - Optional configuration
 * @returns PerformanceMonitor instance
 */
export function createPerformanceMonitor(
  config?: PerformanceMonitorConfig
): PerformanceMonitor {
  return new DefaultPerformanceMonitor(config);
}

/**
 * Utility function to measure async operation duration
 *
 * @param operation - The async operation to measure
 * @param monitor - Performance monitor to record the metric
 * @param methodName - Name of the method for the metric
 * @returns The result of the operation
 */
export async function measureAsync<T>(
  operation: () => Promise<T>,
  monitor: PerformanceMonitor,
  methodName: string
): Promise<T> {
  const stopTimer = monitor.startTimer(methodName);
  let success = true;

  try {
    return await operation();
  } catch (error) {
    success = false;
    throw error;
  } finally {
    const duration = stopTimer();
    monitor.recordMetric({
      methodName,
      duration,
      success,
      timestamp: new Date(),
    });
  }
}
