/**
 * 类型定义统一导出
 */

// API 通用类型
export * from './api';

// 业务类型
// Note: ChatMessage, ChatMessageQueryParams, MuteUserRequest are only exported from chat.ts
// to avoid conflicts with the same names in content.ts
export type {
  FeedModerationStatus,
  FeedReportStatus,
  ChatMessageAuditStatus,
  ContentCategoryStatus,
  FeedReportAction,
  Feed,
  FeedQueryParams,
  FeedReport,
  FeedReportQueryParams,
  ProcessReportRequest,
  ContentCategory,
  ContentCategoryQueryParams,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  ContentStats,
  ContentTrend,
  ContentStatsDTO,
} from './content';

export {
  FEED_MODERATION_STATUS_TEXT,
  FEED_MODERATION_STATUS_COLOR,
  CHAT_AUDIT_STATUS_TEXT,
  CHAT_AUDIT_STATUS_COLOR,
  FEED_REPORT_STATUS_TEXT,
  FEED_REPORT_STATUS_COLOR,
  CATEGORY_STATUS_TEXT,
  CATEGORY_STATUS_COLOR,
  FEED_REPORT_ACTION_TEXT,
} from './content';

export * from './review';

// RBAC 权限管理类型
export * from './permission';

// Dashboard 类型 (排除与 monitor 冲突的类型)
export {
  type DashboardStats,
  type OrderStatusData,
  type TrendData,
  type PaginationParams,
  // Alert 相关类型使用 Dashboard 前缀避免冲突
  type Alert as DashboardAlert,
  type AlertLevel as DashboardAlertLevel,
  type AlertType as DashboardAlertType,
} from './dashboard';

// Monitor 类型 (使用 Monitor 前缀避免冲突)
export {
  type WSMessageType,
  type WSMessage,
  type DBConnections,
  type SystemStatus,
  type OnlineUsers,
  type OrderQueue,
  type Alert as MonitorAlert,
  type AlertLevel as MonitorAlertLevel,
} from './monitor';

// Dispute 纠纷管理类型
export * from './dispute';

// Marketing 营销模块类型
export * from './marketing';

// Payment & Operations 支付运营类型
// Note: Dispute* types, BatchOperationResult, BatchOperationError, AssignmentSource are exported from dispute.ts
// and are not re-exported from here to avoid conflicts
export type {
  CompanyType,
  CompanyStatus,
  ConditionField,
  ConditionOperator,
  RuleStatus,
  RechargeStatus,
  SettlementCompany,
  SettlementCompanyHistory,
  PlayerCompanyAssignment,
  RoutingCondition,
  CollectionEntity,
  RoutingRule,
  RoutingRuleHistory,
  RoutingTestRequest,
  RoutingTestResponse,
  RechargeOption,
  RechargeRecord,
  RechargeStats,
} from './payment';

// Game 游戏模块类型
export * from './game';

// Chat 聊天模块类型
export * from './chat';

// User Block 用户拉黑模块类型
export type {
  BlockUserType,
  BlockStatus,
  UserBlock,
  UserBlockQueryParams,
  UserBlockStats,
  AdminUnblockRequest,
  BatchUnblockRequest,
  BatchDeleteRequest,
  BlockInputItem,
  BatchBlockRequest,
  BatchOperationResult as UserBlockBatchOperationResult,
  BatchOperationError as UserBlockBatchOperationError,
  CheckBlockStatusResponse,
} from './userBlock';

export {
  BLOCK_STATUS_TEXT,
  BLOCK_STATUS_COLOR,
  BLOCK_USER_TYPE_TEXT,
  BLOCK_USER_TYPE_COLOR,
} from './userBlock';

// User 用户模块类型
export * from './user';
