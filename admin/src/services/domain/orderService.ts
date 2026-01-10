/**
 * Order Domain Service
 * Encapsulates all order-related business logic
 *
 * @module services/domain/orderService
 */

import {
  BaseService,
  type ServiceDependencies,
} from './base';
import {
  ServiceErrorCodes,
  type ServiceResult,
  type BatchResult,
  ServiceResultHelper,
} from '../utils';
import type {
  Order,
  OrderQueryParams,
  TrendData,
} from '@/api/admin';

/**
 * Order cancellation check result
 */
export interface CancellationCheckResult {
  allowed: boolean;
  reason?: string;
}

/**
 * Order refund check result
 */
export interface RefundCheckResult {
  allowed: boolean;
  reason?: string;
}

/**
 * Refund calculation result
 */
export interface RefundCalculation {
  /** Original order amount in cents */
  originalAmount: number;
  /** Amount to refund in cents */
  refundAmount: number;
  /** Platform fee in cents (retained by platform) */
  platformFee: number;
  /** Amount to deduct from player earnings in cents */
  playerAmount: number;
}

/**
 * Order statistics
 */
export interface OrderStatistics {
  /** Total number of orders */
  totalOrders: number;
  /** Total revenue in cents */
  totalRevenue: number;
  /** Orders grouped by status */
  ordersByStatus: Record<string, number>;
  /** Average order value in cents */
  averageOrderValue: number;
  /** Completion rate (0-1) */
  completionRate: number;
}

/**
 * Order Service Interface
 */
export interface IOrderService {
  // Query Operations
  getOrders(params?: OrderQueryParams): Promise<ServiceResult<Order[]>>;
  getOrderById(id: number): Promise<ServiceResult<Order>>;

  // Order Operations
  cancelOrder(id: number, reason?: string): Promise<ServiceResult<Order>>;
  refundOrder(id: number, data: { reason: string; amount_cents: number; note?: string }): Promise<ServiceResult<Order>>;

  // Batch Operations
  batchCancel(orderIds: number[], reason?: string): Promise<BatchResult<void>>;
  batchComplete(orderIds: number[]): Promise<BatchResult<void>>;

  // Validation
  canCancel(order: Order): CancellationCheckResult;
  canRefund(order: Order): RefundCheckResult;
  calculateRefund(order: Order, requestedAmount: number): RefundCalculation;

  // Statistics
  computeStatistics(orders: Order[]): OrderStatistics;
  computeTrend(orders: Order[], days: number): TrendData[];
}

/**
 * Order statuses that allow cancellation
 */
const CANCELLABLE_STATUSES: Order['status'][] = ['pending', 'confirmed'];

/**
 * Order statuses that allow refund
 */
const REFUNDABLE_STATUSES: Order['status'][] = ['completed', 'in_progress'];

/**
 * Default platform commission rate (15%)
 */
const DEFAULT_PLATFORM_COMMISSION_RATE = 0.15;

/**
 * Order Service Implementation
 *
 * Provides all order-related business logic including:
 * - Query operations
 * - Cancellation and refund processing
 * - Batch operations
 * - Statistics computation
 */
export class OrderService extends BaseService implements IOrderService {
  constructor(deps: ServiceDependencies = {}) {
    super(deps);
  }

  // ==================== Query Operations ====================

