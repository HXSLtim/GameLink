/**
 * 通知API模块
 */

import { apiClient } from '@/api/client';
import type {
  Notification,
  NotificationType,
  NotificationStatus,
  GetNotificationsParams,
  GetNotificationsResponse,
  UnreadNotificationCount,
  CreateNotificationRequest,
} from '@/shared/types/notification';

/**
 * 获取通知列表
 */
export const getNotifications = async (
  params?: GetNotificationsParams
): Promise<GetNotificationsResponse> => {
  const response = await apiClient.get<{ data: GetNotificationsResponse }>('/notifications', {
    params,
  });
  return response.data;
};

/**
 * 获取通知详情
 */
export const getNotificationById = async (id: number): Promise<Notification> => {
  const response = await apiClient.get<{ data: Notification }>(`/notifications/${id}`);
  return response.data;
};

/**
 * 获取未读通知数量
 */
export const getUnreadNotificationCount = async (): Promise<UnreadNotificationCount> => {
  const response = await apiClient.get<{ data: UnreadNotificationCount }>(
    '/notifications/unread-count'
  );
  return response.data;
};

/**
 * 标记通知为已读
 */
export const markNotificationAsRead = async (id: number): Promise<void> => {
  await apiClient.post(`/notifications/${id}/read`);
};

/**
 * 批量标记通知为已读
 */
export const markNotificationsAsRead = async (ids: number[]): Promise<void> => {
  await apiClient.post('/notifications/read-batch', { ids });
};

/**
 * 标记所有通知为已读
 */
export const markAllNotificationsAsRead = async (): Promise<void> => {
  await apiClient.post('/notifications/read-all');
};

/**
 * 删除通知
 */
export const deleteNotification = async (id: number): Promise<void> => {
  await apiClient.delete(`/notifications/${id}`);
};

/**
 * 批量删除通知
 */
export const deleteNotifications = async (ids: number[]): Promise<void> => {
  await apiClient.post('/notifications/delete-batch', { ids });
};

/**
 * 创建系统通知（管理员）
 */
export const createSystemNotification = async (
  data: CreateNotificationRequest
): Promise<Notification> => {
  const response = await apiClient.post<{ data: Notification }>('/notifications/system', data);
  return response.data;
};

/**
 * 订阅通知（WebSocket）
 * 返回 WebSocket URL
 */
export const subscribeNotifications = async (): Promise<{ wsUrl: string }> => {
  const response = await apiClient.get<{ data: { wsUrl: string } }>('/notifications/subscribe');
  return response.data;
};