/**
 * Game API
 * Handles game list, game ranks, game categories
 */

import { http } from '@/lib/http';
import type {
    Game,
    GameRank,
    PaginatedResponse
} from '@/types/api';

export const gameApi = {
    /**
     * Get game list
     */
    list: (params?: {
        page?: number;
        pageSize?: number;
        category?: string;
        search?: string;
    }) =>
        http.get<PaginatedResponse<Game>>('/game/list', { params }),

    /**
     * Get game detail
     */
    get: (id: number) =>
        http.get<Game>(`/game/${id}`),

    /**
     * Get popular games
     */
    getPopular: (limit: number = 10) =>
        http.get<Game[]>('/game/popular', { params: { limit } }),

    /**
     * Get game categories
     */
    getCategories: () =>
        http.get<Array<{ id: string; name: string; icon: string }>>('/game/categories'),

    /**
     * Get ranks for a game
     */
    getRanks: (gameId: number) =>
        http.get<GameRank[]>(`/game/${gameId}/ranks`),

    /**
     * Search games
     */
    search: (keyword: string, limit: number = 20) =>
        http.get<Game[]>('/game/search', { params: { keyword, limit } }),
};
