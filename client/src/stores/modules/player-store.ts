import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const VerificationStatus = {
    PENDING: 'pending',       // 待审核
    VERIFIED: 'verified',     // 已认证
    REJECTED: 'rejected',     // 已拒绝
    SUSPENDED: 'suspended'    // 已暂停
} as const;

export type VerificationStatus = typeof VerificationStatus[keyof typeof VerificationStatus];

export const OnlineStatus = {
    ONLINE: 'online',         // 在线接单
    BUSY: 'busy',             // 忙碌中
    OFFLINE: 'offline'        // 离线
} as const;

export type OnlineStatus = typeof OnlineStatus[keyof typeof OnlineStatus];

export const EarningsStatus = {
    FROZEN: 'frozen',         // 冻结中 (T+7)
    SETTLED: 'settled',       // 已结算
    DISPUTED: 'disputed',     // 争议中
    REFUNDED: 'refunded'      // 已退款
} as const;

export type EarningsStatus = typeof EarningsStatus[keyof typeof EarningsStatus];

// ============ Interfaces ============

export interface Player {
    id: number;
    userId: number;
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

export interface PlayerDetailProfile {
    id: number;
    userId: number;

    // 认证信息
    realName: string;
    idCard: string;           // 身份证号 (脱敏显示)
    verificationStatus: VerificationStatus;
    verifiedAt?: string;
    rejectionReason?: string;

    // 展示信息
    displayName: string;
    bio: string;
    avatar: string;
    voiceIntro?: string;      // 语音介绍 URL
    gallery: string[];        // 相册

    // 服务信息
    games: PlayerGame[];
    hourlyRateCents: number;
    serviceTimeSlots: TimeSlot[];

    // 状态
    onlineStatus: OnlineStatus;
    acceptingOrders: boolean;
    lastOnlineAt?: string;

    // 统计
    totalOrders: number;
    completedOrders: number;
    rating: number;
    reviewCount: number;

    // 收益
    totalEarningsCents: number;
    monthlyEarningsCents: number;

    createdAt: string;
    updatedAt: string;
}

export interface PlayerGame {
    gameId: number;
    gameName: string;
    gameIcon: string;
    rank: string;
    rankIcon?: string;
    certificate?: string;
}

export interface TimeSlot {
    dayOfWeek: number;        // 0-6 (周日-周六)
    startTime: string;        // "09:00"
    endTime: string;          // "22:00"
}

export interface PlayerApplyRequest {
    realName: string;
    idCard: string;
    games: GameApplication[];
    displayName?: string;
    bio?: string;
    voiceIntro?: string;
    gallery?: string[];
}

export interface GameApplication {
    gameId: number;
    rank: string;
    certificate?: string;
}

export interface ApplicationStatus {
    status: VerificationStatus;
    appliedAt: string;
    reviewedAt?: string;
    rejectionReason?: string;
    canReapply: boolean;
    reapplyAfter?: string;
}

export interface PlayerEarnings {
    totalEarningsCents: number;
    monthlyEarningsCents: number;
    weeklyEarningsCents: number;
    todayEarningsCents: number;

    wallet: {
        balanceCents: number;
        frozenCents: number;
        pendingWithdrawCents: number;
    };

    completedOrders: number;
    averageOrderCents: number;
}

export interface EarningsRecord {
    id: number;
    orderId: number;
    orderNo: string;
    orderAmountCents: number;
    commissionCents: number;
    earningsCents: number;
    status: EarningsStatus;
    settledAt?: string;
    createdAt: string;
}

export interface CommissionResult {
    orderAmountCents: number;
    baseRate: number;
    rankingDiscount: number;
    effectiveRate: number;
    commissionCents: number;
    playerEarningsCents: number;
}

export interface EarningsTrendPoint {
    date: string;
    earningsCents: number;
    orderCount: number;
}

export interface EarningsTrend {
    period: 'week' | 'month' | 'year';
    data: EarningsTrendPoint[];
    totalEarningsCents: number;
    totalOrders: number;
    averageEarningsCents: number;
}

export interface PlayerFilters {
    gameId?: number;
    minPrice?: number;
    maxPrice?: number;
    onlineOnly: boolean;
    sortBy: 'rating' | 'price' | 'orders';
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

// ============ State & Actions ============

export interface PlayerState {
    // 用户端数据
    players: Player[];
    featuredPlayers: Player[];
    currentPlayer: Player | null;
    filters: PlayerFilters;
    pagination: {
        page: number;
        pageSize: number;
        hasMore: boolean;
        total: number;
    };

