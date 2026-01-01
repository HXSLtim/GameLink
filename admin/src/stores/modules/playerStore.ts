/**
 * Player Store - 陪玩师数据管理
 *
 * 功能：
 * - 陪玩师列表缓存（减少重复 API 调用）
 * - CRUD 操作（创建、更新、删除陪玩师）
 * - 分页支持
 * - 筛选功能（按状态、关键词筛选）
 * - 状态管理（available/busy/offline）
 * - 价格调整
 * - 等级管理
 * - 审核状态管理（pending/verified/rejected）
 * - 批量操作支持
 * - 与 authStore 联动（登出时清理数据）
 *
 * @module playerStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { adminApi, type Player as ApiPlayer, type CreatePlayerDto, type UpdatePlayerDto } from '../../api/admin';

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
interface PlayerFilters {
  status?: string; // 审核状态: pending/verified/rejected
  keyword?: string; // 昵称搜索
  verification_status?: 'pending' | 'verified' | 'rejected';
}

/**
 * 陪玩师状态接口
 */
interface PlayerState {
  // ========== State ==========
  /** 陪玩师列表缓存 */
  players: ApiPlayer[];
  /** 当前选中的陪玩师详情 */
  currentPlayer: ApiPlayer | null;
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: string | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: PlayerFilters;

  // ========== Actions ==========

  /**
   * 获取陪玩师列表
   * @param page - 页码（默认 1）
   * @param pageSize - 每页数量（默认 10）
   */
  fetchPlayers: (page?: number, pageSize?: number) => Promise<void>;

  /**
   * 获取单个陪玩师详情
   * @param id - 陪玩师 ID
   */
  fetchPlayer: (id: number) => Promise<void>;

  /**
   * 创建陪玩师
   * @param playerData - 陪玩师数据
   * @throws 创建失败时抛出错误
   */
  createPlayer: (playerData: CreatePlayerDto) => Promise<void>;

  /**
   * 更新陪玩师信息
   * @param id - 陪玩师 ID
   * @param playerData - 更新的陪玩师数据
   * @throws 更新失败时抛出错误
   */
  updatePlayer: (id: number, playerData: UpdatePlayerDto) => Promise<void>;

  /**
   * 删除陪玩师
   * @param id - 陪玩师 ID
   * @throws 删除失败时抛出错误
   */
  deletePlayer: (id: number) => Promise<void>;

  /**
   * 批量删除陪玩师
   * @param playerIds - 陪玩师 ID 数组
   * @throws 删除失败时抛出错误
   */
  batchDeletePlayers: (playerIds: number[]) => Promise<void>;

  /**
   * 更新陪玩师审核状态
   * @param id - 陪玩师 ID
   * @param status - 审核状态（pending/verified/rejected）
   * @param remark - 备注/拒绝原因
   * @throws 更新失败时抛出错误
   */
  updateVerificationStatus: (
    id: number,
    status: 'pending' | 'verified' | 'rejected',
    remark?: string
  ) => Promise<void>;

  /**
   * 批量更新陪玩师状态
   * @param playerIds - 陪玩师 ID 数组
   * @param status - 状态
   * @throws 更新失败时抛出错误
   */
  batchUpdatePlayerStatus: (playerIds: number[], status: string) => Promise<void>;

  /**
   * 更新陪玩师技能标签
   * @param id - 陪玩师 ID
   * @param tags - 技能标签数组
   * @throws 更新失败时抛出错误
   */
  updateSkillTags: (id: number, tags: string[]) => Promise<void>;

  /**
   * 调整陪玩师时薪价格
   * @param id - 陪玩师 ID
   * @param hourlyRateCents - 新价格（分）
   * @throws 更新失败时抛出错误
   */
  updatePrice: (id: number, hourlyRateCents: number) => Promise<void>;

  /**
   * 更新陪玩师等级
   * @param id - 陪玩师 ID
   * @param rank - 等级
   * @throws 更新失败时抛出错误
   */
  updateRank: (id: number, rank: string) => Promise<void>;

  /**
   * 设置筛选条件
   * @param newFilters - 新的筛选条件（部分更新）
   */
  setFilters: (newFilters: Partial<PlayerFilters>) => void;

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
   * 根据 ID 获取陪玩师
   * @param id - 陪玩师 ID
   * @returns 陪玩师对象或 undefined
   */
  getPlayerById: (id: number) => ApiPlayer | undefined;

  /**
   * 获取已审核的陪玩师列表
   * @returns 审核状态为 verified 的陪玩师数组
   */
  getVerifiedPlayers: () => ApiPlayer[];

