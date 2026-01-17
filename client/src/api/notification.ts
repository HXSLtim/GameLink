/**
 * Notification API
 * Handles system notifications, push notifications
 */

import { http } from '@/lib/http';
import type {
    Notification,
    NotificationListParams,
    PaginatedResponse
} from '@/types/api';

export const notificationApi = {
    /**
     * Get notification list
     */
    list: (params: NotificationListParams) =>
        http.get<PaginatedResponse<Notification>>('/notification/list', { params }),

    /**
     * Get notification detail
     */
    get: (id: number) =>
        http.get<Notification>(`/notification/${id}`),

    /**
     * Mark notification as read
     */
    markAsRead: (id: number) =>
        http.post<void>(`/notification/${id}/read`),

    /**
     * Mark all notifications as read
     */
    markAllAsRead: () =>
        http.post<void>('/notification/read-all'),

    /**
     * Delete notification
     */
    delete: (id: number) =>
        http.delete<void>(`/notification/${id}`),

    /**
     * Delete all notifications
     */
    deleteAll: () =>
        http.delete<void>('/notification/all'),

    /**
     * Get unread count
     */
    getUnreadCount: () =>
        http.get<{ count: number }>('/notification/unread-count'),

    /**
     * Get notification settings
     */
    getSettings: () =>
        http.get<{
            orderUpdates: boolean;
            chatMessages: boolean;
            systemAnnouncements: boolean;
            promotions: boolean;
        }>('/notification/settings'),

    /**
     * Update notification settings
     */
    updateSettings: (settings: {
        orderUpdates?: boolean;
        chatMessages?: boolean;
        systemAnnouncements?: boolean;
        promotions?: boolean;
    }) =>
        http.put<void>('/notification/settings', settings),
};
