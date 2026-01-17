/**
 * Performance Monitor Tests
 *
 * Tests for the performance monitoring module including:
 * - Request recording
 * - Token refresh tracking
 * - Queue management metrics
 * - Metrics calculation
 * - Report generation
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { performanceMonitor, getPerformanceReportString } from '../monitor';

describe('Performance Monitor - Request Tracking', () => {
    beforeEach(() => {
        // Reset metrics before each test
        performanceMonitor.reset();
    });

    it('should record successful request', () => {
        performanceMonitor.recordRequest(100, true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(1);
        expect(metrics.errorCount).toBe(0);
        expect(metrics.totalResponseTime).toBe(100);
    });

    it('should record failed request', () => {
        performanceMonitor.recordRequest(200, false);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(1);
        expect(metrics.errorCount).toBe(1);
        expect(metrics.totalResponseTime).toBe(200);
    });

    it('should calculate average response time', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordRequest(200, true);
        performanceMonitor.recordRequest(300, true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.averageResponseTime).toBe(200); // (100 + 200 + 300) / 3
    });

    it('should calculate error rate correctly', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordRequest(100, false);
        performanceMonitor.recordRequest(100, true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.errorRate).toBeCloseTo(33.33, 2); // 1/3, allow floating point precision
    });

    it('should handle zero requests gracefully', () => {
        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(0);
        expect(metrics.averageResponseTime).toBe(0);
        expect(metrics.errorRate).toBe(0);
    });
});

describe('Performance Monitor - Token Refresh Tracking', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should record successful token refresh', () => {
        performanceMonitor.recordTokenRefresh(true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshCount).toBe(1);
        expect(metrics.tokenRefreshSuccessCount).toBe(1);
        expect(metrics.tokenRefreshSuccessRate).toBe(100);
    });

    it('should record failed token refresh', () => {
        performanceMonitor.recordTokenRefresh(false);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshCount).toBe(1);
        expect(metrics.tokenRefreshSuccessCount).toBe(0);
        expect(metrics.tokenRefreshSuccessRate).toBe(0);
    });

    it('should calculate token refresh success rate', () => {
        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordTokenRefresh(false);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshSuccessRate).toBeCloseTo(66.67, 2); // 2/3, allow floating point precision
    });

    it('should handle zero refreshes gracefully', () => {
        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshCount).toBe(0);
        expect(metrics.tokenRefreshSuccessRate).toBe(0);
    });
});

describe('Performance Monitor - Queue Management', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should record queue rejection', () => {
        performanceMonitor.recordQueueRejection();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.queueRejectionCount).toBe(1);
    });

    it('should record multiple queue rejections', () => {
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordQueueRejection();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.queueRejectionCount).toBe(3);
    });

    it('should record queue timeout', () => {
        performanceMonitor.recordQueueTimeout();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.queueTimeoutCount).toBe(1);
    });

    it('should record multiple queue timeouts', () => {
        performanceMonitor.recordQueueTimeout();
        performanceMonitor.recordQueueTimeout();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.queueTimeoutCount).toBe(2);
    });
});

describe('Performance Monitor - Report Generation', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should generate report with all metrics', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordRequest(200, false);
        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordQueueRejection();

        const report = performanceMonitor.generateReport();

        expect(report).toHaveProperty('timestamp');
        expect(report).toHaveProperty('metrics');
        expect(report.metrics.requestCount).toBe(2);
        expect(report.metrics.errorCount).toBe(1);
        expect(report.metrics.tokenRefreshCount).toBe(1);
        expect(report.metrics.queueRejectionCount).toBe(1);
    });

    it('should include ISO timestamp in report', () => {
        const report = performanceMonitor.generateReport();
        const timestamp = new Date(report.timestamp);

        expect(timestamp.toISOString()).toBe(report.timestamp);
    });

    it('should generate readable report string', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordTokenRefresh(true);

        const reportString = getPerformanceReportString();

        expect(reportString).toContain('Performance Report');
        expect(reportString).toContain('Requests:');
        expect(reportString).toContain('1');
        expect(reportString).toContain('Avg Response Time:');
    });
});

describe('Performance Monitor - Reset', () => {
    it('should reset all metrics', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordRequest(200, false);
        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordQueueTimeout();

        performanceMonitor.reset();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(0);
        expect(metrics.errorCount).toBe(0);
        expect(metrics.totalResponseTime).toBe(0);
        expect(metrics.tokenRefreshCount).toBe(0);
        expect(metrics.tokenRefreshSuccessCount).toBe(0);
        expect(metrics.queueRejectionCount).toBe(0);
        expect(metrics.queueTimeoutCount).toBe(0);
    });
});

describe('Performance Monitor - Realistic Scenarios', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should handle high request volume', () => {
        // Simulate 100 requests
        for (let i = 0; i < 100; i++) {
            const duration = Math.random() * 1000; // Random duration 0-1000ms
            const success = Math.random() > 0.05; // 95% success rate
            performanceMonitor.recordRequest(duration, success);
        }

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(100);
        expect(metrics.errorRate).toBeGreaterThan(0);
        expect(metrics.errorRate).toBeLessThan(20); // Should be around 5%
    });

    it('should track token refresh health', () => {
        // Simulate 10 token refreshes with 90% success rate
        for (let i = 0; i < 10; i++) {
            const success = i < 9; // First 9 succeed, last one fails
            performanceMonitor.recordTokenRefresh(success);
        }

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshCount).toBe(10);
        expect(metrics.tokenRefreshSuccessRate).toBe(90);
    });

    it('should detect queue issues', () => {
        // Simulate queue problems
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordQueueTimeout();

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.queueRejectionCount).toBe(2);
        expect(metrics.queueTimeoutCount).toBe(1);

        const report = getPerformanceReportString();
        expect(report).toContain('Queue Rejections:');
        expect(report).toContain('Queue Timeouts:');
    });

    it('should provide comprehensive performance overview', () => {
        // Simulate realistic usage
        for (let i = 0; i < 50; i++) {
            const duration = 50 + Math.random() * 200; // 50-250ms
            performanceMonitor.recordRequest(duration, Math.random() > 0.1);
        }

        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordTokenRefresh(true);

        const metrics = performanceMonitor.getMetrics();
        const report = performanceMonitor.generateReport();

        expect(metrics.requestCount).toBe(50);
        expect(metrics.averageResponseTime).toBeGreaterThan(50);
        expect(metrics.averageResponseTime).toBeLessThan(250);
        expect(metrics.tokenRefreshSuccessRate).toBe(100);
        expect(report.timestamp).toBeDefined();
    });
});

describe('Performance Monitor - Edge Cases', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should handle very fast requests (<1ms)', () => {
        performanceMonitor.recordRequest(0.1, true);
        performanceMonitor.recordRequest(0.5, true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(2);
        expect(metrics.totalResponseTime).toBe(0.6);
    });

    it('should handle very slow requests (>10s)', () => {
        performanceMonitor.recordRequest(15000, true);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.averageResponseTime).toBe(15000);
    });

    it('should handle 100% error rate', () => {
        for (let i = 0; i < 10; i++) {
            performanceMonitor.recordRequest(100, false);
        }

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.errorRate).toBe(100);
    });

    it('should handle 0% error rate', () => {
        for (let i = 0; i < 10; i++) {
            performanceMonitor.recordRequest(100, true);
        }

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.errorRate).toBe(0);
    });

    it('should handle 100% token refresh failure', () => {
        for (let i = 0; i < 5; i++) {
            performanceMonitor.recordTokenRefresh(false);
        }

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.tokenRefreshSuccessRate).toBe(0);
    });
});

describe('Performance Monitor - Data Integrity', () => {
    beforeEach(() => {
        performanceMonitor.reset();
    });

    it('should maintain accuracy across multiple operations', () => {
        // Record mixed operations
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.recordRequest(200, false);
        performanceMonitor.recordTokenRefresh(true);
        performanceMonitor.recordQueueRejection();
        performanceMonitor.recordRequest(150, true);
        performanceMonitor.recordTokenRefresh(false);

        const metrics1 = performanceMonitor.getMetrics();
        expect(metrics1.requestCount).toBe(3);
        expect(metrics1.errorCount).toBe(1);

        // Record more operations
        performanceMonitor.recordRequest(300, true);
        performanceMonitor.recordTokenRefresh(true);

        const metrics2 = performanceMonitor.getMetrics();
        expect(metrics2.requestCount).toBe(4);
        expect(metrics2.tokenRefreshCount).toBe(3);
    });

    it('should not leak data between resets', () => {
        performanceMonitor.recordRequest(100, true);
        performanceMonitor.reset();
        performanceMonitor.recordRequest(200, false);

        const metrics = performanceMonitor.getMetrics();
        expect(metrics.requestCount).toBe(1);
        expect(metrics.totalResponseTime).toBe(200);
        expect(metrics.errorCount).toBe(1);
    });
});
