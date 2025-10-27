import { API_BASE, STORAGE_KEYS } from '../config';
import type { ApiResponse } from '../types/api';
import { ApiError } from '../types/api';
import { storage } from '../utils/storage';
import { retryAsync } from './retry';

function getToken(): string | null {
  return storage.getItem<string>(STORAGE_KEYS.token);
}

function buildQuery(params?: Record<string, unknown>) {
  if (!params) return '';
  const usp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v === undefined || v === null) return;
    usp.set(k, String(v));
  });
  const s = usp.toString();
  return s ? `?${s}` : '';
}

/**
 * 执行 HTTP 请求（带重试机制）
 * @param path API 路径
 * @param options 请求选项
 * @returns Promise<T>
 */
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };
  if (token) headers.Authorization = `Bearer ${token}`;

  // 使用重试机制执行请求
  const res = await retryAsync(
    () =>
      fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
        credentials: 'include',
      }),
    {
      maxRetries: 2, // 重试 2 次（总共 3 次请求）
      initialDelay: 500,
      shouldRetry: (error: unknown, _attempt: number) => {
        // 仅对网络错误和 5xx 错误重试
        if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
          return true; // 网络错误
        }
        if (error instanceof Response) {
          const status = error.status;
          // 服务器错误（5xx）但不包括特定的认证/授权错误
          return status >= 500 && status < 600;
        }
        return false;
      },
      onRetry: (error: unknown, attempt: number, delay: number) => {
        console.warn(`[HTTP] Retrying request to ${path} (attempt ${attempt}, delay ${delay}ms)`);
      },
    },
  );

  let payload: ApiResponse<T> | null = null;
  try {
    payload = (await res.json()) as ApiResponse<T>;
  } catch {
    // not json
  }

  if (!payload) {
    throw new ApiError(res.status, `HTTP ${res.status}`);
  }
  if (!payload.success) {
    throw new ApiError(payload.code, payload.message || 'Unknown error');
  }
  return payload.data;
}

export const http = {
  get<T>(path: string, params?: Record<string, unknown>) {
    return request<T>(`${path}${buildQuery(params)}`, { method: 'GET' });
  },
  post<T>(path: string, body?: unknown) {
    return request<T>(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined });
  },
  put<T>(path: string, body?: unknown) {
    return request<T>(path, { method: 'PUT', body: body ? JSON.stringify(body) : undefined });
  },
  patch<T>(path: string, body?: unknown) {
    return request<T>(path, { method: 'PATCH', body: body ? JSON.stringify(body) : undefined });
  },
  delete<T>(path: string) {
    return request<T>(path, { method: 'DELETE' });
  },
};
