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
        name: '系统设置',
        path: '/admin/settings',
        component: 'Settings',
        icon: 'SettingOutlined',
        order: 99,
        permission: 'admin.settings.view',
        description: '系统参数设置',
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
];

export default { ADMIN_MENUS, ADMIN_PERMISSIONS };
