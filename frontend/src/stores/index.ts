/**
 * 状态管理模块
 * @description
 * 集中导出所有 Zustand store，提供统一的状态管理入口
 *
 * @example
 * // 在组件中使用
 * import { useAuthStore, useAppStore } from '@/stores';
 *
 * // 认证状态
 * const { token, user, setAuth, clearAuth } = useAuthStore();
 *
 * // 应用状态
 * const { showNotification, showLoading, hideLoading } = useAppStore();
 *
 * // 使用快捷 hooks
 * import { useNotification, useLoading } from '@/stores';
 *
 * const showNotification = useNotification();
 * const { showLoading, hideLoading } = useLoading();
 */

// 认证状态管理
export { useAuthStore } from './useAuthStore';
export type { AuthState } from './useAuthStore';

// 应用全局状态管理
export { useAppStore, useNotification, useLoading } from './useAppStore';
export type { AppState, NotificationItem, NotificationType } from './useAppStore';
