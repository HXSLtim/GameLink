/**
 * User Store - 用户数据管理
 *
 * 功能：
 * - 用户列表缓存（减少重复 API 调用）
 * - CRUD 操作（创建、更新、删除用户）
 * - 分页支持
 * - 筛选功能（按状态、角色、关键词筛选）
 * - 与 authStore 联动（登出时清理数据）
 * - 使用 UserService 进行业务逻辑处理
 *
 * @module userStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User as ApiUser, CreateUserDto, UpdateUserDto } from '@/api/admin';
import {
  userService as defaultUserService,
  type IUserService,
  type UserValidationResult,
  type UserExportData,
} from '@/services/domain';
import type { ServiceResult, BatchResult } from '@/services/utils';

/**
 * 分页参数接口
 */
interface Pagination {
  current: number;
  pageSize: number;
  total: number;
}

/**
 * 筛选条件接口
 */
interface UserFilters {
  status?: string;
  role?: string;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}

/**
 * Store 错误信息接口
 */
interface StoreError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

/**
 * 用户状态接口
 */
interface UserState {
  // ========== State ==========
  /** 用户列表缓存 */
  users: ApiUser[];
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: StoreError | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: UserFilters;
  /** 最后一次批量操作结果 */
  lastBatchResult: BatchResult<void> | null;

  // ========== Actions ==========

  /**
   * 获取用户列表
   * @param page - 页码（默认 1）
   * @param pageSize - 每页数量（默认 10）
   */
  fetchUsers: (page?: number, pageSize?: number) => Promise<void>;

  /**
   * 创建用户
   * @param userData - 用户数据
   * @throws 创建失败时抛出错误
   */
  createUser: (userData: CreateUserDto) => Promise<ServiceResult<ApiUser>>;

  /**
   * 更新用户
   * @param id - 用户 ID
   * @param userData - 更新的用户数据
   * @throws 更新失败时抛出错误
   */
  updateUser: (id: number, userData: UpdateUserDto) => Promise<ServiceResult<ApiUser>>;

  /**
   * 删除用户
   * @param id - 用户 ID
   * @throws 删除失败时抛出错误
   */
  deleteUser: (id: number) => Promise<ServiceResult<void>>;

  /**
   * 批量删除用户
   * @param userIds - 用户 ID 数组
   * @returns 批量操作结果
   */
  batchDeleteUsers: (userIds: number[]) => Promise<BatchResult<void>>;

  /**
   * 更新用户状态
   * @param id - 用户 ID
   * @param status - 状态（active/banned/suspended）
   * @throws 更新失败时抛出错误
   */
  updateUserStatus: (id: number, status: string) => Promise<ServiceResult<ApiUser>>;

  /**
   * 批量更新用户状态
   * @param userIds - 用户 ID 数组
   * @param status - 状态
   * @returns 批量操作结果
   */
  batchUpdateUserStatus: (userIds: number[], status: string) => Promise<BatchResult<void>>;

  /**
   * 更新用户角色
   * @param id - 用户 ID
   * @param role - 角色（user/player/admin）
   * @throws 更新失败时抛出错误
   */
  updateUserRole: (id: number, role: string) => Promise<ServiceResult<ApiUser>>;

  /**
   * 批量更新用户角色
   * @param userIds - 用户 ID 数组
   * @param role - 角色
   * @returns 批量操作结果
   */
  batchUpdateUserRole: (userIds: number[], role: string) => Promise<BatchResult<void>>;

  /**
   * 验证用户数据
   * @param data - 用户数据
   * @returns 验证结果
   */
  validateUserData: (data: Partial<CreateUserDto>) => UserValidationResult;

  /**
   * 导出用户数据
   * @param users - 用户列表（可选，默认使用当前缓存的用户）
   * @returns 导出数据
   */
  exportUsers: (users?: ApiUser[]) => UserExportData;

  /**
   * 设置筛选条件
   * @param newFilters - 新的筛选条件（部分更新）
   */
  setFilters: (newFilters: Partial<UserFilters>) => void;

