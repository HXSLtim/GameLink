import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

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

export interface AuthActions {
    login: (credentials: any) => Promise<void>;
    logout: () => Promise<void>;
    refresh: () => Promise<void>;
    updateProfile: (data: Partial<User>) => Promise<void>;
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
                    const data = await http.post<any>('/auth/login', credentials);
                    const { token, refreshToken, user, role, permissions } = data;

                    set({
                        token,
                        refreshToken,
                        user,
                        role: role || 'user',
                        permissions: permissions || [],
                        isAuthenticated: true,
                        loading: false,
                        error: null,
                    });
                } catch (err: any) {
                    console.error("Login failed", err);
                    set({
                        loading: false,
                        error: err.response?.data?.message || err.message || 'Login failed',
                    });
                    throw err;
                }
            },

            logout: async () => {
                try {
                    // Best effort logout
                    await http.post('/auth/logout');
                } catch (e) { /* ignore */ }

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
                    // MOCK implementation for now - normally this would be a real network call
                    // const data = await http.post<{ token: string, user: User }>('/auth/refresh', { refreshToken });

                    // Simulating API delay and response
                    await new Promise(resolve => setTimeout(resolve, 500));

                    // In a real app, backend validates refresh token and issues new access token
                    // Here we just "mock" a success by retaining the current user/refresh token
                    // and generating a "new" mock access token if needed, or just acknowledging success.

                    const newToken = "mock_refreshed_access_token_" + Date.now();

                    console.log("[Auth] Token refreshed successfully");

                    set({
                        token: newToken,
                        loading: false,
                        error: null,
                        isAuthenticated: true
                    });
                } catch (err: any) {
                    console.error("Token refresh failed", err);
                    set({
                        token: null,
                        user: null,
                        isAuthenticated: false,
                        role: 'guest'
                    });
                    throw err;
                }
            },

            updateProfile: async (data) => {
                // TODO: Call API
                set((state) => ({
                    user: state.user ? { ...state.user, ...data } : null
                }));
            },

            switchToPlayerMode: () => {
                if (get().role === 'player' || get().playerProfile) {
                    console.log("Switching to player mode");
                }
            },

            checkAuth: async () => {
                const { token, refreshToken, refresh } = get();

                // If we have no token but have a refresh token (e.g. page reload), try to refresh
                if (!token && refreshToken) {
                    try {
                        await refresh();
                        return;
                    } catch (e) {
                        // Refresh failed, proceed to logout/guest state
                    }
                }

                if (!get().token) {
                    set({ isAuthenticated: false, role: 'guest' });
                    return;
                }

                try {
                    const data = await http.get<any>('/auth/me');
                    const userData = data.user || data;
                    set({
                        user: userData,
                        isAuthenticated: true,
                        role: userData.role || 'user'
                    });
                } catch (err) {
                    console.error("Token validation failed", err);
                    get().logout();
                }
            }
        }),
        {
            name: 'auth-storage',
            partialize: (state) => ({
                // SECURITY FLAGGED: Do NOT persist 'token' (Access Token) to localStorage
                // Only persist refreshToken (simulating httpOnly cookie behavior for this client-only phase)
                // and non-sensitive user preference/profile data
                refreshToken: state.refreshToken,
                user: state.user,
                role: state.role,
                playerProfile: state.playerProfile
            }),
            onRehydrateStorage: () => (state) => {
                if (state) {
                    // On hydration, if we have a refresh token but no access token, 
                    // we immediately try to restore the session.
                    if (state.refreshToken && !state.token) {
                        state.refresh().catch(() => {
                            console.log("Session restoration failed on hydration");
                        });
                    }
                }
            }
        }
    )
);
