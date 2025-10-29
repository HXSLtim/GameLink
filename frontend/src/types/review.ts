import type { BaseEntity } from './user';

/**
 * 评价实体 - 与后端 model.Review 保持一致
 */
export interface Review extends BaseEntity {
  order_id: number;
  reviewer_id: number;
  player_id: number;
  rating: number; // 评分 (1-5)
  comment?: string;

  // 关联信息（API 返回时可能包含）
  reviewer?: {
    id: number;
    name: string;
    avatar_url?: string;
  };
  player?: {
    id: number;
    nickname?: string;
    avatar_url?: string;
  };
  order?: {
    id: number;
    order_no?: string;
    title?: string;
  };
}

/**
 * 评价列表查询参数
 */
export interface ReviewListQuery {
  page?: number;
  page_size?: number;
  order_id?: number;
  reviewer_id?: number;
  player_id?: number;
  min_rating?: number;
  max_rating?: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
  sort_by?: 'created_at' | 'updated_at' | 'rating';
  sort_order?: 'asc' | 'desc';
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
