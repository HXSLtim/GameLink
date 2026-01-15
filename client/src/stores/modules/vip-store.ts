import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

export interface VipBenefit {
    id: string;
    name: string;
    description?: string;
    icon: string;
}

export interface VipLevel {
    id: number;
    level: number;
    name: string;
    icon: string;
    color: string;
}

export interface VipState {
    vipUnlocked: boolean;
    currentLevel: VipLevel | null;
    currentExp: number;
    nextLevelExp: number;
    expProgress: number; // 0-1
    vipExpireAt: string | null;
    benefits: VipBenefit[];
    monthlyTicketsRemaining: number;
    discountRate: number;

    loading: boolean;
    error: string | null;
}

export interface VipActions {
    fetchVipInfo: () => Promise<void>;
    purchaseSubscription: (months: number) => Promise<void>;
}

export const useVipStore = create<VipState & VipActions>()(
    persist(
        (set, get) => ({
            vipUnlocked: false,
            currentLevel: null,
            currentExp: 0,
            nextLevelExp: 0,
            expProgress: 0,
            vipExpireAt: null,
            benefits: [],
            monthlyTicketsRemaining: 0,
            discountRate: 0,

            loading: false,
            error: null,

            fetchVipInfo: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<any>('/user/vip/info');

                    set({
                        vipUnlocked: data.vipUnlocked,
                        currentLevel: data.currentLevel, // backend object matches interface largely
                        currentExp: data.currentExp,
                        nextLevelExp: data.nextLevelExp,
                        expProgress: data.expProgress,
                        vipExpireAt: data.vipExpireAt,
                        benefits: data.benefits || [],
                        monthlyTicketsRemaining: data.monthlyTicketsRemaining,
                        discountRate: data.discountRate,
                        loading: false
                    });
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Failed to fetch VIP info' });
                }
            },

            purchaseSubscription: async (months) => {
                set({ loading: true, error: null });
                try {
                    await http.post('/user/vip/subscribe', { months });
                    // Refresh info after purchase
                    await get().fetchVipInfo();
                } catch (err: any) {
                    set({ loading: false, error: err.message || 'Subscription failed' });
                    throw err;
                }
            }
        }),
        {
            name: 'vip-storage',
            partialize: (state) => ({
                vipUnlocked: state.vipUnlocked,
                currentLevel: state.currentLevel,
                vipExpireAt: state.vipExpireAt
            })
        }
    )
);
