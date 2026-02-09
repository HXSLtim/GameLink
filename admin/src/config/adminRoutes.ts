/**
 * Admin路由配置
 * 用于同步到后端菜单和权限系统
 */

/**
 * 菜单配置接口
 */
export interface MenuConfig {
    /** 菜单名称 */
    name: string;
    /** 路由路径 */
    path: string;
    /** 组件路径 */
    component: string;
    /** 图标名称 */
    icon?: string;
    /** 排序顺序 */
    order: number;
    /** 是否隐藏 */
    hidden?: boolean;
    /** 权限码 */
    permission?: string;
    /** 描述 */
    description?: string;
    /** 重定向地址 */
    redirect?: string;
    /** 子菜单 */
    children?: MenuConfig[];
}

/**
 * 权限配置接口
 */
export interface PermissionConfig {
    /** HTTP方法 */
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
    /** API路径 */
    path: string;
    /** 权限码 */
    code: string;
    /** 权限分组 */
    group: string;
    /** 描述 */
    description: string;
}

/**
 * Admin菜单配置
 */
export const ADMIN_MENUS: MenuConfig[] = [
    {
        name: '仪表盘',
        path: '/admin',
        component: 'Dashboard',
        icon: 'DashboardOutlined',
        order: 0,
        permission: 'admin.dashboard.view',
        description: '系统概览和数据统计',
    },
    {
        name: '系统管理',
        path: '/admin/sys',
        component: 'Layout',
        icon: 'SettingOutlined',
        order: 1,
        description: '系统配置和权限管理',
        children: [
            {
                name: '用户管理',
                path: '/admin/sys/user',
                component: 'User',
                icon: 'UserOutlined',
                order: 0,
                permission: 'admin.users.list',
                description: '管理平台用户',
            },
            {
                name: '角色管理',
                path: '/admin/sys/role',
                component: 'Role',
                icon: 'TeamOutlined',
                order: 1,
                permission: 'admin.roles.list',
                description: '管理系统角色',
            },
            {
                name: '权限管理',
                path: '/admin/sys/permission',
                component: 'Permission',
                icon: 'SafetyCertificateOutlined',
                order: 2,
                permission: 'admin.permissions.list',
                description: '管理系统权限',
            },
            {
                name: '菜单管理',
                path: '/admin/sys/menu',
                component: 'Menu',
                icon: 'MenuOutlined',
                order: 3,
                permission: 'admin.menus.list',
                description: '管理系统菜单',
            },
            {
                name: '用户标签',
                path: '/admin/sys/user-tag',
                component: 'UserTag',
                icon: 'tags',
                order: 4,
                permission: 'admin.user-tags.list',
                description: '管理用户标签',
            },
            {
                name: '用户角色分配',
                path: '/admin/sys/user-role',
                component: 'UserRole',
                icon: 'UserSwitchOutlined',
                order: 5,
                permission: 'admin.roles.assign-user',
                description: '分配用户角色',
            },
            {
                name: '审计日志',
                path: '/admin/sys/log',
                component: 'Audit',
                icon: 'FileSearchOutlined',
                order: 6,
                permission: 'admin.operation-logs.list',
                description: '查看系统操作日志',
            },
        ],
    },
    {
        name: '业务管理',
        path: '/admin/biz',
        component: 'Layout',
        icon: 'AppstoreOutlined',
        order: 2,
        description: '业务数据管理',
        children: [
            {
                name: '游戏管理',
                path: '/admin/biz/game',
                component: 'Game',
                icon: 'AppstoreOutlined',
                order: 0,
                permission: 'admin.games.list',
                description: '管理平台游戏',
            },
            {
                name: '陪玩师管理',
                path: '/admin/biz/player',
                component: 'Player',
                icon: 'TeamOutlined',
                order: 1,
                permission: 'admin.players.list',
                description: '管理陪玩师',
            },
            {
                name: '段位管理',
                path: '/admin/biz/game-rank',
                component: 'GameRank',
                icon: 'TrophyOutlined',
                order: 2,
                permission: 'admin.game-ranks.list',
                description: '管理游戏段位配置',
            },
            {
                name: '段位审核',
                path: '/admin/biz/player-rank',
                component: 'PlayerRank',
                icon: 'AuditOutlined',
                order: 3,
                permission: 'admin.player-ranks.list',
                description: '审核陪玩师段位认证',
            },
            {
                name: '实名审核',
                path: '/admin/biz/player-certification',
                component: 'AdminPlayerCertification',
                icon: 'IdcardOutlined',
                order: 4,
                permission: 'admin.player-certifications.list',
                description: '审核陪玩师实名认证',
            },
            {
                name: '订单管理',
                path: '/admin/biz/order',
                component: 'Order',
                icon: 'ShoppingCartOutlined',
                order: 5,
                permission: 'admin.orders.list',
                description: '管理平台订单',
            },
            {
                name: '服务项目',
                path: '/admin/biz/service',
                component: 'ServiceItem',
                icon: 'GiftOutlined',
                order: 6,
                permission: 'admin.service-items.list',
                description: '管理服务项目',
            },
            {
                name: '纠纷管理',
                path: '/admin/biz/dispute',
                component: 'Dispute',
                icon: 'warning',
                order: 7,
                permission: 'admin.disputes.list',
                description: '管理订单纠纷',
            },
            {
                name: '分流规则',
                path: '/admin/biz/routing-rule',
                component: 'RoutingRule',
                icon: 'branches',
                order: 8,
                permission: 'admin.routing-rules.list',
                description: '管理订单分流规则',
            },
        ],
    },
    {
        name: '财务管理',
        path: '/admin/finance',
        component: 'Layout',
        icon: 'DollarOutlined',
        order: 3,
        description: '财务和结算管理',
        children: [
            {
                name: '提现管理',
                path: '/admin/finance/withdraw',
                component: 'Withdraw',
                icon: 'WalletOutlined',
                order: 1,
                permission: 'admin.withdraws.list',
                description: '管理陪玩师提现申请',
            },
            {
                name: '佣金设置',
                path: '/admin/finance/commission',
                component: 'Commission',
                icon: 'PercentageOutlined',
                order: 2,
                permission: 'admin.commissions.list',
                description: '平台佣金规则和结算',
            },
            {
                name: '结算公司',
                path: '/admin/finance/settlement-company',
                component: 'SettlementCompany',
                icon: 'bank',
                order: 3,
                permission: 'admin.settlement-companies.list',
                description: '管理结算公司',
            },
            {
                name: '排行榜抽成',
                path: '/admin/finance/ranking-commission',
                component: 'RankingCommission',
                icon: 'trophy',
                order: 4,
                permission: 'admin.ranking-commissions.list',
                description: '管理排行榜抽成规则',
            },
            {
                name: '提现分流',
                path: '/admin/finance/withdraw-routing',
                component: 'WithdrawRouting',
                icon: 'swap',
                order: 5,
                permission: 'admin.withdraw-routing.list',
                description: '管理提现分流规则',
            },
        ],
    },
    {
        name: '监控中心',
        path: '/admin/monitor',
        component: 'Layout',
        icon: 'MonitorOutlined',
        order: 4,
        description: '系统监控和数据分析',
        children: [
            {
                name: '实时监控大屏',
                path: '/admin/monitor/dashboard',
                component: 'Monitor',
                icon: 'DashboardOutlined',
                order: 0,
                permission: 'admin.monitor.view',
                description: '实时系统状态监控大屏',
            },
            {
                name: '实时监控',
                path: '/admin/monitor/realtime',
                component: 'RealtimeMonitor',
                icon: 'MonitorOutlined',
                order: 1,
                description: '实时系统状态监控',
            },
            {
                name: '运营分析',
                path: '/admin/monitor/analytics',
                component: 'Analytics',
                icon: 'LineChartOutlined',
                order: 2,
                description: '运营数据分析',
            },
            {
                name: 'KPI 仪表板',
                path: '/admin/monitor/kpi',
                component: 'KPIDashboard',
                icon: 'FundOutlined',
                order: 3,
                description: 'KPI 指标展示',
            },
            {
                name: '用户行为',
                path: '/admin/monitor/user-behavior',
                component: 'UserBehavior',
                icon: 'line-chart',
                order: 4,
                permission: 'admin.user-behavior.list',
                description: '用户行为分析',
            },
            {
                name: '告警管理',
                path: '/admin/monitor/alert',
                component: 'Alert',
                icon: 'AlertOutlined',
                order: 5,
                permission: 'admin.monitor.alerts',
                description: '查看和管理系统告警',
            },
        ],
    },
    {
        name: '内容管理',
        path: '/admin/content',
        component: 'Layout',
        icon: 'FileTextOutlined',
        order: 5,
        description: '动态审核和聊天监控',
        children: [
            {
                name: '动态审核',
                path: '/admin/content/feeds',
                component: 'ContentFeeds',
                icon: 'FileTextOutlined',
                order: 1,
                permission: 'content.feed.list',
                description: '审核用户发布的动态',
            },
            {
                name: '聊天监控',
                path: '/admin/content/chat',
                component: 'ContentChat',
                icon: 'MessageOutlined',
                order: 2,
                permission: 'content.chat.list',
                description: '监控聊天消息',
            },
            {
                name: '举报管理',
                path: '/admin/content/reports',
                component: 'ContentReports',
                icon: 'WarningOutlined',
                order: 3,
                permission: 'content.report.list',
                description: '处理动态举报',
            },
            {
                name: '内容分类',
                path: '/admin/content/categories',
                component: 'ContentCategories',
                icon: 'TagsOutlined',
                order: 4,
                permission: 'content.category.list',
                description: '管理内容分类',
            },
            {
                name: '内容统计',
                path: '/admin/content/stats',
                component: 'ContentStats',
                icon: 'BarChartOutlined',
                order: 5,
                permission: 'content.stats',
                description: '内容数据统计',
            },
        ],
    },
    {
        name: '聊天管理',
        path: '/admin/chat',
        component: 'Layout',
        icon: 'MessageOutlined',
        order: 6,
        description: '聊天会话和消息管理',
        children: [
            {
                name: '聊天室管理',
                path: '/admin/chat/rooms',
                component: 'ChatRooms',
                icon: 'TeamOutlined',
                order: 1,
                permission: 'admin.chat.conversations.list',
                description: '管理所有聊天会话',
            },
            {
                name: '聊天记录管理',
                path: '/admin/chat/records',
                component: 'ChatRecords',
                icon: 'UnorderedListOutlined',
                order: 2,
                permission: 'admin.chat.messages.list',
                description: '查看和管理聊天消息记录',
            },
        ],
    },
    {
        name: '营销管理',
        path: '/admin/marketing',
        component: 'Layout',
        icon: 'GiftOutlined',
        order: 7,
        description: '营销活动和用户增长',
        children: [
            {
                name: 'VIP 管理',
                path: '/admin/marketing/vip',
                component: 'VIP',
                icon: 'crown',
                order: 1,
                permission: 'admin.vip.list',
                description: '管理 VIP 等级和权益',
            },
            {
                name: '优惠券管理',
                path: '/admin/marketing/coupon',
                component: 'Coupon',
                icon: 'gift',
                order: 2,
                permission: 'admin.coupons.list',
                description: '创建和管理优惠券',
            },
            {
                name: '推荐系统',
                path: '/admin/marketing/referral',
                component: 'Referral',
                icon: 'share-alt',
                order: 3,
                permission: 'admin.referrals.list',
                description: '管理推荐和邀请奖励',
            },
            {
                name: '团队管理',
                path: '/admin/marketing/team',
                component: 'Team',
                icon: 'team',
                order: 4,
                permission: 'admin.teams.list',
                description: '管理陪玩师团队',
            },
            {
                name: '活动管理',
                path: '/admin/marketing/activity',
                component: 'Activity',
                icon: 'calendar',
                order: 5,
                permission: 'admin.activities.list',
                description: '创建和管理营销活动',
            },
        ],
    },
    {
        name: '支付管理',
        path: '/admin/payment',
        component: 'Layout',
        icon: 'PayCircleOutlined',
        order: 8,
        description: '支付和充值管理',
        children: [
            {
                name: '充值管理',
                path: '/admin/payment/recharge',
                component: 'Recharge',
                icon: 'wallet',
                order: 1,
                permission: 'admin.recharges.list',
                description: '管理充值选项和记录',
            },
            {
                name: '支付记录',
                path: '/admin/payment/payment-records',
                component: 'PaymentRecords',
                icon: 'transaction',
                order: 2,
                permission: 'admin.payments.list',
                description: '查看和管理支付记录',
            },
        ],
    },
    {
        name: '评价管理',
        path: '/admin/reviews',
        component: 'Layout',
        icon: 'StarOutlined',
        order: 9,
        description: '评价和举报管理',
        children: [
            {
                name: '评价列表',
                path: '/admin/reviews/list',
                component: 'ReviewList',
                icon: 'UnorderedListOutlined',
                order: 1,
                permission: 'admin.reviews.list',
                description: '查看和管理评价',
            },
            {
                name: '评价审核',
                path: '/admin/reviews/moderation',
                component: 'ReviewModeration',
                icon: 'AuditOutlined',
                order: 2,
                permission: 'admin.reviews.pending.list',
                description: '审核待处理评价',
            },
            {
                name: '举报管理',
                path: '/admin/reviews/reports',
                component: 'ReviewReports',
                icon: 'WarningOutlined',
                order: 3,
                permission: 'admin.review-reports.list',
                description: '处理评价举报',
            },
            {
                name: '敏感词管理',
                path: '/admin/reviews/sensitive-words',
                component: 'SensitiveWords',
                icon: 'StopOutlined',
                order: 4,
                permission: 'admin.sensitive-words.list',
                description: '管理敏感词库',
            },
            {
                name: '评价统计',
                path: '/admin/reviews/stats',
                component: 'ReviewStats',
                icon: 'BarChartOutlined',
                order: 5,
                permission: 'admin.reviews.stats.list',
                description: '评价数据统计',
            },
            {
                name: '评价设置',
                path: '/admin/reviews/settings',
                component: 'ReviewSettings',
                icon: 'SettingOutlined',
                order: 6,
                permission: 'admin.review-settings.list',
                description: '评价系统参数设置',
            },
        ],
    },
    {
        name: '系统设置',
        path: '/admin/settings',
        component: 'Settings',
        icon: 'SettingOutlined',
        order: 99,
        permission: 'admin.settings.view',
        description: '系统参数设置（含评价设置）',
    },
    {
        name: '个人中心',
        path: '/admin/profile',
        component: 'Profile',
        icon: 'UserOutlined',
        order: 100,
        hidden: true, // 不在侧边栏显示，只通过用户下拉菜单访问
        description: '个人信息管理和密码修改',
    },
];

