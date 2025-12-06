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
                name: '订单管理',
                path: '/admin/biz/order',
                component: 'Order',
                icon: 'ShoppingCartOutlined',
                order: 2,
                permission: 'admin.orders.list',
                description: '管理平台订单',
            },
            {
                name: '服务项目',
                path: '/admin/biz/service',
                component: 'ServiceItem',
                icon: 'GiftOutlined',
                order: 3,
                permission: 'admin.service-items.list',
                description: '管理服务项目',
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
        name: '评价管理',
        path: '/admin/reviews',
        component: 'Layout',
        icon: 'StarOutlined',
        order: 6,
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
                path: '/admin/review-reports',
                component: 'ReviewReports',
                icon: 'WarningOutlined',
                order: 3,
                permission: 'admin.review-reports.list',
                description: '处理评价举报',
            },
            {
                name: '敏感词管理',
                path: '/admin/sensitive-words',
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
];

export default { ADMIN_MENUS, ADMIN_PERMISSIONS };
