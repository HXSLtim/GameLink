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
 * VIP管理权限
 */
export const VIP_PERMISSIONS = {
    /** VIP等级列表 */
    LIST: 'admin.vip.list',
    /** VIP等级详情 */
    READ: 'admin.vip.read',
    /** 创建VIP等级 */
    CREATE: 'admin.vip.create',
    /** 更新VIP等级 */
    UPDATE: 'admin.vip.update',
    /** 删除VIP等级 */
    DELETE: 'admin.vip.delete',
    /** 设置默认等级 */
    SET_DEFAULT: 'admin.vip.set-default',
    /** 批量更新状态 */
    BATCH_STATUS: 'admin.vip.batch-status',
    /** 批量删除 */
    BATCH_DELETE: 'admin.vip.batch-delete',
    /** VIP配置管理 */
    CONFIG: 'admin.vip.config',
} as const;

/**
/**
 * 团队管理权限
 */
export const TEAM_PERMISSIONS = {
    /** 团队列表 */
    LIST: 'admin.teams.list',
    /** 团队详情 */
    READ: 'admin.teams.read',
    /** 创建团队 */
    CREATE: 'admin.teams.create',
    /** 更新团队 */
    UPDATE: 'admin.teams.update',
    /** 删除团队 */
    DELETE: 'admin.teams.delete',
    /** 团队状态管理 */
    STATUS: 'admin.teams.status',
    /** 批量更新状态 */
    BATCH_STATUS: 'admin.teams.batch-status',
    /** 批量删除 */
    BATCH_DELETE: 'admin.teams.batch-delete',
    /** 成员管理 */
    MEMBERS_MANAGE: 'admin.teams.members-manage',
    /** 转让队长 */
    TRANSFER_LEADER: 'admin.teams.transfer-leader',
} as const;
/**
 * 充值管理权限
 */
export const RECHARGE_PERMISSIONS = {
    /** 充值档位列表 */
    LIST_OPTIONS: 'admin.recharge.options.list',
    /** 充值档位详情 */
    READ_OPTION: 'admin.recharge.options.read',
    /** 创建充值档位 */
    CREATE_OPTION: 'admin.recharge.options.create',
    /** 更新充值档位 */
    UPDATE_OPTION: 'admin.recharge.options.update',
    /** 删除充值档位 */
    DELETE_OPTION: 'admin.recharge.options.delete',
    /** 批量更新状态 */
    BATCH_UPDATE_STATUS: 'admin.recharge.options.batch-status',
    /** 批量删除 */
    BATCH_DELETE: 'admin.recharge.options.batch-delete',
    /** 充值记录列表 */
    LIST_RECORDS: 'admin.recharge.records.list',
    /** 充值记录详情 */
    READ_RECORD: 'admin.recharge.records.read',
    /** 退款 */
    REFUND: 'admin.recharge.records.refund',
    /** 充值统计 */
    STATS: 'admin.recharge.stats',
} as const;

/**
 * 纠纷管理权限
 */
export const DISPUTE_PERMISSIONS = {
    /** 纠纷列表 */
    LIST: 'admin.disputes.list',
    /** 纠纷详情 */
    READ: 'admin.disputes.read',
    /** 分配客服 */
    ASSIGN: 'admin.disputes.assign',
    /** 解决纠纷 */
    RESOLVE: 'admin.disputes.resolve',
    /** 回滚分配 */
    ROLLBACK: 'admin.disputes.rollback',
    /** 查看统计 */
    STATS: 'admin.disputes.stats',
    /** 导出记录 */
    EXPORT: 'admin.disputes.export',
    /** 批量操作 */
    BATCH: 'admin.disputes.batch',
} as const;

/**
 * 优惠券管理权限
 */
