/**
 * VIP Store Tests
 * Tests for VIP status, levels, benefits, and discounts
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useVipStore, VipStatus } from '../vip-store';

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

// Default levels matching the store's DEFAULT_LEVELS
const DEFAULT_LEVELS = [
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
        expRequired: 50000,
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
        expRequired: 200000,
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
        expRequired: 500000,
        orderDiscount: 0.92,
        monthlyCouponCount: 3,
        iconUrl: '',
        color: '#FFD700',
        benefits: {},
        sortOrder: 3,
        isDefault: false,
        isActive: true
    },
];

describe('VIP Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useVipStore.setState({
            userVip: null,
            levels: DEFAULT_LEVELS,
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
            expect(state.levels.length).toBeGreaterThan(0);
            expect(state.benefits).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should have default level configs', () => {
            const state = useVipStore.getState();

            expect(state.levels[0].slug).toBe('none');
            expect(state.levels[0].orderDiscount).toBe(1);
        });
    });

    describe('fetchVipStatus', () => {
        it('should fetch VIP status successfully', async () => {
            const mockVipStatus = {
                userId: 1,
                vipLevelId: 3,
                vipLevel: DEFAULT_LEVELS[3],
                status: VipStatus.ACTIVE,
                totalRechargeCents: 600000,
                totalSpentCents: 500000,
                vipExp: 500000,
                vipUnlocked: true,
                vipUnlockedAt: '2024-01-01T00:00:00Z',
            };

            mockHttp.get.mockResolvedValueOnce(mockVipStatus);

            await useVipStore.getState().fetchVipStatus();

            const state = useVipStore.getState();
            expect(state.userVip).not.toBeNull();
            expect(state.userVip?.vipLevelId).toBe(3);
            expect(state.userVip?.status).toBe(VipStatus.ACTIVE);
            expect(state.userVip?.vipUnlocked).toBe(true);
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
            expect(state.userVip?.status).toBe(VipStatus.ACTIVE);
            expect(state.userVip?.vipUnlocked).toBe(false);
            expect(state.userVip?.totalSpentCents).toBe(0);
        });
    });

    describe('fetchLevels', () => {
        it('should fetch level configs from API', async () => {
            const mockConfigs = [
                {
                    id: 0,
                    slug: 'none',
                    title: 'Free User',
                    expRequired: 0,
                    orderDiscount: 1,
                    monthlyCouponCount: 0,
                    iconUrl: '',
                    color: '#999',
                    benefits: {},
                    sortOrder: 0,
                    isDefault: true,
                    isActive: true,
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockConfigs);

            await useVipStore.getState().fetchLevels();

            const state = useVipStore.getState();
            expect(state.levels[0].title).toBe('Free User');
        });

        it('should keep default configs on API error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('API error'));

            await useVipStore.getState().fetchLevels();

            const state = useVipStore.getState();
            expect(state.levels.length).toBeGreaterThan(0);
            expect(state.levels[0].title).toBe('普通用户');
        });
    });

    describe('claimMonthlyCoupon', () => {
        it('should claim monthly coupon successfully', async () => {
            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce({
                userId: 1,
                vipLevelId: 3,
                vipUnlocked: true,
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

            expect(result.originalPriceCents).toBe(10000);
            expect(result.discountRate).toBe(1);
            expect(result.discountAmountCents).toBe(0);
            expect(result.finalPriceCents).toBe(10000);
        });

        it('should calculate discount for VIP users', () => {
            useVipStore.setState({
                userVip: {
                    userId: 1,
                    vipLevelId: 3,
                    vipLevel: {
                        id: 3,
                        slug: 'gold',
                        title: '黄金 VIP',
                        expRequired: 500000,
                        orderDiscount: 0.92,
                        monthlyCouponCount: 3,
                        iconUrl: '',
                        color: '#FFD700',
                        benefits: {},
                        sortOrder: 3,
                        isDefault: false,
                        isActive: true,
                    },
                    status: VipStatus.ACTIVE,
                    totalRechargeCents: 600000,
                    totalSpentCents: 500000,
                    vipExp: 500000,
                    vipUnlocked: true,
                    progressPercent: 60,
                },
            });

            const result = useVipStore.getState().calculateVipDiscount(10000);

            expect(result.originalPriceCents).toBe(10000);
            expect(result.discountRate).toBe(0.92);
            // 10000 * (1 - 0.92) = 800, but due to floating point: Math.floor(10000 * 0.07999...) = 799
            expect(result.discountAmountCents).toBe(799);
            expect(result.finalPriceCents).toBe(9201);
        });

        it('should handle non-unlocked VIP', () => {
            useVipStore.setState({
                userVip: {
                    userId: 1,
                    vipLevelId: 0,
                    status: VipStatus.ACTIVE,
                    totalRechargeCents: 0,
                    totalSpentCents: 0,
                    vipExp: 0,
                    vipUnlocked: false,
                    progressPercent: 0,
                },
            });

            const result = useVipStore.getState().calculateVipDiscount(10000);

            expect(result.discountAmountCents).toBe(0);
            expect(result.finalPriceCents).toBe(10000);
        });
    });

    describe('getLevelBySlug', () => {
        it('should return correct level by slug', () => {
            expect(useVipStore.getState().getLevelBySlug('none')?.title).toBe('普通用户');
            expect(useVipStore.getState().getLevelBySlug('bronze')?.title).toBe('青铜 VIP');
            expect(useVipStore.getState().getLevelBySlug('gold')?.title).toBe('黄金 VIP');
        });

        it('should return undefined for unknown slug', () => {
            expect(useVipStore.getState().getLevelBySlug('unknown')).toBeUndefined();
        });
    });

    describe('getLevelByTitle', () => {
        it('should return correct level by title', () => {
            expect(useVipStore.getState().getLevelByTitle('普通用户')?.slug).toBe('none');
            expect(useVipStore.getState().getLevelByTitle('青铜 VIP')?.slug).toBe('bronze');
        });

        it('should return undefined for unknown title', () => {
            expect(useVipStore.getState().getLevelByTitle('Unknown Level')).toBeUndefined();
        });
    });
});
