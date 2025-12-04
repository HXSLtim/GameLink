/**
 * 监控模块 API
 */
import apiClient from './client';
import type { ApiResponse } from './admin';
import type {
  SystemStatus,
  OnlineUsers,
  OrderQueue,
  Alert,
  KPITarget,
  KPIOverview,
  KPITrendPoint,
  ActiveUsersData,
  RetentionData,
  PaymentData,
  ConversionFunnel,
} from '@/types/monitor';

/**
 * 告警列表查询参数
 */
export interface AlertQueryParams {
  page?: number;
  page_size?: number;
  level?: string;
  type?: string;
  is_read?: boolean;
  date_from?: string;
  date_to?: string;
}

/**
 * KPI 查询参数
 */
export interface KPIQueryParams {
  period?: 'today' | 'week' | 'month' | 'custom';
  start_date?: string;
  end_date?: string;
  compare?: 'mom' | 'yoy';
}

/**
 * 运营分析查询参数
 */
export interface AnalyticsQueryParams {
  start_date: string;
  end_date: string;
  granularity?: 'day' | 'week' | 'month';
}

/**
 * 监控模块 API
 */
export const monitorApi = {
  // ==================== 实时监控 ====================

  /**
   * 获取系统状态快照
   */
  getSystemStatus: () =>
    apiClient.get<ApiResponse<SystemStatus>>('/admin/monitor/system-status'),

  /**
   * 获取在线用户统计
   */
  getOnlineUsers: () =>
    apiClient.get<ApiResponse<OnlineUsers>>('/admin/monitor/online-users'),

  /**
   * 获取订单队列状态
   */
  getOrderQueue: () =>
    apiClient.get<ApiResponse<OrderQueue>>('/admin/monitor/order-queue'),

  /**
   * 获取告警列表
   */
  getAlerts: (params?: AlertQueryParams) =>
    apiClient.get<ApiResponse<Alert[]>>('/admin/monitor/alerts', { params }),

  /**
   * 标记告警已读
   */
  markAlertRead: (id: string) =>
    apiClient.put<ApiResponse<void>>(`/admin/monitor/alerts/${id}/read`),

  /**
   * 批量标记告警已读
   */
  markAlertsRead: (ids: string[]) =>
    apiClient.put<ApiResponse<void>>('/admin/monitor/alerts/batch-read', { ids }),

  // ==================== 运营分析 ====================

  /**
   * 获取活跃用户数据
   */
  getActiveUsers: (params: AnalyticsQueryParams) =>
    apiClient.get<ApiResponse<ActiveUsersData>>('/admin/analytics/active-users', { params }),

  /**
   * 获取用户留存数据
   */
  getRetention: (params: AnalyticsQueryParams) =>
    apiClient.get<ApiResponse<RetentionData>>('/admin/analytics/retention', { params }),

  /**
   * 获取付费分析数据
   */
  getPaymentAnalytics: (params: AnalyticsQueryParams) =>
    apiClient.get<ApiResponse<PaymentData>>('/admin/analytics/payment', { params }),

  /**
   * 获取转化漏斗数据
   */
  getConversionFunnel: (params: AnalyticsQueryParams) =>
    apiClient.get<ApiResponse<ConversionFunnel>>('/admin/analytics/conversion', { params }),

  // ==================== KPI ====================

  /**
   * 获取 KPI 概览
   */
  getKPIOverview: (params?: KPIQueryParams) =>
    apiClient.get<ApiResponse<KPIOverview>>('/admin/kpi/overview', { params }),

  /**
   * 获取 KPI 趋势
   */
  getKPITrend: (metric: string, params?: KPIQueryParams) =>
    apiClient.get<ApiResponse<KPITrendPoint[]>>('/admin/kpi/trend', {
      params: { metric, ...params },
    }),

  /**
   * 获取 KPI 目标配置
   */
  getKPITargets: (params?: { period_type?: string; metric_name?: string }) =>
    apiClient.get<ApiResponse<KPITarget[]>>('/admin/kpi/targets', { params }),

  /**
   * 创建 KPI 目标
   */
  createKPITarget: (data: Omit<KPITarget, 'id' | 'createdBy' | 'createdAt' | 'updatedAt'>) =>
    apiClient.post<ApiResponse<KPITarget>>('/admin/kpi/targets', data),

  /**
   * 更新 KPI 目标
   */
  updateKPITarget: (id: number, data: Partial<KPITarget>) =>
    apiClient.put<ApiResponse<KPITarget>>(`/admin/kpi/targets/${id}`, data),

  /**
   * 删除 KPI 目标
   */
  deleteKPITarget: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/kpi/targets/${id}`),
};

/**
 * 获取 WebSocket 监控连接 URL
 */
export function getMonitorWSUrl(): string {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1';
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

  // 如果 baseUrl 是完整 URL，提取 host
  if (baseUrl.startsWith('http')) {
    const url = new URL(baseUrl);
    return `${wsProtocol}//${url.host}${url.pathname}/admin/ws/monitor`;
  }

  // 否则使用当前 host
  return `${wsProtocol}//${window.location.host}${baseUrl}/admin/ws/monitor`;
}

export default monitorApi;
