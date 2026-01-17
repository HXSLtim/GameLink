/**
 * Standardized error handling utilities for stores
 *
 * This module provides consistent error handling patterns across all Zustand stores.
 */

import axios, { AxiosError } from 'axios';

// ============ Type Guards ============

/**
 * Type guard to check if error is an AxiosError
 */
export function isAxiosError<T = unknown>(err: unknown): err is AxiosError<T> {
    return axios.isAxiosError(err);
}

// ============ Error Types ============

export interface AppError {
    message: string;
    code?: string;
    status?: number;
    details?: Record<string, unknown>;
}

export interface ApiErrorResponse {
    message?: string;
    error?: string;
    code?: string;
    status?: number;
    details?: Record<string, unknown>;
}

// ============ Error Extraction ============

/**
 * Extracts a user-friendly error message from any error type.
 * Handles Error objects, AxiosError, API error responses, and unknown error types.
 *
 * @param err - The error to extract message from
 * @param fallback - Fallback message if extraction fails
 * @returns A user-friendly error message
 */
export function getErrorMessage(err: unknown, fallback = 'An unexpected error occurred'): string {
    // Handle AxiosError specifically
    if (isAxiosError<ApiErrorResponse>(err)) {
        const responseData = err.response?.data;
        if (responseData?.message && typeof responseData.message === 'string') {
            return responseData.message;
        }
        if (responseData?.error && typeof responseData.error === 'string') {
            return responseData.error;
        }
        // Fall back to axios error message
        if (err.message) {
            return err.message;
        }
    }

    // Standard Error object
    if (err instanceof Error) {
        return err.message || fallback;
    }

    // API error response object (non-axios)
    if (err && typeof err === 'object') {
        const apiError = err as ApiErrorResponse;

        // Try common error message fields
        if (apiError.message && typeof apiError.message === 'string') {
            return apiError.message;
        }
        if (apiError.error && typeof apiError.error === 'string') {
            return apiError.error;
        }
    }

    // String error
    if (typeof err === 'string') {
        return err;
    }

    return fallback;
}

/**
 * Extracts error code from API error response
 */
export function getErrorCode(err: unknown): string | undefined {
    if (err && typeof err === 'object') {
        const apiError = err as ApiErrorResponse;
        return apiError.code;
    }
    return undefined;
}

/**
 * Extracts HTTP status from API error response
 */
export function getErrorStatus(err: unknown): number | undefined {
    if (err && typeof err === 'object') {
        const apiError = err as ApiErrorResponse;
        return apiError.status;
    }
    return undefined;
}

// ============ Error Logging ============

/**
 * Logs error with consistent format for debugging.
 * Only logs in development mode.
 *
 * @param context - Context string (e.g., 'fetchOrders', 'createPayment')
 * @param err - The error to log
 */
export function logError(context: string, err: unknown): void {
    if (import.meta.env.DEV) {
        console.error(`[${context}]`, err);
    }
}

/**
 * Logs warning with consistent format.
 * Only logs in development mode.
 */
export function logWarn(context: string, message: string): void {
    if (import.meta.env.DEV) {
        console.warn(`[${context}]`, message);
    }
}

// ============ Error Checking ============

/**
 * Checks if error is a network error (offline, timeout, etc.)
 */
export function isNetworkError(err: unknown): boolean {
    if (err instanceof Error) {
        const message = err.message.toLowerCase();
        return (
            message.includes('network') ||
            message.includes('timeout') ||
            message.includes('offline') ||
            message.includes('failed to fetch')
        );
    }
    return false;
}

/**
 * Checks if error is an authentication error (401)
 */
export function isAuthError(err: unknown): boolean {
    const status = getErrorStatus(err);
    return status === 401;
}

/**
 * Checks if error is a forbidden error (403)
 */
export function isForbiddenError(err: unknown): boolean {
    const status = getErrorStatus(err);
    return status === 403;
}

/**
 * Checks if error is a not found error (404)
 */
export function isNotFoundError(err: unknown): boolean {
    const status = getErrorStatus(err);
    return status === 404;
}

/**
 * Checks if error is a validation error (400)
 */
export function isValidationError(err: unknown): boolean {
    const status = getErrorStatus(err);
    return status === 400;
}

// ============ Store Error Handler ============

/**
 * Standard error handler for store actions.
 * Returns an object with error message for setting in store state.
 *
 * @param context - Context string for logging
 * @param err - The error
 * @param fallback - Fallback message
 * @param shouldLog - Whether to log the error (default: true)
 * @returns Object with error message
 *
 * @example
 * ```ts
 * catch (err) {
 *   const { error } = handleStoreError('fetchOrders', err, 'Failed to fetch orders');
 *   set({ loading: false, error });
 * }
 * ```
 */
