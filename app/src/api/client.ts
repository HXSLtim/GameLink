import Taro from '@tarojs/taro';
import { encryptRequest, decryptResponse, shouldEncrypt } from '../utils/crypto';

/**
 * API 响应结构
 */
export interface ApiResponse<T = unknown> {
  success: boolean;
  code: number;
  message: string;
  data: T;
}

/**
 * 请求配置
 */
interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  data?: unknown;
  header?: Record<string, string>;
  timeout?: number;
  showLoading?: boolean;
  loadingText?: string;
}

// API 基础 URL
const BASE_URL = process.env.TARO_APP_API_BASE_URL || 'http://localhost:8080/api/v1';

// Token 存储 key
const TOKEN_KEY = 'token';
const REFRESH_TOKEN_KEY = 'refresh_token';

/**
 * 获取存储的 Token
 */
export function getToken(): string {
  return Taro.getStorageSync(TOKEN_KEY) || '';
}

/**
 * 设置 Token
 */
export function setToken(token: string): void {
  Taro.setStorageSync(TOKEN_KEY, token);
}

/**
 * 移除 Token
 */
export function removeToken(): void {
  Taro.removeStorageSync(TOKEN_KEY);
  Taro.removeStorageSync(REFRESH_TOKEN_KEY);
}

/**
 * Token 刷新状态
 */
let isRefreshing = false;
let refreshSubscribers: Array<(token: string) => void> = [];

/**
 * 订阅 Token 刷新
 */
function subscribeTokenRefresh(callback: (token: string) => void): void {
  refreshSubscribers.push(callback);
}

/**
 * 通知所有订阅者
 */
function onTokenRefreshed(token: string): void {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

/**
 * 刷新 Token
 */
async function refreshToken(): Promise<string | null> {
  try {
    const response = await Taro.request<ApiResponse<{ token: string }>>({
      url: `${BASE_URL}/auth/refresh`,
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getToken()}`,
      },
    });

    if (response.data.success && response.data.data?.token) {
      const newToken = response.data.data.token;
      setToken(newToken);
      return newToken;
    }
    return null;
  } catch {
    return null;
  }
}

/**
 * 处理 401 错误
 */
async function handle401Error<T>(config: RequestConfig): Promise<ApiResponse<T>> {
  if (isRefreshing) {
    // 等待 Token 刷新完成
    return new Promise((resolve) => {
      subscribeTokenRefresh((token) => {
        config.header = {
          ...config.header,
          Authorization: `Bearer ${token}`,
        };
        resolve(request<T>(config));
      });
    });
  }

  isRefreshing = true;

  try {
    const newToken = await refreshToken();
    if (newToken) {
      onTokenRefreshed(newToken);
      config.header = {
        ...config.header,
        Authorization: `Bearer ${newToken}`,
      };
      return request<T>(config);
    }

    // 刷新失败，跳转登录
    removeToken();
    Taro.navigateTo({ url: '/pages/login/index' });
    return {
      success: false,
      code: 401,
      message: '登录已过期，请重新登录',
      data: null as T,
    };
  } finally {
    isRefreshing = false;
  }
}

/**
 * 统一请求方法
 */
export async function request<T = unknown>(config: RequestConfig): Promise<ApiResponse<T>> {
  const { url, method = 'GET', data, header = {}, timeout = 10000, showLoading = false, loadingText = '加载中...' } = config;

  // 显示 loading
  if (showLoading) {
    Taro.showLoading({ title: loadingText, mask: true });
  }

  try {
    // 构建请求头
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...header,
    };

    // 添加 Token
    const token = getToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    // 处理请求数据（加密）
    let requestData = data;
    const fullUrl = url.startsWith('http') ? url : `${BASE_URL}${url}`;

    if (data && shouldEncrypt(method, fullUrl)) {
      requestData = encryptRequest(data);
    }

    // 发起请求
    const response = await Taro.request<ApiResponse<T>>({
      url: fullUrl,
      method,
      data: requestData,
      header: headers,
      timeout,
    });

    // 隐藏 loading
    if (showLoading) {
      Taro.hideLoading();
    }

    // 处理 HTTP 状态码
    if (response.statusCode === 401) {
      return handle401Error<T>(config);
    }

    if (response.statusCode >= 400) {
      return {
        success: false,
        code: response.statusCode,
        message: response.data?.message || '请求失败',
        data: null as T,
      };
    }

    // 解密响应数据
    const responseData = decryptResponse<ApiResponse<T>>(response.data);

    return responseData;
  } catch (error) {
    // 隐藏 loading
    if (showLoading) {
      Taro.hideLoading();
    }

    console.error('Request error:', error);

    // 网络错误处理
    const errorMessage = error instanceof Error ? error.message : '网络请求失败';

    return {
      success: false,
      code: -1,
      message: errorMessage,
      data: null as T,
    };
  }
}

/**
 * GET 请求
 */
export function get<T = unknown>(url: string, options?: Omit<RequestConfig, 'url' | 'method'>): Promise<ApiResponse<T>> {
  return request<T>({ ...options, url, method: 'GET' });
}

/**
 * POST 请求
 */
export function post<T = unknown>(url: string, data?: unknown, options?: Omit<RequestConfig, 'url' | 'method' | 'data'>): Promise<ApiResponse<T>> {
  return request<T>({ ...options, url, method: 'POST', data });
}

/**
 * PUT 请求
 */
export function put<T = unknown>(url: string, data?: unknown, options?: Omit<RequestConfig, 'url' | 'method' | 'data'>): Promise<ApiResponse<T>> {
  return request<T>({ ...options, url, method: 'PUT', data });
}

/**
 * DELETE 请求
 */
export function del<T = unknown>(url: string, options?: Omit<RequestConfig, 'url' | 'method'>): Promise<ApiResponse<T>> {
  return request<T>({ ...options, url, method: 'DELETE' });
}

/**
 * PATCH 请求
 */
export function patch<T = unknown>(url: string, data?: unknown, options?: Omit<RequestConfig, 'url' | 'method' | 'data'>): Promise<ApiResponse<T>> {
  return request<T>({ ...options, url, method: 'PATCH', data });
}

export default {
  request,
  get,
  post,
  put,
  del,
  patch,
  getToken,
  setToken,
  removeToken,
};
