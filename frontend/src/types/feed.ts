export interface FeedImage {
  url: string;
  width?: number;
  height?: number;
  sizeBytes?: number;
  order?: number;
}

export type FeedVisibility = 'public' | 'followers' | 'private';

export interface Feed {
  id: number;
  authorId: number;
  content: string;
  visibility: FeedVisibility;
  moderationStatus: 'pending' | 'approved' | 'rejected' | 'removed';
  moderationNote?: string;
  createdAt: string;
  images: FeedImage[];
}

export interface FeedListResponse {
  items: Feed[];
  nextCursor?: string;
}

export interface CreateFeedPayload {
  content: string;
  visibility?: FeedVisibility;
  images: { url: string; sizeBytes?: number }[];
}

export interface ReportFeedPayload {
  reason: string;
}
