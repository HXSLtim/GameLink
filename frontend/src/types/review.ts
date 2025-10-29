import type { BaseEntity } from './user';

/**
 * 评价实体 - 与后端 model.Review 保持一致
 * 使用 camelCase 命名规范
 */
export interface Review extends BaseEntity {
  orderId: number;
  reviewerId: number;
  playerId: number;
  rating: number; // 评分 (1-5)
  comment?: string;

  // 关联信息（API 返回时可能包含）
  reviewer?: {
    id: number;
    name: string;
    avatarUrl?: string;
  };
  player?: {
    id: number;
    nickname?: string;
    avatarUrl?: string;
  };
  order?: {
    id: number;
    orderNo?: string;
    title?: string;
  };
}

/**
 * 评价列表查询参数
 */
export interface ReviewListQuery {
  page?: number;
  pageSize?: number;
  orderId?: number;
  reviewerId?: number;
  playerId?: number;
  minRating?: number;
  maxRating?: number;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'rating';
  sortOrder?: 'asc' | 'desc';
}

/**
 * 创建评价请求
 */
export interface CreateReviewRequest {
  order_id: number;
  reviewer_id: number;
  player_id: number;
  rating: number;
  comment?: string;
}

/**
 * 更新评价请求
 */
export interface UpdateReviewRequest {
  rating?: number;
  comment?: string;
}

/**
 * 评价统计
 */
export interface ReviewStatistics {
  total: number;
  average_rating: number;
  rating_distribution: {
    '1': number;
    '2': number;
    '3': number;
    '4': number;
    '5': number;
  };
}

/**
 * 评分显示文本
 */
export const RATING_TEXT: Record<number, string> = {
  1: '非常差',
  2: '较差',
  3: '一般',
  4: '满意',
  5: '非常满意',
};

/**
 * 获取评分颜色（使用CSS变量，符合黑白设计风格）
 */
export const getRatingColor = (rating: number): string => {
  if (rating >= 4.5) return 'var(--rating-excellent)'; // 优秀 - 黑色/白色
  if (rating >= 3.5) return 'var(--rating-good)'; // 良好 - 深灰/浅灰
  if (rating >= 2.5) return 'var(--rating-average)'; // 一般 - 中灰
  return 'var(--rating-poor)'; // 较差 - 浅灰/深灰
};
