/**
 * Favorite API
 * Handles favorite players, favorite management
 */

import { http } from '@/lib/http';
import type {
    Player,
    PaginatedResponse
} from '@/types/api';

export const favoriteApi = {
    /**
     * Get favorite players
     */
    list: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<Player>>('/favorite/list', { params }),

    /**
     * Add player to favorites
     */
    add: (playerId: number) =>
        http.post<void>('/favorite/add', { playerId }),

    /**
     * Remove player from favorites
     */
    remove: (playerId: number) =>
        http.delete<void>(`/favorite/${playerId}`),

    /**
     * Check if player is favorited
     */
    isFavorited: (playerId: number) =>
        http.get<{ isFavorited: boolean }>(`/favorite/check/${playerId}`),

    /**
     * Get favorite count
     */
    getCount: () =>
        http.get<{ count: number }>('/favorite/count'),
};
