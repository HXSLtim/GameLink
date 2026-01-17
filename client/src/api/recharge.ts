/**
 * Recharge API
 * Handles balance recharge, recharge packages
 */

import { http } from '@/lib/http';
import type {
    RechargePackage,
    CreateRechargeRequest,
    RechargeRecord,
    PaginatedResponse
} from '@/types/api';

export const rechargeApi = {
    /**
     * Get recharge packages
     */
    getPackages: () =>
        http.get<RechargePackage[]>('/recharge/packages'),

    /**
     * Create recharge order
     */
    create: (data: CreateRechargeRequest) =>
        http.post<{
            rechargeId: number;
            paymentId: number;
            amount: number;
        }>('/recharge/create', data),

    /**
     * Get recharge records
     */
    getRecords: (params: {
        page: number;
        pageSize: number;
        status?: string;
    }) =>
        http.get<PaginatedResponse<RechargeRecord>>('/recharge/records', { params }),

    /**
     * Get recharge detail
     */
    get: (id: number) =>
        http.get<RechargeRecord>(`/recharge/${id}`),

    /**
     * Check recharge status
     */
    checkStatus: (id: number) =>
        http.get<{
            status: string;
            completedAt?: string;
        }>(`/recharge/${id}/status`),

    /**
     * Get recharge statistics
     */
    getStats: () =>
        http.get<{
            totalRecharged: number;
            rechargeCount: number;
            averageAmount: number;
        }>('/recharge/stats'),
};
