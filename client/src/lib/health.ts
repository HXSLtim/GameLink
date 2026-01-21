/**
 * Health Check Utilities
 * Provides health check functionality for the client application
 */

interface HealthCheckResult {
    status: 'healthy' | 'degraded' | 'unhealthy';
    checks: {
        name: string;
        status: 'pass' | 'fail' | 'warn';
        message?: string;
        duration?: number;
    }[];
    timestamp: string;
}

export async function performHealthChecks(): Promise<HealthCheckResult> {
    const checks = [];

    // Check 1: API connectivity
    try {
        const start = Date.now();
        const response = await fetch('/api/v1/healthz', {
            method: 'HEAD',
        });
        const duration = Date.now() - start;

        if (response.ok) {
            checks.push({
                name: 'API Connectivity',
                status: 'pass',
                duration,
            });
        } else {
            checks.push({
                name: 'API Connectivity',
                status: 'fail',
                message: `HTTP ${response.status}`,
            });
        }
    } catch (error) {
        checks.push({
            name: 'API Connectivity',
            status: 'fail',
            message: error instanceof Error ? error.message : 'Unknown error',
        });
    }

    // Check 2: LocalStorage availability
    try {
        const testKey = '__health_check__';
        localStorage.setItem(testKey, 'test');
        localStorage.removeItem(testKey);
        checks.push({
            name: 'LocalStorage',
            status: 'pass',
        });
    } catch (_error) {
        checks.push({
            name: 'LocalStorage',
            status: 'fail',
            message: 'LocalStorage unavailable',
        });
    }

    // Check 3: Crypto configuration
    const cryptoEnabled = import.meta.env.VITE_CRYPTO_ENABLED === 'true';
    if (cryptoEnabled) {
        const hasKey = !!import.meta.env.VITE_CRYPTO_SECRET_KEY;
        const hasIV = !!import.meta.env.VITE_CRYPTO_IV;
        if (hasKey && hasIV) {
            checks.push({
                name: 'Crypto Configuration',
                status: 'pass',
            });
        } else {
            checks.push({
                name: 'Crypto Configuration',
                status: 'fail',
                message: 'Missing crypto keys',
            });
        }
    } else {
        checks.push({
            name: 'Crypto Configuration',
            status: 'warn',
            message: 'Encryption disabled',
        });
    }

    // Determine overall status
    const failedChecks = checks.filter(c => c.status === 'fail');
    const warnChecks = checks.filter(c => c.status === 'warn');

    let overallStatus: 'healthy' | 'degraded' | 'unhealthy';
    if (failedChecks.length > 0) {
        overallStatus = 'unhealthy';
    } else if (warnChecks.length > 0) {
        overallStatus = 'degraded';
    } else {
        overallStatus = 'healthy';
    }

    return {
        status: overallStatus,
        checks,
        timestamp: new Date().toISOString(),
    };
}

/**
 * Get health report as string
 */
export async function getHealthReport(): Promise<string> {
    const health = await performHealthChecks();

    const lines = [
        `Client Health Check - ${health.timestamp}`,
        `Overall Status: ${health.status.toUpperCase()}`,
        '',
        'Checks:',
    ];

    health.checks.forEach(check => {
        const icon = check.status === 'pass' ? '✅' : check.status === 'fail' ? '❌' : '⚠️';
        const duration = check.duration ? ` (${duration}ms)` : '';
        const message = check.message ? ` - ${check.message}` : '';
        lines.push(`  ${icon} ${check.name}${duration}${message}`);
    });

    return lines.join('\n');
}
