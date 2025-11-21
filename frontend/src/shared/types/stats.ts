/**
 * 统计数据类型定义
 * 对应后端 internal/service/stats.go 的返回结构
 * 统一使用 camelCase 命名规范
 */

/**
 * Dashboard 总览统计
 * 使用 camelCase 命名规范
 */
export interface DashboardStats {
  // 总量统计
  totalUsers: number;
  totalPlayers: number;
  totalGames: number;
  totalOrders: number;
  totalPaidAmountCents: number;

  // 订单状态分布
  ordersByStatus: Record<string, number>;

  // 支付状态分布
  paymentsByStatus: Record<string, number>;
}

/**
 * 订单统计
 */
export interface OrderStatistics {
  pending?: number;
  confirmed?: number;
  in_progress?: number;
  completed?: number;
  canceled?: number;
  refunded?: number;
}

/**
 * 收入趋势数据点
 */
export interface RevenueTrendPoint {
  date: string; // 日期 (YYYY-MM-DD)
  value: number; // 收入（分）
}

/**
 * 收入趋势数据（后端直接返回数组）
 */
export type RevenueTrendData = RevenueTrendPoint[];

/**
 * 用户增长数据点
 */
export interface UserGrowthPoint {
  date: string; // 日期 (YYYY-MM-DD)
  value: number; // 新增数量
}

/**
 * 用户增长数据（后端直接返回数组）
 */
export type UserGrowthData = UserGrowthPoint[];

/**
 * TOP 陪玩师
 */
export interface TopPlayer {
  playerId: number;
  nickname: string;
  ratingAverage: number;
  ratingCount: number;
}

/**
 * TOP 陪玩师数据（后端直接返回数组）
 */
export type TopPlayersData = TopPlayer[];

/**
 * 审计总览 - 按实体类型统计
 */
export interface AuditOverviewByEntity {
  entity_type: string;
  total_actions: number;
  action_breakdown: Record<string, number>; // action -> count
}

/**
 * 审计总览数据
 */
export interface AuditOverviewData {
  by_entity: AuditOverviewByEntity[];
  total_actions: number;
  top_actors: Array<{
    actor_user_id: number;
    actor_name: string;
    action_count: number;
  }>;
  period_start?: string;
  period_end?: string;
}

/**
 * 审计趋势数据点
 */
export interface AuditTrendPoint {
  date: string; // 日期 (YYYY-MM-DD)
  action_count: number; // 操作数量
  entity_type?: string; // 实体类型（如果有筛选）
  action?: string; // 操作类型（如果有筛选）
}

/**
 * 审计趋势数据
 */
export interface AuditTrendData {
  trend: AuditTrendPoint[];
  total_actions: number;
  period_start: string;
  period_end: string;
  entity_type?: string;
  action?: string;
}

/**
 * 审计实体类型
 */
export type AuditEntityType = 'order' | 'payment' | 'player' | 'game' | 'review' | 'user';

/**
 * 统计查询参数 - 收入趋势
 */
export interface RevenueTrendQuery {
  days?: number; // 天数，默认7天
}

/**
 * 统计查询参数 - 用户增长
 */
export interface UserGrowthQuery {
  days?: number; // 天数，默认7天
}

/**
 * 统计查询参数 - TOP 陪玩师
 */
export interface TopPlayersQuery {
  limit?: number; // 数量，默认10
}

/**
 * 统计查询参数 - 审计总览
 */
export interface AuditOverviewQuery {
  date_from?: string; // 开始时间
  date_to?: string; // 结束时间
}

/**
 * 统计查询参数 - 审计趋势
 */
export interface AuditTrendQuery {
  date_from?: string; // 开始时间
  date_to?: string; // 结束时间
  entity?: AuditEntityType; // 实体类型
  action?: string; // 操作类型
}
