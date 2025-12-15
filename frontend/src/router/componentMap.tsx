import React from 'react';
import LazyLoad from '@/components/common/LazyLoad';
import LayoutOutlet from './LayoutOutlet';

// Lazy load components
const Audit = React.lazy(() => import('@/pages/sys/log'));
const Settings = React.lazy(() => import('@/pages/sys/setting'));
const ServiceItemList = React.lazy(() => import('@/pages/biz/service'));
const ServiceItemCreate = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemEdit = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = React.lazy(() => import('@/pages/biz/service/detail'));
const Menu = React.lazy(() => import('@/pages/sys/menu'));
const Permission = React.lazy(() => import('@/pages/sys/permission'));
const Role = React.lazy(() => import('@/pages/admin/Role'));

// 组件映射表：支持数据库中的组件名称
export const componentMap: Record<string, React.LazyExoticComponent<React.ComponentType> | React.FC> = {
    // 仪表盘
    'Dashboard': React.lazy(() => import('@/pages/admin/Dashboard')),
    
    // 系统管理 - Layout 用于父级菜单容器
    'Layout': LayoutOutlet, // 父级菜单使用 Outlet 渲染子路由
    'User': React.lazy(() => import('@/pages/admin/User')),
    'Role': Role,
    'Permission': Permission,
    'Menu': Menu,
    
    // 业务管理
    'Game': React.lazy(() => import('@/pages/admin/Game')),
    'Player': React.lazy(() => import('@/pages/admin/Player')),
    'Order': React.lazy(() => import('@/pages/admin/Order')),
    'ServiceItem': ServiceItemList,
    
    // 监控中心
    'RealtimeMonitor': React.lazy(() => import('@/pages/admin/Monitor/Realtime')),
    'Analytics': React.lazy(() => import('@/pages/admin/Monitor/Analytics')),
    'KPIDashboard': React.lazy(() => import('@/pages/admin/Monitor/KPI')),
    
    // 内容管理
    'ContentFeeds': React.lazy(() => import('@/pages/admin/Content/Feeds')),
    'ContentChat': React.lazy(() => import('@/pages/admin/Content/ChatMonitor')),
    'ContentReports': React.lazy(() => import('@/pages/admin/Content/Reports')),
    'ContentCategories': React.lazy(() => import('@/pages/admin/Content/Categories')),
    'ContentStats': React.lazy(() => import('@/pages/admin/Content/Stats')),
    
    // 评价管理
    'ReviewList': React.lazy(() => import('@/pages/admin/Review/index')),
    'ReviewModeration': React.lazy(() => import('@/pages/admin/Review/Moderation')),
    'ReviewReports': React.lazy(() => import('@/pages/admin/Review/Reports')),
    'SensitiveWords': React.lazy(() => import('@/pages/admin/Review/SensitiveWords')),
    'ReviewStats': React.lazy(() => import('@/pages/admin/Review/Stats')),
    
    // 系统设置
    'Settings': Settings,
    
    // 兼容旧的路径格式
    'sys/dashboard': React.lazy(() => import('@/pages/admin/Dashboard')),
    'sys/user': React.lazy(() => import('@/pages/admin/User')),
    'sys/role': Role,
    'sys/permission': Permission,
    'sys/menu': Menu,
    'sys/log': Audit,
    'sys/setting': Settings,
    'biz/game': React.lazy(() => import('@/pages/admin/Game')),
    'biz/order': React.lazy(() => import('@/pages/admin/Order')),
    'biz/player': React.lazy(() => import('@/pages/admin/Player')),
    'biz/service': ServiceItemList,
    'biz/service/list': ServiceItemList,
    'biz/service/create': ServiceItemCreate,
    'biz/service/edit': ServiceItemEdit,
    'biz/service/detail': ServiceItemDetail,
    'content/feeds': React.lazy(() => import('@/pages/admin/Content/Feeds')),
    'content/chat': React.lazy(() => import('@/pages/admin/Content/ChatMonitor')),
    'content/reports': React.lazy(() => import('@/pages/admin/Content/Reports')),
    'content/categories': React.lazy(() => import('@/pages/admin/Content/Categories')),
    'content/stats': React.lazy(() => import('@/pages/admin/Content/Stats')),
};

export const getComponent = (componentKey: string) => {
    const Component = componentMap[componentKey];
    if (!Component) {
        console.warn(`Component not found for key: ${componentKey}`);
        return null;
    }
    
    // Layout 组件不需要 LazyLoad 包装
    if (componentKey === 'Layout') {
        return Component as React.FC;
    }
    
    // 其他组件包装 LazyLoad
    const LazyComponent = Component as React.LazyExoticComponent<React.ComponentType>;
    return () => <LazyLoad><LazyComponent /></LazyLoad>;
};
