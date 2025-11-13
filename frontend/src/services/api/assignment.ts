import { apiClient } from '../../api/client';
import type { ListResult } from '../../types/api';
import type {
  PendingAssignment,
  AssignmentCandidate,
  AssignmentDispute,
  CancelAssignmentRequest,
  MediateDisputeRequest,
} from '../../types/assignment';
import type { AssignOrderRequest } from '../../types/order';

export type PendingAssignmentListResponse = ListResult<PendingAssignment>;

export const assignmentApi = {
  getPending: (params?: { page?: number; pageSize?: number }): Promise<PendingAssignmentListResponse> => {
    return apiClient.get('/api/v1/admin/orders/pending-assign', {
      params: {
        page: params?.page,
        page_size: params?.pageSize,
      },
    });
  },
  getCandidates: (orderId: number, params?: { limit?: number }): Promise<AssignmentCandidate[]> => {
    return apiClient.get(`/api/v1/admin/orders/${orderId}/candidates`, {
      params: {
        limit: params?.limit,
      },
    });
  },
  assign: (orderId: number, data: AssignOrderRequest): Promise<void> => {
    return apiClient.post(`/api/v1/admin/orders/${orderId}/assign`, data);
  },
  cancelAssignment: (orderId: number, data: CancelAssignmentRequest): Promise<void> => {
    return apiClient.post(`/api/v1/admin/orders/${orderId}/assign/cancel`, data);
  },
  getDisputes: (orderId: number): Promise<AssignmentDispute[]> => {
    return apiClient.get(`/api/v1/admin/orders/${orderId}/disputes`);
  },
  mediate: (orderId: number, data: MediateDisputeRequest): Promise<AssignmentDispute> => {
    return apiClient.post(`/api/v1/admin/orders/${orderId}/mediate`, data);
  },
};

export type { PendingAssignment, AssignmentCandidate, AssignmentDispute } from '../../types/assignment';
