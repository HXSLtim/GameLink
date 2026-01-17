/**
 * Referral API
 * Handles referral codes, rewards, referral statistics
 */

import { http } from '@/lib/http';
import type {
    ReferralInfo,
    ReferralRecord,
    PaginatedResponse
} from '@/types/api';

export const referralApi = {
    /**
     * Get user's referral info
     */
    getInfo: () =>
        http.get<ReferralInfo>('/referral/info'),

    /**
     * Get referral code
     */
    getCode: () =>
        http.get<{ code: string }>('/referral/code'),

    /**
     * Generate new referral code
     */
    generateCode: () =>
        http.post<{ code: string }>('/referral/generate-code'),

    /**
     * Get referral records
     */
    getRecords: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<ReferralRecord>>('/referral/records', { params }),

    /**
     * Get referral statistics
     */
    getStats: () =>
        http.get<{
            totalReferrals: number;
            totalRewards: number;
            pendingRewards: number;
            successfulReferrals: number;
        }>('/referral/stats'),

    /**
     * Claim referral rewards
     */
    claimRewards: () =>
        http.post<{
            amount: number;
            claimedAt: string;
        }>('/referral/claim-rewards'),

    /**
     * Get referral leaderboard
     */
    getLeaderboard: (limit: number = 10) =>
        http.get<Array<{
            userId: number;
            username: string;
            avatar: string;
            referralCount: number;
            rank: number;
        }>>('/referral/leaderboard', { params: { limit } }),
};
