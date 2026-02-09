import React from 'react';
import LazyLoad from '@/components/common/LazyLoad';
import LayoutOutlet from './LayoutOutlet';

import { logger } from '@/utils/logger';
// Lazy load components
const Audit = React.lazy(() => import('@/pages/sys/log'));
const Settings = React.lazy(() => import('@/pages/sys/setting'));
const ServiceItemList = React.lazy(() => import('@/pages/biz/service'));
const ServiceItemCreate = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = React.lazy(() => import('@/pages/biz/service/detail'));
const Menu = React.lazy(() => import('@/pages/sys/menu'));
const Permission = React.lazy(() => import('@/pages/sys/permission'));
const Role = React.lazy(() => import('@/pages/admin/Role'));

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

// VIP管理
const VIPLevels = React.lazy(() => import('@/pages/admin/VIP'));
const VIPConfig = React.lazy(() => import('@/pages/admin/VIP/Config'));

// 充值管理
const AdminRecharge = React.lazy(() => import('@/pages/admin/Recharge'));
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
    'biz/service/detail': ServiceItemDetail,  // 编辑使用详情页模态弹窗
    'biz/withdraw': React.lazy(() => import('@/pages/admin/Withdraw')),
    'biz/commission': React.lazy(() => import('@/pages/admin/Commission')),
    'finance/withdraw': React.lazy(() => import('@/pages/admin/Withdraw')),
    'finance/commission': React.lazy(() => import('@/pages/admin/Commission')),
    
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

    // VIP管理
    'VIPLevels': VIPLevels,
    'VIPConfig': VIPConfig,
    'VIP': VIPLevels,
    'vip/levels': VIPLevels,
    'vip/config': VIPConfig,
    'admin/vip': VIPLevels,
    'admin/vip/levels': VIPLevels,
    'admin/marketing/vip': VIPLevels,

    // 优惠券管理
    'Coupon': React.lazy(() => import('@/pages/admin/Coupon')),
    'admin/coupon': React.lazy(() => import('@/pages/admin/Coupon')),
    'admin/marketing/coupon': React.lazy(() => import('@/pages/admin/Coupon')),

    // 推荐系统
    'Referral': React.lazy(() => import('@/pages/admin/Referral')),
    'admin/referral': React.lazy(() => import('@/pages/admin/Referral')),
    'admin/marketing/referral': React.lazy(() => import('@/pages/admin/Referral')),

    // 团队管理
    'Team': React.lazy(() => import('@/pages/admin/Team')),
    'admin/team': React.lazy(() => import('@/pages/admin/Team')),
    'admin/marketing/team': React.lazy(() => import('@/pages/admin/Team')),

    // 活动管理
    'Activity': React.lazy(() => import('@/pages/admin/Activity')),
    'admin/activity': React.lazy(() => import('@/pages/admin/Activity')),
    'admin/marketing/activity': React.lazy(() => import('@/pages/admin/Activity')),

    // 充值管理
    'Recharge': AdminRecharge,
    'admin/recharge': AdminRecharge,
    'admin/payment/recharge': AdminRecharge,
    'admin/vip/config': VIPConfig,
    // 支付记录
    'PaymentRecords': React.lazy(() => import('@/pages/admin/PaymentRecords')),
    'admin/payment/records': React.lazy(() => import('@/pages/admin/PaymentRecords')),

    'content/feeds': React.lazy(() => import('@/pages/admin/Content/Feeds')),
    'content/chat': React.lazy(() => import('@/pages/admin/Content/ChatMonitor')),
    'content/reports': React.lazy(() => import('@/pages/admin/Content/Reports')),
    'content/categories': React.lazy(() => import('@/pages/admin/Content/Categories')),
    'content/stats': React.lazy(() => import('@/pages/admin/Content/Stats')),

    // 段位管理
    'GameRank': React.lazy(() => import('@/pages/admin/GameRank')),
    'admin/game-rank': React.lazy(() => import('@/pages/admin/GameRank')),
    'admin/biz/game-rank': React.lazy(() => import('@/pages/admin/GameRank')),

    // 段位审核
    'PlayerRank': React.lazy(() => import('@/pages/admin/PlayerRank')),
    'admin/player-rank': React.lazy(() => import('@/pages/admin/PlayerRank')),
    'admin/biz/player-rank': React.lazy(() => import('@/pages/admin/PlayerRank')),

    // 实名审核（管理端）
    'PlayerCertificationAdmin': React.lazy(() => import('@/pages/admin/PlayerCertification')),
    'AdminPlayerCertification': React.lazy(() => import('@/pages/admin/PlayerCertification')),
    'admin/player-certification': React.lazy(() => import('@/pages/admin/PlayerCertification')),
    'admin/biz/player-certification': React.lazy(() => import('@/pages/admin/PlayerCertification')),

    // 个人中心
    'Profile': React.lazy(() => import('@/pages/admin/Profile')),
    'admin/profile': React.lazy(() => import('@/pages/admin/Profile')),

    // 实时监控大屏
    'Monitor': React.lazy(() => import('@/pages/admin/Monitor')),
    'admin/monitor': React.lazy(() => import('@/pages/admin/Monitor')),

    // 聊天管理
    'ChatRooms': React.lazy(() => import('@/pages/adminChat/rooms')),
    'ChatRecords': React.lazy(() => import('@/pages/adminChat/records')),
    'admin/chat/rooms': React.lazy(() => import('@/pages/adminChat/rooms')),
    'admin/chat/records': React.lazy(() => import('@/pages/adminChat/records')),

    // 审计日志
    'Audit': React.lazy(() => import('@/pages/sys/log')),
    'admin/sys/log': React.lazy(() => import('@/pages/sys/log')),

    // 告警管理
    'Alert': AdminAlert,
    'admin/alert': AdminAlert,

    // 用户角色分配
    'UserRole': React.lazy(() => import('@/pages/sys/user-role')),
    'admin/sys/user-role': React.lazy(() => import('@/pages/sys/user-role')),

    // 支付记录
    'admin/payment/payment-records': React.lazy(() => import('@/pages/admin/PaymentRecords')),

    // 结算公司管理（兼容旧路径）
    'Settlement': React.lazy(() => import('@/pages/admin/Settlement')),
    'admin/settlement': React.lazy(() => import('@/pages/admin/Settlement')),
    'SettlementPlayers': React.lazy(() => import('@/pages/admin/Settlement/Players')),
    'admin/settlement/players': React.lazy(() => import('@/pages/admin/Settlement/Players')),

    // 评价设置
    'ReviewSettings': React.lazy(() => import('@/pages/admin/Review/Settings')),
    'admin/reviews/settings': React.lazy(() => import('@/pages/admin/Review/Settings')),

    // 评价详情
    'ReviewDetail': React.lazy(() => import('@/pages/admin/Review/Detail')),
    'admin/reviews/detail': React.lazy(() => import('@/pages/admin/Review/Detail')),
};

export const getComponent = (componentKey: string) => {
    const Component = componentMap[componentKey];
    if (!Component) {
        logger.warn(`Component not found for key: ${componentKey}`);
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
