/**
 * 用户拉黑管理类型定义
 */

/**
 * 拉黑用户类型
 */
export type BlockUserType = 'user' | 'player';

/**
 * 拉黑状态
 */
export type BlockStatus = 'active' | 'canceled' | 'admin_canceled';

/**
 * 用户拉黑记录
 */
export interface UserBlock {
  id: number;
  blockerId: number;
  blockerType: BlockUserType;
  blockerName?: string;
  blockerAvatar?: string;
  blockedId: number;
  blockedType: BlockUserType;
  blockedName?: string;
  blockedAvatar?: string;
  reason?: string;
  status: BlockStatus;
  canceledAt?: string;
  adminCanceledBy?: number;
  adminCanceledByName?: string;
  adminCanceledRemark?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 拉黑查询参数
 */
export interface UserBlockQueryParams {
  page?: number;
  pageSize?: number;
  blockerId?: number;
  blockedId?: number;
  blockerType?: BlockUserType;
  blockedType?: BlockUserType;
  status?: BlockStatus;
}

/**
 * 拉黑统计
 */
export interface UserBlockStats {
  total: number;
  active: number;
  canceled: number;
  adminCanceled: number;
  todayCount: number;
  userBlocksPlayer: number;
  playerBlocksUser: number;
  userBlocksUser: number;
  playerBlocksPlayer: number;
}

/**
 * 管理员取消拉黑请求
 */
export interface AdminUnblockRequest {
  remark?: string;
}

/**
 * 批量取消拉黑请求
 */
export interface BatchUnblockRequest {
  blockIds: number[];
  remark?: string;
}

/**
 * 批量删除请求
 */
export interface BatchDeleteRequest {
  blockIds: number[];
}

/**
 * 批量拉黑请求项
 */
export interface BlockInputItem {
  blockerId: number;
  blockerType: BlockUserType;
  blockedId: number;
  blockedType: BlockUserType;
  reason?: string;
}

/**
 * 批量拉黑请求
 */
export interface BatchBlockRequest {
  blocks: BlockInputItem[];
}

/**
 * 批量操作结果
 */
export interface BatchOperationResult {
  successCount: number;
  failedCount: number;
  totalCount: number;
  failedItems: BatchOperationError[];
}

/**
 * 批量操作错误
 */
export interface BatchOperationError {
  id: number;
  message: string;
}

/**
 * 检查拉黑状态响应
 */
export interface CheckBlockStatusResponse {
  isBlocked: boolean;
  user1BlockedUser2: boolean;
  user2BlockedUser1: boolean;
}

/**
 * 拉黑状态文本映射
 */
export const BLOCK_STATUS_TEXT: Record<BlockStatus, string> = {
  active: '生效中',
  canceled: '已取消',
  admin_canceled: '管理员解除',
};

/**
 * 拉黑状态颜色映射
 */
export const BLOCK_STATUS_COLOR: Record<BlockStatus, string> = {
  active: 'error',
  canceled: 'default',
  admin_canceled: 'warning',
};

/**
 * 拉黑用户类型文本映射
 */
export const BLOCK_USER_TYPE_TEXT: Record<BlockUserType, string> = {
  user: '用户',
  player: '陪玩师',
};

/**
 * 拉黑用户类型颜色映射
 */
export const BLOCK_USER_TYPE_COLOR: Record<BlockUserType, string> = {
  user: 'blue',
  player: 'purple',
};
