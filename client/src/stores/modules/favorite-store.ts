import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

export interface FavoritePlayer {
    id: number;
    nickname: string;
    avatar: string;
    gameId?: number;
    bio?: string;
    rating?: number;
}

export interface FavoriteStore {
    favorites: FavoritePlayer[];
    loading: boolean;
    error: string | null;
    toggleFavorite: (player: FavoritePlayer) => Promise<void>;
    checkIsFavorite: (playerId: number) => boolean;
    fetchFavorites: () => Promise<void>;
    removeFavorite: (playerId: number) => Promise<void>;
}

export const useFavoriteStore = create<FavoriteStore>()(
    persist(
        (set, get) => ({
            favorites: [],
            loading: false,
            error: null,

            fetchFavorites: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<{ items: FavoritePlayer[] }>('/user/favorites');
                    set({ favorites: data.items || [], loading: false });
                } catch (err) {
                    set({
                        loading: false,
                        error: err instanceof Error ? err.message : 'Failed to fetch favorites'
                    });
                }
            },

            toggleFavorite: async (player) => {
                const { favorites } = get();
                const isFav = favorites.some(f => f.id === player.id);
                const previousFavorites = [...favorites];

                // Optimistic update
                if (isFav) {
                    set({ favorites: favorites.filter(f => f.id !== player.id) });
                } else {
                    set({ favorites: [...favorites, player] });
                }

                try {
                    if (isFav) {
                        await http.delete(`/user/favorites/${player.id}`);
                    } else {
                        await http.post('/user/favorites', { playerId: player.id });
                    }
                } catch (err) {
                    // Rollback on error
                    set({
                        favorites: previousFavorites,
                        error: err instanceof Error ? err.message : 'Failed to toggle favorite'
                    });
                }
            },

            removeFavorite: async (playerId) => {
                const { favorites } = get();
                const previousFavorites = [...favorites];

                // Optimistic update
                set({ favorites: favorites.filter(f => f.id !== playerId) });

                try {
                    await http.delete(`/user/favorites/${playerId}`);
                } catch (err) {
                    // Rollback on error
                    set({
                        favorites: previousFavorites,
                        error: err instanceof Error ? err.message : 'Failed to remove favorite'
                    });
                }
            },

            checkIsFavorite: (playerId) => {
                return get().favorites.some(f => f.id === playerId);
            }
        }),
        {
            name: 'favorite-storage',
            partialize: (state) => ({
                favorites: state.favorites
            })
        }
    )
);
