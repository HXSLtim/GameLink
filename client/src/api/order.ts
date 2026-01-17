/**
 * Order API
 * Handles order creation, listing, status updates
 */

import { http } from '@/lib/http';
import type {
    Order,
    CreateOrderRequest,
    UpdateOrderRequest,
    OrderListParams,
    PaginatedResponse
} from '@/types/api';

export const orderApi = {
    /**
     * Create new order
     */
    create: (data: CreateOrderRequest) =>
        http.post<Order>('/order/create', data),

    /**
     * Get order list with pagination
     */
    list: (params: OrderListParams) =>
        http.get<PaginatedResponse<Order>>('/order/list', { params }),

    /**
     * Get order detail by ID
     */
    get: (id: number) =>
        http.get<Order>(`/order/${id}`),

    /**
     * Update order
     */
    update: (id: number, data: UpdateOrderRequest) =>
        http.put<Order>(`/order/${id}`, data),

    /**
     * Cancel order
     */
    cancel: (id: number, reason: string) =>
        http.post<void>(`/order/${id}/cancel`, { reason }),

    /**
     * Accept order (player)
     */
    accept: (id: number) =>
        http.post<void>(`/order/${id}/accept`),

    /**
     * Reject order (player)
     */
    reject: (id: number, reason: string) =>
        http.post<void>(`/order/${id}/reject`, { reason }),

    /**
     * Start order (player)
     */
    start: (id: number) =>
        http.post<void>(`/order/${id}/start`),

    /**
     * Complete order (player)
     */
    complete: (id: number) =>
        http.post<void>(`/order/${id}/complete`),

    /**
     * Confirm completion (user)
     */
    confirm: (id: number) =>
        http.post<void>(`/order/${id}/confirm`),

    /**
     * Get order statistics
     */
    stats: () =>
        http.get<{
            total: number;
            pending: number;
            inProgress: number;
            completed: number;
            cancelled: number;
        }>('/order/stats'),
};
