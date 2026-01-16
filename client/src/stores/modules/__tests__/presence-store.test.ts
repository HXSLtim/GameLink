import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
    usePresenceStore,
    PlayerPresence,
    PresenceStatus,
    getStatusDisplay,
    getStatusColor,
    isOnlineStatus,
    isAvailableStatus,
} from '../presence-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockPresence: PlayerPresence = {
    id: 1,
    playerId: 100,
    status: 'online',
    currentGameId: 1,
    currentGameName: 'League of Legends',
    customStatus: 'Ready to play!',
    currentOrderId: undefined,
    currentRoomId: undefined,
    lastHeartbeatAt: '2024-01-01T00:00:00Z',
    deviceType: 'web',
    player: {
        id: 100,
        nickname: 'TestPlayer',
        avatar: 'https://example.com/avatar.png',
    },
};

describe('presence-store', () => {
    beforeEach(() => {
        usePresenceStore.getState().reset();
        vi.clearAllMocks();
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    // ========== Helper Functions Tests ==========

    describe('getStatusDisplay', () => {
        it('should return correct display text for each status', () => {
            expect(getStatusDisplay(PresenceStatus.ONLINE)).toBe('在线');
            expect(getStatusDisplay(PresenceStatus.ACCEPTING)).toBe('接单中');
            expect(getStatusDisplay(PresenceStatus.IN_GAME)).toBe('游戏中');
            expect(getStatusDisplay(PresenceStatus.MATCHING)).toBe('匹配中');
            expect(getStatusDisplay(PresenceStatus.RESTING)).toBe('休息中');
            expect(getStatusDisplay(PresenceStatus.OFFLINE)).toBe('离线');
            expect(getStatusDisplay(PresenceStatus.INVISIBLE)).toBe('隐身');
        });

        it('should return "未知" for unknown status', () => {
            expect(getStatusDisplay('unknown' as PresenceStatus)).toBe('未知');
        });
    });

    describe('getStatusColor', () => {
        it('should return correct color for each status', () => {
            expect(getStatusColor(PresenceStatus.ONLINE)).toBe('#22c55e');
            expect(getStatusColor(PresenceStatus.ACCEPTING)).toBe('#3b82f6');
            expect(getStatusColor(PresenceStatus.IN_GAME)).toBe('#a855f7');
            expect(getStatusColor(PresenceStatus.MATCHING)).toBe('#f59e0b');
            expect(getStatusColor(PresenceStatus.RESTING)).toBe('#6b7280');
            expect(getStatusColor(PresenceStatus.OFFLINE)).toBe('#9ca3af');
            expect(getStatusColor(PresenceStatus.INVISIBLE)).toBe('#9ca3af');
        });
    });

    describe('isOnlineStatus', () => {
        it('should return true for online statuses', () => {
            expect(isOnlineStatus(PresenceStatus.ONLINE)).toBe(true);
            expect(isOnlineStatus(PresenceStatus.ACCEPTING)).toBe(true);
            expect(isOnlineStatus(PresenceStatus.IN_GAME)).toBe(true);
            expect(isOnlineStatus(PresenceStatus.MATCHING)).toBe(true);
            expect(isOnlineStatus(PresenceStatus.RESTING)).toBe(true);
        });

        it('should return false for offline/invisible statuses', () => {
            expect(isOnlineStatus(PresenceStatus.OFFLINE)).toBe(false);
            expect(isOnlineStatus(PresenceStatus.INVISIBLE)).toBe(false);
        });
    });

    describe('isAvailableStatus', () => {
        it('should return true for available statuses', () => {
            expect(isAvailableStatus(PresenceStatus.ONLINE)).toBe(true);
            expect(isAvailableStatus(PresenceStatus.ACCEPTING)).toBe(true);
        });

        it('should return false for unavailable statuses', () => {
            expect(isAvailableStatus(PresenceStatus.IN_GAME)).toBe(false);
            expect(isAvailableStatus(PresenceStatus.MATCHING)).toBe(false);
            expect(isAvailableStatus(PresenceStatus.RESTING)).toBe(false);
            expect(isAvailableStatus(PresenceStatus.OFFLINE)).toBe(false);
            expect(isAvailableStatus(PresenceStatus.INVISIBLE)).toBe(false);
        });
    });

    // ========== Store Actions Tests ==========

    describe('fetchMyPresence', () => {
        it('should fetch my presence successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockPresence);

            await usePresenceStore.getState().fetchMyPresence();

            expect(http.get).toHaveBeenCalledWith('/user/presence');
            expect(usePresenceStore.getState().myPresence).toEqual(mockPresence);
            expect(usePresenceStore.getState().loading).toBe(false);
        });

        it('should handle fetch error', async () => {
            const error = new Error('Network error');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await usePresenceStore.getState().fetchMyPresence();

            expect(usePresenceStore.getState().error).toBe('Network error');
            expect(usePresenceStore.getState().loading).toBe(false);
        });
    });

    describe('updatePresence', () => {
        it('should update presence successfully', async () => {
            const updatedPresence = { ...mockPresence, customStatus: 'New status' };
            vi.mocked(http.put).mockResolvedValueOnce(updatedPresence);

            await usePresenceStore.getState().updatePresence({ customStatus: 'New status' });

            expect(http.put).toHaveBeenCalledWith('/user/presence', { customStatus: 'New status' });
            expect(usePresenceStore.getState().myPresence).toEqual(updatedPresence);
        });

        it('should handle update error', async () => {
            const error = new Error('Update failed');
            vi.mocked(http.put).mockRejectedValueOnce(error);

            await expect(
                usePresenceStore.getState().updatePresence({ customStatus: 'New status' })
            ).rejects.toThrow('Update failed');
            expect(usePresenceStore.getState().error).toBe('Update failed');
        });
    });

    describe('setStatus', () => {
        it('should set status successfully', async () => {
            usePresenceStore.setState({ myPresence: mockPresence });
            vi.mocked(http.put).mockResolvedValueOnce(undefined);

            await usePresenceStore.getState().setStatus(PresenceStatus.ACCEPTING);

            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'accepting' });
            expect(usePresenceStore.getState().myPresence?.status).toBe('accepting');
        });

        it('should handle set status error', async () => {
            const error = new Error('Failed to set status');
            vi.mocked(http.put).mockRejectedValueOnce(error);

            await expect(
                usePresenceStore.getState().setStatus(PresenceStatus.ACCEPTING)
            ).rejects.toThrow('Failed to set status');
        });
    });

    describe('quick status setters', () => {
        beforeEach(() => {
            usePresenceStore.setState({ myPresence: mockPresence });
            vi.mocked(http.put).mockResolvedValue(undefined);
        });

        it('goOnline should set status to online', async () => {
            await usePresenceStore.getState().goOnline();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'online' });
        });

        it('goOffline should set status to offline', async () => {
            await usePresenceStore.getState().goOffline();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'offline' });
        });

        it('setAccepting should set status to accepting', async () => {
            await usePresenceStore.getState().setAccepting();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'accepting' });
        });

        it('setResting should set status to resting', async () => {
            await usePresenceStore.getState().setResting();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'resting' });
        });

        it('setInvisible should set status to invisible', async () => {
            await usePresenceStore.getState().setInvisible();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'invisible' });
        });

        it('setMatching should set status to matching', async () => {
            await usePresenceStore.getState().setMatching();
            expect(http.put).toHaveBeenCalledWith('/user/presence/status', { status: 'matching' });
        });
    });

    describe('setInGame', () => {
        it('should set in_game status with game info', async () => {
            const updatedPresence = {
                ...mockPresence,
                status: 'in_game',
                currentGameId: 2,
                currentGameName: 'Valorant',
            };
            vi.mocked(http.put).mockResolvedValueOnce(updatedPresence);

            await usePresenceStore.getState().setInGame(2, 'Valorant');

            expect(http.put).toHaveBeenCalledWith('/user/presence', {
                status: 'in_game',
                currentGameId: 2,
                currentGameName: 'Valorant',
            });
        });
    });

    // ========== Heartbeat Tests ==========

    describe('heartbeat', () => {
        it('startHeartbeat should send initial heartbeat and set interval', async () => {
            vi.mocked(http.post).mockResolvedValue(undefined);

            usePresenceStore.getState().startHeartbeat();

            // Initial heartbeat should be sent
            expect(http.post).toHaveBeenCalledWith('/user/presence/heartbeat');
            expect(usePresenceStore.getState().heartbeatInterval).not.toBeNull();
        });

        it('startHeartbeat should clear existing interval before starting new one', async () => {
            vi.mocked(http.post).mockResolvedValue(undefined);

            usePresenceStore.getState().startHeartbeat();
            const firstInterval = usePresenceStore.getState().heartbeatInterval;

            usePresenceStore.getState().startHeartbeat();
            const secondInterval = usePresenceStore.getState().heartbeatInterval;

            expect(firstInterval).not.toBe(secondInterval);
        });

        it('stopHeartbeat should clear interval', () => {
            vi.mocked(http.post).mockResolvedValue(undefined);

            usePresenceStore.getState().startHeartbeat();
            expect(usePresenceStore.getState().heartbeatInterval).not.toBeNull();

            usePresenceStore.getState().stopHeartbeat();
            expect(usePresenceStore.getState().heartbeatInterval).toBeNull();
        });

        it('sendHeartbeat should call heartbeat endpoint', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);

            await usePresenceStore.getState().sendHeartbeat();

            expect(http.post).toHaveBeenCalledWith('/user/presence/heartbeat');
        });

        it('sendHeartbeat should handle errors silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            vi.mocked(http.post).mockRejectedValueOnce(new Error('Network error'));

            await usePresenceStore.getState().sendHeartbeat();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    // ========== Watch Player Tests ==========

    describe('watchPlayer', () => {
        it('should fetch and store player presence', async () => {
            const otherPresence = { ...mockPresence, playerId: 200 };
            vi.mocked(http.get).mockResolvedValueOnce(otherPresence);

            await usePresenceStore.getState().watchPlayer(200);

            expect(http.get).toHaveBeenCalledWith('/user/players/200/presence');
            expect(usePresenceStore.getState().watchedPresences.get(200)).toEqual(otherPresence);
        });

        it('should not store if fetch returns null', async () => {
            vi.mocked(http.get).mockRejectedValueOnce(new Error('Not found'));

            await usePresenceStore.getState().watchPlayer(200);

            expect(usePresenceStore.getState().watchedPresences.has(200)).toBe(false);
        });
    });

    describe('unwatchPlayer', () => {
        it('should remove player from watched presences', async () => {
            const otherPresence = { ...mockPresence, playerId: 200 };
            vi.mocked(http.get).mockResolvedValueOnce(otherPresence);

            await usePresenceStore.getState().watchPlayer(200);
            expect(usePresenceStore.getState().watchedPresences.has(200)).toBe(true);

            usePresenceStore.getState().unwatchPlayer(200);
            expect(usePresenceStore.getState().watchedPresences.has(200)).toBe(false);
        });
    });

    describe('fetchPlayerPresence', () => {
        it('should fetch single player presence', async () => {
            const otherPresence = { ...mockPresence, playerId: 200 };
            vi.mocked(http.get).mockResolvedValueOnce(otherPresence);

            const result = await usePresenceStore.getState().fetchPlayerPresence(200);

            expect(http.get).toHaveBeenCalledWith('/user/players/200/presence');
            expect(result).toEqual(otherPresence);
        });

        it('should return null on error', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            vi.mocked(http.get).mockRejectedValueOnce(new Error('Not found'));

            const result = await usePresenceStore.getState().fetchPlayerPresence(999);

            expect(result).toBeNull();
            consoleSpy.mockRestore();
        });
    });

    describe('fetchPlayersPresence', () => {
        it('should fetch multiple players presence', async () => {
            const presences = [
                { ...mockPresence, playerId: 200 },
                { ...mockPresence, playerId: 201 },
            ];
            vi.mocked(http.post).mockResolvedValueOnce(presences);

            await usePresenceStore.getState().fetchPlayersPresence([200, 201]);

            expect(http.post).toHaveBeenCalledWith('/user/players/presence', { playerIds: [200, 201] });
            expect(usePresenceStore.getState().watchedPresences.get(200)).toEqual(presences[0]);
            expect(usePresenceStore.getState().watchedPresences.get(201)).toEqual(presences[1]);
        });

        it('should do nothing for empty array', async () => {
            await usePresenceStore.getState().fetchPlayersPresence([]);

            expect(http.post).not.toHaveBeenCalled();
        });

        it('should handle errors silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            vi.mocked(http.post).mockRejectedValueOnce(new Error('Network error'));

            await usePresenceStore.getState().fetchPlayersPresence([200, 201]);

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    // ========== Online Players Tests ==========

    describe('fetchOnlinePlayers', () => {
        it('should fetch online players with refresh', async () => {
            const mockResponse = {
                items: [mockPresence],
                total: 1,
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await usePresenceStore.getState().fetchOnlinePlayers(true);

            expect(http.get).toHaveBeenCalledWith('/user/players/online', {
                params: { page: 1, pageSize: 20 },
            });
            expect(usePresenceStore.getState().onlinePlayers).toEqual([mockPresence]);
            expect(usePresenceStore.getState().loading).toBe(false);
        });

        it('should append players without refresh', async () => {
            usePresenceStore.setState({
                onlinePlayers: [mockPresence],
                pagination: { page: 1, pageSize: 20, total: 2, hasMore: true },
            });
            const newPresence = { ...mockPresence, playerId: 201 };
            const mockResponse = {
                items: [newPresence],
                total: 2,
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await usePresenceStore.getState().fetchOnlinePlayers(false);

            expect(usePresenceStore.getState().onlinePlayers).toHaveLength(2);
        });

        it('should handle fetch error', async () => {
            const error = new Error('Network error');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await usePresenceStore.getState().fetchOnlinePlayers(true);

            expect(usePresenceStore.getState().error).toBe('Network error');
            expect(usePresenceStore.getState().loading).toBe(false);
        });
    });

    describe('fetchOnlineCount', () => {
        it('should fetch online count', async () => {
            vi.mocked(http.get).mockResolvedValueOnce({ count: 42 });

            await usePresenceStore.getState().fetchOnlineCount();

            expect(http.get).toHaveBeenCalledWith('/user/players/online/count');
            expect(usePresenceStore.getState().onlineCount).toBe(42);
        });

        it('should handle error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            vi.mocked(http.get).mockRejectedValueOnce(new Error('Network error'));

            await usePresenceStore.getState().fetchOnlineCount();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    // ========== WebSocket Handler Tests ==========

    describe('handlePresenceUpdate', () => {
        it('should update myPresence if playerId matches', () => {
            usePresenceStore.setState({ myPresence: mockPresence });
            const updatedPresence = { ...mockPresence, status: 'accepting' as const };

            usePresenceStore.getState().handlePresenceUpdate(updatedPresence);

            expect(usePresenceStore.getState().myPresence?.status).toBe('accepting');
        });

        it('should update watched presence if player is being watched', async () => {
            const otherPresence = { ...mockPresence, playerId: 200 };
            vi.mocked(http.get).mockResolvedValueOnce(otherPresence);
            await usePresenceStore.getState().watchPlayer(200);

            const updatedPresence = { ...otherPresence, status: 'in_game' as const };
            usePresenceStore.getState().handlePresenceUpdate(updatedPresence);

            expect(usePresenceStore.getState().watchedPresences.get(200)?.status).toBe('in_game');
        });

        it('should not update if player is not watched and not self', () => {
            usePresenceStore.setState({ myPresence: mockPresence });
            const unknownPresence = { ...mockPresence, playerId: 999 };

            usePresenceStore.getState().handlePresenceUpdate(unknownPresence);

            expect(usePresenceStore.getState().watchedPresences.has(999)).toBe(false);
        });
    });

    // ========== Reset Tests ==========

    describe('reset', () => {
        it('should reset store to initial state', async () => {
            vi.mocked(http.post).mockResolvedValue(undefined);
            vi.mocked(http.get).mockResolvedValue(mockPresence);

            // Set up some state
            usePresenceStore.setState({
                myPresence: mockPresence,
                onlinePlayers: [mockPresence],
                onlineCount: 10,
                loading: true,
                error: 'Some error',
            });
            usePresenceStore.getState().startHeartbeat();
            await usePresenceStore.getState().watchPlayer(200);

            usePresenceStore.getState().reset();

            expect(usePresenceStore.getState().myPresence).toBeNull();
            expect(usePresenceStore.getState().watchedPresences.size).toBe(0);
            expect(usePresenceStore.getState().onlinePlayers).toEqual([]);
            expect(usePresenceStore.getState().onlineCount).toBe(0);
            expect(usePresenceStore.getState().heartbeatInterval).toBeNull();
            expect(usePresenceStore.getState().loading).toBe(false);
            expect(usePresenceStore.getState().error).toBeNull();
        });

        it('should clear heartbeat interval on reset', () => {
            vi.mocked(http.post).mockResolvedValue(undefined);

            usePresenceStore.getState().startHeartbeat();
            expect(usePresenceStore.getState().heartbeatInterval).not.toBeNull();

            usePresenceStore.getState().reset();
            expect(usePresenceStore.getState().heartbeatInterval).toBeNull();
        });
    });
});
