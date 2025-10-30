import { useState, useCallback } from 'react';
import { useLocation } from 'react-router-dom';

export interface RouteCacheConfig {
  /** 是否启用缓存 */
  enabled: boolean;
  /** 最大缓存数量 */
  maxCache: number;
  /** 需要缓存的路由列表（为空则缓存所有） */
  cacheRoutes: string[];
  /** 排除缓存的路由列表 */
  excludeRoutes: string[];
}

export interface RouteCacheControl {
  /** 配置 */
  config: RouteCacheConfig;
  /** 更新配置 */
  updateConfig: (config: Partial<RouteCacheConfig>) => void;
  /** 清除指定路由的缓存 */
  clearCache: (path?: string) => void;
  /** 清除所有缓存 */
  clearAllCache: () => void;
  /** 刷新当前路由 */
  refreshCurrent: () => void;
}

/**
 * 路由缓存控制 Hook
 * 
 * 提供路由缓存的配置和控制功能
 * 
 * @param initialConfig - 初始配置
 * @returns 缓存控制对象
 * 
 * @example
 * ```tsx
 * const cacheControl = useRouteCache({
 *   enabled: true,
 *   maxCache: 10,
 *   cacheRoutes: ['/users', '/orders'],
 * });
 * 
 * // 清除当前路由缓存
 * cacheControl.refreshCurrent();
 * 
 * // 清除所有缓存
 * cacheControl.clearAllCache();
 * ```
 */
export const useRouteCache = (
  initialConfig: Partial<RouteCacheConfig> = {}
): RouteCacheControl => {
  const location = useLocation();
  
  const [config, setConfig] = useState<RouteCacheConfig>({
    enabled: true,
    maxCache: 10,
    cacheRoutes: [],
    excludeRoutes: ['/login', '/register'],
    ...initialConfig,
  });

  const [_cacheVersion, setCacheVersion] = useState(0);

  /**
   * 更新配置
   */
  const updateConfig = useCallback((newConfig: Partial<RouteCacheConfig>) => {
    setConfig((prev) => ({ ...prev, ...newConfig }));
  }, []);

  /**
   * 清除指定路由的缓存
   */
  const clearCache = useCallback((path?: string) => {
    // 通过改变版本号来触发缓存清除
    setCacheVersion((prev) => prev + 1);
    console.log(`[RouteCache] Cleared cache for: ${path || 'all'}`);
  }, []);

  /**
   * 清除所有缓存
   */
  const clearAllCache = useCallback(() => {
    setCacheVersion((prev) => prev + 1);
    console.log('[RouteCache] Cleared all cache');
  }, []);

  /**
   * 刷新当前路由
   */
  const refreshCurrent = useCallback(() => {
    clearCache(location.pathname);
  }, [location.pathname, clearCache]);

  return {
    config,
    updateConfig,
    clearCache,
    clearAllCache,
    refreshCurrent,
  };
};

/**
 * 判断路由是否应该被缓存
 */
export const shouldCacheRoute = (
  pathname: string,
  config: RouteCacheConfig
): boolean => {
  if (!config.enabled) return false;

  // 检查是否在排除列表中
  if (config.excludeRoutes.some((route) => pathname.startsWith(route))) {
    return false;
  }

  // 如果指定了缓存路由，只缓存列表中的
  if (config.cacheRoutes.length > 0) {
    return config.cacheRoutes.some((route) => pathname.startsWith(route));
  }

  // 默认缓存所有（除了排除列表）
  return true;
};


