/**
 * 权限常量定义
 *
 * @description
 * 定义系统中所有权限码常量，与后端 Permission.Code 字段对应
 * 权限码格式：{module}.{resource}.{action}
 *
 * 操作类型说明：
 * - list: 列表查询
 * - read: 详情查询
 * - create: 创建
 * - update: 更新
 * - delete: 删除
 */

/**
 * 超级管理员权限（拥有所有权限）
 */
export const SUPER_ADMIN = '*';

/**
 * 用户管理权限
 */
export const USER_PERMISSIONS = {
    /** 用户列表 */
    LIST: 'admin.users.list',
    /** 用户详情 */
    READ: 'admin.users.read',
    /** 创建用户 */
    CREATE: 'admin.users.create',
    /** 更新用户 */
    UPDATE: 'admin.users.update',
    /** 删除用户 */
    DELETE: 'admin.users.delete',
    /** 用户状态管理 */
    STATUS: 'admin.users.status',
} as const;

/**
 * 游戏管理权限
 */
export const GAME_PERMISSIONS = {
    /** 游戏列表 */
    LIST: 'admin.games.list',
    /** 游戏详情 */
    READ: 'admin.games.read',
    /** 创建游戏 */
    CREATE: 'admin.games.create',
    /** 更新游戏 */
    UPDATE: 'admin.games.update',
    /** 删除游戏 */
    DELETE: 'admin.games.delete',
} as const;

/**
 * 订单管理权限
 */
export const ORDER_PERMISSIONS = {
    /** 订单列表 */
    LIST: 'admin.orders.list',
    /** 订单详情 */
    READ: 'admin.orders.read',
    /** 创建订单 */
    CREATE: 'admin.orders.create',
    /** 更新订单 */
    UPDATE: 'admin.orders.update',
    /** 删除订单 */
    DELETE: 'admin.orders.delete',
    /** 订单取消 */
    CANCEL: 'admin.orders.cancel',
    /** 订单退款 */
    REFUND: 'admin.orders.refund',
} as const;

/**
 * 陪玩师管理权限
 */
export const PLAYER_PERMISSIONS = {
    /** 陪玩师列表 */
    LIST: 'admin.players.list',
    /** 陪玩师详情 */
    READ: 'admin.players.read',
    /** 创建陪玩师 */
    CREATE: 'admin.players.create',
    /** 更新陪玩师 */
    UPDATE: 'admin.players.update',
    /** 删除陪玩师 */
    DELETE: 'admin.players.delete',
    /** 审核陪玩师 */
    AUDIT: 'admin.players.audit',
} as const;

/**
 * 角色管理权限
 */
export const ROLE_PERMISSIONS = {
    /** 角色列表 */
    LIST: 'admin.roles.list',
    /** 角色详情 */
    READ: 'admin.roles.read',
    /** 创建角色 */
    CREATE: 'admin.roles.create',
    /** 更新角色 */
    UPDATE: 'admin.roles.update',
    /** 删除角色 */
    DELETE: 'admin.roles.delete',
    /** 分配权限 */
    ASSIGN_PERMISSIONS: 'admin.roles.permissions',
    /** 分配用户 */
    ASSIGN_USER: 'admin.roles.assign-user',
} as const;

/**
 * 权限管理权限
 */
export const PERMISSION_PERMISSIONS = {
    /** 权限列表 */
    LIST: 'admin.permissions.list',
    /** 权限详情 */
    READ: 'admin.permissions.read',
    /** 创建权限 */
    CREATE: 'admin.permissions.create',
    /** 更新权限 */
    UPDATE: 'admin.permissions.update',
    /** 删除权限 */
    DELETE: 'admin.permissions.delete',
    /** 权限分组 */
    GROUPS: 'admin.permissions.groups',
} as const;

/**
 * 菜单管理权限
 */
export const MENU_PERMISSIONS = {
    /** 菜单列表 */
    LIST: 'admin.menus.list',
    /** 菜单详情 */
    READ: 'admin.menus.read',
    /** 创建菜单 */
    CREATE: 'admin.menus.create',
    /** 更新菜单 */
    UPDATE: 'admin.menus.update',
    /** 删除菜单 */
    DELETE: 'admin.menus.delete',
} as const;

/**
 * 服务项目管理权限
 */
