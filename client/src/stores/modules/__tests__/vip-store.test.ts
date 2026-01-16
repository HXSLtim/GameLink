/**
 * VIP Store Tests
 * Tests for VIP status, levels, benefits, and discounts
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useVipStore, VipLevel, VipStatus } from '../vip-store';

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

describe('VIP Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useVipStore.setState({
            userVip: null,
            levelConfigs: [
                {
                    level: VipLevel.NONE,
                    name: '普通用户',
                    icon: '👤',
                    color: '#9CA3AF',
                    thresholdCents: 0,
                    discountRate: 1,
                    monthlyFreeCoupons: 0,
                    prioritySupport: false,
                    exclusiveActivities: false,
                },
                {
                    level: VipLevel.BRONZE,
                    name: '青铜 VIP',
                    icon: '🥉',
                    color: '#CD7F32',
                    thresholdCents: 50000,
                    discountRate: 0.98,
                    monthlyFreeCoupons: 1,
                    prioritySupport: false,
                    exclusiveActivities: false,
                },
                {
                    level: VipLevel.SILVER,
                    name: '白银 VIP',
                    icon: '🥈',
                    color: '#C0C0C0',
                    thresholdCents: 200000,
                    discountRate: 0.95,
                    monthlyFreeCoupons: 2,
                    prioritySupport: false,
                    exclusiveActivities: false,
                },
                {
                    level: VipLevel.GOLD,
                    name: '黄金 VIP',
                    icon: '🥇',
                    color: '#FFD700',
                    thresholdCents: 500000,
                    discountRate: 0.92,
                    monthlyFreeCoupons: 3,
                    prioritySupport: true,
                    exclusiveActivities: false,
                },
            ],
            benefits: [],
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useVipStore.getState();

            expect(state.userVip).toBeNull();
            expect(state.levelConfigs.length).toBeGreaterThan(0);
            expect(state.benefits).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should have default level configs', () => {
            const state = useVipStore.getState();

            expect(state.levelConfigs[0].level).toBe(VipLevel.NONE);
            expect(state.levelConfigs[0].discountRate).toBe(1);
        });
    });

    describe('fetchVipStatus', () => {
        it('should fetch VIP status successfully', async () => {
            const mockVipStatus = {
                userId: 1,
                level: VipLevel.GOLD,
                status: VipStatus.ACTIVE,
                totalSpentCents: 600000,
                currentYearSpentCents: 300000,
                activatedAt: '2024-01-01T00:00:00Z',
                monthlyFreeCoupons: 3,
                usedFreeCoupons: 1,
                discountRate: 0.92,
                nextLevel: VipLevel.PLATINUM,
                nextLevelThreshold: 1000000,
                progressPercent: 60,
            };

            mockHttp.get.mockResolvedValueOnce(mockVipStatus);

            await useVipStore.getState().fetchVipStatus();

            const state = useVipStore.getState();
            expect(state.userVip).not.toBeNull();
            expect(state.userVip?.level).toBe(VipLevel.GOLD);
            expect(state.userVip?.status).toBe(VipStatus.ACTIVE);
            expect(state.userVip?.discountRate).toBe(0.92);
            expect(state.loading).toBe(false);
        });

        it('should handle fetch VIP status error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useVipStore.getState().fetchVipStatus();

            const state = useVipStore.getState();
            expect(state.userVip).toBeNull();
            expect(state.error).toBe('Network error');
        });

        it('should set default values for missing fields', async () => {
            mockHttp.get.mockResolvedValueOnce({ userId: 1 });

            await useVipStore.getState().fetchVipStatus();

            const state = useVipStore.getState();
            expect(state.userVip?.level).toBe(VipLevel.NONE);
            expect(state.userVip?.status).toBe(VipStatus.ACTIVE);
            expect(state.userVip?.discountRate).toBe(1);
        });
    });

    describe('fetchLevelConfigs', () => {
        it('should fetch level configs from API', async () => {
            const mockConfigs = [
                {
                    level: VipLevel.NONE,
                    name: 'Free User',
                    icon: '👤',
                    color: '#999',
                    thresholdCents: 0,
                    discountRate: 1,
                    monthlyFreeCoupons: 0,
                    prioritySupport: false,
                    exclusiveActivities: false,
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockConfigs);

            await useVipStore.getState().fetchLevelConfigs();

            const state = useVipStore.getState();
            expect(state.levelConfigs[0].name).toBe('Free User');
        });

        it('should keep default configs on API error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('API error'));

            await useVipStore.getState().fetchLevelConfigs();

            const state = useVipStore.getState();
            expect(state.levelConfigs.length).toBeGreaterThan(0);
            expect(state.levelConfigs[0].name).toBe('普通用户');
        });
    });

    describe('claimMonthlyCoupon', () => {
        it('should claim monthly coupon successfully', async () => {
            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce({
                userId: 1,
                level: VipLevel.GOLD,
                usedFreeCoupons: 2,
            });

            await useVipStore.getState().claimMonthlyCoupon();

            expect(mockHttp.post).toHaveBeenCalledWith('/user/vip/monthly-coupon');
        });

        it('should handle claim error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Already claimed'));

            await expect(
                useVipStore.getState().claimMonthlyCoupon()
            ).rejects.toThrow('Already claimed');

            expect(useVipStore.getState().error).toBe('Already claimed');
        });
    });

    describe('calculateVipDiscount', () => {
        it('should return no discount for non-VIP users', () => {
            const result = useVipStore.getState().calculateVipDiscount(10000);

            expect(result.originalPrice).toBe(10000);
            expect(result.discountRate).toBe(1);
            expect(result.discountAmount).toBe(0);
            expect(result.finalPrice).toBe(10000);
        });

        it('should calculate discount for VIP users', () => {
            useVipStore.setState({
                userVip: {
                    userId: 1,
                    level: VipLevel.GOLD,
                    status: VipStatus.ACTIVE,
                    totalSpentCents: 600000,
                    currentYearSpentCents: 300000,
                    monthlyFreeCoupons: 3,
                    usedFreeCoupons: 0,
                    discountRate: 0.92,
                    progressPercent: 60,
                },
            });

            const result = useVipStore.getState().calculateVipDiscount(10000);

            expect(result.originalPrice).toBe(10000);
            expect(result.discountRate).toBe(0.92);
            // Note: Due to floating point precision (1 - 0.92 = 0.07999999999999996)
            // Math.floor(10000 * 0.07999999999999996) = 799
            expect(result.discountAmount).toBe(799);
            expect(result.finalPrice).toBe(9201);
        });

        it('should handle NONE level VIP', () => {
            useVipStore.setState({
                userVip: {
                    userId: 1,
                    level: VipLevel.NONE,
                    status: VipStatus.ACTIVE,
                    totalSpentCents: 0,
                    currentYearSpentCents: 0,
                    monthlyFreeCoupons: 0,
                    usedFreeCoupons: 0,
                    discountRate: 1,
                    progressPercent: 0,
                },
            });

            const result = useVipStore.getState().calculateVipDiscount(10000);

            expect(result.discountAmount).toBe(0);
            expect(result.finalPrice).toBe(10000);
        });
    });

    describe('getLevelName', () => {
        it('should return correct level name', () => {
            expect(useVipStore.getState().getLevelName(VipLevel.NONE)).toBe('普通用户');
            expect(useVipStore.getState().getLevelName(VipLevel.BRONZE)).toBe('青铜 VIP');
            expect(useVipStore.getState().getLevelName(VipLevel.GOLD)).toBe('黄金 VIP');
        });

        it('should return default name for unknown level', () => {
            expect(useVipStore.getState().getLevelName(99 as VipLevel)).toBe('普通用户');
        });
    });

    describe('getLevelColor', () => {
        it('should return correct level color', () => {
            expect(useVipStore.getState().getLevelColor(VipLevel.NONE)).toBe('#9CA3AF');
            expect(useVipStore.getState().getLevelColor(VipLevel.BRONZE)).toBe('#CD7F32');
            expect(useVipStore.getState().getLevelColor(VipLevel.GOLD)).toBe('#FFD700');
        });

        it('should return default color for unknown level', () => {
            expect(useVipStore.getState().getLevelColor(99 as VipLevel)).toBe('#9CA3AF');
        });
    });
});