/**
 * Admin权限配置
 * 这些权限会自动同步到后端
 */
export const ADMIN_PERMISSIONS: PermissionConfig[] = [
    // 仪表盘
    { method: 'GET', path: '/api/v1/admin/dashboard', code: 'admin.dashboard.view', group: '/admin/dashboard', description: '查看仪表盘' },
    { method: 'GET', path: '/api/v1/admin/stats', code: 'admin.stats.read', group: '/admin/stats', description: '查看统计数据' },

    // 用户管理
    { method: 'GET', path: '/api/v1/admin/users', code: 'admin.users.list', group: '/admin/users', description: '获取用户列表' },
    { method: 'GET', path: '/api/v1/admin/users/:id', code: 'admin.users.read', group: '/admin/users', description: '获取用户详情' },
    { method: 'POST', path: '/api/v1/admin/users', code: 'admin.users.create', group: '/admin/users', description: '创建用户' },
    { method: 'PUT', path: '/api/v1/admin/users/:id', code: 'admin.users.update', group: '/admin/users', description: '更新用户' },
    { method: 'DELETE', path: '/api/v1/admin/users/:id', code: 'admin.users.delete', group: '/admin/users', description: '删除用户' },
    { method: 'PATCH', path: '/api/v1/admin/users/:id/status', code: 'admin.users.status', group: '/admin/users', description: '更新用户状态' },

    // 角色管理
    { method: 'GET', path: '/api/v1/admin/roles', code: 'admin.roles.list', group: '/admin/roles', description: '获取角色列表' },
    { method: 'GET', path: '/api/v1/admin/roles/:id', code: 'admin.roles.read', group: '/admin/roles', description: '获取角色详情' },
    { method: 'POST', path: '/api/v1/admin/roles', code: 'admin.roles.create', group: '/admin/roles', description: '创建角色' },
    { method: 'PUT', path: '/api/v1/admin/roles/:id', code: 'admin.roles.update', group: '/admin/roles', description: '更新角色' },
    { method: 'DELETE', path: '/api/v1/admin/roles/:id', code: 'admin.roles.delete', group: '/admin/roles', description: '删除角色' },
    { method: 'PUT', path: '/api/v1/admin/roles/:id/permissions', code: 'admin.roles.permissions', group: '/admin/roles', description: '分配角色权限' },

    // 权限管理
    { method: 'GET', path: '/api/v1/admin/permissions', code: 'admin.permissions.list', group: '/admin/permissions', description: '获取权限列表' },
    { method: 'GET', path: '/api/v1/admin/permissions/:id', code: 'admin.permissions.read', group: '/admin/permissions', description: '获取权限详情' },
    { method: 'POST', path: '/api/v1/admin/permissions', code: 'admin.permissions.create', group: '/admin/permissions', description: '创建权限' },
    { method: 'PUT', path: '/api/v1/admin/permissions/:id', code: 'admin.permissions.update', group: '/admin/permissions', description: '更新权限' },
    { method: 'DELETE', path: '/api/v1/admin/permissions/:id', code: 'admin.permissions.delete', group: '/admin/permissions', description: '删除权限' },

    // 菜单管理
    { method: 'GET', path: '/api/v1/admin/menus', code: 'admin.menus.list', group: '/admin/menus', description: '获取菜单列表' },
    { method: 'GET', path: '/api/v1/admin/menus/:id', code: 'admin.menus.read', group: '/admin/menus', description: '获取菜单详情' },
    { method: 'POST', path: '/api/v1/admin/menus', code: 'admin.menus.create', group: '/admin/menus', description: '创建菜单' },
    { method: 'PUT', path: '/api/v1/admin/menus/:id', code: 'admin.menus.update', group: '/admin/menus', description: '更新菜单' },
    { method: 'DELETE', path: '/api/v1/admin/menus/:id', code: 'admin.menus.delete', group: '/admin/menus', description: '删除菜单' },

    // 游戏管理
    { method: 'GET', path: '/api/v1/admin/games', code: 'admin.games.list', group: '/admin/games', description: '获取游戏列表' },
    { method: 'GET', path: '/api/v1/admin/games/:id', code: 'admin.games.read', group: '/admin/games', description: '获取游戏详情' },
    { method: 'POST', path: '/api/v1/admin/games', code: 'admin.games.create', group: '/admin/games', description: '创建游戏' },
    { method: 'PUT', path: '/api/v1/admin/games/:id', code: 'admin.games.update', group: '/admin/games', description: '更新游戏' },
    { method: 'DELETE', path: '/api/v1/admin/games/:id', code: 'admin.games.delete', group: '/admin/games', description: '删除游戏' },

    // 陪玩师管理
    { method: 'GET', path: '/api/v1/admin/players', code: 'admin.players.list', group: '/admin/players', description: '获取陪玩师列表' },
    { method: 'GET', path: '/api/v1/admin/players/:id', code: 'admin.players.read', group: '/admin/players', description: '获取陪玩师详情' },
    { method: 'PUT', path: '/api/v1/admin/players/:id', code: 'admin.players.update', group: '/admin/players', description: '更新陪玩师' },
    { method: 'DELETE', path: '/api/v1/admin/players/:id', code: 'admin.players.delete', group: '/admin/players', description: '删除陪玩师' },
    { method: 'POST', path: '/api/v1/admin/players/:id/audit', code: 'admin.players.audit', group: '/admin/players', description: '审核陪玩师' },

    // 段位管理
    { method: 'GET', path: '/api/v1/admin/game-ranks', code: 'admin.game-ranks.list', group: '/admin/game-ranks', description: '获取段位列表' },
    { method: 'GET', path: '/api/v1/admin/game-ranks/:id', code: 'admin.game-ranks.read', group: '/admin/game-ranks', description: '获取段位详情' },
    { method: 'POST', path: '/api/v1/admin/game-ranks', code: 'admin.game-ranks.create', group: '/admin/game-ranks', description: '创建段位' },
    { method: 'PUT', path: '/api/v1/admin/game-ranks/:id', code: 'admin.game-ranks.update', group: '/admin/game-ranks', description: '更新段位' },
    { method: 'DELETE', path: '/api/v1/admin/game-ranks/:id', code: 'admin.game-ranks.delete', group: '/admin/game-ranks', description: '删除段位' },

    // 段位审核
    { method: 'GET', path: '/api/v1/admin/player-ranks', code: 'admin.player-ranks.list', group: '/admin/player-ranks', description: '获取段位认证列表' },
    { method: 'GET', path: '/api/v1/admin/player-ranks/:id', code: 'admin.player-ranks.read', group: '/admin/player-ranks', description: '获取段位认证详情' },
    { method: 'PUT', path: '/api/v1/admin/player-ranks/:id/verify', code: 'admin.player-ranks.verify', group: '/admin/player-ranks', description: '审核段位认证' },

    // 实名审核
    { method: 'GET', path: '/api/v1/admin/player-certifications', code: 'admin.player-certifications.list', group: '/admin/player-certifications', description: '获取实名认证列表' },
    { method: 'GET', path: '/api/v1/admin/player-certifications/:id', code: 'admin.player-certifications.read', group: '/admin/player-certifications', description: '获取实名认证详情' },
    { method: 'PUT', path: '/api/v1/admin/player-certifications/:id/verify', code: 'admin.player-certifications.verify', group: '/admin/player-certifications', description: '审核实名认证' },

    // 订单管理
    { method: 'GET', path: '/api/v1/admin/orders', code: 'admin.orders.list', group: '/admin/orders', description: '获取订单列表' },
    { method: 'GET', path: '/api/v1/admin/orders/:id', code: 'admin.orders.read', group: '/admin/orders', description: '获取订单详情' },
    { method: 'PUT', path: '/api/v1/admin/orders/:id', code: 'admin.orders.update', group: '/admin/orders', description: '更新订单' },
    { method: 'POST', path: '/api/v1/admin/orders/:id/cancel', code: 'admin.orders.cancel', group: '/admin/orders', description: '取消订单' },
    { method: 'POST', path: '/api/v1/admin/orders/:id/refund', code: 'admin.orders.refund', group: '/admin/orders', description: '订单退款' },

    // 服务项目管理
    { method: 'GET', path: '/api/v1/admin/service-items', code: 'admin.service-items.list', group: '/admin/service-items', description: '获取服务项目列表' },
    { method: 'GET', path: '/api/v1/admin/service-items/:id', code: 'admin.service-items.read', group: '/admin/service-items', description: '获取服务项目详情' },
    { method: 'POST', path: '/api/v1/admin/service-items', code: 'admin.service-items.create', group: '/admin/service-items', description: '创建服务项目' },
    { method: 'PUT', path: '/api/v1/admin/service-items/:id', code: 'admin.service-items.update', group: '/admin/service-items', description: '更新服务项目' },
    { method: 'DELETE', path: '/api/v1/admin/service-items/:id', code: 'admin.service-items.delete', group: '/admin/service-items', description: '删除服务项目' },

    // 提现管理
    { method: 'GET', path: '/api/v1/admin/withdraws', code: 'admin.withdraws.list', group: '/admin/withdraws', description: '获取提现列表' },
    { method: 'GET', path: '/api/v1/admin/withdraws/:id', code: 'admin.withdraws.read', group: '/admin/withdraws', description: '获取提现详情' },
    { method: 'POST', path: '/api/v1/admin/withdraws/:id/approve', code: 'admin.withdraws.approve', group: '/admin/withdraws', description: '批准提现' },
    { method: 'POST', path: '/api/v1/admin/withdraws/:id/reject', code: 'admin.withdraws.reject', group: '/admin/withdraws', description: '拒绝提现' },
    { method: 'POST', path: '/api/v1/admin/withdraws/:id/complete', code: 'admin.withdraws.complete', group: '/admin/withdraws', description: '完成提现打款' },

    // 佣金管理
    { method: 'GET', path: '/api/v1/admin/commission/stats', code: 'admin.commissions.list', group: '/admin/commission', description: '获取平台统计' },
    { method: 'POST', path: '/api/v1/admin/commission/rules', code: 'admin.commissions.create', group: '/admin/commission', description: '创建佣金规则' },
    { method: 'PUT', path: '/api/v1/admin/commission/rules/:id', code: 'admin.commissions.update', group: '/admin/commission', description: '更新佣金规则' },
    { method: 'POST', path: '/api/v1/admin/commission/settlements/trigger', code: 'admin.commissions.settle', group: '/admin/commission', description: '触发月度结算' },

    // 系统设置
    { method: 'GET', path: '/api/v1/admin/settings', code: 'admin.settings.view', group: '/admin/settings', description: '查看系统设置' },
    { method: 'PUT', path: '/api/v1/admin/settings', code: 'admin.settings.update', group: '/admin/settings', description: '更新系统设置' },

    // 监控中心
    { method: 'GET', path: '/api/v1/admin/monitor/system-status', code: 'admin.monitor.system', group: '/admin/monitor', description: '查看系统状态' },
    { method: 'GET', path: '/api/v1/admin/monitor/online-users', code: 'admin.monitor.online', group: '/admin/monitor', description: '查看在线用户' },
    { method: 'GET', path: '/api/v1/admin/monitor/order-queue', code: 'admin.monitor.orders', group: '/admin/monitor', description: '查看订单队列' },
    { method: 'GET', path: '/api/v1/admin/monitor/alerts', code: 'admin.monitor.alerts', group: '/admin/monitor', description: '查看告警' },
    { method: 'GET', path: '/api/v1/admin/analytics/active-users', code: 'admin.analytics.users', group: '/admin/analytics', description: '查看活跃用户' },
    { method: 'GET', path: '/api/v1/admin/analytics/retention', code: 'admin.analytics.retention', group: '/admin/analytics', description: '查看留存率' },
    { method: 'GET', path: '/api/v1/admin/analytics/payment', code: 'admin.analytics.payment', group: '/admin/analytics', description: '查看支付分析' },
    { method: 'GET', path: '/api/v1/admin/analytics/conversion', code: 'admin.analytics.conversion', group: '/admin/analytics', description: '查看转化漏斗' },
    { method: 'GET', path: '/api/v1/admin/kpi/overview', code: 'admin.kpi.overview', group: '/admin/kpi', description: '查看KPI概览' },
    { method: 'GET', path: '/api/v1/admin/kpi/trend', code: 'admin.kpi.trend', group: '/admin/kpi', description: '查看KPI趋势' },
    { method: 'GET', path: '/api/v1/admin/kpi/targets', code: 'admin.kpi.targets', group: '/admin/kpi', description: '查看KPI目标' },

    // 内容管理 - 动态审核
    { method: 'GET', path: '/api/v1/admin/content/feeds', code: 'content.feed.list', group: '/admin/content', description: '获取动态列表' },
    { method: 'GET', path: '/api/v1/admin/content/feeds/:id', code: 'content.feed.get', group: '/admin/content', description: '获取动态详情' },
    { method: 'PUT', path: '/api/v1/admin/content/feeds/:id/approve', code: 'content.feed.approve', group: '/admin/content', description: '批准动态' },
    { method: 'PUT', path: '/api/v1/admin/content/feeds/:id/reject', code: 'content.feed.reject', group: '/admin/content', description: '拒绝动态' },
    { method: 'DELETE', path: '/api/v1/admin/content/feeds/:id', code: 'content.feed.delete', group: '/admin/content', description: '删除动态' },
    { method: 'POST', path: '/api/v1/admin/content/feeds/batch-approve', code: 'content.feed.batch_approve', group: '/admin/content', description: '批量批准动态' },
    { method: 'POST', path: '/api/v1/admin/content/feeds/batch-reject', code: 'content.feed.batch_reject', group: '/admin/content', description: '批量拒绝动态' },

    // 内容管理 - 聊天监控
    { method: 'GET', path: '/api/v1/admin/content/chat/messages', code: 'content.chat.list', group: '/admin/content', description: '获取聊天消息列表' },
    { method: 'DELETE', path: '/api/v1/admin/content/chat/messages/:id', code: 'content.chat.delete', group: '/admin/content', description: '删除聊天消息' },
    { method: 'POST', path: '/api/v1/admin/content/chat/mute', code: 'content.chat.mute', group: '/admin/content', description: '禁言用户' },
    { method: 'POST', path: '/api/v1/admin/content/chat/unmute', code: 'content.chat.unmute', group: '/admin/content', description: '解除禁言' },

    // 内容管理 - 举报管理
    { method: 'GET', path: '/api/v1/admin/content/reports', code: 'content.report.list', group: '/admin/content', description: '获取举报列表' },
    { method: 'GET', path: '/api/v1/admin/content/reports/:id', code: 'content.report.get', group: '/admin/content', description: '获取举报详情' },
    { method: 'POST', path: '/api/v1/admin/content/reports/:id/process', code: 'content.report.process', group: '/admin/content', description: '处理举报' },

    // 内容管理 - 分类管理
    { method: 'GET', path: '/api/v1/admin/content/categories', code: 'content.category.list', group: '/admin/content', description: '获取分类列表' },
    { method: 'GET', path: '/api/v1/admin/content/categories/:id', code: 'content.category.get', group: '/admin/content', description: '获取分类详情' },
    { method: 'POST', path: '/api/v1/admin/content/categories', code: 'content.category.create', group: '/admin/content', description: '创建分类' },
    { method: 'PUT', path: '/api/v1/admin/content/categories/:id', code: 'content.category.update', group: '/admin/content', description: '更新分类' },
    { method: 'DELETE', path: '/api/v1/admin/content/categories/:id', code: 'content.category.delete', group: '/admin/content', description: '删除分类' },

    // 内容管理 - 统计
    { method: 'GET', path: '/api/v1/admin/content/stats', code: 'content.stats', group: '/admin/content', description: '获取内容统计' },

    // 评价管理 - 权限码格式与后端自动生成的一致: admin.{resource}.{action}
    { method: 'GET', path: '/api/v1/admin/reviews', code: 'admin.reviews.list', group: '/admin/reviews', description: '获取评价列表' },
    { method: 'GET', path: '/api/v1/admin/reviews/:id', code: 'admin.reviews.read', group: '/admin/reviews', description: '获取评价详情' },
    { method: 'GET', path: '/api/v1/admin/reviews/pending', code: 'admin.reviews.pending.list', group: '/admin/reviews', description: '获取待审核评价' },
    { method: 'PUT', path: '/api/v1/admin/reviews/:id/approve', code: 'admin.reviews.approve.update', group: '/admin/reviews', description: '批准评价' },
    { method: 'PUT', path: '/api/v1/admin/reviews/:id/reject', code: 'admin.reviews.reject.update', group: '/admin/reviews', description: '拒绝评价' },
    { method: 'PUT', path: '/api/v1/admin/reviews/batch-approve', code: 'admin.reviews.batch-approve.update', group: '/admin/reviews', description: '批量批准评价' },
    { method: 'PUT', path: '/api/v1/admin/reviews/batch-reject', code: 'admin.reviews.batch-reject.update', group: '/admin/reviews', description: '批量拒绝评价' },
    { method: 'DELETE', path: '/api/v1/admin/reviews/:id', code: 'admin.reviews.delete', group: '/admin/reviews', description: '删除评价' },
    { method: 'PUT', path: '/api/v1/admin/reviews/:id', code: 'admin.reviews.update', group: '/admin/reviews', description: '更新评价' },
    { method: 'GET', path: '/api/v1/admin/reviews/:id/logs', code: 'admin.reviews.logs.list', group: '/admin/reviews', description: '查看评价操作日志' },

    // 评价举报管理
    { method: 'GET', path: '/api/v1/admin/review-reports', code: 'admin.review-reports.list', group: '/admin/review-reports', description: '获取举报列表' },
    { method: 'GET', path: '/api/v1/admin/review-reports/:id', code: 'admin.review-reports.read', group: '/admin/review-reports', description: '获取举报详情' },
    { method: 'PUT', path: '/api/v1/admin/review-reports/:id/handle', code: 'admin.review-reports.handle.update', group: '/admin/review-reports', description: '处理举报' },

    // 敏感词管理
    { method: 'GET', path: '/api/v1/admin/sensitive-words', code: 'admin.sensitive-words.list', group: '/admin/sensitive-words', description: '获取敏感词列表' },
    { method: 'POST', path: '/api/v1/admin/sensitive-words', code: 'admin.sensitive-words.create', group: '/admin/sensitive-words', description: '添加敏感词' },
    { method: 'PUT', path: '/api/v1/admin/sensitive-words/:id', code: 'admin.sensitive-words.update', group: '/admin/sensitive-words', description: '更新敏感词' },
    { method: 'DELETE', path: '/api/v1/admin/sensitive-words/:id', code: 'admin.sensitive-words.delete', group: '/admin/sensitive-words', description: '删除敏感词' },

    // 评价统计
    { method: 'GET', path: '/api/v1/admin/reviews/stats', code: 'admin.reviews.stats.list', group: '/admin/reviews', description: '获取评价统计' },

    // 评价设置
    { method: 'GET', path: '/api/v1/admin/review-settings', code: 'admin.review-settings.list', group: '/admin/review-settings', description: '查看评价设置' },
    { method: 'PUT', path: '/api/v1/admin/review-settings', code: 'admin.review-settings.update', group: '/admin/review-settings', description: '更新评价设置' },

    // 用户标签管理
    { method: 'GET', path: '/api/v1/admin/user-tags', code: 'admin.user-tags.list', group: '/admin/user-tags', description: '获取用户标签列表' },
    { method: 'GET', path: '/api/v1/admin/user-tags/:id', code: 'admin.user-tags.read', group: '/admin/user-tags', description: '获取用户标签详情' },
    { method: 'POST', path: '/api/v1/admin/user-tags', code: 'admin.user-tags.create', group: '/admin/user-tags', description: '创建用户标签' },
    { method: 'PUT', path: '/api/v1/admin/user-tags/:id', code: 'admin.user-tags.update', group: '/admin/user-tags', description: '更新用户标签' },
    { method: 'DELETE', path: '/api/v1/admin/user-tags/:id', code: 'admin.user-tags.delete', group: '/admin/user-tags', description: '删除用户标签' },

    // 纠纷管理
    { method: 'GET', path: '/api/v1/admin/disputes', code: 'admin.disputes.list', group: '/admin/disputes', description: '获取纠纷列表' },
    { method: 'GET', path: '/api/v1/admin/disputes/:id', code: 'admin.disputes.read', group: '/admin/disputes', description: '获取纠纷详情' },
    { method: 'PUT', path: '/api/v1/admin/disputes/:id', code: 'admin.disputes.update', group: '/admin/disputes', description: '更新纠纷' },
    { method: 'POST', path: '/api/v1/admin/disputes/:id/resolve', code: 'admin.disputes.resolve', group: '/admin/disputes', description: '解决纠纷' },

    // 分流规则管理
    { method: 'GET', path: '/api/v1/admin/routing-rules', code: 'admin.routing-rules.list', group: '/admin/routing-rules', description: '获取分流规则列表' },
    { method: 'GET', path: '/api/v1/admin/routing-rules/:id', code: 'admin.routing-rules.read', group: '/admin/routing-rules', description: '获取分流规则详情' },
    { method: 'POST', path: '/api/v1/admin/routing-rules', code: 'admin.routing-rules.create', group: '/admin/routing-rules', description: '创建分流规则' },
    { method: 'PUT', path: '/api/v1/admin/routing-rules/:id', code: 'admin.routing-rules.update', group: '/admin/routing-rules', description: '更新分流规则' },
    { method: 'DELETE', path: '/api/v1/admin/routing-rules/:id', code: 'admin.routing-rules.delete', group: '/admin/routing-rules', description: '删除分流规则' },

    // 结算公司管理
    { method: 'GET', path: '/api/v1/admin/settlement-companies', code: 'admin.settlement-companies.list', group: '/admin/settlement-companies', description: '获取结算公司列表' },
    { method: 'GET', path: '/api/v1/admin/settlement-companies/:id', code: 'admin.settlement-companies.read', group: '/admin/settlement-companies', description: '获取结算公司详情' },
    { method: 'POST', path: '/api/v1/admin/settlement-companies', code: 'admin.settlement-companies.create', group: '/admin/settlement-companies', description: '创建结算公司' },
    { method: 'PUT', path: '/api/v1/admin/settlement-companies/:id', code: 'admin.settlement-companies.update', group: '/admin/settlement-companies', description: '更新结算公司' },
    { method: 'DELETE', path: '/api/v1/admin/settlement-companies/:id', code: 'admin.settlement-companies.delete', group: '/admin/settlement-companies', description: '删除结算公司' },

    // 排行榜抽成管理
    { method: 'GET', path: '/api/v1/admin/ranking-commissions', code: 'admin.ranking-commissions.list', group: '/admin/ranking-commissions', description: '获取排行榜抽成列表' },
    { method: 'GET', path: '/api/v1/admin/ranking-commissions/:id', code: 'admin.ranking-commissions.read', group: '/admin/ranking-commissions', description: '获取排行榜抽成详情' },
    { method: 'POST', path: '/api/v1/admin/ranking-commissions', code: 'admin.ranking-commissions.create', group: '/admin/ranking-commissions', description: '创建排行榜抽成规则' },
    { method: 'PUT', path: '/api/v1/admin/ranking-commissions/:id', code: 'admin.ranking-commissions.update', group: '/admin/ranking-commissions', description: '更新排行榜抽成规则' },
    { method: 'DELETE', path: '/api/v1/admin/ranking-commissions/:id', code: 'admin.ranking-commissions.delete', group: '/admin/ranking-commissions', description: '删除排行榜抽成规则' },

    // 提现分流管理
    { method: 'GET', path: '/api/v1/admin/withdraw-routing', code: 'admin.withdraw-routing.list', group: '/admin/withdraw-routing', description: '获取提现分流规则列表' },
    { method: 'GET', path: '/api/v1/admin/withdraw-routing/:id', code: 'admin.withdraw-routing.read', group: '/admin/withdraw-routing', description: '获取提现分流规则详情' },
    { method: 'POST', path: '/api/v1/admin/withdraw-routing', code: 'admin.withdraw-routing.create', group: '/admin/withdraw-routing', description: '创建提现分流规则' },
    { method: 'PUT', path: '/api/v1/admin/withdraw-routing/:id', code: 'admin.withdraw-routing.update', group: '/admin/withdraw-routing', description: '更新提现分流规则' },
    { method: 'DELETE', path: '/api/v1/admin/withdraw-routing/:id', code: 'admin.withdraw-routing.delete', group: '/admin/withdraw-routing', description: '删除提现分流规则' },

    // 用户行为分析
    { method: 'GET', path: '/api/v1/admin/user-behavior', code: 'admin.user-behavior.list', group: '/admin/user-behavior', description: '获取用户行为数据' },
    { method: 'GET', path: '/api/v1/admin/user-behavior/stats', code: 'admin.user-behavior.stats', group: '/admin/user-behavior', description: '获取用户行为统计' },

    // ==================== 营销管理 ====================

    // VIP 管理
    { method: 'GET', path: '/api/v1/admin/vip/levels', code: 'admin.vip.list', group: '/admin/marketing/vip', description: '获取 VIP 等级列表' },
    { method: 'GET', path: '/api/v1/admin/vip/levels/:id', code: 'admin.vip.read', group: '/admin/marketing/vip', description: '获取 VIP 等级详情' },
    { method: 'POST', path: '/api/v1/admin/vip/levels', code: 'admin.vip.create', group: '/admin/marketing/vip', description: '创建 VIP 等级' },
    { method: 'PUT', path: '/api/v1/admin/vip/levels/:id', code: 'admin.vip.update', group: '/admin/marketing/vip', description: '更新 VIP 等级' },
    { method: 'DELETE', path: '/api/v1/admin/vip/levels/:id', code: 'admin.vip.delete', group: '/admin/marketing/vip', description: '删除 VIP 等级' },
    { method: 'GET', path: '/api/v1/admin/vip/configs', code: 'admin.vip.config.list', group: '/admin/marketing/vip', description: '获取 VIP 配置' },
    { method: 'PUT', path: '/api/v1/admin/vip/configs/:id', code: 'admin.vip.config.update', group: '/admin/marketing/vip', description: '更新 VIP 配置' },

    // 优惠券管理
    { method: 'GET', path: '/api/v1/admin/coupons', code: 'admin.coupons.list', group: '/admin/marketing/coupon', description: '获取优惠券列表' },
    { method: 'GET', path: '/api/v1/admin/coupons/:id', code: 'admin.coupons.read', group: '/admin/marketing/coupon', description: '获取优惠券详情' },
    { method: 'POST', path: '/api/v1/admin/coupons', code: 'admin.coupons.create', group: '/admin/marketing/coupon', description: '创建优惠券' },
    { method: 'PUT', path: '/api/v1/admin/coupons/:id', code: 'admin.coupons.update', group: '/admin/marketing/coupon', description: '更新优惠券' },
    { method: 'DELETE', path: '/api/v1/admin/coupons/:id', code: 'admin.coupons.delete', group: '/admin/marketing/coupon', description: '删除优惠券' },
    { method: 'POST', path: '/api/v1/admin/coupons/:id/toggle', code: 'admin.coupons.toggle', group: '/admin/marketing/coupon', description: '启用/禁用优惠券' },
    { method: 'GET', path: '/api/v1/admin/coupons/:id/usage', code: 'admin.coupons.usage', group: '/admin/marketing/coupon', description: '获取优惠券使用情况' },
    { method: 'POST', path: '/api/v1/admin/coupons/:id/issue', code: 'admin.coupons.issue', group: '/admin/marketing/coupon', description: '发放优惠券' },
    { method: 'GET', path: '/api/v1/admin/user-coupons', code: 'admin.user-coupons.list', group: '/admin/marketing/coupon', description: '获取用户优惠券列表' },
    { method: 'POST', path: '/api/v1/admin/coupons/batch/status', code: 'admin.coupons.batch-status', group: '/admin/marketing/coupon', description: '批量更新优惠券状态' },
    { method: 'POST', path: '/api/v1/admin/coupons/batch/delete', code: 'admin.coupons.batch-delete', group: '/admin/marketing/coupon', description: '批量删除优惠券' },

    // 推荐系统
    { method: 'GET', path: '/api/v1/admin/referrals/configs', code: 'admin.referrals.config', group: '/admin/marketing/referral', description: '获取推荐配置' },
    { method: 'PUT', path: '/api/v1/admin/referrals/configs', code: 'admin.referrals.config.update', group: '/admin/marketing/referral', description: '更新推荐配置' },
    { method: 'GET', path: '/api/v1/admin/referrals/codes', code: 'admin.referral-codes.list', group: '/admin/marketing/referral', description: '获取邀请码列表' },
    { method: 'POST', path: '/api/v1/admin/referrals/codes', code: 'admin.referral-codes.create', group: '/admin/marketing/referral', description: '创建邀请码' },
    { method: 'PUT', path: '/api/v1/admin/referrals/codes/:id', code: 'admin.referral-codes.update', group: '/admin/marketing/referral', description: '更新邀请码' },
    { method: 'DELETE', path: '/api/v1/admin/referrals/codes/:id', code: 'admin.referral-codes.delete', group: '/admin/marketing/referral', description: '删除邀请码' },
    { method: 'GET', path: '/api/v1/admin/referrals', code: 'admin.referrals.list', group: '/admin/marketing/referral', description: '获取推荐列表' },
    { method: 'GET', path: '/api/v1/admin/referrals/:id', code: 'admin.referrals.read', group: '/admin/marketing/referral', description: '获取推荐详情' },
    { method: 'PUT', path: '/api/v1/admin/referrals/:id/status', code: 'admin.referrals.status', group: '/admin/marketing/referral', description: '更新推荐状态' },
    { method: 'POST', path: '/api/v1/admin/referrals/rewards/:id/issue', code: 'admin.referral-rewards.issue', group: '/admin/marketing/referral', description: '发放推荐奖励' },
    { method: 'GET', path: '/api/v1/admin/referrals/stats', code: 'admin.referrals.stats', group: '/admin/marketing/referral', description: '获取推荐统计' },
    { method: 'POST', path: '/api/v1/admin/referrals/codes/batch/status', code: 'admin.referral-codes.batch-status', group: '/admin/marketing/referral', description: '批量更新邀请码状态' },
    { method: 'DELETE', path: '/api/v1/admin/referrals/codes/batch', code: 'admin.referral-codes.batch-delete', group: '/admin/marketing/referral', description: '批量删除邀请码' },

    // 团队管理
    { method: 'GET', path: '/api/v1/admin/teams', code: 'admin.teams.list', group: '/admin/marketing/team', description: '获取团队列表' },
    { method: 'GET', path: '/api/v1/admin/teams/:id', code: 'admin.teams.read', group: '/admin/marketing/team', description: '获取团队详情' },
    { method: 'POST', path: '/api/v1/admin/teams', code: 'admin.teams.create', group: '/admin/marketing/team', description: '创建团队' },
    { method: 'PUT', path: '/api/v1/admin/teams/:id', code: 'admin.teams.update', group: '/admin/marketing/team', description: '更新团队' },
    { method: 'DELETE', path: '/api/v1/admin/teams/:id', code: 'admin.teams.delete', group: '/admin/marketing/team', description: '删除团队' },
    { method: 'PUT', path: '/api/v1/admin/teams/:id/status', code: 'admin.teams.status', group: '/admin/marketing/team', description: '更新团队状态' },
    { method: 'GET', path: '/api/v1/admin/teams/:id/members', code: 'admin.teams.members', group: '/admin/marketing/team', description: '获取团队成员' },
    { method: 'POST', path: '/api/v1/admin/teams/:id/members', code: 'admin.teams.members.add', group: '/admin/marketing/team', description: '添加团队成员' },
    { method: 'DELETE', path: '/api/v1/admin/teams/:id/members/:userId', code: 'admin.teams.members.remove', group: '/admin/marketing/team', description: '移除团队成员' },
    { method: 'PUT', path: '/api/v1/admin/teams/:id/captain', code: 'admin.teams.captain', group: '/admin/marketing/team', description: '转让队长' },
    { method: 'GET', path: '/api/v1/admin/teams/stats', code: 'admin.teams.stats', group: '/admin/marketing/team', description: '获取团队统计' },
    { method: 'POST', path: '/api/v1/admin/teams/batch/status', code: 'admin.teams.batch-status', group: '/admin/marketing/team', description: '批量更新团队状态' },
    { method: 'POST', path: '/api/v1/admin/teams/batch/delete', code: 'admin.teams.batch-delete', group: '/admin/marketing/team', description: '批量删除团队' },

    // 活动管理
    { method: 'GET', path: '/api/v1/admin/activities', code: 'admin.activities.list', group: '/admin/marketing/activity', description: '获取活动列表' },
    { method: 'GET', path: '/api/v1/admin/activities/:id', code: 'admin.activities.read', group: '/admin/marketing/activity', description: '获取活动详情' },
    { method: 'POST', path: '/api/v1/admin/activities', code: 'admin.activities.create', group: '/admin/marketing/activity', description: '创建活动' },
    { method: 'PUT', path: '/api/v1/admin/activities/:id', code: 'admin.activities.update', group: '/admin/marketing/activity', description: '更新活动' },
    { method: 'DELETE', path: '/api/v1/admin/activities/:id', code: 'admin.activities.delete', group: '/admin/marketing/activity', description: '删除活动' },
    { method: 'POST', path: '/api/v1/admin/activities/:id/publish', code: 'admin.activities.publish', group: '/admin/marketing/activity', description: '发布活动' },
    { method: 'POST', path: '/api/v1/admin/activities/:id/unpublish', code: 'admin.activities.unpublish', group: '/admin/marketing/activity', description: '下架活动' },
    { method: 'GET', path: '/api/v1/admin/activities/:id/rewards', code: 'admin.activity-rewards.list', group: '/admin/marketing/activity', description: '获取活动奖励' },

    // ==================== 支付管理 ====================

    // 充值管理
    { method: 'GET', path: '/api/v1/admin/recharges/options', code: 'admin.recharges.list', group: '/admin/payment/recharge', description: '获取充值选项列表' },
    { method: 'GET', path: '/api/v1/admin/recharges/options/:id', code: 'admin.recharges.read', group: '/admin/payment/recharge', description: '获取充值选项详情' },
    { method: 'POST', path: '/api/v1/admin/recharges/options', code: 'admin.recharges.create', group: '/admin/payment/recharge', description: '创建充值选项' },
    { method: 'PUT', path: '/api/v1/admin/recharges/options/:id', code: 'admin.recharges.update', group: '/admin/payment/recharge', description: '更新充值选项' },
    { method: 'DELETE', path: '/api/v1/admin/recharges/options/:id', code: 'admin.recharges.delete', group: '/admin/payment/recharge', description: '删除充值选项' },
    { method: 'POST', path: '/api/v1/admin/recharges/options/:id/toggle', code: 'admin.recharges.toggle', group: '/admin/payment/recharge', description: '启用/禁用充值选项' },
    { method: 'GET', path: '/api/v1/admin/recharges/records', code: 'admin.recharge-records.list', group: '/admin/payment/recharge', description: '获取充值记录' },
    { method: 'GET', path: '/api/v1/admin/recharges/records/:id', code: 'admin.recharge-records.read', group: '/admin/payment/recharge', description: '获取充值记录详情' },
    { method: 'POST', path: '/api/v1/admin/recharges/records/:id/refund', code: 'admin.recharge-records.refund', group: '/admin/payment/recharge', description: '充值退款' },

    // 支付记录
    { method: 'GET', path: '/api/v1/admin/payments', code: 'admin.payments.list', group: '/admin/payment', description: '获取支付记录列表' },
    { method: 'GET', path: '/api/v1/admin/payments/:id', code: 'admin.payments.read', group: '/admin/payment', description: '获取支付记录详情' },
    { method: 'POST', path: '/api/v1/admin/payments/:id/refund', code: 'admin.payments.refund', group: '/admin/payment', description: '发起支付退款' },
    { method: 'GET', path: '/api/v1/admin/payments/stats', code: 'admin.payments.stats', group: '/admin/payment', description: '获取支付统计' },

    // ==================== 聊天管理 ====================

    // 聊天会话管理
    { method: 'GET', path: '/api/v1/admin/chats/conversations', code: 'admin.chats.list', group: '/admin/chats', description: '获取聊天会话列表' },
    { method: 'GET', path: '/api/v1/admin/chats/conversations/:id', code: 'admin.chats.read', group: '/admin/chats', description: '获取会话详情' },
    { method: 'POST', path: '/api/v1/admin/chats/conversations/:id/close', code: 'admin.chats.close', group: '/admin/chats', description: '关闭会话' },
    { method: 'GET', path: '/api/v1/admin/chats/messages', code: 'admin.chat-messages.list', group: '/admin/chats', description: '获取聊天消息' },
    { method: 'DELETE', path: '/api/v1/admin/chats/messages/:id', code: 'admin.chat-messages.delete', group: '/admin/chats', description: '删除聊天消息' },
    { method: 'POST', path: '/api/v1/admin/chats/messages/broadcast', code: 'admin.chats.broadcast', group: '/admin/chats', description: '发送系统广播' },
    { method: 'GET', path: '/api/v1/admin/chats/stats', code: 'admin.chats.stats', group: '/admin/chats', description: '获取聊天统计' },
];

export default { ADMIN_MENUS, ADMIN_PERMISSIONS };
