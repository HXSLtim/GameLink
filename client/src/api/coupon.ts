/**
 * Coupon API
 * Handles coupon listing, claiming, usage
 */

import { http } from '@/lib/http';
import type {
    Coupon,
    CouponListParams,
    PaginatedResponse
} from '@/types/api';

export const couponApi = {
    /**
     * Get available coupons (coupon center)
     */
    getAvailable: (params?: {
        page?: number;
        pageSize?: number;
        type?: string;
    }) =>
        http.get<PaginatedResponse<Coupon>>('/coupon/available', { params }),

    /**
     * Get user's coupons
     */
    getMyCoupons: (params: CouponListParams) =>
        http.get<PaginatedResponse<Coupon>>('/coupon/my', { params }),

    /**
     * Get coupon detail
     */
    get: (id: number) =>
        http.get<Coupon>(`/coupon/${id}`),

    /**
     * Claim coupon
     */
    claim: (couponId: number) =>
        http.post<{ userCouponId: number }>('/coupon/claim', { couponId }),

    /**
     * Claim coupon by code
     */
    claimByCode: (code: string) =>
        http.post<{ userCouponId: number }>('/coupon/claim-by-code', { code }),

    /**
     * Check if coupon is applicable to order
     */
    checkApplicable: (couponId: number, orderData: {
        serviceItemId: number;
        amount: number;
    }) =>
        http.post<{
            applicable: boolean;
            discountAmount: number;
            reason?: string;
        }>('/coupon/check-applicable', { couponId, ...orderData }),

    /**
     * Get best coupon for order
     */
    getBestCoupon: (orderData: {
        serviceItemId: number;
        amount: number;
    }) =>
        http.post<Coupon | null>('/coupon/best', orderData),
};