    // 陪玩师端数据
    myProfile: PlayerDetailProfile | null;
    applicationStatus: ApplicationStatus | null;
    earnings: PlayerEarnings | null;
    earningsRecords: EarningsRecord[];
    earningsTrend: EarningsTrend | null;

    // 状态
    loading: boolean;
    error: string | null;
}

export interface PlayerActions {
    // 用户端 Actions
    fetchPlayers: (refresh?: boolean) => Promise<void>;
    fetchPlayerById: (id: number) => Promise<void>;
    fetchFeaturedPlayers: () => Promise<void>;
    setFilters: (filters: Partial<PlayerFilters>) => void;
    resetFilters: () => void;
    setPage: (page: number) => void;

    // 陪玩师申请 Actions
    applyToBePlayer: (data: PlayerApplyRequest) => Promise<void>;
    fetchApplicationStatus: () => Promise<void>;
    reapply: (data: PlayerApplyRequest) => Promise<void>;

    // 陪玩师资料 Actions
    fetchMyProfile: () => Promise<void>;
    updateProfile: (data: Partial<PlayerDetailProfile>) => Promise<void>;
    updateOnlineStatus: (status: OnlineStatus) => Promise<void>;
    sendHeartbeat: () => Promise<void>;

    // 收益 Actions
    fetchEarnings: () => Promise<void>;
    fetchEarningsRecords: (page?: number) => Promise<void>;
    fetchEarningsTrend: (period?: 'week' | 'month' | 'year') => Promise<void>;
    calculateCommission: (orderAmountCents: number) => Promise<CommissionResult>;
}

const INITIAL_FILTERS: PlayerFilters = {
    onlineOnly: false,
    sortBy: 'rating',
};

// Debounce helper
function debounce<T extends (...args: Parameters<T>) => void>(func: T, wait: number): T {
    let timeout: ReturnType<typeof setTimeout>;
    return function (this: unknown, ...args: Parameters<T>) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    } as T;
}

// PERFORMANCE FIX: Create debounced fetch function ONCE outside the store
// to prevent memory leak from creating new debounce instances on every filter change
let debouncedFetchPlayers: (() => void) | null = null;

const getDebouncedFetchPlayers = (fetchFn: () => void) => {
    if (!debouncedFetchPlayers) {
        debouncedFetchPlayers = debounce(fetchFn, 500);
    }
    return debouncedFetchPlayers;
};

// ============ Store ============

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
        myProfile: null,
        applicationStatus: null,
        earnings: null,
        earningsRecords: [],
        earningsTrend: null,
        loading: false,
        error: null,

