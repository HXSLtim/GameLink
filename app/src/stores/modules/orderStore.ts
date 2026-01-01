// Order Store - Taro App
// Order management state for user side
// Features: My orders, create order, order details, payment, review, dispute

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import Taro from '@tarojs/taro';
import { get, post, put } from '../../api/client';
import type {
  Order,
  OrderStatus,
  CreateOrderRequest,
  OrderDraft,
  OrderListParams,
  PaymentRequest,
  ReviewRequest,
  DisputeRequest,
} from '../types';

/**
 * Pagination info
 */
interface PaginationInfo {
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
}

/**
 * Order store state and actions
 */
interface OrderState {
  // State
  orders: Order[];
  currentOrder: Order | null;
  orderDraft: OrderDraft | null;
  loading: boolean;
  pagination: PaginationInfo;

  // Filters
  activeFilter: OrderStatus | 'all';

  // Actions - Order List
  fetchOrders: (params?: OrderListParams) => Promise<void>;
  loadMoreOrders: () => Promise<void>;
  setFilter: (status: OrderStatus | 'all') => void;
  refreshOrders: () => Promise<void>;

  // Actions - Order Detail
  fetchOrderDetail: (orderId: number) => Promise<void>;
  setCurrentOrder: (order: Order | null) => void;

  // Actions - Create Order
  createOrderDraft: (draft: OrderDraft) => void;
  updateOrderDraft: (updates: Partial<OrderDraft>) => void;
  clearOrderDraft: () => void;
  submitOrder: () => Promise<number>; // Returns order ID

  // Actions - Payment
  payOrder: (request: PaymentRequest) => Promise<void>;
  cancelOrder: (orderId: number, reason?: string) => Promise<void>;

  // Actions - Review & Dispute
  reviewOrder: (request: ReviewRequest) => Promise<void>;
  fileDispute: (request: DisputeRequest) => Promise<void>;

  // Selectors (computed values)
  getOrdersByStatus: (status: OrderStatus) => Order[];
  getOrderById: (orderId: number) => Order | undefined;
  activeOrders: () => Order[];
  completedOrders: () => Order[];
}

/**
 * Initial pagination state
 */
const initialPagination: PaginationInfo = {
  total: 0,
  page: 1,
  pageSize: 10,
  hasMore: true,
};

/**
 * Order store with persistence for draft only
 */
