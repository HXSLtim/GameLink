/**
 * Standardized error handling utilities for stores
 *
 * This module provides consistent error handling patterns across all Zustand stores.
 */

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
 * Handles Error objects, API error responses, and unknown error types.
 *
 * @param err - The error to extract message from
 * @param fallback - Fallback message if extraction fails
 * @returns A user-friendly error message
 */
export function getErrorMessage(err: unknown, fallback = 'An unexpected error occurred'): string {
    // Standard Error object
    if (err instanceof Error) {
        return err.message || fallback;
    }

    // API error response object
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
