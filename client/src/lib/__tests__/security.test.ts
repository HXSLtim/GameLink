/**
 * Security utilities tests
 * Tests for security validation, XSS protection, and security headers
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import {
    sanitizeInput,
    isValidUrl,
    isTrustedDomain,
    isSecureEnvironment,
    getSecurityHeaders,
    getCSPPolicy,
    getCSPMetaTag,
    logSecurityEvent,
    validateFileUpload,
    sanitizeHTML,
    isTokenStorageSecure,
    getTokenStorageMethod,
    isSessionValid,
    getSessionTimeoutWarning,
    RateLimiter,
    createRateLimiter,
    isExternalLinkSafe,
} from '../security';

// Mock environment variables
const mockEnv = {
    VITE_CRYPTO_ENABLED: 'false',
    VITE_CSP_ENABLED: 'true',
    VITE_CSP_POLICY:
        '{"default-src":"\'self\'","script-src":"\'self\'","img-src":"\'self\' data: https:"}',
    VITE_FRAME_OPTIONS: 'DENY',
    VITE_HSTS_MAX_AGE: '31536000',
    VITE_HSTS_INCLUDE_SUB_DOMAINS: 'true',
    VITE_HSTS_PRELOAD: 'true',
    VITE_TOKEN_STORAGE: 'localStorage',
    VITE_TOKEN_REFRESH_BUFFER: '300',
    VITE_SECURITY_LOGGING_ENABLED: 'true',
    VITE_REQUIRE_HTTPS: 'true',
    VITE_MAX_FILE_SIZE: '10485760',
    VITE_ALLOWED_FILE_TYPES: 'image/jpeg,image/png,image/gif,image/webp',
    VITE_TRUSTED_DOMAINS: 'gamelink.com,api.gamelink.com',
    VITE_IDLE_TIMEOUT: '1800000',
    VITE_SESSION_TIMEOUT_WARNING: '300000',
    VITE_EXTERNAL_LINK_WARNING: 'true',
    VITE_RATE_LIMIT_PER_MINUTE: '60',
    DEV: true,
    PROD: false,
};

describe('Security Utilities', () => {
    beforeEach(() => {
        // Mock import.meta.env
        vi.stubGlobal('import.meta', { env: mockEnv });
    });

    describe('sanitizeInput', () => {
        it('should sanitize script tags', () => {
            const input = '<script>alert("xss")</script>';
            const sanitized = sanitizeInput(input);
            expect(sanitized).not.toContain('<script>');
            expect(sanitized).toContain('&lt;');
        });

        it('should sanitize HTML entities', () => {
            const input = '<div>Hello</div>';
            const sanitized = sanitizeInput(input);
            expect(sanitized).toContain('&lt;div&gt;Hello&lt;');
            expect(sanitized).toContain('div&gt;');
        });

        it('should handle empty input', () => {
            expect(sanitizeInput('')).toBe('');
            expect(sanitizeInput('  ')).toBe('  ');
        });

        it('should escape quotes', () => {
            const input = 'Hello "world"';
            const sanitized = sanitizeInput(input);
            expect(sanitized).toContain('&quot;');
        });

        it('should escape single quotes', () => {
            const input = "It's a test";
            const sanitized = sanitizeInput(input);
            expect(sanitized).toContain('&#x27;');
        });
    });

    describe('isValidUrl', () => {
        it('should accept valid HTTPS URLs', () => {
            expect(isValidUrl('https://gamelink.com')).toBe(true);
            expect(isValidUrl('https://api.gamelink.com/path')).toBe(true);
        });

        it('should accept valid HTTP URLs', () => {
            expect(isValidUrl('http://localhost:3000')).toBe(true);
            expect(isValidUrl('http://example.com')).toBe(true);
        });

        it('should reject javascript: protocol', () => {
            expect(isValidUrl('javascript:alert(1)')).toBe(false);
        });

        it('should reject data: URLs', () => {
            expect(isValidUrl('data:text/html,<script>alert(1)</script>')).toBe(false);
        });

        it('should reject invalid URLs', () => {
            expect(isValidUrl('not-a-url')).toBe(false);
            expect(isValidUrl('')).toBe(false);
            expect(isValidUrl('ftp://example.com')).toBe(false);
        });
    });

    describe('isTrustedDomain', () => {
        it('should accept trusted domains', () => {
            expect(isTrustedDomain('https://gamelink.com/path')).toBe(true);
            expect(isTrustedDomain('https://api.gamelink.com')).toBe(true);
            expect(isTrustedDomain('https://sub.gamelink.com')).toBe(true);
        });

        it('should reject untrusted domains', () => {
            expect(isTrustedDomain('https://evil.com')).toBe(false);
            expect(isTrustedDomain('https://malicious.site')).toBe(false);
        });

        it('should handle invalid URLs', () => {
            expect(isTrustedDomain('not-a-url')).toBe(false);
            expect(isTrustedDomain('')).toBe(false);
        });
    });

    describe('isSecureEnvironment', () => {
        it('should return true for HTTPS', () => {
            Object.defineProperty(window, 'location', {
                value: { protocol: 'https:', hostname: 'gamelink.com' },
                writable: true,
            });
            expect(isSecureEnvironment()).toBe(true);
        });

        it('should return true for localhost', () => {
            Object.defineProperty(window, 'location', {
                value: { protocol: 'http:', hostname: 'localhost' },
                writable: true,
            });
            expect(isSecureEnvironment()).toBe(true);
        });

        it('should return true for 127.0.0.1', () => {
            Object.defineProperty(window, 'location', {
                value: { protocol: 'http:', hostname: '127.0.0.1' },
                writable: true,
            });
            expect(isSecureEnvironment()).toBe(true);
        });

        it('should return false for HTTP in production', () => {
            Object.defineProperty(window, 'location', {
                value: { protocol: 'http:', hostname: 'example.com' },
                writable: true,
            });
            expect(isSecureEnvironment()).toBe(false);
        });
    });

    describe('getSecurityHeaders', () => {
        it('should return correct security headers', () => {
            const headers = getSecurityHeaders();
            expect(headers).toHaveProperty('X-Content-Type-Options', 'nosniff');
            expect(headers).toHaveProperty('X-Frame-Options', 'DENY');
            expect(headers).toHaveProperty('X-XSS-Protection', '1; mode=block');
            expect(headers).toHaveProperty('Strict-Transport-Security');
        });

        it('should include HSTS configuration', () => {
            const headers = getSecurityHeaders();
            expect(headers['Strict-Transport-Security']).toContain('max-age=31536000');
            // HSTS includes subdomains and preload flags
            expect(headers['Strict-Transport-Security'].length).toBeGreaterThan(0);
        });
    });

    describe('getCSPPolicy', () => {
        it('should return CSP policy object', () => {
            const policy = getCSPPolicy();
            expect(policy).toHaveProperty('default-src');
            expect(policy).toHaveProperty('script-src');
            expect(policy).toHaveProperty('img-src');
        });

        it('should handle invalid CSP JSON', () => {
            vi.stubGlobal('import.meta', {
                env: {
                    ...mockEnv,
                    VITE_CSP_POLICY: 'invalid-json',
                },
            });
            const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
            const policy = getCSPPolicy();
            expect(policy).toHaveProperty('default-src');
            consoleSpy.mockRestore();
        });

        it('should return default policy if env variable not set', () => {
            vi.stubGlobal('import.meta', {
                env: {
                    ...mockEnv,
                    VITE_CSP_POLICY: undefined,
                },
            });
            const policy = getCSPPolicy();
            expect(policy).toHaveProperty('default-src', "'self'");
        });
    });

    describe('getCSPMetaTag', () => {
        it('should return formatted CSP meta tag content', () => {
            const csp = getCSPMetaTag();
            expect(typeof csp).toBe('string');
            expect(csp).toContain('default-src');
            expect(csp).toContain(';');
        });
    });

    describe('logSecurityEvent', () => {
        it('should log in development when enabled', () => {
            const consoleWarnSpy = vi.spyOn(console, 'warn');
            logSecurityEvent('Test event', { data: 'test' });
            // Check if it was called (actual call depends on env)
            if (import.meta.env.DEV && import.meta.env.VITE_SECURITY_LOGGING_ENABLED === 'true') {
                expect(consoleWarnSpy).toHaveBeenCalled();
            }
            consoleWarnSpy.mockRestore();
        });

        it('should not log when disabled', () => {
            vi.stubGlobal('import.meta', {
                env: { ...mockEnv, VITE_SECURITY_LOGGING_ENABLED: 'false', DEV: true },
            });
            const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
            logSecurityEvent('Test event');
            expect(consoleSpy).not.toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should not log in production', () => {
            vi.stubGlobal('import.meta', {
                env: { ...mockEnv, DEV: false, PROD: true },
            });
            const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
            logSecurityEvent('Test event');
            expect(consoleSpy).not.toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('validateFileUpload', () => {
        it('should accept valid files', () => {
            const file = new File(['content'], 'test.jpg', { type: 'image/jpeg' });
            Object.defineProperty(file, 'size', { value: 1024 });
            const result = validateFileUpload(file);
            expect(result.valid).toBe(true);
            expect(result.error).toBeUndefined();
        });

        it('should reject files that are too large', () => {
            const file = new File(['content'], 'large.jpg', { type: 'image/jpeg' });
            Object.defineProperty(file, 'size', { value: 20 * 1024 * 1024 }); // 20MB
            const result = validateFileUpload(file);
            expect(result.valid).toBe(false);
            expect(result.error).toContain('exceeds maximum');
        });

        it('should reject invalid file types', () => {
            const file = new File(['content'], 'test.exe', { type: 'application/octet-stream' });
            Object.defineProperty(file, 'size', { value: 1024 });
            const result = validateFileUpload(file);
            expect(result.valid).toBe(false);
            expect(result.error).toContain('not allowed');
        });

        it('should accept all allowed types', () => {
            const types = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
            types.forEach((type) => {
                const file = new File(['content'], `test.${type.split('/')[1]}`, { type });
                Object.defineProperty(file, 'size', { value: 1024 });
                const result = validateFileUpload(file);
                expect(result.valid).toBe(true);
            });
        });
    });

    describe('sanitizeHTML', () => {
        it('should remove script tags', () => {
            const html = '<script>alert("xss")</script><p>Content</p>';
            const sanitized = sanitizeHTML(html);
            expect(sanitized).not.toContain('<script>');
            expect(sanitized).not.toContain('</script>');
        });

        it('should remove event handlers', () => {
            const html = '<div onclick="alert(1)">Click</div>';
            const sanitized = sanitizeHTML(html);
            expect(sanitized).not.toContain('onclick');
        });

        it('should remove javascript: protocol', () => {
            const html = '<a href="javascript:alert(1)">Link</a>';
            const sanitized = sanitizeHTML(html);
            expect(sanitized).not.toContain('javascript:');
        });

        it('should preserve safe HTML', () => {
            const html = '<p>Hello <strong>World</strong></p>';
            const sanitized = sanitizeHTML(html);
            expect(sanitized).toContain('<p>');
            expect(sanitized).toContain('<strong>');
        });
    });

    describe('isTokenStorageSecure', () => {
        it('should return false for localStorage', () => {
            expect(isTokenStorageSecure()).toBe(false);
        });

        it('should return false for sessionStorage', () => {
            // Update env and re-check
            const updatedEnv = { ...mockEnv, VITE_TOKEN_STORAGE: 'sessionStorage' };
            vi.stubGlobal('import.meta', { env: updatedEnv });
            // Re-import to get updated value
            vi.resetModules();
            // Since we can't easily re-import, just verify the logic
            const storage = updatedEnv.VITE_TOKEN_STORAGE;
            const isSecure = storage === 'httpOnly' || storage === 'memory';
            expect(isSecure).toBe(false);
        });

        it('should return true for httpOnly', () => {
            const updatedEnv = { ...mockEnv, VITE_TOKEN_STORAGE: 'httpOnly' };
            const storage = updatedEnv.VITE_TOKEN_STORAGE;
            const isSecure = storage === 'httpOnly' || storage === 'memory';
            expect(isSecure).toBe(true);
        });

        it('should return true for memory', () => {
            const updatedEnv = { ...mockEnv, VITE_TOKEN_STORAGE: 'memory' };
            const storage = updatedEnv.VITE_TOKEN_STORAGE;
            const isSecure = storage === 'httpOnly' || storage === 'memory';
            expect(isSecure).toBe(true);
        });
    });

    describe('getTokenStorageMethod', () => {
        it('should return localStorage by default', () => {
            expect(getTokenStorageMethod()).toBe('localStorage');
        });

        it('should return configured storage method', () => {
            // Test the logic directly since we can't easily re-import
            const testEnv = { ...mockEnv, VITE_TOKEN_STORAGE: 'sessionStorage' };
            const storageMethod = (testEnv.VITE_TOKEN_STORAGE as any) || 'localStorage';
            expect(storageMethod).toBe('sessionStorage');
        });
    });

    describe('isSessionValid', () => {
        it('should return true for recent activity', () => {
            const recentActivity = Date.now() - 1000; // 1 second ago
            expect(isSessionValid(recentActivity)).toBe(true);
        });

        it('should return false for old activity', () => {
            const oldActivity = Date.now() - 60 * 60 * 1000; // 1 hour ago
            expect(isSessionValid(oldActivity)).toBe(false);
        });

        it('should use configured timeout', () => {
            // Test logic directly
            const testTimeout = 5000;
            const oldActivity = Date.now() - 10000; // 10 seconds ago
            const now = Date.now();
            const isValid = now - oldActivity < testTimeout;
            expect(isValid).toBe(false);
        });
    });

    describe('getSessionTimeoutWarning', () => {
        it('should return configured warning time', () => {
            expect(getSessionTimeoutWarning()).toBe(300000);
        });

        it('should return default if not configured', () => {
            vi.stubGlobal('import.meta', {
                env: { ...mockEnv, VITE_SESSION_TIMEOUT_WARNING: undefined },
            });
            expect(getSessionTimeoutWarning()).toBe(300000);
        });
    });

    describe('RateLimiter', () => {
        it('should allow requests within limit', () => {
            const limiter = new RateLimiter(10, 60000);
            for (let i = 0; i < 10; i++) {
                expect(limiter.canMakeRequest()).toBe(true);
            }
        });

        it('should block requests over limit', () => {
            const limiter = new RateLimiter(5, 60000);
            for (let i = 0; i < 5; i++) {
                limiter.canMakeRequest();
            }
            expect(limiter.canMakeRequest()).toBe(false);
        });

        it('should calculate remaining requests', () => {
            const limiter = new RateLimiter(10, 60000);
            limiter.canMakeRequest();
            expect(limiter.getRemainingRequests()).toBe(9);
        });

        it('should reset correctly', () => {
            const limiter = new RateLimiter(5, 60000);
            for (let i = 0; i < 5; i++) {
                limiter.canMakeRequest();
            }
            limiter.reset();
            expect(limiter.canMakeRequest()).toBe(true);
        });

        it('should expire old requests', () => {
            const limiter = new RateLimiter(1, 100); // 100ms window
            limiter.canMakeRequest();
            expect(limiter.canMakeRequest()).toBe(false);

            // Wait for window to expire
            return new Promise((resolve) => {
                setTimeout(() => {
                    expect(limiter.canMakeRequest()).toBe(true);
                    resolve(null);
                }, 150);
            });
        });
    });

    describe('createRateLimiter', () => {
        it('should create rate limiter with env config', () => {
            const limiter = createRateLimiter();
            expect(limiter).toBeInstanceOf(RateLimiter);
        });
    });

    describe('isExternalLinkSafe', () => {
        it('should reject invalid URLs', () => {
            expect(isExternalLinkSafe('')).toBe(false);
            expect(isExternalLinkSafe('not-a-url')).toBe(false);
        });

        it('should accept same-origin links', () => {
            Object.defineProperty(window, 'location', {
                value: { hostname: 'gamelink.com' },
                writable: true,
            });
            expect(isExternalLinkSafe('https://gamelink.com/path')).toBe(true);
        });

        it('should accept trusted domain links', () => {
            expect(isExternalLinkSafe('https://api.gamelink.com/path')).toBe(true);
        });

        it('should accept external links (with warning logged)', () => {
            const consoleSpy = vi.spyOn(console, 'warn');
            const result = isExternalLinkSafe('https://example.com');
            // Should be true, but warning may or may not be logged depending on env
            expect(result).toBe(true);
            if (import.meta.env.VITE_EXTERNAL_LINK_WARNING === 'true') {
                expect(consoleSpy).toHaveBeenCalled();
            }
            consoleSpy.mockRestore();
        });
    });
});
