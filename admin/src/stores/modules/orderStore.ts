/**
 * Order Store - 订单数据管理
 *
 * 功能：
 * - 订单列表缓存（带分页）
 * - 订单详情缓存
 * - 订单状态更新（pending → confirmed → in_progress → completed）
 * - 筛选功能（状态、日期范围、用户、陪玩师、订单号）
 * - 批量操作（批量取消、批量完成）
 * - 订单取消和退款
 *
 * @module orderStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { adminApi, type Order as ApiOrder, type OrderQueryParams } from '@/api/admin';

/**
 * 分页参数接口
 */
interface Pagination {
  current: number;
  pageSize: number;
  total: number;
}

/**
 * 订单筛选条件接口
 */
interface OrderFilters {
  status?: string;
  userId?: number;
  playerId?: number;
  orderNumber?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * 订单状态接口
 */
interface OrderState {
  // ========== State ==========
  /** 订单列表缓存 */
  orders: ApiOrder[];
  /** 当前选中的订单 ID 列表（用于批量操作） */
  selectedOrderIds: number[];
  /** 加载状态 */
  loading: boolean;
  /** 错误信息 */
  error: string | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: OrderFilters;

  // ========== Actions ==========

  /**
   * 获取订单列表
   * @param page - 页码（默认 1）
   * @param pageSize - 每页数量（默认 10）
   */
  fetchOrders: (page?: number, pageSize?: number) => Promise<void>;

  /**
   * 根据 ID 获取订单详情
   * @param id - 订单 ID
   * @returns 订单对象或 undefined
   */
  getOrderById: (id: number) => ApiOrder | undefined;

  /**
   * 更新订单状态
   * @param id - 订单 ID
   * @param status - 新状态
   * @throws 更新失败时抛出错误
   */
  updateOrderStatus: (
    id: number,
    status: ApiOrder['status']
  ) => Promise<void>;

  /**
   * 取消订单
   * @param id - 订单 ID
   * @param note - 取消备注
   * @throws 取消失败时抛出错误
   */
  cancelOrder: (id: number, note?: string) => Promise<void>;

  /**
   * 退款订单
   * @param id - 订单 ID
   * @param data - 退款数据（包含原因和金额）
   * @throws 退款失败时抛出错误
   */
  refundOrder: (
    id: number,
    data: { reason: string; amount_cents: number; note?: string }
  ) => Promise<void>;

  /**
   * 批量取消订单
   * @param orderIds - 订单 ID 数组
   * @param reason - 取消原因
   * @throws 批量操作失败时抛出错误
   */
  batchCancelOrders: (
    orderIds: number[],
    reason?: string
  ) => Promise<void>;

  /**
   * 批量完成订单
   * @param orderIds - 订单 ID 数组
   * @throws 批量操作失败时抛出错误
   */
  batchCompleteOrders: (orderIds: number[]) => Promise<void>;

  /**
   * 设置选中的订单 ID 列表
   * @param ids - 订单 ID 数组
   */
  setSelectedOrders: (ids: number[]) => void;

  /**
   * 清空选中的订单
   */
  clearSelectedOrders: () => void;

  /**
   * 设置筛选条件
   * @param newFilters - 新的筛选条件（部分更新）
   */
  setFilters: (newFilters: Partial<OrderFilters>) => void;

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
   * 根据状态筛选订单
   * @param status - 订单状态
   * @returns 指定状态的订单数组
   */
  getOrdersByStatus: (status: ApiOrder['status']) => ApiOrder[];

  /**
   * 根据用户 ID 筛选订单
   * @param userId - 用户 ID
   * @returns 该用户的订单数组
   */
  getOrdersByUser: (userId: number) => ApiOrder[];

  /**
   * 根据陪玩师 ID 筛选订单
   * @param playerId - 陪玩师 ID
   * @returns 该陪玩师的订单数组
   */
  getOrdersByPlayer: (playerId: number) => ApiOrder[];

  /**
   * 获取待处理订单数
   * @returns 状态为 pending 或 confirmed 的订单数量
   */
  getPendingOrdersCount: () => number;

  /**
   * 获取进行中的订单数
   * @returns 状态为 in_progress 的订单数量
   */
  getInProgressOrdersCount: () => number;
}

/**
 * 订单 Store 实现
 */
export const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      // ========== Initial State ==========
      orders: [],
      selectedOrderIds: [],
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
       * 获取订单列表
       */
      fetchOrders: async (page = 1, pageSize = 10) => {
        set({ loading: true, error: null });

        try {
          const { filters } = get();

          // 构建查询参数
          const params: OrderQueryParams = {
            page,
            page_size: pageSize,
            status: filters.status,
            userId: filters.userId,
            orderNumber: filters.orderNumber,
            dateFrom: filters.dateFrom,
            dateTo: filters.dateTo,
          };

          // 调用 API
          const response = await adminApi.getOrders(params);

          // 更新状态
          set({
            orders: response.data.data || [],
            pagination: {
              current: page,
              pageSize,
              total: response.data.pagination?.total || 0,
            },
            loading: false,
          });
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '获取订单列表失败';
          set({
            error: errorMessage,
            loading: false,
          });
          // 向上层抛出错误，让调用方处理
          throw error;
        }
      },

