/**
 * Coupon Store Tests
 * Tests for coupon management, claiming, and discount calculation
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useCouponStore, CouponType, CouponStatus } from '../coupon-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

describe('Coupon Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useCouponStore.setState({
            myCoupons: [],
            availableCoupons: [],
            couponCounts: {
                available: 0,
                used: 0,
                expired: 0,
            },
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useCouponStore.getState();

            expect(state.myCoupons).toEqual([]);
            expect(state.availableCoupons).toEqual([]);
            expect(state.couponCounts.available).toBe(0);
            expect(state.couponCounts.used).toBe(0);
            expect(state.couponCounts.expired).toBe(0);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchMyCoupons', () => {
        it('should fetch my coupons successfully', async () => {
            const mockCoupons = [
                {
                    id: 1,
                    templateId: 1,
                    userId: 1,
                    name: '新用户优惠券',
                    type: CouponType.AMOUNT,
                    amountCents: 1000,
                    minOrderCents: 5000,
                    status: CouponStatus.AVAILABLE,
                    validFrom: '2024-01-01T00:00:00Z',
                    validUntil: '2024-12-31T23:59:59Z',
                    createdAt: '2024-01-01T00:00:00Z',
                },
                {
                    id: 2,
                    templateId: 2,
                    userId: 1,
                    name: '95折优惠券',
                    type: CouponType.DISCOUNT,
                    discountRate: 0.95,
                    minOrderCents: 10000,
                    status: CouponStatus.AVAILABLE,
                    validFrom: '2024-01-01T00:00:00Z',
                    validUntil: '2024-12-31T23:59:59Z',
                    createdAt: '2024-01-01T00:00:00Z',
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockCoupons);

            await useCouponStore.getState().fetchMyCoupons();

            const state = useCouponStore.getState();
            expect(state.myCoupons).toHaveLength(2);
            expect(state.myCoupons[0].name).toBe('新用户优惠券');
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should fetch coupons with status filter', async () => {
            mockHttp.get.mockResolvedValueOnce([]);

            await useCouponStore.getState().fetchMyCoupons(CouponStatus.USED);

            expect(mockHttp.get).toHaveBeenCalledWith('/user/coupons', {
                params: { status: CouponStatus.USED },
            });
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useCouponStore.getState().fetchMyCoupons();

            const state = useCouponStore.getState();
            expect(state.myCoupons).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });
    });

    describe('fetchAvailableCoupons', () => {
        it('should fetch available coupon templates', async () => {
            const mockTemplates = [
                {
                    id: 1,
                    name: '新用户专享',
                    type: CouponType.AMOUNT,
                    amountCents: 2000,
                    minOrderCents: 10000,
                    validDays: 30,
                    totalCount: 1000,
                    remainingCount: 500,
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockTemplates);

            await useCouponStore.getState().fetchAvailableCoupons();

            const state = useCouponStore.getState();
            expect(state.availableCoupons).toHaveLength(1);
            expect(state.availableCoupons[0].name).toBe('新用户专享');
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('API error'));

            await useCouponStore.getState().fetchAvailableCoupons();

            const state = useCouponStore.getState();
            expect(state.error).toBe('API error');
        });
    });

    describe('fetchCouponCounts', () => {
        it('should fetch coupon counts', async () => {
            mockHttp.get.mockResolvedValueOnce({
                available: 5,
                used: 10,
                expired: 3,
            });

            await useCouponStore.getState().fetchCouponCounts();

            const state = useCouponStore.getState();
            expect(state.couponCounts.available).toBe(5);
            expect(state.couponCounts.used).toBe(10);
            expect(state.couponCounts.expired).toBe(3);
        });

        it('should handle error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { });
            mockHttp.get.mockRejectedValueOnce(new Error('API error'));

            await useCouponStore.getState().fetchCouponCounts();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('claimCoupon', () => {
        it('should claim coupon successfully', async () => {
            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce([]); // fetchMyCoupons
            mockHttp.get.mockResolvedValueOnce([]); // fetchAvailableCoupons
            mockHttp.get.mockResolvedValueOnce({ available: 1, used: 0, expired: 0 }); // fetchCouponCounts

            await useCouponStore.getState().claimCoupon(1);

            expect(mockHttp.post).toHaveBeenCalledWith('/user/coupons/1/claim');
        });

        it('should handle claim error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Already claimed'));

            await expect(
                useCouponStore.getState().claimCoupon(1)
            ).rejects.toThrow('Already claimed');

            const state = useCouponStore.getState();
            expect(state.error).toBe('Already claimed');
        });
    });

    describe('getCouponDetail', () => {
        it('should get coupon detail', async () => {
            const mockCoupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '测试优惠券',
                type: CouponType.AMOUNT,
                amountCents: 1000,
                minOrderCents: 5000,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            mockHttp.get.mockResolvedValueOnce(mockCoupon);

            const result = await useCouponStore.getState().getCouponDetail(1);

            expect(result).toEqual(mockCoupon);
            expect(mockHttp.get).toHaveBeenCalledWith('/user/coupons/1');
        });

        it('should return null on error', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { });
            mockHttp.get.mockRejectedValueOnce(new Error('Not found'));

            const result = await useCouponStore.getState().getCouponDetail(999);

            expect(result).toBeNull();
            consoleSpy.mockRestore();
        });
    });

    describe('getApplicableCoupons', () => {
        const now = new Date();
        const validFrom = new Date(now.getTime() - 86400000).toISOString(); // Yesterday
        const validUntil = new Date(now.getTime() + 86400000).toISOString(); // Tomorrow

        beforeEach(() => {
            useCouponStore.setState({
                myCoupons: [
                    {
                        id: 1,
                        templateId: 1,
                        userId: 1,
                        name: '满50减10',
                        type: CouponType.AMOUNT,
                        amountCents: 1000,
                        minOrderCents: 5000,
                        status: CouponStatus.AVAILABLE,
                        validFrom,
                        validUntil,
                        createdAt: validFrom,
                    },
                    {
                        id: 2,
                        templateId: 2,
                        userId: 1,
                        name: '满100减20',
                        type: CouponType.AMOUNT,
                        amountCents: 2000,
                        minOrderCents: 10000,
                        status: CouponStatus.AVAILABLE,
                        validFrom,
                        validUntil,
                        createdAt: validFrom,
                    },
                    {
                        id: 3,
                        templateId: 3,
                        userId: 1,
                        name: '已使用券',
                        type: CouponType.AMOUNT,
                        amountCents: 500,
                        minOrderCents: 3000,
                        status: CouponStatus.USED,
                        validFrom,
                        validUntil,
                        createdAt: validFrom,
                    },
                ],
            });
        });

        it('should return applicable coupons for order amount', () => {
            const result = useCouponStore.getState().getApplicableCoupons(8000);

            expect(result).toHaveLength(1);
            expect(result[0].name).toBe('满50减10');
        });

        it('should return multiple applicable coupons', () => {
            const result = useCouponStore.getState().getApplicableCoupons(15000);

            expect(result).toHaveLength(2);
        });

        it('should filter out used coupons', () => {
            const result = useCouponStore.getState().getApplicableCoupons(5000);

            const usedCoupon = result.find(c => c.status === CouponStatus.USED);
            expect(usedCoupon).toBeUndefined();
        });

        it('should filter by game ID', () => {
            useCouponStore.setState({
                myCoupons: [
                    {
                        id: 1,
                        templateId: 1,
                        userId: 1,
                        name: '游戏专属券',
                        type: CouponType.AMOUNT,
                        amountCents: 1000,
                        minOrderCents: 5000,
                        applicableGames: [1, 2],
                        status: CouponStatus.AVAILABLE,
                        validFrom,
                        validUntil,
                        createdAt: validFrom,
                    },
                ],
            });

            const result1 = useCouponStore.getState().getApplicableCoupons(10000, 1);
            expect(result1).toHaveLength(1);

            const result2 = useCouponStore.getState().getApplicableCoupons(10000, 99);
            expect(result2).toHaveLength(0);
        });

        it('should return empty for insufficient order amount', () => {
            const result = useCouponStore.getState().getApplicableCoupons(1000);

            expect(result).toHaveLength(0);
        });
    });

    describe('calculateCouponDiscount', () => {
        it('should calculate amount coupon discount', () => {
            const coupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '满50减10',
                type: CouponType.AMOUNT,
                amountCents: 1000,
                minOrderCents: 5000,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            const result = useCouponStore.getState().calculateCouponDiscount(10000, coupon);

            expect(result.applicable).toBe(true);
            expect(result.discountAmount).toBe(1000);
            expect(result.finalAmount).toBe(9000);
        });

        it('should calculate discount coupon', () => {
            const coupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '95折券',
                type: CouponType.DISCOUNT,
                discountRate: 0.95,
                minOrderCents: 5000,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            const result = useCouponStore.getState().calculateCouponDiscount(10000, coupon);

            expect(result.applicable).toBe(true);
            // 10000 * (1 - 0.95) = 500
            expect(result.discountAmount).toBe(500);
            expect(result.finalAmount).toBe(9500);
        });

        it('should calculate free coupon', () => {
            const coupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '免单券',
                type: CouponType.FREE,
                minOrderCents: 0,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            const result = useCouponStore.getState().calculateCouponDiscount(10000, coupon);

            expect(result.applicable).toBe(true);
            expect(result.discountAmount).toBe(10000);
            expect(result.finalAmount).toBe(0);
        });

        it('should reject coupon for insufficient order amount', () => {
            const coupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '满100减20',
                type: CouponType.AMOUNT,
                amountCents: 2000,
                minOrderCents: 10000,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            const result = useCouponStore.getState().calculateCouponDiscount(5000, coupon);

            expect(result.applicable).toBe(false);
            expect(result.reason).toContain('100.00');
            expect(result.discountAmount).toBe(0);
            expect(result.finalAmount).toBe(5000);
        });

        it('should cap amount discount at order total', () => {
            const coupon = {
                id: 1,
                templateId: 1,
                userId: 1,
                name: '满10减50',
                type: CouponType.AMOUNT,
                amountCents: 5000,
                minOrderCents: 1000,
                status: CouponStatus.AVAILABLE,
                validFrom: '2024-01-01',
                validUntil: '2024-12-31',
                createdAt: '2024-01-01',
            };

            const result = useCouponStore.getState().calculateCouponDiscount(3000, coupon);

            expect(result.applicable).toBe(true);
            expect(result.discountAmount).toBe(3000); // Capped at order amount
            expect(result.finalAmount).toBe(0);
        });
    });
});
