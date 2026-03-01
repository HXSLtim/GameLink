import { adminApi, type ApiResponse } from '@/api/admin';
import type {
  BatchNotificationRequest,
  BatchPointsRequest,
  BatchRoleRequest,
  BatchStatusRequest,
} from '@/types/user';

type VoidApiResponse = ApiResponse<void>;

const unwrapResponse = async (
  request: Promise<{ data: VoidApiResponse }>
): Promise<VoidApiResponse> => {
  const response = await request;
  return response.data;
};

export const batchUpdateUserRole = async (
  payload: BatchRoleRequest
): Promise<VoidApiResponse> => unwrapResponse(adminApi.batchUpdateUserRole(payload));

export const batchUpdateUserStatus = async (
  payload: BatchStatusRequest
): Promise<VoidApiResponse> => unwrapResponse(adminApi.batchUpdateUserStatus(payload));

export const batchAddUserPoints = async (
  payload: BatchPointsRequest
): Promise<VoidApiResponse> => unwrapResponse(adminApi.batchAddUserPoints(payload));

export const batchSendNotification = async (
  payload: BatchNotificationRequest
): Promise<VoidApiResponse> => unwrapResponse(adminApi.batchSendNotification(payload));

export const userBatchService = {
  batchUpdateRole: batchUpdateUserRole,
  batchUpdateStatus: batchUpdateUserStatus,
  batchAddPoints: batchAddUserPoints,
  batchSendNotification,
};
