import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const VipLevel = {
    NONE: 0,           // 非 VIP
    BRONZE: 1,         // 青铜 VIP
    SILVER: 2,         // 白银 VIP
    GOLD: 3,           // 黄金 VIP
    PLATINUM: 4,       // 铂金 VIP
    DIAMOND: 5         // 钻石 VIP
} as const;

export type VipLevel = typeof VipLevel[keyof typeof VipLevel];

export const VipStatus = {
    ACTIVE: 'active',       // 有效
    EXPIRED: 'expired',     // 已过期
    SUSPENDED: 'suspended'  // 已暂停
} as const;

export type VipStatus = typeof VipStatus[keyof typeof VipStatus];

// ============ Interfaces ============

export interface VipLevelConfig {
    level: VipLevel;
    name: string;
    icon: string;
    color: string;

    // 解锁条件
    thresholdCents: number;          // 累计消费门槛

    // 权益
    discountRate: number;            // 订单折扣 (如 0.95 = 95折)
    monthlyFreeCoupons: number;      // 每月免费优惠券
    prioritySupport: boolean;        // 优先客服
    exclusiveActivities: boolean;    // 专属活动
}

export interface UserVip {
    userId: number;
    level: VipLevel;
    status: VipStatus;

    // 累计消费
    totalSpentCents: number;
    currentYearSpentCents: number;

    // 有效期
    activatedAt?: string;
    expiresAt?: string;

    // 权益使用
    monthlyFreeCoupons: number;      // 本月免费优惠券数
    usedFreeCoupons: number;         // 已使用数
    discountRate: number;            // 折扣率

    // 升级进度
    nextLevel?: VipLevel;
    nextLevelThreshold?: number;     // 下一等级门槛
    progressPercent: number;         // 升级进度百分比
}

export interface VipBenefit {
    id: string;
    name: string;
    description: string;
    icon: string;
    unlockLevel: VipLevel;
    isUnlocked: boolean;
}

export interface DiscountResult {
    originalPrice: number;
    discountRate: number;
    discountAmount: number;
    finalPrice: number;
}

// ============ State & Actions ============

export interface VipState {
    userVip: UserVip | null;
    levelConfigs: VipLevelConfig[];
    benefits: VipBenefit[];
    loading: boolean;
    error: string | null;
}

export interface VipActions {
    fetchVipStatus: () => Promise<void>;
    fetchLevelConfigs: () => Promise<void>;
    claimMonthlyCoupon: () => Promise<void>;
    calculateVipDiscount: (originalPriceCents: number) => DiscountResult;
    getLevelName: (level: VipLevel) => string;
    getLevelColor: (level: VipLevel) => string;
}

// ============ Default Level Configs ============

const DEFAULT_LEVEL_CONFIGS: VipLevelConfig[] = [
    {
        level: VipLevel.NONE,
        name: '普通用户',
        icon: '👤',
        color: '#9CA3AF',
        thresholdCents: 0,
        discountRate: 1,
        monthlyFreeCoupons: 0,
        prioritySupport: false,
        exclusiveActivities: false
    },
    {
        level: VipLevel.BRONZE,
        name: '青铜 VIP',
        icon: '🥉',
        color: '#CD7F32',
        thresholdCents: 50000,      // ¥500
        discountRate: 0.98,
        monthlyFreeCoupons: 1,
        prioritySupport: false,
        exclusiveActivities: false
    },
    {
        level: VipLevel.SILVER,
        name: '白银 VIP',
        icon: '🥈',
        color: '#C0C0C0',
        thresholdCents: 200000,     // ¥2,000
        discountRate: 0.95,
        monthlyFreeCoupons: 2,
        prioritySupport: false,
        exclusiveActivities: false
    },
    {
        level: VipLevel.GOLD,
        name: '黄金 VIP',
        icon: '🥇',
        color: '#FFD700',
        thresholdCents: 500000,     // ¥5,000
        discountRate: 0.92,
        monthlyFreeCoupons: 3,
        prioritySupport: true,
        exclusiveActivities: false
    },
    {
        level: VipLevel.PLATINUM,
        name: '铂金 VIP',
        icon: '💎',
        color: '#E5E4E2',
        thresholdCents: 1000000,    // ¥10,000
        discountRate: 0.90,
        monthlyFreeCoupons: 5,
        prioritySupport: true,
        exclusiveActivities: true
    },
    {
        level: VipLevel.DIAMOND,
        name: '钻石 VIP',
        icon: '👑',
        color: '#B9F2FF',
        thresholdCents: 3000000,    // ¥30,000
        discountRate: 0.88,
        monthlyFreeCoupons: 8,
        prioritySupport: true,
        exclusiveActivities: true
    }
];

