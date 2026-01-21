/**
 * Activity Store Tests
 * Tests for activity management and participation
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useActivityStore, ActivityType, ActivityStatus } from '../activity-store';

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

// Helper function to create mock activity
function createMockActivity(overrides = {}) {
    const now = new Date();
    const startDate = new Date(now.getTime() - 24 * 60 * 60 * 1000); // Yesterday
    const endDate = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000); // 30 days from now

    return {
        id: 1,
        name: 'Test Activity',
        type: ActivityType.FESTIVAL,
        status: ActivityStatus.ACTIVE,
        bannerUrl: '/banner.jpg',
        startAt: startDate.toISOString(),
        endAt: endDate.toISOString(),
        rules: [],
        currentParticipants: 10,
        userLimit: 3,
        maxParticipants: 100,
        rewards: [],
        ...overrides,
    };
}

describe('Activity Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        useActivityStore.setState({
            activities: [],
            currentActivity: null,
            myParticipations: [],
            loading: false,
            error: null,
        });
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useActivityStore.getState();

            expect(state.activities).toEqual([]);
            expect(state.currentActivity).toBeNull();
            expect(state.myParticipations).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchActivities', () => {
        it('should fetch activities successfully', async () => {
            const mockActivities = [
                {
                    id: 1,
                    name: 'New Year Event',
                    type: ActivityType.FESTIVAL,
                    status: ActivityStatus.ACTIVE,
                    bannerUrl: '/banner1.jpg',
                    startAt: '2024-01-01T00:00:00Z',
                    endAt: '2024-01-31T23:59:59Z',
                    rules: [],
                    currentParticipants: 50,
                    userLimit: 3,
                    rewards: [],
                },
                {
                    id: 2,
                    name: 'Weekend Bonus',
                    type: ActivityType.RECHARGE,
                    status: ActivityStatus.ACTIVE,
                    bannerUrl: '/banner2.jpg',
                    startAt: '2024-01-06T00:00:00Z',
                    endAt: '2024-01-07T23:59:59Z',
                    rules: [],
                    currentParticipants: 100,
                    userLimit: 1,
                    rewards: [],
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockActivities);

            await useActivityStore.getState().fetchActivities();

            const state = useActivityStore.getState();
            expect(state.activities).toHaveLength(2);
            expect(state.activities[0].name).toBe('New Year Event');
            expect(state.loading).toBe(false);
            expect(mockHttp.get).toHaveBeenCalledWith('/activities', {
                params: { status: 'active' }
            });
        });

        it('should fetch activities by type', async () => {
            mockHttp.get.mockResolvedValueOnce([]);

            await useActivityStore.getState().fetchActivities(ActivityType.RECHARGE);

            expect(mockHttp.get).toHaveBeenCalledWith('/activities', {
                params: { status: 'active', type: 'recharge' }
            });
        });

        it('should handle fetch activities error', async () => {
            const errorMessage = 'Network error';
            mockHttp.get.mockRejectedValueOnce(new Error(errorMessage));

            await useActivityStore.getState().fetchActivities();

            const state = useActivityStore.getState();
            expect(state.error).toBe(errorMessage);
            expect(state.loading).toBe(false);
            expect(state.activities).toEqual([]);
        });
    });

    describe('fetchActivityDetail', () => {
        it('should fetch activity detail successfully', async () => {
            const mockActivity = {
                id: 1,
                name: 'Flash Sale',
                type: ActivityType.FLASH,
                status: ActivityStatus.ACTIVE,
                bannerUrl: '/flash.jpg',
                detailUrl: '/detail/1',
                startAt: '2024-01-01T00:00:00Z',
                endAt: '2024-01-01T23:59:59Z',
                rules: [
                    { type: 'min_order', value: 5000, description: 'Minimum order ¥50' }
                ],
                currentParticipants: 200,
                userLimit: 1,
                rewards: [
                    { type: 'discount', value: '20%', description: '20% off' }
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockActivity);

            await useActivityStore.getState().fetchActivityDetail(1);

            const state = useActivityStore.getState();
            expect(state.currentActivity).toEqual(mockActivity);
            expect(state.loading).toBe(false);
            expect(mockHttp.get).toHaveBeenCalledWith('/activities/1');
        });

        it('should handle fetch activity detail error', async () => {
            const errorMessage = 'Activity not found';
            mockHttp.get.mockRejectedValueOnce(new Error(errorMessage));

            await useActivityStore.getState().fetchActivityDetail(999);

            const state = useActivityStore.getState();
            expect(state.error).toBe(errorMessage);
            expect(state.currentActivity).toBeNull();
        });
    });

    describe('fetchMyParticipation', () => {
        it('should fetch participation successfully', async () => {
            const mockParticipation = {
                activityId: 1,
                userId: 100,
                participatedAt: '2024-01-01T10:00:00Z',
                participationCount: 1,
                rewardsReceived: [
                    { type: 'coupon', value: 'SAVE10', description: '¥10 coupon' }
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockParticipation);

            const result = await useActivityStore.getState().fetchMyParticipation(1);

            expect(result).toEqual(mockParticipation);
            expect(useActivityStore.getState().myParticipations).toContain(mockParticipation);
        });

        it('should return null for 404 (not participated)', async () => {
            const error = new Error('Not found');
            (error as any).status = 404;
            mockHttp.get.mockRejectedValueOnce(error);

            const result = await useActivityStore.getState().fetchMyParticipation(1);

            expect(result).toBeNull();
            expect(useActivityStore.getState().myParticipations).not.toContain(
                expect.objectContaining({ activityId: 1 })
            );
        });

        it('should handle other errors in fetchMyParticipation', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Server error'));

            const result = await useActivityStore.getState().fetchMyParticipation(1);

            expect(result).toBeNull();
        });
    });

    describe('joinActivity', () => {
        it('should join activity successfully', async () => {
            const mockActivity = {
                id: 1,
                name: 'Test Activity',
                type: ActivityType.DISCOUNT,
                status: ActivityStatus.ACTIVE,
                bannerUrl: '/test.jpg',
                startAt: '2024-01-01T00:00:00Z',
                endAt: '2024-12-31T23:59:59Z',
                rules: [],
                currentParticipants: 10,
                userLimit: 3,
                rewards: [],
            };

            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce(mockActivity); // fetchActivityDetail
            mockHttp.get.mockResolvedValueOnce({
                activityId: 1,
                participationCount: 1,
                rewardsReceived: [],
            } as any); // fetchMyParticipation

            await useActivityStore.getState().joinActivity(1);

            const state = useActivityStore.getState();
            expect(state.currentActivity).toEqual(mockActivity);
            expect(state.loading).toBe(false);
            expect(mockHttp.post).toHaveBeenCalledWith('/activities/1/join');
        });

        it('should handle join activity error', async () => {
            const errorMessage = 'Activity full';
            mockHttp.post.mockRejectedValueOnce(new Error(errorMessage));

            await expect(useActivityStore.getState().joinActivity(1))
                .rejects.toThrow(errorMessage);

            const state = useActivityStore.getState();
            expect(state.error).toBe(errorMessage);
            expect(state.loading).toBe(false);
        });
    });

    describe('getActiveActivities', () => {
        it('should return active activities within time range', () => {
            const now = new Date();
            const tomorrow = new Date(now);
            tomorrow.setDate(tomorrow.getDate() + 1);
            const yesterday = new Date(now);
            yesterday.setDate(yesterday.getDate() - 1);

            useActivityStore.setState({
                activities: [
                    {
                        id: 1,
                        name: 'Active Event',
                        type: ActivityType.DISCOUNT,
                        status: ActivityStatus.ACTIVE,
                        bannerUrl: '/active.jpg',
                        startAt: yesterday.toISOString(),
                        endAt: tomorrow.toISOString(),
                        rules: [],
                        currentParticipants: 50,
                        userLimit: 3,
                        rewards: [],
                    },
                    {
                        id: 2,
                        name: 'Ended Event',
                        type: ActivityType.FESTIVAL,
                        status: ActivityStatus.ENDED,
                        bannerUrl: '/ended.jpg',
                        startAt: '2024-01-01T00:00:00Z',
                        endAt: '2024-01-01T23:59:59Z',
                        rules: [],
                        currentParticipants: 100,
                        userLimit: 1,
                        rewards: [],
                    },
                    {
                        id: 3,
                        name: 'Future Event',
                        type: ActivityType.NEW_USER,
                        status: ActivityStatus.ACTIVE,
                        bannerUrl: '/future.jpg',
                        startAt: tomorrow.toISOString(),
                        endAt: new Date(tomorrow.getTime() + 86400000).toISOString(),
                        rules: [],
                        currentParticipants: 0,
                        userLimit: 5,
                        rewards: [],
                    },
                ],
            });

            const active = useActivityStore.getState().getActiveActivities();

            expect(active).toHaveLength(1);
            expect(active[0].id).toBe(1);
            expect(active[0].name).toBe('Active Event');
        });

        it('should return empty array when no activities', () => {
            const active = useActivityStore.getState().getActiveActivities();
            expect(active).toEqual([]);
        });

        it('should filter out activities with wrong status', () => {
            const now = new Date();
            const tomorrow = new Date(now);
            tomorrow.setDate(tomorrow.getDate() + 1);

            useActivityStore.setState({
                activities: [
                    {
                        id: 1,
                        name: 'Draft Event',
                        type: ActivityType.DISCOUNT,
                        status: ActivityStatus.DRAFT,
                        bannerUrl: '/draft.jpg',
                        startAt: now.toISOString(),
                        endAt: tomorrow.toISOString(),
                        rules: [],
                        currentParticipants: 0,
                        userLimit: 1,
                        rewards: [],
                    },
                    {
                        id: 2,
                        name: 'Scheduled Event',
                        type: ActivityType.RECHARGE,
                        status: ActivityStatus.SCHEDULED,
                        bannerUrl: '/scheduled.jpg',
                        startAt: now.toISOString(),
                        endAt: tomorrow.toISOString(),
                        rules: [],
                        currentParticipants: 0,
                        userLimit: 1,
                        rewards: [],
                    },
                ],
            });

            const active = useActivityStore.getState().getActiveActivities();

            expect(active).toHaveLength(0);
        });
    });

    describe('canParticipate', () => {
        const now = new Date();
        const tomorrow = new Date(now);
        tomorrow.setDate(tomorrow.getDate() + 1);
        const yesterday = new Date(now);
        yesterday.setDate(yesterday.getDate() - 1);

        const createMockActivity = (overrides = {}) => ({
            id: 1,
            name: 'Test Activity',
            type: ActivityType.DISCOUNT,
            status: ActivityStatus.ACTIVE,
            bannerUrl: '/test.jpg',
            startAt: yesterday.toISOString(),
            endAt: tomorrow.toISOString(),
            rules: [],
            currentParticipants: 50,
            userLimit: 3,
            rewards: [],
            ...overrides,
        });

        it('should allow participation for valid activity', () => {
            const activity = createMockActivity();
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(true);
            expect(result.reason).toBeUndefined();
        });

        it('should reject if activity is not active', () => {
            const activity = createMockActivity({ status: ActivityStatus.ENDED });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(false);
            expect(result.reason).toBe('活动未开始或已结束');
        });

        it('should reject if activity has not started', () => {
            const future = new Date();
            future.setDate(future.getDate() + 7);
            const activity = createMockActivity({ startAt: future.toISOString() });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(false);
            expect(result.reason).toBe('活动尚未开始');
        });

        it('should reject if activity has ended', () => {
            const past = new Date();
            past.setDate(past.getDate() - 7);
            const activity = createMockActivity({ endAt: past.toISOString() });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(false);
            expect(result.reason).toBe('活动已结束');
        });

        it('should reject if participant limit reached', () => {
            const activity = createMockActivity({
                maxParticipants: 100,
                currentParticipants: 100,
            });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(false);
            expect(result.reason).toBe('参与人数已满');
        });

        it('should reject if user limit reached', () => {
            const activity = createMockActivity({ userLimit: 2 });
            const participation = {
                activityId: 1,
                userId: 100,
                participatedAt: '2024-01-01T00:00:00Z',
                participationCount: 2,
                rewardsReceived: [],
            };
            const result = useActivityStore.getState().canParticipate(activity, participation);

            expect(result.can).toBe(false);
            expect(result.reason).toBe('已达到参与次数上限');
        });

        it('should allow if user has not reached limit', () => {
            const activity = createMockActivity({ userLimit: 3 });
            const participation = {
                activityId: 1,
                userId: 100,
                participatedAt: '2024-01-01T00:00:00Z',
                participationCount: 1,
                rewardsReceived: [],
            };
            const result = useActivityStore.getState().canParticipate(activity, participation);

            expect(result.can).toBe(true);
        });

        it('should allow if no participation record', () => {
            const activity = createMockActivity({ userLimit: 1 });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(true);
        });
    });

    describe('Activity Types', () => {
        it('should have all activity types', () => {
            expect(ActivityType.DISCOUNT).toBe('discount');
            expect(ActivityType.RECHARGE).toBe('recharge');
            expect(ActivityType.NEW_USER).toBe('new_user');
            expect(ActivityType.FESTIVAL).toBe('festival');
            expect(ActivityType.FLASH).toBe('flash');
        });
    });

    describe('Activity Status', () => {
        it('should have all activity statuses', () => {
            expect(ActivityStatus.DRAFT).toBe('draft');
            expect(ActivityStatus.SCHEDULED).toBe('scheduled');
            expect(ActivityStatus.ACTIVE).toBe('active');
            expect(ActivityStatus.ENDED).toBe('ended');
            expect(ActivityStatus.CANCELED).toBe('canceled');
        });
    });

    describe('Edge Cases', () => {
        it('should handle empty activities list', () => {
            mockHttp.get.mockResolvedValueOnce([]);

            useActivityStore.getState().fetchActivities();

            const state = useActivityStore.getState();
            expect(state.activities).toEqual([]);
        });

        it('should handle activity without maxParticipants', () => {
            const activity = createMockActivity({
                maxParticipants: undefined,
                currentParticipants: 999999,
            });
            const result = useActivityStore.getState().canParticipate(activity, undefined as any);

            expect(result.can).toBe(true);
        });
    });
});
