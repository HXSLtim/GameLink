import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';
import { getErrorMessage } from '@/types/api';
import type { LoginResponse, RegisterResponse, RefreshResponse, MeResponse } from '@/types/api';

// --- Types ---
export interface User {
    id: string; // or number, aligning with existing logic
    username: string;
    avatar: string;
    email?: string;
    nickname?: string;
}

export interface PlayerProfile {
    id: number;
    level: number;
    gameId: number;
    rating: number;
}

export interface AuthState {
    // User Info
    user: User | null;
    token: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;

    // Role Info
    role: 'guest' | 'user' | 'player' | 'admin';
    permissions: string[];

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
                        user,
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
                        user,
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
                const { refreshToken } = get();
                if (!refreshToken) {
                    throw new Error('No refresh token available');
                }

                try {
                    // Use public endpoint which doesn't require auth header
                    const data = await http.post<RefreshResponse>('/public/auth/refresh', { refreshToken });
                    const { token: newToken, refreshToken: newRefreshToken, user } = data;


                    set({
                        token: newToken,
                        refreshToken: newRefreshToken || refreshToken, // Use new one if provided
                        user: user || get().user,
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
                }
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
                // TODO: Implement player mode switching logic
                // Currently a no-op placeholder for future implementation
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
                    set({
                        user: userData as User,
                        isAuthenticated: true,
                        role: (userData.role as AuthState['role']) || 'user'
                    });
                } catch {
                    get().logout();
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
                playerProfile: state.playerProfile,
            }),
            onRehydrateStorage: () => (state) => {
                if (state) {
                    if (state.refreshToken && !state.token) {
                        state.refresh().catch(() => {
                        });
                    }
                }
            }
        }
    )
);
