export type NotificationPriority = 'low' | 'normal' | 'high';

export interface NotificationEvent {
  id: number;
  title: string;
  message: string;
  priority: NotificationPriority;
  channel: string;
  referenceType?: string;
  referenceId?: number;
  readAt?: string;
  createdAt: string;
}

export interface NotificationListResponse {
  items: NotificationEvent[];
  page: number;
  pageSize: number;
  total: number;
  unreadCount: number;
}

export interface MarkNotificationReadPayload {
  ids: number[];
}
