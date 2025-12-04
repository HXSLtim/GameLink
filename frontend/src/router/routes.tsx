/**
 * 路由配置
 */
import { lazy, Suspense } from 'react';
import type { RouteConfig } from './types';
import { Spin } from 'antd';

// 布局组件
import ClientLayout from '@/layouts/ClientLayout';
import CompanionLayout from '@/layouts/CompanionLayout';
import AdminLayout from '@/layouts/AdminLayout/index';

// 直接导入的页面
import Login from '@/pages/auth/Login';
import Register from '@/pages/auth/Register';
import NotFound from '@/pages/NotFound';
import ClientHome from '@/pages/client/Home';
import CompanionDashboard from '@/pages/companion/Dashboard';

// 懒加载Admin页面
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));
const UserPage = lazy(() => import('@/pages/admin/User'));
const UserBehavior = lazy(() => import('@/pages/admin/User/Behavior'));
const UserTags = lazy(() => import('@/pages/admin/User/Tags'));
const UserLevel = lazy(() => import('@/pages/admin/User/Level'));
const UserPortrait = lazy(() => import('@/pages/admin/User/Portrait'));
const RolePage = lazy(() => import('@/pages/admin/Role'));
const GamePage = lazy(() => import('@/pages/admin/Game'));
const OrderPage = lazy(() => import('@/pages/admin/Order'));
const PlayerPage = lazy(() => import('@/pages/admin/Player'));

// 兼容旧页面（如果存在）
const AdminSettings = lazy(() => import('@/pages/sys/setting'));
const AdminAudit = lazy(() => import('@/pages/sys/log'));
const AdminMenu = lazy(() => import('@/pages/sys/menu'));
const AdminPermission = lazy(() => import('@/pages/sys/permission'));

// 监控模块页面
const RealtimeMonitor = lazy(() => import('@/pages/admin/Monitor/Realtime'));
const AnalyticsPage = lazy(() => import('@/pages/admin/Monitor/Analytics'));
const KPIDashboard = lazy(() => import('@/pages/admin/Monitor/KPI'));

// Service Item Pages
const ServiceItemList = lazy(() => import('@/pages/biz/service'));
const ServiceItemForm = lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = lazy(() => import('@/pages/biz/service/detail'));

/**
 * 懒加载包装器
 */
const LazyLoad = ({ children }: { children: React.ReactNode }) => (
    <Suspense
        fallback={
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', minHeight: 300 }}>
                <Spin size="large" />
            </div>
        }
    >
        {children}
    </Suspense>
);

export const routes: RouteConfig[] = [
    // 登录页
    {
        path: '/login',
        element: <Login />,
        meta: { title: '登录' }
    },
    // 注册页
    {
        path: '/register',
        element: <Register />,
        meta: { title: '注册' }
    },
    // 客户端
    {
        path: '/',
        element: <ClientLayout />,
        children: [
            {
                path: '',
                element: <ClientHome />,
                meta: { title: '首页' }
            }
        ]
    },
    // 陪玩师端
    {
        path: '/companion',
        element: <CompanionLayout />,
        meta: { roles: ['COMPANION'], requiresAuth: true, title: '陪玩师工作台' },
        children: [
            {
                path: '',
                element: <CompanionDashboard />,
                meta: { title: '工作台' }
            }
        ]
    },
    // 管理端
    {
        path: '/admin',
        element: <AdminLayout />,
        meta: { roles: ['ADMIN', 'CS', 'FINANCE'], requiresAuth: true, title: '管理后台' },
        children: [
            // 仪表盘
            {
                index: true,
                element: <LazyLoad><Dashboard /></LazyLoad>,
                meta: { title: '仪表盘' }
            },
            // 系统管理
            {
                path: 'sys/user',
                element: <LazyLoad><UserPage /></LazyLoad>,
                meta: { title: '用户管理' }
            },
            {
                path: 'sys/user/behavior',
                element: <LazyLoad><UserBehavior /></LazyLoad>,
                meta: { title: '用户行为分析' }
            },
            {
                path: 'sys/user/tags',
                element: <LazyLoad><UserTags /></LazyLoad>,
                meta: { title: '用户标签管理' }
            },
            {
                path: 'sys/user/level',
                element: <LazyLoad><UserLevel /></LazyLoad>,
                meta: { title: '用户等级管理' }
            },
            {
                path: 'sys/user/portrait',
                element: <LazyLoad><UserPortrait /></LazyLoad>,
                meta: { title: '用户画像分析' }
            },
            {
                path: 'sys/role',
                element: <LazyLoad><RolePage /></LazyLoad>,
                meta: { title: '角色管理' }
            },
            {
                path: 'sys/permission',
                element: <LazyLoad><AdminPermission /></LazyLoad>,
                meta: { title: '权限管理' }
            },
            {
                path: 'sys/menu',
                element: <LazyLoad><AdminMenu /></LazyLoad>,
                meta: { title: '菜单管理' }
            },
            {
                path: 'sys/log',
                element: <LazyLoad><AdminAudit /></LazyLoad>,
                meta: { title: '审计日志' }
            },
            // 业务管理
            {
                path: 'biz/game',
                element: <LazyLoad><GamePage /></LazyLoad>,
                meta: { title: '游戏管理' }
            },
            {
                path: 'biz/player',
                element: <LazyLoad><PlayerPage /></LazyLoad>,
                meta: { title: '陪玩师管理' }
            },
            {
                path: 'biz/order',
                element: <LazyLoad><OrderPage /></LazyLoad>,
                meta: { title: '订单管理' }
            },
            // 服务项目
            {
                path: 'biz/service',
                element: <LazyLoad><ServiceItemList /></LazyLoad>,
                meta: { title: '服务项目管理' }
            },
            {
                path: 'biz/service/create',
                element: <LazyLoad><ServiceItemForm /></LazyLoad>,
                meta: { title: '新建服务项目' }
            },
            {
                path: 'biz/service/:id',
                element: <LazyLoad><ServiceItemDetail /></LazyLoad>,
                meta: { title: '服务项目详情' }
            },
            {
                path: 'biz/service/:id/edit',
                element: <LazyLoad><ServiceItemForm /></LazyLoad>,
                meta: { title: '编辑服务项目' }
            },
            // 系统设置
            {
                path: 'settings',
                element: <LazyLoad><AdminSettings /></LazyLoad>,
                meta: { title: '系统设置' }
            },
            // 监控模块
            {
                path: 'monitor/realtime',
                element: <LazyLoad><RealtimeMonitor /></LazyLoad>,
                meta: { title: '实时监控' }
            },
            {
                path: 'monitor/analytics',
                element: <LazyLoad><AnalyticsPage /></LazyLoad>,
                meta: { title: '运营分析' }
            },
            {
                path: 'monitor/kpi',
                element: <LazyLoad><KPIDashboard /></LazyLoad>,
                meta: { title: 'KPI 仪表板' }
            }
        ]
    },
    // 404页面
    {
        path: '*',
        element: <NotFound />
    }
];
