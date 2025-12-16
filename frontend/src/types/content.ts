/**
 * 内容管理模块类型定义
 */

// ==================== 动态相关类型 ====================

/** 动态审核状态 */
export type FeedModerationStatus = 'pending' | 'approved' | 'rejected' | 'deleted';

/** 动态记录 */
export interface Feed {
  id: number;
  authorId: number;
  authorName?: string;
  authorAvatar?: string;
  content: string;
  images?: string[];
  categoryId?: number;
  categoryName?: string;
  moderationStatus: FeedModerationStatus;
  moderatedBy?: number;
  moderatorName?: string;
  moderatedAt?: string;
  rejectionReason?: string;
  likeCount: number;
  commentCount: number;
  shareCount: number;
  createdAt: string;
  updatedAt: string;
}

/** 动态查询参数 */
export interface FeedQueryParams {
  page?: number;
  pageSize?: number;
  authorId?: number;
  categoryId?: number;
  keyword?: string;
  moderationStatus?: FeedModerationStatus;
  dateFrom?: string;
  dateTo?: string;
}

// ==================== 聊天消息相关类型 ====================

/** 聊天消息审核状态 */
export type ChatMessageAuditStatus = 'pending' | 'approved' | 'rejected' | 'deleted';

/** 聊天消息 */
export interface ChatMessage {
  id: number;
  groupId: number;
  groupName?: string;
  senderId: number;
  senderName?: string;
  senderAvatar?: string;
  content: string;
  messageType: string;
  auditStatus: ChatMessageAuditStatus;
  isDeleted?: boolean;
  flaggedWords?: string[];
  createdAt: string;
}

/** 聊天消息查询参数 */
export interface ChatMessageQueryParams {
  page?: number;
  pageSize?: number;
  groupId?: number;
  senderId?: number;
  auditStatus?: ChatMessageAuditStatus;
  dateFrom?: string;
  dateTo?: string;
}

/** 禁言请求 */
export interface MuteUserRequest {
  groupId: number;
  userId: number;
  duration: number; // 分钟
  reason?: string;
}

// ==================== 动态举报相关类型 ====================

/** 动态举报状态 */
export type FeedReportStatus = 'pending' | 'approved' | 'rejected';

/** 动态举报记录 */
export interface FeedReport {
  id: number;
  feedId: number;
  feed?: Feed;
  reporterId: number;
  reporterName?: string;
  reason: string;
  status: FeedReportStatus;
  handledBy?: number;
  handlerName?: string;
  handledAt?: string;
  handlingNote?: string;
  createdAt: string;
  updatedAt: string;
}

/** 动态举报查询参数 */
export interface FeedReportQueryParams {
  page?: number;
  pageSize?: number;
  feedId?: number;
  reporterId?: number;
  status?: FeedReportStatus;
  dateFrom?: string;
  dateTo?: string;
}

/** 举报处理动作 */
export type FeedReportAction = 'delete' | 'warn' | 'reject';

/** 处理举报请求 */
export interface ProcessReportRequest {
  action: FeedReportAction;
  note?: string;
}

// ==================== 内容分类相关类型 ====================

/** 内容分类状态 */
export type ContentCategoryStatus = 'active' | 'inactive';

/** 内容分类 */
export interface ContentCategory {
  id: number;
  name: string;
  description?: string;
  sortOrder: number;
  status: ContentCategoryStatus;
  feedCount?: number;
  createdAt: string;
  updatedAt: string;
}

/** 内容分类查询参数 */
export interface ContentCategoryQueryParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: ContentCategoryStatus;
}

/** 创建分类请求 */
export interface CreateCategoryRequest {
  name: string;
  description?: string;
  sortOrder?: number;
  status?: ContentCategoryStatus;
}

/** 更新分类请求 */
export interface UpdateCategoryRequest {
  name?: string;
  description?: string;
  sortOrder?: number;
  status?: ContentCategoryStatus;
}

// ==================== 内容统计相关类型 ====================

/** 内容统计概览 */
export interface ContentStats {
  totalFeeds: number;
  pendingFeeds: number;
  approvedFeeds: number;
  rejectedFeeds: number;
  totalMessages: number;
  flaggedMessages: number;
  totalReports: number;
  pendingReports: number;
  reportHandleRate: number;
}

/** 内容趋势数据 */
export interface ContentTrend {
  date: string;
  feedCount: number;
  messageCount: number;
  reportCount: number;
}

/** 内容统计DTO */
export interface ContentStatsDTO {
  stats: ContentStats;
  trend: ContentTrend[];
}

// ==================== 常量映射 ====================

/** 动态审核状态显示文本 */
export const FEED_MODERATION_STATUS_TEXT: Record<FeedModerationStatus, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  deleted: '已删除',
};

/** 动态审核状态颜色 */
export const FEED_MODERATION_STATUS_COLOR: Record<FeedModerationStatus, string> = {
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
  deleted: 'default',
};

/** 聊天消息审核状态显示文本 */
export const CHAT_AUDIT_STATUS_TEXT: Record<ChatMessageAuditStatus, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  deleted: '已删除',
};

/** 聊天消息审核状态颜色 */
export const CHAT_AUDIT_STATUS_COLOR: Record<ChatMessageAuditStatus, string> = {
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
  deleted: 'default',
};

/** 动态举报状态显示文本 */
export const FEED_REPORT_STATUS_TEXT: Record<FeedReportStatus, string> = {
  pending: '待处理',
  approved: '已通过',
  rejected: '已驳回',
};

/** 动态举报状态颜色 */
export const FEED_REPORT_STATUS_COLOR: Record<FeedReportStatus, string> = {
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
};

/** 内容分类状态显示文本 */
export const CATEGORY_STATUS_TEXT: Record<ContentCategoryStatus, string> = {
  active: '启用',
  inactive: '禁用',
};

/** 内容分类状态颜色 */
export const CATEGORY_STATUS_COLOR: Record<ContentCategoryStatus, string> = {
  active: 'green',
  inactive: 'default',
};

/** 举报处理动作显示文本 */
export const FEED_REPORT_ACTION_TEXT: Record<FeedReportAction, string> = {
  delete: '删除动态',
  warn: '警告发布者',
  reject: '驳回举报',
};
