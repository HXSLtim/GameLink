/**
 * Service Item API
 * Handles service items, pricing, availability
 */

import { http } from '@/lib/http';
import type {
    ServiceItem,
    PaginatedResponse
} from '@/types/api';

// 重新导出 ServiceItem 类型，方便页面使用
export type { ServiceItem } from '@/types/api';

export const itemApi = {
    /**
     * Get service item list
     */
    list: (params?: {
        page?: number;
        pageSize?: number;
        gameId?: number;
        playerId?: number;
    }) =>
        http.get<PaginatedResponse<ServiceItem>>('/item/list', { params }),

    /**
     * Get service item detail
     */
    get: (id: number) =>
        http.get<ServiceItem>(`/item/${id}`),

    /**
     * Get service items by game
     */
    getByGame: (gameId: number) =>
        http.get<ServiceItem[]>(`/item/game/${gameId}`),

    /**
     * Get service items by player
     */
    getByPlayer: (playerId: number) =>
        http.get<ServiceItem[]>(`/item/player/${playerId}`),

    /**
     * Get popular service items
     */
    getPopular: (limit: number = 10) =>
        http.get<ServiceItem[]>('/item/popular', { params: { limit } }),

    /**
     * Get service items with filters
     * @param params Filter parameters
     */
    getServiceItems: (params?: {
        subCategory?: string; // 'service' | 'gift'
        gameId?: number;
        type?: string;
        category?: string;
        keyword?: string;
    }) =>
        http.get<ServiceItem[]>('/public/service-items', { params }),

    /**
     * Get gift items
     */
    getGifts: (params?: {
        category?: string;
    }) =>
        http.get<ServiceItem[]>('/public/service-items', { 
            params: { ...params, subCategory: 'gift' } 
        }),
};
