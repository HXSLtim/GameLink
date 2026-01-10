/**
 * API Client
 * 统一的 HTTP 请求客户端
 */
import Taro from '@tarojs/taro';
import type { ApiResponse } from '@/types/api';

interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: Record<string, unknown>;
  headers?: Record<string, string>;
  skipAuth?: boolean;
}

/**
 * API 错误类
 */
export class ApiError extends Error {
  code: number;
  traceId?: string;

  constructor(message: string, code: number, traceId?: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.traceId = traceId;
  }
}

class ApiClient {
  private baseURL: string;
  private timeout: number = 10000;
  private isRefreshing: boolean = false;
  private refreshPromise: Promise<boolean> | null = null;

  constructor() {
    this.baseURL = process.env.TARO_APP_API_BASE_URL || 'http://localhost:8080/api/v1';
  }

  /**
   * 获取存储的 Token
   */
  private getToken(): string | null {
    try {
      return Taro.getStorageSync('token') || null;
    } catch {
      return null;
    }
  }

  /**
   * 获取刷新 Token
   */
  private getRefreshToken(): string | null {
    try {
      return Taro.getStorageSync('refreshToken') || null;
    } catch {
      return null;
    }
  }

  /**
   * 保存 Token
   */
  private saveTokens(token: string, refreshToken?: string): void {
    Taro.setStorageSync('token', token);
    if (refreshToken) {
      Taro.setStorageSync('refreshToken', refreshToken);
    }
  }

  /**
   * 清除 Token
   */
  private clearTokens(): void {
    Taro.removeStorageSync('token');
    Taro.removeStorageSync('refreshToken');
    Taro.removeStorageSync('userInfo');
  }

  /**
   * 刷新 Token
   */
  async refreshToken(): Promise<boolean> {
    // 防止并发刷新
    if (this.isRefreshing && this.refreshPromise) {
      return this.refreshPromise;
    }

    const refreshTokenValue = this.getRefreshToken();
    if (!refreshTokenValue) {
      return false;
    }

    this.isRefreshing = true;
    this.refreshPromise = this.doRefreshToken(refreshTokenValue);

    try {
      return await this.refreshPromise;
    } finally {
      this.isRefreshing = false;
      this.refreshPromise = null;
    }
  }

  /**
   * 执行 Token 刷新
   */
  private async doRefreshToken(refreshTokenValue: string): Promise<boolean> {
    try {
      const response = await Taro.request({
        url: `${this.baseURL}/auth/refresh`,
        method: 'POST',
        data: { refreshToken: refreshTokenValue },
        header: { 'Content-Type': 'application/json' },
        timeout: this.timeout,
      });

      const result = response.data as ApiResponse<{ token: string; refreshToken?: string }>;

      if (result.success && result.data?.token) {
        this.saveTokens(result.data.token, result.data.refreshToken);
        return true;
      }
      return false;
    } catch {
      return false;
    }
  }

  /**
   * 跳转到登录页
   */
  private redirectToLogin(): void {
    this.clearTokens();
    Taro.redirectTo({ url: '/pages/login/index' });
  }

  /**
   * 显示错误提示
   */
  private showError(message: string): void {
    Taro.showToast({
      title: message,
      icon: 'none',
      duration: 2000,
    });
  }

  /**
   * 通用请求方法
   */
  async request<T>(config: RequestConfig): Promise<ApiResponse<T>> {
    const { url, method = 'GET', data, headers = {}, skipAuth = false } = config;

    // 添加认证 Token
    if (!skipAuth) {
      const token = this.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    try {
      const response = await Taro.request({
        url: `${this.baseURL}${url}`,
        method,
        data,
        header: {
          'Content-Type': 'application/json',
          ...headers,
        },
        timeout: this.timeout,
      });

      const result = response.data as ApiResponse<T>;

      // 处理 401 未授权
      if (result.code === 401 && !skipAuth) {
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // 重试原请求
          return this.request<T>(config);
        }
        // 刷新失败，跳转登录
        this.redirectToLogin();
        throw new ApiError('登录已过期，请重新登录', 401, result.traceId);
      }

      // 处理其他错误
      if (!result.success) {
        const errorMessage = result.message || '请求失败';
        // 显示错误提示（排除 401）
        if (result.code !== 401) {
          this.showError(errorMessage);
        }
        throw new ApiError(errorMessage, result.code, result.traceId);
      }

      return result;
    } catch (error) {
      // 已经是 ApiError，直接抛出
      if (error instanceof ApiError) {
        throw error;
      }

      // 网络错误
      const networkError = error as { errMsg?: string };
      const message = networkError.errMsg?.includes('timeout')
        ? '请求超时，请稍后重试'
        : '网络错误，请检查网络连接';

      this.showError(message);
      throw new ApiError(message, -1);
    }
  }

  /**
   * GET 请求
   */
  get<T>(url: string, params?: Record<string, unknown>, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'GET',
      data: params,
      ...config,
    });
  }

  /**
   * POST 请求
   */
  post<T>(url: string, data?: Record<string, unknown>, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'POST',
      data,
      ...config,
    });
  }

  /**
   * PUT 请求
   */
  put<T>(url: string, data?: Record<string, unknown>, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'PUT',
      data,
      ...config,
    });
  }

  /**
   * DELETE 请求
   */
  delete<T>(url: string, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'DELETE',
      ...config,
    });
  }
}

export default new ApiClient();
