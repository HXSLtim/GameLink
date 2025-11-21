/**
 * 管理后台路由配置
 */
import React, { lazy } from 'react';
import type { RouteObject } from 'react-router-dom';

// 懒加载页面组件
const AdminLayout = lazy(() => import('./layouts/AdminLayout'));
const Dashboard = lazy(() => import('./pages/Dashboard/Dashboard'));

export const adminRoutes: RouteObject[] = [
  {
    path: '/admin',
    element: <AdminLayout />,
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'users',
        element: <div>用户管理</div>,
      },
      {
        path: 'orders',
        element: <div>订单管理</div>,
      },
      {
        path: 'finance',
        element: <div>财务管理</div>,
      },
      {
        path: 'settings',
        element: <div>系统设置</div>,
      },
      {
        path: 'profile',
        element: <div>个人中心</div>,
      },
    ],
  },
];