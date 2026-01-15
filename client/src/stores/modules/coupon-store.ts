import { create } from 'zustand';
import { http } from '@/lib/http';

// ============ Enums ============

export const CouponType = {
    DISCOUNT: 'discount',      // 折扣券
    AMOUNT: 'amount',          // 满减券
    FREE: 'free'               // 免单券
} as const;

export type CouponType = typeof CouponType[keyof typeof CouponType];

export const CouponStatus = {
    AVAILABLE: 'available',    // 可用
    LOCKED: 'locked',          // 已锁定 (下单中)
    USED: 'used',              // 已使用
    EXPIRED: 'expired'         // 已过期
} as const;

export type CouponStatus = typeof CouponStatus[keyof typeof CouponStatus];

// ============ Interfaces ============

export interface CouponTemplate {
    id: number;
    name: string;
    type: CouponType;

    // 优惠规则
    discountRate?: number;       // 折扣率 (折扣券)
    amountCents?: number;        // 优惠金额 (满减券)
    minOrderCents: number;       // 最低订单金额

    // 使用限制
    applicableGames?: number[];  // 适用游戏
    applicablePlayers?: number[]; // 适用陪玩师

    // 有效期
    validDays: number;           // 领取后有效天数
    startAt?: string;
    endAt?: string;

    // 库存
    totalCount: number;
    remainingCount: number;
}

export interface Coupon {
    id: number;
    templateId: number;
    userId: number;

    // 优惠信息 (从模板复制)
    name: string;
    type: CouponType;
    discountRate?: number;
    amountCents?: number;
    minOrderCents: number;

    // 使用限制
    applicableGames?: number[];
    applicablePlayers?: number[];

    // 状态
    status: CouponStatus;
    lockedOrderId?: number;
    usedOrderId?: number;

    // 有效期
    validFrom: string;
    validUntil: string;

    createdAt: string;
}

export interface CouponDiscountResult {
    applicable: boolean;
    reason?: string;
    discountAmount: number;
    finalAmount: number;
}

// ============ State & Actions ============

export interface CouponState {
    myCoupons: Coupon[];
    availableCoupons: CouponTemplate[];
    couponCounts: {
        available: number;
        used: number;
        expired: number;
    };
    loading: boolean;
    error: string | null;
}

export interface CouponActions {
    fetchMyCoupons: (status?: CouponStatus) => Promise<void>;
    fetchAvailableCoupons: () => Promise<void>;
    fetchCouponCounts: () => Promise<void>;
    claimCoupon: (templateId: number) => Promise<void>;
    getCouponDetail: (id: number) => Promise<Coupon | null>;
    getApplicableCoupons: (orderAmountCents: number, gameId?: number) => Coupon[];
    calculateCouponDiscount: (orderAmountCents: number, coupon: Coupon) => CouponDiscountResult;
}

// ============ Store ============

export const useCouponStore = create<CouponState & CouponActions>((set, get) => ({
    myCoupons: [],
    availableCoupons: [],
    couponCounts: {
        available: 0,
        used: 0,
        expired: 0
    },
    loading: false,
    error: null,

    fetchMyCoupons: async (status) => {
        set({ loading: true, error: null });
        try {
            const params: Record<string, any> = {};
            if (status) params.status = status;

            const data = await http.get<Coupon[]>('/user/coupons', { params });
            set({ myCoupons: data, loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to fetch coupons' });
        }
    },

    fetchAvailableCoupons: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<CouponTemplate[]>('/coupons/available');
            set({ availableCoupons: data, loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to fetch available coupons' });
        }
    },

    fetchCouponCounts: async () => {
        try {
            const data = await http.get<{ available: number; used: number; expired: number }>('/user/coupons/count');
            set({ couponCounts: data });
        } catch (err: any) {
            console.error('Failed to fetch coupon counts:', err);
        }
    },

    claimCoupon: async (templateId) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/user/coupons/${templateId}/claim`);

            // 刷新列表
            await Promise.all([
                get().fetchMyCoupons(),
                get().fetchAvailableCoupons(),
                get().fetchCouponCounts()
            ]);

            set({ loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to claim coupon' });
            throw err;
        }
    },

    getCouponDetail: async (id) => {
        try {
            const data = await http.get<Coupon>(`/user/coupons/${id}`);
            return data;
        } catch (err: any) {
            console.error('Failed to get coupon detail:', err);
            return null;
        }
    },

    getApplicableCoupons: (orderAmountCents, gameId) => {
        const { myCoupons } = get();
        return myCoupons.filter(coupon => {
            // 检查状态
            if (coupon.status !== CouponStatus.AVAILABLE) return false;

            // 检查最低消费
            if (orderAmountCents < coupon.minOrderCents) return false;

            // 检查适用游戏
            if (coupon.applicableGames?.length && gameId) {
                if (!coupon.applicableGames.includes(gameId)) return false;
            }

            // 检查有效期
            const now = new Date();
            const validFrom = new Date(coupon.validFrom);
            const validUntil = new Date(coupon.validUntil);
            if (now < validFrom || now > validUntil) return false;

            return true;
        });
    },

    calculateCouponDiscount: (orderAmountCents, coupon) => {
        // 检查最低消费
        if (orderAmountCents < coupon.minOrderCents) {
            return {
                applicable: false,
                reason: `订单金额需满 ¥${(coupon.minOrderCents / 100).toFixed(2)}`,
                discountAmount: 0,
                finalAmount: orderAmountCents
            };
        }

        let discountAmount = 0;

        switch (coupon.type) {
            case CouponType.DISCOUNT:
                discountAmount = Math.floor(orderAmountCents * (1 - (coupon.discountRate || 1)));
                break;
            case CouponType.AMOUNT:
                discountAmount = Math.min(coupon.amountCents || 0, orderAmountCents);
                break;
            case CouponType.FREE:
                discountAmount = orderAmountCents;
                break;
        }

        return {
            applicable: true,
            discountAmount,
            finalAmount: orderAmountCents - discountAmount
        };
    }
}));
