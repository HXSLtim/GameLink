/**
 * Player Store Tests
 * Tests for player listing, filtering, and player-side features
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { usePlayerStore, VerificationStatus, OnlineStatus } from '../player-store';

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

describe('Player Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        usePlayerStore.setState({
            players: [],
            featuredPlayers: [],
            currentPlayer: null,
            filters: {
                onlineOnly: false,
                sortBy: 'rating',
            },
            pagination: {
                page: 1,
                pageSize: 20,
                hasMore: true,
                total: 0,
            },
            myProfile: null,
            applicationStatus: null,
            earnings: null,
            earningsRecords: [],
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = usePlayerStore.getState();

            expect(state.players).toEqual([]);
            expect(state.featuredPlayers).toEqual([]);
            expect(state.currentPlayer).toBeNull();
            expect(state.filters.onlineOnly).toBe(false);
            expect(state.filters.sortBy).toBe('rating');
            expect(state.pagination.page).toBe(1);
            expect(state.pagination.pageSize).toBe(20);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchPlayers', () => {
        it('should fetch players successfully', async () => {
            const mockPlayers = {
                players: [
                    {
                        id: 1,
                        userId: 101,
                        nickname: 'Player1',
                        avatar: 'avatar1.jpg',
                        ratingAverage: 4.5,
                        price: 5000,
                        gameId: 1,
                        gameName: 'Valorant',
                        tags: ['Pro'],
                        onlineStatus: 'online',
                    },
                    {
                        id: 2,
                        userId: 102,
                        nickname: 'Player2',
                        avatar: 'avatar2.jpg',
                        ratingAverage: 4.8,
                        price: 6000,
                        gameId: 1,
                        gameName: 'Valorant',
                        tags: ['Friendly'],
                        onlineStatus: 'offline',
                    },
                ],
                total: 2,
            };

            mockHttp.get.mockResolvedValueOnce(mockPlayers);

            await usePlayerStore.getState().fetchPlayers(true);

            const state = usePlayerStore.getState();
            expect(state.players).toHaveLength(2);
            expect(state.players[0].nickname).toBe('Player1');
            expect(state.players[0].online).toBe(true);
            expect(state.players[1].online).toBe(false);
            expect(state.pagination.total).toBe(2);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should handle fetch players error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await usePlayerStore.getState().fetchPlayers(true);

            const state = usePlayerStore.getState();
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });

        it('should append players when not refreshing', async () => {
            // Set initial players
            usePlayerStore.setState({
                players: [
                    {
                        id: 1,
                        userId: 101,
                        username: 'Player1',
                        nickname: 'Player1',
                        avatar: 'avatar1.jpg',
                        rating: 4.5,
                        price: 5000,
                        gameId: 1,
                        gameName: 'Valorant',
                        tags: ['Pro'],
                        online: true,
                        orderCount: 10,
                    },
                ],
                pagination: {
                    page: 1,
                    pageSize: 20,
                    hasMore: true,
                    total: 50,
                },
            });

            const mockNewPlayers = {
                players: [
                    {
                        id: 2,
                        userId: 102,
                        nickname: 'Player2',
                        avatar: 'avatar2.jpg',
                        ratingAverage: 4.8,
                    },
                ],
                total: 50,
            };

            mockHttp.get.mockResolvedValueOnce(mockNewPlayers);

            await usePlayerStore.getState().fetchPlayers(false);

            const state = usePlayerStore.getState();
            expect(state.players).toHaveLength(2);
        });
    });

    describe('fetchPlayerById', () => {
        it('should fetch player details successfully', async () => {
            const mockPlayer = {
                id: 1,
                userId: 101,
                nickname: 'Player1',
                avatar: 'avatar1.jpg',
                rating: 4.5,
                price: 5000,
                gameId: 1,
                gameName: 'Valorant',
                tags: ['Pro'],
                online: true,
                orderCount: 100,
            };

            mockHttp.get.mockResolvedValueOnce(mockPlayer);

            await usePlayerStore.getState().fetchPlayerById(1);

            const state = usePlayerStore.getState();
            expect(state.currentPlayer).toEqual(mockPlayer);
            expect(state.loading).toBe(false);
        });

        it('should handle fetch player error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Player not found'));

            await usePlayerStore.getState().fetchPlayerById(999);

            const state = usePlayerStore.getState();
            expect(state.error).toBe('Player not found');
            expect(state.loading).toBe(false);
        });
    });

    describe('Filters', () => {
        it('should set filters correctly', () => {
            // Use a spy to prevent actual API call
            vi.useFakeTimers();

            usePlayerStore.getState().setFilters({ gameId: 1, onlineOnly: true });

            const state = usePlayerStore.getState();
            expect(state.filters.gameId).toBe(1);
            expect(state.filters.onlineOnly).toBe(true);
            expect(state.pagination.page).toBe(1);
            expect(state.players).toEqual([]);

            vi.useRealTimers();
        });

        it('should reset filters to initial state', () => {
            vi.useFakeTimers();

            // Set some filters first
            usePlayerStore.setState({
                filters: {
                    gameId: 1,
                    minPrice: 1000,
                    maxPrice: 10000,
                    onlineOnly: true,
                    sortBy: 'price',
                },
            });

            mockHttp.get.mockResolvedValueOnce({ players: [], total: 0 });
            usePlayerStore.getState().resetFilters();

            const state = usePlayerStore.getState();
            expect(state.filters.onlineOnly).toBe(false);
            expect(state.filters.sortBy).toBe('rating');
            expect(state.filters.gameId).toBeUndefined();

            vi.useRealTimers();
        });
    });

    describe('calculateCommission', () => {
        it('should calculate commission correctly', async () => {
            const mockResult = {
                orderAmountCents: 10000,
                baseRate: 20,
                rankingDiscount: 0,
                effectiveRate: 20,
                commissionCents: 2000,
                playerEarningsCents: 8000
            };
            mockHttp.post.mockResolvedValueOnce(mockResult);

            const result = await usePlayerStore.getState().calculateCommission(10000);

            expect(result.orderAmountCents).toBe(10000);
            expect(result.baseRate).toBe(20);
            expect(result.effectiveRate).toBe(20);
            expect(result.commissionCents).toBe(2000);
            expect(result.playerEarningsCents).toBe(8000);
        });

        it('should handle zero amount', async () => {
            const mockResult = {
                orderAmountCents: 0,
                baseRate: 20,
                rankingDiscount: 0,
                effectiveRate: 20,
                commissionCents: 0,
                playerEarningsCents: 0
            };
            mockHttp.post.mockResolvedValueOnce(mockResult);

            const result = await usePlayerStore.getState().calculateCommission(0);

            expect(result.orderAmountCents).toBe(0);
            expect(result.commissionCents).toBe(0);
            expect(result.playerEarningsCents).toBe(0);
        });

        it('should handle large amounts', async () => {
            const mockResult = {
                orderAmountCents: 1000000,
                baseRate: 20,
                rankingDiscount: 0,
                effectiveRate: 20,
                commissionCents: 200000,
                playerEarningsCents: 800000
            };
            mockHttp.post.mockResolvedValueOnce(mockResult);

            const result = await usePlayerStore.getState().calculateCommission(1000000);

            expect(result.commissionCents).toBe(200000);
            expect(result.playerEarningsCents).toBe(800000);
        });

        it('should return fallback on API error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('API error'));

            const result = await usePlayerStore.getState().calculateCommission(10000);

            // Should return safe fallback with zero earnings
            expect(result.orderAmountCents).toBe(10000);
            expect(result.playerEarningsCents).toBe(0);
            expect(result.commissionCents).toBe(10000);
        });
    });

    describe('Player Application', () => {
        it('should submit application successfully', async () => {
            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce({
                status: VerificationStatus.PENDING,
                appliedAt: '2024-01-01T00:00:00Z',
                canReapply: false,
            });

            const applicationData = {
                realName: 'John Doe',
                idCard: '123456789012345678',
                games: [{ gameId: 1, rank: 'Diamond' }],
            };

            await usePlayerStore.getState().applyToBePlayer(applicationData);

            expect(mockHttp.post).toHaveBeenCalledWith('/player/apply', applicationData);
            const state = usePlayerStore.getState();
            expect(state.applicationStatus?.status).toBe(VerificationStatus.PENDING);
        });

        it('should handle application error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Application failed'));

            const applicationData = {
                realName: 'John Doe',
                idCard: '123456789012345678',
                games: [{ gameId: 1, rank: 'Diamond' }],
            };

            await expect(
                usePlayerStore.getState().applyToBePlayer(applicationData)
            ).rejects.toThrow('Application failed');

            const state = usePlayerStore.getState();
            expect(state.error).toBe('Application failed');
        });
    });

    describe('Player Profile', () => {
        it('should fetch my profile successfully', async () => {
            const mockProfile = {
                id: 1,
                userId: 101,
                displayName: 'ProPlayer',
                bio: 'Professional gamer',
                avatar: 'avatar.jpg',
                onlineStatus: OnlineStatus.ONLINE,
                rating: 4.9,
                totalOrders: 500,
            };

            mockHttp.get.mockResolvedValueOnce(mockProfile);

            await usePlayerStore.getState().fetchMyProfile();

            const state = usePlayerStore.getState();
            expect(state.myProfile?.displayName).toBe('ProPlayer');
            expect(state.loading).toBe(false);
        });

        it('should update online status', async () => {
            usePlayerStore.setState({
                myProfile: {
                    id: 1,
                    userId: 101,
                    realName: 'John',
                    idCard: '***',
                    verificationStatus: VerificationStatus.VERIFIED,
                    displayName: 'ProPlayer',
                    bio: 'Pro gamer',
                    avatar: 'avatar.jpg',
                    gallery: [],
                    games: [],
                    hourlyRateCents: 5000,
                    serviceTimeSlots: [],
                    onlineStatus: OnlineStatus.OFFLINE,
                    acceptingOrders: true,
                    totalOrders: 100,
                    completedOrders: 95,
                    rating: 4.8,
                    reviewCount: 80,
                    totalEarningsCents: 500000,
                    monthlyEarningsCents: 50000,
                    createdAt: '2024-01-01',
                    updatedAt: '2024-01-01',
                },
            });

            mockHttp.put.mockResolvedValueOnce({});

            await usePlayerStore.getState().updateOnlineStatus(OnlineStatus.ONLINE);

            const state = usePlayerStore.getState();
            expect(state.myProfile?.onlineStatus).toBe(OnlineStatus.ONLINE);
        });
    });

    describe('Earnings', () => {
        it('should fetch earnings successfully', async () => {
            const mockEarnings = {
                totalEarningsCents: 1000000,
                monthlyEarningsCents: 100000,
                weeklyEarningsCents: 25000,
                todayEarningsCents: 5000,
                wallet: {
                    balanceCents: 50000,
                    frozenCents: 10000,
                    pendingWithdrawCents: 0,
                },
                completedOrders: 200,
                averageOrderCents: 5000,
            };

            mockHttp.get.mockResolvedValueOnce(mockEarnings);

            await usePlayerStore.getState().fetchEarnings();

            const state = usePlayerStore.getState();
            expect(state.earnings?.totalEarningsCents).toBe(1000000);
            expect(state.earnings?.wallet.balanceCents).toBe(50000);
        });

        it('should fetch earnings records', async () => {
            const mockRecords = {
                items: [
                    {
                        id: 1,
                        orderId: 101,
                        orderNo: 'ORD001',
                        orderAmountCents: 10000,
                        commissionCents: 2000,
                        earningsCents: 8000,
                        status: 'settled',
                        createdAt: '2024-01-01',
                    },
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockRecords);

            await usePlayerStore.getState().fetchEarningsRecords();

            const state = usePlayerStore.getState();
            expect(state.earningsRecords).toHaveLength(1);
            expect(state.earningsRecords[0].earningsCents).toBe(8000);
        });
    });

    describe('Pagination', () => {
        it('should set page and fetch players', async () => {
            mockHttp.get.mockResolvedValueOnce({ players: [], total: 0 });

            usePlayerStore.getState().setPage(2);

            const state = usePlayerStore.getState();
            expect(state.pagination.page).toBe(2);
            expect(mockHttp.get).toHaveBeenCalled();
        });
    });
});