      /**
       * 根据 ID 获取订单详情
       */
      getOrderById: (id: number) => {
        return get().orders.find((o) => o.id === id);
      },

      /**
       * 更新订单状态
       */
      updateOrderStatus: async (id: number, status: ApiOrder['status']) => {
        set({ loading: true, error: null });

        try {
          // 注意：后端 API 可能不直接提供更新状态的接口
          // 这里假设通过其他操作（如取消、完成）来间接更新状态
          // 如果后端有专门的更新状态接口，可以在这里调用

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              o.id === id ? { ...o, status } : o
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '更新订单状态失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 取消订单
       */
      cancelOrder: async (id: number, note?: string) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 取消订单
          const response = await adminApi.cancelOrder(id, note);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              o.id === id ? response.data.data : o
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '取消订单失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 退款订单
       */
      refundOrder: async (
        id: number,
        data: { reason: string; amount_cents: number; note?: string }
      ) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 退款
          const response = await adminApi.refundOrder(id, data);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              o.id === id ? response.data.data : o
            ),
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '退款订单失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 批量取消订单
       */
      batchCancelOrders: async (orderIds: number[], reason?: string) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 批量取消
          await adminApi.batchCancelOrders(orderIds, reason);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              orderIds.includes(o.id) ? { ...o, status: 'cancelled' as const } : o
            ),
            selectedOrderIds: [], // 清空选中
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量取消订单失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 批量完成订单
       */
      batchCompleteOrders: async (orderIds: number[]) => {
        set({ loading: true, error: null });

        try {
          // 调用 API 批量完成
          await adminApi.batchCompleteOrders(orderIds);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              orderIds.includes(o.id) ? { ...o, status: 'completed' as const } : o
            ),
            selectedOrderIds: [], // 清空选中
            loading: false,
          }));
        } catch (error: unknown) {
          const errorMessage =
            error instanceof Error ? error.message : '批量完成订单失败';
          set({
            error: errorMessage,
            loading: false,
          });
          throw error;
        }
      },

      /**
       * 设置选中的订单
       */
      setSelectedOrders: (ids: number[]) => {
        set({ selectedOrderIds: ids });
      },

      /**
       * 清空选中的订单
       */
      clearSelectedOrders: () => {
        set({ selectedOrderIds: [] });
      },

      /**
       * 设置筛选条件
       */
      setFilters: (newFilters: Partial<OrderFilters>) => {
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
          orders: [],
          selectedOrderIds: [],
          loading: false,
          error: null,
          pagination: { current: 1, pageSize: 10, total: 0 },
          filters: {},
        });
      },

      // ========== Selectors ==========

      /**
       * 根据状态筛选订单
       */
      getOrdersByStatus: (status: ApiOrder['status']) => {
        return get().orders.filter((o) => o.status === status);
      },

      /**
       * 根据用户 ID 筛选订单
       */
      getOrdersByUser: (userId: number) => {
        return get().orders.filter((o) => o.userId === userId);
      },

      /**
       * 根据陪玩师 ID 筛选订单
       */
      getOrdersByPlayer: (playerId: number) => {
        return get().orders.filter((o) => o.playerId === playerId);
      },

      /**
       * 获取待处理订单数
       */
      getPendingOrdersCount: () => {
        return get().orders.filter(
          (o) => o.status === 'pending' || o.status === 'confirmed'
        ).length;
      },

      /**
       * 获取进行中的订单数
       */
      getInProgressOrdersCount: () => {
        return get().orders.filter((o) => o.status === 'in_progress').length;
      },
    }),
    {
      name: 'order-storage',
      // 只缓存部分数据，避免 localStorage 过大
      partialize: (state) => ({
        // 最多缓存最近 50 条订单数据
        orders: state.orders.slice(0, 50),
        // 保留筛选条件
        filters: state.filters,
      }),
    }
  )
);

/**
 * 与 authStore 联动的辅助函数
 * 在 authStore 的 logout 方法中调用此函数清理订单数据
 */
export const clearOrderStore = () => {
  useOrderStore.getState().reset();
};

/**
 * 导出便捷的 hooks
 */
export const useOrders = () => useOrderStore((state) => state.orders);
export const useOrderLoading = () => useOrderStore((state) => state.loading);
export const useOrderError = () => useOrderStore((state) => state.error);
export const useOrderPagination = () => useOrderStore((state) => state.pagination);
export const useOrderFilters = () => useOrderStore((state) => state.filters);
export const useSelectedOrderIds = () => useOrderStore((state) => state.selectedOrderIds);
