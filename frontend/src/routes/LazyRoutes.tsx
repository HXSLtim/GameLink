import { lazy } from 'react';
import { PageSkeleton } from 'components';

/**
 * 懒加载页面组件
 * 统一使用 named exports 避免配置差异
 */
export const Dashboard = lazy(() =>
  import('pages/Dashboard').then((module) => ({
    default: module.Dashboard,
  })),
);

export const Users = lazy(() =>
  import('pages/Users').then((module) => ({
    default: module.Users,
  })),
);

export const Orders = lazy(() =>
  import('pages/Orders').then((module) => ({
    default: module.Orders,
  })),
);

export const Permissions = lazy(() =>
  import('pages/Permissions').then((module) => ({
    default: module.Permissions,
  })),
);

export const Login = lazy(() =>
  import('pages/Login').then((module) => ({
    default: module.Login,
  })),
);

/**
 * 路由加载组件（骨架屏）
 * 在页面懒加载时显示，提供更好的用户体验
 */
export const RouteLoading: React.FC = () => <PageSkeleton />;
