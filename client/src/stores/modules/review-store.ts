import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

// ============ Interfaces ============

export interface Review {
    id: number;
    orderId: number;
    orderNo: string;
    userId: number;
    userName: string;
    userAvatar: string;
    playerId: number;
    rating: number;
    content: string;
    tags: string[];
    images?: string[];
    reply?: string;
    repliedAt?: string;
    isAnonymous: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface ReviewStats {
    totalReviews: number;
    averageRating: number;
    ratingDistribution: {
        1: number;
        2: number;
        3: number;
        4: number;
        5: number;
    };
    topTags: {
        tag: string;
        count: number;
    }[];
    recentTrend: {
        period: string;
        averageRating: number;
        count: number;
    }[];
}

export interface CreateReviewRequest {
    orderId: number;
    rating: number;
    content: string;
    tags?: string[];
    images?: string[];
    isAnonymous?: boolean;
}

export interface ReplyReviewRequest {
    content: string;
}

// ============ State & Actions ============

export interface ReviewState {
    // For players - received reviews
    receivedReviews: Review[];
    reviewStats: ReviewStats | null;

    // For users - my reviews
    myReviews: Review[];

    // For player detail page
    playerReviews: Review[];

    loading: boolean;
    error: string | null;
    pagination: {
        page: number;
        pageSize: number;
        total: number;
        hasMore: boolean;
    };
}

export interface ReviewActions {
    // Player actions - received reviews
    fetchReceivedReviews: (refresh?: boolean) => Promise<void>;
    fetchReviewStats: () => Promise<void>;
    replyToReview: (reviewId: number, request: ReplyReviewRequest) => Promise<void>;

    // User actions - my reviews
    fetchMyReviews: (refresh?: boolean) => Promise<void>;
    createReview: (request: CreateReviewRequest) => Promise<void>;
    updateReview: (reviewId: number, request: Partial<CreateReviewRequest>) => Promise<void>;
    deleteReview: (reviewId: number) => Promise<void>;

    // Public - player reviews for detail page
    fetchPlayerReviews: (playerId: number, refresh?: boolean) => Promise<void>;

    // Helpers
    canReview: (orderId: number) => Promise<boolean>;
    getAverageRating: () => number;
}

// ============ Store ============

export const useReviewStore = create<ReviewState & ReviewActions>((set, get) => ({
    receivedReviews: [],
    reviewStats: null,
    myReviews: [],
    playerReviews: [],
    loading: false,
    error: null,
    pagination: {
        page: 1,
        pageSize: 20,
        total: 0,
        hasMore: true
    },

    // ========== Player Actions ==========

    fetchReceivedReviews: async (refresh = false) => {
        set({ loading: true, error: null });
        const { pagination, receivedReviews } = get();
        const currentPage = refresh ? 1 : pagination.page;

        try {
            const data = await http.get<{ items: Review[]; total: number }>('/player/reviews', {
                params: { page: currentPage, pageSize: pagination.pageSize }
            });

            const newReviews = data.items || [];
            const total = data.total || 0;

            set({
                receivedReviews: refresh ? newReviews : [...receivedReviews, ...newReviews],
                pagination: {
                    ...pagination,
                    page: currentPage,
                    total,
                    hasMore: receivedReviews.length + newReviews.length < total
                },
                loading: false
            });
        } catch (err) {
            logError('fetchReceivedReviews', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch reviews') });
        }
    },

    fetchReviewStats: async () => {
        try {
            const data = await http.get<ReviewStats>('/player/reviews/stats');
            set({ reviewStats: data });
        } catch (err) {
            logError('fetchReviewStats', err);
        }
    },

    replyToReview: async (reviewId, request) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/player/reviews/${reviewId}/reply`, request);
            // Update the review in local state
            set((state) => ({
                receivedReviews: state.receivedReviews.map(r =>
                    r.id === reviewId
                        ? { ...r, reply: request.content, repliedAt: new Date().toISOString() }
                        : r
                ),
                loading: false
            }));
        } catch (err) {
            logError('replyToReview', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to reply to review') });
            throw err;
        }
    },

    // ========== User Actions ==========

    fetchMyReviews: async (refresh = false) => {
        set({ loading: true, error: null });
        const { pagination, myReviews } = get();
        const currentPage = refresh ? 1 : pagination.page;

        try {
            const data = await http.get<{ items: Review[]; total: number }>('/user/reviews', {
                params: { page: currentPage, pageSize: pagination.pageSize }
            });

            const newReviews = data.items || [];
            const total = data.total || 0;

            set({
                myReviews: refresh ? newReviews : [...myReviews, ...newReviews],
                pagination: {
                    ...pagination,
                    page: currentPage,
                    total,
                    hasMore: myReviews.length + newReviews.length < total
                },
                loading: false
            });
        } catch (err) {
            logError('fetchMyReviews', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch my reviews') });
        }
    },

    createReview: async (request) => {
        set({ loading: true, error: null });
        try {
            const newReview = await http.post<Review>('/user/reviews', request);
            set((state) => ({
                myReviews: [newReview, ...state.myReviews],
                loading: false
            }));
        } catch (err) {
            logError('createReview', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to create review') });
            throw err;
        }
    },

    updateReview: async (reviewId, request) => {
        set({ loading: true, error: null });
        try {
            const updatedReview = await http.put<Review>(`/user/reviews/${reviewId}`, request);
            set((state) => ({
                myReviews: state.myReviews.map(r => r.id === reviewId ? updatedReview : r),
                loading: false
            }));
        } catch (err) {
            logError('updateReview', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to update review') });
            throw err;
        }
    },

    deleteReview: async (reviewId) => {
        const previousReviews = [...get().myReviews];

        // Optimistic update
        set((state) => ({
            myReviews: state.myReviews.filter(r => r.id !== reviewId)
        }));

        try {
            await http.delete(`/user/reviews/${reviewId}`);
        } catch (err) {
            // Rollback on error
            logError('deleteReview', err);
            set({
                myReviews: previousReviews,
                error: getErrorMessage(err, 'Failed to delete review')
            });
            throw err;
        }
    },

    // ========== Public Actions ==========

    fetchPlayerReviews: async (playerId, refresh = false) => {
        set({ loading: true, error: null });
        const { pagination, playerReviews } = get();
        const currentPage = refresh ? 1 : pagination.page;

        try {
            const data = await http.get<{ items: Review[]; total: number }>(
                `/public/players/${playerId}/reviews`,
                { params: { page: currentPage, pageSize: pagination.pageSize } }
            );

            const newReviews = data.items || [];
            const total = data.total || 0;

            set({
                playerReviews: refresh ? newReviews : [...playerReviews, ...newReviews],
                pagination: {
                    ...pagination,
                    page: currentPage,
                    total,
                    hasMore: playerReviews.length + newReviews.length < total
                },
                loading: false
            });
        } catch (err) {
            logError('fetchPlayerReviews', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch player reviews') });
        }
    },

    // ========== Helpers ==========

    canReview: async (orderId) => {
        try {
            const data = await http.get<{ canReview: boolean }>(`/user/orders/${orderId}/can-review`);
            return data.canReview;
        } catch (err) {
            logError('canReview', err);
            return false;
        }
    },

    getAverageRating: () => {
        const { reviewStats } = get();
        return reviewStats?.averageRating || 0;
    }
}));
