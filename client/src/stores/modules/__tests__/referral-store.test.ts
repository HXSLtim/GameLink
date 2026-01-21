/**
 * Referral Store Tests
 * Tests for referral system and rewards
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useReferralStore, ReferralStatus } from '../referral-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
};

describe('Referral Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useReferralStore.setState({
            referralInfo: null,
            records: [],
            shareContent: null,
            rules: [],
            loading: false,
            error: null,
        });
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useReferralStore.getState();
            expect(state.referralInfo).toBeNull();
            expect(state.records).toEqual([]);
            expect(state.shareContent).toBeNull();
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchReferralInfo', () => {
        it('should fetch referral info successfully', async () => {
            const mockInfo = {
                userId: 100,
                referralCode: 'ABC123',
                totalReferrals: 5,
                successfulReferrals: 3,
                totalRewardsCents: 15000,
                rewardPerReferral: 5000,
                refereeReward: 1000,
            };
            mockHttp.get.mockResolvedValueOnce(mockInfo);

            await useReferralStore.getState().fetchReferralInfo();

            const state = useReferralStore.getState();
            expect(state.referralInfo).toEqual(mockInfo);
            expect(state.loading).toBe(false);
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Failed'));

            await useReferralStore.getState().fetchReferralInfo();

            const state = useReferralStore.getState();
            expect(state.error).toBe('Failed');
            expect(state.loading).toBe(false);
        });
    });

    describe('fetchReferralRecords', () => {
        it('should fetch records successfully', async () => {
            const mockRecords = {
                items: [
                    {
                        id: 1,
                        status: ReferralStatus.COMPLETED,
                        rewardCents: 5000,
                    },
                    {
                        id: 2,
                        status: ReferralStatus.PENDING,
                        rewardCents: 5000,
                    },
                ],
            };
            mockHttp.get.mockResolvedValueOnce(mockRecords);

            await useReferralStore.getState().fetchReferralRecords(1);

            expect(useReferralStore.getState().records).toHaveLength(2);
            expect(mockHttp.get).toHaveBeenCalledWith('/user/referral/records', {
                params: { page: 1, pageSize: 20 }
            });
        });

        it('should handle empty items', async () => {
            mockHttp.get.mockResolvedValueOnce({ items: null });

            await useReferralStore.getState().fetchReferralRecords();

            expect(useReferralStore.getState().records).toEqual([]);
        });
    });

    describe('getRewardsSummary', () => {
        it('should calculate summary correctly', () => {
            const now = new Date();
            const thisMonth = new Date(now.getFullYear(), now.getMonth(), 1);

            useReferralStore.setState({
                records: [
                    {
                        id: 1,
                        referrerId: 100,
                        refereeId: 101,
                        refereeNickname: 'User1',
                        status: ReferralStatus.REWARDED,
                        rewardCents: 1000,
                        rewardedAt: thisMonth.toISOString(),
                        createdAt: thisMonth.toISOString(),
                    },
                    {
                        id: 2,
                        referrerId: 100,
                        refereeId: 102,
                        refereeNickname: 'User2',
                        status: ReferralStatus.COMPLETED,
                        rewardCents: 2000,
                        createdAt: thisMonth.toISOString(),
                    },
                ],
            });

            const summary = useReferralStore.getState().getRewardsSummary();

            expect(summary.total).toBe(10);
            expect(summary.pending).toBe(20);
            expect(summary.thisMonth).toBe(10);
        });

        it('should handle empty records', () => {
            useReferralStore.setState({ records: [] });

            const summary = useReferralStore.getState().getRewardsSummary();

            expect(summary.total).toBe(0);
            expect(summary.pending).toBe(0);
            expect(summary.thisMonth).toBe(0);
        });
    });

    describe('ReferralStatus', () => {
        it('should have all status values', () => {
            expect(ReferralStatus.PENDING).toBe('pending');
            expect(ReferralStatus.COMPLETED).toBe('completed');
            expect(ReferralStatus.REWARDED).toBe('rewarded');
            expect(ReferralStatus.EXPIRED).toBe('expired');
        });
    });
});
