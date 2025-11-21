/**
 * 应用路由配置 - 整合三端路由
 */
import { RouteObject } from 'react-router-dom';
import { adminRoutes } from '@/apps/admin/routes';
import { playerRoutes } from '@/apps/player/routes';
import { userRoutes } from '@/apps/user/routes';

// 导出所有路由配置
export const appRoutes: RouteObject[] = [
  ...adminRoutes,
  ...playerRoutes,
  ...userRoutes,
];