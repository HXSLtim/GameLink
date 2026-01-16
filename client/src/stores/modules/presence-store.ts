import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const PresenceStatus = {
    ONLINE: 'online',           // 在线空闲
    ACCEPTING: 'accepting',     // 接单中
    IN_GAME: 'in_game',         // 游戏中
    MATCHING: 'matching',       // 匹配中
    RESTING: 'resting',         // 休息中
    OFFLINE: 'offline',         // 离线
    INVISIBLE: 'invisible',     // 隐身
} as const;

export type PresenceStatus = typeof PresenceStatus[keyof typeof PresenceStatus];

// ============ Interfaces ============

export interface PlayerPresence {
    id: number;
    playerId: number;
    status: PresenceStatus;
    currentGameId?: number;
    currentGameName?: string;
    customStatus?: string;
    currentOrderId?: number;
    currentRoomId?: number;
    lastHeartbeatAt: string;
    deviceType?: string;
    // Relations
    player?: {
        id: number;
        nickname: string;
        avatar: string;
    };
}

export interface UpdatePresenceRequest {
    status?: PresenceStatus;
    currentGameId?: number;
    currentGameName?: string;
    customStatus?: string;
    deviceType?: string;
}

export interface PresenceListResponse {
    items: PlayerPresence[];
    total: number;
}

// ============ Helper Functions ============

export function getStatusDisplay(status: PresenceStatus): string {
    switch (status) {
        case PresenceStatus.ONLINE:
            return '在线';
        case PresenceStatus.ACCEPTING:
            return '接单中';
        case PresenceStatus.IN_GAME:
            return '游戏中';
        case PresenceStatus.MATCHING:
            return '匹配中';
        case PresenceStatus.RESTING:
            return '休息中';
        case PresenceStatus.OFFLINE:
            return '离线';
        case PresenceStatus.INVISIBLE:
            return '隐身';
        default:
            return '未知';
    }
}

export function getStatusColor(status: PresenceStatus): string {
    switch (status) {
        case PresenceStatus.ONLINE:
            return '#22c55e'; // green
        case PresenceStatus.ACCEPTING:
            return '#3b82f6'; // blue
        case PresenceStatus.IN_GAME:
            return '#a855f7'; // purple
        case PresenceStatus.MATCHING:
            return '#f59e0b'; // amber
        case PresenceStatus.RESTING:
            return '#6b7280'; // gray
        case PresenceStatus.OFFLINE:
            return '#9ca3af'; // light gray
        case PresenceStatus.INVISIBLE:
            return '#9ca3af'; // light gray
        default:
            return '#9ca3af';
    }
}

export function isOnlineStatus(status: PresenceStatus): boolean {
    return status !== PresenceStatus.OFFLINE && status !== PresenceStatus.INVISIBLE;
}

export function isAvailableStatus(status: PresenceStatus): boolean {
    return status === PresenceStatus.ONLINE || status === PresenceStatus.ACCEPTING;
}

// ============ State & Actions ============

export interface PresenceState {
    // My presence
    myPresence: PlayerPresence | null;

    // Watched presences (other players)
    watchedPresences: Map<number, PlayerPresence>;

    // Online players list
    onlinePlayers: PlayerPresence[];
    onlineCount: number;

    // Pagination for online list
    pagination: {
        page: number;
        pageSize: number;
        total: number;
        hasMore: boolean;
    };

    // Heartbeat interval
    heartbeatInterval: ReturnType<typeof setInterval> | null;

    // Status
    loading: boolean;
    error: string | null;
}

export interface PresenceActions {
    // My presence actions
    fetchMyPresence: () => Promise<void>;
    updatePresence: (data: UpdatePresenceRequest) => Promise<void>;
    setStatus: (status: PresenceStatus) => Promise<void>;
    setCustomStatus: (customStatus: string) => Promise<void>;

