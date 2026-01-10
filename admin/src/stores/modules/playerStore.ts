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
 * - 使用 PlayerService 进行业务逻辑处理
 * - 与 authStore 联动（登出时清理数据）
 *
 * @module playerStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Player as ApiPlayer, CreatePlayerDto, UpdatePlayerDto, Order } from '@/api/admin';
import {
  playerService as defaultPlayerService,
  type IPlayerService,
  type PlayerStatistics,
  type EarningsCalculation,
  type VerificationCheckResult,
  type SkillTagValidationResult,
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
interface PlayerFilters {
  status?: string; // 审核状态: pending/verified/rejected
  keyword?: string; // 昵称搜索
  verification_status?: 'pending' | 'verified' | 'rejected';
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
  error: StoreError | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: PlayerFilters;
  /** 最后一次批量操作结果 */
  lastBatchResult: BatchResult<void> | null;

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
  fetchPlayer: (id: number) => Promise<ServiceResult<ApiPlayer>>;

  /**
   * 创建陪玩师
   * @param playerData - 陪玩师数据
   * @returns 服务结果
   */
  createPlayer: (playerData: CreatePlayerDto) => Promise<ServiceResult<ApiPlayer>>;

  /**
   * 更新陪玩师信息
   * @param id - 陪玩师 ID
   * @param playerData - 更新的陪玩师数据
   * @returns 服务结果
   */
  updatePlayer: (id: number, playerData: UpdatePlayerDto) => Promise<ServiceResult<ApiPlayer>>;

  /**
   * 删除陪玩师
   * @param id - 陪玩师 ID
   * @returns 服务结果
   */
  deletePlayer: (id: number) => Promise<ServiceResult<void>>;

  /**
   * 批量删除陪玩师
   * @param playerIds - 陪玩师 ID 数组
   * @returns 批量操作结果
   */
  batchDeletePlayers: (playerIds: number[]) => Promise<BatchResult<void>>;

  /**
   * 更新陪玩师审核状态
   * @param id - 陪玩师 ID
   * @param status - 审核状态（pending/verified/rejected）
   * @param remark - 备注/拒绝原因
   * @returns 服务结果
   */
  updateVerificationStatus: (
    id: number,
    status: 'pending' | 'verified' | 'rejected',
    remark?: string
  ) => Promise<ServiceResult<ApiPlayer>>;

  /**
   * 批量更新陪玩师状态
   * @param playerIds - 陪玩师 ID 数组
   * @param status - 状态
   * @returns 批量操作结果
   */
  batchUpdatePlayerStatus: (playerIds: number[], status: string) => Promise<BatchResult<void>>;

  /**
   * 更新陪玩师技能标签
   * @param id - 陪玩师 ID
   * @param tags - 技能标签数组
   * @returns 服务结果
   */
  updateSkillTags: (id: number, tags: string[]) => Promise<ServiceResult<void>>;

  /**
   * 检查是否可以更新审核状态
   * @param player - 陪玩师对象
   * @param newStatus - 新状态
   * @returns 验证结果
   */
  canVerify: (player: ApiPlayer, newStatus: string) => VerificationCheckResult;

  /**
   * 验证技能标签
   * @param tags - 技能标签数组
   * @returns 验证结果
   */
  validateSkillTags: (tags: string[]) => SkillTagValidationResult;

  /**
   * 解析技能标签字符串
   * @param tagsString - 逗号分隔的标签字符串
   * @returns 标签数组
   */
  parseSkillTags: (tagsString: string) => string[];

  /**
   * 计算陪玩师收益
   * @param order - 订单对象
   * @returns 收益计算结果
   */
  calculateEarnings: (order: Order) => EarningsCalculation;

  /**
   * 计算陪玩师统计
   * @param player - 陪玩师对象
   * @param orders - 订单列表
   * @returns 统计数据
   */
  computeStatistics: (player: ApiPlayer, orders: Order[]) => PlayerStatistics;

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

  /**
   * 清除错误状态
   */
  clearError: () => void;

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
 * 创建 PlayerStore 的工厂函数
 * 支持依赖注入，便于测试
 */
export const createPlayerStore = (playerService: IPlayerService = defaultPlayerService) => {
  return create<PlayerState>()(
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
        lastBatchResult: null,

        // ========== Actions ==========

        /**
         * 获取陪玩师列表
         */
        fetchPlayers: async (page = 1, pageSize = 10) => {
          set({ loading: true, error: null });

          const { filters } = get();

          // 使用 PlayerService 获取陪玩师
          const result = await playerService.getPlayers({
            page,
            page_size: pageSize,
            keyword: filters.keyword,
            status: filters.verification_status,
          });

          if (result.success && result.data) {
            set({
              players: result.data,
              pagination: {
                current: page,
                pageSize,
                total: result.data.length, // Note: API should return total in pagination
              },
              loading: false,
            });
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '获取陪玩师列表失败' };
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
         * 获取单个陪玩师详情
         */
        fetchPlayer: async (id: number) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 获取陪玩师
          const result = await playerService.getPlayerById(id);

          if (result.success && result.data) {
            set({
              currentPlayer: result.data,
              loading: false,
            });
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '获取陪玩师详情失败' };
            set({
              error: {
                code: error.code,
                message: error.message,
                details: error.details,
              },
              loading: false,
              currentPlayer: null,
            });
          }

          return result;
        },

        /**
         * 创建陪玩师
         */
        createPlayer: async (playerData: CreatePlayerDto) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 创建陪玩师
          const result = await playerService.createPlayer(playerData);

          if (result.success && result.data) {
            // 将新陪玩师添加到列表开头
            set((state) => ({
              players: [result.data!, ...state.players],
              pagination: {
                ...state.pagination,
                total: state.pagination.total + 1,
              },
              loading: false,
            }));
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '创建陪玩师失败' };
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
        },

        /**
         * 更新陪玩师信息
         */
        updatePlayer: async (id: number, playerData: UpdatePlayerDto) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 更新陪玩师
          const result = await playerService.updatePlayer(id, playerData);

          if (result.success && result.data) {
            // 更新本地缓存中的陪玩师
            set((state) => ({
              players: state.players.map((p) =>
                p.id === id ? result.data! : p
              ),
              currentPlayer: state.currentPlayer?.id === id ? result.data! : state.currentPlayer,
              loading: false,
            }));
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新陪玩师信息失败' };
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
        },

        /**
         * 删除陪玩师
         */
        deletePlayer: async (id: number) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 删除陪玩师
          const result = await playerService.deletePlayer(id);

          if (result.success) {
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
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '删除陪玩师失败' };
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
        },

        /**
         * 批量删除陪玩师
         */
        batchDeletePlayers: async (playerIds: number[]) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 批量删除
          const result = await playerService.batchDelete(playerIds);

          // 获取成功删除的陪玩师 ID
          const successfulIds = result.results
            .filter((r) => r.success)
            .map((r) => playerIds[r.index]);

          // 从本地缓存中移除成功删除的陪玩师
          set((state) => ({
            players: state.players.filter((p) => !successfulIds.includes(p.id)),
            currentPlayer: state.currentPlayer && successfulIds.includes(state.currentPlayer.id)
              ? null
              : state.currentPlayer,
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
                message: `批量删除部分失败：${failedCount} 个陪玩师删除失败`,
                details: { failedCount, results: result.results },
              },
            });
          }

          return result;
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

          // 使用 PlayerService 更新审核状态
          const result = await playerService.verifyPlayer(id, status, remark);

          if (result.success && result.data) {
            // 更新本地缓存
            set((state) => ({
              players: state.players.map((p) =>
                p.id === id ? result.data! : p
              ),
              currentPlayer: state.currentPlayer?.id === id ? result.data! : state.currentPlayer,
              loading: false,
            }));
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新审核状态失败' };
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
        },

        /**
         * 批量更新陪玩师状态
         */
        batchUpdatePlayerStatus: async (playerIds: number[], status: string) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 批量更新状态
          const result = await playerService.batchUpdateStatus(playerIds, status);

          // 获取成功更新的陪玩师 ID
          const successfulIds = result.results
            .filter((r) => r.success)
            .map((r) => playerIds[r.index]);

          // 更新本地缓存
          set((state) => ({
            players: state.players.map((p) =>
              successfulIds.includes(p.id)
                ? { ...p, verificationStatus: status as ApiPlayer['verificationStatus'] }
                : p
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
                message: `批量更新状态部分失败：${failedCount} 个陪玩师更新失败`,
                details: { failedCount, results: result.results },
              },
            });
          }

          return result;
        },

        /**
         * 更新陪玩师技能标签
         */
        updateSkillTags: async (id: number, tags: string[]) => {
          set({ loading: true, error: null });

          // 使用 PlayerService 更新技能标签
          const result = await playerService.updateSkillTags(id, tags);

          if (result.success) {
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
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '更新技能标签失败' };
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
        },

        /**
         * 检查是否可以更新审核状态
         */
        canVerify: (player: ApiPlayer, newStatus: string) => {
          return playerService.canVerify(player, newStatus);
        },

        /**
         * 验证技能标签
         */
        validateSkillTags: (tags: string[]) => {
          return playerService.validateSkillTags(tags);
        },

        /**
         * 解析技能标签字符串
         */
        parseSkillTags: (tagsString: string) => {
          return playerService.parseSkillTags(tagsString);
        },

        /**
         * 计算陪玩师收益
         */
        calculateEarnings: (order: Order) => {
          return playerService.calculateEarnings(order);
        },

        /**
         * 计算陪玩师统计
         */
        computeStatistics: (player: ApiPlayer, orders: Order[]) => {
          return playerService.computeStatistics(player, orders);
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
};

/**
 * 默认 PlayerStore 实例
 */
export const usePlayerStore = createPlayerStore();

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
export const useLastPlayerBatchResult = () => usePlayerStore((state) => state.lastBatchResult);

// Re-export types for convenience
export type { PlayerStatistics, EarningsCalculation, VerificationCheckResult, SkillTagValidationResult };
