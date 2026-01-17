/**
 * Dispute API
 * Handles dispute filing, resolution, customer service
 */

import { http } from '@/lib/http';
import type {
    Dispute,
    CreateDisputeRequest,
    DisputeListParams,
    PaginatedResponse
} from '@/types/api';

export const disputeApi = {
    /**
     * Create dispute
     */
    create: (data: CreateDisputeRequest) =>
        http.post<Dispute>('/dispute/create', data),

    /**
     * Get dispute list
     */
    list: (params: DisputeListParams) =>
        http.get<PaginatedResponse<Dispute>>('/dispute/list', { params }),

    /**
     * Get dispute detail
     */
    get: (id: number) =>
        http.get<Dispute>(`/dispute/${id}`),

    /**
     * Add evidence to dispute
     */
    addEvidence: (id: number, data: {
        description: string;
        images?: string[];
    }) =>
        http.post<void>(`/dispute/${id}/evidence`, data),

    /**
     * Cancel dispute
     */
    cancel: (id: number, reason: string) =>
        http.post<void>(`/dispute/${id}/cancel`, { reason }),

    /**
     * Accept resolution
     */
    acceptResolution: (id: number) =>
        http.post<void>(`/dispute/${id}/accept`),

    /**
     * Reject resolution
     */
    rejectResolution: (id: number, reason: string) =>
        http.post<void>(`/dispute/${id}/reject`, { reason }),

    /**
     * Get dispute statistics
     */
    getStats: () =>
        http.get<{
            total: number;
            pending: number;
            resolved: number;
            rejected: number;
        }>('/dispute/stats'),

    /**
     * Upload dispute evidence image
     */
    uploadEvidence: (file: File) => {
        const formData = new FormData();
        formData.append('evidence', file);
        return http.post<{ url: string }>('/dispute/upload/evidence', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },
};
