/**
 * Admin API Tests
 *
 * Tests for admin operations including user management, game management,
 * player management, order management, and statistics
 *
 * Requirements:
 * - CRUD operations for users, games, players, orders
 * - Batch operations support
 * - Filtering and pagination
 * - Error handling for all endpoints
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  adminApi,
  type User,
  type Game,
  type Player,
  type Order,
  type CreateUserDto,
  type UpdateUserDto,
  type UserQueryParams,
  type CreateGameDto,
  type UpdateGameDto,
  type CreatePlayerDto,
  type UpdatePlayerDto,
  type OrderQueryParams,
  type ApiResponse,
} from './admin';
import apiClient from './client';
import type { AxiosResponse } from 'axios';

// Mock the API client
vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn(), handlers: [] },
      response: { use: vi.fn(), handlers: [] },
    },
    defaults: {
      baseURL: '/api/v1',
      timeout: 10000,
      headers: {},
    },
  },
}));

// Mock crypto utilities
vi.mock('@/utils/crypto', () => ({
  encryptRequest: vi.fn((data) => data),
  shouldEncrypt: vi.fn(() => false),
}));

describe('Admin API - User Management', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockUsers: User[] = [
    {
      id: 1,
      name: 'Admin User',
      email: 'admin@gamelink.com',
      phone: '+1234567890',
      role: 'admin',
      status: 'active',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 2,
      name: 'Player User',
      email: 'player@test.com',
      phone: '+0987654321',
      role: 'player',
      status: 'active',
      createdAt: '2024-01-02T00:00:00Z',
    },
  ];

  describe('getUsers', () => {
    it('should fetch users without parameters', async () => {
      const axiosResponse: AxiosResponse<ApiResponse<User[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockUsers,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getUsers();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/users', { params: undefined });
      expect(result.data.data).toEqual(mockUsers);
    });

    it('should fetch users with query parameters', async () => {
      const params: UserQueryParams = {
        page: 1,
        page_size: 10,
        keyword: 'admin',
        role: ['admin'],
        status: ['active'],
      };

      const axiosResponse: AxiosResponse<ApiResponse<User[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockUsers,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      await adminApi.getUsers(params);

      expect(apiClient.get).toHaveBeenCalledWith('/admin/users', { params });
    });

    it('should handle 401 unauthorized error', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            message: 'Unauthorized',
          },
        },
      };

      vi.mocked(apiClient.get).mockRejectedValueOnce(error);

      await expect(adminApi.getUsers()).rejects.toMatchObject({
        response: { status: 401 },
      });
    });

    it('should handle 403 forbidden error', async () => {
      const error = {
        response: {
          status: 403,
          data: {
            success: false,
            message: 'Forbidden',
          },
        },
      };

      vi.mocked(apiClient.get).mockRejectedValueOnce(error);

      await expect(adminApi.getUsers()).rejects.toMatchObject({
        response: { status: 403 },
      });
    });
  });

  describe('getUser', () => {
    it('should fetch a single user by ID', async () => {
      const userId = 1;
      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockUsers[0],
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getUser(userId);

      expect(apiClient.get).toHaveBeenCalledWith(`/admin/users/${userId}`);
      expect(result.data.data).toEqual(mockUsers[0]);
    });

    it('should handle 404 not found error', async () => {
      const userId = 999;
      const error = {
        response: {
          status: 404,
          data: {
            success: false,
            message: 'User not found',
          },
        },
      };

      vi.mocked(apiClient.get).mockRejectedValueOnce(error);

      await expect(adminApi.getUser(userId)).rejects.toMatchObject({
        response: { status: 404 },
      });
    });
  });

  describe('createUser', () => {
    it('should create a new user successfully', async () => {
      const newUser: CreateUserDto = {
        name: 'New User',
        email: 'newuser@test.com',
        phone: '+1555555555',
        password: 'SecurePassword123!',
        role: 'user',
        status: 'active',
      };

      const createdUser: User = {
        id: 3,
        ...newUser,
        createdAt: '2024-01-03T00:00:00Z',
      };

      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 201,
          message: 'User created successfully',
          data: createdUser,
        },
        status: 201,
        statusText: 'Created',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.createUser(newUser);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/users', newUser);
      expect(result.data.data.id).toBe(3);
    });

    it('should handle validation error for invalid email', async () => {
      const invalidUser: CreateUserDto = {
        name: 'Test',
        email: 'invalid-email',
        phone: '+1555555555',
        password: 'password',
        role: 'user',
        status: 'active',
      };

      const error = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Invalid email format',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(adminApi.createUser(invalidUser)).rejects.toMatchObject({
        response: { status: 400 },
      });
    });

    it('should handle duplicate email error', async () => {
      const duplicateUser: CreateUserDto = {
        name: 'Duplicate',
        email: 'admin@gamelink.com',
        phone: '+1555555555',
        password: 'password',
        role: 'user',
        status: 'active',
      };

      const error = {
        response: {
          status: 409,
          data: {
            success: false,
            message: 'Email already exists',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(adminApi.createUser(duplicateUser)).rejects.toMatchObject({
        response: { status: 409 },
      });
    });
  });

  describe('updateUser', () => {
    it('should update an existing user', async () => {
      const userId = 1;
      const updateData: UpdateUserDto = {
        name: 'Updated Name',
        email: 'updated@test.com',
        phone: '+1999999999',
        avatarUrl: 'https://example.com/new-avatar.jpg',
        role: 'admin',
        status: 'active',
      };

      const updatedUser: User = {
        ...mockUsers[0],
        ...updateData,
      };

      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 200,
          message: 'User updated successfully',
          data: updatedUser,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.updateUser(userId, updateData);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/users/${userId}`, updateData);
      expect(result.data.data.name).toBe('Updated Name');
    });

    it('should handle update with password change', async () => {
      const userId = 1;
      const updateWithPassword: UpdateUserDto = {
        name: 'Admin User',
        email: 'admin@gamelink.com',
        phone: '+1234567890',
        role: 'admin',
        status: 'active',
        password: 'NewSecurePassword123!',
      };

      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 200,
          message: 'User updated successfully',
          data: mockUsers[0],
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updateUser(userId, updateWithPassword);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/users/${userId}`, updateWithPassword);
    });
  });

  describe('deleteUser', () => {
    it('should delete a user successfully', async () => {
      const userId = 2;

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'User deleted successfully',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.delete).mockResolvedValueOnce(axiosResponse);

      await adminApi.deleteUser(userId);

      expect(apiClient.delete).toHaveBeenCalledWith(`/admin/users/${userId}`);
    });

    it('should handle deletion of non-existent user', async () => {
      const userId = 999;
      const error = {
        response: {
          status: 404,
          data: {
            success: false,
            message: 'User not found',
          },
        },
      };

      vi.mocked(apiClient.delete).mockRejectedValueOnce(error);

      await expect(adminApi.deleteUser(userId)).rejects.toMatchObject({
        response: { status: 404 },
      });
    });
  });

  describe('updateUserStatus', () => {
    it('should update user status', async () => {
      const userId = 1;
      const newStatus = 'banned';

      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 200,
          message: 'Status updated',
          data: { ...mockUsers[0], status: newStatus as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updateUserStatus(userId, newStatus);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/users/${userId}/status`, {
        status: newStatus,
      });
    });
  });

  describe('updateUserRole', () => {
    it('should update user role', async () => {
      const userId = 2;
      const newRole = 'admin';

      const axiosResponse: AxiosResponse<ApiResponse<User>> = {
        data: {
          success: true,
          code: 200,
          message: 'Role updated',
          data: { ...mockUsers[1], role: newRole as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updateUserRole(userId, newRole);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/users/${userId}/role`, {
        role: newRole,
      });
    });
  });

  describe('Batch operations', () => {
    it('should batch update user roles', async () => {
      const data = {
        userIds: [1, 2, 3],
        role: 'player',
      };

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch update successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchUpdateUserRole(data);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/users/batch/role', data);
    });

    it('should batch update user status', async () => {
      const data = {
        userIds: [1, 2, 3],
        status: 'suspended',
      };

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch update successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchUpdateUserStatus(data);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/users/batch/status', data);
    });

    it('should batch delete users', async () => {
      const userIds = [1, 2, 3];

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch delete successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchDeleteUsers(userIds);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/users/batch/delete', {
        userIds,
      });
    });
  });
});

describe('Admin API - Game Management', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockGames: Game[] = [
    {
      id: 1,
      key: 'league-of-legends',
      name: 'League of Legends',
      category: 'MOBA',
      categoryId: 1,
      iconUrl: 'https://example.com/lol-icon.png',
      coverUrl: 'https://example.com/lol-cover.png',
      description: 'Popular MOBA game',
      isActive: true,
      sortOrder: 1,
    },
    {
      id: 2,
      key: 'valorant',
      name: 'Valorant',
      category: 'FPS',
      categoryId: 2,
      iconUrl: 'https://example.com/val-icon.png',
      coverUrl: 'https://example.com/val-cover.png',
      description: 'Tactical shooter',
      isActive: true,
      sortOrder: 2,
    },
  ];

  describe('getGames', () => {
    it('should fetch games without parameters', async () => {
      const axiosResponse: AxiosResponse<ApiResponse<Game[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockGames,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getGames();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/games', { params: undefined });
      expect(result.data.data).toEqual(mockGames);
    });

    it('should fetch games with status filter', async () => {
      const params = { status: 'active', page_size: 100 };

      const axiosResponse: AxiosResponse<ApiResponse<Game[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockGames.filter(g => g.isActive),
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      await adminApi.getGames(params);

      expect(apiClient.get).toHaveBeenCalledWith('/admin/games', { params });
    });
  });

  describe('createGame', () => {
    it('should create a new game successfully', async () => {
      const newGame: CreateGameDto = {
        key: 'apex-legends',
        name: 'Apex Legends',
        category: 'Battle Royale',
        categoryId: 3,
        iconUrl: 'https://example.com/apex-icon.png',
        description: 'Battle royale game',
        isActive: true,
        sortOrder: 3,
      };

      const createdGame: Game = {
        id: 3,
        ...newGame,
      };

      const axiosResponse: AxiosResponse<ApiResponse<Game>> = {
        data: {
          success: true,
          code: 201,
          message: 'Game created successfully',
          data: createdGame,
        },
        status: 201,
        statusText: 'Created',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.createGame(newGame);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/games', newGame);
      expect(result.data.data.id).toBe(3);
    });

    it('should handle duplicate game key error', async () => {
      const duplicateGame: CreateGameDto = {
        key: 'league-of-legends',
        name: 'Duplicate LoL',
        category: 'MOBA',
      };

      const error = {
        response: {
          status: 409,
          data: {
            success: false,
            message: 'Game key already exists',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(adminApi.createGame(duplicateGame)).rejects.toMatchObject({
        response: { status: 409 },
      });
    });
  });

  describe('updateGame', () => {
    it('should update an existing game', async () => {
      const gameId = 1;
      const updateData: UpdateGameDto = {
        name: 'League of Legends (Updated)',
        description: 'Updated description',
        isActive: false,
      };

      const axiosResponse: AxiosResponse<ApiResponse<Game>> = {
        data: {
          success: true,
          code: 200,
          message: 'Game updated successfully',
          data: { ...mockGames[0], ...updateData },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updateGame(gameId, updateData);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/games/${gameId}`, updateData);
    });
  });

  describe('deleteGame', () => {
    it('should delete a game successfully', async () => {
      const gameId = 1;

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Game deleted successfully',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.delete).mockResolvedValueOnce(axiosResponse);

      await adminApi.deleteGame(gameId);

      expect(apiClient.delete).toHaveBeenCalledWith(`/admin/games/${gameId}`);
    });

    it('should handle deletion of game with active orders', async () => {
      const gameId = 1;
      const error = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Cannot delete game with active orders',
          },
        },
      };

      vi.mocked(apiClient.delete).mockRejectedValueOnce(error);

      await expect(adminApi.deleteGame(gameId)).rejects.toMatchObject({
        response: { status: 400 },
      });
    });
  });

  describe('batchDeleteGames', () => {
    it('should batch delete games', async () => {
      const gameIds = [1, 2];

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch delete successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchDeleteGames(gameIds);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/games/batch/delete', {
        gameIds,
      });
    });
  });
});

describe('Admin API - Player Management', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockPlayers: Player[] = [
    {
      id: 1,
      userId: 2,
      nickname: 'ProGamer123',
      bio: 'Professional player',
      rank: 'diamond',
      hourlyRateCents: 5000,
      mainGameId: 1,
      verificationStatus: 'verified',
      ratingAverage: 4.8,
      ratingCount: 120,
      skillTags: ['support', 'carry'],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  describe('getPlayers', () => {
    it('should fetch players without parameters', async () => {
      const axiosResponse: AxiosResponse<ApiResponse<Player[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockPlayers,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getPlayers();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/players', { params: undefined });
      expect(result.data.data).toEqual(mockPlayers);
    });

    it('should fetch players with status filter', async () => {
      const params = { status: 'verified', page: 1, page_size: 10 };

      const axiosResponse: AxiosResponse<ApiResponse<Player[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockPlayers.filter(p => p.verificationStatus === 'verified'),
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      await adminApi.getPlayers(params);

      expect(apiClient.get).toHaveBeenCalledWith('/admin/players', { params });
    });
  });

  describe('updatePlayerVerification', () => {
    it('should approve player verification', async () => {
      const playerId = 1;
      const status = 'verified';
      const remark = 'Player verified successfully';

      const axiosResponse: AxiosResponse<ApiResponse<Player>> = {
        data: {
          success: true,
          code: 200,
          message: 'Verification updated',
          data: { ...mockPlayers[0], verificationStatus: status as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updatePlayerVerification(playerId, status, remark);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/players/${playerId}/verification`, {
        verification_status: status,
        remark,
      });
    });

    it('should reject player verification', async () => {
      const playerId = 1;
      const status = 'rejected';
      const remark = 'Insufficient documentation';

      const axiosResponse: AxiosResponse<ApiResponse<Player>> = {
        data: {
          success: true,
          code: 200,
          message: 'Verification rejected',
          data: { ...mockPlayers[0], verificationStatus: status as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updatePlayerVerification(playerId, status, remark);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/players/${playerId}/verification`, {
        verification_status: status,
        remark,
      });
    });
  });

  describe('updatePlayerSkillTags', () => {
    it('should update player skill tags', async () => {
      const playerId = 1;
      const tags = ['support', 'jungle', 'mid'];

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Skill tags updated',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.put).mockResolvedValueOnce(axiosResponse);

      await adminApi.updatePlayerSkillTags(playerId, tags);

      expect(apiClient.put).toHaveBeenCalledWith(`/admin/players/${playerId}/skill-tags`, {
        tags,
      });
    });
  });
});

describe('Admin API - Order Management', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockOrders: Order[] = [
    {
      id: 1,
      orderNo: 'ORD20240101001',
      userId: 1,
      playerId: 2,
      gameId: 1,
      title: 'Coaching Session',
      description: '1-hour coaching',
      totalPriceCents: 5000,
      currency: 'CNY',
      status: 'pending',
      scheduledStart: '2024-01-02T10:00:00Z',
      scheduledEnd: '2024-01-02T11:00:00Z',
      completedAt: '',
      cancelReason: '',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  describe('getOrders', () => {
    it('should fetch orders without parameters', async () => {
      const axiosResponse: AxiosResponse<ApiResponse<Order[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockOrders,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getOrders();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/orders', { params: undefined });
      expect(result.data.data).toEqual(mockOrders);
    });

    it('should fetch orders with filters', async () => {
      const params: OrderQueryParams = {
        page: 1,
        page_size: 20,
        status: 'pending',
        userId: 1,
        orderNumber: 'ORD20240101001',
      };

      const axiosResponse: AxiosResponse<ApiResponse<Order[]>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockOrders,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      await adminApi.getOrders(params);

      expect(apiClient.get).toHaveBeenCalledWith('/admin/orders', { params });
    });
  });

  describe('cancelOrder', () => {
    it('should cancel an order successfully', async () => {
      const orderId = 1;
      const note = 'Customer requested cancellation';

      const axiosResponse: AxiosResponse<ApiResponse<Order>> = {
        data: {
          success: true,
          code: 200,
          message: 'Order cancelled',
          data: { ...mockOrders[0], status: 'cancelled' as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.cancelOrder(orderId, note);

      expect(apiClient.post).toHaveBeenCalledWith(`/admin/orders/${orderId}/cancel`, {
        note,
      });
    });

    it('should handle cancellation of completed order', async () => {
      const orderId = 1;
      const error = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Cannot cancel completed order',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(adminApi.cancelOrder(orderId)).rejects.toMatchObject({
        response: { status: 400 },
      });
    });
  });

  describe('refundOrder', () => {
    it('should process refund successfully', async () => {
      const orderId = 1;
      const refundData = {
        reason: 'Service not provided',
        amount_cents: 5000,
        note: 'Full refund approved',
      };

      const axiosResponse: AxiosResponse<ApiResponse<Order>> = {
        data: {
          success: true,
          code: 200,
          message: 'Refund processed',
          data: { ...mockOrders[0], status: 'refunded' as any },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.refundOrder(orderId, refundData);

      expect(apiClient.post).toHaveBeenCalledWith(
        `/admin/orders/${orderId}/refund`,
        refundData
      );
    });

    it('should handle partial refund', async () => {
      const orderId = 1;
      const partialRefundData = {
        reason: 'Partial service',
        amount_cents: 2500,
        note: '50% refund',
      };

      const axiosResponse: AxiosResponse<ApiResponse<Order>> = {
        data: {
          success: true,
          code: 200,
          message: 'Partial refund processed',
          data: mockOrders[0],
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.refundOrder(orderId, partialRefundData);

      expect(partialRefundData.amount_cents).toBeLessThan(mockOrders[0].totalPriceCents);
    });
  });

  describe('Batch operations', () => {
    it('should batch cancel orders', async () => {
      const orderIds = [1, 2, 3];
      const reason = 'System maintenance';

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch cancel successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchCancelOrders(orderIds, reason);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/orders/batch/cancel', {
        orderIds,
        reason,
      });
    });

    it('should batch complete orders', async () => {
      const orderIds = [1, 2, 3];

      const axiosResponse: AxiosResponse<ApiResponse<void>> = {
        data: {
          success: true,
          code: 200,
          message: 'Batch complete successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await adminApi.batchCompleteOrders(orderIds);

      expect(apiClient.post).toHaveBeenCalledWith('/admin/orders/batch/complete', {
        orderIds,
      });
    });
  });
});

describe('Admin API - Dashboard & Statistics', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getDashboardStats', () => {
    it('should fetch dashboard statistics', async () => {
      const mockStats = {
        totalUsers: 1000,
        totalPlayers: 150,
        totalGames: 50,
        totalOrders: 500,
        ordersByStatus: {
          pending: 50,
          confirmed: 100,
          in_progress: 75,
          completed: 250,
          cancelled: 25,
        },
        paymentsByStatus: {
          pending: 30,
          completed: 450,
          failed: 20,
        },
        totalPaidAmountCents: 2500000,
      };

      const axiosResponse: AxiosResponse<ApiResponse<typeof mockStats>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockStats,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getDashboardStats();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/stats/dashboard');
      expect(result.data.data.totalUsers).toBe(1000);
    });
  });

  describe('getUserStats', () => {
    it('should fetch user statistics', async () => {
      const mockUserStats = {
        total: 1000,
        byRole: {
          user: 800,
          player: 150,
          admin: 50,
        },
        byStatus: {
          active: 950,
          banned: 30,
          suspended: 20,
        },
        recentRegistrations: 50,
      };

      const axiosResponse: AxiosResponse<ApiResponse<typeof mockUserStats>> = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockUserStats,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await adminApi.getUserStats();

      expect(apiClient.get).toHaveBeenCalledWith('/admin/users/stats');
      expect(result.data.data.byRole.player).toBe(150);
    });
  });
});

describe('Admin API - Error Handling', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should handle network errors', async () => {
    const networkError = new Error('Network Error');
    vi.mocked(apiClient.get).mockRejectedValueOnce(networkError);

    await expect(adminApi.getUsers()).rejects.toThrow('Network Error');
  });

  it('should handle timeout errors', async () => {
    const timeoutError = new Error('timeout of 10000ms exceeded');
    vi.mocked(apiClient.get).mockRejectedValueOnce(timeoutError);

    await expect(adminApi.getUsers()).rejects.toThrow();
  });

  it('should handle 500 internal server error', async () => {
    const error = {
      response: {
        status: 500,
        data: {
          success: false,
          message: 'Internal server error',
        },
      },
    };

    vi.mocked(apiClient.get).mockRejectedValueOnce(error);

    await expect(adminApi.getUsers()).rejects.toMatchObject({
      response: { status: 500 },
    });
  });
});
