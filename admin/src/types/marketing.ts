/**
 * Marketing Module Types
 * 营销模块类型定义
 *
 * Contains types for:
 * - VIP (会员等级系统)
 * - Coupon (优惠券系统)
 * - Referral (推荐系统)
 * - Team (团队系统)
 * - Activity (活动系统)
 */

// ============================================================================
// VIP Types (会员等级)
// ============================================================================

/**
 * VIP Level Configuration
 * VIP等级配置
 */
export interface VIPLevel {
    id: number;
    slug: string;                    // 唯一标识符
    title: string;                   // 等级名称
    expRequired: number;             // 升级所需经验值
    orderDiscount: number;           // 订单折扣比例 (0-1)
    monthlyCouponTemplateId?: number; // 每月赠送优惠券模板ID
    monthlyCouponCount: number;       // 每月优惠券数量
    iconUrl: string;                 // 等级图标
    color: string;                   // 等级颜色
    benefits: string;                // 权益描述 (JSON或文本)
    sortOrder: number;               // 排序
    isDefault: boolean;              // 是否默认等级
    isActive: boolean;               // 是否启用
    createdAt: string;
    updatedAt: string;
}

/**
 * VIP System Configuration
 * VIP系统配置
 */
export interface VIPConfig {
    id: number;
    configKey: string;               // 配置键
    configValue: string;             // 配置值
    description: string;             // 描述
    createdAt: string;
    updatedAt: string;
}

/**
 * VIP Config Key Constants
 * VIP配置键常量
 */
export const VIP_CONFIG_KEYS = {
    UNLOCK_BY_CONSUME: 'unlock_by_consume',    // 累计消费解锁门槛（分）
    UNLOCK_BY_RECHARGE: 'unlock_by_recharge',  // 累计充值解锁门槛（分）
    EXPIRE_DAYS: 'expire_days',                // VIP过期天数（0=永久）
} as const;

export type VIPConfigKey = typeof VIP_CONFIG_KEYS[keyof typeof VIP_CONFIG_KEYS];

// ============================================================================
// Coupon Types (优惠券)
// ============================================================================

/**
 * Coupon Type
 * 优惠券类型
 */
export type CouponType = 'deduct' | 'discount';

/**
 * Coupon Scope
 * 优惠券适用范围
 */
export type CouponScope = 'all' | 'game' | 'item';

/**
 * Coupon Source
 * 优惠券来源
 */
export type CouponSource = 'new_user' | 'link' | 'vip' | 'activity' | 'manual' | 'referral' | 'team';

/**
 * Coupon State
 * 优惠券状态
 */
export type CouponState = 'available' | 'locked' | 'used' | 'expired';

/**
 * Validity Type
 * 有效期类型
 */
export type ValidityType = 'days' | 'fixed';

/**
 * Coupon Template
 * 优惠券模板
 */
export interface CouponTemplate {
    id: number;
    name: string;
    type: CouponType;
    source: CouponSource;
    description?: string;

    // 折扣配置 (Discount Configuration)
    minAmountCents: number;          // 最低使用金额（分）
    deductAmountCents: number;       // 减免金额（分）- deduct类型
    discountRate: number;            // 折扣率 (0-1) - discount类型
    maxDiscountCents: number;        // 最大折扣金额（分）

    // 适用范围 (Applicability)
    scope: CouponScope;
    gameIds: string;                 // 适用游戏ID列表 (JSON数组字符串)
    itemIds: string;                 // 适用服务项ID列表 (JSON数组字符串)

    // 有效期 (Validity)
    validityType: ValidityType;
    validityDays: number;            // 有效天数
    fixedExpireAt?: string;          // 固定过期时间

    // 领取配置 (Claim Configuration)
    totalCount: number;              // 总发行量
    claimedCount: number;            // 已领取数量
    perUserLimit: number;            // 每用户限领
    claimLink?: string;              // 领取链接

    isActive: boolean;
    createdAt: string;
    updatedAt: string;
}

/**
 * User Coupon
 * 用户优惠券
 */
export interface Coupon {
    id: number;
    templateId: number;
    userId: number;
    state: CouponState;

    // 模板字段 (Denormalized from template)
    name: string;
    type: CouponType;
    source: CouponSource;
    minAmountCents: number;
    deductAmountCents: number;
    discountRate: number;
    maxDiscountCents: number;
    scope: CouponScope;
    gameIds: string;
    itemIds: string;

    // 时间字段 (Time Fields)
    claimedAt?: string;
    expireAt: string;
    usedAt?: string;

