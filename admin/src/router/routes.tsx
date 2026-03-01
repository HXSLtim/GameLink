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
const AdminLogin = lazy(() => import('@/pages/admin/Login'));
const NotFound = lazy(() => import('@/pages/NotFound'));
const Forbidden = lazy(() => import('@/pages/Forbidden'));

// 懒加载Admin页面
const Dashboard = lazy(() => import('@/pages/admin/Dashboard'));
const UserPage = lazy(() => import('@/pages/admin/User'));
const UserBehavior = lazy(() => import('@/pages/admin/User/Behavior'));
const UserTags = lazy(() => import('@/pages/admin/User/Tags'));
const UserLevel = lazy(() => import('@/pages/admin/User/Level'));
const UserPortrait = lazy(() => import('@/pages/admin/User/Portrait'));
const BatchPage = lazy(() => import('@/pages/admin/Batch'));
const RolePage = lazy(() => import('@/pages/admin/Role'));
const RolePermissionConfig = lazy(() => import('@/pages/admin/Role/PermissionConfig'));
const GamePage = lazy(() => import('@/pages/admin/Game'));
const OrderPage = lazy(() => import('@/pages/admin/Order'));
const PlayerPage = lazy(() => import('@/pages/admin/Player'));

// 段位管理相关页面
const GameRankPage = lazy(() => import('@/pages/admin/GameRank'));
const PlayerRankPage = lazy(() => import('@/pages/admin/PlayerRank'));
const PlayerCertificationPage = lazy(() => import('@/pages/admin/PlayerCertification'));

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
const AdminAlert = lazy(() => import('@/pages/admin/Alert'));

// 评价管理模块页面
const ReviewList = lazy(() => import('@/pages/admin/Review/index'));
const ReviewDetail = lazy(() => import('@/pages/admin/Review/Detail'));
const ReviewModeration = lazy(() => import('@/pages/admin/Review/Moderation'));
const ReviewReports = lazy(() => import('@/pages/admin/Review/Reports'));
const SensitiveWords = lazy(() => import('@/pages/admin/Review/SensitiveWords'));
const ReviewStats = lazy(() => import('@/pages/admin/Review/Stats'));
const ReviewSettings = lazy(() => import('@/pages/admin/Review/Settings'));

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

// VIP管理页面
const VIPLevels = lazy(() => import('@/pages/admin/VIP'));
const VIPConfig = lazy(() => import('@/pages/admin/VIP/Config'));

// 营销模块页面
const CouponPage = lazy(() => import('@/pages/admin/Coupon'));
const ActivityPage = lazy(() => import('@/pages/admin/Activity'));
const RechargePage = lazy(() => import('@/pages/admin/Recharge'));
const TeamPage = lazy(() => import('@/pages/admin/Team'));
const ReferralPage = lazy(() => import('@/pages/admin/Referral'));

// 财务管理页面
const WithdrawPage = lazy(() => import('@/pages/admin/Withdraw'));
const CommissionPage = lazy(() => import('@/pages/admin/Commission'));
const RankingCommissionPage = lazy(() => import('@/pages/admin/RankingCommission'));

// 业务管理页面
const DisputePage = lazy(() => import('@/pages/admin/Dispute'));
const RoutingPage = lazy(() => import('@/pages/admin/Routing'));

// 支付管理页面
const PaymentRecordsPage = lazy(() => import('@/pages/admin/PaymentRecords'));

// 结算公司管理页面
const SettlementPage = lazy(() => import('@/pages/admin/Settlement'));
const SettlementPlayersPage = lazy(() => import('@/pages/admin/Settlement/Players'));

// 对账管理页面
const ReconciliationPage = lazy(() => import('@/pages/admin/Reconciliation'));

// 聊天管理页面
const ChatRecordsPage = lazy(() => import('@/pages/adminChat/records'));
const ChatRoomsPage = lazy(() => import('@/pages/adminChat/rooms'));

// 实时监控页面
const MonitorPage = lazy(() => import('@/pages/admin/Monitor'));

// 个人中心页面
const ProfilePage = lazy(() => import('@/pages/admin/Profile'));