export function handleStoreError(
    context: string,
    err: unknown,
    fallback: string,
    shouldLog = true
): { error: string } {
    if (shouldLog) {
        logError(context, err);
    }

    return {
        error: getErrorMessage(err, fallback)
    };
}

// ============ Error Code Mapping (Backend Error Codes) ============

/**
 * Map backend error codes to user-friendly Chinese messages
 * Aligned with backend error code ranges
 */
export const ERROR_CODE_MESSAGES: Record<number, string> = {
    // Authentication errors (40000-40099)
    40001: '用户名或密码错误',
    40002: '账号已被禁用',
    40003: '账号已被封禁',
    40004: 'Token已过期，请重新登录',
    40005: '无效的Token',
    40006: '验证码错误',
    40007: '验证码已过期',
    40008: '手机号已被注册',
    40009: '邮箱已被注册',

    // Authorization errors (40100-40199)
    40101: '您没有权限执行此操作',
    40102: '需要管理员权限',
    40103: '需要陪玩师权限',
    40104: '账号未认证',

    // Validation errors (40200-40299)
    40201: '手机号格式不正确',
    40202: '邮箱格式不正确',
    40203: '密码长度不符合要求',
    40204: '密码必须包含大小写字母和数字',
    40205: '必填字段不能为空',
    40206: '参数格式不正确',
    40207: '文件大小超出限制',
    40208: '文件类型不支持',

    // Business logic errors (40300-40399)
    40301: '余额不足，请先充值',
    40302: '订单不存在',
    40303: '订单状态不允许此操作',
    40304: '陪玩师不可用',
    40305: '服务项目不存在',
    40306: '优惠券不可用',
    40307: '优惠券已使用',
    40308: '优惠券已过期',
    40309: '提现金额不足',
    40310: '提现金额超出限制',
    40311: '已达到每日提现次数限制',
    40312: '该用户已被拉黑',
    40313: '不能对自己进行此操作',
    40314: '聊天消息发送失败',
    40315: '纠纷已存在',
    40316: '纠纷处理超时',

    // Resource not found (40400-40499)
    40401: '用户不存在',
    40402: '陪玩师不存在',
    40403: '订单不存在',
    40404: '资源不存在',

    // Server errors (50000-50099)
    50001: '服务器繁忙，请稍后重试',
    50002: '数据库连接失败',
    50003: '缓存服务异常',
    50004: '第三方服务异常',
    50005: '文件上传失败',
};

/**
 * Get user-friendly message from error code
 */
export function getMessageByCode(code: number): string {
    return ERROR_CODE_MESSAGES[code] || '操作失败，请稍后重试';
}

/**
 * Check if error code is in specific range
 */
export function isErrorCodeInRange(code: number, min: number, max: number): boolean {
    return code >= min && code < max;
}

/**
 * Check if error is authentication error by code
 */
export function isAuthErrorByCode(code: number): boolean {
    return isErrorCodeInRange(code, 40000, 40100);
}

/**
 * Check if error is authorization error by code
 */
export function isAuthorizationErrorByCode(code: number): boolean {
    return isErrorCodeInRange(code, 40100, 40200);
}

/**
 * Check if error is validation error by code
 */
export function isValidationErrorByCode(code: number): boolean {
    return isErrorCodeInRange(code, 40200, 40300);
}

/**
 * Check if error is business logic error by code
 */
export function isBusinessErrorByCode(code: number): boolean {
    return isErrorCodeInRange(code, 40300, 40400);
}

// ============ Error Messages (i18n-ready) ============

export const ErrorMessages = {
    // Generic
    UNEXPECTED: 'An unexpected error occurred',
    NETWORK: 'Network error. Please check your connection.',
    TIMEOUT: 'Request timed out. Please try again.',

    // Auth
    AUTH_REQUIRED: 'Please log in to continue',
    AUTH_EXPIRED: 'Your session has expired. Please log in again.',
    FORBIDDEN: 'You do not have permission to perform this action',

    // Data
    NOT_FOUND: 'The requested resource was not found',
    VALIDATION: 'Please check your input and try again',

    // Operations
    FETCH_FAILED: 'Failed to load data',
    CREATE_FAILED: 'Failed to create',
    UPDATE_FAILED: 'Failed to update',
    DELETE_FAILED: 'Failed to delete',

    // Specific
    ORDER_FETCH_FAILED: 'Failed to fetch orders',
    ORDER_CREATE_FAILED: 'Failed to create order',
    PAYMENT_FAILED: 'Payment failed',
    WALLET_FETCH_FAILED: 'Failed to fetch wallet',
    PLAYER_FETCH_FAILED: 'Failed to fetch players',
    DISPUTE_CREATE_FAILED: 'Failed to create dispute',
} as const;
