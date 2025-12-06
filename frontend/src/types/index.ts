/**
 * 类型定义统一导出
 */

// API 通用类型
export * from './api';

// 业务类型
export * from './content';
export * from './review';

// Dashboard 类型 (排除与 monitor 冲突的类型)
export {
  type DashboardStats,
  type OrderStats,
  type UserStats,
  type RevenueStats,
  type TrendData,
  type PaginationParams,
  // Alert 相关类型使用 Dashboard 前缀避免冲突
  type Alert as DashboardAlert,
  type AlertLevel as DashboardAlertLevel,
  type AlertType as DashboardAlertType,
} from './dashboard';

// Monitor 类型 (使用 Monitor 前缀避免冲突)
export {
  type WSMessageType,
  type WSMessage,
  type DBConnections,
  type SystemStatus,
  type OnlineUsers,
  type OrderQueue,
  type Alert as MonitorAlert,
  type AlertLevel as MonitorAlertLevel,
} from './monitor';
