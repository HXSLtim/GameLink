/**
 * Dispute Management Types
 * Matches backend model.OrderDispute and related structures
 */

/**
 * Dispute Status Enum
 * Matches backend model.DisputeStatus
 */
export type DisputeStatus =
  | 'pending'   // 待处理
  | 'assigned'  // 已指派
  | 'mediating' // 调解中
  | 'resolved'  // 已解决
  | 'rejected'  // 已驳回
  | 'canceled'; // 已取消

/**
 * Dispute Resolution Enum
 * Matches backend model.DisputeResolution
 */
export type DisputeResolution =
  | 'refund'    // 全额退款
  | 'partial'   // 部分退款
  | 'reassign'  // 重新指派
  | 'reject'    // 驳回
  | 'pending';  // 待决定

/**
 * Dispute Initiator Type
 * Matches backend model.DisputeInitiatorType
 */
export type DisputeInitiatorType = 'user' | 'player';

/**
 * Dispute Type
 * Matches backend model.DisputeType
 */
export type DisputeType =
  | 'service_quality'      // 服务质量问题
  | 'bad_attitude'         // 态度问题
  | 'incomplete_service'   // 未完成服务
  | 'user_not_cooperative' // 用户不配合/不听指挥
  | 'user_harassment'      // 用户骚扰
  | 'other';               // 其他

/**
 * Assignment Source
 * Matches backend model.AssignmentSource
 */
export type AssignmentSource = 'system' | 'manual' | 'team';

/**
 * Order Dispute Interface
 * Matches backend model.OrderDispute
 */
export interface Dispute {
  id: number;
  orderId: number;
  orderNo?: string;
  initiatorId: number;
  initiatorName?: string;
  initiatorType: DisputeInitiatorType;
  type: DisputeType;
  status: DisputeStatus;
  reason: string;
  evidenceUrls?: string[];
  evidenceText?: string;
  chatSnapshotId?: number;

  // 双客服机制 (Dual-CS Mechanism)
  originalServiceId?: number;
  originalServiceName?: string;
  assignedServiceId?: number;
  assignedServiceName?: string;

  // SLA 信息 (SLA Information)
  slaDeadline?: string;
  slaBreached: boolean;
  slaBreachedAt?: string;

  // 处理信息 (Resolution Information)
  resolution?: DisputeResolution;
  resolvedBy?: number;
  resolvedByName?: string;
  resolvedAt?: string;
  resolveRemark?: string;

  // 回退信息 (Rollback Information)
  rolledBackAt?: string;
  rolledBackByUserId?: number;
  rollbackReason?: string;

  // 追踪信息 (Trace Information)
  traceId: string;

  createdAt: string;
  updatedAt: string;
}

/**
 * Dispute Statistics
 */
export interface DisputeStats {
  total: number;
  pending: number;
  assigned: number;
  mediating: number;
  resolved: number;
  rejected: number;
  canceled: number;
  slaBreached: number;
}

/**
 * Dispute List Query Parameters
 */
export interface DisputeQueryParams {
  page?: number;
  pageSize?: number;
  status?: DisputeStatus;
  orderNo?: string;
  initiatorType?: DisputeInitiatorType;
}

/**
 * Dispute List Response
 */
export interface DisputeListResponse {
  disputes: Dispute[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * Assign Dispute Payload
 */
export interface AssignDisputePayload {
  assignedServiceId: number;
  originalServiceId?: number;
}

/**
 * Rollback Assignment Payload
 */
export interface RollbackAssignmentPayload {
  rollbackReason: string;
}

/**
 * Resolve Dispute Payload
 */
export interface ResolveDisputePayload {
  resolution: DisputeResolution;
  resolveRemark: string;
}

/**
 * Batch Assign Disputes Payload
 */
export interface BatchAssignDisputesPayload {
  disputeIds: number[];
  assignedServiceId: number;
  originalServiceId?: number;
}

/**
 * Batch Update Disputes Status Payload
 */
export interface BatchUpdateDisputesStatusPayload {
  disputeIds: number[];
  status: 'assigned' | 'mediating' | 'canceled';
}

/**
 * Batch Close Disputes Payload
 */
export interface BatchCloseDisputesPayload {
  disputeIds: number[];
  resolution: 'refund' | 'partial' | 'reject';
  resolveRemark: string;
}

/**
 * Batch Operation Result
 */
export interface BatchOperationResult {
  success: boolean;
  message: string;
  successCount: number;
  failedCount: number;
  errors?: BatchOperationError[];
}

/**
 * Batch Operation Error
 */
export interface BatchOperationError {
  disputeId: number;
  error: string;
}

/**
 * Dispute Type Display Names (for UI rendering)
 */
export const DISPUTE_TYPE_LABELS: Record<DisputeType, string> = {
  service_quality: '服务质量问题',
  bad_attitude: '态度问题',
  incomplete_service: '未完成服务',
  user_not_cooperative: '用户不配合/不听指挥',
  user_harassment: '用户骚扰',
  other: '其他',
};

/**
 * Dispute Status Display Names (for UI rendering)
 */
export const DISPUTE_STATUS_LABELS: Record<DisputeStatus, string> = {
  pending: '待处理',
  assigned: '已指派',
  mediating: '调解中',
  resolved: '已解决',
  rejected: '已驳回',
  canceled: '已取消',
};

/**
 * Dispute Resolution Display Names (for UI rendering)
 */
export const DISPUTE_RESOLUTION_LABELS: Record<DisputeResolution, string> = {
  refund: '全额退款',
  partial: '部分退款',
  reassign: '重新指派',
  reject: '驳回',
  pending: '待决定',
};

/**
 * Dispute Status Colors (for Ant Design Tags)
 */
export const DISPUTE_STATUS_COLORS: Record<DisputeStatus, string> = {
  pending: 'warning',
  assigned: 'processing',
  mediating: 'blue',
  resolved: 'success',
  rejected: 'error',
  canceled: 'default',
};
