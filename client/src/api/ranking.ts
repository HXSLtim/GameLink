/**
 * Ranking API
 * Handles player rankings, leaderboards
 */

import { http } from '@/lib/http';
import type {
    RankingPlayer,
    RankingConfig
} from '@/types/api';

export const rankingApi = {
    /**
     * Get ranking leaderboard
     */
    getLeaderboard: (params: {
        period: 'daily' | 'weekly' | 'monthly';
        limit?: number;
    }) =>
        http.get<RankingPlayer[]>('/ranking/leaderboard', { params }),

    /**
     * Get player's ranking
     */
    getPlayerRanking: (playerId?: number) =>
        http.get<{
            rank: number;
            score: number;
            period: string;
        }>(`/ranking/player${playerId ? `/${playerId}` : ''}`),

    /**
     * Get ranking configuration
     */
    getConfig: () =>
        http.get<RankingConfig>('/ranking/config'),

    /**
     * Get ranking history
     */
    getHistory: (playerId: number, params: {
        startDate: string;
        endDate: string;
    }) =>
        http.get<Array<{
            date: string;
            rank: number;
            score: number;
        }>>(`/ranking/player/${playerId}/history`, { params }),
};
