/**
 * LFG (Looking For Group) API
 * Handles LFG posts, team finding, matchmaking
 */

import { http } from '@/lib/http';
import type {
    LFGPost,
    CreateLFGRequest,
    LFGListParams,
    PaginatedResponse
} from '@/types/api';

export const lfgApi = {
    /**
     * Get LFG post list
     */
    list: (params: LFGListParams) =>
        http.get<PaginatedResponse<LFGPost>>('/lfg/list', { params }),

    /**
     * Get LFG post detail
     */
    get: (id: number) =>
        http.get<LFGPost>(`/lfg/${id}`),

    /**
     * Create LFG post
     */
    create: (data: CreateLFGRequest) =>
        http.post<LFGPost>('/lfg/create', data),

    /**
     * Update LFG post
     */
    update: (id: number, data: Partial<CreateLFGRequest>) =>
        http.put<LFGPost>(`/lfg/${id}`, data),

    /**
     * Delete LFG post
     */
    delete: (id: number) =>
        http.delete<void>(`/lfg/${id}`),

    /**
     * Apply to join LFG
     */
    apply: (lfgId: number, message?: string) =>
        http.post<void>(`/lfg/${lfgId}/apply`, { message }),

    /**
     * Accept application (LFG owner only)
     */
    acceptApplication: (lfgId: number, applicantId: number) =>
        http.post<void>(`/lfg/${lfgId}/accept/${applicantId}`),

    /**
     * Reject application (LFG owner only)
     */
    rejectApplication: (lfgId: number, applicantId: number) =>
        http.post<void>(`/lfg/${lfgId}/reject/${applicantId}`),

    /**
     * Get LFG applications
     */
    getApplications: (lfgId: number) =>
        http.get<Array<{
            userId: number;
            username: string;
            avatar: string;
            message: string;
            appliedAt: string;
        }>>(`/lfg/${lfgId}/applications`),

    /**
     * Get user's LFG posts
     */
    getMyPosts: (params: { page: number; pageSize: number }) =>
        http.get<PaginatedResponse<LFGPost>>('/lfg/my', { params }),
};
