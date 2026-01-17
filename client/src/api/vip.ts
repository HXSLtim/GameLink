/**
 * VIP API
 * Handles VIP levels, benefits, upgrades
 */

import { http } from '@/lib/http';
import type {
    VipLevel,
    VipBenefit,
    UserVipInfo
} from '@/types/api';

export const vipApi = {
    /**
     * Get all VIP levels
     */
    getLevels: () =>
        http.get<VipLevel[]>('/vip/levels'),

    /**
     * Get VIP level detail
     */
    getLevel: (levelId: number) =>
        http.get<VipLevel>(`/vip/level/${levelId}`),

    /**
     * Get user VIP info
     */
    getUserVipInfo: () =>
        http.get<UserVipInfo>('/vip/user'),

    /**
     * Get VIP benefits
     */
    getBenefits: (levelId?: number) =>
        http.get<VipBenefit[]>('/vip/benefits', { params: { levelId } }),

    /**
     * Claim monthly coupon (VIP benefit)
     */
    claimMonthlyCoupon: () =>
        http.post<{ couponId: number }>('/vip/claim-monthly-coupon'),

    /**
     * Get VIP upgrade requirements
     */
    getUpgradeRequirements: (targetLevelId: number) =>
        http.get<{
            currentExp: number;
            requiredExp: number;
            canUpgrade: boolean;
        }>(`/vip/upgrade-requirements/${targetLevelId}`),
};
