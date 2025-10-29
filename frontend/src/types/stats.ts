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
  total: number;
  pending: number;
  confirmed: number;
  in_progress: number;
  completed: number;
  cancelled: number;
  today_orders: number;
  today_revenue: number;
}

/**
 * 收入趋势数据点
 */
export interface RevenueTrendPoint {
  date: string; // 日期 (YYYY-MM-DD)
  revenue_cents: number; // 收入（分）
  order_count: number; // 订单数量
  avg_order_value: number; // 平均订单价值
}

/**
 * 收入趋势数据
 */
export interface RevenueTrendData {
  trend: RevenueTrendPoint[];
  total_revenue_cents: number;
  total_orders: number;
  avg_daily_revenue_cents: number;
  period_start: string;
  period_end: string;
}

/**
 * 用户增长数据点
 */
export interface UserGrowthPoint {
  date: string; // 日期 (YYYY-MM-DD)
  new_users: number; // 新增用户
  new_players: number; // 新增陪玩师
  total_users: number; // 累计用户
  total_players: number; // 累计陪玩师
}

/**
 * 用户增长数据
 */
export interface UserGrowthData {
  trend: UserGrowthPoint[];
  total_new_users: number;
  total_new_players: number;
  avg_daily_new_users: number;
  period_start: string;
  period_end: string;
}

/**
 * TOP 陪玩师
 */
export interface TopPlayer {
  player_id: number;
  user_id: number;
  player_name: string;
  avatar_url?: string;
  order_count: number;
  total_revenue_cents: number;
  avg_rating: number;
  review_count: number;
  completion_rate: number;
  main_game?: string;
}

/**
 * TOP 陪玩师数据
 */
export interface TopPlayersData {
  players: TopPlayer[];
  period_start: string;
  period_end: string;
}

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
