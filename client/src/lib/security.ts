/**
 * Security utilities for GameLink Client
 * Provides security validation, XSS protection, and security headers
 */

/**
 * Validate crypto configuration
 * @throws Error if configuration is invalid
 */
export function validateCryptoConfig(): void {
    const enabled = import.meta.env.VITE_CRYPTO_ENABLED === 'true';

    if (enabled) {
        const secretKey = import.meta.env.VITE_CRYPTO_SECRET_KEY;
        const iv = import.meta.env.VITE_CRYPTO_IV;

        if (!secretKey || secretKey.length !== 32) {
            throw new Error(
                'VITE_CRYPTO_SECRET_KEY must be exactly 32 characters. ' +
                    `Current length: ${secretKey?.length || 0}`
            );
        }

        if (!iv || iv.length !== 16) {
            throw new Error(
                'VITE_CRYPTO_IV must be exactly 16 characters. ' +
                    `Current length: ${iv?.length || 0}`
            );
        }
    }
}

/**
 * Sanitize user input to prevent XSS attacks
 * @param input - Raw user input
 * @returns Sanitized string
 */
export function sanitizeInput(input: string): string {
    if (!input) return '';

    // Remove potentially dangerous HTML/JS
    return input
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#x27;')
        .replace(/\//g, '&#x2F;');
}

/**
 * Validate URL to prevent XSS
 * @param url - URL to validate
 * @returns true if URL is safe
 */
export function isValidUrl(url: string): boolean {
    if (!url) return false;

    try {
        const parsed = new URL(url);
        // Only allow http and https protocols
        return ['http:', 'https:'].includes(parsed.protocol);
    } catch {
        return false;
    }
}

/**
 * Check if domain is trusted
 * @param url - URL to check
 * @returns true if domain is in trusted list
 */
export function isTrustedDomain(url: string): boolean {
    try {
        const parsed = new URL(url);
        const trustedDomains = import.meta.env.VITE_TRUSTED_DOMAINS?.split(',') || [
            'gamelink.com',
            'api.gamelink.com',
        ];

        return trustedDomains.some((domain) => parsed.hostname.endsWith(domain));
    } catch {
        return false;
    }
}

/**
 * Check if current environment is secure
 * @returns true if using HTTPS (or localhost)
 */
export function isSecureEnvironment(): boolean {
    if (typeof window === 'undefined') return true; // SSR

    return (
        window.location.protocol === 'https:' ||
        window.location.hostname === 'localhost' ||
        window.location.hostname === '127.0.0.1'
    );
}

/**
 * Get security headers configuration
 * @returns Object containing security headers
 */
export function getSecurityHeaders() {
    return {
        'X-Content-Type-Options': 'nosniff',
        'X-Frame-Options': import.meta.env.VITE_FRAME_OPTIONS || 'DENY',
        'X-XSS-Protection': '1; mode=block',
        'Strict-Transport-Security': `max-age=${import.meta.env.VITE_HSTS_MAX_AGE || 31536000}${
            import.meta.env.VITE_HSTS_INCLUDE_SUB_DOMAINS === 'true' ? '; includeSubDomains' : ''
        }${import.meta.env.VITE_HSTS_PRELOAD === 'true' ? '; preload' : ''}`,
    };
}

/**
 * Get CSP policy as object
 * @returns CSP policy object
 */
export function getCSPPolicy(): Record<string, string> {
    const defaultPolicy = {
        'default-src': "'self'",
        'script-src': "'self' 'unsafe-inline' 'unsafe-eval'",
        'style-src': "'self' 'unsafe-inline'",
        'img-src': "'self' data: https:",
        'connect-src': "'self'",
        'font-src': "'self' data:",
        'object-src': "'none'",
        'base-uri': "'self'",
        'form-action': "'self'",
        'frame-ancestors': "'none'",
    };

    try {
        const policyJson = import.meta.env.VITE_CSP_POLICY;
        if (policyJson) {
            return JSON.parse(policyJson);
        }
    } catch (error) {
        logSecurityEvent('Failed to parse CSP policy', error);
    }

    return defaultPolicy;
}

/**
 * Generate CSP meta tag content
 * @returns CSP meta tag content string
 */
export function getCSPMetaTag(): string {
    const policy = getCSPPolicy();
    return Object.entries(policy)
        .map(([directive, value]) => `${directive} ${value}`)
        .join('; ');
}

/**
 * Log security event (development only)
 * @param event - Security event description
 * @param data - Additional data
 */
export function logSecurityEvent(event: string, data?: unknown): void {
    if (import.meta.env.DEV && import.meta.env.VITE_SECURITY_LOGGING_ENABLED === 'true') {
        console.warn(`[Security] ${event}`, data);
    }
}

/**
 * Validate file upload
 * @param file - File to validate
 * @returns true if file is valid
 */
export function validateFileUpload(file: File): { valid: boolean; error?: string } {
    const maxSize = Number(import.meta.env.VITE_MAX_FILE_SIZE) || 10485760; // 10MB default
    const allowedTypes = import.meta.env.VITE_ALLOWED_FILE_TYPES?.split(',') || [
        'image/jpeg',
        'image/png',
        'image/gif',
        'image/webp',
    ];

    // Check file size
    if (file.size > maxSize) {
        return {
            valid: false,
            error: `File size exceeds maximum allowed size of ${maxSize / 1024 / 1024}MB`,
        };
    }

    // Check file type
    if (!allowedTypes.includes(file.type)) {
        return {
            valid: false,
            error: `File type ${file.type} is not allowed`,
        };
    }

    return { valid: true };
}

/**
 * Sanitize HTML content (basic version)
 * For production, consider using a library like DOMPurify
 * @param html - Raw HTML content
 * @returns Sanitized HTML
 */
export function sanitizeHTML(html: string): string {
    // Basic sanitization - remove script tags and event handlers
    return html
        .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
        .replace(/on\w+="[^"]*"/gi, '')
        .replace(/on\w+='[^']*'/gi, '')
        .replace(/javascript:/gi, '');
}

/**
 * Check if token storage is secure
 * @returns true if using secure storage method
 */
export function isTokenStorageSecure(): boolean {
    const storage = import.meta.env.VITE_TOKEN_STORAGE || 'localStorage';
    return storage === 'httpOnly' || storage === 'memory';
}

/**
 * Get token storage method
 * @returns Storage method name
 */
export function getTokenStorageMethod(): 'localStorage' | 'sessionStorage' | 'memory' {
    return (import.meta.env.VITE_TOKEN_STORAGE as any) || 'localStorage';
}

/**
 * Validate session timeout
 * @param lastActivity - Last activity timestamp
 * @returns true if session is still valid
 */
export function isSessionValid(lastActivity: number): boolean {
    const idleTimeout = Number(import.meta.env.VITE_IDLE_TIMEOUT) || 1800000; // 30 minutes default
    const now = Date.now();
    return now - lastActivity < idleTimeout;
}

/**
 * Get session timeout warning time
 * @returns Warning time in milliseconds before timeout
 */
export function getSessionTimeoutWarning(): number {
    return Number(import.meta.env.VITE_SESSION_TIMEOUT_WARNING) || 300000; // 5 minutes default
}

/**
 * Rate limiter for API requests
 */
export class RateLimiter {
    private requests: number[] = [];
    private maxRequests: number;
    private windowMs: number;

    constructor(maxRequests: number = 60, windowMs: number = 60000) {
        this.maxRequests = maxRequests;
        this.windowMs = windowMs;
    }

    /**
     * Check if request is allowed
     * @returns true if within rate limit
     */
    canMakeRequest(): boolean {
        const now = Date.now();
        // Remove old requests outside the window
        this.requests = this.requests.filter((time) => now - time < this.windowMs);

        if (this.requests.length < this.maxRequests) {
            this.requests.push(now);
            return true;
        }

        return false;
    }

    /**
     * Get remaining requests
     * @returns Number of remaining requests
     */
    getRemainingRequests(): number {
        const now = Date.now();
        this.requests = this.requests.filter((time) => now - time < this.windowMs);
        return Math.max(0, this.maxRequests - this.requests.length);
    }

    /**
     * Reset rate limiter
     */
    reset(): void {
        this.requests = [];
    }
}

/**
 * Create a rate limiter instance
 * @returns Rate limiter instance
 */
export function createRateLimiter(): RateLimiter {
    const maxRequests = Number(import.meta.env.VITE_RATE_LIMIT_PER_MINUTE) || 60;
    return new RateLimiter(maxRequests, 60000); // 1 minute window
}

/**
 * Security check for external links
 * @param url - URL to check
 * @returns true if link is safe to open
 */
export function isExternalLinkSafe(url: string): boolean {
    if (!url) return false;

    if (!isValidUrl(url)) {
        logSecurityEvent('Invalid URL detected', { url });
        return false;
    }

    // Check if it's an external link
    try {
        const parsed = new URL(url);
        const currentHost = typeof window !== 'undefined' ? window.location.hostname : '';

        // Same origin is safe
        if (parsed.hostname === currentHost) {
            return true;
        }

        // Check trusted domains
        if (isTrustedDomain(url)) {
            return true;
        }

        // External link - should show warning
        if (import.meta.env.VITE_EXTERNAL_LINK_WARNING === 'true') {
            logSecurityEvent('External link detected', { url });
        }

        return true;
    } catch {
        return false;
    }
}

/**
 * Validate environment on app start
 * @throws Error if security configuration is invalid
 */
export function validateSecurityConfig(): void {
    // Check crypto configuration
    validateCryptoConfig();

    // Check HTTPS requirement in production
    if (
        import.meta.env.PROD &&
        import.meta.env.VITE_REQUIRE_HTTPS === 'true' &&
        !isSecureEnvironment()
    ) {
        console.error('Production environment requires HTTPS');
    }

    // Check token storage
    const tokenStorage = getTokenStorageMethod();
    if (!['localStorage', 'sessionStorage', 'memory'].includes(tokenStorage)) {
        throw new Error(`Invalid token storage method: ${tokenStorage}`);
    }

    // Log security configuration
    logSecurityEvent('Security configuration validated', {
        cryptoEnabled: import.meta.env.VITE_CRYPTO_ENABLED,
        tokenStorage,
        cspEnabled: import.meta.env.VITE_CSP_ENABLED,
    });
}

/**
 * Content Security Policy types
 */
export interface CSPDirective {
    name: string;
    value: string;
}

/**
 * Build CSP header value
 * @param directives - Array of CSP directives
 * @returns CSP header value
 */
export function buildCSPHeader(directives: CSPDirective[]): string {
    return directives.map((d) => `${d.name} ${d.value}`).join('; ');
}
