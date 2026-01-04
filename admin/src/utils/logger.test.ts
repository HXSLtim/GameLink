/**
 * Logger Utility Tests
 *
 * Coverage Target: 90%+
 *
 * Test Scenarios:
 * 1. Environment-aware logging (development vs production)
 * 2. All log levels (info, warn, error, debug)
 * 3. Specialized logging methods (api, apiResponse, userAction, lifecycle)
 * 4. Error handling and stack traces
 * 5. Context formatting
 * 6. Production mode suppression
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import { logger } from './logger';

// Mock console methods
const consoleMocks = {
    log: vi.spyOn(console, 'log').mockImplementation(() => {}),
    warn: vi.spyOn(console, 'warn').mockImplementation(() => {}),
    error: vi.spyOn(console, 'error').mockImplementation(() => {}),
};

describe('logger utility', () => {
    // Mock import.meta.env.MODE
    const originalEnv = import.meta.env.MODE;

    beforeEach(() => {
        vi.clearAllMocks();
        consoleMocks.log.mockClear();
        consoleMocks.warn.mockClear();
        consoleMocks.error.mockClear();
    });

    describe('environment-aware logging', () => {
        it('should log in development mode', () => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });

            logger.info('test message');

            expect(consoleMocks.log).toHaveBeenCalled();
        });

        it('should suppress logs in production mode', () => {
            vi.stubGlobal('import.meta', { env: { MODE: 'production' } });

            logger.info('test message');
            logger.warn('test warning');
            logger.error('test error');
            logger.debug('test debug');

            expect(consoleMocks.log).not.toHaveBeenCalled();
            expect(consoleMocks.warn).not.toHaveBeenCalled();
            expect(consoleMocks.error).not.toHaveBeenCalled();
        });

        it('should log in test mode', () => {
            vi.stubGlobal('import.meta', { env: { MODE: 'test' } });

            logger.info('test message');

            expect(consoleMocks.log).toHaveBeenCalled();
        });
    });

    describe('info logging', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should log info messages', () => {
            logger.info('User logged in');

            expect(consoleMocks.log).toHaveBeenCalledWith(
                expect.stringContaining('[INFO]'),
                'User logged in',
                ''
            );
        });

        it('should include context', () => {
            logger.info('User action', { userId: 123, action: 'login' });

            expect(consoleMocks.log).toHaveBeenCalledWith(
                expect.stringContaining('[INFO]'),
                'User action',
                { userId: 123, action: 'login' }
            );
        });

        it('should handle array context', () => {
            logger.info('Multiple items', [1, 2, 3]);

            expect(consoleMocks.log).toHaveBeenCalledWith(
                expect.stringContaining('[INFO]'),
                'Multiple items',
                [1, 2, 3]
            );
        });

        it('should handle null/undefined context', () => {
            logger.info('Test', null);
            logger.info('Test', undefined);

            expect(consoleMocks.log).toHaveBeenCalledTimes(2);
        });
    });

    describe('warn logging', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should log warning messages', () => {
            logger.warn('API rate limit approaching');

            expect(consoleMocks.warn).toHaveBeenCalledWith(
                expect.stringContaining('[WARN]'),
                'API rate limit approaching',
                ''
            );
        });

        it('should include context in warnings', () => {
            logger.warn('Rate limit', { remaining: 5, limit: 100 });

            expect(consoleMocks.warn).toHaveBeenCalledWith(
                expect.stringContaining('[WARN]'),
                'Rate limit',
                { remaining: 5, limit: 100 }
            );
        });
    });

    describe('error logging', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should log error messages', () => {
            logger.error('Failed to fetch data');

            expect(consoleMocks.error).toHaveBeenCalledWith(
                expect.stringContaining('[ERROR]'),
                'Failed to fetch data',
                '',
                ''
            );
        });

        it('should log error with Error object', () => {
            const error = new Error('Network error');
            logger.error('Request failed', error);

            expect(consoleMocks.error).toHaveBeenCalledWith(
                expect.stringContaining('[ERROR]'),
                'Request failed',
                '',
                error
            );
        });

        it('should log error with context', () => {
            logger.error('Request failed', { url: '/api/users', status: 500 });

            expect(consoleMocks.error).toHaveBeenCalledWith(
                expect.stringContaining('[ERROR]'),
                'Request failed',
                { url: '/api/users', status: 500 },
                ''
            );
        });

        it('should log error with both error and context', () => {
            const error = new Error('Validation failed');
            logger.error('Form submission failed', error, { field: 'email' });

            expect(consoleMocks.error).toHaveBeenCalledWith(
                expect.stringContaining('[ERROR]'),
                'Form submission failed',
                { field: 'email' },
                error
            );
        });

        it('should include stack trace in development mode', () => {
            const error = new Error('Test error');
            logger.error('Test', error);

            expect(consoleMocks.error).toHaveBeenCalled();
            // Check that stack trace is logged
            const calls = consoleMocks.error.mock.calls;
            const stackCall = calls.find(call =>
                Array.isArray(call) && call.some(arg =>
                    typeof arg === 'string' && arg.includes('Stack trace:')
                )
            );
            expect(stackCall).toBeDefined();
        });
    });

    describe('debug logging', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should log debug messages in development', () => {
            logger.debug('Debug info');

            expect(consoleMocks.log).toHaveBeenCalledWith(
                expect.stringContaining('[DEBUG]'),
                'Debug info',
                ''
            );
        });

        it('should not log debug messages in production', () => {
            vi.stubGlobal('import.meta', { env: { MODE: 'production' } });

            logger.debug('Debug info');

            expect(consoleMocks.log).not.toHaveBeenCalled();
        });
    });

    describe('specialized logging methods', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        describe('api logging', () => {
            it('should log API requests', () => {
                logger.api('POST', '/api/users', { body: { name: 'John' } });

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('POST'),
                    { url: '/api/users', body: { name: 'John' } }
                );
            });

            it('should log API requests without context', () => {
                logger.api('GET', '/api/users');

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('GET'),
                    { url: '/api/users' }
                );
            });

            it('should handle non-object context', () => {
                logger.api('GET', '/api/users', 'string context');

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('GET'),
                    { url: '/api/users' }
                );
            });
        });

        describe('apiResponse logging', () => {
            it('should log API responses', () => {
                logger.apiResponse('GET', '/api/users', 200, { duration: 150 });

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('GET'),
                    { url: '/api/users', status: 200, duration: 150 }
                );
            });

            it('should log API responses without context', () => {
                logger.apiResponse('POST', '/api/users', 201);

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('POST'),
                    { url: '/api/users', status: 201 }
                );
            });
        });

        describe('userAction logging', () => {
            it('should log user actions', () => {
                logger.userAction('click_button', { buttonId: 'submit' });

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[INFO]'),
                    'User action: click_button',
                    { buttonId: 'submit' }
                );
            });

            it('should log user actions without context', () => {
                logger.userAction('page_view');

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[INFO]'),
                    'User action: page_view',
                    ''
                );
            });
        });

        describe('lifecycle logging', () => {
            it('should log component mount', () => {
                logger.lifecycle('UserProfile', 'mount', { props: { userId: 1 } });

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('mount'),
                    { component: 'UserProfile', props: { userId: 1 } }
                );
            });

            it('should log component unmount', () => {
                logger.lifecycle('UserProfile', 'unmount');

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('unmount'),
                    { component: 'UserProfile' }
                );
            });

            it('should log component update', () => {
                logger.lifecycle('UserProfile', 'update', { changedProps: ['userId'] });

                expect(consoleMocks.log).toHaveBeenCalledWith(
                    expect.stringContaining('[DEBUG]'),
                    expect.stringContaining('update'),
                    { component: 'UserProfile', changedProps: ['userId'] }
                );
            });
        });
    });

    describe('timestamp formatting', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should include ISO timestamp in logs', () => {
            logger.info('Test message');

            const call = consoleMocks.log.mock.calls[0];
            const logMessage = call[0] as string;

            expect(logMessage).toMatch(/\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z\]/);
        });

        it('should include correct log level in timestamp prefix', () => {
            logger.info('info');
            logger.warn('warn');
            logger.error('error');
            logger.debug('debug');

            const infoCall = consoleMocks.log.mock.calls[0];
            const warnCall = consoleMocks.warn.mock.calls[0];
            const errorCall = consoleMocks.error.mock.calls[0];
            const debugCall = consoleMocks.log.mock.calls[1];

            expect(infoCall[0]).toContain('[INFO]');
            expect(warnCall[0]).toContain('[WARN]');
            expect(errorCall[0]).toContain('[ERROR]');
            expect(debugCall[0]).toContain('[DEBUG]');
        });
    });

    describe('property-based tests', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        /**
         * Property: Logger should handle any string message
         */
        it('should handle any string message', () => {
            fc.assert(
                fc.property(fc.string(), (message) => {
                    expect(() => logger.info(message)).not.toThrow();
                    expect(() => logger.warn(message)).not.toThrow();
                    expect(() => logger.debug(message)).not.toThrow();
                    expect(() => logger.error(message)).not.toThrow();
                    return true;
                }),
                { numRuns: 50 }
            );
        });

        /**
         * Property: Logger should handle various context types
         */
        it('should handle various context types', () => {
            fc.assert(
                fc.property(
                    fc.oneof(
                        fc.object(),
                        fc.array(fc.anything()),
                        fc.string(),
                        fc.nat(),
                        fc.boolean(),
                        fc.constant(null),
                        fc.constant(undefined)
                    ),
                    (context) => {
                        expect(() => logger.info('test', context)).not.toThrow();
                        expect(() => logger.warn('test', context)).not.toThrow();
                        expect(() => logger.debug('test', context)).not.toThrow();
                        return true;
                    }
                ),
                { numRuns: 50 }
            );
        });

        /**
         * Property: Error logging should handle Error-like objects
         */
        it('should handle Error-like objects', () => {
            fc.assert(
                fc.property(
                    fc.record({
                        message: fc.string(),
                        name: fc.string(),
                        stack: fc.string(),
                    }, { requiredKeys: ['message'] }),
                    (errorLike) => {
                        expect(() => logger.error('test', errorLike)).not.toThrow();
                        return true;
                    }
                ),
                { numRuns: 20 }
            );
        });

        /**
         * Property: API methods should handle any HTTP method
         */
        it('should handle any HTTP method string', () => {
            fc.assert(
                fc.property(fc.stringMatching(/^[A-Z]+$/), (method) => {
                    expect(() => logger.api(method, '/api/test')).not.toThrow();
                    expect(() => logger.apiResponse(method, '/api/test', 200)).not.toThrow();
                    return true;
                }),
                { numRuns: 20 }
            );
        });

        /**
         * Property: Logger should handle special characters in messages
         */
        it('should handle special characters', () => {
            fc.assert(
                fc.property(
                    fc.stringOf(fc.constantFrom(...'!@#$%^&*()[]{}\'"\\|;:,.<>?/~` \n\r\t'.split(''))),
                    (message) => {
                        expect(() => logger.info(message)).not.toThrow();
                        return true;
                    }
                ),
                { numRuns: 30 }
            );
        });
    });

    describe('edge cases', () => {
        beforeEach(() => {
            vi.stubGlobal('import.meta', { env: { MODE: 'development' } });
        });

        it('should handle empty strings', () => {
            expect(() => logger.info('')).not.toThrow();
            expect(() => logger.warn('')).not.toThrow();
            expect(() => logger.error('')).not.toThrow();
            expect(() => logger.debug('')).not.toThrow();
        });

        it('should handle very long messages', () => {
            const longMessage = 'x'.repeat(10000);
            expect(() => logger.info(longMessage)).not.toThrow();
        });

        it('should handle unicode characters', () => {
            expect(() => logger.info('你好世界🌍')).not.toThrow();
            expect(() => logger.info('🚀 🎉 ⭐')).not.toThrow();
            expect(() => logger.info('مرحبا بالعالم')).not.toThrow();
        });

        it('should handle circular references in context', () => {
            const circular: Record<string, unknown> = { a: 1 };
            circular.self = circular;

            expect(() => logger.info('test', circular)).not.toThrow();
        });

        it('should handle nested objects', () => {
            const nested = {
                level1: {
                    level2: {
                        level3: {
                            value: 'deep',
                        },
                    },
                },
            };

            expect(() => logger.info('test', nested)).not.toThrow();
        });

        it('should handle arrays of mixed types', () => {
            const mixed = [1, 'string', null, undefined, { key: 'value' }, [1, 2, 3]];
            expect(() => logger.info('test', mixed)).not.toThrow();
        });
    });

    describe('production mode behavior', () => {
        it('should not produce any console output in production', () => {
            vi.stubGlobal('import.meta', { env: { MODE: 'production' } });

            logger.info('info');
            logger.warn('warn');
            logger.error('error');
            logger.debug('debug');
            logger.api('GET', '/api/test');
            logger.apiResponse('GET', '/api/test', 200);
            logger.userAction('test');
            logger.lifecycle('TestComponent', 'mount');

            expect(consoleMocks.log).not.toHaveBeenCalled();
            expect(consoleMocks.warn).not.toHaveBeenCalled();
            expect(consoleMocks.error).not.toHaveBeenCalled();
        });
    });
});
