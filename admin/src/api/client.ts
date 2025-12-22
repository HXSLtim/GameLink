import axios from 'axios';
import { encryptRequest, shouldEncrypt } from '../utils/crypto';

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor
apiClient.interceptors.request.use(
    (config) => {
        // 添加 JWT Token
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        
        // 加密请求体（如果需要）
        if (config.data && shouldEncrypt(config.method || 'GET', config.url || '')) {
            config.data = encryptRequest(config.data);
        }
        
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response interceptor
interface FailedRequest {
    resolve: (token: string | null) => void;
    reject: (error: unknown) => void;
}

let isRefreshing = false;
let failedQueue: FailedRequest[] = [];

const processQueue = (error: unknown, token: string | null = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });

    failedQueue = [];
};

apiClient.interceptors.response.use(
    (response) => {
        // Return full response to maintain AxiosResponse<T> type contract
        return response;
    },
    async (error) => {
        const originalRequest = error.config;
        
        // 检查是否在登录页面 - 使用多种方式确保检测准确
        const pathname = window.location.pathname;
        const isOnLoginPage = pathname.includes('/login') || 
                              pathname === '/admin/login' || 
                              pathname.endsWith('/login');
        
        // 检查是否是登录请求
        const requestUrl = originalRequest?.url || '';
        const isLoginRequest = requestUrl.includes('/auth/login') || 
                               requestUrl.includes('login') ||
                               requestUrl === '/auth/login';
        
        // 在登录页面或登录请求，所有错误都直接返回，不尝试刷新 token
        if (isOnLoginPage || isLoginRequest) {
            return Promise.reject(error);
        }

        if (error.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                return new Promise(function (resolve, reject) {
                    failedQueue.push({ resolve, reject });
                }).then((token) => {
                    originalRequest.headers['Authorization'] = 'Bearer ' + token;
                    return apiClient(originalRequest);
                }).catch(err => {
                    return Promise.reject(err);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                // Import authApi dynamically to avoid circular dependency if possible,
                // but here we might need to use a direct axios call or ensure authApi is safe.
                // Using a direct call to avoid circular dependency issues with auth.ts importing client.ts
                const response = await axios.post<{
                    success: boolean;
                    code: number;
                    message: string;
                    data: { token: string }
                }>(
                    (import.meta.env.VITE_API_BASE_URL || '/api/v1') + '/auth/refresh',
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('token')}`
                        }
                    }
                );

                const { token } = response.data.data;

                localStorage.setItem('token', token);
                apiClient.defaults.headers.common['Authorization'] = 'Bearer ' + token;
                processQueue(null, token);

                originalRequest.headers['Authorization'] = 'Bearer ' + token;
                return apiClient(originalRequest);
            } catch (err) {
                processQueue(err, null);
                localStorage.removeItem('token');
                localStorage.removeItem('user_role');
                localStorage.removeItem('user_info');
                // 如果当前已经在登录页面，不需要重定向
                if (!window.location.pathname.includes('/login')) {
                    window.location.href = '/admin/login';
                }
                return Promise.reject(err);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

export default apiClient;
