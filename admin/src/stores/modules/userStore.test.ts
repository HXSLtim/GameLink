/**
 * User Store Tests
 * Tests user management, CRUD operations, filtering, and pagination
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useUserStore } from './userStore';
import type { User as ApiUser } from '@/api/admin';
import * as adminApi from '@/api/admin';

// Mock admin API
vi.mock('@/api/admin', () => ({
  adminApi: {
    getUsers: vi.fn(),
    createUser: vi.fn(),
    updateUser: vi.fn(),
    deleteUser: vi.fn(),
    batchDeleteUsers: vi.fn(),
    updateUserStatus: vi.fn(),
    batchUpdateUserStatus: vi.fn(),
    updateUserRole: vi.fn(),
    batchUpdateUserRole: vi.fn(),
  },
}));

describe('userStore', () => {
  const mockUsers: ApiUser[] = [
    {
      id: 1,
      name: 'Admin User',
      email: 'admin@gamelink.com',
      phone: '13800138001',
      avatarUrl: 'https://example.com/avatar1.jpg',
      role: 'admin',
      status: 'active',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 2,
      name: 'Player User',
      email: 'player@gamelink.com',
      phone: '13800138002',
      avatarUrl: 'https://example.com/avatar2.jpg',
      role: 'player',
      status: 'active',
      createdAt: '2024-01-02T00:00:00Z',
      updatedAt: '2024-01-02T00:00:00Z',
    },
    {
      id: 3,
      name: 'Banned User',
      email: 'banned@gamelink.com',
      phone: '13800138003',
      avatarUrl: 'https://example.com/avatar3.jpg',
      role: 'user',
      status: 'banned',
      createdAt: '2024-01-03T00:00:00Z',
      updatedAt: '2024-01-03T00:00:00Z',
    },
  ];

  beforeEach(() => {
    // Reset store state before each test
    useUserStore.getState().reset();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('State Initialization', () => {
    it('should initialize with default state', () => {
      const { result } = renderHook(() => useUserStore());

      expect(result.current.users).toEqual([]);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.pagination).toEqual({
        current: 1,
        pageSize: 10,
        total: 0,
      });
      expect(result.current.filters).toEqual({});
    });
  });

  describe('Fetch Users', () => {
    it('should fetch users successfully', async () => {
      const mockResponse = {
        data: {
          data: mockUsers,
          pagination: {
            total: 3,
          },
        },
      };

      vi.mocked(adminApi.adminApi.getUsers).mockResolvedValue(mockResponse as {
        data: { data: ApiUser[]; pagination: { total: number } };
      });

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUsers(1, 10);
      });

      expect(result.current.users).toEqual(mockUsers);
      expect(result.current.pagination).toEqual({
        current: 1,
        pageSize: 10,
        total: 3,
      });
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('should fetch users with filters', async () => {
      const mockResponse = {
        data: {
          data: [mockUsers[0]], // Only admin user
          pagination: {
            total: 1,
          },
        },
      };

      vi.mocked(adminApi.adminApi.getUsers).mockResolvedValue(mockResponse as {
        data: { data: ApiUser[]; pagination: { total: number } };
      });

      const { result } = renderHook(() => useUserStore());

      // Set filters
      act(() => {
        result.current.setFilters({ role: 'admin', status: 'active' });
      });

      await act(async () => {
        await result.current.fetchUsers(1, 10);
      });

      expect(adminApi.adminApi.getUsers).toHaveBeenCalledWith({
        page: 1,
        page_size: 10,
        keyword: undefined,
        role: ['admin'],
        status: ['active'],
        date_from: undefined,
        date_to: undefined,
      });

      expect(result.current.users.length).toBe(1);
      expect(result.current.users[0].role).toBe('admin');
    });

    it('should handle fetch users failure', async () => {
      const mockError = new Error('Failed to fetch users');

      vi.mocked(adminApi.adminApi.getUsers).mockRejectedValue(mockError);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          await result.current.fetchUsers();
        })
      ).rejects.toThrow('Failed to fetch users');
    });

    it('should handle empty response data', async () => {
      const mockResponse = {
        data: {
          data: null,
          pagination: {
            total: 0,
          },
        },
      };

      vi.mocked(adminApi.adminApi.getUsers).mockResolvedValue(mockResponse as {
        data: { data: ApiUser[]; pagination: { total: number } };
      });

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUsers();
      });

      expect(result.current.users).toEqual([]);
      expect(result.current.pagination.total).toBe(0);
    });
  });

  describe('Create User', () => {
    it('should create user successfully', async () => {
      const newUser: Partial<ApiUser> = {
        name: 'New User',
        email: 'new@gamelink.com',
        phone: '13800138004',
        role: 'user',
        status: 'active',
      };

      const mockResponse = {
        data: {
          data: {
            id: 4,
            ...newUser,
            createdAt: '2024-01-04T00:00:00Z',
            updatedAt: '2024-01-04T00:00:00Z',
          } as ApiUser,
        },
      };

      vi.mocked(adminApi.adminApi.createUser).mockResolvedValue(mockResponse as {
        data: { data: ApiUser };
      });

      const { result } = renderHook(() => useUserStore());

      // Initialize with some users
      act(() => {
        useUserStore.setState({ users: mockUsers, pagination: { current: 1, pageSize: 10, total: 3 } });
      });

      await act(async () => {
        await result.current.createUser(newUser);
      });

      expect(result.current.users.length).toBe(4);
      expect(result.current.users[0]).toEqual(mockResponse.data.data);
      expect(result.current.pagination.total).toBe(4);
      expect(result.current.loading).toBe(false);
    });

    it('should handle create user failure', async () => {
      const newUser: Partial<ApiUser> = {
        name: 'New User',
        email: 'new@gamelink.com',
        phone: '13800138004',
        role: 'user',
        status: 'active',
      };

      const mockError = new Error('Failed to create user');

      vi.mocked(adminApi.adminApi.createUser).mockRejectedValue(mockError);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          await result.current.createUser(newUser);
        })
      ).rejects.toThrow('Failed to create user');
    });
  });

  describe('Update User', () => {
    it('should update user successfully', async () => {
      const updatedData: Partial<ApiUser> = {
        name: 'Updated Admin User',
        email: 'updated@gamelink.com',
      };

      const mockResponse = {
        data: {
          data: {
            ...mockUsers[0],
            ...updatedData,
          } as ApiUser,
        },
      };

      vi.mocked(adminApi.adminApi.updateUser).mockResolvedValue(mockResponse as {
        data: { data: ApiUser };
      });

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUser(1, updatedData);
      });

      const updatedUser = result.current.users.find((u) => u.id === 1);
      expect(updatedUser?.name).toBe('Updated Admin User');
      expect(updatedUser?.email).toBe('updated@gamelink.com');
      expect(result.current.loading).toBe(false);
    });

    it('should handle update user failure', async () => {
      const mockError = new Error('Failed to update user');

      vi.mocked(adminApi.adminApi.updateUser).mockRejectedValue(mockError);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          await result.current.updateUser(1, { name: 'Updated' });
        })
      ).rejects.toThrow('Failed to update user');
    });
  });

  describe('Delete User', () => {
    it('should delete user successfully', async () => {
      vi.mocked(adminApi.adminApi.deleteUser).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          pagination: { current: 1, pageSize: 10, total: 3 },
        });
      });

      await act(async () => {
        await result.current.deleteUser(1);
      });

      expect(result.current.users.length).toBe(2);
      expect(result.current.users.find((u) => u.id === 1)).toBeUndefined();
      expect(result.current.pagination.total).toBe(2);
      expect(result.current.loading).toBe(false);
    });

    it('should handle delete user failure', async () => {
      const mockError = new Error('Failed to delete user');

      vi.mocked(adminApi.adminApi.deleteUser).mockRejectedValue(mockError);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          await result.current.deleteUser(1);
        })
      ).rejects.toThrow('Failed to delete user');
    });
  });

  describe('Batch Delete Users', () => {
    it('should batch delete users successfully', async () => {
      vi.mocked(adminApi.adminApi.batchDeleteUsers).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          pagination: { current: 1, pageSize: 10, total: 3 },
        });
      });

      await act(async () => {
        await result.current.batchDeleteUsers([1, 2]);
      });

      expect(result.current.users.length).toBe(1);
      expect(result.current.users[0].id).toBe(3);
      expect(result.current.pagination.total).toBe(1);
      expect(adminApi.adminApi.batchDeleteUsers).toHaveBeenCalledWith([1, 2]);
    });

    it('should handle batch delete failure', async () => {
      const mockError = new Error('Failed to batch delete');

      vi.mocked(adminApi.adminApi.batchDeleteUsers).mockRejectedValue(mockError);

      const { result } = renderHook(() => useUserStore());

      await expect(
        act(async () => {
          await result.current.batchDeleteUsers([1, 2]);
        })
      ).rejects.toThrow('Failed to batch delete');
    });
  });

  describe('Update User Status', () => {
    it('should update user status successfully', async () => {
      vi.mocked(adminApi.adminApi.updateUserStatus).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUserStatus(1, 'banned');
      });

      const updatedUser = result.current.users.find((u) => u.id === 1);
      expect(updatedUser?.status).toBe('banned');
      expect(adminApi.adminApi.updateUserStatus).toHaveBeenCalledWith(1, 'banned');
    });

    it('should batch update user status successfully', async () => {
      vi.mocked(adminApi.adminApi.batchUpdateUserStatus).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.batchUpdateUserStatus([1, 2], 'suspended');
      });

      expect(result.current.users[0].status).toBe('suspended');
      expect(result.current.users[1].status).toBe('suspended');
      expect(adminApi.adminApi.batchUpdateUserStatus).toHaveBeenCalledWith({
        userIds: [1, 2],
        status: 'suspended',
      });
    });
  });

  describe('Update User Role', () => {
    it('should update user role successfully', async () => {
      vi.mocked(adminApi.adminApi.updateUserRole).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUserRole(2, 'admin');
      });

      const updatedUser = result.current.users.find((u) => u.id === 2);
      expect(updatedUser?.role).toBe('admin');
      expect(adminApi.adminApi.updateUserRole).toHaveBeenCalledWith(2, 'admin');
    });

    it('should batch update user role successfully', async () => {
      vi.mocked(adminApi.adminApi.batchUpdateUserRole).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.batchUpdateUserRole([2, 3], 'admin');
      });

      expect(result.current.users[1].role).toBe('admin');
      expect(result.current.users[2].role).toBe('admin');
      expect(adminApi.adminApi.batchUpdateUserRole).toHaveBeenCalledWith({
        userIds: [2, 3],
        role: 'admin',
      });
    });
  });

  describe('Filters', () => {
    it('should set filters', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.setFilters({ role: 'admin', status: 'active' });
      });

      expect(result.current.filters).toEqual({
        role: 'admin',
        status: 'active',
      });
      expect(result.current.pagination.current).toBe(1); // Should reset to first page
    });

    it('should clear filters', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.setFilters({ role: 'admin' });
        result.current.clearFilters();
      });

      expect(result.current.filters).toEqual({});
      expect(result.current.pagination).toEqual({
        current: 1,
        pageSize: 10,
        total: 0,
      });
    });

    it('should partially update filters', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        result.current.setFilters({ role: 'admin', status: 'active' });
        result.current.setFilters({ keyword: 'test' });
      });

      expect(result.current.filters).toEqual({
        role: 'admin',
        status: 'active',
        keyword: 'test',
      });
    });
  });

  describe('Selectors', () => {
    it('should get user by id', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      const user = result.current.getUserById(1);
      expect(user).toEqual(mockUsers[0]);

      const nonExistent = result.current.getUserById(999);
      expect(nonExistent).toBeUndefined();
    });

    it('should get active users', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      const activeUsers = result.current.getActiveUsers();
      expect(activeUsers.length).toBe(2);
      expect(activeUsers.every((u) => u.status === 'active')).toBe(true);
    });

    it('should get users by role', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      const admins = result.current.getUsersByRole('admin');
      expect(admins.length).toBe(1);
      expect(admins[0].role).toBe('admin');

      const players = result.current.getUsersByRole('player');
      expect(players.length).toBe(1);
      expect(players[0].role).toBe('player');
    });

    it('should get users by status', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      const activeUsers = result.current.getUsersByStatus('active');
      expect(activeUsers.length).toBe(2);

      const bannedUsers = result.current.getUsersByStatus('banned');
      expect(bannedUsers.length).toBe(1);
      expect(bannedUsers[0].status).toBe('banned');
    });
  });

  describe('Reset', () => {
    it('should reset store to initial state', () => {
      const { result } = renderHook(() => useUserStore());

      // Set some state
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          loading: true,
          error: 'Test error',
          filters: { role: 'admin' },
          pagination: { current: 2, pageSize: 20, total: 100 },
        });
      });

      expect(result.current.users.length).toBe(3);
      expect(result.current.loading).toBe(true);
      expect(result.current.error).toBe('Test error');

      // Reset
      act(() => {
        result.current.reset();
      });

      expect(result.current.users).toEqual([]);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.filters).toEqual({});
      expect(result.current.pagination).toEqual({
        current: 1,
        pageSize: 10,
        total: 0,
      });
    });
  });

  describe('Edge Cases', () => {
    it('should handle empty user list for selectors', () => {
      const { result } = renderHook(() => useUserStore());

      expect(result.current.getUserById(1)).toBeUndefined();
      expect(result.current.getActiveUsers()).toEqual([]);
      expect(result.current.getUsersByRole('admin')).toEqual([]);
      expect(result.current.getUsersByStatus('active')).toEqual([]);
    });

    it('should prevent negative total count on delete', async () => {
      vi.mocked(adminApi.adminApi.deleteUser).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with one user
      act(() => {
        useUserStore.setState({
          users: [mockUsers[0]],
          pagination: { current: 1, pageSize: 10, total: 0 },
        });
      });

      await act(async () => {
        await result.current.deleteUser(1);
      });

      expect(result.current.pagination.total).toBe(0); // Should not go negative
    });

    it('should handle pagination correctly with batch operations', async () => {
      vi.mocked(adminApi.adminApi.batchDeleteUsers).mockResolvedValue({} as Record<string, never>);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          pagination: { current: 1, pageSize: 10, total: 3 },
        });
      });

      await act(async () => {
        await result.current.batchDeleteUsers([1, 2, 3]); // Delete all
      });

      expect(result.current.pagination.total).toBe(0);
      expect(result.current.users).toEqual([]);
    });
  });
});
