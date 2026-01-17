/**
 * Enhanced HTTP client for GameLink
 * Features:
 * - Proactive JWT token refresh (5-minute buffer)
 * - Request encryption (AES-256-CBC)
 * - Unified error handling
 * - Request queue during token refresh with size limits and timeout
 * - Performance monitoring
 * - Auto-unwrap API responses
 */

import axios, {
    type AxiosInstance,
    type AxiosRequestConfig,
    type AxiosResponse,
    type AxiosError,
    type InternalAxiosRequestConfig
} from 'axios';
import { useAuthStore } from '@/stores';
import { encryptRequest, shouldEncrypt } from './crypto';
import {
    TOKEN_REFRESH_BUFFER_SECONDS,
    MAX_QUEUE_SIZE,
    QUEUE_TIMEOUT_MS,
    QUEUE_CLEANUP_INTERVAL_MS,
    HTTP_TIMEOUT_MS,
} from './constants';
import { performanceMonitor } from './monitor';

/**
 * Standard API response wrapper
 */
interface ApiResponse<T> {
    success: boolean;
    code: number;
    message: string;
    data: T;
}

/**
 * Request metadata for performance tracking
 */
interface RequestMetadata {
    startTime: number;
}

/**
 * Extended Axios config with metadata
 */
interface AxiosConfigWithMetadata extends AxiosRequestConfig {
    metadata?: RequestMetadata;
}

/**
 * JWT payload structure
 */
interface JWTPayload {
    exp: number;  // Expiration timestamp (seconds)
    iat: number;  // Issued at timestamp (seconds)
    sub: number;  // User ID
}

/**
 * Parse JWT token to extract payload
 */
function parseJWT(token: string): JWTPayload | null {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) return null;

        const payload = JSON.parse(atob(parts[1]));
        return payload as JWTPayload;
    } catch {
        return null;
    }
}

/**
 * Check if token is expiring soon (within buffer seconds)
 * Default buffer: 300 seconds (5 minutes)
 */
function isTokenExpiringSoon(token: string, bufferSeconds = TOKEN_REFRESH_BUFFER_SECONDS): boolean {
    const payload = parseJWT(token);
    if (!payload?.exp) return true;

    const now = Math.floor(Date.now() / 1000);
    return payload.exp - now < bufferSeconds;
}

/**
 * Check if token is already expired
 */
function isTokenExpired(token: string): boolean {
    const payload = parseJWT(token);
    if (!payload?.exp) return true;

    const now = Math.floor(Date.now() / 1000);
    return now >= payload.exp;
}

// Export JWT utilities for external use
export { parseJWT, isTokenExpiringSoon, isTokenExpired };

// Get API base URL from environment
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

/**
 * Queued request with timestamp for timeout tracking
 */
interface QueuedRequest {
    resolve: (value?: unknown) => void;
    reject: (reason?: unknown) => void;
    timestamp: number;
}

// Request queue for token refresh with enhanced tracking
let isRefreshing = false;
let failedQueue: QueuedRequest[] = [];

/**
 * Process queued requests after token refresh
 *
 * @param error - Error if refresh failed, null if successful
 */
const processQueue = (error: Error | null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve();
        }
    });
    failedQueue = [];
};

/**
 * Add a request to the queue with size limit checking
 *
 * @returns Promise that resolves when token refresh completes
 * @throws Error if queue is full
 */
function enqueueRequest(): Promise<void> {
    // Check queue size limit
    if (failedQueue.length >= MAX_QUEUE_SIZE) {
        performanceMonitor.recordQueueRejection();
        return Promise.reject(new Error('Token refresh queue is full'));
    }

    return new Promise((resolve, reject) => {
        failedQueue.push({
            resolve,
            reject,
            timestamp: Date.now(),
        });
    });
}

/**
 * Clean up expired queue items
 * Called periodically to remove requests that have waited too long
 */
function cleanupExpiredQueueItems(): void {
    const now = Date.now();
    const initialLength = failedQueue.length;

    failedQueue = failedQueue.filter(item => {
        if (now - item.timestamp > QUEUE_TIMEOUT_MS) {
            item.reject(new Error('Request queue timeout'));
            performanceMonitor.recordQueueTimeout();
            return false;
        }
        return true;
    });

    // Log if items were removed
    if (import.meta.env.DEV && failedQueue.length < initialLength) {
        console.warn(
            `[HTTP] Cleaned up ${initialLength - failedQueue.length} expired queue items`
        );
    }
}

// Start periodic queue cleanup
if (typeof window !== 'undefined') {
    const cleanupInterval = setInterval(cleanupExpiredQueueItems, QUEUE_CLEANUP_INTERVAL_MS);

    // Clear interval on page unload to prevent memory leaks
    window.addEventListener('beforeunload', () => {
        clearInterval(cleanupInterval);
    });
}

/**
 * Enhanced HTTP client class
 */
export class HttpClient {
    private instance: AxiosInstance;

