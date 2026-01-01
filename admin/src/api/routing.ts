/**
 * 支付路由规则管理 API
 * Reference: api/internal/handler/admin/routingRule.go
 */

import apiClient from './client';
import type { ApiResponse, Pagination } from '../types/api';

// ==================== Type Definitions ====================

/**
 * 条件字段枚举
 * @Enum game_type, service_type, order_amount, region
 */
export type ConditionField = 'game_type' | 'service_type' | 'order_amount' | 'region';

/**
 * 条件操作符枚举
 * @Enum eq, neq, in, not_in, gt, lt, between
 */
export type ConditionOperator = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'lt' | 'between';

/**
 * 规则状态枚举
 */
export type RuleStatus = 'active' | 'inactive';

/**
 * 路由条件
 */
interface RoutingCondition {
  field: ConditionField;
  operator: ConditionOperator;
  value: string | number | string[] | number[];
}

/**
 * 收款主体
 */
interface CollectionEntity {
  id: number;
  name: string;
  creditCode: string;
  taxRegistrationNo?: string;
  status: 'active' | 'inactive';
  isDefault: boolean;
  totalCollectionCents: number;
  transactionCount: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * 支付路由规则
 */
interface RoutingRule {
  id: number;
  name: string;
  priority: number;
  conditions: RoutingCondition[];
  targetEntityId: number;
  status: RuleStatus;
  description?: string;
  createdBy: number;
  updatedBy?: number;
  targetEntity?: CollectionEntity;
  createdAt: string;
  updatedAt: string;
}

/**
 * 路由规则历史记录
 */
interface RoutingRuleHistory {
  id: number;
  routingRuleId: number;
  fieldName: string;
  oldValue: string;
  newValue: string;
  changedBy: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * 路由测试请求
 */
interface RoutingTestRequest {
  gameType?: string;
  serviceType?: string;
  amountCents?: number;
  region?: string;
}

/**
 * 路由测试响应
 */
interface RoutingTestResponse {
  matchedRuleId?: number;
  matchedRuleName?: string;
  collectionEntityId: number;
  entityName: string;
  merchantNo: string;
  isDefault: boolean;
  matchDetails?: RoutingCondition[];
}

/**
 * 批量操作响应
 */
interface BatchOperationResponse {
  successCount: number;
  failedCount: number;
  errors?: Array<{ id: number; error: string }>;
}

// ==================== Request/Response DTOs ====================

/**
 * 创建路由规则请求
 */
interface CreateRoutingRuleDto {
  name: string;
  priority: number;
  conditions: RoutingCondition[];
  targetEntityId: number;
  description?: string;
}

/**
 * 更新路由规则请求
 */
interface UpdateRoutingRuleDto {
  name?: string;
  priority?: number;
  conditions?: RoutingCondition[];
  targetEntityId?: number;
  description?: string;
}

/**
 * 路由规则查询参数
 */
interface RoutingRuleQueryParams {
  page?: number;
  page_size?: number;
  status?: RuleStatus;
  targetEntityId?: number;
  keyword?: string;
}

/**
 * 切换状态请求
 */
interface ToggleStatusPayload {
  enabled: boolean;
}

/**
 * 批量更新状态请求
 */
interface BatchUpdateStatusDto {
  rule_ids: number[];
  is_active: boolean;
}

/**
 * 批量删除请求
 */
interface BatchDeleteDto {
  rule_ids: number[];
}

/**
 * 设置默认主体请求
 */
interface SetDefaultEntityDto {
  entityId: number;
}

/**
 * 重新排序优先级请求
 */
interface ReorderPrioritiesDto {
  ruleIds: number[];
}

// ==================== API Functions ====================

/**
 * 支付路由规则管理 API
 */
export const routingApi = {
  /**
   * 获取路由规则列表
   * GET /admin/routing-rules
   */
  getRoutingRules: (params?: RoutingRuleQueryParams) =>
    apiClient.get<ApiResponse<RoutingRule[]>>('/admin/routing-rules', { params }),

  /**
   * 获取路由规则详情
   * GET /admin/routing-rules/:id
   */
  getRoutingRuleDetail: (id: number) =>
    apiClient.get<ApiResponse<RoutingRule>>(`/admin/routing-rules/${id}`),

  /**
   * 创建路由规则
   * POST /admin/routing-rules
   */
  createRoutingRule: (data: CreateRoutingRuleDto) =>
    apiClient.post<ApiResponse<RoutingRule>>('/admin/routing-rules', data),

  /**
   * 更新路由规则
   * PUT /admin/routing-rules/:id
   */
  updateRoutingRule: (id: number, data: UpdateRoutingRuleDto) =>
    apiClient.put<ApiResponse<RoutingRule>>(`/admin/routing-rules/${id}`, data),

  /**
   * 删除路由规则
   * DELETE /admin/routing-rules/:id
   */
  deleteRoutingRule: (id: number) =>
    apiClient.delete<ApiResponse<void>>(`/admin/routing-rules/${id}`),

  /**
   * 启用/禁用路由规则
   * POST /admin/routing-rules/:id/toggle
   */
  toggleRoutingRuleStatus: (id: number, enabled: boolean) =>
    apiClient.post<ApiResponse<void>>(`/admin/routing-rules/${id}/toggle`, { enabled }),

  /**
   * 获取路由规则修改历史
   * GET /admin/routing-rules/:id/history
   */
  getRoutingRuleHistory: (id: number) =>
    apiClient.get<ApiResponse<RoutingRuleHistory[]>>(`/admin/routing-rules/${id}/history`),

  /**
   * 批量更新路由规则状态
   * POST /admin/routing-rules/batch/status
   */
  batchUpdateRoutingRuleStatus: (data: BatchUpdateStatusDto) =>
    apiClient.post<ApiResponse<BatchOperationResponse>>('/admin/routing-rules/batch/status', data),

  /**
   * 批量删除路由规则
   * POST /admin/routing-rules/batch/delete
   */
  batchDeleteRoutingRules: (data: BatchDeleteDto) =>
    apiClient.post<ApiResponse<BatchOperationResponse>>('/admin/routing-rules/batch/delete', data),

  /**
   * 获取默认收款主体
   * GET /admin/routing-rules/default-entity
   */
  getDefaultEntity: () =>
    apiClient.get<ApiResponse<CollectionEntity>>('/admin/routing-rules/default-entity'),

  /**
   * 设置默认收款主体
   * POST /admin/routing-rules/set-default
   */
  setDefaultEntity: (data: SetDefaultEntityDto) =>
    apiClient.post<ApiResponse<void>>('/admin/routing-rules/set-default', data),

  /**
   * 测试路由规则
   * POST /admin/routing-rules/test
   */
  testRouting: (data: RoutingTestRequest) =>
    apiClient.post<ApiResponse<RoutingTestResponse>>('/admin/routing-rules/test', data),

  /**
   * 重新排序路由规则优先级
   * POST /admin/routing-rules/reorder
   */
  reorderPriorities: (data: ReorderPrioritiesDto) =>
    apiClient.post<ApiResponse<void>>('/admin/routing-rules/reorder', data),
};

// ==================== Exports ====================

export default routingApi;

// Re-export commonly used types for convenience
export type {
  RoutingRule,
  RoutingCondition,
  RoutingRuleHistory,
  CollectionEntity,
  RoutingTestRequest,
  RoutingTestResponse,
  BatchOperationResponse,
  CreateRoutingRuleDto,
  UpdateRoutingRuleDto,
  RoutingRuleQueryParams,
  ToggleStatusPayload,
  BatchUpdateStatusDto,
  BatchDeleteDto,
  SetDefaultEntityDto,
  ReorderPrioritiesDto,
};
