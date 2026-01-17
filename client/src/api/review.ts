/**
 * Review API
 * Handles order reviews, ratings, comments
 */

import { http } from '@/lib/http';
import type {
    Review,
    CreateReviewRequest,
    ReviewListParams,
    PaginatedResponse
} from '@/types/api';

export const reviewApi = {
    /**
     * Create review for order
     */
    create: (data: CreateReviewRequest) =>
        http.post<Review>('/review/create', data),

    /**
     * Get review list
     */
    list: (params: ReviewListParams) =>
        http.get<PaginatedResponse<Review>>('/review/list', { params }),

    /**
     * Get review detail
     */
    get: (id: number) =>
        http.get<Review>(`/review/${id}`),

    /**
     * Get reviews for player
     */
    getPlayerReviews: (playerId: number, params: {
        page: number;
        pageSize: number;
        rating?: number;
    }) =>
        http.get<PaginatedResponse<Review>>(`/review/player/${playerId}`, { params }),

    /**
     * Get reviews by user
     */
    getUserReviews: (params: { page: number; pageSize: number }) =>
        http.get<PaginatedResponse<Review>>('/review/user', { params }),

    /**
     * Update review
     */
    update: (id: number, data: Partial<CreateReviewRequest>) =>
        http.put<Review>(`/review/${id}`, data),

    /**
     * Delete review
     */
    delete: (id: number) =>
        http.delete<void>(`/review/${id}`),

    /**
     * Reply to review (player only)
     */
    reply: (id: number, content: string) =>
        http.post<void>(`/review/${id}/reply`, { content }),

    /**
     * Upload review images
     */
    uploadImages: (files: File[]) => {
        const formData = new FormData();
        files.forEach(file => formData.append('images', file));
        return http.post<{ urls: string[] }>('/review/upload/images', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },
};
