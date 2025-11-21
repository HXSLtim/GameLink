/**
 * 陪玩师端路由配置
 */
import React, { lazy } from 'react';
import type { RouteObject } from 'react-router-dom';

// 懒加载页面组件
const PlayerLayout = lazy(() => import('./layouts/PlayerLayout'));
const Dashboard = lazy(() => import('./pages/Dashboard/Dashboard'));

export const playerRoutes: RouteObject[] = [
  {
    path: '/player',
    element: <PlayerLayout />,
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'orders',
        element: <div>我的订单</div>,
      },
      {
        path: 'earnings',
        element: <div>收益管理</div>,
      },
      {
        path: 'profile',
        element: <div>个人资料</div>,
      },
    ],
  },
];