export const COUPON_PERMISSIONS = {
    /** 优惠券模板列表 */
    LIST_TEMPLATES: 'admin.coupons.templates.list',
    /** 优惠券模板详情 */
    READ_TEMPLATE: 'admin.coupons.templates.read',
    /** 创建优惠券模板 */
    CREATE_TEMPLATE: 'admin.coupons.templates.create',
    /** 更新优惠券模板 */
    UPDATE_TEMPLATE: 'admin.coupons.templates.update',
    /** 删除优惠券模板 */
    DELETE_TEMPLATE: 'admin.coupons.templates.delete',
    /** 发放优惠券 */
    ISSUE: 'admin.coupons.issue',
    /** 用户优惠券列表 */
    LIST_USER_COUPONS: 'admin.coupons.user.list',
    /** 用户优惠券详情 */
    READ_USER_COUPON: 'admin.coupons.user.read',
    /** 作废优惠券 */
    VOID: 'admin.coupons.void',
    /** 批量删除模板 */
    BATCH_DELETE_TEMPLATES: 'admin.coupons.templates.batch-delete',
    /** 批量更新状态 */
    BATCH_UPDATE_STATUS: 'admin.coupons.batch-status',
    /** 优惠券统计 */
    STATS: 'admin.coupons.stats',
} as const;

/**
 * 推荐系统权限
 */
export const REFERRAL_PERMISSIONS = {
    /** 推荐关系列表 */
    LIST: 'admin.referrals.list',
    /** 推荐详情 */
    READ: 'admin.referrals.read',
    /** 推荐统计 */
    STATS: 'admin.referrals.stats',
    /** 邀请码列表 */
    LIST_CODES: 'admin.referrals.codes.list',
    /** 创建邀请码 */
    CREATE_CODE: 'admin.referrals.codes.create',
    /** 删除邀请码 */
    DELETE_CODE: 'admin.referrals.codes.delete',
    /** 批量删除邀请码 */
    BATCH_DELETE_CODES: 'admin.referrals.codes.batch-delete',
    /** 奖励配置 */
    CONFIG: 'admin.referrals.config',
    /** 奖励发放记录 */
    LIST_REWARDS: 'admin.referrals.rewards.list',
    /** 手动发放奖励 */
    ISSUE_REWARD: 'admin.referrals.rewards.issue',
} as const;

/**
 * 路由规则权限
 */
export const ROUTING_PERMISSIONS = {
    /** 路由规则列表 */
    LIST: 'admin.routing.list',
    /** 路由规则详情 */
    READ: 'admin.routing.read',
    /** 创建路由规则 */
    CREATE: 'admin.routing.create',
    /** 更新路由规则 */
    UPDATE: 'admin.routing.update',
    /** 删除路由规则 */
    DELETE: 'admin.routing.delete',
    /** 复制规则 */
    COPY: 'admin.routing.copy',
    /** 测试规则 */
    TEST: 'admin.routing.test',
    /** 批量更新状态 */
    BATCH_UPDATE_STATUS: 'admin.routing.batch-status',
} as const;

/**
 * 活动管理权限
 */
export const ACTIVITY_PERMISSIONS = {
    /** 活动列表 */
    LIST: 'admin.activities.list',
    /** 活动详情 */
    READ: 'admin.activities.read',
    /** 创建活动 */
    CREATE: 'admin.activities.create',
    /** 更新活动 */
    UPDATE: 'admin.activities.update',
    /** 删除活动 */
    DELETE: 'admin.activities.delete',
    /** 活动状态管理 */
    STATUS: 'admin.activities.status',
    /** 奖励管理 */
    MANAGE_REWARDS: 'admin.activities.rewards.manage',
    /** 活动统计 */
    STATS: 'admin.activities.stats',
    /** 批量操作 */
    BATCH: 'admin.activities.batch',
} as const;

/**
 * 结算公司权限
 */
