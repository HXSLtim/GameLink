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

// 用户端组件
const UserHome = React.lazy(() => import('@/pages/user/Home'));
const UserOrders = React.lazy(() => import('@/pages/user/Orders'));

// 陪玩师端组件
const PlayerHome = React.lazy(() => import('@/pages/player/Home'));
const PlayerOrders = React.lazy(() => import('@/pages/player/Orders'));
const PlayerEarnings = React.lazy(() => import('@/pages/player/Earnings'));

// 用户端扩展组件
const UserWallet = React.lazy(() => import('@/pages/user/Wallet'));
const UserRanking = React.lazy(() => import('@/pages/user/Ranking'));

// 管理端扩展组件
const AdminService = React.lazy(() => import('@/pages/admin/Service'));
const AdminAlert = React.lazy(() => import('@/pages/admin/Alert'));
const AdminDispute = React.lazy(() => import('@/pages/admin/Dispute'));
const AdminUserTag = React.lazy(() => import('@/pages/admin/UserTag'));
const AdminSettlementCompany = React.lazy(() => import('@/pages/admin/SettlementCompany'));
const AdminRankingCommission = React.lazy(() => import('@/pages/admin/RankingCommission'));
const AdminRoutingRule = React.lazy(() => import('@/pages/admin/RoutingRule'));
const AdminUserBehavior = React.lazy(() => import('@/pages/admin/UserBehavior'));
const AdminWithdrawRouting = React.lazy(() => import('@/pages/admin/WithdrawRouting'));

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
    
    // 财务管理
    'Withdraw': React.lazy(() => import('@/pages/admin/Withdraw')),
    'Commission': React.lazy(() => import('@/pages/admin/Commission')),
    
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
    'biz/withdraw': React.lazy(() => import('@/pages/admin/Withdraw')),
    'biz/commission': React.lazy(() => import('@/pages/admin/Commission')),
    'finance/withdraw': React.lazy(() => import('@/pages/admin/Withdraw')),
    'finance/commission': React.lazy(() => import('@/pages/admin/Commission')),
    
    // 用户端
    'UserHome': UserHome,
    'UserOrders': UserOrders,
    'user/home': UserHome,
    'user/orders': UserOrders,
    
    // 陪玩师端
    'PlayerHome': PlayerHome,
    'PlayerOrders': PlayerOrders,
    'PlayerEarnings': PlayerEarnings,
    'player/home': PlayerHome,
    'player/orders': PlayerOrders,
    'player/earnings': PlayerEarnings,
    
    // 用户端扩展
    'UserWallet': UserWallet,
    'UserRanking': UserRanking,
    'user/wallet': UserWallet,
    'user/ranking': UserRanking,
    
    // 管理端扩展
    'AdminService': AdminService,
    'AdminAlert': AdminAlert,
    'AdminDispute': AdminDispute,
    'AdminUserTag': AdminUserTag,
    'AdminSettlementCompany': AdminSettlementCompany,
    'AdminRankingCommission': AdminRankingCommission,
    'AdminRoutingRule': AdminRoutingRule,
    'AdminUserBehavior': AdminUserBehavior,
    'AdminWithdrawRouting': AdminWithdrawRouting,
    'Service': AdminService,
    'Alert': AdminAlert,
    'Dispute': AdminDispute,
    'UserTag': AdminUserTag,
    'SettlementCompany': AdminSettlementCompany,
    'RankingCommission': AdminRankingCommission,
    'RoutingRule': AdminRoutingRule,
    'UserBehavior': AdminUserBehavior,
    'WithdrawRouting': AdminWithdrawRouting,
    'admin/service': AdminService,
    'admin/alert': AdminAlert,
    'admin/dispute': AdminDispute,
    'admin/user-tag': AdminUserTag,
    'admin/settlement-company': AdminSettlementCompany,
    'admin/ranking-commission': AdminRankingCommission,
    'admin/routing-rule': AdminRoutingRule,
    'admin/user-behavior': AdminUserBehavior,
    'admin/withdraw-routing': AdminWithdrawRouting,
    
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
