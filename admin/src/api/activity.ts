import apiClient from './client';
import type { ApiResponse, Pagination } from '@/types/api';

// Re-export for backward compatibility
export type { ApiResponse };

// ============================================================================
// Type Definitions
// ============================================================================

/**
 * Activity Status (活动状态)
 */
export type ActivityStatus = 'draft' | 'preheat' | 'active' | 'paused' | 'ended' | 'canceled';

/**
 * Activity Type (活动类型)
 */
export type ActivityType = 'coupon' | 'discount' | 'gift';

/**
 * Activity Interface
 * Matches backend model.Activity
 */
export interface Activity {
    id: number;
    name: string;
    description?: string;
    type: ActivityType;
    status: ActivityStatus;
    coverUrl?: string;
    bannerUrl?: string;

    // Time control
    preheatAt?: string;
    startAt: string;
    endAt: string;

    // Participation limits
    totalLimit: number;
    dailyLimit: number;
    perUserLimit: number;

    // Statistics
    totalParticipants: number;
    todayParticipants: number;
    totalClaimed: number;

    // Configuration
    allowVipStack: boolean;
    rules?: string;
    sortOrder: number;
    isVisible: boolean;

    // Relations
    rewards?: ActivityReward[];

    // Timestamps
    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Reward Interface
 * Matches backend model.ActivityReward
 */
export interface ActivityReward {
    id: number;
    activityId: number;
    couponTemplateId: number;
    couponCount: number;
    probability: number;
    totalStock: number;
    remainingStock: number;
    sortOrder: number;

    // Relations
    activity?: Activity;
    couponTemplate?: {
        id: number;
        name: string;
        type: string;
    };

    // Timestamps
    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Participation Interface
 * Matches backend model.ActivityParticipation
 */
export interface ActivityParticipation {
    id: number;
    activityId: number;
    userId: number;
    rewardId: number;
    couponIds?: string;
    claimedAt: string;
    clientIp?: string;

    // Relations
    activity?: Activity;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    reward?: ActivityReward;