export const SETTLEMENT_PERMISSIONS = {
    /** 结算公司列表 */
    LIST: 'admin.settlement.list',
    /** 结算公司详情 */
    READ: 'admin.settlement.read',
    /** 创建结算公司 */
    CREATE: 'admin.settlement.create',
    /** 更新结算公司 */
    UPDATE: 'admin.settlement.update',
    /** 删除结算公司 */
    DELETE: 'admin.settlement.delete',
    /** 分配陪玩师 */
    ASSIGN_PLAYER: 'admin.settlement.assign-player',
    /** 移除陪玩师 */
    REMOVE_PLAYER: 'admin.settlement.remove-player',
    /** 结算统计 */
    STATS: 'admin.settlement.stats',
} as const;

/**
 * 聊天管理权限
 */
export const CHAT_PERMISSIONS = {
    /** 聊天记录列表 */
    LIST: 'admin.chat.list',
    /** 聊天记录详情 */
    READ: 'admin.chat.read',
    /** 聊天房间列表 */
    LIST_ROOMS: 'admin.chat.rooms.list',
    /** 聊天房间详情 */
    READ_ROOM: 'admin.chat.rooms.read',
    /** 禁用房间 */
    DISABLE_ROOM: 'admin.chat.rooms.disable',
    /** 聊天统计 */
    STATS: 'admin.chat.stats',
} as const;

/**
 * 排行榜管理权限
 */
export const GAME_RANK_PERMISSIONS = {
    /** 排行榜列表 */
    LIST: 'admin.ranks.list',
    /** 排行榜详情 */
    READ: 'admin.ranks.read',
    /** 创建排行榜 */
    CREATE: 'admin.ranks.create',
    /** 更新排行榜 */
    UPDATE: 'admin.ranks.update',
    /** 删除排行榜 */
    DELETE: 'admin.ranks.delete',
    /** 手动刷新排行 */
    REFRESH: 'admin.ranks.refresh',
} as const;

/** 系统管理权限
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
    VIP: VIP_PERMISSIONS,
    TEAM: TEAM_PERMISSIONS,
    RECHARGE: RECHARGE_PERMISSIONS,
    DISPUTE: DISPUTE_PERMISSIONS,
    COUPON: COUPON_PERMISSIONS,
    REFERRAL: REFERRAL_PERMISSIONS,
    ROUTING: ROUTING_PERMISSIONS,
    ACTIVITY: ACTIVITY_PERMISSIONS,
    SETTLEMENT: SETTLEMENT_PERMISSIONS,
    CHAT: CHAT_PERMISSIONS,
    GAME_RANK: GAME_RANK_PERMISSIONS,
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
    | typeof VIP_PERMISSIONS[keyof typeof VIP_PERMISSIONS]
    | typeof RECHARGE_PERMISSIONS[keyof typeof RECHARGE_PERMISSIONS]
    | typeof TEAM_PERMISSIONS[keyof typeof TEAM_PERMISSIONS]
    | typeof DISPUTE_PERMISSIONS[keyof typeof DISPUTE_PERMISSIONS]
    | typeof COUPON_PERMISSIONS[keyof typeof COUPON_PERMISSIONS]
    | typeof REFERRAL_PERMISSIONS[keyof typeof REFERRAL_PERMISSIONS]
    | typeof ROUTING_PERMISSIONS[keyof typeof ROUTING_PERMISSIONS]
    | typeof ACTIVITY_PERMISSIONS[keyof typeof ACTIVITY_PERMISSIONS]
    | typeof SETTLEMENT_PERMISSIONS[keyof typeof SETTLEMENT_PERMISSIONS]
    | typeof CHAT_PERMISSIONS[keyof typeof CHAT_PERMISSIONS]
    | typeof GAME_RANK_PERMISSIONS[keyof typeof GAME_RANK_PERMISSIONS]
    | typeof DASHBOARD_PERMISSIONS[keyof typeof DASHBOARD_PERMISSIONS]
    | typeof SYSTEM_PERMISSIONS[keyof typeof SYSTEM_PERMISSIONS];

export default PERMISSIONS;
