import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuthStore } from '../auth-store';
import { http } from '@/lib/http';

// Mock http client
vi.mock('@/lib/http', () => ({
    http: {
        post: vi.fn(),
        get: vi.fn(),
        put: vi.fn(),
    },
}));

describe('Auth Store', () => {
    beforeEach(() => {
        // Reset store state
        useAuthStore.setState({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
            role: 'guest',
            loading: false,
            error: null,
        });
        vi.clearAllMocks();
    });

    it('should set loading state correctly during login', async () => {
        const mockResponse = {
            token: 'fake-token',
            refreshToken: 'fake-refresh',
            user: { id: '1', username: 'test', avatar: 'avatar.png' },
            role: 'user',
            permissions: [],
        };
        (http.post as ReturnType<typeof vi.fn>).mockImplementation(() => new Promise(resolve => setTimeout(() => resolve(mockResponse), 10)));

        const loginPromise = useAuthStore.getState().login({ username: 'test', password: 'password' });

        expect(useAuthStore.getState().loading).toBe(true);
        await loginPromise;
        expect(useAuthStore.getState().loading).toBe(false);
    });

    it('should login successfully', async () => {
        const mockResponse = {
            token: 'fake-token',
            refreshToken: 'fake-refresh',
            user: { id: '1', username: 'test', avatar: 'avatar.png' },
            role: 'user',
            permissions: [],
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().login({ username: 'test', password: 'password' });

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(true);
        expect(state.token).toBe('fake-token');
        expect(state.user).toEqual(mockResponse.user);
        expect(state.error).toBeNull();
        expect(http.post).toHaveBeenCalledWith('/auth/login', { username: 'test', password: 'password' });
    });

    it('should handle login error', async () => {
        const errorMessage = 'Invalid credentials';
        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error(errorMessage));

        await expect(useAuthStore.getState().login({ username: 'test', password: 'wrong' }))
            .rejects.toThrow(errorMessage);

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.token).toBeNull();
        expect(state.error).toBe(errorMessage);
    });

    it('should logout successfully', async () => {
        // Setup initial state
        useAuthStore.setState({
            token: 'token',
            isAuthenticated: true,
            user: { id: '1', username: 'test', avatar: '' }
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useAuthStore.getState().logout();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.token).toBeNull();
        expect(state.user).toBeNull();
        expect(http.post).toHaveBeenCalledWith('/auth/logout');
    });

    it('should refresh token successfully', async () => {
        useAuthStore.setState({
            token: 'old-token',
            refreshToken: 'old-refresh',
        });

        const mockResponse = {
            token: 'new-token',
            refreshToken: 'new-refresh',
            user: { id: '1', username: 'test', avatar: '' }
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().refresh();

        const state = useAuthStore.getState();
        expect(state.token).toBe('new-token');
        expect(state.refreshToken).toBe('new-refresh');
        expect(http.post).toHaveBeenCalledWith('/public/auth/refresh', { refreshToken: 'old-refresh' });
    });

    it('should handle refresh token error', async () => {
        useAuthStore.setState({
            token: 'old-token',
            refreshToken: 'old-refresh',
            isAuthenticated: true,
        });

        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Refresh failed'));

        await expect(useAuthStore.getState().refresh()).rejects.toThrow();

        const state = useAuthStore.getState();
        expect(state.token).toBeNull();
        expect(state.isAuthenticated).toBe(false);
    });
});
