import apiClient from './client';
import type {
  ApiResponse,
  Dispute,
  DisputeStats,
  DisputeQueryParams,
  DisputeListResponse,
  AssignDisputePayload,
  RollbackAssignmentPayload,
  ResolveDisputePayload,
  BatchAssignDisputesPayload,
  BatchUpdateDisputesStatusPayload,
  BatchCloseDisputesPayload,
  BatchOperationResult,
} from '../types';

/**
 * Dispute API Module
 * Provides methods for managing order disputes
 *
 * API Endpoints:
 * - GET    /admin/disputes             - List disputes with filters
 * - GET    /admin/disputes/pending     - List pending disputes
 * - GET    /admin/disputes/:id         - Get dispute detail
 * - GET    /admin/disputes/stats       - Get dispute statistics
 * - POST   /admin/disputes/:id/assign  - Assign dispute to CS
 * - POST   /admin/disputes/:id/rollback - Rollback assignment
 * - POST   /admin/disputes/:id/resolve - Resolve dispute
 * - POST   /admin/disputes/batch/assign - Batch assign disputes
 * - PUT    /admin/disputes/batch/status - Batch update status
 * - POST   /admin/disputes/batch/close - Batch close disputes
 */
export const disputeApi = {
  /**
   * Get dispute list with optional filters
   * GET /admin/disputes
   * @param params - Query parameters (page, pageSize, status, orderNo, initiatorType)
   */
  getDisputes: (params?: DisputeQueryParams) =>
    apiClient.get<ApiResponse<DisputeListResponse>>('/admin/disputes', { params }),

  /**
   * Get pending disputes list
   * GET /admin/disputes/pending
   * @param params - Query parameters (page, pageSize)
   */
  getPendingDisputes: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<DisputeListResponse>>('/admin/disputes/pending', { params }),

  /**
   * Get dispute detail by ID
   * GET /admin/disputes/:id
   * @param id - Dispute ID
   */
  getDisputeDetail: (id: number) =>
    apiClient.get<ApiResponse<Dispute>>(`/admin/disputes/${id}`),

  /**
   * Get dispute statistics
   * GET /admin/disputes/stats
   * Returns counts for each dispute status
   */
  getDisputeStats: () =>
    apiClient.get<ApiResponse<DisputeStats>>('/admin/disputes/stats'),

  /**
   * Assign a dispute to a customer service representative
   * POST /admin/disputes/:id/assign
   * @param id - Dispute ID
   * @param data - Assignment payload with assignedServiceId and optional originalServiceId
   */
  assignDispute: (id: number, data: AssignDisputePayload) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/admin/disputes/${id}/assign`, data),

  /**
   * Rollback a dispute assignment
   * POST /admin/disputes/:id/rollback
   * @param id - Dispute ID
   * @param data - Rollback payload with reason
   */
  rollbackAssignment: (id: number, data: RollbackAssignmentPayload) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/admin/disputes/${id}/rollback`, data),

  /**
   * Resolve a dispute with a decision
   * POST /admin/disputes/:id/resolve
   * @param id - Dispute ID
   * @param data - Resolution payload with resolution type and remark
   */
  resolveDispute: (id: number, data: ResolveDisputePayload) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/admin/disputes/${id}/resolve`, data),

  /**
   * Batch assign disputes to a CS agent
   * POST /admin/disputes/batch/assign
   * @param data - Batch assignment payload with dispute IDs and service ID
   */
  batchAssignDisputes: (data: BatchAssignDisputesPayload) =>
    apiClient.post<ApiResponse<BatchOperationResult>>('/admin/disputes/batch/assign', data),

  /**
   * Batch update disputes status
   * PUT /admin/disputes/batch/status
   * @param data - Batch status update payload with dispute IDs and new status
   */
  batchUpdateDisputesStatus: (data: BatchUpdateDisputesStatusPayload) =>
    apiClient.put<ApiResponse<BatchOperationResult>>('/admin/disputes/batch/status', data),

  /**
   * Batch close disputes with resolution
   * POST /admin/disputes/batch/close
   * @param data - Batch close payload with dispute IDs, resolution, and remark
   */
  batchCloseDisputes: (data: BatchCloseDisputesPayload) =>
    apiClient.post<ApiResponse<BatchOperationResult>>('/admin/disputes/batch/close', data),
};

export default disputeApi;