    // Timestamps
    createdAt: string;
    updatedAt: string;
}

/**
 * Activity Statistics
 */
export interface ActivityStats {
    activityId: number;
    totalParticipants: number;
    todayParticipants: number;
    totalClaimed: number;
    remainingStock?: number;
}

/**
 * All Activities Statistics Overview
 */
export interface AllActivityStats {
    totalActivities: number;
    activeActivities: number;
    draftActivities: number;
    endedActivities: number;
    totalParticipants: number;
    totalClaimed: number;
}

/**
 * Batch Operation Result
 */
export interface BatchOperationResult {
    successCount: number;
    failedCount: number;
    failedIds: number[];
    errors: string[];
}

// ============================================================================
// Request DTOs
// ============================================================================

/**
 * Create Activity Request
 */
export interface CreateActivityDto {
    name: string;
    description?: string;
    type: ActivityType;
    coverUrl?: string;
    bannerUrl?: string;
    preheatAt?: string;
    startAt: string;
    endAt: string;
    totalLimit?: number;
    dailyLimit?: number;
    perUserLimit?: number;
    allowVipStack?: boolean;
    rules?: string;
    sortOrder?: number;
    isVisible?: boolean;
    status?: ActivityStatus;
}

/**
 * Update Activity Request
 */
export type UpdateActivityDto = Partial<CreateActivityDto>;

/**
 * Activity Query Parameters
 */
export interface ActivityQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    type?: ActivityType;
    status?: ActivityStatus;
    isVisible?: boolean;
    startTime?: string;
    endTime?: string;
}

/**
 * Update Activity Status Request
 */
export interface UpdateActivityStatusDto {
    status: ActivityStatus;
}

/**
 * Create Reward Request
 */
export interface CreateRewardDto {
    activityId: number;
    couponTemplateId: number;
    couponCount: number;
    probability: number;
    totalStock?: number;
    sortOrder?: number;
}

/**
 * Update Reward Request
 */
export type UpdateRewardDto = Partial<CreateRewardDto>;

/**
 * Participation Query Parameters
 */
export interface ParticipationQueryParams {
    page?: number;
    page_size?: number;
    activityId?: number;
    userId?: number;
    startTime?: string;
    endTime?: string;
}

/**
 * Batch Delete Activities Request
 */
export interface BatchDeleteActivitiesDto {
    ids: number[];
}

/**
 * Batch Update Activities Status Request
 */
export interface BatchUpdateActivitiesStatusDto {
    ids: number[];
    status: ActivityStatus;
}

/**
 * Batch Publish Activities Request
 */
export interface BatchPublishActivitiesDto {
    ids: number[];
    isVisible: boolean;
}

// ============================================================================
// API Client
// ============================================================================

/**
 * Activity API
 * Provides methods for managing activities, rewards, and participations
 */
export const activityApi = {
    // ========================================================================
    // Activity Management
    // ========================================================================

    /**
     * Get activities list
     * GET /admin/activities
     */
    getActivities: (params?: ActivityQueryParams) =>
        apiClient.get<ApiResponse<Activity[]>>('/admin/activities', { params }),

    /**
     * Get activity detail
     * GET /admin/activities/:id
     */
    getActivityDetail: (id: number) =>
        apiClient.get<ApiResponse<Activity>>(`/admin/activities/${id}`),

    /**
     * Create new activity
     * POST /admin/activities
     */
    createActivity: (data: CreateActivityDto) =>
        apiClient.post<ApiResponse<Activity>>('/admin/activities', data),

    /**
     * Update activity
     * PUT /admin/activities/:id
     */
    updateActivity: (id: number, data: UpdateActivityDto) =>
        apiClient.put<ApiResponse<Activity>>(`/admin/activities/${id}`, data),

    /**
     * Delete activity
     * DELETE /admin/activities/:id
     */
    deleteActivity: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/activities/${id}`),

    /**
     * Update activity status
     * PUT /admin/activities/:id/status
     */
    updateActivityStatus: (id: number, data: UpdateActivityStatusDto) =>
        apiClient.put<ApiResponse<{ message: string }>>(`/admin/activities/${id}/status`, data),

    /**
     * Publish activity (convenience method)
     * Updates status to 'active' and isVisible to true
     */
    publishActivity: (id: number) =>
        activityApi.updateActivity(id, { status: 'active', isVisible: true }),

    /**
     * Unpublish activity (convenience method)
     * Sets isVisible to false
     */
    unpublishActivity: (id: number) =>
        activityApi.updateActivity(id, { isVisible: false }),

    // ========================================================================
    // Batch Operations
    // ========================================================================

    /**
     * Batch delete activities
     * DELETE /admin/activities/batch
     */
    batchDeleteActivities: (data: BatchDeleteActivitiesDto) =>
        apiClient.delete<ApiResponse<BatchOperationResult>>('/admin/activities/batch', { data }),

    /**
     * Batch update activities status
     * PUT /admin/activities/batch/status
     */
    batchUpdateActivitiesStatus: (data: BatchUpdateActivitiesStatusDto) =>
        apiClient.put<ApiResponse<BatchOperationResult>>('/admin/activities/batch/status', data),

    /**
     * Batch publish/unpublish activities
     * PUT /admin/activities/batch/publish
     */
    batchPublishActivities: (data: BatchPublishActivitiesDto) =>
        apiClient.put<ApiResponse<BatchOperationResult>>('/admin/activities/batch/publish', data),

    // ========================================================================
    // Reward Management
    // ========================================================================

    /**
     * Get rewards by activity ID
     * GET /admin/activities/:id/rewards
     */
    getActivityRewards: (activityId: number) =>
        apiClient.get<ApiResponse<ActivityReward[]>>(`/admin/activities/${activityId}/rewards`),

    /**
     * Get reward detail
     * GET /admin/activities/:id/rewards/:rewardId
     */
    getRewardDetail: (activityId: number, rewardId: number) =>
        apiClient.get<ApiResponse<ActivityReward>>(`/admin/activities/${activityId}/rewards/${rewardId}`),

    /**
     * Create new reward
     * POST /admin/activities/rewards
     */
    createReward: (data: CreateRewardDto) =>
        apiClient.post<ApiResponse<ActivityReward>>('/admin/activities/rewards', data),

    /**
     * Update reward
     * PUT /admin/activities/rewards/:rewardId
     */
    updateReward: (rewardId: number, data: UpdateRewardDto) =>
        apiClient.put<ApiResponse<ActivityReward>>(`/admin/activities/rewards/${rewardId}`, data),

    /**
     * Delete reward
     * DELETE /admin/activities/rewards/:rewardId
     */
    deleteReward: (rewardId: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/activities/rewards/${rewardId}`),

