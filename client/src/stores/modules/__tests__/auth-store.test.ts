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
            permissions: [],
            viewMode: 'user',
            isPlayer: false,
            playerProfile: null,
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
            user: { id: '1', username: 'test', avatar: 'avatar.png', nickname: 'Test User' },
            role: 'user',
            permissions: [],
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().login({ username: 'test', password: 'password' });

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(true);
        expect(state.token).toBe('fake-token');
        expect(state.user?.id).toBe(1);
        expect(state.user?.username).toBe('test');
        expect(state.user?.avatar).toBe('avatar.png');
        expect(state.user?.name).toBe('Test User');
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

    it('should register successfully', async () => {
        const mockResponse = {
            token: 'new-token',
            refreshToken: 'new-refresh',
            user: { id: '2', username: 'newuser', avatar: 'default.png', nickname: 'New User' },
            role: 'user',
            permissions: [],
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().register({
            email: 'test@example.com',
            password: 'password123',
            name: 'Test User',
        });

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(true);
        expect(state.token).toBe('new-token');
        expect(state.user?.username).toBe('newuser');
        expect(state.user?.name).toBe('New User');
        expect(http.post).toHaveBeenCalledWith('/auth/register', {
            email: 'test@example.com',
            password: 'password123',
            name: 'Test User',
        });
    });

    it('should handle register error', async () => {
        const errorMessage = 'Email already exists';
        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error(errorMessage));

        await expect(useAuthStore.getState().register({
            email: 'test@example.com',
            password: 'password123',
            name: 'Test',
        })).rejects.toThrow(errorMessage);

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.error).toBe(errorMessage);
    });

    it('should logout successfully', async () => {
        useAuthStore.setState({
            token: 'token',
            isAuthenticated: true,
            user: { id: 1, username: 'test', avatar: '', name: 'test' },
            role: 'user',
            permissions: ['read'],
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useAuthStore.getState().logout();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.token).toBeNull();
        expect(state.user).toBeNull();
        expect(state.role).toBe('guest');
        expect(state.permissions).toEqual([]);
        expect(state.playerProfile).toBeNull();
        expect(http.post).toHaveBeenCalledWith('/auth/logout');
    });

    it('should logout even when API call fails', async () => {
        useAuthStore.setState({
            token: 'token',
            isAuthenticated: true,
            user: { id: 1, username: 'test', avatar: '', name: 'test' },
        });

        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

        await useAuthStore.getState().logout();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.token).toBeNull();
        expect(state.user).toBeNull();
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

    it('should use old refresh token when new one not provided', async () => {
        useAuthStore.setState({
            token: 'old-token',
            refreshToken: 'old-refresh',
        });

        const mockResponse = {
            token: 'new-token',
            user: { id: '1', username: 'test', avatar: '' }
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().refresh();

        const state = useAuthStore.getState();
        expect(state.token).toBe('new-token');
        expect(state.refreshToken).toBe('old-refresh');
    });

    it('should handle refresh token error', async () => {
        useAuthStore.setState({
            token: 'old-token',
            refreshToken: 'old-refresh',
            isAuthenticated: true,
            role: 'user',
        });

        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Refresh failed'));

        await expect(useAuthStore.getState().refresh()).rejects.toThrow();

        const state = useAuthStore.getState();
        expect(state.token).toBeNull();
        expect(state.refreshToken).toBeNull();
        expect(state.isAuthenticated).toBe(false);
        expect(state.role).toBe('guest');
    });

    it('should throw error when no refresh token available', async () => {
        useAuthStore.setState({
            token: 'expired-token',
            refreshToken: null,
        });

        await expect(useAuthStore.getState().refresh()).rejects.toThrow('No refresh token available');
    });

    it('should prevent concurrent refresh calls', async () => {
        useAuthStore.setState({
            token: 'old-token',
            refreshToken: 'old-refresh',
        });

        const mockResponse = {
            token: 'new-token',
            refreshToken: 'new-refresh',
        };
        (http.post as ReturnType<typeof vi.fn>).mockImplementation(() =>
            new Promise(resolve => setTimeout(() => resolve(mockResponse), 100))
        );

        const promise1 = useAuthStore.getState().refresh();
        const promise2 = useAuthStore.getState().refresh();

        await Promise.all([promise1, promise2]);

        expect(http.post).toHaveBeenCalledTimes(1);
    });

    it('should update profile successfully', async () => {
        useAuthStore.setState({
            user: { id: 1, username: 'test', avatar: 'old.png', name: 'Test' },
        });

        const updatedUser = { id: 1, username: 'test', avatar: 'new.png', name: 'Updated' };
        (http.put as ReturnType<typeof vi.fn>).mockResolvedValue(updatedUser);

        await useAuthStore.getState().updateProfile({ avatar: 'new.png' });

        const state = useAuthStore.getState();
        expect(state.user?.avatar).toBe('new.png');
        expect(state.user?.name).toBe('Updated');
        expect(http.put).toHaveBeenCalledWith('/user/profile', { avatar: 'new.png' });
    });

    it('should handle update profile error', async () => {
        useAuthStore.setState({
            user: { id: 1, username: 'test', avatar: '', name: 'test' },
        });

        (http.put as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Update failed'));

        await expect(useAuthStore.getState().updateProfile({ name: 'New' }))
            .rejects.toThrow('Update failed');

        const state = useAuthStore.getState();
        expect(state.error).toBe('Update failed');
        expect(state.loading).toBe(false);
    });

    it('should change password successfully', async () => {
        (http.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useAuthStore.getState().changePassword({
            oldPassword: 'oldpass',
            newPassword: 'newpass',
        });

        expect(http.put).toHaveBeenCalledWith('/user/password', {
            oldPassword: 'oldpass',
            newPassword: 'newpass',
        });
        expect(useAuthStore.getState().loading).toBe(false);
    });

    it('should change password without old password', async () => {
        (http.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useAuthStore.getState().changePassword({
            newPassword: 'newpass',
        });

        expect(http.put).toHaveBeenCalledWith('/user/password', {
            newPassword: 'newpass',
        });
    });

    it('should handle change password error', async () => {
        (http.put as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Wrong password'));

        await expect(useAuthStore.getState().changePassword({
            oldPassword: 'wrong',
            newPassword: 'newpass',
        })).rejects.toThrow('Wrong password');

        expect(useAuthStore.getState().error).toBe('Wrong password');
    });

    it('should switch to player mode when isPlayer is true', () => {
        useAuthStore.setState({ isPlayer: true, viewMode: 'user' });

        useAuthStore.getState().switchToPlayerMode();

        expect(useAuthStore.getState().viewMode).toBe('player');
    });

    it('should not switch to player mode when isPlayer is false', () => {
        useAuthStore.setState({ isPlayer: false, viewMode: 'user' });

        useAuthStore.getState().switchToPlayerMode();

        expect(useAuthStore.getState().viewMode).toBe('user');
    });

    it('should switch to user mode', () => {
        useAuthStore.setState({ viewMode: 'player' });

        useAuthStore.getState().switchToUserMode();

        expect(useAuthStore.getState().viewMode).toBe('user');
    });

    it('should set isPlayer flag', () => {
        useAuthStore.setState({ isPlayer: false });

        useAuthStore.getState().setIsPlayer(true);

        expect(useAuthStore.getState().isPlayer).toBe(true);
    });

    it('should switch to user mode when isPlayer is set to false', () => {
        useAuthStore.setState({ isPlayer: true, viewMode: 'player' });

        useAuthStore.getState().setIsPlayer(false);

        expect(useAuthStore.getState().isPlayer).toBe(false);
        expect(useAuthStore.getState().viewMode).toBe('user');
    });

    it('should check auth successfully with valid token', async () => {
        useAuthStore.setState({
            token: 'valid-token',
            refreshToken: 'refresh-token',
        });

        const mockResponse = {
            user: { id: '123', username: 'testuser', avatar: 'avatar.png', nickname: 'Test User' },
            role: 'user',
        };
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(true);
        expect(state.user?.id).toBe(123);
        expect(state.user?.username).toBe('testuser');
        expect(state.user?.name).toBe('Test User');
        expect(http.get).toHaveBeenCalledWith('/auth/me');
    });

    it('should check auth and refresh when token is missing but refresh token exists', async () => {
        useAuthStore.setState({
            token: null,
            refreshToken: 'valid-refresh',
        });

        const mockRefreshResponse = {
            token: 'new-token',
            refreshToken: 'new-refresh',
            user: { id: '1', username: 'test', avatar: '' },
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockRefreshResponse);

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.token).toBe('new-token');
        expect(state.isAuthenticated).toBe(true);
    });

    it('should set guest state when no token and refresh fails', async () => {
        useAuthStore.setState({
            token: null,
            refreshToken: 'invalid-refresh',
            isAuthenticated: true,
            role: 'user',
        });

        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Refresh failed'));

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.role).toBe('guest');
    });

    it('should set guest state when no token available', async () => {
        useAuthStore.setState({
            token: null,
            refreshToken: null,
        });

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.role).toBe('guest');
        expect(http.get).not.toHaveBeenCalled();
    });

    it('should logout on 401 error during checkAuth', async () => {
        useAuthStore.setState({
            token: 'expired-token',
            isAuthenticated: true,
            user: { id: 1, username: 'test', avatar: '', name: 'test' },
        });

        const error = new Error('Unauthorized');
        (error as any).response = { status: 401 };
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(error);

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(false);
        expect(state.user).toBeNull();
    });

    it('should keep session on non-401 error during checkAuth', async () => {
        useAuthStore.setState({
            token: 'valid-token',
            isAuthenticated: true,
            user: { id: 1, username: 'test', avatar: '', name: 'test' },
        });

        const error = new Error('Network error');
        (error as any).response = { status: 500 };
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(error);

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.isAuthenticated).toBe(true);
        expect(state.user).not.toBeNull();
    });

    it('should handle API response without user field in checkAuth', async () => {
        useAuthStore.setState({
            token: 'valid-token',
        });

        const mockResponse = {
            id: '123',
            username: 'testuser',
            avatar: 'avatar.png',
            nickname: 'Test User',
            role: 'user',
        };
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().checkAuth();

        const state = useAuthStore.getState();
        expect(state.user?.id).toBe(123);
        expect(state.user?.username).toBe('testuser');
    });

    it('should map username to name when nickname is missing', async () => {
        const mockResponse = {
            token: 'token',
            refreshToken: 'refresh',
            user: { id: '1', username: 'testuser', avatar: '' },
            role: 'user',
            permissions: [],
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().login({ username: 'testuser', password: 'pass' });

        expect(useAuthStore.getState().user?.name).toBe('testuser');
    });

    it('should handle admin role', async () => {
        const mockResponse = {
            token: 'admin-token',
            refreshToken: 'refresh',
            user: { id: '1', username: 'admin', avatar: '' },
            role: 'admin',
            permissions: ['all'],
        };
        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

        await useAuthStore.getState().login({ username: 'admin', password: 'pass' });

        expect(useAuthStore.getState().role).toBe('admin');
        expect(useAuthStore.getState().permissions).toEqual(['all']);
    });
});
