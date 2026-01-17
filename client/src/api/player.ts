/**
 * Player API
 * Handles player profiles, search, certification
 */

import { http } from '@/lib/http';
import type {
    Player,
    PlayerProfile,
    UpdatePlayerProfileRequest,
    PlayerSearchParams,
    PaginatedResponse
} from '@/types/api';

export const playerApi = {
    /**
     * Get player list with filters
     */
    list: (params: PlayerSearchParams) =>
        http.get<PaginatedResponse<Player>>('/player/list', { params }),

    /**
     * Get player detail by ID
     */
    get: (id: number) =>
        http.get<Player>(`/player/${id}`),

    /**
     * Get player profile (for authenticated player)
     */
    getProfile: () =>
        http.get<PlayerProfile>('/player/profile'),

    /**
     * Update player profile
     */
    updateProfile: (data: UpdatePlayerProfileRequest) =>
        http.put<PlayerProfile>('/player/profile', data),

    /**
     * Apply to become a player
     */
    apply: (data: {
        realName: string;
        idCard: string;
        phone: string;
        gameIds: number[];
        introduction: string;
    }) =>
        http.post<void>('/player/apply', data),

    /**
     * Get player application status
     */
    getApplicationStatus: () =>
        http.get<{
            status: 'pending' | 'approved' | 'rejected';
            reason?: string;
            appliedAt: string;
        }>('/player/application/status'),

    /**
     * Get player statistics
     */
    getStats: (playerId?: number) =>
        http.get<{
            totalOrders: number;
            completedOrders: number;
            rating: number;
            reviewCount: number;
            totalEarnings: number;
        }>(`/player/${playerId ? playerId + '/' : ''}stats`),

    /**
     * Get player reviews
     */
    getReviews: (playerId: number, params: { page: number; pageSize: number }) =>
        http.get<PaginatedResponse<any>>(`/player/${playerId}/reviews`, { params }),

    /**
     * Toggle player online status
     */
    toggleOnline: (online: boolean) =>
        http.post<void>('/player/online', { online }),

    /**
     * Get recommended players
     */
    getRecommended: (limit: number = 10) =>
        http.get<Player[]>('/player/recommended', { params: { limit } }),
};
