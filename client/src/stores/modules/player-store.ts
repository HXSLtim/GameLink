import { create } from 'zustand';
import { http } from '@/lib/http';

export interface Player {
    id: number;
    userId: number; // Linked user account
    username: string;
    nickname: string;
    avatar: string;
    rating: number;
    price: number;
    gameId: number;
    gameName: string;
    tags: string[];
    online: boolean;
    orderCount: number;
}

// Raw API response structure
export interface RawPlayer {
    id: number;
    userId: number;
    nickname: string;
    avatar: string;
    ratingAverage?: number;
    ratingCount?: number;
    price?: number;
    gameId?: number;
    gameName?: string;
    tags?: string[];
    onlineStatus?: string;
}

export interface PlayerResponse {
    players: RawPlayer[];
    total: number;
}

export interface PlayerFilters {
    gameId?: number;
    minPrice?: number;
    maxPrice?: number;
    onlineOnly: boolean;
    sortBy: 'rating' | 'price' | 'orders';
}

export interface PlayerState {
    // Data
    players: Player[];
    featuredPlayers: Player[];
    currentPlayer: Player | null;

    // Filters
    filters: PlayerFilters;

    // Pagination
    pagination: {
        page: number;
        pageSize: number;
        hasMore: boolean;
        total: number;
    };

    // Status
    loading: boolean;
    error: string | null;
}

export interface PlayerActions {
    fetchPlayers: (refresh?: boolean) => Promise<void>;
    fetchPlayerById: (id: number) => Promise<void>;
    fetchFeaturedPlayers: () => Promise<void>;
    setFilters: (filters: Partial<PlayerFilters>) => void;
    resetFilters: () => void;
    setPage: (page: number) => void;
}

const INITIAL_FILTERS: PlayerFilters = {
    onlineOnly: false,
    sortBy: 'rating',
};

import { subscribeWithSelector } from 'zustand/middleware';

// Debounce helper
function debounce<T extends (...args: any[]) => void>(func: T, wait: number): T {
    let timeout: ReturnType<typeof setTimeout>;
    return function (this: any, ...args: Parameters<T>) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    } as T;
}

export const usePlayerStore = create<PlayerState & PlayerActions>()(
    subscribeWithSelector((set, get) => ({
        // Initial State
        players: [],
        featuredPlayers: [],
        currentPlayer: null,
        filters: INITIAL_FILTERS,
        pagination: {
            page: 1,
            pageSize: 20,
            hasMore: true,
            total: 0,
        },
        loading: false,
        error: null,

        // Actions
        fetchPlayers: async (refresh = false) => {
            set({ loading: true, error: null });
            const { filters, pagination, players } = get();
            const currentPage = refresh ? 1 : pagination.page;

            try {
                const params = {
                    page: currentPage,
                    pageSize: pagination.pageSize,
                    ...filters,
                };

                // Use generics for strict typing
                const data = await http.get<PlayerResponse>('/public/players', { params });

                const rawPlayers = data.players || [];
                const total = data.total || 0;

                const newPlayers: Player[] = rawPlayers.map((p) => ({
                    id: p.id,
                    userId: p.userId,
                    username: p.nickname || `user_${p.id}`,
                    nickname: p.nickname,
                    avatar: p.avatar,
                    rating: p.ratingAverage || 5.0,
                    price: p.price !== undefined ? p.price : 0,
                    gameId: p.gameId || 1,
                    gameName: p.gameName || 'Valorant',
                    tags: p.tags || ['Pro', 'Friendly', 'Mic ON'],
                    online: p.onlineStatus === 'online',
                    orderCount: p.ratingCount || 0
                }));

                set({
                    players: refresh ? newPlayers : [...players, ...newPlayers],
                    pagination: {
                        ...pagination,
                        page: currentPage,
                        total,
                        hasMore: players.length + newPlayers.length < total
                    },
                    loading: false,
                });
            } catch (err) {
                const errorMessage = err instanceof Error ? err.message : 'Failed to fetch players';
                console.error("Fetch players failed", err);
                set({ loading: false, error: errorMessage });
            }
        },

        fetchPlayerById: async (id: number) => {
            set({ loading: true, error: null });
            try {
                const data = await http.get<Player>(`/public/players/${id}`);
                set({ currentPlayer: data, loading: false });
            } catch (err: any) {
                set({ loading: false, error: err.message || 'Failed to fetch player details' });
            }
        },

        fetchFeaturedPlayers: async () => {
            try {
                const data = await http.get<Player[]>('/players/featured');
                set({ featuredPlayers: data });
            } catch {
                console.warn("Failed to fetch featured players");
            }
        },

        // Debounced Filter Update
        setFilters: (newFilters) => {
            // 1. Update state immediately for UI responsiveness (optimistic)
            set((state) => ({
                filters: { ...state.filters, ...newFilters },
                pagination: { ...state.pagination, page: 1 },
                players: [] // Clear for visual feedback or keep for smooth transition? Clearing usually better for filter change
            }));

            // 2. Debounce the actual fetch to prevent spamming
            const debouncedFetch = debounce(() => {
                get().fetchPlayers(true);
            }, 500);

            debouncedFetch();
        },

        resetFilters: () => {
            set({ filters: INITIAL_FILTERS });
            get().fetchPlayers(true);
        },

        setPage: (page) => {
            set((state) => ({ pagination: { ...state.pagination, page } }));
            get().fetchPlayers(false);
        }
    }))
);
