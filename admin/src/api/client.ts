import axios from 'axios';
import { encryptRequest, shouldEncrypt } from '../utils/crypto';

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// ============ Token 过期检测工具 ============

interface JWTPayload {
    exp: number;  // 过期时间戳（秒）
    iat: number;  // 签发时间戳（秒）
    sub: number;  // 用户ID
}

/**
 * 解析 JWT Token 获取 payload
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
 * 检查 Token 是否即将过期（提前 5 分钟刷新）
 */
function isTokenExpiringSoon(token: string, bufferSeconds = 300): boolean {
    const payload = parseJWT(token);
    if (!payload?.exp) return true; // 无法解析视为过期
    
    const now = Math.floor(Date.now() / 1000);
    return payload.exp - now < bufferSeconds;
}

/**
 * 检查 Token 是否已过期
 */
function isTokenExpired(token: string): boolean {
    const payload = parseJWT(token);
    if (!payload?.exp) return true;
    
    const now = Math.floor(Date.now() / 1000);
    return now >= payload.exp;
}

// 导出工具函数供外部使用
export { parseJWT, isTokenExpiringSoon, isTokenExpired };

// ============ Token 刷新队列管理 ============

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

/**
 * 执行 Token 刷新
 */
async function refreshToken(): Promise<string> {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
        throw new Error('No token to refresh');
    }

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
                Authorization: `Bearer ${currentToken}`
            }
        }
    );

    const tokenData = response.data?.data;
    if (!tokenData?.token) {
        throw new Error('Invalid refresh token response');
    }

    const { token } = tokenData;
    
    // 更新存储
    localStorage.setItem('token', token);
    apiClient.defaults.headers.common['Authorization'] = 'Bearer ' + token;
    
    return token;
}

/**
 * 清除认证状态并重定向到登录页
 */
function clearAuthAndRedirect() {
    localStorage.removeItem('token');
    localStorage.removeItem('user_role');
    localStorage.removeItem('user_info');
    localStorage.removeItem('auth-storage');
    
    if (!window.location.pathname.includes('/login')) {
        window.location.href = '/admin/login';
    }
}

// ============ Request Interceptor ============

apiClient.interceptors.request.use(
    async (config) => {
        const token = localStorage.getItem('token');
        
        if (token) {
            // 检查是否是登录/刷新请求（这些不需要检查过期）
            const url = config.url || '';
            const isAuthRequest = url.includes('/auth/login') || url.includes('/auth/refresh');
            
            if (!isAuthRequest) {
                // 主动检测：Token 即将过期时提前刷新
                if (isTokenExpiringSoon(token)) {
                    // 如果已经在刷新中，等待刷新完成
                    if (isRefreshing) {
                        return new Promise((resolve, reject) => {
                            failedQueue.push({
                                resolve: (newToken) => {
                                    config.headers.Authorization = `Bearer ${newToken}`;
                                    resolve(config);
                                },
                                reject: (err) => reject(err)
                            });
                        });
                    }
                    
                    // 开始刷新
                    isRefreshing = true;
                    try {
                        const newToken = await refreshToken();
                        processQueue(null, newToken);
                        config.headers.Authorization = `Bearer ${newToken}`;
                    } catch (err) {
                        processQueue(err, null);
                        // Token 刷新失败，但不阻止请求（让 401 响应处理）
                        config.headers.Authorization = `Bearer ${token}`;
                    } finally {
                        isRefreshing = false;
                    }
                } else {
                    config.headers.Authorization = `Bearer ${token}`;
                }
            } else {
                config.headers.Authorization = `Bearer ${token}`;
            }
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

// ============ Response Interceptor ============

apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        
        // 检查是否在登录页面
        const pathname = window.location.pathname;
        const isOnLoginPage = pathname.includes('/login');
        
        // 检查是否是登录/刷新请求
        const requestUrl = originalRequest?.url || '';
        const isAuthRequest = requestUrl.includes('/auth/login') || requestUrl.includes('/auth/refresh');
        
        // 登录页面或认证请求，直接返回错误
        if (isOnLoginPage || isAuthRequest) {
            return Promise.reject(error);
        }

        // 401 错误且未重试过
        if (error.response?.status === 401 && !originalRequest._retry) {
            // 如果正在刷新，加入队列等待
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({
                        resolve: (token) => {
                            originalRequest.headers['Authorization'] = 'Bearer ' + token;
                            resolve(apiClient(originalRequest));
                        },
                        reject
                    });
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const newToken = await refreshToken();
                processQueue(null, newToken);
                originalRequest.headers['Authorization'] = 'Bearer ' + newToken;
                return apiClient(originalRequest);
            } catch (err) {
                processQueue(err, null);
                clearAuthAndRedirect();
                return Promise.reject(err);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

export default apiClient;
