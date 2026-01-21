/**
 * Monitor Store - 实时监控数据状态管理
 * 处理系统状态、在线用户、订单队列等实时数据
 */
import { create } from 'zustand';
import type {
  SystemStatus,
  OnlineUsers,
  OrderQueue,
  Alert,
} from '@/utils/websocket';

// ==================== 状态定义 ====================

interface MonitorState {
  // 系统状态
  systemStatus: SystemStatus | null;
  systemStatusLoading: boolean;
  systemStatusLastUpdate: string | null;

  // 在线用户
  onlineUsers: OnlineUsers | null;
  onlineUsersLoading: boolean;

  // 订单队列
  orderQueue: OrderQueue | null;
  orderQueueLoading: boolean;

  // 警告列表
  alerts: Alert[];
  alertsLoading: boolean;

  // WebSocket 连接状态
  wsConnected: boolean;

  // ==================== Actions ====================

  // 更新系统状态
  setSystemStatus: (status: SystemStatus) => void;

  // 更新在线用户
  setOnlineUsers: (data: OnlineUsers) => void;

  // 更新订单队列
  setOrderQueue: (data: OrderQueue) => void;

  // 添加警告
  addAlert: (alert: Alert) => void;

  // 标记警告为已读
  markAlertAsRead: (alertId: string) => void;

  // 清除警告
  clearAlert: (alertId: string) => void;

  // 清除所有警告
  clearAllAlerts: () => void;

  // 设置 WebSocket 连接状态
  setWsConnected: (connected: boolean) => void;

  // 重置所有状态
  reset: () => void;
}

// ==================== Store 实现 ====================

export const useMonitorStore = create<MonitorState>((set) => ({
  // 初始状态
  systemStatus: null,
  systemStatusLoading: false,
  systemStatusLastUpdate: null,

  onlineUsers: null,
  onlineUsersLoading: false,

  orderQueue: null,
  orderQueueLoading: false,

  alerts: [],
  alertsLoading: false,

  wsConnected: false,

  // ==================== Actions 实现 ====================

  /**
   * 更新系统状态
   */
  setSystemStatus: (status) => {
    set({
      systemStatus: status,
      systemStatusLoading: false,
      systemStatusLastUpdate: new Date().toISOString(),
    });
  },

  /**
   * 更新在线用户
   */
  setOnlineUsers: (data) => {
    set({
      onlineUsers: data,
      onlineUsersLoading: false,
    });
  },

  /**
   * 更新订单队列
   */
  setOrderQueue: (data) => {
    set({
      orderQueue: data,
      orderQueueLoading: false,
    });
  },

  /**
   * 添加警告
   */
  addAlert: (alert) => {
    set((state) => ({
      alerts: [alert, ...state.alerts].slice(0, 100), // 最多保留 100 条
      alertsLoading: false,
    }));
  },

  /**
   * 标记警告为已读
   */
  markAlertAsRead: (alertId) => {
    set((state) => ({
      alerts: state.alerts.map((alert) =>
        alert.id === alertId ? { ...alert, isRead: true } : alert
      ),
    }));
  },

  /**
   * 清除警告
   */
  clearAlert: (alertId) => {
    set((state) => ({
      alerts: state.alerts.filter((alert) => alert.id !== alertId),
    }));
  },

  /**
   * 清除所有警告
   */
  clearAllAlerts: () => {
    set({ alerts: [] });
  },

  /**
   * 设置 WebSocket 连接状态
   */
  setWsConnected: (connected) => {
    set({ wsConnected: connected });
  },

  /**
   * 重置所有状态
   */
  reset: () => {
    set({
      systemStatus: null,
      systemStatusLoading: false,
      systemStatusLastUpdate: null,
      onlineUsers: null,
      onlineUsersLoading: false,
      orderQueue: null,
      orderQueueLoading: false,
      alerts: [],
      alertsLoading: false,
      wsConnected: false,
    });
  },
}));

// ==================== 选择器 ====================

/**
 * 获取未读警告数量
 */
export const selectUnreadAlertsCount = (): number => {
  return useMonitorStore.getState().alerts.filter((a) => !a.isRead).length;
};

/**
 * 获取高级别警告
 */
export const selectHighLevelAlerts = (): Alert[] => {
  return useMonitorStore.getState().alerts.filter((a) => a.level === 'high' && !a.isRead);
};

/**
 * 获取系统健康状态
 */
export const selectSystemHealth = (): 'healthy' | 'degraded' | 'critical' => {
  const status = useMonitorStore.getState().systemStatus;
  return status?.status ?? 'healthy';
};
