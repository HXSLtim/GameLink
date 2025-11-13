import { apiClient } from '../../api/client';
import type {
  NotificationListResponse,
  NotificationPriority,
  MarkNotificationReadPayload,
} from '../../types/notification';

export const notificationApi = {
  list: (params?: { page?: number; pageSize?: number; unread?: boolean; priority?: NotificationPriority[] }) => {
    return apiClient.get<NotificationListResponse>('/api/v1/notifications', { params });
  },
  markRead: (payload: MarkNotificationReadPayload) => {
    return apiClient.post('/api/v1/notifications/read', payload);
  },
  unreadCount: (): Promise<{ unread: number }> => {
    return apiClient.get('/api/v1/notifications/unread-count');
  },
};
