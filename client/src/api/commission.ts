/**
 * Commission API
 * Handles commission rules, commission history
 */

import { http } from '@/lib/http';
import type {
    CommissionRule,
    CommissionRecord,
    PaginatedResponse
} from '@/types/api';

export const commissionApi = {
    /**
     * Get commission rules
     */
    getRules: () =>
        http.get<CommissionRule[]>('/commission/rules'),

    /**
     * Get player's commission rule
     */
    getPlayerRule: (playerId?: number) =>
        http.get<CommissionRule>(`/commission/player${playerId ? `/${playerId}` : ''}`),

    /**
     * Get commission records
     */
    getRecords: (params: {
        page: number;
        pageSize: number;
        startDate?: string;
        endDate?: string;
    }) =>
        http.get<PaginatedResponse<CommissionRecord>>('/commission/records', { params }),

    /**
     * Get commission statistics
     */
    getStats: (params?: {
        startDate?: string;
        endDate?: string;
    }) =>
        http.get<{
            totalCommission: number;
            averageRate: number;
            orderCount: number;
        }>('/commission/stats', { params }),
};
