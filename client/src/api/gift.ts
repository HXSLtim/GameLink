/**
 * Gift API
 * Handles gift orders, gift sending
 */

import { http } from '@/lib/http';
import type {
    Gift,
    CreateGiftOrderRequest,
    GiftOrder,
    PaginatedResponse
} from '@/types/api';

export const giftApi = {
    /**
     * Get available gifts
     */
    getAvailableGifts: () =>
        http.get<Gift[]>('/gift/available'),

    /**
     * Get gift detail
     */
    getGift: (id: number) =>
        http.get<Gift>(`/gift/${id}`),

    /**
     * Send gift (create gift order)
     */
    sendGift: (data: CreateGiftOrderRequest) =>
        http.post<GiftOrder>('/gift/send', data),

    /**
     * Get sent gifts
     */
    getSentGifts: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<GiftOrder>>('/gift/sent', { params }),

    /**
     * Get received gifts
     */
    getReceivedGifts: (params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<GiftOrder>>('/gift/received', { params }),

    /**
     * Get gift statistics
     */
    getStats: () =>
        http.get<{
            totalSent: number;
            totalReceived: number;
            totalSpent: number;
            totalEarned: number;
        }>('/gift/stats'),
};