export const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      // Initial state
      orders: [],
      currentOrder: null,
      orderDraft: null,
      loading: false,
      pagination: initialPagination,
      activeFilter: 'all',

      /**
       * Fetch user orders with optional filters
       * GET /api/v1/user/orders
       */
      fetchOrders: async (params: OrderListParams = {}) => {
        const { status, page = 1, pageSize = 10 } = params;

        set({ loading: true });

        try {
          const queryParams = new URLSearchParams();
          if (status && status !== 'all') {
            queryParams.append('status', status);
          }
          queryParams.append('page', String(page));
          queryParams.append('pageSize', String(pageSize));

          const response = await get<any>(`/user/orders?${queryParams.toString()}`);

          if (response.success && response.data) {
            const { items, total, page: currentPage, pageSize: currentPageSize } = response.data;

            // Add computed fields to orders
            const ordersWithComputed = (items || []).map((order: Order) => ({
              ...order,
              amountYuan: order.amount / 100,
              canCancel: ['pending', 'paid'].includes(order.status),
              canReview: order.status === 'completed',
              canDispute: ['completed', 'in_progress'].includes(order.status),
            }));

            set({
              orders: ordersWithComputed,
              pagination: {
                total,
                page: currentPage,
                pageSize: currentPageSize,
                hasMore: currentPage * currentPageSize < total,
              },
              loading: false,
            });
          } else {
            throw new Error(response.message || 'Failed to fetch orders');
          }
        } catch (error: any) {
          console.error('Fetch orders error:', error);

          Taro.showToast({
            title: error.message || '获取订单列表失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
        }
      },

      /**
       * Load more orders (infinite scroll)
       */
      loadMoreOrders: async () => {
        const { pagination, activeFilter, orders } = get();

        if (!pagination.hasMore || get().loading) {
          return;
        }

        const nextPage = pagination.page + 1;
        set({ loading: true });

        try {
          const queryParams = new URLSearchParams();
          if (activeFilter !== 'all') {
            queryParams.append('status', activeFilter);
          }
          queryParams.append('page', String(nextPage));
          queryParams.append('pageSize', String(pagination.pageSize));

          const response = await get<any>(`/user/orders?${queryParams.toString()}`);

          if (response.success && response.data) {
            const { items, total } = response.data;

            const newOrders = (items || []).map((order: Order) => ({
              ...order,
              amountYuan: order.amount / 100,
              canCancel: ['pending', 'paid'].includes(order.status),
              canReview: order.status === 'completed',
              canDispute: ['completed', 'in_progress'].includes(order.status),
            }));

            set({
              orders: [...orders, ...newOrders],
              pagination: {
                ...pagination,
                page: nextPage,
                hasMore: nextPage * pagination.pageSize < total,
              },
              loading: false,
            });
          }
        } catch (error: any) {
          console.error('Load more orders error:', error);
          set({ loading: false });
        }
      },

      /**
       * Set filter for order list
       */
      setFilter: (status: OrderStatus | 'all') => {
        set({ activeFilter: status, pagination: initialPagination });
        get().fetchOrders({ status, page: 1 });
      },

      /**
       * Refresh current order list
       */
      refreshOrders: async () => {
        const { activeFilter, pagination } = get();
        await get().fetchOrders({
          status: activeFilter,
          page: 1,
          pageSize: pagination.pageSize,
        });
      },

      /**
       * Fetch order detail by ID
       * GET /api/v1/user/orders/:id
       */
      fetchOrderDetail: async (orderId: number) => {
        set({ loading: true });

        try {
          const response = await get<any>(`/user/orders/${orderId}`);

          if (response.success && response.data) {
            const order = {
              ...response.data,
              amountYuan: response.data.amount / 100,
              canCancel: ['pending', 'paid'].includes(response.data.status),
              canReview: response.data.status === 'completed',
              canDispute: ['completed', 'in_progress'].includes(response.data.status),
            };

            set({ currentOrder: order, loading: false });
            return order;
          } else {
            throw new Error(response.message || 'Failed to fetch order detail');
          }
        } catch (error: any) {
          console.error('Fetch order detail error:', error);

          Taro.showToast({
            title: error.message || '获取订单详情失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * Set current order manually
       */
      setCurrentOrder: (order: Order | null) => {
        set({ currentOrder: order });
      },

      /**
       * Create order draft (local state, not saved to server)
       */
      createOrderDraft: (draft: OrderDraft) => {
        set({ orderDraft: draft });

        // Save to storage for persistence
        Taro.setStorageSync('orderDraft', draft);
      },

      /**
       * Update order draft
       */
      updateOrderDraft: (updates: Partial<OrderDraft>) => {
        const { orderDraft } = get();
        if (!orderDraft) {
          return;
        }

        const updatedDraft = { ...orderDraft, ...updates };
        set({ orderDraft: updatedDraft });

        Taro.setStorageSync('orderDraft', updatedDraft);
      },

      /**
       * Clear order draft
       */
      clearOrderDraft: () => {
        set({ orderDraft: null });
        Taro.removeStorageSync('orderDraft');
      },

      /**
       * Submit order to server
       * POST /api/v1/user/orders
       */
      submitOrder: async () => {
        const { orderDraft } = get();

        if (!orderDraft) {
          Taro.showToast({
            title: '请先选择服务项目',
            icon: 'none',
            duration: 2000,
          });
          throw new Error('No order draft');
        }

        set({ loading: true });

        try {
          const request: CreateOrderRequest = {
            itemId: orderDraft.itemId,
            playerId: orderDraft.playerId,
            quantity: orderDraft.quantity,
            scheduledStart: orderDraft.scheduledStart,
            gameIds: orderDraft.gameIds,
            remark: orderDraft.remark,
          };

          const response = await post<any>('/user/orders', request);

          if (response.success && response.data) {
            // Clear draft after successful submission
            get().clearOrderDraft();

            set({ loading: false });

            // Navigate to order detail page
            Taro.navigateTo({
              url: `/pages/order-detail/index?id=${response.data.id}`,
            });

            return response.data.id;
          } else {
            throw new Error(response.message || 'Failed to create order');
          }
        } catch (error: any) {
          console.error('Submit order error:', error);

          Taro.showToast({
            title: error.message || '创建订单失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * Pay for order
       * POST /api/v1/user/orders/:id/pay
       */
      payOrder: async (request: PaymentRequest) => {
        set({ loading: true });

        try {
          // Show loading
          Taro.showLoading({
            title: '支付中...',
            mask: true,
          });

          const response = await post<any>(`/user/orders/${request.orderId}/pay`, {
            method: request.method,
            amount: request.amount,
          });

          Taro.hideLoading();

          if (response.success && response.data) {
            const { paymentParams } = response.data;

            // Use Taro payment interface
            if (request.method === 'wechat' || request.method === 'alipay') {
              await Taro.requestPayment({
                ...paymentParams,
                success: () => {
                  Taro.showToast({
                    title: '支付成功',
                    icon: 'success',
                    duration: 2000,
                  });

                  // Refresh order detail
                  get().fetchOrderDetail(request.orderId);

                  // Navigate to order detail page
                  setTimeout(() => {
                    Taro.redirectTo({
                      url: `/pages/order-detail/index?id=${request.orderId}`,
                    });
                  }, 1500);
                },
                fail: (err: any) => {
                  if (err.errMsg.includes('cancel')) {
                    Taro.showToast({
                      title: '取消支付',
                      icon: 'none',
                      duration: 2000,
                    });
                  } else {
                    Taro.showToast({
                      title: '支付失败',
                      icon: 'none',
                      duration: 2000,
                    });
                  }
                },
              });
            } else if (request.method === 'wallet') {
              // Wallet payment succeeds directly
              Taro.showToast({
                title: '支付成功',
                icon: 'success',
                duration: 2000,
              });

              get().fetchOrderDetail(request.orderId);

              setTimeout(() => {
                Taro.redirectTo({
                  url: `/pages/order-detail/index?id=${request.orderId}`,
                });
              }, 1500);
            }

            set({ loading: false });
          } else {
            throw new Error(response.message || 'Failed to initiate payment');
          }
        } catch (error: any) {
          console.error('Pay order error:', error);

          Taro.hideLoading();

          Taro.showToast({
            title: error.message || '支付失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * Cancel order
       * POST /api/v1/user/orders/:id/cancel
       */
      cancelOrder: async (orderId: number, reason?: string) => {
        // Show confirmation dialog
        const { confirm } = await Taro.showModal({
          title: '确认取消订单',
          content: reason || '确定要取消此订单吗？',
          confirmText: '确认取消',
          cancelText: '再想想',
        });

        if (!confirm) {
          return;
        }

        set({ loading: true });

        try {
          const response = await post<any>(`/user/orders/${orderId}/cancel`, {
            reason,
          });

          if (response.success) {
            Taro.showToast({
              title: '订单已取消',
              icon: 'success',
              duration: 2000,
            });

            // Refresh order detail
            await get().fetchOrderDetail(orderId);

            // Refresh order list
            await get().refreshOrders();

            set({ loading: false });
          } else {
            throw new Error(response.message || 'Failed to cancel order');
          }
        } catch (error: any) {
          console.error('Cancel order error:', error);

          Taro.showToast({
            title: error.message || '取消订单失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * Review completed order
       * POST /api/v1/user/orders/:id/review
       */
      reviewOrder: async (request: ReviewRequest) => {
        set({ loading: true });

        try {
          const response = await post<any>(
            `/user/orders/${request.orderId}/review`,
            {
              playerId: request.playerId,
              rating: request.rating,
              content: request.content,
              tags: request.tags,
            }
          );

          if (response.success) {
            Taro.showToast({
              title: '评价成功',
              icon: 'success',
              duration: 2000,
            });

            // Refresh order detail
            await get().fetchOrderDetail(request.orderId);

            set({ loading: false });
          } else {
            throw new Error(response.message || 'Failed to submit review');
          }
        } catch (error: any) {
          console.error('Review order error:', error);

          Taro.showToast({
            title: error.message || '评价失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      /**
       * File dispute for order
       * POST /api/v1/user/orders/:id/dispute
       */
      fileDispute: async (request: DisputeRequest) => {
        set({ loading: true });

        try {
          const response = await post<any>(
            `/user/orders/${request.orderId}/dispute`,
            {
              type: request.type,
              reason: request.reason,
              evidenceUrls: request.evidenceUrls,
              evidenceText: request.evidenceText,
            }
          );

          if (response.success) {
            Taro.showToast({
              title: '投诉已提交',
              icon: 'success',
              duration: 2000,
            });

            // Refresh order detail
            await get().fetchOrderDetail(request.orderId);

            set({ loading: false });
          } else {
            throw new Error(response.message || 'Failed to file dispute');
          }
        } catch (error: any) {
          console.error('File dispute error:', error);

          Taro.showToast({
            title: error.message || '提交投诉失败',
            icon: 'none',
            duration: 2000,
          });

          set({ loading: false });
          throw error;
        }
      },

      // Selectors (computed values)

      /**
       * Get orders by status
       */
      getOrdersByStatus: (status: OrderStatus) => {
        const { orders } = get();
        return orders.filter((order) => order.status === status);
      },

      /**
       * Get order by ID
       */
      getOrderById: (orderId: number) => {
        const { orders } = get();
        return orders.find((order) => order.id === orderId);
      },

      /**
       * Get active (in-progress) orders
       */
      activeOrders: () => {
        const { orders } = get();
        return orders.filter((order) => ['pending', 'paid', 'in_progress'].includes(order.status));
      },

      /**
       * Get completed orders
       */
      completedOrders: () => {
        const { orders } = get();
        return orders.filter((order) => order.status === 'completed');
      },
    }),
    {
      name: 'order-storage',
      // Only persist order draft, not the orders list or current order
      partialize: (state) => ({
        orderDraft: state.orderDraft,
      }),
    }
  )
);

export type { OrderState, PaginationInfo };