export const routes: RouteConfig[] = [
    // 管理端登录页
    {
        path: '/admin/login',
        element: <LazyLoad><AdminLogin /></LazyLoad>,
        meta: { title: '管理后台登录' }
    },
    // 重定向首页到管理后台，/login 到管理端登录
    {
        path: '/',
        element: <Navigate to="/admin" replace />
    },
    {
        path: '/login',
        element: <Navigate to="/admin/login" replace />
    },
    // 管理端
    {
        path: '/admin',
        element: <AdminLayout />,
        meta: { roles: ['ADMIN'], requiresAuth: true, title: '管理后台' },
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
                path: 'sys/batch',
                element: <LazyLoad><BatchPage /></LazyLoad>,
                meta: { title: '批量操作', permission: 'admin.users.update' }
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
            // 段位管理
            {
                path: 'biz/game-rank',
                element: <LazyLoad><GameRankPage /></LazyLoad>,
                meta: { title: '段位管理', permission: 'admin.game-ranks.list' }
            },
            // 段位审核
            {
                path: 'biz/player-rank',
                element: <LazyLoad><PlayerRankPage /></LazyLoad>,
                meta: { title: '段位审核', permission: 'admin.player-ranks.list' }
            },
            // 实名审核
            {
                path: 'biz/player-certification',
                element: <LazyLoad><PlayerCertificationPage /></LazyLoad>,
                meta: { title: '实名审核', permission: 'admin.player-certifications.list' }
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
            {
                path: 'alert',
                element: <LazyLoad><AdminAlert /></LazyLoad>,
                meta: { title: '告警管理', permission: 'admin.alert.list' }
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
                path: 'reviews/reports',
                element: <LazyLoad><ReviewReports /></LazyLoad>,
                meta: { title: '举报管理', permission: 'admin.review-reports.list' }
            },
            {
                path: 'reviews/sensitive-words',
                element: <LazyLoad><SensitiveWords /></LazyLoad>,
                meta: { title: '敏感词管理', permission: 'admin.sensitive-words.list' }
            },
            {
                path: 'reviews/stats',
                element: <LazyLoad><ReviewStats /></LazyLoad>,
                meta: { title: '评价统计', permission: 'admin.reviews.stats.list' }
            },
            {
                path: 'reviews/settings',
                element: <LazyLoad><ReviewSettings /></LazyLoad>,
                meta: { title: '评价设置', permission: 'admin.reviews.settings.update' }
            },
            // 内容管理模块
            {
                path: 'content/feeds',
                element: <LazyLoad><ContentFeeds /></LazyLoad>,
                meta: { title: '动态审核', permission: 'admin.content.feeds.list' }
            },
            {
                path: 'content/chat-monitor',
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
            },
            // 聊天管理模块
            {
                path: 'chat/records',
                element: <LazyLoad><ChatRecordsPage /></LazyLoad>,
                meta: { title: '聊天记录管理', permission: 'admin.chat.messages.list' }
            },
            {
                path: 'chat/rooms',
                element: <LazyLoad><ChatRoomsPage /></LazyLoad>,
                meta: { title: '聊天室管理', permission: 'admin.chat.conversations.list' }
            },
            // 实时监控大屏
            {
                path: 'monitor/dashboard',
                element: <LazyLoad><MonitorPage /></LazyLoad>,
                meta: { title: '实时监控大屏', permission: 'admin.monitor.view' }
            },
            // VIP管理模块
            {
                path: 'vip',
                element: <LazyLoad><VIPLevels /></LazyLoad>,
                meta: { title: 'VIP等级管理', permission: 'admin.vip.list' }
            },
            {
                path: 'vip/levels',
                element: <LazyLoad><VIPLevels /></LazyLoad>,
                meta: { title: 'VIP等级管理', permission: 'admin.vip.list' }
            },
            {
                path: 'vip/config',
                element: <LazyLoad><VIPConfig /></LazyLoad>,
                meta: { title: 'VIP系统配置', permission: 'admin.vip.config' }
            },
            // 营销模块
            {
                path: 'coupon',
                element: <LazyLoad><CouponPage /></LazyLoad>,
                meta: { title: '优惠券管理', permission: 'admin.coupons.list' }
            },
            {
                path: 'activity',
                element: <LazyLoad><ActivityPage /></LazyLoad>,
                meta: { title: '活动管理', permission: 'admin.activities.list' }
            },
            {
                path: 'recharge',
                element: <LazyLoad><RechargePage /></LazyLoad>,
                meta: { title: '充值套餐管理', permission: 'admin.recharges.list' }
            },
            {
                path: 'team',
                element: <LazyLoad><TeamPage /></LazyLoad>,
                meta: { title: '战队管理', permission: 'admin.teams.list' }
            },
            {
                path: 'referral',
                element: <LazyLoad><ReferralPage /></LazyLoad>,
                meta: { title: '推荐管理', permission: 'admin.referrals.list' }
            },
            // 财务管理模块
            {
                path: 'finance/withdraw',
                element: <LazyLoad><WithdrawPage /></LazyLoad>,
                meta: { title: '提现管理', permission: 'admin.withdraws.list' }
            },
            {
                path: 'finance/commission',
                element: <LazyLoad><CommissionPage /></LazyLoad>,
                meta: { title: '佣金管理', permission: 'admin.commissions.list' }
            },
            {
                path: 'finance/ranking-commission',
                element: <LazyLoad><RankingCommissionPage /></LazyLoad>,
                meta: { title: '排名佣金配置', permission: 'admin.ranking-commissions.list' }
            },
            // 支付管理模块
            {
                path: 'payment/records',
                element: <LazyLoad><PaymentRecordsPage /></LazyLoad>,
                meta: { title: '支付记录', permission: 'admin.payments.list' }
            },
            // 业务管理模块 - 纠纷和分流
            {
                path: 'biz/dispute',
                element: <LazyLoad><DisputePage /></LazyLoad>,
                meta: { title: '纠纷管理', permission: 'admin.disputes.list' }
            },
            {
                path: 'biz/routing',
                element: <LazyLoad><RoutingPage /></LazyLoad>,
                meta: { title: '分流规则管理', permission: 'admin.routing.list' }
            },
            // 结算公司管理模块
            {
                path: 'settlement',
                element: <LazyLoad><SettlementPage /></LazyLoad>,
                meta: { title: '结算公司管理', permission: 'admin.settlement.list' }
            },
            {
                path: 'settlement/companies',
                element: <LazyLoad><SettlementPage /></LazyLoad>,
                meta: { title: '结算公司管理', permission: 'admin.settlement.list' }
            },
            {
                path: 'settlement/players',
                element: <LazyLoad><SettlementPlayersPage /></LazyLoad>,
                meta: { title: '陪玩师归属管理', permission: 'admin.settlement.players' }
            },
            // 对账管理模块
            {
                path: 'finance/settlement-company',
                element: <LazyLoad><SettlementPage /></LazyLoad>,
                meta: { title: '结算公司管理', permission: 'admin.settlement-companies.list' }
            },
            {
                path: 'finance/reconciliation',
                element: <LazyLoad><ReconciliationPage /></LazyLoad>,
                meta: { title: '对账管理', permission: 'admin.reconciliations.list' }
            },
            // 个人中心
            {
                path: 'profile',
                element: <LazyLoad><ProfilePage /></LazyLoad>,
                meta: { title: '个人中心' }
            },
            // 营销管理（动态菜单路径兼容）
            {
                path: 'marketing/vip',
                element: <LazyLoad><VIPLevels /></LazyLoad>,
                meta: { title: 'VIP等级管理', permission: 'admin.vip.list' }
            },
            {
                path: 'marketing/coupon',
                element: <LazyLoad><CouponPage /></LazyLoad>,
                meta: { title: '优惠券管理', permission: 'admin.coupons.list' }
            },
            {
                path: 'marketing/referral',
                element: <LazyLoad><ReferralPage /></LazyLoad>,
                meta: { title: '推荐管理', permission: 'admin.referrals.list' }
            },
            {
                path: 'marketing/team',
                element: <LazyLoad><TeamPage /></LazyLoad>,
                meta: { title: '战队管理', permission: 'admin.teams.list' }
            },
            {
                path: 'marketing/activity',
                element: <LazyLoad><ActivityPage /></LazyLoad>,
                meta: { title: '活动管理', permission: 'admin.activities.list' }
            },
            // 支付记录（动态菜单路径）
            {
                path: 'payment/payment-records',
                element: <LazyLoad><PaymentRecordsPage /></LazyLoad>,
                meta: { title: '支付记录', permission: 'admin.payments.list' }
            },
            // 告警管理（动态菜单路径）
            {
                path: 'monitor/alert',
                element: <LazyLoad><AdminAlert /></LazyLoad>,
                meta: { title: '告警管理', permission: 'admin.monitor.alerts' }
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
