import { create } from 'zustand';
import { persist } from 'zustand/middleware';
// import { http } from '@/lib/http';

export interface VipBenefit {
    id: string;
    name: string;
    description: string;
    icon: string;
}

export interface VipLevel {
    level: number;
    name: string;
    requiredExp: number;
    benefits: VipBenefit[];
}

export interface VipState {
    currentLevel: number;
    currentExp: number;
    nextLevelExp: number;
    isSubscriber: boolean; // Monthly subscription (SVIP)
    subscriptionExpireAt?: string;
    benefits: VipBenefit[];
    loading: boolean;
    error: string | null;
}

export interface VipActions {
    fetchVipInfo: () => Promise<void>;
    purchaseSubscription: (months: number) => Promise<void>;
}

export const useVipStore = create<VipState & VipActions>()(
    persist(
        (set) => ({
            currentLevel: 0,
            currentExp: 0,
            nextLevelExp: 100,
            isSubscriber: false,
            benefits: [],
            loading: false,
            error: null,

            fetchVipInfo: async () => {
                set({ loading: true, error: null });
                try {
                    // await http.get('/vip/info');
                    await new Promise(resolve => setTimeout(resolve, 500));

                    set({
                        currentLevel: 1,
                        currentExp: 450,
                        nextLevelExp: 1000,
                        isSubscriber: true,
                        subscriptionExpireAt: new Date(Date.now() + 86400000 * 15).toISOString(), // 15 days left
                        benefits: [
                            { id: 'b1', name: 'No Ads', description: 'Enjoy ad-free experience', icon: 'block' },
                            { id: 'b2', name: 'Priority Support', description: '24/7 dedicated support', icon: 'headset' },
                        ],
                        loading: false
                    });
                } catch (err: any) {
                    set({ loading: false, error: err.message });
                }
            },

            purchaseSubscription: async (months) => {
                set({ loading: true, error: null });
                try {
                    // await http.post('/vip/subscribe', { months });
                    await new Promise(resolve => setTimeout(resolve, 1000));

                    // Update mock state
                    set(state => ({
                        isSubscriber: true,
                        // Extend time roughly
                        subscriptionExpireAt: new Date(new Date(state.subscriptionExpireAt || Date.now()).getTime() + months * 30 * 86400000).toISOString(),
                        loading: false
                    }));
                } catch (err: any) {
                    set({ loading: false, error: err.message });
                }
            }
        }),
        {
            name: 'vip-storage',
            partialize: (state) => ({
                currentLevel: state.currentLevel,
                isSubscriber: state.isSubscriber,
                subscriptionExpireAt: state.subscriptionExpireAt
            })
        }
    )
);
