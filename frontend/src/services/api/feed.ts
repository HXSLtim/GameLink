import { apiClient } from '../../api/client';
import type { CreateFeedPayload, FeedListResponse, ReportFeedPayload } from '../../types/feed';

export const feedApi = {
  list: (params?: { cursor?: string; limit?: number }): Promise<FeedListResponse> => {
    return apiClient.get('/api/v1/user/feeds', { params });
  },
  create: (payload: CreateFeedPayload) => {
    return apiClient.post('/api/v1/user/feeds', payload);
  },
  report: (feedId: number, payload: ReportFeedPayload) => {
    return apiClient.post(`/api/v1/user/feeds/${feedId}/report`, payload);
  },
};
