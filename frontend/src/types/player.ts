import { BaseEntity, User } from './user';

/**
 * 认证状态枚举
 */
export type VerificationStatus = 'pending' | 'verified' | 'rejected';

/**
 * 玩家实体 - 与后端model.Player保持一致
 */
export interface Player extends BaseEntity {
  user_id: number;
  nickname?: string;
  bio?: string;
  rating_average: number;
  rating_count: number;
  hourly_rate_cents: number;
  main_game_id?: number;
  verification_status: VerificationStatus;
}

/**
 * 玩家详细信息（包含用户信息）
 */
export interface PlayerDetail extends Player {
  user?: User;
}

/**
 * 玩家列表查询参数
 */
export interface PlayerListQuery {
  page?: number;
  page_size?: number;
  user_id?: number;
  main_game_id?: number;
  verification_status?: VerificationStatus;
  min_rating?: number;
  max_rating?: number;
  min_hourly_rate?: number;
  max_hourly_rate?: number;
  keyword?: string;
  sort_by?: 'created_at' | 'updated_at' | 'rating_average' | 'hourly_rate_cents' | 'rating_count';
  sort_order?: 'asc' | 'desc';
}

/**
 * 创建玩家请求
 */
export interface CreatePlayerRequest {
  user_id: number;
  nickname?: string;
  bio?: string;
  hourly_rate_cents: number;
  main_game_id?: number;
}

/**
 * 更新玩家请求
 */
export interface UpdatePlayerRequest {
  nickname?: string;
  bio?: string;
  hourly_rate_cents?: number;
  main_game_id?: number;
  verification_status?: VerificationStatus;
}

/**
 * 认证状态显示文本
 */
export const VERIFICATION_STATUS_TEXT: Record<VerificationStatus, string> = {
  pending: '待认证',
  verified: '已认证',
  rejected: '认证失败',
};

/**
 * 认证状态徽章类型
 */
export const VERIFICATION_STATUS_BADGE: Record<
  VerificationStatus,
  'success' | 'warning' | 'error' | 'processing'
> = {
  pending: 'processing',
  verified: 'success',
  rejected: 'error',
};

/**
 * 玩家等级计算
 */
export const getPlayerLevel = (ratingCount: number): string => {
  if (ratingCount >= 100) return '金牌打手';
  if (ratingCount >= 50) return '银牌打手';
  if (ratingCount >= 20) return '铜牌打手';
  if (ratingCount >= 5) return '新手打手';
  return '见习打手';
};

/**
 * 玩家评分显示
 */
export const formatRating = (rating: number): string => {
  return rating.toFixed(1);
};

/**
 * 费率格式化（分转元）
 */
export const formatHourlyRate = (cents: number, currency: string = 'CNY'): string => {
  const yuan = cents / 100;
  return `${yuan.toFixed(2)} ${currency}/小时`;
};
