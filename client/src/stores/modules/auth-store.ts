import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';
import { getErrorMessage } from '@/types/api';
import type { LoginResponse, RegisterResponse, RefreshResponse, MeResponse } from '@/types/api';

// Promise lock to prevent concurrent refresh calls
let refreshPromise: Promise<void> | null = null;

// --- Types ---
export interface User {
    id: number; // Changed to number to match backend uint64
    username: string;
    avatar: string;
    email?: string;
    name: string; // Display name (mapped from backend 'nickname')
}

// Helper to map API user response to User type
function mapApiUserToUser(apiUser: {
    id: string;
    username: string;
    avatar: string;
    email?: string;
    nickname?: string;
}): User {
    return {
        id: parseInt(apiUser.id, 10) || 0,
        username: apiUser.username,
        avatar: apiUser.avatar,
        email: apiUser.email,
        name: apiUser.nickname || apiUser.username, // Fallback to username if no nickname
    };
}

export interface PlayerProfile {
    id: number;
    level: number;
    gameId: number;
    rating: number;
}

export type ViewMode = 'user' | 'player';

export interface AuthState {
    // User Info
    user: User | null;
    token: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;

    // Role Info
    role: 'guest' | 'user' | 'player' | 'admin';
    permissions: string[];

    // View Mode (用户视图 or 陪玩视图)
    viewMode: ViewMode;
    isPlayer: boolean; // 用户是否已认证为陪玩

    // Player Profile (if applicable)
    playerProfile: PlayerProfile | null;

    // Status
    loading: boolean;
    error: string | null;
}

export interface RegisterPayload {
    phone?: string;
    email?: string;
    password: string;
    name: string;
}

export interface AuthActions {
    login: (credentials: { username: string; password: string }) => Promise<void>;
    register: (credentials: RegisterPayload) => Promise<void>;
    logout: () => Promise<void>;
    refresh: () => Promise<void>;
    updateProfile: (data: Partial<User>) => Promise<void>;
    changePassword: (data: { oldPassword?: string; newPassword: string }) => Promise<void>;
    switchToPlayerMode: () => void;
    switchToUserMode: () => void;
    setIsPlayer: (isPlayer: boolean) => void;
    checkAuth: () => Promise<void>;
}

