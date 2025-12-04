import apiClient from './client';

export interface Notification {
    id: number;
    title: string;
    content: string;
    type: 'system' | 'activity';
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
    };
}

export const userApi = {
    // Notifications
    getNotifications: (params?: NotificationQueryParams) => apiClient.get<ApiResponse<Notification[]>>('/user/notifications', { params }),
    getUnreadCount: () => apiClient.get<ApiResponse<{ count: number }>>('/user/notifications/unread-count'),
    markAsRead: (id: number) => apiClient.put<ApiResponse<void>>(`/user/notifications/${id}/read`),
    markAllAsRead: () => apiClient.put<ApiResponse<void>>('/user/notifications/read-all'),
    deleteNotification: (id: number) => apiClient.delete<ApiResponse<void>>(`/user/notifications/${id}`),
};