  /**
   * 获取待审核的陪玩师列表
   * @returns 审核状态为 pending 的陪玩师数组
   */
  getPendingPlayers: () => ApiPlayer[];

  /**
   * 根据审核状态筛选陪玩师
   * @param status - 审核状态
   * @returns 指定审核状态的陪玩师数组
   */
  getPlayersByVerificationStatus: (status: 'pending' | 'verified' | 'rejected') => ApiPlayer[];

  /**
   * 根据等级筛选陪玩师
   * @param rank - 等级
   * @returns 指定等级的陪玩师数组
   */
  getPlayersByRank: (rank: string) => ApiPlayer[];

  /**
   * 搜索陪玩师（按昵称）
   * @param keyword - 关键词
   * @returns 匹配的陪玩师数组
   */
  searchPlayers: (keyword: string) => ApiPlayer[];
}

/**
 * 陪玩师 Store 实现
 */
export const usePlayerStore = create<PlayerState>()(
  persist(
    (set, get) => ({
      // ========== Initial State ==========
      players: [],
      currentPlayer: null,
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
       * 获取陪玩师列表
       */
      fetchPlayers: async (page = 1, pageSize = 10) => {
        set({ loading: true, error: null });

        try {
          const { filters } = get();

          // 调用 API
          const response = await adminApi.getPlayers({
            page,
            page_size: pageSize,
            keyword: filters.keyword,
            status: filters.verification_status,
          });

          // 更新状态
          set({
            players: response.data.data || [],
            pagination: {
              current: page,
              pageSize,
              total: response.data.pagination?.total || (response.data.data?.length || 0),
            },
            loading: false,
          });
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '获取陪玩师列表失败';
          set({
            error: errorMessage,
            loading: false,
          });
          // 向上层抛出错误，让调用方处理
          throw error;
        }
      },

      /**
       * 获取单个陪玩师详情
       */
      fetchPlayer: async (id: number) => {
        set({ loading: true, error: null });

        try {
          // 调用 API
          const response = await adminApi.getPlayer(id);

          // 更新状态
          set({
            currentPlayer: response.data.data,
            loading: false,
          });
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '获取陪玩师详情失败';
          set({
            error: errorMessage,
            loading: false,
            currentPlayer: null,
          });
          throw error;
        }
      },

      /**
       * 创建陪玩师
       */
      createPlayer: async (playerData: CreatePlayerDto) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 创建陪玩师
          const response = await adminApi.createPlayer(playerData);

          // 将新陪玩师添加到列表开头
          set((state) => ({
            players: [response.data.data, ...state.players],
            pagination: {
              ...state.pagination,
              total: state.pagination.total + 1,
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '创建陪玩师失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 更新陪玩师信息
       */
      updatePlayer: async (id: number, playerData: UpdatePlayerDto) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新陪玩师
          const response = await adminApi.updatePlayer(id, playerData);

          // 更新本地缓存中的陪玩师
          set((state) => ({
            players: state.players.map((p) =>
              p.id === id ? response.data.data : p
            ),
            currentPlayer: state.currentPlayer?.id === id ? response.data.data : state.currentPlayer,
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新陪玩师信息失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 删除陪玩师
       */
      deletePlayer: async (id: number) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 删除陪玩师
          await adminApi.deletePlayer(id);

          // 从本地缓存中移除
          set((state) => ({
            players: state.players.filter((p) => p.id !== id),
            currentPlayer: state.currentPlayer?.id === id ? null : state.currentPlayer,
            pagination: {
              ...state.pagination,
              total: Math.max(0, state.pagination.total - 1),
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '删除陪玩师失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 批量删除陪玩师
       */
      batchDeletePlayers: async (playerIds: number[]) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 批量删除
          await adminApi.batchDeletePlayers(playerIds);

          // 从本地缓存中移除
          set((state) => ({
            players: state.players.filter((p) => !playerIds.includes(p.id)),
            currentPlayer: state.currentPlayer && playerIds.includes(state.currentPlayer.id)
              ? null
              : state.currentPlayer,
            pagination: {
              ...state.pagination,
              total: Math.max(0, state.pagination.total - playerIds.length),
            },
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量删除陪玩师失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 更新陪玩师审核状态
       */
      updateVerificationStatus: async (
        id: number,
        status: 'pending' | 'verified' | 'rejected',
        remark?: string
      ) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新审核状态
          const response = await adminApi.updatePlayerVerification(id, status, remark);

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              p.id === id ? response.data.data : p
            ),
            currentPlayer: state.currentPlayer?.id === id ? response.data.data : state.currentPlayer,
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新审核状态失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 批量更新陪玩师状态
       */
      batchUpdatePlayerStatus: async (playerIds: number[], status: string) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 批量更新状态
          await adminApi.batchUpdatePlayerStatus({ playerIds, status });

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              playerIds.includes(p.id)
                ? { ...p, verificationStatus: status as ApiPlayer['verificationStatus'] }
                : p
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量更新陪玩师状态失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 更新陪玩师技能标签
       */
      updateSkillTags: async (id: number, tags: string[]) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新技能标签
          await adminApi.updatePlayerSkillTags(id, tags);

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              p.id === id ? { ...p, skillTags: tags } : p
            ),
            currentPlayer: state.currentPlayer?.id === id
              ? { ...state.currentPlayer, skillTags: tags }
              : state.currentPlayer,
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新技能标签失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 调整陪玩师时薪价格
       */
      updatePrice: async (id: number, hourlyRateCents: number) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新价格
          const response = await adminApi.updatePlayer(id, {
            hourlyRateCents,
            verificationStatus: 'verified', // 保持当前审核状态
          });

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              p.id === id ? response.data.data : p
            ),
            currentPlayer: state.currentPlayer?.id === id ? response.data.data : state.currentPlayer,
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '调整价格失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 更新陪玩师等级
       */
      updateRank: async (id: number, rank: string) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 更新等级
          const response = await adminApi.updatePlayer(id, {
            rank,
            verificationStatus: 'verified', // 保持当前审核状态
          });

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              p.id === id ? response.data.data : p
            ),
            currentPlayer: state.currentPlayer?.id === id ? response.data.data : state.currentPlayer,
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新等级失败';
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
      setFilters: (newFilters: Partial<PlayerFilters>) => {
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
          players: [],
          currentPlayer: null,
          loading: false,
          error: null,
          pagination: { current: 1, pageSize: 10, total: 0 },
          filters: {},
        });
      },

      // ========== Selectors ==========

      /**
       * 根据 ID 获取陪玩师
       */
      getPlayerById: (id: number) => {
        return get().players.find((p) => p.id === id);
      },

      /**
       * 获取已审核的陪玩师
       */
      getVerifiedPlayers: () => {
        return get().players.filter((p) => p.verificationStatus === 'verified');
      },

      /**
       * 获取待审核的陪玩师
       */
      getPendingPlayers: () => {
        return get().players.filter((p) => p.verificationStatus === 'pending');
      },

      /**
       * 根据审核状态筛选陪玩师
       */
      getPlayersByVerificationStatus: (status: 'pending' | 'verified' | 'rejected') => {
        return get().players.filter((p) => p.verificationStatus === status);
      },

      /**
       * 根据等级筛选陪玩师
       */
      getPlayersByRank: (rank: string) => {
        return get().players.filter((p) => p.rank === rank);
      },

      /**
       * 搜索陪玩师（按昵称）
       */
      searchPlayers: (keyword: string) => {
        const lowerKeyword = keyword.toLowerCase();
        return get().players.filter((p) =>
          p.nickname.toLowerCase().includes(lowerKeyword)
        );
      },
    }),
    {
      name: 'player-storage',
      // 只缓存部分数据，避免 localStorage 过大
      partialize: (state) => ({
        // 最多缓存最近 50 条陪玩师数据
        players: state.players.slice(0, 50),
        // 保留筛选条件
        filters: state.filters,
      }),
    }
  )
);

/**
 * 与 authStore 联动的辅助函数
 * 在 authStore 的 logout 方法中调用此函数清理陪玩师数据
 */
export const clearPlayerStore = () => {
  usePlayerStore.getState().reset();
};

/**
 * 导出便捷的 hooks
 */
export const usePlayers = () => usePlayerStore((state) => state.players);
export const useCurrentPlayer = () => usePlayerStore((state) => state.currentPlayer);
export const usePlayerLoading = () => usePlayerStore((state) => state.loading);
export const usePlayerError = () => usePlayerStore((state) => state.error);
export const usePlayerPagination = () => usePlayerStore((state) => state.pagination);
export const usePlayerFilters = () => usePlayerStore((state) => state.filters);
export const useVerifiedPlayers = () => usePlayerStore((state) => state.getVerifiedPlayers());
export const usePendingPlayers = () => usePlayerStore((state) => state.getPendingPlayers());
