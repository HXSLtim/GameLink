/**
 * User Block API Module
 * Provides methods for managing user blocking relationships
 *
 * API Endpoints:
 * - GET    /admin/user-blocks              - List user blocks with filters
 * - GET    /admin/user-blocks/:id          - Get user block detail
 * - GET    /admin/user-blocks/stats        - Get block statistics
 * - GET    /admin/user-blocks/check        - Check block status between users
 * - GET    /admin/users/:id/blocks         - Get user's block list
 * - GET    /admin/users/:id/blocked-by     - Get users who blocked this user
 * - POST   /admin/user-blocks/:id/unblock  - Admin force unblock
 * - POST   /admin/user-blocks/batch/unblock - Batch unblock
 * - DELETE /admin/user-blocks/:id          - Delete block record
 * - POST   /admin/user-blocks/batch/delete - Batch delete blocks
 * - POST   /admin/user-blocks/batch        - Batch create blocks
 */
import apiClient from './client';
import type { ApiResponse } from '../types/api';
import type {
  UserBlock,
  UserBlockQueryParams,
  UserBlockStats,
  AdminUnblockRequest,
  BatchUnblockRequest,
  BatchDeleteRequest,
  BatchBlockRequest,
  BatchOperationResult,
  CheckBlockStatusResponse,
} from '../types/userBlock';

/**
 * User Block API
 */
export const userBlockApi = {
  /**
   * Get user block list with optional filters
   * GET /admin/user-blocks
   * @param params - Query parameters (page, pageSize, blockerId, blockedId, blockerType, blockedType, status)
   */
  getUserBlocks: (params?: UserBlockQueryParams) =>
    apiClient.get<ApiResponse<{ blocks: UserBlock[]; total: number }>>('/admin/user-blocks', {
      params,
    }),

  /**
   * Get user block detail by ID
   * GET /admin/user-blocks/:id
   * @param id - Block record ID
   */
  getUserBlockDetail: (id: number) =>
    apiClient.get<ApiResponse<UserBlock>>(`/admin/user-blocks/${id}`),

  /**
   * Get user block statistics
   * GET /admin/user-blocks/stats
   * Returns counts for total, active, canceled, admin_canceled blocks
   */
  getUserBlockStats: () =>
    apiClient.get<ApiResponse<UserBlockStats>>('/admin/user-blocks/stats'),

  /**
   * Check block status between two users
   * GET /admin/user-blocks/check?userId1=1&userId2=2
   * @param userId1 - First user ID
   * @param userId2 - Second user ID
   */
  checkBlockStatus: (userId1: number, userId2: number) =>
    apiClient.get<ApiResponse<CheckBlockStatusResponse>>('/admin/user-blocks/check', {
      params: { userId1, userId2 },
    }),

  /**
   * Get user's block list (who this user blocked)
   * GET /admin/users/:id/blocks
   * @param id - User ID
   */
  getUserBlocksByUser: (id: number) =>
    apiClient.get<ApiResponse<UserBlock[]>>(`/admin/users/${id}/blocks`),

  /**
   * Get users who blocked this user
   * GET /admin/users/:id/blocked-by
   * @param id - User ID
   */
  getUserBlockedByList: (id: number) =>
    apiClient.get<ApiResponse<UserBlock[]>>(`/admin/users/${id}/blocked-by`),

  /**
   * Admin force unblock (cancel a block)
   * POST /admin/user-blocks/:id/unblock
   * @param id - Block record ID
   * @param data - Optional remark for the unblock action
   */
  adminUnblock: (id: number, data?: AdminUnblockRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/admin/user-blocks/${id}/unblock`, data || {}),

  /**
   * Batch unblock multiple blocks
   * POST /admin/user-blocks/batch/unblock
   * @param data - Batch unblock request with block IDs and optional remark
   */
  batchUnblock: (data: BatchUnblockRequest) =>
    apiClient.post<ApiResponse<BatchOperationResult>>('/admin/user-blocks/batch/unblock', data),

  /**
   * Delete a block record
   * DELETE /admin/user-blocks/:id
   * @param id - Block record ID
   */
  deleteUserBlock: (id: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(`/admin/user-blocks/${id}`),

  /**
   * Batch delete block records
   * POST /admin/user-blocks/batch/delete
   * @param data - Batch delete request with block IDs
   */
  batchDelete: (data: BatchDeleteRequest) =>
    apiClient.post<ApiResponse<BatchOperationResult>>('/admin/user-blocks/batch/delete', data),

  /**
   * Batch create blocks
   * POST /admin/user-blocks/batch
   * @param data - Batch block request with block items
   */
  batchBlock: (data: BatchBlockRequest) =>
    apiClient.post<ApiResponse<BatchOperationResult>>('/admin/user-blocks/batch', data),
};

export default userBlockApi;
