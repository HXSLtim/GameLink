/**
 * 路由配置
 */
import { lazy } from 'react';
import { Navigate } from 'react-router-dom';
import type { RouteConfig } from './types';
import LazyLoad from '@/components/common/LazyLoad';

// 布局组件 - 保持直接导入（核心布局）
import AdminLayout from '@/layouts/AdminLayout/index';

// 懒加载基础页面
const Auth = lazy(() => import('@/pages/auth/Auth'));
const NotFound = lazy(() => import('@/pages/NotFound'));
const Forbidden = lazy(() => import('@/pages/Forbidden'));

// 懒加载Admin页面
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));
const UserPage = lazy(() => import('@/pages/admin/User'));
const UserBehavior = lazy(() => import('@/pages/admin/User/Behavior'));
const UserTags = lazy(() => import('@/pages/admin/User/Tags'));
const UserLevel = lazy(() => import('@/pages/admin/User/Level'));
const UserPortrait = lazy(() => import('@/pages/admin/User/Portrait'));
const RolePage = lazy(() => import('@/pages/admin/Role'));
const RolePermissionConfig = lazy(() => import('@/pages/admin/Role/PermissionConfig'));
const GamePage = lazy(() => import('@/pages/admin/Game'));
const OrderPage = lazy(() => import('@/pages/admin/Order'));
const PlayerPage = lazy(() => import('@/pages/admin/Player'));

// 兼容旧页面（如果存在）
const AdminSettings = lazy(() => import('@/pages/sys/setting'));
const AdminAudit = lazy(() => import('@/pages/sys/log'));
const AdminMenu = lazy(() => import('@/pages/sys/menu'));
const AdminPermission = lazy(() => import('@/pages/sys/permission'));
const UserRolePage = lazy(() => import('@/pages/sys/user-role'));

// 监控模块页面
const RealtimeMonitor = lazy(() => import('@/pages/admin/Monitor/Realtime'));
const AnalyticsPage = lazy(() => import('@/pages/admin/Monitor/Analytics'));
const KPIDashboard = lazy(() => import('@/pages/admin/Monitor/KPI'));

// 评价管理模块页面
const ReviewList = lazy(() => import('@/pages/admin/Review/index'));
const ReviewDetail = lazy(() => import('@/pages/admin/Review/Detail'));
const ReviewModeration = lazy(() => import('@/pages/admin/Review/Moderation'));
const ReviewReports = lazy(() => import('@/pages/admin/Review/Reports'));
const SensitiveWords = lazy(() => import('@/pages/admin/Review/SensitiveWords'));
const ReviewStats = lazy(() => import('@/pages/admin/Review/Stats'));

// 内容管理模块页面
const ContentFeeds = lazy(() => import('@/pages/admin/Content/Feeds'));
const ContentChatMonitor = lazy(() => import('@/pages/admin/Content/ChatMonitor'));
const ContentReports = lazy(() => import('@/pages/admin/Content/Reports'));
const ContentCategories = lazy(() => import('@/pages/admin/Content/Categories'));
const ContentStats = lazy(() => import('@/pages/admin/Content/Stats'));

// Notifications
const AdminNotificationsPage = lazy(() => import('@/pages/admin/Notifications'));

// Service Item Pages
const ServiceItemList = lazy(() => import('@/pages/biz/service'));
const ServiceItemForm = lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = lazy(() => import('@/pages/biz/service/detail'));

export const routes: RouteConfig[] = [
    // 登录/注册页（Tab 切换）
    {
        path: '/login',
        element: <LazyLoad><Auth /></LazyLoad>,
        meta: { title: '登录' }
    },
    {
        path: '/register',
        element: <LazyLoad><Auth /></LazyLoad>,
        meta: { title: '注册' }
    },
    // 重定向首页到管理端
    {
        path: '/',
        element: <Navigate to="/admin" replace />
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
                path: 'sys/role/:id/permissions',
                element: <LazyLoad><RolePermissionConfig /></LazyLoad>,
                meta: { title: '角色权限配置', permission: 'admin.roles.permissions' }
            },
            {
                path: 'sys/permission',
                element: <LazyLoad><AdminPermission /></LazyLoad>,
                meta: { title: '权限管理' }
            },
            {
                path: 'sys/user-role',
                element: <LazyLoad><UserRolePage /></LazyLoad>,
                meta: { title: '用户角色分配', permission: 'admin.roles.assign-user' }
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
            {
                path: 'notifications',
                element: <LazyLoad><AdminNotificationsPage /></LazyLoad>,
                meta: { title: '通知中心' }
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
            },
            // 评价管理模块
            {
                path: 'reviews/list',
                element: <LazyLoad><ReviewList /></LazyLoad>,
                meta: { title: '评价列表', permission: 'admin.reviews.list' }
            },
            {
                path: 'reviews/detail/:id',
                element: <LazyLoad><ReviewDetail /></LazyLoad>,
                meta: { title: '评价详情', permission: 'admin.reviews.read' }
            },
            {
                path: 'reviews/moderation',
                element: <LazyLoad><ReviewModeration /></LazyLoad>,
                meta: { title: '评价审核', permission: 'admin.reviews.pending.list' }
            },
            {
                path: 'review-reports',
                element: <LazyLoad><ReviewReports /></LazyLoad>,
                meta: { title: '举报管理', permission: 'admin.review-reports.list' }
            },
            {
                path: 'sensitive-words',
                element: <LazyLoad><SensitiveWords /></LazyLoad>,
                meta: { title: '敏感词管理', permission: 'admin.sensitive-words.list' }
            },
            {
                path: 'reviews/stats',
                element: <LazyLoad><ReviewStats /></LazyLoad>,
                meta: { title: '评价统计', permission: 'admin.reviews.stats.list' }
            },
            // 内容管理模块
            {
                path: 'content/feeds',
                element: <LazyLoad><ContentFeeds /></LazyLoad>,
                meta: { title: '动态审核', permission: 'admin.content.feeds.list' }
            },
            {
                path: 'content/chat',
                element: <LazyLoad><ContentChatMonitor /></LazyLoad>,
                meta: { title: '聊天监控', permission: 'admin.content.chat.list' }
            },
            {
                path: 'content/reports',
                element: <LazyLoad><ContentReports /></LazyLoad>,
                meta: { title: '举报管理', permission: 'admin.content.reports.list' }
            },
            {
                path: 'content/categories',
                element: <LazyLoad><ContentCategories /></LazyLoad>,
                meta: { title: '内容分类', permission: 'admin.content.categories.list' }
            },
            {
                path: 'content/stats',
                element: <LazyLoad><ContentStats /></LazyLoad>,
                meta: { title: '内容统计', permission: 'admin.content.stats.list' }
            }
        ]
    },
    // 403 禁止访问页面
    {
        path: '/403',
        element: <LazyLoad><Forbidden /></LazyLoad>,
        meta: { title: '无权限访问' }
    },
    // 404页面
    {
        path: '*',
        element: <LazyLoad><NotFound /></LazyLoad>
    }
];