  /**
   * 清除所有筛选条件
   */
  clearFilters: () => void;

  /**
   * 重置 store 到初始状态
   */
  reset: () => void;

  /**
   * 清除错误状态
   */
  clearError: () => void;

  // ========== Selectors ==========

  /**
   * 根据 ID 获取用户
   * @param id - 用户 ID
   * @returns 用户对象或 undefined
   */
  getUserById: (id: number) => ApiUser | undefined;

  /**
   * 获取活跃用户列表
   * @returns 状态为 active 的用户数组
   */
  getActiveUsers: () => ApiUser[];

  /**
   * 根据角色筛选用户
   * @param role - 角色名称
   * @returns 指定角色的用户数组
   */
  getUsersByRole: (role: string) => ApiUser[];

  /**
   * 根据状态筛选用户
   * @param status - 状态名称
   * @returns 指定状态的用户数组
   */
  getUsersByStatus: (status: string) => ApiUser[];
}

/**
 * 创建 UserStore 的工厂函数
 * 支持依赖注入，便于测试
 */
export const createUserStore = (userService: IUserService = defaultUserService) => {
  return create<UserState>()(
    persist(
      (set, get) => ({
        // ========== Initial State ==========
        users: [],
        loading: false,
        error: null,
        pagination: {
          current: 1,
          pageSize: 10,
          total: 0,
        },
        filters: {},
        lastBatchResult: null,

        // ========== Actions ==========

        /**
         * 获取用户列表
         */
        fetchUsers: async (page = 1, pageSize = 10) => {
          set({ loading: true, error: null });

          const { filters } = get();

          // 使用 UserService 获取用户
          const result = await userService.getUsers({
            page,
            page_size: pageSize,
            keyword: filters.keyword,
            role: filters.role ? [filters.role] : undefined,
            status: filters.status ? [filters.status] : undefined,
            date_from: filters.date_from,
            date_to: filters.date_to,
          });

          if (result.success && result.data) {
            set({
              users: result.data,
              pagination: {
                current: page,
                pageSize,
                total: result.data.length, // Note: API should return total in pagination
              },
              loading: false,
            });
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '获取用户列表失败' };
            set({
              error: {
                code: error.code,
                message: error.message,
                details: error.details,
              },
              loading: false,
            });
            throw new Error(error.message);
          }
        },

        /**
         * 创建用户
         */
        createUser: async (userData: CreateUserDto) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 创建用户
            const result = await userService.createUser(userData);

            if (result.success && result.data) {
              // 将新用户添加到列表开头
              set((state) => ({
                users: [result.data!, ...state.users],
                pagination: {
                  ...state.pagination,
                  total: state.pagination.total + 1,
                },
                loading: false,
              }));
            } else {
              const error = result.error || { code: 'UNKNOWN_ERROR', message: '创建用户失败' };
              set({
                error: {
                  code: error.code,
                  message: error.message,
                  details: error.details,
                },
                loading: false,
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '创建用户失败';
            set({
              error: { code: 'CREATE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 更新用户
         */
        updateUser: async (id: number, userData: UpdateUserDto) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 更新用户
            const result = await userService.updateUser(id, userData);

            if (result.success && result.data) {
              // 更新本地缓存中的用户
              set((state) => ({
                users: state.users.map((u) =>
                  u.id === id ? result.data! : u
                ),
                loading: false,
              }));
            } else {
              const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新用户失败' };
              set({
                error: {
                  code: error.code,
                  message: error.message,
                  details: error.details,
                },
                loading: false,
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '更新用户失败';
            set({
              error: { code: 'UPDATE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 删除用户
         */
        deleteUser: async (id: number) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 删除用户
            const result = await userService.deleteUser(id);

            if (result.success) {
              // 从本地缓存中移除
              set((state) => ({
                users: state.users.filter((u) => u.id !== id),
                pagination: {
                  ...state.pagination,
                  total: Math.max(0, state.pagination.total - 1),
                },
                loading: false,
              }));
            } else {
              const error = result.error || { code: 'UNKNOWN_ERROR', message: '删除用户失败' };
              set({
                error: {
                  code: error.code,
                  message: error.message,
                  details: error.details,
                },
                loading: false,
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '删除用户失败';
            set({
              error: { code: 'DELETE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 批量删除用户
         */
        batchDeleteUsers: async (userIds: number[]) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 批量删除
            const result = await userService.batchDelete(userIds);

            // 获取成功删除的用户 ID
            const successfulIds = result.results
              .filter((r) => r.success)
              .map((r) => userIds[r.index]);

            // 从本地缓存中移除成功删除的用户
            set((state) => ({
              users: state.users.filter((u) => !successfulIds.includes(u.id)),
              pagination: {
                ...state.pagination,
                total: Math.max(0, state.pagination.total - successfulIds.length),
              },
              loading: false,
              lastBatchResult: result,
            }));

            // 如果有失败的操作，设置错误信息
            if (!result.success) {
              const failedCount = result.failed;
              set({
                error: {
                  code: 'BATCH_PARTIAL_FAILURE',
                  message: `批量删除部分失败：${failedCount} 个用户删除失败`,
                  details: { failedCount, results: result.results },
                },
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '批量删除用户失败';
            set({
              error: { code: 'BATCH_DELETE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 更新用户状态
         */
        updateUserStatus: async (id: number, status: string) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 更新状态
            const result = await userService.updateUserStatus(id, status);

            if (result.success && result.data) {
              // 更新本地缓存
              set((state) => ({
                users: state.users.map((u) =>
                  u.id === id ? { ...u, status: status as ApiUser['status'] } : u
                ),
                loading: false,
              }));
            } else {
              const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新用户状态失败' };
              set({
                error: {
                  code: error.code,
                  message: error.message,
                  details: error.details,
                },
                loading: false,
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '更新用户状态失败';
            set({
              error: { code: 'UPDATE_STATUS_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 批量更新用户状态
         */
        batchUpdateUserStatus: async (userIds: number[], status: string) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 批量更新状态
            const result = await userService.batchUpdateStatus(userIds, status);

            // 获取成功更新的用户 ID
            const successfulIds = result.results
              .filter((r) => r.success)
              .map((r) => userIds[r.index]);

            // 更新本地缓存
            set((state) => ({
              users: state.users.map((u) =>
                successfulIds.includes(u.id)
                  ? { ...u, status: status as ApiUser['status'] }
                  : u
              ),
              loading: false,
              lastBatchResult: result,
            }));

            // 如果有失败的操作，设置错误信息
            if (!result.success) {
              const failedCount = result.failed;
              set({
                error: {
                  code: 'BATCH_PARTIAL_FAILURE',
                  message: `批量更新状态部分失败：${failedCount} 个用户更新失败`,
                  details: { failedCount, results: result.results },
                },
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '批量更新用户状态失败';
            set({
              error: { code: 'BATCH_UPDATE_STATUS_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 更新用户角色
         */
        updateUserRole: async (id: number, role: string) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 更新角色
            const result = await userService.updateUserRole(id, role);

            if (result.success && result.data) {
              // 更新本地缓存
              set((state) => ({
                users: state.users.map((u) =>
                  u.id === id ? { ...u, role: role as 'user' | 'player' | 'admin' } : u
                ),
                loading: false,
              }));
            } else {
              const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新用户角色失败' };
              set({
                error: {
                  code: error.code,
                  message: error.message,
                  details: error.details,
                },
                loading: false,
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '更新用户角色失败';
            set({
              error: { code: 'UPDATE_ROLE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 批量更新用户角色
         */
        batchUpdateUserRole: async (userIds: number[], role: string) => {
          set({ loading: true, error: null });

          try {
            // 使用 UserService 批量更新角色
            const result = await userService.batchUpdateRole(userIds, role);

            // 获取成功更新的用户 ID
            const successfulIds = result.results
              .filter((r) => r.success)
              .map((r) => userIds[r.index]);

            // 更新本地缓存
            set((state) => ({
              users: state.users.map((u) =>
                successfulIds.includes(u.id)
                  ? { ...u, role: role as 'user' | 'player' | 'admin' }
                  : u
              ),
              loading: false,
              lastBatchResult: result,
            }));

            // 如果有失败的操作，设置错误信息
            if (!result.success) {
              const failedCount = result.failed;
              set({
                error: {
                  code: 'BATCH_PARTIAL_FAILURE',
                  message: `批量更新角色部分失败：${failedCount} 个用户更新失败`,
                  details: { failedCount, results: result.results },
                },
              });
            }

            return result;
          } catch (error: unknown) {
            const errorMessage =
              error instanceof Error ? error.message : '批量更新用户角色失败';
            set({
              error: { code: 'BATCH_UPDATE_ROLE_ERROR', message: errorMessage },
              loading: false,
            });
            throw error;
          }
        },

        /**
         * 验证用户数据
         */
        validateUserData: (data: Partial<CreateUserDto>) => {
          return userService.validateUserData(data);
        },

        /**
         * 导出用户数据
         */
        exportUsers: (users?: ApiUser[]) => {
          const usersToExport = users || get().users;
          return userService.exportUsers(usersToExport);
        },

        /**
         * 设置筛选条件
         */
        setFilters: (newFilters: Partial<UserFilters>) => {
          set((state) => ({
            filters: { ...state.filters, ...newFilters },
            // 筛选条件变化时重置到第一页
            pagination: { ...state.pagination, current: 1 },
          }));
        },

        /**
         * 清除筛选条件
         */
        clearFilters: () => {
          set({
            filters: {},
            pagination: { current: 1, pageSize: 10, total: 0 },
          });
        },

        /**
         * 重置 store
         */
        reset: () => {
          set({
            users: [],
            loading: false,
            error: null,
            pagination: { current: 1, pageSize: 10, total: 0 },
            filters: {},
            lastBatchResult: null,
          });
        },

        /**
         * 清除错误状态
         */
        clearError: () => {
          set({ error: null });
        },

        // ========== Selectors ==========

        /**
         * 根据 ID 获取用户
         */
        getUserById: (id: number) => {
          return get().users.find((u) => u.id === id);
        },

        /**
         * 获取活跃用户
         */
        getActiveUsers: () => {
          return get().users.filter((u) => u.status === 'active');
        },

        /**
         * 根据角色筛选用户
         */
        getUsersByRole: (role: string) => {
          return get().users.filter((u) => u.role === role);
        },

        /**
         * 根据状态筛选用户
         */
        getUsersByStatus: (status: string) => {
          return get().users.filter((u) => u.status === status);
        },
      }),
      {
        name: 'user-storage',
        // 只缓存部分数据，避免 localStorage 过大
        partialize: (state) => ({
          // 最多缓存最近 50 条用户数据
          users: state.users.slice(0, 50),
          // 保留筛选条件
          filters: state.filters,
        }),
      }
    )
  );
};

/**
 * 默认 UserStore 实例
 */
export const useUserStore = createUserStore();

/**
 * 与 authStore 联动的辅助函数
 * 在 authStore 的 logout 方法中调用此函数清理用户数据
 */
export const clearUserStore = () => {
  useUserStore.getState().reset();
};

/**
 * 导出便捷的 hooks
 */
export const useUsers = () => useUserStore((state) => state.users);
export const useUserLoading = () => useUserStore((state) => state.loading);
export const useUserError = () => useUserStore((state) => state.error);
export const useUserPagination = () => useUserStore((state) => state.pagination);
export const useUserFilters = () => useUserStore((state) => state.filters);
export const useLastBatchResult = () => useUserStore((state) => state.lastBatchResult);

// Re-export types for convenience
export type { UserValidationResult, UserExportData };