    // 锁定信息 (Lock Information)
    lockedByOrderId?: number;
    lockedAt?: string;

    // 使用信息 (Usage Information)
    usedOrderId?: number;
    discountCents: number;           // 实际优惠金额

    createdAt: string;
    updatedAt: string;
}

/**
 * Coupon with Relations
 * 带关联信息的优惠券
 */
export interface CouponWithTemplate extends Coupon {
    template?: CouponTemplate;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
}

/**
 * Coupon Statistics
 * 优惠券统计
 */
export interface CouponStats {
    totalTemplates: number;
    activeTemplates: number;
    totalCoupons: number;
    availableCoupons: number;
    usedCoupons: number;
    expiredCoupons: number;
    totalDiscountCents: number;
}

// ============================================================================
// Referral Types (推荐系统)
// ============================================================================

/**
 * Referral Type
 * 推荐类型
 */
export type ReferralType = 'user' | 'player';

/**
 * Referral Status
 * 推荐状态
 */
export type ReferralStatus = 'pending' | 'completed' | 'canceled';

/**
 * Reward Type
 * 奖励类型
 */
export type RewardType = 'referrer' | 'referee';

/**
 * Referral Reward Status
 * 推荐奖励状态
 */
export type ReferralRewardStatus = 'pending' | 'issued' | 'failed';

/**
 * Referral Config
 * 推荐系统配置
 */
export interface ReferralConfig {
    id: number;
    configKey: string;
    configValue: string;
    description: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * Referral Config Key Constants
 * 推荐配置键常量
 */
export const REFERRAL_CONFIG_KEYS = {
    ENABLED: 'referral_enabled',              // 推荐功能开关
    REFERRER_REWARD_CENTS: 'referrer_reward_cents',  // 推荐人奖励（分）
    REFEREE_REWARD_CENTS: 'referee_reward_cents',    // 被推荐人奖励（分）
    MAX_REWARD_USES: 'max_reward_uses',       // 每个邀请码最大使用次数
    EXPIRE_DAYS: 'expire_days',               // 邀请码有效期（天）
} as const;

export type ReferralConfigKey = typeof REFERRAL_CONFIG_KEYS[keyof typeof REFERRAL_CONFIG_KEYS];

/**
 * Referral Code
 * 推荐码
 */
export interface ReferralCode {
    id: number;
    code: string;
    ownerId: number;
    type: ReferralType;
    maxUses: number;                  // 最大使用次数
    usedCount: number;                // 已使用次数
    expiresAt: string;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
    owner?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
}

/**
 * Referral Record
 * 推荐记录
 */
export interface Referral {
    id: number;
    referrerId: number;               // 推荐人ID
    refereeId: number;                // 被推荐人ID
    codeId: number;                   // 使用的推荐码ID
    type: ReferralType;
    status: ReferralStatus;
    completedAt?: string;
    cancelReason?: string;
    createdAt: string;
    updatedAt: string;
    referrer?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    referee?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    code?: ReferralCode;
}

/**
 * Referral Reward
 * 推荐奖励
 */
export interface ReferralReward {
    id: number;
    referralId: number;
    userId: number;                   // 奖励接收用户ID
    type: RewardType;                // 推荐人/被推荐人
    amountCents: number;             // 奖励金额（分）
    status: ReferralRewardStatus;
    issuedAt?: string;
    failedAt?: string;
    failureReason?: string;
    createdAt: string;
    updatedAt: string;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    referral?: Referral;
}

/**
 * Referral Statistics
 * 推荐系统统计
 */
export interface ReferralStats {
    totalReferrals: number;
    completedReferrals: number;
    pendingReferrals: number;
    canceledReferrals: number;
    totalRewardsCents: number;
    issuedRewardsCents: number;
    pendingRewardsCents: number;
    failedRewardsCents: number;
    activeCodes: number;
    totalCodes: number;
}

// ============================================================================
// Team Types (团队系统)
// ============================================================================

/**
 * Team Status
 * 团队状态
 */
export type TeamStatus = 'active' | 'busy' | 'inactive';

/**
 * Team Member Role
 * 团队成员角色
 */
export type TeamMemberRole = 'leader' | 'member';

/**
 * Team Member Status
 * 团队成员状态
 */
export type TeamMemberStatus = 'active' | 'left' | 'kicked';

/**
 * Team Invite Status
 * 团队邀请状态
 */
export type TeamInviteStatus = 'pending' | 'accepted' | 'rejected' | 'expired';

/**
 * Team
 * 团队
 */
export interface Team {
    id: number;
    name: string;
    description?: string;
    avatarUrl?: string;
    leaderId: number;                 // 队长ID
    leader?: {
        id: number;
        nickname: string;
        avatar?: string;
        rank?: string;
    };
    status: TeamStatus;
    maxMembers: number;               // 最大成员数
    memberCount: number;              // 当前成员数
    incomeShareType: 'equal' | 'custom'; // 收益分配方式
    leaderBonusRate: number;          // 队长额外加成比例
    totalOrderCount: number;          // 总订单数
    totalIncomeCents: number;         // 总收益（分）
    currentOrderId?: number;          // 当前进行中的订单ID
    createdAt: string;
    updatedAt: string;
}

/**
 * Team Member
 * 团队成员
 */
export interface TeamMember {
    id: number;
    teamId: number;
    playerId: number;
    player?: {
        id: number;
        nickname: string;
        avatar?: string;
        rank?: string;
    };
    role: TeamMemberRole;
    status: TeamMemberStatus;
    sortOrder: number;
    orderCount: number;               // 订单数
    incomeCents: number;              // 收益（分）
    joinedAt: string;
    leftAt?: string;
}

/**
 * Team Invite
 * 团队邀请
 */
export interface TeamInvite {
    id: number;
    teamId: number;
    playerId: number;
    player?: {
        id: number;
        nickname: string;
        avatar?: string;
    };
    inviterId: number;
    inviter?: {
        id: number;
        nickname: string;
        avatar?: string;
    };
    status: TeamInviteStatus;
    expireAt: string;
    message?: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * Team Statistics
 * 团队统计
 */
export interface TeamStats {
    totalTeams: number;
    activeTeams: number;
    busyTeams: number;
    inactiveTeams: number;
    totalMembers: number;
    totalIncomeCents: number;
}

// ============================================================================
// Activity Types (活动系统)
// ============================================================================

/**
 * Activity Status
 * 活动状态
 */
export type ActivityStatus = 'draft' | 'preheat' | 'active' | 'paused' | 'ended' | 'canceled';

/**
 * Activity Type
 * 活动类型
 */
export type ActivityType = 'coupon' | 'discount' | 'gift';

/**
 * Activity
 * 活动
 */
export interface Activity {
    id: number;
    name: string;
    description?: string;
    type: ActivityType;
    status: ActivityStatus;
    coverUrl?: string;
    bannerUrl?: string;

