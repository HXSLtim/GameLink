/**
 * 通知类型
 */
export type NotificationType = 'info' | 'success' | 'warning' | 'error' | 'order' | 'payment' | 'review';

/**
 * 通知状态
 */
export type NotificationStatus = 'unread' | 'read';

/**
 * 通知优先级
 */
export type NotificationPriority = 'low' | 'normal' | 'high';

/**
 * 通知实体（主类型）
 */
export interface Notification {
  id: number;
  userId: number;
  type: NotificationType;
  title: string;
  message: string;
  priority: NotificationPriority;
  status: NotificationStatus;
  readAt?: string;
  actionUrl?: string;
  actionText?: string;
  createdAt: string;
  updatedAt?: string;
}

/**
 * 通知事件
 */
export interface NotificationEvent extends Notification {
  channel: string;
  referenceType?: string;
  referenceId?: number;
}

/**
 * 获取通知列表参数
 */
export interface GetNotificationsParams {
  page?: number;
  page_size?: number;
  userId?: number;
  type?: NotificationType;
  status?: NotificationStatus;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * 获取通知列表响应
 */
export interface GetNotificationsResponse {
  items: Notification[];
  page: number;
  page_size: number;
  total: number;
  unread_count: number;
}

/**
 * 未读通知数量
 */
export interface UnreadNotificationCount {
  total: number;
  byType: Partial<Record<NotificationType, number>>;
}

/**
 * 创建通知请求
 */
export interface CreateNotificationRequest {
  userId: number;
  type: NotificationType;
  title: string;
  message: string;
  priority?: NotificationPriority;
  actionUrl?: string;
  actionText?: string;
}

/**
 * 通知列表响应
 */
export interface NotificationListResponse {
  items: NotificationEvent[];
  page: number;
  pageSize: number;
  total: number;
  unreadCount: number;
}

/**
 * 标记通知已读请求
 */
export interface MarkNotificationReadPayload {
  ids: number[];
}