export const SERVICE_ITEM_PERMISSIONS = {
    /** 服务项目列表 */
    LIST: 'admin.service-items.list',
    /** 服务项目详情 */
    READ: 'admin.service-items.read',
    /** 创建服务项目 */
    CREATE: 'admin.service-items.create',
    /** 更新服务项目 */
    UPDATE: 'admin.service-items.update',
    /** 删除服务项目 */
    DELETE: 'admin.service-items.delete',
    /** 批量更新状态 */
    BATCH_STATUS: 'admin.service-items.batch-status',
} as const;

/**
 * 佣金管理权限
 */
export const COMMISSION_PERMISSIONS = {
    /** 佣金列表 */
    LIST: 'admin.commissions.list',
    /** 佣金详情 */
    READ: 'admin.commissions.read',
    /** 创建佣金配置 */
    CREATE: 'admin.commissions.create',
    /** 更新佣金配置 */
    UPDATE: 'admin.commissions.update',
    /** 删除佣金配置 */
    DELETE: 'admin.commissions.delete',
    /** 触发结算 */
    SETTLE: 'admin.commissions.settle',
} as const;

/**
 * 提现管理权限
 */
export const WITHDRAW_PERMISSIONS = {
    /** 提现列表 */
    LIST: 'admin.withdraws.list',
    /** 提现详情 */
    READ: 'admin.withdraws.read',
    /** 审批提现 */
    APPROVE: 'admin.withdraws.approve',
    /** 拒绝提现 */
    REJECT: 'admin.withdraws.reject',
    /** 完成提现打款 */
    COMPLETE: 'admin.withdraws.complete',
} as const;

/**
 * 仪表板权限
 */
export const DASHBOARD_PERMISSIONS = {
    /** 查看仪表板 */
    VIEW: 'admin.dashboard.view',
    /** 统计数据 */
    STATS: 'admin.stats.read',
} as const;

/**
 * 系统管理权限
 */
export const SYSTEM_PERMISSIONS = {
    /** 系统信息 */
    INFO: 'admin.system.info',
    /** 系统配置 */
    CONFIG: 'admin.system.config',
} as const;

/**
 * 所有权限集合（便于权限组件使用）
 */
export const PERMISSIONS = {
    SUPER_ADMIN,
    USER: USER_PERMISSIONS,
    GAME: GAME_PERMISSIONS,
    ORDER: ORDER_PERMISSIONS,
    PLAYER: PLAYER_PERMISSIONS,
    ROLE: ROLE_PERMISSIONS,
    PERMISSION: PERMISSION_PERMISSIONS,
    MENU: MENU_PERMISSIONS,
    SERVICE_ITEM: SERVICE_ITEM_PERMISSIONS,
    COMMISSION: COMMISSION_PERMISSIONS,
    WITHDRAW: WITHDRAW_PERMISSIONS,
    DASHBOARD: DASHBOARD_PERMISSIONS,
    SYSTEM: SYSTEM_PERMISSIONS,
} as const;

/**
 * 权限类型定义
 */
export type PermissionCode =
    | typeof SUPER_ADMIN
    | typeof USER_PERMISSIONS[keyof typeof USER_PERMISSIONS]
    | typeof GAME_PERMISSIONS[keyof typeof GAME_PERMISSIONS]
    | typeof ORDER_PERMISSIONS[keyof typeof ORDER_PERMISSIONS]
    | typeof PLAYER_PERMISSIONS[keyof typeof PLAYER_PERMISSIONS]
    | typeof ROLE_PERMISSIONS[keyof typeof ROLE_PERMISSIONS]
    | typeof PERMISSION_PERMISSIONS[keyof typeof PERMISSION_PERMISSIONS]
    | typeof MENU_PERMISSIONS[keyof typeof MENU_PERMISSIONS]
    | typeof SERVICE_ITEM_PERMISSIONS[keyof typeof SERVICE_ITEM_PERMISSIONS]
    | typeof COMMISSION_PERMISSIONS[keyof typeof COMMISSION_PERMISSIONS]
    | typeof WITHDRAW_PERMISSIONS[keyof typeof WITHDRAW_PERMISSIONS]
    | typeof DASHBOARD_PERMISSIONS[keyof typeof DASHBOARD_PERMISSIONS]
    | typeof SYSTEM_PERMISSIONS[keyof typeof SYSTEM_PERMISSIONS];

export default PERMISSIONS;