    // Quick status setters
    goOnline: () => Promise<void>;
    goOffline: () => Promise<void>;
    setAccepting: () => Promise<void>;
    setResting: () => Promise<void>;
    setInvisible: () => Promise<void>;
    setInGame: (gameId?: number, gameName?: string) => Promise<void>;
    setMatching: () => Promise<void>;

    // Heartbeat
    startHeartbeat: () => void;
    stopHeartbeat: () => void;
    sendHeartbeat: () => Promise<void>;

    // Watch other players
    watchPlayer: (playerId: number) => Promise<void>;
    unwatchPlayer: (playerId: number) => void;
    fetchPlayerPresence: (playerId: number) => Promise<PlayerPresence | null>;
    fetchPlayersPresence: (playerIds: number[]) => Promise<void>;

    // Online players list
    fetchOnlinePlayers: (refresh?: boolean) => Promise<void>;
    fetchOnlineCount: () => Promise<void>;

    // WebSocket handlers
    handlePresenceUpdate: (data: PlayerPresence) => void;

    // Cleanup
    reset: () => void;
}

const HEARTBEAT_INTERVAL = 30000; // 30 seconds

// ============ Store ============

export const usePresenceStore = create<PresenceState & PresenceActions>()(
    subscribeWithSelector((set, get) => ({
        // Initial State
        myPresence: null,
        watchedPresences: new Map(),
        onlinePlayers: [],
        onlineCount: 0,
        pagination: {
            page: 1,
            pageSize: 20,
            total: 0,
            hasMore: true,
        },
        heartbeatInterval: null,
        loading: false,
        error: null,

        // ========== My Presence Actions ==========

        fetchMyPresence: async () => {
            set({ loading: true, error: null });
            try {
                const data = await http.get<PlayerPresence>('/user/presence');
                set({ myPresence: data, loading: false });
            } catch (err) {
                set({
                    loading: false,
                    error: err instanceof Error ? err.message : 'Failed to fetch presence',
                });
            }
        },

        updatePresence: async (data) => {
            set({ loading: true, error: null });
            try {
                const updated = await http.put<PlayerPresence>('/user/presence', data);
                set({ myPresence: updated, loading: false });
            } catch (err) {
                set({
                    loading: false,
                    error: err instanceof Error ? err.message : 'Failed to update presence',
                });
                throw err;
            }
        },

        setStatus: async (status) => {
            try {
                await http.put('/user/presence/status', { status });
                set((state) => ({
                    myPresence: state.myPresence
                        ? { ...state.myPresence, status }
                        : null,
                }));
            } catch (err) {
                console.error('Failed to set status:', err);
                throw err;
            }
        },

        setCustomStatus: async (customStatus) => {
            await get().updatePresence({ customStatus });
        },

        // ========== Quick Status Setters ==========

        goOnline: async () => {
            await get().setStatus(PresenceStatus.ONLINE);
        },

        goOffline: async () => {
            await get().setStatus(PresenceStatus.OFFLINE);
        },

        setAccepting: async () => {
            await get().setStatus(PresenceStatus.ACCEPTING);
        },

        setResting: async () => {
            await get().setStatus(PresenceStatus.RESTING);
        },

        setInvisible: async () => {
            await get().setStatus(PresenceStatus.INVISIBLE);
        },

        setInGame: async (gameId, gameName) => {
            await get().updatePresence({
                status: PresenceStatus.IN_GAME,
                currentGameId: gameId,
                currentGameName: gameName,
            });
        },

        setMatching: async () => {
            await get().setStatus(PresenceStatus.MATCHING);
        },

        // ========== Heartbeat ==========

        startHeartbeat: () => {
            const { heartbeatInterval } = get();
            if (heartbeatInterval) {
                clearInterval(heartbeatInterval);
            }

            // Send initial heartbeat
            get().sendHeartbeat();

            // Set up interval
            const interval = setInterval(() => {
                get().sendHeartbeat();
            }, HEARTBEAT_INTERVAL);

            set({ heartbeatInterval: interval });
        },

        stopHeartbeat: () => {
            const { heartbeatInterval } = get();
            if (heartbeatInterval) {
                clearInterval(heartbeatInterval);
                set({ heartbeatInterval: null });
            }
        },

        sendHeartbeat: async () => {
            try {
                await http.post('/user/presence/heartbeat');
            } catch (err) {
                console.error('Heartbeat failed:', err);
            }
        },

        // ========== Watch Other Players ==========

        watchPlayer: async (playerId) => {
            const presence = await get().fetchPlayerPresence(playerId);
            if (presence) {
                set((state) => {
                    const newMap = new Map(state.watchedPresences);
                    newMap.set(playerId, presence);
                    return { watchedPresences: newMap };
                });
            }
        },

        unwatchPlayer: (playerId) => {
            set((state) => {
                const newMap = new Map(state.watchedPresences);
                newMap.delete(playerId);
                return { watchedPresences: newMap };
            });
        },

        fetchPlayerPresence: async (playerId) => {
            try {
                const data = await http.get<PlayerPresence>(`/user/players/${playerId}/presence`);
                return data;
            } catch (err) {
                console.error('Failed to fetch player presence:', err);
                return null;
            }
        },

        fetchPlayersPresence: async (playerIds) => {
            if (playerIds.length === 0) return;

            try {
                const data = await http.post<PlayerPresence[]>('/user/players/presence', {
                    playerIds,
                });

                set((state) => {
                    const newMap = new Map(state.watchedPresences);
                    for (const presence of data) {
                        newMap.set(presence.playerId, presence);
                    }
                    return { watchedPresences: newMap };
                });
            } catch (err) {
                console.error('Failed to fetch players presence:', err);
            }
        },

        // ========== Online Players List ==========

        fetchOnlinePlayers: async (refresh = false) => {
            set({ loading: true, error: null });
            const { pagination, onlinePlayers } = get();
            const currentPage = refresh ? 1 : pagination.page;

            try {
                const data = await http.get<PresenceListResponse>('/user/players/online', {
                    params: {
                        page: currentPage,
                        pageSize: pagination.pageSize,
                    },
                });

                const newPlayers = data.items || [];
                const total = data.total || 0;

                set({
                    onlinePlayers: refresh ? newPlayers : [...onlinePlayers, ...newPlayers],
                    pagination: {
                        ...pagination,
                        page: currentPage,
                        total,
                        hasMore: onlinePlayers.length + newPlayers.length < total,
                    },
                    loading: false,
                });
            } catch (err) {
                set({
                    loading: false,
                    error: err instanceof Error ? err.message : 'Failed to fetch online players',
                });
            }
        },

        fetchOnlineCount: async () => {
            try {
                const data = await http.get<{ count: number }>('/user/players/online/count');
                set({ onlineCount: data.count });
            } catch (err) {
                console.error('Failed to fetch online count:', err);
            }
        },

        // ========== WebSocket Handlers ==========

        handlePresenceUpdate: (data) => {
            const { myPresence, watchedPresences } = get();

            // Update my presence if it's mine
            if (myPresence && data.playerId === myPresence.playerId) {
                set({ myPresence: data });
            }

            // Update watched presence if we're watching this player
            if (watchedPresences.has(data.playerId)) {
                set((state) => {
                    const newMap = new Map(state.watchedPresences);
                    newMap.set(data.playerId, data);
                    return { watchedPresences: newMap };
                });
            }
        },

        // ========== Cleanup ==========

        reset: () => {
            const { heartbeatInterval } = get();
            if (heartbeatInterval) {
                clearInterval(heartbeatInterval);
            }

            set({
                myPresence: null,
                watchedPresences: new Map(),
                onlinePlayers: [],
                onlineCount: 0,
                pagination: {
                    page: 1,
                    pageSize: 20,
                    total: 0,
                    hasMore: true,
                },
                heartbeatInterval: null,
                loading: false,
                error: null,
            });
        },
    }))
);
