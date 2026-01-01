// Player Store - Taro App
// Player browsing and filtering state with offline support

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import Taro from '@tarojs/taro';
import { get, post } from '../../api/client';
import type {
  Player,
  PlayerDetail,
  PlayerFilters,
  PlayerListResponse,
  PlayerImage,
  PlayerReview,
} from '../types';

/**
 * Player store state and actions
 */
interface PlayerState {
  // Lists
  players: Player[];
  featuredPlayers: Player[];

  // Current view
  currentPlayer: PlayerDetail | null;
  playerImages: PlayerImage[];
  playerReviews: PlayerReview[];

  // Pagination
  currentPage: number;
  pageSize: number;
  total: number;
  hasMore: boolean;

  // Loading states
  loading: boolean;
  loadingDetail: boolean;
  loadingMore: boolean;
  refreshing: boolean;

  // Filters
  filters: PlayerFilters;

  // Favorites
  favoritePlayerIds: Set<number>;

  // Online status cache (playerId -> isOnline)
  onlineStatusCache: Record<number, boolean>;

  // Actions - List
  fetchPlayers: (refresh?: boolean) => Promise<void>;
  fetchFeaturedPlayers: () => Promise<void>;
  loadMorePlayers: () => Promise<void>;
  refreshPlayers: () => Promise<void>;

  // Actions - Detail
  fetchPlayerDetail: (playerId: number) => Promise<void>;
  fetchPlayerImages: (playerId: number) => Promise<void>;
  fetchPlayerReviews: (playerId: number, page?: number) => Promise<void>;

  // Actions - Filters
  setFilters: (filters: Partial<PlayerFilters>) => void;
  resetFilters: () => void;
  applyFilters: () => Promise<void>;

  // Actions - Favorites
  toggleFavorite: (playerId: number) => Promise<void>;
  isFavorite: (playerId: number) => boolean;
  fetchFavorites: () => Promise<void>;

  // Actions - Online Status
  updateOnlineStatus: (playerId: number, isOnline: boolean) => void;
  checkOnlineStatus: (playerId: number) => Promise<boolean>;

  // Selectors (computed values)
  getFilteredPlayers: () => Player[];
  getAvailablePlayers: () => Player[];
  getPlayersByPriceRange: (min: number, max: number) => Player[];
  getPlayersByGame: (game: string) => Player[];
  getPlayersByTag: (tag: string) => Player[];
  getTopRatedPlayers: (limit?: number) => Player[];
}

/**
 * Initial filters state
 */
const initialFilters: PlayerFilters = {
  minPrice: undefined,
  maxPrice: undefined,
  games: [],
  tags: [],
  rank: undefined,
  status: undefined,
  keyword: undefined,
};

/**
 * Player store with persistence for favorites and cache
 */
