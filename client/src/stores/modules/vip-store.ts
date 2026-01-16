import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const VipStatus = {
    ACTIVE: 'active',       // 有效
    EXPIRED: 'expired',     // 已过期
    SUSPENDED: 'suspended'  // 已暂停
} as const;

export type VipStatus = typeof VipStatus[keyof typeof VipStatus];

// ============ Interfaces ============

// Matches backend VipLevel model from api/internal/model/vip.go
export interface VipLevel {
    id: number;
    slug: string;                // 等级标识 (vip1, vip2, svip1, etc.)
    title: string;               // 等级名称
    expRequired: number;         // 升级所需累计消费/经验（分）
    orderDiscount: number;       // 下单永久折扣 (0.98 = 98折, 1.0 = 无折扣)
    monthlyCouponTemplateId?: number;
    monthlyCouponCount: number;   // 每月发放数量
    iconUrl: string;             // 等级图标URL
    color: string;               // 等级颜色
    benefits: Record<string, unknown>; // 其他权益描述 (JSON)
    sortOrder: number;           // 排序（越小越靠前）
    isDefault: boolean;          // 是否默认等级
    isActive: boolean;           // 是否启用
    createdAt?: string;
    updatedAt?: string;
}

export interface UserVip {
    userId: number;
    vipLevelId?: number;
    vipLevel?: VipLevel;         // 关联的VIP等级详情
    status: VipStatus;

    // 累计消费/经验
    totalRechargeCents: number;   // 累计充值（分）
    totalSpentCents: number;      // 累计消费（分）
    vipExp: number;              // VIP经验（分）

    // 有效期
    vipUnlocked: boolean;
    vipUnlockedAt?: string;
    vipExpireAt?: string;

    // 权益
    currentMonthCouponClaimedAt?: string; // 本月月度券已领取时间

    // 升级进度
    nextLevel?: VipLevel;
    nextLevelThreshold?: number;     // 下一等级门槛
    progressPercent: number;         // 当前等级进度百分比
}

export interface VipBenefit {
    id: string;
    name: string;
    description: string;
    icon: string;
    requiredLevel: string;  // VIP level slug required
    isUnlocked: boolean;
}

export interface DiscountResult {
    originalPrice: number;
    originalPriceCents: number;
    discountRate: number;
    discountAmount: number;
    discountAmountCents: number;
    finalPrice: number;
    finalPriceCents: number;
}

// ============ State & Actions ============

export interface VipState {
    userVip: UserVip | null;
    levels: VipLevel[];
    benefits: VipBenefit[];
    loading: boolean;
    error: string | null;
}

export interface VipActions {
    fetchVipStatus: () => Promise<void>;
    fetchLevels: () => Promise<void>;
    claimMonthlyCoupon: () => Promise<void>;
    calculateVipDiscount: (priceCents: number) => DiscountResult;
    getLevelBySlug: (slug: string) => VipLevel | undefined;
    getLevelByTitle: (title: string) => VipLevel | undefined;
}

// ============ Default Level Configs (Fallback) ============

const DEFAULT_LEVELS: VipLevel[] = [
    {
        id: 0,
        slug: 'none',
        title: '普通用户',
        expRequired: 0,
        orderDiscount: 1.0,
        monthlyCouponCount: 0,
        iconUrl: '',
        color: '#9CA3AF',
        benefits: {},
        sortOrder: 0,
        isDefault: true,
        isActive: true
    },
    {
        id: 1,
        slug: 'bronze',
        title: '青铜 VIP',
        expRequired: 50000,      // ¥500
        orderDiscount: 0.98,
        monthlyCouponCount: 1,
        iconUrl: '',
        color: '#CD7F32',
        benefits: {},
        sortOrder: 1,
        isDefault: false,
        isActive: true
    },
    {
        id: 2,
        slug: 'silver',
        title: '白银 VIP',
        expRequired: 200000,     // ¥2,000
        orderDiscount: 0.95,
        monthlyCouponCount: 2,
        iconUrl: '',
        color: '#C0C0C0',
        benefits: {},
        sortOrder: 2,
        isDefault: false,
        isActive: true
    },
    {
        id: 3,
        slug: 'gold',
        title: '黄金 VIP',
        expRequired: 500000,     // ¥5,000
        orderDiscount: 0.92,
        monthlyCouponCount: 3,
        iconUrl: '',
        color: '#FFD700',
        benefits: {},
        sortOrder: 3,
        isDefault: false,
        isActive: true
    },
    {
        id: 4,
        slug: 'platinum',
        title: '铂金 VIP',
        expRequired: 1000000,    // ¥10,000
        orderDiscount: 0.90,
        monthlyCouponCount: 5,
        iconUrl: '',
        color: '#E5E4E2',
        benefits: {},
        sortOrder: 4,
        isDefault: false,
        isActive: true
    },
    {
        id: 5,
        slug: 'diamond',
        title: '钻石 VIP',
        expRequired: 3000000,    // ¥30,000
        orderDiscount: 0.88,
        monthlyCouponCount: 8,
        iconUrl: '',
        color: '#B9F2FF',
        benefits: {},
        sortOrder: 5,
        isDefault: false,
        isActive: true
    }
];

