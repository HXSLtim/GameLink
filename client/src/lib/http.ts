import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type AxiosError } from 'axios';
import { useAuthStore } from '@/stores';


class HttpClient {
    private instance: AxiosInstance;

    constructor() {
        this.instance = axios.create({
            baseURL: '/api/v1',
            timeout: 10000,
            withCredentials: true, // Enable sending cookies (Refresh Token, CSRF Token)
            xsrfCookieName: 'csrf_token', // Standard CSRF cookie name
            xsrfHeaderName: 'X-CSRF-Token', // Standard CSRF header name
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
            (error) => {
                return Promise.reject(error);
            }
        );

        // Response Interceptor
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => {
                // You can unwrap response data here if your API always wraps in { data: ... }
                // For now, returning the full response or just response.data depending on convention
                return response;
            },
            async (error: AxiosError) => {
                if (error.response) {
                    const { status } = error.response;

                    if (status === 401) {
                        // Token expired or invalid
                        console.warn('Unauthorized access, logging out...');
                        const { logout } = useAuthStore.getState();
                        await logout();
                        // Optionally redirect to login, but store logout usually handles state clearing
                        // which might trigger a redirect via ProtectedRoute
                    }
                }
                return Promise.reject(error);
            }
        );
    }

    // Helper to unwrap response
    private unwrap<T>(response: AxiosResponse<any>): T {
        const body = response.data;
        // Check for standard API wrapper: { success, code, data }
        if (body && typeof body === 'object' && body.success === true && 'data' in body) {
            return body.data as T;
        }
        return body as T;
    }

    // Generic request methods
    public async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.get<T>(url, config);
        return this.unwrap<T>(response);
    }

    public async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.post<T>(url, data, config);
        return this.unwrap<T>(response);
    }

    public async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.put<T>(url, data, config);
        return this.unwrap<T>(response);
    }

    public async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
        const response = await this.instance.delete<T>(url, config);
        return this.unwrap<T>(response);
    }
}

export const http = new HttpClient();
