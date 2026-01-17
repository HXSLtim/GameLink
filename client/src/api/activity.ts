/**
 * Activity API
 * Handles marketing activities, events, promotions
 */

import { http } from '@/lib/http';
import type {
    Activity,
    ActivityListParams,
    PaginatedResponse
} from '@/types/api';

export const activityApi = {
    /**
     * Get activity list
     */
    list: (params: ActivityListParams) =>
        http.get<PaginatedResponse<Activity>>('/activity/list', { params }),

    /**
     * Get activity detail
     */
    get: (id: number) =>
        http.get<Activity>(`/activity/${id}`),

    /**
     * Get ongoing activities
     */
    getOngoing: (limit: number = 10) =>
        http.get<Activity[]>('/activity/ongoing', { params: { limit } }),

    /**
     * Participate in activity
     */
    participate: (activityId: number, data?: Record<string, any>) =>
        http.post<{
            success: boolean;
            reward?: any;
        }>(`/activity/${activityId}/participate`, data),

    /**
     * Get user's activity participation records
     */
    getParticipationRecords: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<any>>('/activity/my-participations', { params }),

    /**
     * Check if user can participate
     */
    checkEligibility: (activityId: number) =>
        http.get<{
            eligible: boolean;
            reason?: string;
        }>(`/activity/${activityId}/check-eligibility`),
};