// ============ Store ============

export const useVipStore = create<VipState & VipActions>()(
    persist(
        (set, get) => ({
            userVip: null,
            levels: DEFAULT_LEVELS,
            benefits: [],
            loading: false,
            error: null,

            fetchVipStatus: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<{
                        userId: number;
                        vipLevelId?: number;
                        vipLevel?: VipLevel;
                        status?: VipStatus;
                        totalRechargeCents?: number;
                        totalSpentCents?: number;
                        vipExp?: number;
                        vipUnlocked?: boolean;
                        vipUnlockedAt?: string;
                        vipExpireAt?: string;
                        currentMonthCouponClaimedAt?: string;
                    }>('/user/vip/status');

                    set({
                        userVip: {
                            userId: data.userId,
                            vipLevelId: data.vipLevelId,
                            vipLevel: data.vipLevel,
                            status: data.status || VipStatus.ACTIVE,
                            totalRechargeCents: data.totalRechargeCents || 0,
                            totalSpentCents: data.totalSpentCents || 0,
                            vipExp: data.vipExp || 0,
                            vipUnlocked: data.vipUnlocked || false,
                            vipUnlockedAt: data.vipUnlockedAt,
                            vipExpireAt: data.vipExpireAt,
                            currentMonthCouponClaimedAt: data.currentMonthCouponClaimedAt,
                            nextLevel: undefined,
                            progressPercent: 0
                        },
                        loading: false
                    });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch VIP status' });
                }
            },

            fetchLevels: async () => {
                try {
                    const data = await http.get<VipLevel[]>('/user/vip/levels');
                    if (data && data.length > 0) {
                        set({ levels: data });
                    }
                } catch {
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
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to claim monthly coupon' });
                    throw err;
                }
            },

            calculateVipDiscount: (priceCents) => {
                const { userVip, levels } = get();

                if (!userVip?.vipLevel || !userVip.vipUnlocked) {
                    return {
                        originalPrice: priceCents / 100,
                        originalPriceCents: priceCents,
                        discountRate: 1,
                        discountAmount: 0,
                        discountAmountCents: 0,
                        finalPrice: priceCents / 100,
                        finalPriceCents: priceCents
                    };
                }

                const discountRate = userVip.vipLevel.orderDiscount;

                if (discountRate >= 1) {
                    return {
                        originalPrice: priceCents / 100,
                        originalPriceCents: priceCents,
                        discountRate: 1,
                        discountAmount: 0,
                        discountAmountCents: 0,
                        finalPrice: priceCents / 100,
                        finalPriceCents: priceCents
                    };
                }

                const discountCents = Math.floor(priceCents * (1 - discountRate));
                const finalCents = priceCents - discountCents;

                return {
                    originalPrice: priceCents / 100,
                    originalPriceCents: priceCents,
                    discountRate,
                    discountAmount: discountCents / 100,
                    discountAmountCents: discountCents,
                    finalPrice: finalCents / 100,
                    finalPriceCents: finalCents
                };
            },

            getLevelBySlug: (slug) => {
                const { levels } = get();
                return levels.find(l => l.slug === slug);
            },

            getLevelByTitle: (title) => {
                const { levels } = get();
                return levels.find(l => l.title === title);
            }
        }),
        {
            name: 'vip-storage',
            partialize: (state) => ({
                // Only persist level configs (public data)
                levels: state.levels,
                // Persist only non-sensitive VIP info for UI display
                userVip: state.userVip ? {
                    userId: state.userVip.userId,
                    vipLevelId: state.userVip.vipLevelId,
                    vipLevel: state.userVip.vipLevel,
                    status: state.userVip.status,
                    vipUnlocked: state.userVip.vipUnlocked,
                    progressPercent: state.userVip.progressPercent
                    // Excluded: totalRechargeCents, totalSpentCents, vipExp, dates
                } : null
            })
        }
    )
);