    // ========================================================================
    // Participation Records
    // ========================================================================

    /**
     * Get participation records list
     * GET /admin/activities/participations
     */
    getParticipations: (params?: ParticipationQueryParams) =>
        apiClient.get<ApiResponse<ActivityParticipation[]>>('/admin/activities/participations', { params }),

    // ========================================================================
    // Statistics
    // ========================================================================

    /**
     * Get activity statistics
     * GET /admin/activities/:id/stats
     */
    getActivityStats: (id: number) =>
        apiClient.get<ApiResponse<ActivityStats>>(`/admin/activities/${id}/stats`),

    /**
     * Get all activities statistics overview
     * GET /admin/activities/stats
     */
    getAllActivityStats: () =>
        apiClient.get<ApiResponse<AllActivityStats>>('/admin/activities/stats'),
};

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Format activity type display name
 */
export const getActivityTypeLabel = (type: ActivityType): string => {
    const labels: Record<ActivityType, string> = {
        coupon: '优惠券发放',
        discount: '限时折扣',
        gift: '赠品活动',
    };
    return labels[type] || type;
};

/**
 * Format activity status display name
 */
export const getActivityStatusLabel = (status: ActivityStatus): string => {
    const labels: Record<ActivityStatus, string> = {
        draft: '草稿',
        preheat: '预热中',
        active: '进行中',
        paused: '已暂停',
        ended: '已结束',
        canceled: '已取消',
    };
    return labels[status] || status;
};

/**
 * Format activity status color (for Ant Design tags)
 */
export const getActivityStatusColor = (status: ActivityStatus): string => {
    const colors: Record<ActivityStatus, string> = {
        draft: 'default',
        preheat: 'blue',
        active: 'green',
        paused: 'orange',
        ended: 'default',
        canceled: 'red',
    };
    return colors[status] || 'default';
};

/**
 * Check if activity is currently active
 */
export const isActivityActive = (activity: Activity): boolean => {
    const now = new Date();
    const startAt = new Date(activity.startAt);
    const endAt = new Date(activity.endAt);
    return activity.status === 'active' && now >= startAt && now <= endAt;
};

/**
 * Check if activity is in preheat period
 */
export const isActivityInPreheat = (activity: Activity): boolean => {
    if (!activity.preheatAt) return false;
    const now = new Date();
    const preheatAt = new Date(activity.preheatAt);
    const startAt = new Date(activity.startAt);
    return now >= preheatAt && now < startAt;
};

/**
 * Check if activity has ended
 */
export const isActivityEnded = (activity: Activity): boolean => {
    const now = new Date();
    const endAt = new Date(activity.endAt);
    return now > endAt;
};

/**
 * Calculate stock percentage
 */
export const calculateStockPercentage = (reward: ActivityReward): number => {
    if (reward.totalStock === 0) return 100; // Unlimited stock
    if (reward.totalStock <= 0) return 0;
    return Math.round((reward.remainingStock / reward.totalStock) * 100);
};

/**
 * Format participant count
 */
export const formatParticipantCount = (count: number): string => {
    if (count >= 10000) {
        return `${(count / 10000).toFixed(1)}万`;
    }
    if (count >= 1000) {
        return `${(count / 1000).toFixed(1)}k`;
    }
    return count.toString();
};

/**
 * Parse JSON array string to array
 */
export const parseCouponIds = (jsonString: string): number[] => {
    try {
        const parsed = JSON.parse(jsonString || '[]');
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return [];
    }
};

/**
 * Get activity time range display text
 */
export const getActivityTimeRange = (activity: Activity): string => {
    const start = new Date(activity.startAt);
    const end = new Date(activity.endAt);
    const format = (date: Date) => {
        return `${date.getMonth() + 1}/${date.getDate()}`;
    };
    return `${format(start)} - ${format(end)}`;
};

/**
 * Check if activity can be edited
 */
export const canEditActivity = (activity: Activity): boolean => {
    return activity.status === 'draft' || activity.status === 'preheat';
};

/**
 * Check if activity can be deleted
 */
export const canDeleteActivity = (activity: Activity): boolean => {
    return activity.status === 'draft' || activity.status === 'canceled' || activity.status === 'ended';
};

// ============================================================================
// Default export
// ============================================================================

export default activityApi;
