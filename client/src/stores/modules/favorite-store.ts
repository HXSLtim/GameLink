import { create } from 'zustand';
import { persist } from 'zustand/middleware';
// import { http } from '@/lib/http';

export interface FavoritePlayer {
    id: number;
    nickname: string;
    avatar: string;
    gameId?: number;
}

export interface FavoriteStore {
    favorites: FavoritePlayer[];
    loading: boolean;
    toggleFavorite: (player: FavoritePlayer) => Promise<void>;
    checkIsFavorite: (playerId: number) => boolean;
    fetchFavorites: () => Promise<void>;
}

export const useFavoriteStore = create<FavoriteStore>()(
    persist(
        (set, get) => ({
            favorites: [],
            loading: false,

            fetchFavorites: async () => {
                set({ loading: true });
                try {
                    // await http.get('/favorites');
                    await new Promise(resolve => setTimeout(resolve, 300));
                    set({ loading: false });
                } catch (e) {
                    set({ loading: false });
                }
            },

            toggleFavorite: async (player) => {
                const { favorites } = get();
                const isFav = favorites.some(f => f.id === player.id);

                // Optimistic update
                if (isFav) {
                    set({ favorites: favorites.filter(f => f.id !== player.id) });
                } else {
                    set({ favorites: [...favorites, player] });
                }

                try {
                    // await http.post('/favorites/toggle', { playerId: player.id });
                } catch (e) {
                    // Rollback on error
                    set({ favorites });
                }
            },

            checkIsFavorite: (playerId) => {
                return get().favorites.some(f => f.id === playerId);
            }
        }),
        {
            name: 'favorite-storage'
        }
    )
);
