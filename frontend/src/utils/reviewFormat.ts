/**
 * 评价管理模块格式化工具
 */
import dayjs from 'dayjs';
import type {
  ReviewStatus,
  ReviewReportStatus,
  SensitiveWordCategory,
  SensitiveWordSeverity,
  ReviewSortBy,
} from '@/types/review';
import {
  REVIEW_STATUS_TEXT,
  REVIEW_STATUS_COLOR,
  REPORT_STATUS_TEXT,
  REPORT_STATUS_COLOR,
  SENSITIVE_WORD_CATEGORY_TEXT,
  SENSITIVE_WORD_CATEGORY_COLOR,
  SENSITIVE_WORD_SEVERITY_TEXT,
  SENSITIVE_WORD_SEVERITY_COLOR,
  SORT_BY_TEXT,
} from '@/types/review';

/**
 * 格式化评分为星级显示文本
 * @param rating 评分 (1-5)
 * @returns 星级文本，如 "★★★★★"
 */
export const formatRatingStars = (rating: number): string => {
  const fullStars = Math.floor(rating);
  const emptyStars = 5 - fullStars;
  return '★'.repeat(fullStars) + '☆'.repeat(emptyStars);
};

/**
 * 格式化评分为数字显示
 * @param rating 评分
 * @param precision 小数位数
 * @returns 格式化后的评分，如 "4.5分"
 */
export const formatRatingNumber = (rating: number, precision = 1): string => {
  return `${rating.toFixed(precision)}分`;
};

/**
 * 格式化日期时间
 * @param dateStr 日期字符串
 * @param format 格式化模板
 * @returns 格式化后的日期时间
 */
export const formatDateTime = (
  dateStr: string | undefined | null,
  format = 'YYYY-MM-DD HH:mm:ss'
): string => {
  if (!dateStr) return '-';
  return dayjs(dateStr).format(format);
};

/**
 * 格式化日期
 * @param dateStr 日期字符串
 * @returns 格式化后的日期
 */
export const formatDate = (dateStr: string | undefined | null): string => {
  return formatDateTime(dateStr, 'YYYY-MM-DD');
};

/**
 * 格式化相对时间
 * @param dateStr 日期字符串
 * @returns 相对时间，如 "3天前"
 */
export const formatRelativeTime = (dateStr: string | undefined | null): string => {
  if (!dateStr) return '-';
  const date = dayjs(dateStr);
  const now = dayjs();
  const diffMinutes = now.diff(date, 'minute');
  const diffHours = now.diff(date, 'hour');
  const diffDays = now.diff(date, 'day');

  if (diffMinutes < 1) return '刚刚';
  if (diffMinutes < 60) return `${diffMinutes}分钟前`;
  if (diffHours < 24) return `${diffHours}小时前`;
  if (diffDays < 30) return `${diffDays}天前`;
  return formatDate(dateStr);
};

/**
 * 获取评价状态文本
 * @param status 评价状态
 * @returns 状态文本
 */
export const getReviewStatusText = (status: ReviewStatus): string => {
  return REVIEW_STATUS_TEXT[status] || status;
};

/**
 * 获取评价状态颜色
 * @param status 评价状态
 * @returns 颜色值
 */
export const getReviewStatusColor = (status: ReviewStatus): string => {
  return REVIEW_STATUS_COLOR[status] || 'default';
};

/**
 * 获取举报状态文本
 * @param status 举报状态
 * @returns 状态文本
 */
export const getReportStatusText = (status: ReviewReportStatus): string => {
  return REPORT_STATUS_TEXT[status] || status;
};

/**
 * 获取举报状态颜色
 * @param status 举报状态
 * @returns 颜色值
 */
export const getReportStatusColor = (status: ReviewReportStatus): string => {
  return REPORT_STATUS_COLOR[status] || 'default';
};

/**
 * 获取敏感词分类文本
 * @param category 分类
 * @returns 分类文本
 */
export const getSensitiveWordCategoryText = (category: SensitiveWordCategory): string => {
  return SENSITIVE_WORD_CATEGORY_TEXT[category] || category;
};

/**
 * 获取敏感词分类颜色
 * @param category 分类
 * @returns 颜色值
 */
export const getSensitiveWordCategoryColor = (category: SensitiveWordCategory): string => {
  return SENSITIVE_WORD_CATEGORY_COLOR[category] || 'default';
};

/**
 * 获取敏感词严重程度文本
 * @param severity 严重程度
 * @returns 严重程度文本
 */
export const getSensitiveWordSeverityText = (severity: SensitiveWordSeverity): string => {
  return SENSITIVE_WORD_SEVERITY_TEXT[severity] || severity;
};

/**
 * 获取敏感词严重程度颜色
 * @param severity 严重程度
 * @returns 颜色值
 */
export const getSensitiveWordSeverityColor = (severity: SensitiveWordSeverity): string => {
  return SENSITIVE_WORD_SEVERITY_COLOR[severity] || 'default';
};

/**
 * 获取排序方式文本
 * @param sortBy 排序方式
 * @returns 排序方式文本
 */
export const getSortByText = (sortBy: ReviewSortBy): string => {
  return SORT_BY_TEXT[sortBy] || sortBy;
};

/**
 * 截断文本
 * @param text 原始文本
 * @param maxLength 最大长度
 * @returns 截断后的文本
 */
export const truncateText = (text: string | undefined | null, maxLength = 50): string => {
  if (!text) return '-';
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + '...';
};

/**
 * 格式化图片数量
 * @param count 图片数量
 * @returns 格式化后的文本
 */
export const formatImageCount = (count: number): string => {
  if (count === 0) return '无图片';
  return `${count}张图片`;
};

/**
 * 格式化百分比
 * @param value 数值
 * @param precision 小数位数
 * @returns 格式化后的百分比
 */
export const formatPercentage = (value: number, precision = 1): string => {
  return `${value.toFixed(precision)}%`;
};