        // ========== 用户端 Actions ==========

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
                set({ loading: false, error: errorMessage });
            }
        },

        fetchPlayerById: async (id: number) => {
            set({ loading: true, error: null });
            try {
                const data = await http.get<Player>(`/public/players/${id}`);
                set({ currentPlayer: data, loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch player details' });
            }
        },

        fetchFeaturedPlayers: async () => {
            try {
                const data = await http.get<Player[]>('/public/players/featured');
                set({ featuredPlayers: data });
            } catch {
                console.warn("Failed to fetch featured players");
            }
        },

        setFilters: (newFilters) => {
            set((state) => ({
                filters: { ...state.filters, ...newFilters },
                pagination: { ...state.pagination, page: 1 },
                players: []
            }));

            // PERFORMANCE FIX: Use memoized debounce function to prevent memory leak
            const debouncedFetch = getDebouncedFetchPlayers(() => {
                get().fetchPlayers(true);
            });
            debouncedFetch();
        },

        resetFilters: () => {
            set({ filters: INITIAL_FILTERS });
            get().fetchPlayers(true);
        },

        setPage: (page) => {
            set((state) => ({ pagination: { ...state.pagination, page } }));
            get().fetchPlayers(false);
        },

        // ========== 陪玩师申请 Actions ==========

        applyToBePlayer: async (data) => {
            set({ loading: true, error: null });
            try {
                await http.post('/player/apply', data);
                await get().fetchApplicationStatus();
                set({ loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to submit application' });
                throw err;
            }
        },

        fetchApplicationStatus: async () => {
            try {
                const data = await http.get<ApplicationStatus>('/player/application');
                set({ applicationStatus: data });
            } catch (err) {
                if (err && typeof err === 'object' && 'status' in err && err.status !== 404) {
                    console.error('Failed to fetch application status:', err);
                }
            }
        },

        reapply: async (data) => {
            set({ loading: true, error: null });
            try {
                await http.post('/player/reapply', data);
                await get().fetchApplicationStatus();
                set({ loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to reapply' });
                throw err;
            }
        },

        // ========== 陪玩师资料 Actions ==========

        fetchMyProfile: async () => {
            set({ loading: true, error: null });
            try {
                const data = await http.get<PlayerDetailProfile>('/player/profile');
                set({ myProfile: data, loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch profile' });
            }
        },

        updateProfile: async (data) => {
            set({ loading: true, error: null });
            try {
                const updated = await http.put<PlayerDetailProfile>('/player/profile', data);
                set({ myProfile: updated, loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to update profile' });
                throw err;
            }
        },

        updateOnlineStatus: async (status) => {
            try {
                await http.put('/player/status', { onlineStatus: status });
                set((state) => ({
                    myProfile: state.myProfile
                        ? { ...state.myProfile, onlineStatus: status }
                        : null
                }));
            } catch (err) {
                console.error('Failed to update online status:', err);
                throw err;
            }
        },

        sendHeartbeat: async () => {
            try {
                await http.post('/player/heartbeat');
            } catch (err) {
                console.error('Heartbeat failed:', err);
            }
        },

        // ========== 收益 Actions ==========

        fetchEarnings: async () => {
            set({ loading: true, error: null });
            try {
                const data = await http.get<PlayerEarnings>('/player/earnings');
                set({ earnings: data, loading: false });
            } catch (err) {
                set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch earnings' });
            }
        },

        fetchEarningsRecords: async (page = 1) => {
            try {
                const data = await http.get<{ items: EarningsRecord[] }>('/player/earnings/records', {
                    params: { page, pageSize: 20 }
                });
                set({ earningsRecords: data.items || [] });
            } catch (err) {
                console.error('Failed to fetch earnings records:', err);
            }
        },

        fetchEarningsTrend: async (period = 'month') => {
            try {
                const data = await http.get<EarningsTrend>('/player/earnings/trend', {
                    params: { period }
                });
                set({ earningsTrend: data });
            } catch (err) {
                console.error('Failed to fetch earnings trend:', err);
            }
        },

        calculateCommission: async (orderAmountCents) => {
            // SECURITY: Commission calculation MUST be server-authoritative
            // Never trust client-side calculations for financial data
            try {
                const result = await http.post<CommissionResult>('/player/commission/calculate', {
                    orderAmountCents
                });
                return result;
            } catch (err) {
                console.error('Failed to calculate commission:', err);
                // Return a safe fallback that shows no earnings (server will validate anyway)
                return {
                    orderAmountCents,
                    baseRate: 0,
                    rankingDiscount: 0,
                    effectiveRate: 0,
                    commissionCents: orderAmountCents,
                    playerEarningsCents: 0
                };
            }
        }
    }))
);
