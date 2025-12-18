/**
 * 评价管理模块类型定义
 */

// ==================== 评价相关类型 ====================

/** 评价状态 */
export type ReviewStatus = 'pending' | 'approved' | 'rejected' | 'deleted';

/** 评价记录 */
export interface Review {
  id: number;
  orderId: number;
  reviewerId: number;
  reviewerName?: string;
  playerId: number;
  playerName?: string;
  rating: number;
  comment: string;
  images: string[];
  status: ReviewStatus;
  isReported: boolean;
  rejectionReason?: string;
  createdAt: string;
  updatedAt: string;
  replies?: ReviewReply[];
}

/** 评价查询参数 */
export interface ReviewQueryParams {
  page?: number;
  pageSize?: number;
  orderId?: number;
  reviewerId?: number;
  playerId?: number;
  minRating?: number;
  maxRating?: number;
  status?: ReviewStatus;
  isReported?: boolean;
  startTime?: string;
  endTime?: string;
  keyword?: string;
}

/** 评价回复 */
export interface ReviewReply {
  id: number;
  reviewId: number;
  userId: number;
  userName?: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

// ==================== 举报相关类型 ====================

/** 举报状态 */
export type ReviewReportStatus = 'pending' | 'approved' | 'rejected';

/** 举报记录 */
export interface ReviewReport {
  id: number;
  reviewId: number;
  review?: Review;
  reporterId: number;
  reporterName?: string;
  reason: string;
  evidence?: string;
  status: ReviewReportStatus;
  handledBy?: number;
  handlerName?: string;
  handledAt?: string;
  handlingNote?: string;
  createdAt: string;
  updatedAt: string;
}

/** 举报查询参数 */
export interface ReviewReportQueryParams {
  page?: number;
  pageSize?: number;
  status?: ReviewReportStatus;
  reviewId?: number;
  reporterId?: number;
  startTime?: string;
  endTime?: string;
}

/** 举报处理动作 */
export type ReportHandleAction = 'delete' | 'warn' | 'reject';

// ==================== 敏感词相关类型 ====================

/** 敏感词分类 */
export type SensitiveWordCategory = 'political' | 'pornographic' | 'violent' | 'advertising' | 'other';

/** 敏感词严重程度 */
export type SensitiveWordSeverity = 'low' | 'medium' | 'high';

/** 敏感词 */
export interface SensitiveWord {
  id: number;
  word: string;
  category: SensitiveWordCategory;
  severity: SensitiveWordSeverity;
  createdAt: string;
  updatedAt: string;
}

/** 敏感词查询参数 */
export interface SensitiveWordQueryParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  category?: SensitiveWordCategory;
  severity?: SensitiveWordSeverity;
}

/** 敏感词检测结果 */
export interface DetectSensitiveWordsResult {
  hasSensitiveWords: boolean;
  detectedWords: {
    word: string;
    category: SensitiveWordCategory;
    severity: SensitiveWordSeverity;
    positions: number[];
  }[];
  highlightedContent?: string;
}

// ==================== 统计相关类型 ====================

/** 评价统计概览 */
export interface ReviewStats {
  totalCount: number;
  averageRating: number;
  ratingDistribution: {
    rating: number;
    count: number;
    percentage: number;
  }[];
  pendingCount: number;
  approvedCount: number;
  rejectedCount: number;
  reportedCount: number;
}

/** 评价趋势数据 */
export interface ReviewTrend {
  date: string;
  count: number;
  averageRating: number;
}

/** 陪玩师排行 */
export interface TopPlayer {
  playerId: number;
  playerName: string;
  avatarUrl?: string;
  reviewCount: number;
  averageRating: number;
  rank: number;
}

/** 游戏统计 */
export interface GameStats {
  gameId: number;
  gameName: string;
  gameIcon?: string;
  reviewCount: number;
  averageRating: number;
}

// ==================== 展示设置相关类型 ====================

/** 评价排序方式 */
export type ReviewSortBy = 'time' | 'score' | 'likes';

/** 评价展示设置 */
export interface ReviewDisplaySettings {
  id?: number;
  sortBy: ReviewSortBy;
  minScore: number;
  showAnonymous: boolean;
  pageSize: number;
  autoApprove: boolean;
  autoApproveMinRating: number;
  updatedAt?: string;
}

// ==================== 操作日志相关类型 ====================

/** 操作类型 */
export type OperationAction = 
  | 'create' 
  | 'approve' 
  | 'reject' 
  | 'delete' 
  | 'reply' 
  | 'update_reply' 
  | 'delete_reply'
  | 'report'
  | 'handle_report';

/** 操作日志 */
export interface OperationLog {
  id: number;
  entityType: string;
  entityId: number;
  action: OperationAction | string;
  actorId: number;
  actorName?: string;
  beforeState?: string;
  afterState?: string;
  note?: string;
  createdAt: string;
}

// ==================== 表单相关类型 ====================

/** 创建/编辑敏感词表单 */
export interface SensitiveWordFormData {
  word: string;
  category: SensitiveWordCategory;
  severity: SensitiveWordSeverity;
}

/** 拒绝评价表单 */
export interface RejectReviewFormData {
  reason: string;
}

/** 处理举报表单 */
export interface HandleReportFormData {
  action: ReportHandleAction;
  note?: string;
}

/** 更新展示设置表单 */
export interface UpdateSettingsFormData {
  sortBy: ReviewSortBy;
  minScore: number;
  showAnonymous: boolean;
  pageSize: number;
  autoApprove?: boolean;
  autoApproveMinRating?: number;
}

// ==================== 常量映射 ====================

/** 评价状态显示文本 */
export const REVIEW_STATUS_TEXT: Record<ReviewStatus, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  deleted: '已删除',
};

/** 评价状态颜色 */
export const REVIEW_STATUS_COLOR: Record<ReviewStatus, string> = {
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
  deleted: 'default',
};

/** 举报状态显示文本 */
export const REPORT_STATUS_TEXT: Record<ReviewReportStatus, string> = {
  pending: '待处理',
  approved: '已通过',
  rejected: '已驳回',
};

/** 举报状态颜色 */
export const REPORT_STATUS_COLOR: Record<ReviewReportStatus, string> = {
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
};

/** 敏感词分类显示文本 */
export const SENSITIVE_WORD_CATEGORY_TEXT: Record<SensitiveWordCategory, string> = {
  political: '政治',
  pornographic: '色情',
  violent: '暴力',
  advertising: '广告',
  other: '其他',
};

/** 敏感词分类颜色 */
export const SENSITIVE_WORD_CATEGORY_COLOR: Record<SensitiveWordCategory, string> = {
  political: 'red',
  pornographic: 'magenta',
  violent: 'volcano',
  advertising: 'orange',
  other: 'default',
};

/** 敏感词严重程度显示文本 */
export const SENSITIVE_WORD_SEVERITY_TEXT: Record<SensitiveWordSeverity, string> = {
  low: '低',
  medium: '中',
  high: '高',
};

/** 敏感词严重程度颜色 */
export const SENSITIVE_WORD_SEVERITY_COLOR: Record<SensitiveWordSeverity, string> = {
  low: 'green',
  medium: 'orange',
  high: 'red',
};

/** 排序方式显示文本 */
export const SORT_BY_TEXT: Record<ReviewSortBy, string> = {
  time: '时间',
  score: '评分',
  likes: '点赞数',
};

/** 举报处理动作显示文本 */
export const REPORT_ACTION_TEXT: Record<ReportHandleAction, string> = {
  delete: '删除评价',
  warn: '警告评价者',
  reject: '驳回举报',
};
