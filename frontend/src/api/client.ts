import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { storage } from '../utils/storage';
import { STORAGE_KEYS } from '../config';
import { cryptoMiddleware } from '../middleware/crypto';

/**
 * 后端统一响应格式
 */
interface ApiResponse<T = unknown> {
  success: boolean;
  code: number;
  message: string;
  data: T;
}

/**
 * 创建 axios 实例
 */
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * 请求拦截器 - 添加 Token 和加密
 */
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 1. 添加 Token
    const token = storage.getItem<string>(STORAGE_KEYS.token);
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 2. 加密请求数据
    return cryptoMiddleware.requestInterceptor(config);
  },
  (error: AxiosError) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  },
);

/**
 * 响应拦截器 - 解密和统一处理响应格式
 */
apiClient.interceptors.response.use(
  (response): any => {
    // 1. 先解密响应数据
    const decryptedResponse = cryptoMiddleware.responseInterceptor(response);

    // 2. 处理后端统一格式 { success, code, message, data, pagination }
    const apiResponse = decryptedResponse.data;

    if (apiResponse.success) {
      // 如果有 pagination 字段，说明是列表响应，需要转换格式
      if (apiResponse.pagination) {
        return {
          list: apiResponse.data,
          total: apiResponse.pagination.total,
          page: apiResponse.pagination.page,
          page_size: apiResponse.pagination.page_size,
          total_pages: apiResponse.pagination.total_pages,
          has_next: apiResponse.pagination.has_next,
          has_prev: apiResponse.pagination.has_prev,
        };
      }
      // 普通响应：直接返回 data 部分
      return apiResponse.data;
    } else {
      // 失败：抛出错误
      const error = new Error(apiResponse.message || '请求失败');
      return Promise.reject(error);
    }
  },
  (error: AxiosError) => {
    // HTTP 错误处理
    if (error.response) {
      const status = error.response.status;
      const apiResponse = error.response.data as ApiResponse;

      switch (status) {
        case 401:
          // Token 过期或无效
          console.error('认证失败，请重新登录');
          storage.removeItem(STORAGE_KEYS.token);
          storage.removeItem(STORAGE_KEYS.user);
          // 跳转到登录页（避免循环跳转）
          if (window.location.pathname !== '/login') {
            window.location.href = '/login';
          }
          break;

        case 403:
          console.error('权限不足');
          break;

        case 404:
          console.error('资源不存在');
          break;

        case 500:
          console.error('服务器错误');
          break;

        default:
          console.error('请求错误:', apiResponse?.message || error.message);
      }

      // 返回错误消息
      return Promise.reject(new Error(apiResponse?.message || `请求失败 (${status})`));
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('网络错误，请检查网络连接');
      return Promise.reject(new Error('网络错误，请检查网络连接'));
    } else {
      // 其他错误
      console.error('请求配置错误:', error.message);
      return Promise.reject(error);
    }
  },
);

/**
 * 导出配置好的 axios 实例
 */
export default apiClient;
