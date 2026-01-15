import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/stores';

// Standard API response wrapper
interface ApiResponse<T> {
    success: boolean;
    code: number;
    message: string;
    data: T;
}

// Get API base URL from environment variable
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

// Request queue for token refresh
let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value?: unknown) => void;
    reject: (reason?: unknown) => void;
}> = [];

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

export class HttpClient {
    private instance: AxiosInstance;

    constructor() {
        this.instance = axios.create({
            baseURL: API_BASE_URL,
            timeout: 15000,
            withCredentials: true,
            xsrfCookieName: 'csrf_token',
            xsrfHeaderName: 'X-CSRF-Token',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        this.setupInterceptors();
    }

    private setupInterceptors() {
        // Request Interceptor
        this.instance.interceptors.request.use(
            (config) => {
                const token = useAuthStore.getState().token;
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
                return config;
            },
            (error) => Promise.reject(error)
        );

        // Response Interceptor
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => response,
            async (error: AxiosError) => {
                const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

                if (!error.response) {
                    // Network error
                    return Promise.reject(new Error('网络连接失败，请检查网络'));
                }

                const { status } = error.response;

                // Handle 401 - Token expired
                if (status === 401 && !originalRequest._retry) {
                    if (isRefreshing) {
                        // Wait for token refresh
                        return new Promise((resolve, reject) => {
                            failedQueue.push({ resolve, reject });
                        }).then(() => {
                            return this.instance(originalRequest);
                        });
                    }

                    originalRequest._retry = true;
                    isRefreshing = true;

                    try {
                        await useAuthStore.getState().refresh();
                        processQueue(null);
                        return this.instance(originalRequest);
                    } catch (refreshError) {
                        processQueue(refreshError as Error);
                        await useAuthStore.getState().logout();
                        return Promise.reject(refreshError);
                    } finally {
                        isRefreshing = false;
                    }
                }

                // Handle other errors
                if (status === 403) {
                    console.warn('Access forbidden');
                }

                return Promise.reject(error);
            }
        );
    }

    // Helper to unwrap response
    private unwrap<T>(response: AxiosResponse<ApiResponse<T> | T>): T {
        const body = response.data as ApiResponse<T>;
        if (body && typeof body === 'object' && body.success === true && 'data' in body) {
            return body.data;
        }
        return response.data as T;
    }

    public async get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.get<ApiResponse<T> | T>(url, config);
        return this.unwrap<T>(response);
    }

    public async post<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.post<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    public async put<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.put<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    public async patch<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.patch<ApiResponse<T> | T>(url, data, config);
        return this.unwrap<T>(response);
    }

    public async delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.delete<ApiResponse<T> | T>(url, config);
        return this.unwrap<T>(response);
    }
}

export const http = new HttpClient();
