/**
 * User Store - 用户数据管理
 *
 * 功能：
 * - 用户列表缓存（减少重复 API 调用）
 * - CRUD 操作（创建、更新、删除用户）
 * - 分页支持
 * - 筛选功能（按状态、角色、关键词筛选）
 * - 与 authStore 联动（登出时清理数据）
 *
 * @module userStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { adminApi, type User as ApiUser } from '@/api/admin';

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
 * 用户状态接口
 */
interface UserState {
  // ========== State ==========
  /** 用户列表缓存 */
  users: ApiUser[];
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: string | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: UserFilters;

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
  createUser: (userData: Partial<ApiUser>) => Promise<void>;

  /**
   * 更新用户
   * @param id - 用户 ID
   * @param userData - 更新的用户数据
   * @throws 更新失败时抛出错误
   */
  updateUser: (id: number, userData: Partial<ApiUser>) => Promise<void>;

  /**
   * 删除用户
   * @param id - 用户 ID
   * @throws 删除失败时抛出错误
   */
  deleteUser: (id: number) => Promise<void>;

  /**
   * 批量删除用户
   * @param userIds - 用户 ID 数组
   * @throws 删除失败时抛出错误
   */
  batchDeleteUsers: (userIds: number[]) => Promise<void>;

  /**
   * 更新用户状态
   * @param id - 用户 ID
   * @param status - 状态（active/banned/suspended）
   * @throws 更新失败时抛出错误
   */
  updateUserStatus: (id: number, status: string) => Promise<void>;

  /**
   * 批量更新用户状态
   * @param userIds - 用户 ID 数组
   * @param status - 状态
   * @throws 更新失败时抛出错误
   */
  batchUpdateUserStatus: (userIds: number[], status: string) => Promise<void>;

  /**
   * 更新用户角色
   * @param id - 用户 ID
   * @param role - 角色（user/player/admin）
   * @throws 更新失败时抛出错误
   */
  updateUserRole: (id: number, role: string) => Promise<void>;

  /**
   * 批量更新用户角色
   * @param userIds - 用户 ID 数组
   * @param role - 角色
   * @throws 更新失败时抛出错误
   */
  batchUpdateUserRole: (userIds: number[], role: string) => Promise<void>;

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
 * 用户 Store 实现
 */
export const useUserStore = create<UserState>()(
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

      // ========== Actions ==========

      /**
       * 获取用户列表
       */
      fetchUsers: async (page = 1, pageSize = 10) => {
        set({ loading: true, error: null });

        try {
          const { filters } = get();

          // 调用 API
          const response = await adminApi.getUsers({
            page,
            page_size: pageSize,
            keyword: filters.keyword,
            role: filters.role ? [filters.role] : undefined,
            status: filters.status ? [filters.status] : undefined,
            date_from: filters.date_from,
            date_to: filters.date_to,
          });

          // 更新状态
          set({
            users: response.data.data || [],
            pagination: {
              current: page,
              pageSize,
              total: response.data.pagination?.total || 0,
            },
            loading: false,
          });
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '获取用户列表失败';
          set({
            error: errorMessage,
            loading: false,
          });
          // 向上层抛出错误，让调用方处理
          throw error;
        }
      },

      /**
       * 创建用户
       */
      createUser: async (userData: Partial<ApiUser>) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 创建用户
          const response = await adminApi.createUser({
            name: userData.name || '',
            email: userData.email || '',
            phone: userData.phone || '',
            password: '', // 需要调用方提供
            avatarUrl: userData.avatarUrl,
            role: userData.role as 'user' | 'player' | 'admin' || 'user',
            status: userData.status as 'active' | 'banned' | 'suspended' || 'active',
          });

          // 将新用户添加到列表开头
          set((state) => ({
            users: [response.data.data, ...state.users],
            pagination: {
              ...state.pagination,
              total: state.pagination.total + 1,
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '创建用户失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 更新用户
       */
      updateUser: async (id: number, userData: Partial<ApiUser>) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新用户
          const response = await adminApi.updateUser(id, {
            name: userData.name || '',
            email: userData.email || '',
            phone: userData.phone || '',
            avatarUrl: userData.avatarUrl,
            role: userData.role as 'user' | 'player' | 'admin',
            status: userData.status as 'active' | 'banned' | 'suspended',
          });

          // 更新本地缓存中的用户
          set((state) => ({
            users: state.users.map((u) =>
              u.id === id ? response.data.data : u
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新用户失败';
          set({
            error: errorMessage,
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
          // 调用 API 删除用户
          await adminApi.deleteUser(id);

          // 从本地缓存中移除
          set((state) => ({
            users: state.users.filter((u) => u.id !== id),
            pagination: {
              ...state.pagination,
              total: Math.max(0, state.pagination.total - 1),
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '删除用户失败';
          set({
            error: errorMessage,
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
          // 调用 API 批量删除
          await adminApi.batchDeleteUsers(userIds);

          // 从本地缓存中移除
          set((state) => ({
            users: state.users.filter((u) => !userIds.includes(u.id)),
            pagination: {
              ...state.pagination,
              total: Math.max(0, state.pagination.total - userIds.length),
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量删除用户失败';
          set({
            error: errorMessage,
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
          // 调用 API 更新状态
          await adminApi.updateUserStatus(id, status);

          // 更新本地缓存
          set((state) => ({
            users: state.users.map((u) =>
              u.id === id ? { ...u, status: status as ApiUser['status'] } : u
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新用户状态失败';
          set({
            error: errorMessage,
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
          // 调用 API 批量更新状态
          await adminApi.batchUpdateUserStatus({ userIds, status });

          // 更新本地缓存
          set((state) => ({
            users: state.users.map((u) =>
              userIds.includes(u.id)
                ? { ...u, status: status as ApiUser['status'] }
                : u
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量更新用户状态失败';
          set({
            error: errorMessage,
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
          // 调用 API 更新角色
          await adminApi.updateUserRole(id, role);

          // 更新本地缓存
          set((state) => ({
            users: state.users.map((u) =>
              u.id === id ? { ...u, role } : u
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新用户角色失败';
          set({
            error: errorMessage,
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
          // 调用 API 批量更新角色
          await adminApi.batchUpdateUserRole({ userIds, role });

          // 更新本地缓存
          set((state) => ({
            users: state.users.map((u) =>
              userIds.includes(u.id) ? { ...u, role } : u
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量更新用户角色失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
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
        });
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
