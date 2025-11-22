/**
 * 评价相关接口定义
 */

import type { BaseEntity } from '@/shared/types/api';

/**
 * 评价实体
 */
export interface Review extends BaseEntity {
  id: number;
  orderId: number;
  reviewerId: number;
  revieweeId: number;
  rating: number;
  content?: string;
  tags?: string[];
  reviewer: {
    id: number;
    username: string;
    avatar?: string;
  };
  reviewee: {
    id: number;
    username: string;
    avatar?: string;
  };
}

/**
 * 创建评价请求
 */
export interface CreateReviewRequest {
  orderId: number;
  revieweeId: number;
  rating: number;
  content?: string;
  tags?: string[];
}

/**
 * 更新评价请求
 */
export interface UpdateReviewRequest {
  rating?: number;
  content?: string;
  tags?: string[];
}

/**
 * 评价列表查询参数
 */
export interface GetReviewsParams {
  orderId?: number;
  reviewerId?: number;
  revieweeId?: number;
  rating?: number;
  page?: number;
  pageSize?: number;
}

/**
 * 评价列表响应
 */
export interface GetReviewsResponse {
  list: Review[];
  total: number;
  page: number;
  pageSize: number;
  averageRating: number;
  ratingDistribution: {
    [key: number]: number;
  };
}

/**
 * 评价统计
 */
export interface ReviewStats {
  total: number;
  averageRating: number;
  fiveStar: number;
  fourStar: number;
  threeStar: number;
  twoStar: number;
  oneStar: number;
}