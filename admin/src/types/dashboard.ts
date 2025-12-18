/**
 * Dashboard Module Type Definitions
 * 仪表盘模块类型定义
 */

// ============================================================================
// Time Range Types
// ============================================================================

export type TimeRange = '7d' | '30d' | '90d';

export interface TimeRangeOption {
  label: string;
  value: TimeRange;
}

// ============================================================================
// Trend Data Types
// ============================================================================

export interface TrendData {
  current: number;
  previous: number;
  changePercentage: number;
  isUp: boolean;
}

// ============================================================================
// Order Status Types
// ============================================================================

export type OrderStatus = 
  | 'pending'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'refunded';

export const OrderStatus = {
  PENDING: 'pending' as const,
  CONFIRMED: 'confirmed' as const,
  IN_PROGRESS: 'in_progress' as const,
  COMPLETED: 'completed' as const,
  CANCELLED: 'cancelled' as const,
  REFUNDED: 'refunded' as const,
};

export interface OrderStatusData {
  status: OrderStatus;
  count: number;
  percentage: number;
}

// ============================================================================
// Revenue Data Types
// ============================================================================

export interface RevenueData {
  date: string;
  revenue: number;
  netRevenue?: number;
}

export interface RevenueTrend {
  data: RevenueData[];
  total: number;
  trend: TrendData;
}

// ============================================================================
// User Growth Types
// ============================================================================

export interface UserGrowthData {
  date: string;
  newUsers: number;
  totalUsers: number;
  isAnomaly?: boolean;
}

export interface UserGrowth {
  data: UserGrowthData[];
  trend: TrendData;
}

// ============================================================================
// Order Types
// ============================================================================

export interface Order {
  id: number;
  orderId: string;
  userName: string;
  playerName: string;
  serviceItem: string;
  amount: number;
  status: OrderStatus;
  createdAt: string;
  isNew?: boolean;
}

// ============================================================================
// Player Ranking Types
// ============================================================================

export interface PlayerRanking {
  id: number;
  rank: number;
  avatar: string;
  nickname: string;
  orderCount: number;
  rating: number;
  revenue: number;
  rankChange?: number; // 正数表示上升，负数表示下降
}

// ============================================================================
// Real-time Monitoring Types
// ============================================================================

export interface RealtimeMetrics {
  timestamp: Date;
  onlineUsers: number;
  queueLength: number;
  cpuUsage: number;
  memoryUsage: number;
  diskUsage: number;
  networkIn: number;
  networkOut: number;
}

// ============================================================================
// Alert Types
// ============================================================================

export type AlertType = 'system' | 'business' | 'security';

export const AlertType = {
  SYSTEM: 'system' as const,
  BUSINESS: 'business' as const,
  SECURITY: 'security' as const,
};

export type AlertLevel = 'info' | 'warning' | 'error' | 'critical';

export const AlertLevel = {
  INFO: 'info' as const,
  WARNING: 'warning' as const,
  ERROR: 'error' as const,
  CRITICAL: 'critical' as const,
};

export interface Alert {
  id: number;
  type: AlertType;
  level: AlertLevel;
  title: string;
  message: string;
  timestamp: Date;
  isRead: boolean;
  isResolved: boolean;
}

// ============================================================================
// KPI Metrics Types
// ============================================================================

export interface KPIMetrics {
  conversionRate: number; // 转化率 (%)
  repeatPurchaseRate: number; // 复购率 (%)
  userRetentionRate: number; // 用户留存率 (%)
  averageOrderValue: number; // 平均客单价
  dau: number; // 日活用户数
  mau: number; // 月活用户数
  paymentRate: number; // 付费率 (%)
  arpu: number; // 平均每用户收入
}

export interface KPIMetric {
  name: string;
  value: number;
  unit: string;
  target?: number;
  monthOverMonth?: number; // 环比变化
  yearOverYear?: number; // 同比变化
  isBelowTarget?: boolean;
}

// ============================================================================
// Operational Data Types
// ============================================================================

export interface OperationalData {
  dau: number;
  mau: number;
  paymentRate: number;
  averageOrderValue: number;
  trend: OperationalTrendData[];
}

export interface OperationalTrendData {
  date: string;
  dau: number;
  mau: number;
  paymentRate: number;
  averageOrderValue: number;
  isAnomaly?: boolean;
  anomalyReason?: string;
}

// ============================================================================
// Dashboard Stats Types
// ============================================================================

export interface DashboardStats {
  totalUsers: number;
  totalOrders: number;
  totalRevenue: number;
  activePlayers: number;
  userGrowth: TrendData;
  revenueTrend: TrendData;
  orderStatusDistribution: OrderStatusData[];
  recentOrders: Order[];
  topPlayers: PlayerRanking[];
}

// ============================================================================
// API Request/Response Types
// ============================================================================

export interface DashboardStatsRequest {
  timeRange: TimeRange;
}

export interface TrendRequest {
  timeRange: TimeRange;
  startDate?: string;
  endDate?: string;
}

export interface AlertRequest {
  page: number;
  pageSize: number;
  type?: AlertType;
  level?: AlertLevel;
  isRead?: boolean;
  isResolved?: boolean;
}

export interface KPIRequest {
  timeRange: TimeRange;
  startDate?: string;
  endDate?: string;
}

export interface OperationalRequest {
  timeRange: TimeRange;
  compareRange?: TimeRange;
}

// ============================================================================
// Component Props Types
// ============================================================================

export interface StatCardProps {
  title: string;
  value: number | string;
  trend?: {
    value: number;
    isUp: boolean;
  };
  icon: React.ReactNode;
  loading?: boolean;
  onClick?: () => void;
}

export interface RevenueChartProps {
  data: RevenueData[];
  timeRange: TimeRange;
  loading?: boolean;
}

export interface OrderStatusPieProps {
  data: OrderStatusData[];
  loading?: boolean;
}

export interface UserGrowthChartProps {
  data: UserGrowthData[];
  timeRange: TimeRange;
  loading?: boolean;
}

export interface RecentOrdersProps {
  orders: Order[];
  loading?: boolean;
  onOrderClick?: (orderId: string) => void;
}

export interface TopPlayersProps {
  players: PlayerRanking[];
  loading?: boolean;
  onPlayerClick?: (playerId: number) => void;
}

export interface RealtimeMonitorProps {
  wsUrl: string;
  onAlert?: (alert: Alert) => void;
}

export interface AlertBannerProps {
  alerts: Alert[];
  onAlertClick?: (alertId: number) => void;
  onAlertClose?: (alertId: number) => void;
}

export interface KPIPanelProps {
  metrics: KPIMetric[];
  timeRange: TimeRange;
  loading?: boolean;
  onMetricClick?: (metricName: string) => void;
}

export interface OperationalOverviewProps {
  data: OperationalData;
  timeRange: TimeRange;
  loading?: boolean;
  onExport?: () => void;
}

export interface DashboardProps {
  timeRange: TimeRange;
  onTimeRangeChange: (range: TimeRange) => void;
}

// ============================================================================
// State Types
// ============================================================================

export interface DashboardState {
  loading: boolean;
  stats: DashboardStats | null;
  error: Error | null;
}

export interface RealtimeMonitorState {
  connected: boolean;
  metrics: RealtimeMetrics | null;
}

// ============================================================================
// Utility Types
// ============================================================================

export interface LoadingState {
  isLoading: boolean;
  error: Error | null;
}

export interface PaginationParams {
  page: number;
  pageSize: number;
}

// ApiResponse 已移至 @/types/api，这里保留导出以保持向后兼容
export type { ApiResponse } from './api';
