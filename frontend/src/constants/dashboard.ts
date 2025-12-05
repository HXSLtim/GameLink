/**
 * Dashboard Module Constants
 * 仪表盘模块常量定义
 */

import type { TimeRangeOption } from '@/types/dashboard';
import { OrderStatus, AlertType, AlertLevel } from '@/types/dashboard';

// ============================================================================
// Time Range Constants
// ============================================================================

export const TIME_RANGE_OPTIONS: TimeRangeOption[] = [
  { label: '最近7天', value: '7d' },
  { label: '最近30天', value: '30d' },
  { label: '最近90天', value: '90d' },
];

export const DEFAULT_TIME_RANGE = '7d';

// ============================================================================
// Order Status Constants
// ============================================================================

export const ORDER_STATUS_LABELS: Record<OrderStatus, string> = {
  [OrderStatus.PENDING]: '待确认',
  [OrderStatus.CONFIRMED]: '已确认',
  [OrderStatus.IN_PROGRESS]: '进行中',
  [OrderStatus.COMPLETED]: '已完成',
  [OrderStatus.CANCELLED]: '已取消',
  [OrderStatus.REFUNDED]: '已退款',
};

export const ORDER_STATUS_COLORS: Record<OrderStatus, string> = {
  [OrderStatus.PENDING]: '#faad14',
  [OrderStatus.CONFIRMED]: '#1890ff',
  [OrderStatus.IN_PROGRESS]: '#13c2c2',
  [OrderStatus.COMPLETED]: '#52c41a',
  [OrderStatus.CANCELLED]: '#8c8c8c',
  [OrderStatus.REFUNDED]: '#f5222d',
};

// ============================================================================
// Alert Constants
// ============================================================================

export const ALERT_TYPE_LABELS: Record<AlertType, string> = {
  [AlertType.SYSTEM]: '系统异常',
  [AlertType.BUSINESS]: '业务异常',
  [AlertType.SECURITY]: '安全告警',
};

export const ALERT_LEVEL_LABELS: Record<AlertLevel, string> = {
  [AlertLevel.INFO]: '信息',
  [AlertLevel.WARNING]: '警告',
  [AlertLevel.ERROR]: '错误',
  [AlertLevel.CRITICAL]: '严重',
};

export const ALERT_LEVEL_COLORS: Record<AlertLevel, string> = {
  [AlertLevel.INFO]: '#1890ff',
  [AlertLevel.WARNING]: '#faad14',
  [AlertLevel.ERROR]: '#ff4d4f',
  [AlertLevel.CRITICAL]: '#cf1322',
};

// ============================================================================
// Chart Colors
// ============================================================================

export const CHART_COLORS = {
  primary: '#1890ff',
  secondary: '#52c41a',
  warning: '#faad14',
  danger: '#f5222d',
  info: '#13c2c2',
  success: '#52c41a',
};

export const REVENUE_CHART_COLORS = {
  revenue: '#1890ff',
  netRevenue: '#52c41a',
};

export const USER_GROWTH_CHART_COLORS = {
  newUsers: '#1890ff',
  totalUsers: '#52c41a',
  anomaly: '#f5222d',
};

// ============================================================================
// Monitoring Thresholds
// ============================================================================

export const MONITORING_THRESHOLDS = {
  cpu: {
    warning: 70,
    critical: 90,
  },
  memory: {
    warning: 75,
    critical: 90,
  },
  disk: {
    warning: 80,
    critical: 95,
  },
  queueLength: {
    warning: 100,
    critical: 500,
  },
};

// ============================================================================
// Refresh Intervals (milliseconds)
// ============================================================================

export const REFRESH_INTERVALS = {
  realtime: 1000, // 1秒
  dashboard: 30000, // 30秒
  kpi: 60000, // 1分钟
};

// ============================================================================
// WebSocket Configuration
// ============================================================================

export const WEBSOCKET_CONFIG = {
  reconnectionDelay: 5000, // 5秒
  reconnectionAttempts: 3,
  timeout: 10000, // 10秒
};

// ============================================================================
// Pagination Constants
// ============================================================================

export const DEFAULT_PAGE_SIZE = 10;
export const MAX_RECENT_ORDERS = 10;
export const MAX_TOP_PLAYERS = 10;

// ============================================================================
// Number Format Constants
// ============================================================================

export const NUMBER_FORMAT = {
  locale: 'zh-CN',
  currency: 'CNY',
  decimalPlaces: 2,
};

// ============================================================================
// Animation Duration (milliseconds)
// ============================================================================

export const ANIMATION_DURATION = {
  fast: 200,
  normal: 300,
  slow: 500,
  number: 1000, // 数字滚动动画
};

// ============================================================================
// KPI Target Values
// ============================================================================

export const KPI_TARGETS = {
  conversionRate: 5, // 5%
  repeatPurchaseRate: 30, // 30%
  userRetentionRate: 40, // 40%
  paymentRate: 10, // 10%
};

// ============================================================================
// Export Format Constants
// ============================================================================

export const EXPORT_FORMATS = {
  excel: 'xlsx',
  csv: 'csv',
  pdf: 'pdf',
};

// ============================================================================
// Date Format Constants
// ============================================================================

export const DATE_FORMATS = {
  date: 'YYYY-MM-DD',
  datetime: 'YYYY-MM-DD HH:mm:ss',
  time: 'HH:mm:ss',
  monthDay: 'MM-DD',
  yearMonth: 'YYYY-MM',
};
