import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuthStore } from './auth-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        post: vi.fn(),
        get: vi.fn(),
        put: vi.fn(),
    },
}));

import { http } from '@/lib/http';

describe('auth-store', () => {
    beforeEach(() => {
        // Reset store state
        useAuthStore.setState({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
            role: 'guest',
            permissions: [],
            playerProfile: null,
            loading: false,
            error: null,
        });
        vi.clearAllMocks();
    });

    describe('login', () => {
        it('should login successfully and update state', async () => {
            const mockResponse = {
                token: 'test-token',
                refreshToken: 'test-refresh-token',
                user: { id: '1', username: 'testuser', avatar: '' },
                role: 'user',
                permissions: ['read'],
            };

            vi.mocked(http.post).mockResolvedValueOnce(mockResponse);

            const { login } = useAuthStore.getState();
            await login({ username: 'testuser', password: 'password123' });

            const state = useAuthStore.getState();
            expect(state.isAuthenticated).toBe(true);
            expect(state.token).toBe('test-token');
            expect(state.refreshToken).toBe('test-refresh-token');
            expect(state.user?.username).toBe('testuser');
            expect(state.role).toBe('user');
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should handle login failure', async () => {
            const mockError = {
                response: { data: { message: 'Invalid credentials' } },
            };

            vi.mocked(http.post).mockRejectedValueOnce(mockError);

            const { login } = useAuthStore.getState();

            await expect(login({ username: 'testuser', password: 'wrong' })).rejects.toEqual(mockError);

            const state = useAuthStore.getState();
            expect(state.isAuthenticated).toBe(false);
            expect(state.token).toBeNull();
            expect(state.error).toBe('Invalid credentials');
            expect(state.loading).toBe(false);
        });
    });

    describe('register', () => {
        it('should register successfully and auto-login', async () => {
            const mockResponse = {
                token: 'new-token',
                refreshToken: 'new-refresh-token',
                user: { id: '2', username: 'newuser', avatar: '' },
                role: 'user',
            };

            vi.mocked(http.post).mockResolvedValueOnce(mockResponse);

            const { register } = useAuthStore.getState();
            await register({
                email: 'test@example.com',
                password: 'password123',
                name: 'Test User',
            });

            const state = useAuthStore.getState();
            expect(state.isAuthenticated).toBe(true);
            expect(state.token).toBe('new-token');
            expect(state.user?.username).toBe('newuser');
        });
    });

    describe('logout', () => {
        it('should clear auth state on logout', async () => {
            // Set initial authenticated state
            useAuthStore.setState({
                user: { id: '1', username: 'testuser', avatar: '' },
                token: 'test-token',
                refreshToken: 'test-refresh-token',
                isAuthenticated: true,
                role: 'user',
            });

            vi.mocked(http.post).mockResolvedValueOnce({});

            const { logout } = useAuthStore.getState();
            await logout();

            const state = useAuthStore.getState();
            expect(state.isAuthenticated).toBe(false);
            expect(state.token).toBeNull();
            expect(state.refreshToken).toBeNull();
            expect(state.user).toBeNull();
            expect(state.role).toBe('guest');
        });
    });

    describe('refresh', () => {
        it('should refresh token successfully', async () => {
            useAuthStore.setState({
                refreshToken: 'old-refresh-token',
                user: { id: '1', username: 'testuser', avatar: '' },
            });

            const mockResponse = {
                token: 'new-access-token',
                refreshToken: 'new-refresh-token',
                user: { id: '1', username: 'testuser', avatar: '' },
            };

            vi.mocked(http.post).mockResolvedValueOnce(mockResponse);

            const { refresh } = useAuthStore.getState();
            await refresh();

            const state = useAuthStore.getState();
            expect(state.token).toBe('new-access-token');
            expect(state.refreshToken).toBe('new-refresh-token');
            expect(state.isAuthenticated).toBe(true);
        });

        it('should throw error if no refresh token', async () => {
            useAuthStore.setState({ refreshToken: null });

            const { refresh } = useAuthStore.getState();
            await expect(refresh()).rejects.toThrow('No refresh token available');
        });

        it('should clear state on refresh failure', async () => {
            useAuthStore.setState({
                refreshToken: 'old-refresh-token',
                token: 'old-token',
                isAuthenticated: true,
            });

            vi.mocked(http.post).mockRejectedValueOnce(new Error('Token expired'));

            const { refresh } = useAuthStore.getState();
            await expect(refresh()).rejects.toThrow();

            const state = useAuthStore.getState();
            expect(state.token).toBeNull();
            expect(state.refreshToken).toBeNull();
            expect(state.isAuthenticated).toBe(false);
            expect(state.role).toBe('guest');
        });
    });
});