    // 时间控制 (Time Control)
    preheatAt?: string;               // 预热时间
    startAt: string;                  // 开始时间
    endAt: string;                    // 结束时间

    // 参与限制 (Participation Limits)
    totalLimit: number;               // 总参与人数限制
    dailyLimit: number;               // 每日参与人数限制
    perUserLimit: number;             // 每用户参与次数限制

    // 统计数据 (Statistics)
    totalParticipants: number;        // 总参与人数
    todayParticipants: number;        // 今日参与人数
    totalClaimed: number;             // 总领取数

    // 配置 (Configuration)
    allowVipStack: boolean;           // 允许VIP叠加
    rules?: string;                   // 活动规则
    sortOrder: number;
    isVisible: boolean;               // 是否可见

    // 关联 (Relations)
    rewards?: ActivityReward[];

    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Reward
 * 活动奖励
 */
export interface ActivityReward {
    id: number;
    activityId: number;
    couponTemplateId: number;         // 优惠券模板ID
    couponCount: number;              // 优惠券数量
    probability: number;              // 中奖概率 (0-1)
    totalStock: number;               // 总库存
    remainingStock: number;           // 剩余库存
    sortOrder: number;

    // 关联 (Relations)
    activity?: Activity;
    couponTemplate?: {
        id: number;
        name: string;
        type: string;
    };

    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Participation
 * 活动参与记录
 */
export interface ActivityParticipation {
    id: number;
    activityId: number;
    userId: number;
    rewardId: number;
    couponIds?: string;               // 获得的优惠券ID列表 (JSON数组)
    claimedAt: string;
    clientIp?: string;

    // 关联 (Relations)
    activity?: Activity;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    reward?: ActivityReward;

    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Statistics
 * 活动统计
 */
export interface ActivityStats {
    activityId: number;
    totalParticipants: number;
    todayParticipants: number;
    totalClaimed: number;
    remainingStock?: number;
}

/**
 * All Activities Statistics Overview
 * 所有活动统计概览
 */
export interface AllActivityStats {
    totalActivities: number;
    activeActivities: number;
    draftActivities: number;
    endedActivities: number;
    totalParticipants: number;
    totalClaimed: number;
}
