/**
 * API 请求 Hook
 * 封装了 API 调用的通用逻辑，包括加载状态、错误处理等
 */

import { useState, useCallback } from 'react';

export interface UseApiOptions<T> {
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
  onFinally?: () => void;
}

export interface UseApiResult<T, P extends any[] = any[]> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  execute: (...args: P) => Promise<T | null>;
  reset: () => void;
}

/**
 * 封装 API 调用的 Hook
 * 
 * @param apiFunction API 函数
 * @param options 配置选项
 * @returns 包含数据、加载状态、错误和执行函数的响应对象
 * 
 * @example
 * ```typescript
 * const { data, loading, error, execute } = useApi(authApi.login);
 * 
 * const handleLogin = async () => {
 *   const result = await execute({ username, password });
 *   if (result) {
 *     console.log('登录成功', result);
 *   }
 * };
 * ```
 */
export function useApi<T, P extends any[] = any[]>(
  apiFunction: (...args: P) => Promise<T>,
  options: UseApiOptions<T> = {}
): UseApiResult<T, P> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(
    async (...args: P): Promise<T | null> => {
      try {
        setLoading(true);
        setError(null);
        
        const result = await apiFunction(...args);
        setData(result);
        options.onSuccess?.(result);
        
        return result;
      } catch (err) {
        const apiError = err instanceof Error ? err : new Error('Unknown error');
        setError(apiError);
        options.onError?.(apiError);
        
        return null;
      } finally {
        setLoading(false);
        options.onFinally?.();
      }
    },
    [apiFunction, options]
  );

  const reset = useCallback(() => {
    setData(null);
    setLoading(false);
    setError(null);
  }, []);

  return {
    data,
    loading,
    error,
    execute,
    reset,
  };
}

/**
 * 用于 GET 请求的 Hook，支持自动加载
 * 
 * @param apiFunction API 函数
 * @param params 请求参数
 * @param options 配置选项
 * @returns 包含数据、加载状态、错误和刷新函数的响应对象
 * 
 * @example
 * ```typescript
 * const { data, loading, error, refresh } = useApiQuery(playerApi.getPlayerById, 1);
 * ```
 */
export function useApiQuery<T, P extends any[] = any[]>(
  apiFunction: (...args: P) => Promise<T>,
  ...params: P
) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(
    async (...args: P) => {
      try {
        setLoading(true);
        setError(null);
        
        const result = await apiFunction(...args);
        setData(result);
        
        return result;
      } catch (err) {
        const apiError = err instanceof Error ? err : new Error('Unknown error');
        setError(apiError);
        
        return null;
      } finally {
        setLoading(false);
      }
    },
    [apiFunction]
  );

  const refresh = useCallback(() => {
    execute(...params);
  }, [execute, params]);

  // 自动执行一次
  useState(() => {
    execute(...params);
  });

  return {
    data,
    loading,
    error,
    refresh,
    execute,
  };
}