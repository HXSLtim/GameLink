/**
 * 应用全局状态管理
 * 使用 Zustand 实现轻量级状态管理
 */
import { create } from 'zustand';
import type { OrderStatus } from '@/shared/types/order';

/**
 * 通知类型
 */
export type NotificationType = 'info' | 'success' | 'warning' | 'error';

/**
 * 通知项
 */
export interface NotificationItem {
  /**
   * 唯一 ID
   */
  id: string;
  /**
   * 通知类型
   */
  type: NotificationType;
  /**
   * 标题
   */
  title: string;
  /**
   * 内容
   */
  message: string;
  /**
   * 持续时间（毫秒）
   */
  duration?: number;
}

/**
 * 通知项
 */
export interface NotificationItem {
  /**
   * 唯一 ID
   */
  id: string;
  /**
   * 通知类型
   */
  type: NotificationType;
  /**
   * 标题
   */
  title: string;
  /**
   * 内容
   */
  message: string;
  /**
   * 持续时间（毫秒）
   */
  duration?: number;
}

/**
 * 应用状态
 */
export interface AppState {
  // UI 状态
  isSidebarCollapsed: boolean;
  activeMenuKey: string;
  isMobile: boolean;

  // 通知系统
  notifications: NotificationItem[];

  // 加载状态
  isLoading: boolean;
  loadingMessage?: string;

  // 数据缓存
  cachedFilters: {
    orderStatus?: OrderStatus;
    orderDateRange?: [string, string];
    playerGameId?: number;
  };

  // 操作
  toggleSidebar: (force?: boolean) => void;
  setActiveMenu: (key: string) => void;
  setMobile: (isMobile: boolean) => void;

  // 通知操作
  showNotification: (
    type: NotificationType,
    title: string,
    message: string,
    duration?: number
  ) => void;
  hideNotification: (id: string) => void;
  clearNotifications: () => void;

  // 加载状态操作
  showLoading: (message?: string) => void;
  hideLoading: () => void;

  // 数据缓存操作
  updateCachedFilters: (filters: Partial<AppState['cachedFilters']>) => void;
  clearCachedFilters: () => void;
}

/**
 * 生成唯一 ID
 */
const generateId = () =>
  `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;

/**
 * 应用 Store
 */
export const useAppStore = create<AppState>((set, get) => ({
  // 初始状态
  isSidebarCollapsed: false,
  activeMenuKey: '',
  isMobile: false,
  notifications: [],
  isLoading: false,
  loadingMessage: undefined,
  cachedFilters: {},

  // 切换侧边栏
  toggleSidebar: (force?: boolean) => {
    set((state) => ({
      isSidebarCollapsed:
        typeof force === 'boolean' ? force : !state.isSidebarCollapsed,
    }));
  },

  // 设置激活菜单
  setActiveMenu: (key: string) => {
    set({ activeMenuKey: key });
  },

  // 设置移动端状态
  setMobile: (isMobile: boolean) => {
    set({ isMobile });
  },

  // 显示通知
  showNotification: (
    type: NotificationType,
    title: string,
    message: string,
    duration: number = 4000
  ) => {
    const id = generateId();
    const notification: NotificationItem = {
      id,
      type,
      title,
      message,
      duration,
    };

    set((state) => ({
      notifications: [...state.notifications, notification],
    }));

    // 自动隐藏
    if (duration > 0) {
      setTimeout(() => {
        get().hideNotification(id);
      }, duration);
    }
  },

  // 隐藏通知
  hideNotification: (id: string) => {
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
    }));
  },

  // 清除所有通知
  clearNotifications: () => {
    set({ notifications: [] });
  },

  // 显示加载状态
  showLoading: (message?: string) => {
    set({
      isLoading: true,
      loadingMessage: message,
    });
  },

  // 隐藏加载状态
  hideLoading: () => {
    set({
      isLoading: false,
      loadingMessage: undefined,
    });
  },

  // 更新缓存的过滤器
  updateCachedFilters: (filters: Partial<AppState['cachedFilters']>) => {
    set((state) => ({
      cachedFilters: { ...state.cachedFilters, ...filters },
    }));
  },

  // 清除缓存的过滤器
  clearCachedFilters: () => {
    set({ cachedFilters: {} });
  },
}));

/**
 * Store 快捷访问 hooks
 */

/**
 * 获取通知显示函数
 */
export const useNotification = () => {
  const show = useAppStore((state) => state.showNotification);
  return show;
};

/**
 * 获取加载状态控制函数
 */
export const useLoading = () => {
  const showLoading = useAppStore((state) => state.showLoading);
  const hideLoading = useAppStore((state) => state.hideLoading);
  const isLoading = useAppStore((state) => state.isLoading);
  const message = useAppStore((state) => state.loadingMessage);

  return { showLoading, hideLoading, isLoading, message };
};
