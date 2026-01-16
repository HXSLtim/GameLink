/**
 * Favorite Store Tests
 * Tests for favorite player management with optimistic updates
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useFavoriteStore } from '../favorite-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

describe('Favorite Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useFavoriteStore.setState({
            favorites: [],
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useFavoriteStore.getState();

            expect(state.favorites).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchFavorites', () => {
        it('should fetch favorites successfully', async () => {
            const mockFavorites = {
                items: [
                    {
                        id: 1,
                        nickname: 'Player1',
                        avatar: 'avatar1.jpg',
                        gameId: 1,
                        rating: 4.5,
                    },
                    {
                        id: 2,
                        nickname: 'Player2',
                        avatar: 'avatar2.jpg',
                        gameId: 2,
                        rating: 4.8,
                    },
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockFavorites);

            await useFavoriteStore.getState().fetchFavorites();

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(2);
            expect(state.favorites[0].nickname).toBe('Player1');
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should handle empty response', async () => {
            mockHttp.get.mockResolvedValueOnce({ items: null });

            await useFavoriteStore.getState().fetchFavorites();

            const state = useFavoriteStore.getState();
            expect(state.favorites).toEqual([]);
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useFavoriteStore.getState().fetchFavorites();

            const state = useFavoriteStore.getState();
            expect(state.favorites).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });
    });

    describe('toggleFavorite', () => {
        it('should add player to favorites (optimistic update)', async () => {
            mockHttp.post.mockResolvedValueOnce({});

            const player = {
                id: 1,
                nickname: 'NewPlayer',
                avatar: 'avatar.jpg',
            };

            await useFavoriteStore.getState().toggleFavorite(player);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(1);
            expect(state.favorites[0].id).toBe(1);
            expect(mockHttp.post).toHaveBeenCalledWith('/user/favorites', { playerId: 1 });
        });

        it('should remove player from favorites (optimistic update)', async () => {
            useFavoriteStore.setState({
                favorites: [
                    { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' },
                ],
            });

            mockHttp.delete.mockResolvedValueOnce({});

            const player = { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' };

            await useFavoriteStore.getState().toggleFavorite(player);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(0);
            expect(mockHttp.delete).toHaveBeenCalledWith('/user/favorites/1');
        });

        it('should rollback on add error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('API error'));

            const player = {
                id: 1,
                nickname: 'NewPlayer',
                avatar: 'avatar.jpg',
            };

            await useFavoriteStore.getState().toggleFavorite(player);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(0);
            expect(state.error).toBe('API error');
        });

        it('should rollback on remove error', async () => {
            const originalPlayer = { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' };
            useFavoriteStore.setState({
                favorites: [originalPlayer],
            });

            mockHttp.delete.mockRejectedValueOnce(new Error('API error'));

            await useFavoriteStore.getState().toggleFavorite(originalPlayer);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(1);
            expect(state.favorites[0].id).toBe(1);
            expect(state.error).toBe('API error');
        });
    });

    describe('removeFavorite', () => {
        it('should remove favorite (optimistic update)', async () => {
            useFavoriteStore.setState({
                favorites: [
                    { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' },
                    { id: 2, nickname: 'Player2', avatar: 'avatar2.jpg' },
                ],
            });

            mockHttp.delete.mockResolvedValueOnce({});

            await useFavoriteStore.getState().removeFavorite(1);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(1);
            expect(state.favorites[0].id).toBe(2);
            expect(mockHttp.delete).toHaveBeenCalledWith('/user/favorites/1');
        });

        it('should rollback on error', async () => {
            const players = [
                { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' },
                { id: 2, nickname: 'Player2', avatar: 'avatar2.jpg' },
            ];
            useFavoriteStore.setState({ favorites: players });

            mockHttp.delete.mockRejectedValueOnce(new Error('API error'));

            await useFavoriteStore.getState().removeFavorite(1);

            const state = useFavoriteStore.getState();
            expect(state.favorites).toHaveLength(2);
            expect(state.error).toBe('API error');
        });
    });

    describe('checkIsFavorite', () => {
        it('should return true for favorited player', () => {
            useFavoriteStore.setState({
                favorites: [
                    { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' },
                ],
            });

            expect(useFavoriteStore.getState().checkIsFavorite(1)).toBe(true);
        });

        it('should return false for non-favorited player', () => {
            useFavoriteStore.setState({
                favorites: [
                    { id: 1, nickname: 'Player1', avatar: 'avatar1.jpg' },
                ],
            });

            expect(useFavoriteStore.getState().checkIsFavorite(2)).toBe(false);
        });

        it('should return false when favorites is empty', () => {
            expect(useFavoriteStore.getState().checkIsFavorite(1)).toBe(false);
        });
    });
});
