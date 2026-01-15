import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const ReferralStatus = {
    PENDING: 'pending',         // 待完成首单
    COMPLETED: 'completed',     // 已完成
    REWARDED: 'rewarded',       // 已发放奖励
    EXPIRED: 'expired'          // 已过期
} as const;

export type ReferralStatus = typeof ReferralStatus[keyof typeof ReferralStatus];

// ============ Interfaces ============

export interface ReferralInfo {
    userId: number;
    referralCode: string;        // 推荐码

    // 统计
    totalReferrals: number;      // 总推荐人数
    successfulReferrals: number; // 成功推荐 (完成首单)
    totalRewardsCents: number;   // 累计奖励

    // 奖励规则
    rewardPerReferral: number;   // 每成功推荐奖励 (分)
    refereeReward: number;       // 被推荐人奖励 (分)
}

export interface ReferralRecord {
    id: number;
    referrerId: number;
    refereeId: number;
    refereeNickname: string;
    refereeAvatar?: string;

    // 状态
    status: ReferralStatus;
    rewardCents: number;
    rewardedAt?: string;

    createdAt: string;
}

export interface ShareContent {
    title: string;
    description: string;
    imageUrl: string;
    link: string;
}

// ============ State & Actions ============

export interface ReferralState {
    referralInfo: ReferralInfo | null;
    records: ReferralRecord[];
    shareContent: ShareContent | null;
    rules: string[];
    loading: boolean;
    error: string | null;
}

export interface ReferralActions {
    fetchReferralInfo: () => Promise<void>;
    fetchReferralRecords: (page?: number) => Promise<void>;
    fetchShareContent: () => Promise<void>;
    generateReferralLink: () => string;
    copyReferralCode: () => Promise<boolean>;
    getRewardsSummary: () => { total: number; pending: number; thisMonth: number };
}

// ============ Store ============

export const useReferralStore = create<ReferralState & ReferralActions>()(
    persist(
        (set, get) => ({
            referralInfo: null,
            records: [],
            shareContent: null,
            rules: [
                '邀请好友注册并完成首单，双方均可获得奖励',
                '被邀请人需在注册后30天内完成首单',
                '奖励将在被邀请人完成首单后7天内发放到钱包',
                '每位用户每月最多可获得10次推荐奖励',
                '推荐奖励不可与其他活动叠加使用'
            ],
            loading: false,
            error: null,

            fetchReferralInfo: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<ReferralInfo>('/user/referral');
                    set({ referralInfo: data, loading: false });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch referral info' });
                }
            },

            fetchReferralRecords: async (page = 1) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<{ items: ReferralRecord[] }>('/user/referral/records', {
                        params: { page, pageSize: 20 }
                    });
                    set({ records: data.items || [], loading: false });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch referral records' });
                }
            },

            fetchShareContent: async () => {
                try {
                    const { referralInfo } = get();
                    if (!referralInfo) return;

                    // 构建分享内容
                    const shareContent: ShareContent = {
                        title: '来 GameLink 一起玩游戏吧！',
                        description: `使用我的邀请码 ${referralInfo.referralCode} 注册，首单立享优惠！`,
                        imageUrl: '/share-banner.png',
                        link: `${window.location.origin}/register?ref=${referralInfo.referralCode}`
                    };
                    set({ shareContent });
                } catch (err) {
                    console.error('Failed to generate share content:', err);
                }
            },

            generateReferralLink: () => {
                const { referralInfo } = get();
                if (!referralInfo) return '';
                return `${window.location.origin}/register?ref=${referralInfo.referralCode}`;
            },

            copyReferralCode: async () => {
                const { referralInfo } = get();
                if (!referralInfo) return false;

                try {
                    await navigator.clipboard.writeText(referralInfo.referralCode);
                    return true;
                } catch (err) {
                    console.error('Failed to copy referral code:', err);
                    return false;
                }
            },

            getRewardsSummary: () => {
                const { records } = get();
                const now = new Date();
                const thisMonth = new Date(now.getFullYear(), now.getMonth(), 1);

                let total = 0;
                let pending = 0;
                let thisMonthRewards = 0;

                records.forEach(record => {
                    if (record.status === ReferralStatus.REWARDED) {
                        total += record.rewardCents;
                        if (record.rewardedAt && new Date(record.rewardedAt) >= thisMonth) {
                            thisMonthRewards += record.rewardCents;
                        }
                    } else if (record.status === ReferralStatus.COMPLETED) {
                        pending += record.rewardCents;
                    }
                });

                return {
                    total: total / 100,
                    pending: pending / 100,
                    thisMonth: thisMonthRewards / 100
                };
            }
        }),
        {
            name: 'referral-storage',
            partialize: (state) => ({
                referralInfo: state.referralInfo
            })
        }
    )
);
