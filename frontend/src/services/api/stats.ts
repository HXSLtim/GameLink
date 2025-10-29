import { apiClient } from '../../api/client';
import type {
  DashboardStats,
  OrderStatistics,
  RevenueTrendData,
  RevenueTrendQuery,
  UserGrowthData,
  UserGrowthQuery,
  TopPlayersData,
  TopPlayersQuery,
  AuditOverviewData,
  AuditOverviewQuery,
  AuditTrendData,
  AuditTrendQuery,
} from '../../types/stats';

/**
 * 统计 API 服务
 * 对应后端 /admin/stats/* 路由
 */
export const statsApi = {
  /**
   * 获取 Dashboard 总览统计
   * GET /admin/stats/dashboard
   */
  getDashboard: (): Promise<DashboardStats> => {
    return apiClient.get('/api/v1/admin/stats/dashboard');
  },

  /**
   * 获取订单状态统计
   * GET /admin/stats/orders
   */
  getOrderStats: (): Promise<OrderStatistics> => {
    return apiClient.get('/api/v1/admin/stats/orders');
  },

  /**
   * 获取收入趋势（按日）
   * GET /admin/stats/revenue-trend
   * @param query.days - 天数，默认7天
   */
  getRevenueTrend: (query?: RevenueTrendQuery): Promise<RevenueTrendData> => {
    return apiClient.get('/api/v1/admin/stats/revenue-trend', {
      params: query,
    });
  },

  /**
   * 获取用户增长趋势（按日）
   * GET /admin/stats/user-growth
   * @param query.days - 天数，默认7天
   */
  getUserGrowth: (query?: UserGrowthQuery): Promise<UserGrowthData> => {
    return apiClient.get('/api/v1/admin/stats/user-growth', {
      params: query,
    });
  },

  /**
   * 获取 TOP 陪玩师排行
   * GET /admin/stats/top-players
   * @param query.limit - 数量，默认10
   */
  getTopPlayers: (query?: TopPlayersQuery): Promise<TopPlayersData> => {
    return apiClient.get('/api/v1/admin/stats/top-players', {
      params: query,
    });
  },

  /**
   * 获取审计总览（按实体/动作汇总）
   * GET /admin/stats/audit/overview
   * @param query.date_from - 开始时间
   * @param query.date_to - 结束时间
   */
  getAuditOverview: (query?: AuditOverviewQuery): Promise<AuditOverviewData> => {
    return apiClient.get('/api/v1/admin/stats/audit/overview', {
      params: query,
    });
  },

  /**
   * 获取审计趋势（按日）
   * GET /admin/stats/audit/trend
   * @param query.date_from - 开始时间
   * @param query.date_to - 结束时间
   * @param query.entity - 实体类型 (order/payment/player/game/review/user)
   * @param query.action - 操作类型
   */
  getAuditTrend: (query?: AuditTrendQuery): Promise<AuditTrendData> => {
    return apiClient.get('/api/v1/admin/stats/audit/trend', {
      params: query,
    });
  },
};

// 导出类型以便组件使用
export type {
  DashboardStats,
  OrderStatistics,
  RevenueTrendData,
  UserGrowthData,
  TopPlayersData,
  AuditOverviewData,
  AuditTrendData,
} from '../../types/stats';
