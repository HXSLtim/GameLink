import { create } from 'zustand';
import { http } from '@/lib/http';

// ============================================================================
// Types
// ============================================================================

export type LFGRequestType = 'find_player' | 'find_team';
export type LFGRequestStatus = 'pending' | 'matched' | 'expired' | 'canceled';

export interface LFGRequest {
    id: number;
    userId: number;
    userNickname?: string;
    userAvatarUrl?: string;
    gameId: number;
    gameName?: string;
    requestType: LFGRequestType;
    title: string;
    description?: string;
    requiredPlayers: number;
    minRank?: string;
    maxPriceCents?: number;
    status: LFGRequestStatus;
    expiresAt: string;
    matchedRoomId?: number;
    createdAt: string;
    updatedAt: string;
}

export interface CreateLFGRequest {
    gameId: number;
    requestType: LFGRequestType;
    title?: string;
    description?: string;
    requiredPlayers?: number;
    minRank?: string;
    maxPriceCents?: number;
    expireMinutes?: number;
}

export interface LFGListOptions {
    page?: number;
    pageSize?: number;
    gameId?: number;
    requestType?: LFGRequestType;
}

export interface LFGState {
    // State
    requests: LFGRequest[];
    myRequests: LFGRequest[];
    activeRequest: LFGRequest | null;
    matches: LFGRequest[];
    pendingCount: number;
    isLoading: boolean;
    error: string | null;
    pagination: {
        page: number;
        pageSize: number;
        total: number;
    };

    // Actions
    fetchRequests: (options?: LFGListOptions) => Promise<void>;
    fetchPendingRequests: (gameId?: number) => Promise<void>;
    fetchMyRequests: (status?: LFGRequestStatus) => Promise<void>;
    fetchActiveRequest: () => Promise<void>;
    fetchRequest: (requestId: number) => Promise<LFGRequest>;
    createRequest: (data: CreateLFGRequest) => Promise<LFGRequest>;
    cancelRequest: (requestId: number) => Promise<void>;
    acceptRequest: (requestId: number) => Promise<any>;
    findMatches: (requestId: number, limit?: number) => Promise<void>;
    fetchPendingCount: (gameId?: number) => Promise<void>;
    reset: () => void;
}

// ============================================================================
// Initial State
// ============================================================================

const initialState = {
    requests: [],
    myRequests: [],
    activeRequest: null,
    matches: [],
    pendingCount: 0,
    isLoading: false,
    error: null,
    pagination: {
        page: 1,
        pageSize: 20,
        total: 0,
    },
};

// ============================================================================
// Store
// ============================================================================

export const useLFGStore = create<LFGState>((set) => ({
    ...initialState,

    // Fetch LFG requests list
    fetchRequests: async (options?: LFGListOptions) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams();
            if (options?.page) params.append('page', String(options.page));
            if (options?.pageSize) params.append('pageSize', String(options.pageSize));
            if (options?.gameId) params.append('gameId', String(options.gameId));
            if (options?.requestType) params.append('requestType', options.requestType);

            const response = await http.get<{
                items: LFGRequest[];
                pagination: { page: number; pageSize: number; total: number };
            }>(`/user/lfg?${params.toString()}`);

            set({
                requests: response.items || [],
                pagination: response.pagination || initialState.pagination,
                isLoading: false,
            });
        } catch (error: any) {
            const message = error?.message || 'Failed to fetch LFG requests';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch pending requests
    fetchPendingRequests: async (gameId?: number) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams();
            if (gameId) params.append('gameId', String(gameId));

            const response = await http.get<{
                items: LFGRequest[];
                pagination: { page: number; pageSize: number; total: number };
            }>(`/user/lfg/pending?${params.toString()}`);

            set({
                requests: response.items || [],
                pagination: response.pagination || initialState.pagination,
                isLoading: false,
            });
        } catch (error: any) {
            const message = error?.message || 'Failed to fetch pending requests';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch my requests
    fetchMyRequests: async (status?: LFGRequestStatus) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams();
            if (status) params.append('status', status);

            const response = await http.get<LFGRequest[]>(
                `/user/lfg/my?${params.toString()}`
            );
            set({ myRequests: response || [], isLoading: false });
        } catch (error: any) {
            const message = error?.message || 'Failed to fetch my requests';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch active request
    fetchActiveRequest: async () => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.get<LFGRequest>('/user/lfg/active');
            set({ activeRequest: response, isLoading: false });
        } catch (error: any) {
            // 404 means no active request, which is fine
            if (error?.status === 404) {
                set({ activeRequest: null, isLoading: false });
                return;
            }
            const message = error?.message || 'Failed to fetch active request';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch single request
    fetchRequest: async (requestId: number) => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.get<LFGRequest>(`/user/lfg/${requestId}`);
            set({ isLoading: false });
            return response;
        } catch (error: any) {
            const message = error?.message || 'Failed to fetch request';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Create LFG request
    createRequest: async (data: CreateLFGRequest) => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.post<LFGRequest>('/user/lfg', data);
            set((state) => ({
                requests: [response, ...state.requests],
                myRequests: [response, ...state.myRequests],
                activeRequest: response,
                isLoading: false,
            }));
            return response;
        } catch (error: any) {
            const message = error?.message || 'Failed to create request';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Cancel request
    cancelRequest: async (requestId: number) => {
        try {
            await http.delete(`/user/lfg/${requestId}`);
            set((state) => ({
                requests: state.requests.filter((r) => r.id !== requestId),
                myRequests: state.myRequests.filter((r) => r.id !== requestId),
                activeRequest:
                    state.activeRequest?.id === requestId
                        ? null
                        : state.activeRequest,
            }));
        } catch (error: any) {
            const message = error?.message || 'Failed to cancel request';
            set({ error: message });
            throw error;
        }
    },

    // Accept request (player accepts a user's LFG request)
    acceptRequest: async (requestId: number) => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.post(`/user/lfg/${requestId}/accept`);
            // Remove from pending list
            set((state) => ({
                requests: state.requests.filter((r) => r.id !== requestId),
                isLoading: false,
            }));
            return response;
        } catch (error: any) {
            const message = error?.message || 'Failed to accept request';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Find matches for a request
    findMatches: async (requestId: number, limit?: number) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams();
            if (limit) params.append('limit', String(limit));

            const response = await http.get<LFGRequest[]>(
                `/user/lfg/${requestId}/matches?${params.toString()}`
            );
            set({ matches: response || [], isLoading: false });
        } catch (error: any) {
            const message = error?.message || 'Failed to find matches';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch pending count
    fetchPendingCount: async (gameId?: number) => {
        try {
            const params = new URLSearchParams();
            if (gameId) params.append('gameId', String(gameId));

            const response = await http.get<{ count: number }>(
                `/user/lfg/count?${params.toString()}`
            );
            set({ pendingCount: response.count || 0 });
        } catch (error: any) {
            console.error('[LFG] Failed to fetch pending count:', error);
        }
    },

    // Reset state
    reset: () => {
        set(initialState);
    },
}));

// ============================================================================
// Selectors
// ============================================================================

export const selectLFGRequests = (state: LFGState) => state.requests;
export const selectMyLFGRequests = (state: LFGState) => state.myRequests;
export const selectActiveRequest = (state: LFGState) => state.activeRequest;
export const selectLFGMatches = (state: LFGState) => state.matches;
export const selectPendingCount = (state: LFGState) => state.pendingCount;
export const selectLFGIsLoading = (state: LFGState) => state.isLoading;
export const selectLFGError = (state: LFGState) => state.error;
