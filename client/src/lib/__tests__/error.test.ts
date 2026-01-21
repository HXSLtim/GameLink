/**
 * Error Handling Utilities Tests
 *
 * Comprehensive tests for error handling including:
 * - Error message extraction from various error types
 * - Error code mapping (40+ error codes)
 * - Error type identification (auth, validation, network, etc.)
 * - Axios error handling
 * - Store error handling helpers
 * - Error logging utilities
 */

/* eslint-disable react-hooks/rules-of-hooks */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { AxiosError } from 'axios';
import {
    isAxiosError,
    getErrorMessage,
    getErrorCode,
    getErrorStatus,
    logError,
    logWarn,
    isNetworkError,
    isAuthError,
    isForbiddenError,
    isNotFoundError,
    isValidationError,
    handleStoreError,
    ERROR_CODE_MESSAGES,
    getMessageByCode,
    isErrorCodeInRange,
    isAuthErrorByCode,
    isAuthorizationErrorByCode,
    isValidationErrorByCode,
    isBusinessErrorByCode,
    ErrorMessages,
    type ApiErrorResponse,
} from '../error';

/* eslint-disable react-hooks/rules-of-hooks */
// vi.mock callback uses 'use' parameter which triggers false positive from react-hooks plugin
vi.mock('../error', async (importOriginal) => {
    const actual = await importOriginal<typeof import('../error')>();
    return {
        ...actual,
    };
});

