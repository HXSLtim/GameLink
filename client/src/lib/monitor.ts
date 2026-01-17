/**
 * Performance Monitoring Module
 *
 * Tracks API performance metrics including:
 * - Request count and error rates
 * - Response times
 * - Token refresh success rates
 * - Queue management metrics
 *
 * Usage:
 * ```ts
 * import { performanceMonitor } from '@/lib/monitor';
 *
 * // Record a request
 * performanceMonitor.recordRequest(duration, success);
 *
 * // Get metrics report
 * const report = performanceMonitor.getMetrics();
 * ```
 */

interface PerformanceMetrics {
    requestCount: number;
    errorCount: number;
    totalResponseTime: number;
    tokenRefreshCount: number;
    tokenRefreshSuccessCount: number;
    queueRejectionCount: number;
    queueTimeoutCount: number;
}

interface DetailedMetrics extends PerformanceMetrics {
    averageResponseTime: number;
    errorRate: number;
    tokenRefreshSuccessRate: number;
}

interface PerformanceReport {
    timestamp: string;
    metrics: DetailedMetrics;
}

/**
 * Performance Monitor Class
 *
 * Thread-safe singleton for tracking HTTP client performance
 */
class PerformanceMonitor {
    private metrics: PerformanceMetrics = {
        requestCount: 0,
        errorCount: 0,
        totalResponseTime: 0,
        tokenRefreshCount: 0,
        tokenRefreshSuccessCount: 0,
        queueRejectionCount: 0,
        queueTimeoutCount: 0,
    };

    /**
     * Record an HTTP request
     *
     * @param duration - Request duration in milliseconds
     * @param success - Whether the request succeeded
     */
    recordRequest(duration: number, success: boolean): void {
        this.metrics.requestCount++;
        this.metrics.totalResponseTime += duration;

        if (!success) {
            this.metrics.errorCount++;
        }

        // Log warnings for slow requests in development
        if (import.meta.env.DEV && duration > 3000) {
            console.warn(
                `[Performance] Slow request detected: ${duration}ms`,
                this.getQuickStats()
            );
        }
    }

    /**
     * Record a token refresh attempt
     *
     * @param success - Whether the refresh succeeded
     */
    recordTokenRefresh(success: boolean): void {
        this.metrics.tokenRefreshCount++;

        if (success) {
            this.metrics.tokenRefreshSuccessCount++;
        } else {
            // Log failed refresh attempts
            console.error('[Performance] Token refresh failed', this.getQuickStats());
        }
    }

    /**
     * Record a queue rejection (queue full)
     */
    recordQueueRejection(): void {
        this.metrics.queueRejectionCount++;

        if (import.meta.env.DEV) {
            console.warn('[Performance] Request queue full, rejecting new request');
        }
    }

    /**
     * Record a queue timeout (request waited too long)
     */
    recordQueueTimeout(): void {
        this.metrics.queueTimeoutCount++;

        if (import.meta.env.DEV) {
            console.warn('[Performance] Request queue timeout');
        }
    }

    /**
     * Get current metrics with computed values
     */
    getMetrics(): DetailedMetrics {
        const { requestCount, errorCount, totalResponseTime, tokenRefreshCount, tokenRefreshSuccessCount } =
            this.metrics;

        return {
            ...this.metrics,
            averageResponseTime: requestCount > 0 ? totalResponseTime / requestCount : 0,
            errorRate: requestCount > 0 ? (errorCount / requestCount) * 100 : 0,
            tokenRefreshSuccessRate:
                tokenRefreshCount > 0 ? (tokenRefreshSuccessCount / tokenRefreshCount) * 100 : 0,
        };
    }

    /**
     * Get a quick stats summary for logging
     */
    private getQuickStats(): string {
        const m = this.metrics;
        return `Requests: ${m.requestCount}, Errors: ${m.errorCount}, Avg: ${
            this.getMetrics().averageResponseTime.toFixed(2)
        }ms`;
    }

    /**
     * Reset all metrics (useful for testing or periodic resets)
     */
    reset(): void {
        this.metrics = {
            requestCount: 0,
            errorCount: 0,
            totalResponseTime: 0,
            tokenRefreshCount: 0,
            tokenRefreshSuccessCount: 0,
            queueRejectionCount: 0,
            queueTimeoutCount: 0,
        };
    }

    /**
     * Generate a detailed performance report
     */
    generateReport(): PerformanceReport {
        return {
            timestamp: new Date().toISOString(),
            metrics: this.getMetrics(),
        };
    }
}

// Export singleton instance
export const performanceMonitor = new PerformanceMonitor();

/**
 * Get performance report as a printable string
 */
export function getPerformanceReportString(): string {
    const report = performanceMonitor.generateReport();
    const m = report.metrics;

    return `
Performance Report - ${report.timestamp}
==========================================
Requests:              ${m.requestCount}
Errors:                ${m.errorCount} (${m.errorRate.toFixed(2)}%)
Avg Response Time:     ${m.averageResponseTime.toFixed(2)}ms
Token Refreshes:       ${m.tokenRefreshCount}
Token Refresh Success: ${m.tokenRefreshSuccessRate.toFixed(2)}%
Queue Rejections:      ${m.queueRejectionCount}
Queue Timeouts:        ${m.queueTimeoutCount}
==========================================
`;
}

/**
 * Development-only periodic reporting
 */
if (import.meta.env.DEV) {
    // Log performance report every 60 seconds
    setInterval(() => {
        const report = performanceMonitor.generateReport();
        const m = report.metrics;

        // Only log if there's actual activity
        if (m.requestCount > 0 || m.tokenRefreshCount > 0) {
            console.log(getPerformanceReportString());
        }
    }, 60000); // 60 seconds

    // Log on page unload
    window.addEventListener('beforeunload', () => {
        const report = performanceMonitor.generateReport();
        const m = report.metrics;

        if (m.requestCount > 0 || m.tokenRefreshCount > 0) {
            console.log('[Performance] Session Summary:', getPerformanceReportString());
        }
    });
}
