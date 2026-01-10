/**
 * 监控模块类型定义
 */

/**
 * WebSocket 消息类型
 */
export type WSMessageType =
  | 'system_status'
  | 'online_users'
  | 'order_queue'
  | 'alert'
  | 'ping'
  | 'pong';

/**
 * WebSocket 消息结构
 */
export interface WSMessage<T = unknown> {
  type: WSMessageType;
  timestamp: string;
  data?: T;
}

/**
 * 数据库连接状态
 */
export interface DBConnections {
  active: number;
  idle: number;
  max: number;
}

/**
 * 系统状态
 */
export interface SystemStatus {
  cpuUsage: number;
  memoryUsage: number;
  memoryTotal: number;
  memoryUsed: number;
  goroutines: number;
  dbConnections: DBConnections;
  uptime: number;
  requestsPerSec: number;
  status: 'healthy' | 'degraded' | 'critical';
}

/**
 * 在线用户统计
 */
export interface OnlineUsers {
  total: number;
  peak: number;
  byRole: {
    user?: number;
    player?: number;
    admin?: number;
  };
  updatedAt: string;
}

/**
 * 订单队列状态
 */
export interface OrderQueue {
  pending: number;
  processing: number;
  completed: number;
  processingSpeed: number;
  averageWaitTime: number;
  hasBacklog: boolean;
}

/**
 * 告警级别
 */
export type AlertLevel = 'high' | 'medium' | 'low';

/**
 * 告警类型
 */
export type AlertType = 'system' | 'business' | 'security';

/**
 * 告警信息
 */
export interface Alert {
  id: string;
  level: AlertLevel;
  type: AlertType;
  title: string;
  message: string;
  source: string;
  createdAt: string;
  isRead: boolean;
}

/**
 * KPI 周期类型
 */
export type KPIPeriodType = 'daily' | 'weekly' | 'monthly';

/**
 * KPI 目标配置
 */
export interface KPITarget {
  id: number;
  periodType: KPIPeriodType;
  metricName: string;
  targetValue: number;
  startDate: string;
  endDate: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * KPI 概览数据
 */
export interface KPIOverview {
  gmv: KPIMetric;
  orders: KPIMetric;
  newUsers: KPIMetric;
  newPlayers: KPIMetric;
  dau: KPIMetric;
  retention: KPIMetric;
  repurchaseRate: KPIMetric;
}

/**
 * KPI 指标
 */
export interface KPIMetric {
  name: string;
  currentValue: number;
  targetValue: number;
  previousValue: number;
  completionRate: number;
  changeRate: number;
  trend: 'up' | 'down' | 'flat';
}

/**
 * KPI 趋势数据点
 */
export interface KPITrendPoint {
  date: string;
  value: number;
  target?: number;
}

/**
 * 运营分析 - 活跃用户数据
 */
export interface ActiveUsersData {
  dau: number;
  wau: number;
  mau: number;
  dauMauRatio: number;
  trend: TrendPoint[];
}

/**
 * 运营分析 - 留存数据
 */
export interface RetentionData {
  day1: number;
  day7: number;
  day30: number;
  matrix: RetentionMatrixRow[];
}

/**
 * 留存矩阵行
 */
export interface RetentionMatrixRow {
  cohortDate: string;
  cohortSize: number;
  retention: number[];
}

/**
 * 运营分析 - 付费数据
 */
export interface PaymentData {
  payingRate: number;
  arpu: number;
  arppu: number;
  avgOrderValue: number;
  trend: TrendPoint[];
}

/**
 * 运营分析 - 转化漏斗
 */
export interface ConversionFunnel {
  steps: FunnelStep[];
}

/**
 * 漏斗步骤
 */
export interface FunnelStep {
  name: string;
  value: number;
  rate: number;
}

/**
 * 趋势数据点
 */
export interface TrendPoint {
  date: string;
  value: number;
  type?: string;
}

/**
 * 对比类型
 */
export type CompareType = 'mom' | 'yoy'; // 环比 | 同比

/**
 * 日期范围
 */
export interface DateRange {
  start: string;
  end: string;
}

/**
 * WebSocket 连接状态
 */
export type WSConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

/**
 * WebSocket Hook 选项
 */
export interface UseWebSocketOptions {
  url: string;
  onMessage?: (data: WSMessage) => void;
  onError?: (error: Event) => void;
  onOpen?: () => void;
  onClose?: () => void;
  reconnectInterval?: number;
  maxRetries?: number;
  autoConnect?: boolean;
}

/**
 * WebSocket Hook 返回值
 */
export interface UseWebSocketReturn {
  connected: boolean;
  connectionState: WSConnectionState;
  send: (data: unknown) => void;
  connect: () => void;
  disconnect: () => void;
  lastMessage: WSMessage | null;
}

/**
 * 运营分析 - 活跃用户数据
 */
export interface ActiveUsersData {
  dau: number;
  wau: number;
  mau: number;
  dauMauRatio: number;
  trend: TrendPoint[];
}

/**
 * 运营分析 - 留存数据
 */
export interface RetentionData {
  day1: number;
  day7: number;
  day30: number;
  matrix: RetentionMatrixRow[];
}

/**
 * 留存矩阵行
 */
export interface RetentionMatrixRow {
  cohortDate: string;
  cohortSize: number;
  retention: number[];
}

/**
 * 运营分析 - 付费数据
 */
export interface PaymentData {
  payingRate: number;
  arpu: number;
  arppu: number;
  avgOrderValue: number;
  trend: TrendPoint[];
}

/**
 * 运营分析 - 转化漏斗
 */
export interface ConversionFunnel {
  steps: FunnelStep[];
}

/**
 * 漏斗步骤
 */
export interface FunnelStep {
  name: string;
  value: number;
  rate: number;
}
