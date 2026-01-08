/**
 * User Store Tests
 * Tests user management, CRUD operations, filtering, and pagination
 * Uses mocked UserService for testing
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { createUserStore } from './userStore';
import type { User as ApiUser } from '@/api/admin';
import type { IUserService, UserValidationResult, UserExportData } from '@/services/domain';
import type { ServiceResult, BatchResult } from '@/services/utils';

// Create mock UserService
const createMockUserService = (): IUserService => ({
  getUsers: vi.fn(),
  getUserById: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  updateUserStatus: vi.fn(),
  updateUserRole: vi.fn(),
  batchUpdateStatus: vi.fn(),
  batchUpdateRole: vi.fn(),
  batchDelete: vi.fn(),
  validateUserData: vi.fn(),
  validateEmail: vi.fn(),
  validatePhone: vi.fn(),
  validatePassword: vi.fn(),
  exportUsers: vi.fn(),
});

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

  let mockUserService: IUserService;
  let useUserStore: ReturnType<typeof createUserStore>;

  beforeEach(() => {
    // Create fresh mock service and store for each test
    mockUserService = createMockUserService();
    useUserStore = createUserStore(mockUserService);
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
      expect(result.current.lastBatchResult).toBeNull();
    });
  });

  describe('Fetch Users', () => {
    it('should fetch users successfully', async () => {
      const mockResult: ServiceResult<ApiUser[]> = {
        success: true,
        data: mockUsers,
      };

      vi.mocked(mockUserService.getUsers).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUsers(1, 10);
      });

      expect(result.current.users).toEqual(mockUsers);
      expect(result.current.pagination.current).toBe(1);
      expect(result.current.pagination.pageSize).toBe(10);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('should fetch users with filters', async () => {
      const filteredUsers = [mockUsers[0]]; // Only admin user
      const mockResult: ServiceResult<ApiUser[]> = {
        success: true,
        data: filteredUsers,
      };

      vi.mocked(mockUserService.getUsers).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Set filters
      act(() => {
        result.current.setFilters({ role: 'admin', status: 'active' });
      });

      await act(async () => {
        await result.current.fetchUsers(1, 10);
      });

      expect(mockUserService.getUsers).toHaveBeenCalledWith({
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
      const mockResult: ServiceResult<ApiUser[]> = {
        success: false,
        error: {
          code: 'FETCH_ERROR',
          message: 'Failed to fetch users',
        },
      };

      vi.mocked(mockUserService.getUsers).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      let thrownError: Error | null = null;
      try {
        await act(async () => {
          await result.current.fetchUsers();
        });
      } catch (e) {
        thrownError = e as Error;
      }

      // Verify error was thrown
      expect(thrownError).not.toBeNull();
      expect(thrownError?.message).toBe('Failed to fetch users');

      // The error should be set in the store
      expect(result.current.error).not.toBeNull();
      expect(result.current.error?.code).toBe('FETCH_ERROR');
      expect(result.current.error?.message).toBe('Failed to fetch users');
    });

    it('should handle empty response data', async () => {
      const mockResult: ServiceResult<ApiUser[]> = {
        success: true,
        data: [],
      };

      vi.mocked(mockUserService.getUsers).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        await result.current.fetchUsers();
      });

      expect(result.current.users).toEqual([]);
    });
  });

  describe('Create User', () => {
    it('should create user successfully', async () => {
      const newUserData = {
        name: 'New User',
        email: 'new@gamelink.com',
        phone: '13800138004',
        password: 'Password123!',
        role: 'user' as const,
        status: 'active' as const,
      };

      const createdUser: ApiUser = {
        id: 4,
        ...newUserData,
        avatarUrl: '',
        createdAt: '2024-01-04T00:00:00Z',
        updatedAt: '2024-01-04T00:00:00Z',
      };

      const mockResult: ServiceResult<ApiUser> = {
        success: true,
        data: createdUser,
      };

      vi.mocked(mockUserService.createUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with some users
      act(() => {
        useUserStore.setState({ users: mockUsers, pagination: { current: 1, pageSize: 10, total: 3 } });
      });

      await act(async () => {
        await result.current.createUser(newUserData);
      });

      expect(result.current.users.length).toBe(4);
      expect(result.current.users[0]).toEqual(createdUser);
      expect(result.current.pagination.total).toBe(4);
      expect(result.current.loading).toBe(false);
    });

    it('should handle create user failure with validation error', async () => {
      const newUserData = {
        name: 'New User',
        email: 'invalid-email',
        phone: '13800138004',
        password: 'weak',
        role: 'user' as const,
        status: 'active' as const,
      };

      const mockResult: ServiceResult<ApiUser> = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'User data validation failed',
          details: {
            errors: [
              { field: 'email', message: 'Invalid email format' },
              { field: 'password', message: 'Password too weak' },
            ],
          },
        },
      };

      vi.mocked(mockUserService.createUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        const createResult = await result.current.createUser(newUserData);
        expect(createResult.success).toBe(false);
        expect(createResult.error?.code).toBe('VALIDATION_ERROR');
      });

      expect(result.current.error?.code).toBe('VALIDATION_ERROR');
    });
  });

  describe('Update User', () => {
    it('should update user successfully', async () => {
      const updatedData = {
        name: 'Updated Admin User',
        email: 'updated@gamelink.com',
        phone: '13800138001',
      };

      const updatedUser: ApiUser = {
        ...mockUsers[0],
        ...updatedData,
      };

      const mockResult: ServiceResult<ApiUser> = {
        success: true,
        data: updatedUser,
      };

      vi.mocked(mockUserService.updateUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUser(1, updatedData);
      });

      const user = result.current.users.find((u) => u.id === 1);
      expect(user?.name).toBe('Updated Admin User');
      expect(user?.email).toBe('updated@gamelink.com');
      expect(result.current.loading).toBe(false);
    });

    it('should handle update user failure', async () => {
      const mockResult: ServiceResult<ApiUser> = {
        success: false,
        error: {
          code: 'NOT_FOUND',
          message: 'User not found',
        },
      };

      vi.mocked(mockUserService.updateUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        const updateResult = await result.current.updateUser(999, { name: 'Updated' });
        expect(updateResult.success).toBe(false);
      });

      expect(result.current.error?.code).toBe('NOT_FOUND');
    });
  });

  describe('Delete User', () => {
    it('should delete user successfully', async () => {
      const mockResult: ServiceResult<void> = {
        success: true,
      };

      vi.mocked(mockUserService.deleteUser).mockResolvedValue(mockResult);

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
      const mockResult: ServiceResult<void> = {
        success: false,
        error: {
          code: 'FORBIDDEN',
          message: 'Cannot delete admin user',
        },
      };

      vi.mocked(mockUserService.deleteUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      await act(async () => {
        const deleteResult = await result.current.deleteUser(1);
        expect(deleteResult.success).toBe(false);
      });

      expect(result.current.error?.code).toBe('FORBIDDEN');
    });
  });

  describe('Batch Delete Users', () => {
    it('should batch delete users successfully', async () => {
      const mockResult: BatchResult<void> = {
        success: true,
        total: 2,
        succeeded: 2,
        failed: 0,
        results: [
          { index: 0, success: true },
          { index: 1, success: true },
        ],
      };

      vi.mocked(mockUserService.batchDelete).mockResolvedValue(mockResult);

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
      expect(result.current.lastBatchResult).toEqual(mockResult);
    });

    it('should handle partial batch delete failure', async () => {
      const mockResult: BatchResult<void> = {
        success: false,
        total: 2,
        succeeded: 1,
        failed: 1,
        results: [
          { index: 0, success: true },
          { index: 1, success: false, error: { code: 'FORBIDDEN', message: 'Cannot delete' } },
        ],
      };

      vi.mocked(mockUserService.batchDelete).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          pagination: { current: 1, pageSize: 10, total: 3 },
        });
      });

      await act(async () => {
        const batchResult = await result.current.batchDeleteUsers([1, 2]);
        expect(batchResult.success).toBe(false);
        expect(batchResult.succeeded).toBe(1);
        expect(batchResult.failed).toBe(1);
      });

      // Only user 1 should be deleted (index 0)
      expect(result.current.users.length).toBe(2);
      expect(result.current.error?.code).toBe('BATCH_PARTIAL_FAILURE');
    });
  });

  describe('Update User Status', () => {
    it('should update user status successfully', async () => {
      const updatedUser: ApiUser = { ...mockUsers[0], status: 'banned' };
      const mockResult: ServiceResult<ApiUser> = {
        success: true,
        data: updatedUser,
      };

      vi.mocked(mockUserService.updateUserStatus).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUserStatus(1, 'banned');
      });

      const user = result.current.users.find((u) => u.id === 1);
      expect(user?.status).toBe('banned');
    });

    it('should batch update user status successfully', async () => {
      const mockResult: BatchResult<void> = {
        success: true,
        total: 2,
        succeeded: 2,
        failed: 0,
        results: [
          { index: 0, success: true },
          { index: 1, success: true },
        ],
      };

      vi.mocked(mockUserService.batchUpdateStatus).mockResolvedValue(mockResult);

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
    });
  });

  describe('Update User Role', () => {
    it('should update user role successfully', async () => {
      const updatedUser: ApiUser = { ...mockUsers[1], role: 'admin' };
      const mockResult: ServiceResult<ApiUser> = {
        success: true,
        data: updatedUser,
      };

      vi.mocked(mockUserService.updateUserRole).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      await act(async () => {
        await result.current.updateUserRole(2, 'admin');
      });

      const user = result.current.users.find((u) => u.id === 2);
      expect(user?.role).toBe('admin');
    });

    it('should batch update user role successfully', async () => {
      const mockResult: BatchResult<void> = {
        success: true,
        total: 2,
        succeeded: 2,
        failed: 0,
        results: [
          { index: 0, success: true },
          { index: 1, success: true },
        ],
      };

      vi.mocked(mockUserService.batchUpdateRole).mockResolvedValue(mockResult);

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
    });
  });

  describe('Validation', () => {
    it('should validate user data using service', () => {
      const mockValidationResult: UserValidationResult = {
        valid: false,
        errors: [
          { field: 'email', message: 'Invalid email format' },
        ],
      };

      vi.mocked(mockUserService.validateUserData).mockReturnValue(mockValidationResult);

      const { result } = renderHook(() => useUserStore());

      const validation = result.current.validateUserData({ email: 'invalid' });

      expect(validation.valid).toBe(false);
      expect(validation.errors).toHaveLength(1);
      expect(mockUserService.validateUserData).toHaveBeenCalledWith({ email: 'invalid' });
    });
  });

  describe('Export', () => {
    it('should export users using service', () => {
      const mockExportData: UserExportData = {
        headers: ['ID', '姓名', '邮箱'],
        rows: [['1', 'Admin User', 'admin@gamelink.com']],
      };

      vi.mocked(mockUserService.exportUsers).mockReturnValue(mockExportData);

      const { result } = renderHook(() => useUserStore());

      // Initialize with users
      act(() => {
        useUserStore.setState({ users: mockUsers });
      });

      const exportData = result.current.exportUsers();

      expect(exportData.headers).toEqual(['ID', '姓名', '邮箱']);
      expect(mockUserService.exportUsers).toHaveBeenCalledWith(mockUsers);
    });

    it('should export specific users when provided', () => {
      const mockExportData: UserExportData = {
        headers: ['ID', '姓名', '邮箱'],
        rows: [['1', 'Admin User', 'admin@gamelink.com']],
      };

      vi.mocked(mockUserService.exportUsers).mockReturnValue(mockExportData);

      const { result } = renderHook(() => useUserStore());

      const specificUsers = [mockUsers[0]];
      const exportData = result.current.exportUsers(specificUsers);

      expect(mockUserService.exportUsers).toHaveBeenCalledWith(specificUsers);
      expect(exportData).toEqual(mockExportData);
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

  describe('Reset and Clear Error', () => {
    it('should reset store to initial state', () => {
      const { result } = renderHook(() => useUserStore());

      // Set some state
      act(() => {
        useUserStore.setState({
          users: mockUsers,
          loading: true,
          error: { code: 'TEST_ERROR', message: 'Test error' },
          filters: { role: 'admin' },
          pagination: { current: 2, pageSize: 20, total: 100 },
          lastBatchResult: { success: true, total: 1, succeeded: 1, failed: 0, results: [] },
        });
      });

      expect(result.current.users.length).toBe(3);
      expect(result.current.loading).toBe(true);
      expect(result.current.error?.message).toBe('Test error');

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
      expect(result.current.lastBatchResult).toBeNull();
    });

    it('should clear error state', () => {
      const { result } = renderHook(() => useUserStore());

      act(() => {
        useUserStore.setState({
          error: { code: 'TEST_ERROR', message: 'Test error' },
        });
      });

      expect(result.current.error).not.toBeNull();

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
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
      const mockResult: ServiceResult<void> = {
        success: true,
      };

      vi.mocked(mockUserService.deleteUser).mockResolvedValue(mockResult);

      const { result } = renderHook(() => useUserStore());

      // Initialize with one user and zero total
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
      const mockResult: BatchResult<void> = {
        success: true,
        total: 3,
        succeeded: 3,
        failed: 0,
        results: [
          { index: 0, success: true },
          { index: 1, success: true },
          { index: 2, success: true },
        ],
      };

      vi.mocked(mockUserService.batchDelete).mockResolvedValue(mockResult);

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