// ============ Store ============

export const useVipStore = create<VipState & VipActions>()(
    persist(
        (set, get) => ({
            userVip: null,
            levelConfigs: DEFAULT_LEVEL_CONFIGS,
            benefits: [],
            loading: false,
            error: null,

            fetchVipStatus: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<any>('/user/vip/status');

                    const userVip: UserVip = {
                        userId: data.userId,
                        level: data.level || VipLevel.NONE,
                        status: data.status || VipStatus.ACTIVE,
                        totalSpentCents: data.totalSpentCents || 0,
                        currentYearSpentCents: data.currentYearSpentCents || 0,
                        activatedAt: data.activatedAt,
                        expiresAt: data.expiresAt,
                        monthlyFreeCoupons: data.monthlyFreeCoupons || 0,
                        usedFreeCoupons: data.usedFreeCoupons || 0,
                        discountRate: data.discountRate || 1,
                        nextLevel: data.nextLevel,
                        nextLevelThreshold: data.nextLevelThreshold,
                        progressPercent: data.progressPercent || 0
                    };

                    set({ userVip, loading: false });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch VIP status' });
                }
            },

            fetchLevelConfigs: async () => {
                try {
                    const data = await http.get<VipLevelConfig[]>('/vip/levels');
                    if (data && data.length > 0) {
                        set({ levelConfigs: data });
                    }
                } catch (err: any) {
                    // 使用默认配置
                    console.warn('Using default VIP level configs');
                }
            },

            claimMonthlyCoupon: async () => {
                set({ loading: true, error: null });
                try {
                    await http.post('/user/vip/monthly-coupon');
                    // 刷新 VIP 状态
                    await get().fetchVipStatus();
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to claim monthly coupon' });
                    throw err;
                }
            },

            calculateVipDiscount: (originalPriceCents) => {
                const { userVip, levelConfigs } = get();

                if (!userVip || userVip.level === VipLevel.NONE) {
                    return {
                        originalPrice: originalPriceCents,
                        discountRate: 1,
                        discountAmount: 0,
                        finalPrice: originalPriceCents
                    };
                }

                const config = levelConfigs.find(c => c.level === userVip.level);
                const discountRate = config?.discountRate || 1;

                if (discountRate >= 1) {
                    return {
                        originalPrice: originalPriceCents,
                        discountRate: 1,
                        discountAmount: 0,
                        finalPrice: originalPriceCents
                    };
                }

                const discountAmount = Math.floor(originalPriceCents * (1 - discountRate));
                const finalPrice = originalPriceCents - discountAmount;

                return {
                    originalPrice: originalPriceCents,
                    discountRate,
                    discountAmount,
                    finalPrice
                };
            },

            getLevelName: (level) => {
                const { levelConfigs } = get();
                const config = levelConfigs.find(c => c.level === level);
                return config?.name || '普通用户';
            },

            getLevelColor: (level) => {
                const { levelConfigs } = get();
                const config = levelConfigs.find(c => c.level === level);
                return config?.color || '#9CA3AF';
            }
        }),
        {
            name: 'vip-storage',
            partialize: (state) => ({
                userVip: state.userVip,
                levelConfigs: state.levelConfigs
            })
        }
    )
);