export const usePlayerStore = create<PlayerState>()(
  persist(
    (set, get) => ({
      // Initial state
      players: [],
      featuredPlayers: [],
      currentPlayer: null,
      playerImages: [],
      playerReviews: [],
      currentPage: 1,
      pageSize: 20,
      total: 0,
      hasMore: true,
      loading: false,
      loadingDetail: false,
      loadingMore: false,
      refreshing: false,
      filters: initialFilters,
      favoritePlayerIds: new Set<number>(),
      onlineStatusCache: {},

      /**
       * Fetch player list with pagination
       * GET /api/v1/player/list
       */
      fetchPlayers: async (refresh = false) => {
        const state = get();
        const page = refresh ? 1 : state.currentPage;

        set({ loading: true, refreshing: refresh });

        try {
          const { filters, pageSize } = state;
          const params: Record<string, any> = {
            page,
            pageSize,
          };

          // Apply filters
          if (filters.minPrice !== undefined) params.minPrice = filters.minPrice * 100; // Convert to cents
          if (filters.maxPrice !== undefined) params.maxPrice = filters.maxPrice * 100;
          if (filters.games && filters.games.length > 0) params.games = filters.games.join(',');
          if (filters.tags && filters.tags.length > 0) params.tags = filters.tags.join(',');
          if (filters.rank !== undefined) params.rank = filters.rank;
          if (filters.status) params.status = filters.status;
          if (filters.keyword) params.keyword = filters.keyword;

          const response = await get<PlayerListResponse>('/player/list', { data: params });

          if (response.success && response.data) {
            const players = response.data.players.map((player) => ({
              ...player,
              isFavorite: state.favoritePlayerIds.has(player.id),
              isOnline: state.onlineStatusCache[player.id] ?? player.isOnline ?? false,
            }));

            set({
              players: refresh ? players : [...state.players, ...players],
              currentPage: response.data.page,
              total: response.data.total,
              hasMore: players.length >= pageSize,
              loading: false,
              refreshing: false,
            });
          } else {
            throw new Error(response.message || 'Failed to fetch players');
          }
        } catch (error: any) {
          console.error('Fetch players error:', error);

          // Show error toast
          Taro.showToast({
            title: error.message || '获取陪玩师列表失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false, refreshing: false });
        }
      },

      /**
       * Fetch featured/recommended players
       * GET /api/v1/player/featured
       */
      fetchFeaturedPlayers: async () => {
        try {
          const state = get();
          const response = await get<Player[]>('/player/featured');

          if (response.success && response.data) {
            const featuredPlayers = response.data.map((player) => ({
              ...player,
              isFavorite: state.favoritePlayerIds.has(player.id),
              isOnline: state.onlineStatusCache[player.id] ?? player.isOnline ?? false,
            }));

            set({ featuredPlayers });
          }
        } catch (error: any) {
          console.error('Fetch featured players error:', error);

          // Silently fail for featured players (non-critical)
        }
      },

      /**
       * Load more players (pagination)
       */
      loadMorePlayers: async () => {
        const state = get();
        if (state.loadingMore || !state.hasMore) return;

        set({ loadingMore: true, currentPage: state.currentPage + 1 });

        try {
          await get().fetchPlayers(false);

          set({ loadingMore: false });
        } catch (error) {
          set({ loadingMore: false, currentPage: state.currentPage }); // Revert page
        }
      },

      /**
       * Refresh player list (pull-to-refresh)
       */
      refreshPlayers: async () => {
        set({ currentPage: 1 });
        await get().fetchPlayers(true);
      },

      /**
       * Fetch player detail
       * GET /api/v1/player/:id
       */
      fetchPlayerDetail: async (playerId: number) => {
        set({ loadingDetail: true, currentPlayer: null });

        try {
          const state = get();
          const response = await get<PlayerDetail>(`/player/${playerId}`);

          if (response.success && response.data) {
            const playerDetail = {
              ...response.data,
              isFavorite: state.favoritePlayerIds.has(playerId),
              isOnline: state.onlineStatusCache[playerId] ?? response.data.isOnline,
            };

            set({ currentPlayer: playerDetail, loadingDetail: false });
          } else {
            throw new Error(response.message || 'Failed to fetch player detail');
          }
        } catch (error: any) {
          console.error('Fetch player detail error:', error);

          // Show error toast
          Taro.showToast({
            title: error.message || '获取陪玩师详情失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loadingDetail: false });
        }
      },

      /**
       * Fetch player images/portfolio
       * GET /api/v1/player/:id/images
       */
      fetchPlayerImages: async (playerId: number) => {
        try {
          const response = await get<PlayerImage[]>(`/player/${playerId}/images`);

          if (response.success && response.data) {
            set({ playerImages: response.data });
          } else {
            set({ playerImages: [] });
          }
        } catch (error: any) {
          console.error('Fetch player images error:', error);
          set({ playerImages: [] });
        }
      },

      /**
       * Fetch player reviews
       * GET /api/v1/player/:id/reviews
       */
      fetchPlayerReviews: async (playerId: number, page = 1) => {
        try {
          const response = await get<PlayerReview[]>(`/player/${playerId}/reviews`, {
            data: { page, pageSize: 20 },
          });

          if (response.success && response.data) {
            set({ playerReviews: response.data });
          } else {
            set({ playerReviews: [] });
          }
        } catch (error: any) {
          console.error('Fetch player reviews error:', error);
          set({ playerReviews: [] });
        }
      },

      /**
       * Set filters
       */
      setFilters: (filters: Partial<PlayerFilters>) => {
        set((state) => ({
          filters: { ...state.filters, ...filters },
        }));
      },

      /**
       * Reset all filters
       */
      resetFilters: () => {
        set({ filters: initialFilters });
      },

      /**
       * Apply filters and fetch players
       */
      applyFilters: async () => {
        set({ currentPage: 1, players: [] });
        await get().fetchPlayers(true);
      },

      /**
       * Toggle favorite status
       * POST /api/v1/user/favorite/:playerId
       */
      toggleFavorite: async (playerId: number) => {
        const state = get();
        const isFav = state.favoritePlayerIds.has(playerId);

        try {
          // Optimistic update
          const newFavorites = new Set(state.favoritePlayerIds);
          if (isFav) {
            newFavorites.delete(playerId);
          } else {
            newFavorites.add(playerId);
          }
          set({ favoritePlayerIds: newFavorites });

          // API call
          const response = await post(isFav ? `/user/favorite/${playerId}` : `/user/favorite/${playerId}`, {
            method: isFav ? 'DELETE' : 'POST',
          });

          if (!response.success) {
            // Revert on failure
            set({ favoritePlayerIds: state.favoritePlayerIds });
            throw new Error(response.message || '操作失败');
          }

          // Update player lists with new favorite status
          set((state) => ({
            players: state.players.map((p) =>
              p.id === playerId ? { ...p, isFavorite: !isFav } : p
            ),
            featuredPlayers: state.featuredPlayers.map((p) =>
              p.id === playerId ? { ...p, isFavorite: !isFav } : p
            ),
            currentPlayer:
              state.currentPlayer?.id === playerId
                ? { ...state.currentPlayer, isFavorite: !isFav }
                : state.currentPlayer,
          }));

          // Show success toast
          Taro.showToast({
            title: isFav ? '已取消收藏' : '已收藏',
            icon: 'success',
            duration: 1500,
          });
        } catch (error: any) {
          console.error('Toggle favorite error:', error);

          // Show error toast
          Taro.showToast({
            title: error.message || '操作失败',
            icon: 'none',
            duration: 2000,
          });
        }
      },

      /**
       * Check if player is favorite
       */
      isFavorite: (playerId: number) => {
        return get().favoritePlayerIds.has(playerId);
      },

      /**
       * Fetch user's favorite players
       * GET /api/v1/user/favorites
       */
      fetchFavorites: async () => {
        try {
          const response = await get<number[]>('/user/favorites');

          if (response.success && response.data) {
            set({ favoritePlayerIds: new Set(response.data) });
          }
        } catch (error: any) {
          console.error('Fetch favorites error:', error);
          // Silently fail (non-critical)
        }
      },

      /**
       * Update online status cache
       */
      updateOnlineStatus: (playerId: number, isOnline: boolean) => {
        set((state) => ({
          onlineStatusCache: { ...state.onlineStatusCache, [playerId]: isOnline },
        }));

        // Update player lists
        set((state) => ({
          players: state.players.map((p) =>
            p.id === playerId ? { ...p, isOnline } : p
          ),
          featuredPlayers: state.featuredPlayers.map((p) =>
            p.id === playerId ? { ...p, isOnline } : p
          ),
          currentPlayer:
            state.currentPlayer?.id === playerId
              ? { ...state.currentPlayer, isOnline }
              : state.currentPlayer,
        }));
      },

      /**
       * Check player online status via API
       * GET /api/v1/player/:id/online
       */
      checkOnlineStatus: async (playerId: number) => {
        try {
          const response = await get<{ isOnline: boolean }>(`/player/${playerId}/online`);

          if (response.success && response.data) {
            get().updateOnlineStatus(playerId, response.data.isOnline);
            return response.data.isOnline;
          }
          return false;
        } catch (error) {
          console.error('Check online status error:', error);
          return false;
        }
      },

      // Selectors (computed values)

      /**
       * Get players with current filters applied
       */
      getFilteredPlayers: () => {
        const state = get();
        const { players, filters } = state;

        return players.filter((player) => {
          // Price filter
          if (filters.minPrice !== undefined && player.pricePerHour < filters.minPrice) {
            return false;
          }
          if (filters.maxPrice !== undefined && player.pricePerHour > filters.maxPrice) {
            return false;
          }

          // Rank filter
          if (filters.rank !== undefined && player.rank !== filters.rank) {
            return false;
          }

          // Status filter
          if (filters.status && player.status !== filters.status) {
            return false;
          }

          // Games filter
          if (filters.games && filters.games.length > 0) {
            if (!player.games || !filters.games.some((g) => player.games?.includes(g))) {
              return false;
            }
          }

          // Tags filter
          if (filters.tags && filters.tags.length > 0) {
            if (!filters.tags.some((t) => player.tags.includes(t))) {
              return false;
            }
          }

          // Keyword filter
          if (filters.keyword) {
            const keyword = filters.keyword.toLowerCase();
            const matchNickname = player.nickname.toLowerCase().includes(keyword);
            const matchBio = player.bio?.toLowerCase().includes(keyword);
            if (!matchNickname && !matchBio) {
              return false;
            }
          }

          return true;
        });
      },

      /**
       * Get available players only
       */
      getAvailablePlayers: () => {
        const state = get();
        return state.players.filter((player) => player.status === 'available');
      },

      /**
       * Get players by price range
       */
      getPlayersByPriceRange: (min: number, max: number) => {
        const state = get();
        return state.players.filter(
          (player) => player.pricePerHour >= min && player.pricePerHour <= max
        );
      },

      /**
       * Get players by game
       */
      getPlayersByGame: (game: string) => {
        const state = get();
        return state.players.filter((player) => player.games?.includes(game));
      },

      /**
       * Get players by tag
       */
      getPlayersByTag: (tag: string) => {
        const state = get();
        return state.players.filter((player) => player.tags.includes(tag));
      },

      /**
       * Get top rated players
       */
      getTopRatedPlayers: (limit = 10) => {
        const state = get();
        return [...state.players]
          .filter((player) => player.rating !== undefined)
          .sort((a, b) => (b.rating || 0) - (a.rating || 0))
          .slice(0, limit);
      },
    }),
    {
      name: 'player-storage',
      // Partialize: only persist favorites and online status cache
      // Don't persist: loading states, lists (fetch fresh data)
      partialize: (state) => ({
        favoritePlayerIds: Array.from(state.favoritePlayerIds), // Convert Set to Array for persistence
        onlineStatusCache: state.onlineStatusCache,
      }),
      // Hydrate: convert Array back to Set
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.favoritePlayerIds = new Set(state.favoritePlayerIds as unknown as number[]);
        }
      },
    }
  )
);

// Export types
export type { PlayerState };
