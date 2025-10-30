import { useMemo } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import type { BreadcrumbItem } from '../components/Breadcrumb/Breadcrumb';

/**
 * 路由配置映射
 */
const ROUTE_MAP: Record<string, string> = {
  dashboard: '仪表盘',
  orders: '订单管理',
  games: '游戏管理',
  players: '陪玩师管理',
  users: '用户管理',
  payments: '支付管理',
  reviews: '评价管理',
  reports: '数据报表',
  permissions: '权限管理',
  settings: '系统设置',
};

/**
 * 获取路由标题
 */
const getRouteTitle = (path: string, params?: Record<string, string>): string => {
  // 处理动态路由参数
  if (params && Object.keys(params).length > 0) {
    // 如果是详情页，显示 "详情"
    if (params.id) {
      return '详情';
    }
  }

  return ROUTE_MAP[path] || path;
};

/**
 * 生成面包屑路径
 */
const generateBreadcrumbs = (pathname: string, params: Record<string, string>): BreadcrumbItem[] => {
  const paths = pathname.split('/').filter(Boolean);
  const breadcrumbs: BreadcrumbItem[] = [];

  // 始终添加首页
  breadcrumbs.push({
    label: '首页',
    path: '/dashboard',
  });

  // 构建面包屑路径
  let currentPath = '';
  paths.forEach((path, index) => {
    currentPath += `/${path}`;

    // 跳过参数值（通常是数字或ID）
    if (params[path] || /^\d+$/.test(path)) {
      // 这是参数值，添加详情标签但不作为链接
      breadcrumbs.push({
        label: getRouteTitle(path, params),
        // 最后一项不需要路径
        path: index < paths.length - 1 ? currentPath : undefined,
      });
    } else {
      breadcrumbs.push({
        label: getRouteTitle(path, params),
        // 最后一项不需要路径
        path: index < paths.length - 1 ? currentPath : undefined,
      });
    }
  });

  return breadcrumbs;
};

/**
 * 面包屑 Hook
 * 
 * 根据当前路由自动生成面包屑导航
 * 
 * @returns 面包屑项数组
 * 
 * @example
 * ```tsx
 * const breadcrumbs = useBreadcrumb();
 * return <Breadcrumb items={breadcrumbs} />;
 * ```
 */
export const useBreadcrumb = (): BreadcrumbItem[] => {
  const location = useLocation();
  const params = useParams();

  const breadcrumbs = useMemo(() => {
    // 首页不显示面包屑
    if (location.pathname === '/' || location.pathname === '/dashboard') {
      return [];
    }

    return generateBreadcrumbs(location.pathname, params as Record<string, string>);
  }, [location.pathname, params]);

  return breadcrumbs;
};

/**
 * 自定义面包屑 Hook
 * 
 * 允许手动指定面包屑项
 * 
 * @param items - 自定义面包屑项
 * @returns 面包屑项数组
 * 
 * @example
 * ```tsx
 * const breadcrumbs = useCustomBreadcrumb([
 *   { label: '首页', path: '/dashboard' },
 *   { label: '用户管理', path: '/users' },
 *   { label: '用户详情' },
 * ]);
 * ```
 */
export const useCustomBreadcrumb = (items: BreadcrumbItem[]): BreadcrumbItem[] => {
  return useMemo(() => items, [items]);
};