  /**
   * Get orders with optional filtering
   */
  async getOrders(params?: OrderQueryParams): Promise<ServiceResult<Order[]>> {
    return this.withLogging('getOrders', { params: params ?? {} }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getOrders(params);
        return response.data.data;
      }, 'Failed to fetch orders');
    });
  }

  /**
   * Get a single order by ID
   */
  async getOrderById(id: number): Promise<ServiceResult<Order>> {
    return this.withLogging('getOrderById', { id }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getOrder(id);
        return response.data.data;
      }, `Failed to fetch order ${id}`);
    });
  }

  // ==================== Order Operations ====================

  /**
   * Cancel an order
   */
  async cancelOrder(id: number, reason?: string): Promise<ServiceResult<Order>> {
    return this.withLogging('cancelOrder', { id, reason }, async () => {
      // First fetch the order to validate
      const orderResult = await this.getOrderById(id);
      if (!orderResult.success || !orderResult.data) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_NOT_FOUND,
          message: `Order ${id} not found`,
        });
      }

      // Check if cancellation is allowed
      const cancellationCheck = this.canCancel(orderResult.data);
      if (!cancellationCheck.allowed) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_CANNOT_CANCEL,
          message: cancellationCheck.reason || 'Order cannot be canceled',
          details: { orderId: id, currentStatus: orderResult.data.status },
        });
      }

      return this.wrapAsync(async () => {
        const response = await this.api.cancelOrder(id, reason);
        return response.data.data;
      }, `Failed to cancel order ${id}`);
    });
  }

  /**
   * Process a refund for an order
   */
  async refundOrder(
    id: number,
    data: { reason: string; amount_cents: number; note?: string }
  ): Promise<ServiceResult<Order>> {
    return this.withLogging('refundOrder', { id, amount: data.amount_cents }, async () => {
      // First fetch the order to validate
      const orderResult = await this.getOrderById(id);
      if (!orderResult.success || !orderResult.data) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_NOT_FOUND,
          message: `Order ${id} not found`,
        });
      }

      const order = orderResult.data;

      // Check if refund is allowed
      const refundCheck = this.canRefund(order);
      if (!refundCheck.allowed) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_CANNOT_REFUND,
          message: refundCheck.reason || 'Order cannot be refunded',
          details: { orderId: id, currentStatus: order.status },
        });
      }

      // Validate refund amount
      if (data.amount_cents <= 0) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_INVALID_REFUND_AMOUNT,
          message: 'Refund amount must be greater than 0',
          details: { requestedAmount: data.amount_cents },
        });
      }

      if (data.amount_cents > order.totalPriceCents) {
        return ServiceResultHelper.failure({
          code: ServiceErrorCodes.ORDER_INVALID_REFUND_AMOUNT,
          message: 'Refund amount cannot exceed original payment amount',
          details: {
            requestedAmount: data.amount_cents,
            originalAmount: order.totalPriceCents,
          },
        });
      }

      return this.wrapAsync(async () => {
        const response = await this.api.refundOrder(id, data);
        return response.data.data;
      }, `Failed to refund order ${id}`);
    });
  }

  // ==================== Batch Operations ====================

  /**
   * Batch cancel orders
   */
  async batchCancel(orderIds: number[], reason?: string): Promise<BatchResult<void>> {
    if (orderIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging('batchCancel', { orderIds, reason }, async () => {
      return this.executeBatch(
        orderIds,
        async (orderId) => {
          // Fetch order first to validate
          const orderResult = await this.getOrderById(orderId);
          if (!orderResult.success || !orderResult.data) {
            throw new Error(`Order ${orderId} not found`);
          }

          const cancellationCheck = this.canCancel(orderResult.data);
          if (!cancellationCheck.allowed) {
            throw new Error(cancellationCheck.reason || 'Order cannot be canceled');
          }

          await this.api.cancelOrder(orderId, reason);
        },
        'batchCancel'
      );
    });
  }

  /**
   * Batch complete orders
   */
  async batchComplete(orderIds: number[]): Promise<BatchResult<void>> {
    if (orderIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging('batchComplete', { orderIds }, async () => {
      return this.executeBatch(
        orderIds,
        async (orderId) => {
          await this.api.batchCompleteOrders([orderId]);
        },
        'batchComplete'
      );
    });
  }

  // ==================== Validation ====================

  /**
   * Check if an order can be canceled
   * Only pending and confirmed orders can be canceled
   */
  canCancel(order: Order): CancellationCheckResult {
    if (CANCELLABLE_STATUSES.includes(order.status)) {
      return { allowed: true };
    }

    const statusMessages: Record<string, string> = {
      in_progress: 'Order is already in progress and cannot be canceled',
      completed: 'Order is already completed and cannot be canceled',
      canceled: 'Order is already canceled',
      refunded: 'Order has been refunded and cannot be canceled',
    };

    return {
      allowed: false,
      reason: statusMessages[order.status] || `Order with status '${order.status}' cannot be canceled`,
    };
  }

  /**
   * Check if an order can be refunded
   * Only completed and in_progress orders can be refunded
   */
  canRefund(order: Order): RefundCheckResult {
    if (REFUNDABLE_STATUSES.includes(order.status)) {
      return { allowed: true };
    }

    const statusMessages: Record<string, string> = {
      pending: 'Order is still pending and should be canceled instead of refunded',
      confirmed: 'Order is confirmed but not started, should be canceled instead',
      canceled: 'Order is already canceled',
      refunded: 'Order has already been refunded',
    };

    return {
      allowed: false,
      reason: statusMessages[order.status] || `Order with status '${order.status}' cannot be refunded`,
    };
  }

  /**
   * Calculate refund breakdown
   * Ensures refund amount does not exceed original payment
   * Calculates platform fee and player amount
   */
  calculateRefund(order: Order, requestedAmount: number): RefundCalculation {
    const originalAmount = order.totalPriceCents;
    
    // Clamp refund amount to not exceed original
    const refundAmount = Math.min(Math.max(0, requestedAmount), originalAmount);
    
    // Calculate platform fee (what platform keeps from the refund)
    // Platform fee is proportional to the refund amount
    const platformFee = Math.round(refundAmount * DEFAULT_PLATFORM_COMMISSION_RATE);
    
    // Player amount is what gets deducted from player's earnings
    const playerAmount = refundAmount - platformFee;

    return {
      originalAmount,
      refundAmount,
      platformFee,
      playerAmount,
    };
  }

  // ==================== Statistics ====================

  /**
   * Compute statistics from a list of orders
   */
  computeStatistics(orders: Order[]): OrderStatistics {
    if (orders.length === 0) {
      return {
        totalOrders: 0,
        totalRevenue: 0,
        ordersByStatus: {},
        averageOrderValue: 0,
        completionRate: 0,
      };
    }

    // Count orders by status
    const ordersByStatus: Record<string, number> = {};
    let totalRevenue = 0;
    let completedCount = 0;

    for (const order of orders) {
      // Count by status
      ordersByStatus[order.status] = (ordersByStatus[order.status] || 0) + 1;

      // Sum revenue (only from completed orders)
      if (order.status === 'completed') {
        totalRevenue += order.totalPriceCents;
        completedCount++;
      }
    }

    const totalOrders = orders.length;
    const averageOrderValue = completedCount > 0 ? Math.round(totalRevenue / completedCount) : 0;
    const completionRate = totalOrders > 0 ? completedCount / totalOrders : 0;

    return {
      totalOrders,
      totalRevenue,
      ordersByStatus,
      averageOrderValue,
      completionRate,
    };
  }

  /**
   * Compute trend data from orders over a specified number of days
   */
  computeTrend(orders: Order[], days: number): TrendData[] {
    const trendMap = new Map<string, number>();
    const now = new Date();

    // Initialize all days with 0
    for (let i = days - 1; i >= 0; i--) {
      const date = new Date(now);
      date.setDate(date.getDate() - i);
      const dateStr = this.formatDateForTrend(date);
      trendMap.set(dateStr, 0);
    }

    // Count orders per day
    for (const order of orders) {
      const orderDate = new Date(order.createdAt);
      const dateStr = this.formatDateForTrend(orderDate);
      
      if (trendMap.has(dateStr)) {
        trendMap.set(dateStr, (trendMap.get(dateStr) || 0) + 1);
      }
    }

    // Convert to array
    const trend: TrendData[] = [];
    for (const [date, value] of trendMap) {
      trend.push({ date, value });
    }

    return trend;
  }

  // ==================== Private Helpers ====================

  /**
   * Format date for trend data (YYYY-MM-DD)
   */
  private formatDateForTrend(date: Date): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
}

/**
 * Default OrderService instance
 */
export const orderService = new OrderService();
