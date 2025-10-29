import React, { useState, useEffect, ReactNode, useRef, cloneElement, isValidElement } from 'react';
import { useLocation } from 'react-router-dom';
import styles from './RouteCache.module.less';

export interface RouteCacheItem {
  /** 路由路径 */
  path: string;
  /** 缓存的组件 */
  element: ReactNode;
  /** 最后访问时间 */
  lastAccess: number;
}

export interface RouteCacheProps {
  /** 子组件 */
  children: ReactNode;
  /** 最大缓存数量 */
  maxCache?: number;
  /** 缓存的路由列表（为空则缓存所有） */
  cacheRoutes?: string[];
  /** 排除缓存的路由列表 */
  excludeRoutes?: string[];
  /** 是否启用缓存 */
  enabled?: boolean;
}

/**
 * 路由缓存组件
 * 
 * 实现简化的路由缓存功能，保持列表页的滚动位置和搜索状态
 * 
 * @component
 * @example
 * ```tsx
 * <RouteCache 
 *   maxCache={10} 
 *   cacheRoutes={['/users', '/orders']}
 *   excludeRoutes={['/login']}
 * >
 *   <Outlet />
 * </RouteCache>
 * ```
 */
export const RouteCache: React.FC<RouteCacheProps> = ({
  children,
  maxCache = 10,
  cacheRoutes = [],
  excludeRoutes = ['/login', '/register'],
  enabled = true,
}) => {
  const location = useLocation();
  const [cacheMap, setCacheMap] = useState<Map<string, RouteCacheItem>>(new Map());
  const currentPath = location.pathname;

  // 判断当前路由是否应该缓存
  const shouldCache = () => {
    if (!enabled) return false;

    // 检查是否在排除列表中
    if (excludeRoutes.some((route) => currentPath.startsWith(route))) {
      return false;
    }

    // 如果指定了缓存路由，只缓存列表中的
    if (cacheRoutes.length > 0) {
      return cacheRoutes.some((route) => currentPath.startsWith(route));
    }

    // 默认缓存所有（除了排除列表）
    return true;
  };

  useEffect(() => {
    if (shouldCache()) {
      setCacheMap((prev) => {
        const newMap = new Map(prev);

        // 更新或添加当前路由的缓存
        newMap.set(currentPath, {
          path: currentPath,
          element: children,
          lastAccess: Date.now(),
        });

        // 检查缓存数量限制
        if (newMap.size > maxCache) {
          // 删除最旧的缓存（除了当前路由）
          let oldestPath: string | null = null;
          let oldestTime = Date.now();

          newMap.forEach((item, path) => {
            if (path !== currentPath && item.lastAccess < oldestTime) {
              oldestPath = path;
              oldestTime = item.lastAccess;
            }
          });

          if (oldestPath) {
            newMap.delete(oldestPath);
          }
        }

        return newMap;
      });
    }
  }, [currentPath, children, maxCache, shouldCache]);

  if (!enabled || !shouldCache()) {
    return <>{children}</>;
  }

  // 渲染所有缓存的组件
  return (
    <div className={styles.cacheContainer}>
      {Array.from(cacheMap.entries()).map(([path, item]) => (
        <div
          key={path}
          className={styles.cacheItem}
          style={{
            display: path === currentPath ? 'block' : 'none',
          }}
        >
          {item.element}
        </div>
      ))}
    </div>
  );
};