describe('Error Handling Utilities', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Mock console methods
        vi.spyOn(console, 'error').mockImplementation(() => {});
        vi.spyOn(console, 'warn').mockImplementation(() => {});
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('isAxiosError', () => {
        it('should identify AxiosError correctly', () => {
            const axiosError = new AxiosError('Test error');
            expect(isAxiosError(axiosError)).toBe(true);
        });

        it('should return false for regular Error', () => {
            const error = new Error('Regular error');
            expect(isAxiosError(error)).toBe(false);
        });

        it('should return false for non-error objects', () => {
            expect(isAxiosError(null)).toBe(false);
            expect(isAxiosError(undefined)).toBe(false);
            expect(isAxiosError({ message: 'error' })).toBe(false);
            expect(isAxiosError('error string')).toBe(false);
        });

        it('should work with typed AxiosError', () => {
            const axiosError = new AxiosError<ApiErrorResponse>('Test');
            expect(isAxiosError<ApiErrorResponse>(axiosError)).toBe(true);
        });
    });

    describe('getErrorMessage', () => {
        it('should extract message from AxiosError with response.data.message', () => {
            const axiosError = new AxiosError<ApiErrorResponse>('Network Error');
            axiosError.response = {
                data: { message: 'Invalid credentials' },
                status: 401,
                statusText: 'Unauthorized',
                headers: {},
                config: {} as any,
            };

            const message = getErrorMessage(axiosError);
            expect(message).toBe('Invalid credentials');
        });

        it('should extract message from AxiosError with response.data.error', () => {
            const axiosError = new AxiosError<ApiErrorResponse>('Network Error');
            axiosError.response = {
                data: { error: 'Forbidden access' },
                status: 403,
                statusText: 'Forbidden',
                headers: {},
                config: {} as any,
            };

            const message = getErrorMessage(axiosError);
            expect(message).toBe('Forbidden access');
        });

        it('should fall back to axios error message', () => {
            const axiosError = new AxiosError('Network timeout');
            axiosError.response = {
                data: {},
                status: 408,
                statusText: 'Timeout',
                headers: {},
                config: {} as any,
            };

            const message = getErrorMessage(axiosError);
            expect(message).toBe('Network timeout');
        });

        it('should extract message from standard Error object', () => {
            const error = new Error('Standard error message');
            expect(getErrorMessage(error)).toBe('Standard error message');
        });

        it('should return fallback for Error without message', () => {
            const error = new Error('');
            expect(getErrorMessage(error, 'Fallback message')).toBe('Fallback message');
        });

        it('should extract message from API error object', () => {
            const apiError = { message: 'API error occurred' };
            expect(getErrorMessage(apiError)).toBe('API error occurred');
        });

        it('should extract error field from API error object', () => {
            const apiError = { error: 'Error field message' };
            expect(getErrorMessage(apiError)).toBe('Error field message');
        });

        it('should return string error as-is', () => {
            expect(getErrorMessage('String error')).toBe('String error');
        });

        it('should return fallback for unknown error types', () => {
            expect(getErrorMessage(null, 'Unknown error')).toBe('Unknown error');
            expect(getErrorMessage(undefined, 'Default fallback')).toBe('Default fallback');
        });

        it('should return fallback for object without message fields', () => {
            const obj = { randomField: 'value' };
            expect(getErrorMessage(obj, 'No message')).toBe('No message');
        });

        it('should prioritize message over error field', () => {
            const apiError = { message: 'Primary message', error: 'Secondary error' };
            expect(getErrorMessage(apiError)).toBe('Primary message');
        });

        it('should handle null response.data in AxiosError', () => {
            const axiosError = new AxiosError('Network Error');
            axiosError.response = {
                data: null as any,
                status: 500,
                statusText: 'Internal Server Error',
                headers: {},
                config: {} as any,
            };

            const message = getErrorMessage(axiosError);
            expect(message).toBe('Network Error');
        });

        it('should handle undefined response.data', () => {
            const axiosError = new AxiosError('Request failed');
            axiosError.response = {
                data: undefined,
                status: 404,
                statusText: 'Not Found',
                headers: {},
                config: {} as any,
            };

            const message = getErrorMessage(axiosError);
            expect(message).toBe('Request failed');
        });
    });

    describe('getErrorCode', () => {
        it('should extract code from API error response', () => {
            const error = { code: 'ERR_001', message: 'Test error' };
            expect(getErrorCode(error)).toBe('ERR_001');
        });

        it('should return undefined for error without code', () => {
            const error = { message: 'Error without code' };
            expect(getErrorCode(error)).toBeUndefined();
        });

        it('should return undefined for non-object errors', () => {
            expect(getErrorCode(null)).toBeUndefined();
            expect(getErrorCode(undefined)).toBeUndefined();
            expect(getErrorCode('string')).toBeUndefined();
            expect(getErrorCode(123)).toBeUndefined();
        });
    });

    describe('getErrorStatus', () => {
        it('should extract status from API error response', () => {
            const error = { status: 404, message: 'Not found' };
            expect(getErrorStatus(error)).toBe(404);
        });

        it('should return undefined for error without status', () => {
            const error = { message: 'Error without status' };
            expect(getErrorStatus(error)).toBeUndefined();
        });

        it('should return undefined for non-object errors', () => {
            expect(getErrorStatus(null)).toBeUndefined();
            expect(getErrorStatus(undefined)).toBeUndefined();
            expect(getErrorStatus('string')).toBeUndefined();
        });
    });

    describe('Error Logging', () => {
        it('should log error in development mode', () => {
            const error = new Error('Test error');
            logError('testContext', error);

            expect(console.error).toHaveBeenCalledWith('[testContext]', error);
        });

        it('should log warning in development mode', () => {
            logWarn('testContext', 'Warning message');

            expect(console.warn).toHaveBeenCalledWith('[testContext]', 'Warning message');
        });

        it('should format context with brackets', () => {
            const error = new Error('Test');
            logError('fetchData', error);

            expect(console.error).toHaveBeenCalledWith('[fetchData]', error);
        });
    });

    describe('Error Type Checking', () => {
        describe('isNetworkError', () => {
            it('should identify network error by message', () => {
                const error = new Error('Network connection failed');
                expect(isNetworkError(error)).toBe(true);
            });

            it('should identify timeout error', () => {
                const error = new Error('Request timeout');
                expect(isNetworkError(error)).toBe(true);
            });

            it('should identify offline error', () => {
                const error = new Error('Device is offline');
                expect(isNetworkError(error)).toBe(true);
            });

            it('should identify failed to fetch error', () => {
                const error = new Error('Failed to fetch');
                expect(isNetworkError(error)).toBe(true);
            });

            it('should be case insensitive', () => {
                const error = new Error('NETWORK ERROR');
                expect(isNetworkError(error)).toBe(true);
            });

            it('should return false for non-network errors', () => {
                const error = new Error('Validation failed');
                expect(isNetworkError(error)).toBe(false);
            });

            it('should return false for non-Error objects', () => {
                expect(isNetworkError(null)).toBe(false);
                expect(isNetworkError('string')).toBe(false);
                expect(isNetworkError({ message: 'network' })).toBe(false);
            });
        });

        describe('isAuthError', () => {
            it('should identify 401 error', () => {
                const error = { status: 401, message: 'Unauthorized' };
                expect(isAuthError(error)).toBe(true);
            });

            it('should return false for other status codes', () => {
                expect(isAuthError({ status: 403 })).toBe(false);
                expect(isAuthError({ status: 404 })).toBe(false);
                expect(isAuthError({ status: 500 })).toBe(false);
            });

            it('should return false for errors without status', () => {
                expect(isAuthError({ message: 'Unauthorized' })).toBe(false);
                expect(isAuthError(null)).toBe(false);
            });
        });

        describe('isForbiddenError', () => {
            it('should identify 403 error', () => {
                const error = { status: 403, message: 'Forbidden' };
                expect(isForbiddenError(error)).toBe(true);
            });

            it('should return false for other status codes', () => {
                expect(isForbiddenError({ status: 401 })).toBe(false);
                expect(isForbiddenError({ status: 404 })).toBe(false);
            });
        });

        describe('isNotFoundError', () => {
            it('should identify 404 error', () => {
                const error = { status: 404, message: 'Not found' };
                expect(isNotFoundError(error)).toBe(true);
            });

            it('should return false for other status codes', () => {
                expect(isNotFoundError({ status: 400 })).toBe(false);
                expect(isNotFoundError({ status: 500 })).toBe(false);
            });
        });

        describe('isValidationError', () => {
            it('should identify 400 error', () => {
                const error = { status: 400, message: 'Validation failed' };
                expect(isValidationError(error)).toBe(true);
            });

            it('should return false for other status codes', () => {
                expect(isValidationError({ status: 401 })).toBe(false);
                expect(isValidationError({ status: 404 })).toBe(false);
            });
        });
    });

    describe('Error Code Mapping', () => {
        describe('ERROR_CODE_MESSAGES', () => {
            it('should have authentication error codes', () => {
                expect(ERROR_CODE_MESSAGES[40001]).toBe('用户名或密码错误');
                expect(ERROR_CODE_MESSAGES[40004]).toBe('Token已过期，请重新登录');
                expect(ERROR_CODE_MESSAGES[40005]).toBe('无效的Token');
            });

            it('should have authorization error codes', () => {
                expect(ERROR_CODE_MESSAGES[40101]).toBe('您没有权限执行此操作');
                expect(ERROR_CODE_MESSAGES[40102]).toBe('需要管理员权限');
            });

            it('should have validation error codes', () => {
                expect(ERROR_CODE_MESSAGES[40201]).toBe('手机号格式不正确');
                expect(ERROR_CODE_MESSAGES[40203]).toBe('密码长度不符合要求');
            });

            it('should have business logic error codes', () => {
                expect(ERROR_CODE_MESSAGES[40301]).toBe('余额不足，请先充值');
                expect(ERROR_CODE_MESSAGES[40302]).toBe('订单不存在');
                expect(ERROR_CODE_MESSAGES[40306]).toBe('优惠券不可用');
            });

            it('should have not found error codes', () => {
                expect(ERROR_CODE_MESSAGES[40401]).toBe('用户不存在');
                expect(ERROR_CODE_MESSAGES[40402]).toBe('陪玩师不存在');
            });

            it('should have server error codes', () => {
                expect(ERROR_CODE_MESSAGES[50001]).toBe('服务器繁忙，请稍后重试');
                expect(ERROR_CODE_MESSAGES[50002]).toBe('数据库连接失败');
            });

            it('should have 40+ error codes defined', () => {
                const codesCount = Object.keys(ERROR_CODE_MESSAGES).length;
                expect(codesCount).toBeGreaterThanOrEqual(40);
            });
        });

        describe('getMessageByCode', () => {
            it('should return message for known error code', () => {
                expect(getMessageByCode(40001)).toBe('用户名或密码错误');
                expect(getMessageByCode(40301)).toBe('余额不足，请先充值');
            });

            it('should return default message for unknown error code', () => {
                expect(getMessageByCode(99999)).toBe('操作失败，请稍后重试');
                expect(getMessageByCode(0)).toBe('操作失败，请稍后重试');
            });
        });

        describe('isErrorCodeInRange', () => {
            it('should return true for code in range', () => {
                expect(isErrorCodeInRange(40001, 40000, 40100)).toBe(true);
                expect(isErrorCodeInRange(40500, 40000, 50000)).toBe(true);
            });

            it('should return false for code below range', () => {
                expect(isErrorCodeInRange(39999, 40000, 50000)).toBe(false);
            });

            it('should return false for code at or above max', () => {
                expect(isErrorCodeInRange(50000, 40000, 50000)).toBe(false);
                expect(isErrorCodeInRange(50001, 40000, 50000)).toBe(false);
            });

            it('should handle edge cases', () => {
                expect(isErrorCodeInRange(40000, 40000, 40100)).toBe(true);
                expect(isErrorCodeInRange(40000, 40001, 40100)).toBe(false);
            });
        });

        describe('Error Type Classification by Code', () => {
            describe('isAuthErrorByCode', () => {
                it('should identify authentication errors (40000-40100)', () => {
                    expect(isAuthErrorByCode(40001)).toBe(true);
                    expect(isAuthErrorByCode(40099)).toBe(true);
                    expect(isAuthErrorByCode(40000)).toBe(true);
                });

                it('should reject non-authentication errors', () => {
                    expect(isAuthErrorByCode(39999)).toBe(false);
                    expect(isAuthErrorByCode(40100)).toBe(false);
                    expect(isAuthErrorByCode(40201)).toBe(false);
                });
            });

            describe('isAuthorizationErrorByCode', () => {
                it('should identify authorization errors (40100-40200)', () => {
                    expect(isAuthorizationErrorByCode(40101)).toBe(true);
                    expect(isAuthorizationErrorByCode(40199)).toBe(true);
                    expect(isAuthorizationErrorByCode(40100)).toBe(true);
                });

                it('should reject non-authorization errors', () => {
                    expect(isAuthorizationErrorByCode(40099)).toBe(false);
                    expect(isAuthorizationErrorByCode(40200)).toBe(false);
                });
            });

            describe('isValidationErrorByCode', () => {
                it('should identify validation errors (40200-40300)', () => {
                    expect(isValidationErrorByCode(40201)).toBe(true);
                    expect(isValidationErrorByCode(40299)).toBe(true);
                    expect(isValidationErrorByCode(40200)).toBe(true);
                });

                it('should reject non-validation errors', () => {
                    expect(isValidationErrorByCode(40199)).toBe(false);
                    expect(isValidationErrorByCode(40300)).toBe(false);
                });
            });

            describe('isBusinessErrorByCode', () => {
                it('should identify business logic errors (40300-40400)', () => {
                    expect(isBusinessErrorByCode(40301)).toBe(true);
                    expect(isBusinessErrorByCode(40399)).toBe(true);
                    expect(isBusinessErrorByCode(40300)).toBe(true);
                });

                it('should reject non-business errors', () => {
                    expect(isBusinessErrorByCode(40299)).toBe(false);
                    expect(isBusinessErrorByCode(40400)).toBe(false);
                });
            });
        });
    });

    describe('handleStoreError', () => {
        it('should extract error message and log error', () => {
            const error = new Error('Store action failed');
            const result = handleStoreError('testAction', error, 'Fallback');

            expect(result).toEqual({ error: 'Store action failed' });
            expect(console.error).toHaveBeenCalledWith('[testAction]', error);
        });

        it('should use fallback message when error message extraction fails', () => {
            const error = null;
            const result = handleStoreError('testAction', error, 'Default error');

            expect(result).toEqual({ error: 'Default error' });
        });

        it('should support disabling logging', () => {
            const error = new Error('Test error');
            const result = handleStoreError('testAction', error, 'Fallback', false);

            expect(result).toEqual({ error: 'Test error' });
            expect(console.error).not.toHaveBeenCalled();
        });

        it('should handle AxiosError correctly', () => {
            const axiosError = new AxiosError<ApiErrorResponse>('Request failed');
            axiosError.response = {
                data: { message: 'API error occurred' },
                status: 400,
                statusText: 'Bad Request',
                headers: {},
                config: {} as any,
            };

            const result = handleStoreError('apiCall', axiosError, 'Fallback');

            expect(result).toEqual({ error: 'API error occurred' });
        });

        it('should log context with error', () => {
            const error = new Error('Context test');
            handleStoreError('fetchUserData', error, 'Fallback');

            expect(console.error).toHaveBeenCalledWith('[fetchUserData]', error);
        });
    });

    describe('ErrorMessages Constants', () => {
        it('should have all required error message categories', () => {
            expect(ErrorMessages.UNEXPECTED).toBeTruthy();
            expect(ErrorMessages.NETWORK).toBeTruthy();
            expect(ErrorMessages.TIMEOUT).toBeTruthy();
            expect(ErrorMessages.AUTH_REQUIRED).toBeTruthy();
            expect(ErrorMessages.AUTH_EXPIRED).toBeTruthy();
            expect(ErrorMessages.FORBIDDEN).toBeTruthy();
        });

        it('should have operation-specific messages', () => {
            expect(ErrorMessages.FETCH_FAILED).toBeTruthy();
            expect(ErrorMessages.CREATE_FAILED).toBeTruthy();
            expect(ErrorMessages.UPDATE_FAILED).toBeTruthy();
            expect(ErrorMessages.DELETE_FAILED).toBeTruthy();
        });

        it('should have domain-specific messages', () => {
            expect(ErrorMessages.ORDER_FETCH_FAILED).toBeTruthy();
            expect(ErrorMessages.ORDER_CREATE_FAILED).toBeTruthy();
            expect(ErrorMessages.PAYMENT_FAILED).toBeTruthy();
            expect(ErrorMessages.WALLET_FETCH_FAILED).toBeTruthy();
            expect(ErrorMessages.PLAYER_FETCH_FAILED).toBeTruthy();
            expect(ErrorMessages.DISPUTE_CREATE_FAILED).toBeTruthy();
        });

        it('should have const assertion for type safety', () => {
            // This test verifies the `as const` assertion is working
            const messages = ErrorMessages;
            expect(typeof messages.UNEXPECTED).toBe('string');
        });
    });

    describe('Complex Error Scenarios', () => {
        it('should handle error with both message and error fields', () => {
            const apiError = {
                message: 'Primary message',
                error: 'Secondary message',
                code: 'ERR_001',
            };
            expect(getErrorMessage(apiError)).toBe('Primary message');
        });

        it('should handle AxiosError without response', () => {
            const axiosError = new AxiosError('Network Error');
            expect(getErrorMessage(axiosError)).toBe('Network Error');
        });

        it('should handle error with details field', () => {
            const error = {
                message: 'Validation failed',
                details: { field: 'email', reason: 'invalid format' },
            };
            expect(getErrorMessage(error)).toBe('Validation failed');
        });

        it('should handle empty string message', () => {
            const apiError = { message: '' };
            expect(getErrorMessage(apiError, 'Fallback')).toBe('Fallback');
        });

        it('should handle null message field', () => {
            const apiError = { message: null as any };
            expect(getErrorMessage(apiError, 'Fallback')).toBe('Fallback');
        });
    });

    describe('Error Type Edge Cases', () => {
        it('should handle number as error', () => {
            expect(getErrorMessage(404, 'Not found')).toBe('Not found');
        });

        it('should handle boolean as error', () => {
            expect(getErrorMessage(false, 'Error')).toBe('Error');
        });

        it('should handle array as error', () => {
            expect(getErrorMessage(['error1', 'error2'], 'Error')).toBe('Error');
        });

        it('should handle function as error', () => {
            expect(getErrorMessage(() => {}, 'Error')).toBe('Error');
        });

        it('should handle Date object as error', () => {
            expect(getErrorMessage(new Date(), 'Error')).toBe('Error');
        });
    });

    describe('Real-world Error Scenarios', () => {
        it('should handle login failure error', () => {
            const error = {
                message: '用户名或密码错误',
                code: '40001',
                status: 401,
            };
            const message = getErrorMessage(error);
            expect(message).toBe('用户名或密码错误');
        });

        it('should handle token expiry error', () => {
            const error = {
                message: 'Token已过期，请重新登录',
                code: '40004',
                status: 401,
            };
            const message = getErrorMessage(error);
            expect(message).toBe('Token已过期，请重新登录');
        });

        it('should handle insufficient balance error', () => {
            const error = {
                message: '余额不足，请先充值',
                code: '40301',
                status: 400,
            };
            const message = getErrorMessage(error);
            expect(message).toBe('余额不足，请先充值');
        });

        it('should handle network timeout', () => {
            const error = new Error('timeout of 5000ms exceeded');
            expect(isNetworkError(error)).toBe(true);
        });

        it('should handle connection refused', () => {
            const error = new Error('network connection refused');
            expect(isNetworkError(error)).toBe(true);
        });
    });

    describe('Multiple Error Type Identification', () => {
        it('should correctly classify error by both status and code', () => {
            const error = {
                status: 400,
                code: '40201',
                message: '手机号格式不正确',
            };

            expect(isValidationError(error)).toBe(true);
            expect(isValidationErrorByCode(40201)).toBe(true);
        });

        it('should handle mixed error classification', () => {
            expect(isAuthErrorByCode(40001)).toBe(true);
            expect(isAuthorizationErrorByCode(40001)).toBe(false);
            expect(isValidationErrorByCode(40001)).toBe(false);
            expect(isBusinessErrorByCode(40001)).toBe(false);
        });
    });
});
