import apiClient from './client';

export interface Notification {
    id: number;
    title: string;
    message: string;
    priority: 'high' | 'normal' | 'low';
    channel: string;
    referenceType: string;
    referenceId?: number;
    readAt?: string;
    isRead: boolean;
    createdAt: string;
}

export interface NotificationQueryParams {
    page?: number;
    page_size?: number;
    type?: string;
}

export interface ApiResponse<T> {
    success: boolean;
    data: T;
    message?: string;
    pagination?: {
        total: number;
        page: number;
        pageSize: number;
        unreadCount?: number;
    };
    traceId?: string;
}

export interface NotificationListResponse {
    items: Notification[];
    page: number;
    pageSize: number;
    total: number;
    unreadCount: number;
}

export const userApi = {
    // Notifications
    getNotifications: (params?: NotificationQueryParams) => apiClient.get<ApiResponse<NotificationListResponse>>('/notifications', { params }),
    getUnreadCount: () => apiClient.get<ApiResponse<{ count: number }>>('/notifications/unread-count'),
    markAsRead: (id: number) => apiClient.post<ApiResponse<void>>('/notifications/read', { ids: [id] }),
    markAllAsRead: () => apiClient.post<ApiResponse<void>>('/notifications/read-all'),
    deleteNotification: (id: number) => apiClient.delete<ApiResponse<void>>(`/notifications/${id}`),
};
