/**
 * Block API
 * Handles user blocking, block list management
 */

import { http } from '@/lib/http';
import type {
    BlockedUser,
    PaginatedResponse
} from '@/types/api';

export const blockApi = {
    /**
     * Get blocked users list
     */
    list: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<BlockedUser>>('/block/list', { params }),

    /**
     * Block a user
     */
    block: (userId: number, reason?: string) =>
        http.post<void>('/block/add', { userId, reason }),

    /**
     * Unblock a user
     */
    unblock: (userId: number) =>
        http.delete<void>(`/block/${userId}`),

    /**
     * Check if user is blocked
     */
    isBlocked: (userId: number) =>
        http.get<{ isBlocked: boolean }>(`/block/check/${userId}`),

    /**
     * Get block count
     */
    getCount: () =>
        http.get<{ count: number }>('/block/count'),
};