// --- Store ---
export const useAuthStore = create<AuthState & AuthActions>()(
    persist(
        (set, get) => ({
            // Initial State
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
            role: 'guest',
            permissions: [],
            viewMode: 'user',
            isPlayer: false,
            playerProfile: null,
            loading: false,
            error: null,

            // Actions
            login: async (credentials) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.post<LoginResponse>('/auth/login', {
                        username: credentials.username,
                        password: credentials.password
                    });
                    const { token, refreshToken, user, role, permissions } = data;

                    set({
                        token,
                        refreshToken,
                        user: mapApiUserToUser(user),
                        role: (role as AuthState['role']) || 'user',
                        permissions: permissions || [],
                        isAuthenticated: true,
                        loading: false,
                        error: null,
                    });
                } catch (err: unknown) {
                    console.error("Login failed", err);
                    set({
                        loading: false,
                        error: getErrorMessage(err),
                    });
                    throw err;
                }
            },

            register: async (credentials: RegisterPayload) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.post<RegisterResponse>('/auth/register', credentials);
                    // Automatically login after register
                    const { token, refreshToken, user, role, permissions } = data;

                    set({
                        token,
                        refreshToken,
                        user: mapApiUserToUser(user),
                        role: (role as AuthState['role']) || 'user',
                        permissions: permissions || [],
                        isAuthenticated: true,
                        loading: false,
                        error: null,
                    });
                } catch (err: unknown) {
                    console.error("Registration failed", err);
                    set({
                        loading: false,
                        error: getErrorMessage(err),
                    });
                    throw err;
                }
            },

            logout: async () => {
                try {
                    // Best effort logout
                    await http.post('/auth/logout');
                } catch { /* ignore */ }

                set({
                    user: null,
                    token: null,
                    refreshToken: null,
                    isAuthenticated: false,
                    role: 'guest',
                    permissions: [],
                    playerProfile: null,
                });
            },

            refresh: async () => {
                // Use Promise lock to prevent concurrent refresh calls (race condition fix)
                if (refreshPromise) {
                    return refreshPromise;
                }

                const { refreshToken } = get();
                if (!refreshToken) {
                    throw new Error('No refresh token available');
                }

                refreshPromise = (async () => {
                    try {
                        // Use public endpoint which doesn't require auth header
                        const data = await http.post<RefreshResponse>('/public/auth/refresh', { refreshToken });
                        const { token: newToken, refreshToken: newRefreshToken, user } = data;

                        set({
                            token: newToken,
                            refreshToken: newRefreshToken || refreshToken, // Use new one if provided
                            user: user ? mapApiUserToUser(user) : get().user,
                            loading: false,
                            error: null,
                            isAuthenticated: true
                        });
                    } catch (err: unknown) {
                        // Clear auth state on refresh failure
                        set({
                            token: null,
                            refreshToken: null,
                            user: null,
                            isAuthenticated: false,
                            role: 'guest'
                        });
                        throw err;
                    } finally {
                        refreshPromise = null;
                    }
                })();

                return refreshPromise;
            },

            updateProfile: async (data) => {
                set({ loading: true, error: null });
                try {
                    const updatedUser = await http.put<User>('/user/profile', data);
                    set((state) => ({
                        user: state.user ? { ...state.user, ...updatedUser } : updatedUser,
                        loading: false
                    }));
                } catch (err: unknown) {
                    set({ loading: false, error: getErrorMessage(err) });
                    throw err;
                }
            },

            changePassword: async (passwords) => {
                set({ loading: true, error: null });
                try {
                    await http.put('/user/password', passwords);
                    set({ loading: false });
                } catch (err: unknown) {
                    set({ loading: false, error: getErrorMessage(err) });
                    throw err;
                }
            },

            switchToPlayerMode: () => {
                const { isPlayer } = get();
                if (isPlayer) {
                    set({ viewMode: 'player' });
                }
            },

            switchToUserMode: () => {
                set({ viewMode: 'user' });
            },

            setIsPlayer: (isPlayer: boolean) => {
                set({ isPlayer });
                // 如果不再是陪玩，自动切换回用户视图
                if (!isPlayer) {
                    set({ viewMode: 'user' });
                }
            },

            checkAuth: async () => {
                const { token, refreshToken, refresh } = get();

                // If we have no token but have a refresh token (e.g. page reload), try to refresh
                if (!token && refreshToken) {
                    try {
                        await refresh();
                        return;
                    } catch {
                        // Refresh failed, proceed to logout/guest state
                    }
                }

                if (!get().token) {
                    set({ isAuthenticated: false, role: 'guest' });
                    return;
                }

                try {
                    const data = await http.get<MeResponse>('/auth/me');
                    const userData = data.user || data;
                    // Map the API response to User type
                    const mappedUser: User = {
                        id: parseInt(userData.id || '0', 10),
                        username: userData.username || '',
                        avatar: userData.avatar || '',
                        email: userData.email,
                        name: userData.nickname || userData.username || '',
                    };
                    set({
                        user: mappedUser,
                        isAuthenticated: true,
                        role: (userData.role as AuthState['role']) || 'user'
                    });
                } catch (err: unknown) {
                    // Only logout if it's an authentication error
                    if (err && typeof err === 'object' && 'response' in err) {
                        const axiosErr = err as { response?: { status?: number } };
                        if (axiosErr.response?.status === 401) {
                            get().logout();
                        } else {
                            console.warn('CheckAuth failed but not 401, keeping session:', err);
                        }
                    } else {
                        console.warn('CheckAuth failed with non-axios error, keeping session:', err);
                    }
                }
            }
        }),
        {
            name: 'auth-storage',
            partialize: (state) => ({
                // Always persist auth state for session restoration
                token: state.token,
                refreshToken: state.refreshToken,
                user: state.user,
                role: state.role,
                viewMode: state.viewMode,
                isPlayer: state.isPlayer,
                playerProfile: state.playerProfile,
            }),
            onRehydrateStorage: () => (state) => {
                if (state) {
                    // If token exists but is expired, or there's only a refresh token, try to refresh
                    if ((state.token && !state.isAuthenticated) || (state.refreshToken && !state.token)) {
                        state.refresh().catch(() => {
                            // Refresh failed, user will be logged out
                            // Use useAuthStore.setState() since we're outside the store creator
                            useAuthStore.setState({
                                user: null,
                                token: null,
                                refreshToken: null,
                                isAuthenticated: false,
                                role: 'guest'
                            });
                        });
                    }
                }
            }
        }
    )
);
