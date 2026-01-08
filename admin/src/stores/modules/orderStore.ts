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
 * - 使用 OrderService 进行业务逻辑处理
 *
 * @module orderStore
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Order as ApiOrder, OrderQueryParams } from '@/api/admin';
import {
  orderService as defaultOrderService,
  type IOrderService,
  type OrderStatistics,
  type RefundCalculation,
  type CancellationCheckResult,
  type RefundCheckResult,
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
 * Store 错误信息接口
 */
interface StoreError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
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
  error: StoreError | null;
  /** 分页信息 */
  pagination: Pagination;
  /** 筛选条件 */
  filters: OrderFilters;
  /** 最后一次批量操作结果 */
  lastBatchResult: BatchResult<void> | null;

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
   */
  updateOrderStatus: (
    id: number,
    status: ApiOrder['status']
  ) => void;

  /**
   * 取消订单
   * @param id - 订单 ID
   * @param note - 取消备注
   * @returns 服务结果
   */
  cancelOrder: (id: number, note?: string) => Promise<ServiceResult<ApiOrder>>;

  /**
   * 退款订单
   * @param id - 订单 ID
   * @param data - 退款数据（包含原因和金额）
   * @returns 服务结果
   */
  refundOrder: (
    id: number,
    data: { reason: string; amount_cents: number; note?: string }
  ) => Promise<ServiceResult<ApiOrder>>;

  /**
   * 批量取消订单
   * @param orderIds - 订单 ID 数组
   * @param reason - 取消原因
   * @returns 批量操作结果
   */
  batchCancelOrders: (
    orderIds: number[],
    reason?: string
  ) => Promise<BatchResult<void>>;

  /**
   * 批量完成订单
   * @param orderIds - 订单 ID 数组
   * @returns 批量操作结果
   */
  batchCompleteOrders: (orderIds: number[]) => Promise<BatchResult<void>>;

  /**
   * 检查订单是否可以取消
   * @param order - 订单对象
   * @returns 取消检查结果
   */
  canCancelOrder: (order: ApiOrder) => CancellationCheckResult;

  /**
   * 检查订单是否可以退款
   * @param order - 订单对象
   * @returns 退款检查结果
   */
  canRefundOrder: (order: ApiOrder) => RefundCheckResult;

  /**
   * 计算退款金额
   * @param order - 订单对象
   * @param requestedAmount - 请求退款金额
   * @returns 退款计算结果
   */
  calculateRefund: (order: ApiOrder, requestedAmount: number) => RefundCalculation;

  /**
   * 计算订单统计
   * @param orders - 订单列表（可选，默认使用当前缓存的订单）
   * @returns 订单统计
   */
  computeStatistics: (orders?: ApiOrder[]) => OrderStatistics;

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

  /**
   * 清除错误状态
   */
  clearError: () => void;

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
 * 创建 OrderStore 的工厂函数
 * 支持依赖注入，便于测试
 */
export const createOrderStore = (orderService: IOrderService = defaultOrderService) => {
  return create<OrderState>()(
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
        lastBatchResult: null,

        // ========== Actions ==========

        /**
         * 获取订单列表
         */
        fetchOrders: async (page = 1, pageSize = 10) => {
          set({ loading: true, error: null });

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

          // 使用 OrderService 获取订单
          const result = await orderService.getOrders(params);

          if (result.success && result.data) {
            set({
              orders: result.data,
              pagination: {
                current: page,
                pageSize,
                total: result.data.length, // Note: API should return total in pagination
              },
              loading: false,
            });
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '获取订单列表失败' };
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
         * 根据 ID 获取订单详情
         */
        getOrderById: (id: number) => {
          return get().orders.find((o) => o.id === id);
        },

        /**
         * 更新订单状态（本地缓存）
         */
        updateOrderStatus: (id: number, status: ApiOrder['status']) => {
          set((state) => ({
            orders: state.orders.map((o) =>
              o.id === id ? { ...o, status } : o
            ),
          }));
        },

        /**
         * 取消订单
         */
        cancelOrder: async (id: number, note?: string) => {
          set({ loading: true, error: null });

          // 使用 OrderService 取消订单
          const result = await orderService.cancelOrder(id, note);

          if (result.success && result.data) {
            // 更新本地缓存
            set((state) => ({
              orders: state.orders.map((o) =>
                o.id === id ? result.data! : o
              ),
              loading: false,
            }));
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '取消订单失败' };
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
         * 退款订单
         */
        refundOrder: async (
          id: number,
          data: { reason: string; amount_cents: number; note?: string }
        ) => {
          set({ loading: true, error: null });

          // 使用 OrderService 退款
          const result = await orderService.refundOrder(id, data);

          if (result.success && result.data) {
            // 更新本地缓存
            set((state) => ({
              orders: state.orders.map((o) =>
                o.id === id ? result.data! : o
              ),
              loading: false,
            }));
          } else {
            const error = result.error || { code: 'UNKNOWN_ERROR', message: '退款订单失败' };
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
         * 批量取消订单
         */
        batchCancelOrders: async (orderIds: number[], reason?: string) => {
          set({ loading: true, error: null });

          // 使用 OrderService 批量取消
          const result = await orderService.batchCancel(orderIds, reason);

          // 获取成功取消的订单 ID
          const successfulIds = result.results
            .filter((r) => r.success)
            .map((r) => orderIds[r.index]);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              successfulIds.includes(o.id) ? { ...o, status: 'canceled' as const } : o
            ),
            selectedOrderIds: [], // 清空选中
            loading: false,
            lastBatchResult: result,
          }));

          // 如果有失败的操作，设置错误信息
          if (!result.success) {
            const failedCount = result.failed;
            set({
              error: {
                code: 'BATCH_PARTIAL_FAILURE',
                message: `批量取消部分失败：${failedCount} 个订单取消失败`,
                details: { failedCount, results: result.results },
              },
            });
          }

          return result;
        },

        /**
         * 批量完成订单
         */
        batchCompleteOrders: async (orderIds: number[]) => {
          set({ loading: true, error: null });

          // 使用 OrderService 批量完成
          const result = await orderService.batchComplete(orderIds);

          // 获取成功完成的订单 ID
          const successfulIds = result.results
            .filter((r) => r.success)
            .map((r) => orderIds[r.index]);

          // 更新本地缓存
          set((state) => ({
            orders: state.orders.map((o) =>
              successfulIds.includes(o.id) ? { ...o, status: 'completed' as const } : o
            ),
            selectedOrderIds: [], // 清空选中
            loading: false,
            lastBatchResult: result,
          }));

          // 如果有失败的操作，设置错误信息
          if (!result.success) {
            const failedCount = result.failed;
            set({
              error: {
                code: 'BATCH_PARTIAL_FAILURE',
                message: `批量完成部分失败：${failedCount} 个订单完成失败`,
                details: { failedCount, results: result.results },
              },
            });
          }

          return result;
        },

        /**
         * 检查订单是否可以取消
         */
        canCancelOrder: (order: ApiOrder) => {
          return orderService.canCancel(order);
        },

        /**
         * 检查订单是否可以退款
         */
        canRefundOrder: (order: ApiOrder) => {
          return orderService.canRefund(order);
        },

        /**
         * 计算退款金额
         */
        calculateRefund: (order: ApiOrder, requestedAmount: number) => {
          return orderService.calculateRefund(order, requestedAmount);
        },

        /**
         * 计算订单统计
         */
        computeStatistics: (orders?: ApiOrder[]) => {
          const ordersToCompute = orders || get().orders;
          return orderService.computeStatistics(ordersToCompute);
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
};

/**
 * 默认 OrderStore 实例
 */
export const useOrderStore = createOrderStore();

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
export const useLastOrderBatchResult = () => useOrderStore((state) => state.lastBatchResult);

// Re-export types for convenience
export type { OrderStatistics, RefundCalculation, CancellationCheckResult, RefundCheckResult };