    constructor() {
        this.instance = axios.create({
            baseURL: API_BASE_URL,
            timeout: HTTP_TIMEOUT_MS,
            withCredentials: true,
            xsrfCookieName: 'csrf_token',
            xsrfHeaderName: 'X-CSRF-Token',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        this.setupInterceptors();
    }

    /**
     * Setup request and response interceptors
     */
    private setupInterceptors() {
        // Request Interceptor
        this.instance.interceptors.request.use(
            async (config: AxiosConfigWithMetadata) => {
                // Add timestamp for performance monitoring
                config.metadata = { startTime: Date.now() };

                const token = useAuthStore.getState().token;

                if (token) {
                    const url = config.url || '';
                    const isAuthRequest = url.includes('/auth/login') ||
                                        url.includes('/auth/register') ||
                                        url.includes('/auth/refresh');

                    // Proactive token refresh (prevent 401 storms)
                    if (!isAuthRequest && isTokenExpiringSoon(token)) {
                        if (isRefreshing) {
                            // Wait for ongoing refresh
                            try {
                                await enqueueRequest();
                                // Use new token after refresh
                                const newToken = useAuthStore.getState().token;
                                if (newToken) {
                                    config.headers.Authorization = `Bearer ${newToken}`;
                                }
                            } catch {
                                // Queue rejected or timed out, still try with old token
                                config.headers.Authorization = `Bearer ${token}`;
                            }
                        } else {
                            // Start refresh
                            isRefreshing = true;
                            try {
                                await useAuthStore.getState().refresh();
                                performanceMonitor.recordTokenRefresh(true);
                                processQueue(null);
                                const newToken = useAuthStore.getState().token;
                                if (newToken) {
                                    config.headers.Authorization = `Bearer ${newToken}`;
                                }
                            } catch (err) {
                                performanceMonitor.recordTokenRefresh(false);
                                processQueue(err as Error);
                                // Use old token as fallback
                                config.headers.Authorization = `Bearer ${token}`;
                            } finally {
                                isRefreshing = false;
                            }
                        }
                    } else {
                        config.headers.Authorization = `Bearer ${token}`;
                    }
                }

                // Encrypt request body if needed
                if (config.data && shouldEncrypt(config.method || 'GET', config.url || '')) {
                    config.data = encryptRequest(config.data);
                }

                return config;
            },
            (error) => Promise.reject(error)
        );

        // Response Interceptor
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => {
                // Record successful request performance
                const metadata = (response.config as AxiosConfigWithMetadata).metadata;
                const duration = metadata ? Date.now() - metadata.startTime : 0;
                performanceMonitor.recordRequest(duration, true);

                return response;
            },
            async (error: AxiosError) => {
                // Record failed request performance
                const metadata = error.config ? (error.config as AxiosConfigWithMetadata).metadata : undefined;
                const duration = metadata ? Date.now() - metadata.startTime : 0;
                performanceMonitor.recordRequest(duration, false);

                const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

                if (!error.response) {
                    // Network error
                    return Promise.reject(new Error('网络连接失败，请检查网络'));
                }

                const { status } = error.response;

                // Handle 401 - Token expired (reactive fallback)
                if (status === 401 && !originalRequest._retry) {
                    if (isRefreshing) {
                        // Wait for token refresh
                        try {
                            await enqueueRequest();
                            return this.instance(originalRequest);
                        } catch (err) {
                            // Queue rejected or timed out
                            return Promise.reject(err);
                        }
                    }

                    originalRequest._retry = true;
                    isRefreshing = true;

                    try {
                        await useAuthStore.getState().refresh();
                        performanceMonitor.recordTokenRefresh(true);
                        processQueue(null);
                        return this.instance(originalRequest);
                    } catch (refreshError) {
                        performanceMonitor.recordTokenRefresh(false);
                        processQueue(refreshError as Error);
                        await useAuthStore.getState().logout();
                        return Promise.reject(refreshError);
                    } finally {
                        isRefreshing = false;
                    }
                }

                // Handle 403 - Forbidden
                if (status === 403) {
                    console.warn('Access forbidden:', error.response.data);
                }

                return Promise.reject(error);
            }
        );
    }

    /**
     * Unwrap API response to extract data
     */
    private unwrap<T>(response: AxiosResponse<ApiResponse<T> | T>): T {
        const body = response.data as ApiResponse<T>;
        if (body && typeof body === 'object' && body.success === true && 'data' in body) {
            return body.data;
        }
        return response.data as T;
    }

    /**
     * GET request
     */
    public async get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.get<ApiResponse<T> | T>(url, config);
        return this.unwrap<T>(response);
    }

    /**
     * POST request
     */
    public async post<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.post<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    /**
     * PUT request
     */
    public async put<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.put<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    /**
     * PATCH request
     */
    public async patch<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.patch<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    /**
     * DELETE request
     */
    public async delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.delete<ApiResponse<T> | T>(url, config);
        return this.unwrap<T>(response);
    }
}

// Export singleton instance
export const http = new HttpClient();
