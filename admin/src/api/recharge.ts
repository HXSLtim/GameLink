import apiClient from './client';
import type { ApiResponse, Pagination } from './admin';

// ============================================================================
// Recharge Option Types (充值档位)
// ============================================================================

/**
 * Recharge option (充值档位)
 */
export interface RechargeOption {
    id: number;
    name: string;
    amountCents: number;
    bonusCents: number;
    originalCents?: number;
    discountPercent?: number;
    description?: string;
    tag?: string;
    iconUrl?: string;
    sortOrder: number;
    isActive: boolean;
    isRecommended: boolean;
    couponTemplateId?: number;
    couponCount: number;
    minVipLevel?: number;
    perUserLimit: number;
    totalLimit: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * Create recharge option DTO
 */
export interface CreateRechargeOptionDto {
    name: string;
    amountCents: number;
    bonusCents?: number;
    originalCents?: number;
    discountPercent?: number;
    description?: string;
    tag?: string;
    iconUrl?: string;
    sortOrder?: number;
    isActive?: boolean;
    isRecommended?: boolean;
    couponTemplateId?: number;
    couponCount?: number;
    minVipLevel?: number;
    perUserLimit?: number;
    totalLimit?: number;
}

/**
 * Update recharge option DTO
 */
export interface UpdateRechargeOptionDto {
    name: string;
    amountCents: number;
    bonusCents?: number;
    originalCents?: number;
    discountPercent?: number;
    description?: string;
    tag?: string;
    iconUrl?: string;
    sortOrder?: number;
    isActive?: boolean;
    isRecommended?: boolean;
    couponTemplateId?: number;
    couponCount?: number;
    minVipLevel?: number;
    perUserLimit?: number;
    totalLimit?: number;
}

/**
 * Recharge option query parameters
 */
export interface RechargeOptionQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    isActive?: boolean;
    isRecommended?: boolean;
}

/**
 * Batch update status DTO
 */
export interface BatchUpdateOptionStatusDto {
    ids: number[];
    isActive: boolean;
}

/**
 * Batch delete DTO
 */
export interface BatchDeleteOptionDto {
    ids: number[];
}

// ============================================================================
// Recharge Record Types (充值记录)
// ============================================================================

/**
 * Recharge status enum
 */
export type RechargeStatus = 'pending' | 'paid' | 'failed' | 'refunded' | 'expired';

/**
 * Recharge record (充值记录)
 */
export interface RechargeRecord {
    id: number;
    orderNo: string;
    userId: number;
    optionId: number;
    amountCents: number;
    bonusCents: number;
    totalCents: number;
    status: RechargeStatus;
    paymentChannel?: string;
    paymentNo?: string;
    paidAt?: string;
    refundedAt?: string;
    refundReason?: string;
    createdAt: string;
    updatedAt: string;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    option?: {
        id: number;
        name: string;
        amountCents: number;
        bonusCents: number;
    };
}

/**
 * Recharge record query parameters
 */
export interface RechargeRecordQueryParams {
    page?: number;
    page_size?: number;
    userId?: number;
    optionId?: number;
    status?: RechargeStatus;
    paymentChannel?: string;
    orderNo?: string;
    startTime?: string;
    endTime?: string;
}

/**
 * Refund request DTO
 */
export interface RefundRecordDto {
    reason: string;
}

// ============================================================================
// Recharge Statistics Types (充值统计)
// ============================================================================

/**
 * Recharge statistics
 */
export interface RechargeStats {
    totalOrders: number;
    totalAmountCents: number;
    totalBonusCents: number;
    paidOrders: number;
    pendingOrders: number;
    failedOrders: number;
    refundedOrders: number;
    todayOrders: number;
    todayAmountCents: number;
    monthOrders: number;
    monthAmountCents: number;
}

// ============================================================================
// Recharge API
// ============================================================================

export const rechargeApi = {
    // ------------------------------------------------------------------------
    // Recharge Option Management (充值档位管理)
    // ------------------------------------------------------------------------

    /**
     * Get recharge options list with pagination
     * GET /api/v1/admin/recharge/options
     */
    getRechargeOptions: (params?: RechargeOptionQueryParams) =>
        apiClient.get<ApiResponse<RechargeOption[]>>('/admin/recharge/options', { params }),

    /**
     * Get recharge option detail by ID
     * GET /api/v1/admin/recharge/options/:id
     */
    getRechargeOptionDetail: (id: number) =>
        apiClient.get<ApiResponse<RechargeOption>>(`/admin/recharge/options/${id}`),

    /**
     * Create a new recharge option
     * POST /api/v1/admin/recharge/options
     */
    createRechargeOption: (data: CreateRechargeOptionDto) =>
        apiClient.post<ApiResponse<RechargeOption>>('/admin/recharge/options', data),

    /**
     * Update recharge option by ID
     * PUT /api/v1/admin/recharge/options/:id
     */
    updateRechargeOption: (id: number, data: UpdateRechargeOptionDto) =>
        apiClient.put<ApiResponse<RechargeOption>>(`/admin/recharge/options/${id}`, data),

    /**
     * Delete recharge option by ID
     * DELETE /api/v1/admin/recharge/options/:id
     */
    deleteRechargeOption: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/recharge/options/${id}`),

    /**
     * Toggle recharge option status (enable/disable)
     * POST /api/v1/admin/recharge/options/batch-status
     */
    toggleRechargeOptionStatus: (ids: number[], isActive: boolean) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/recharge/options/batch-status', {
            ids,
            isActive,
        }),

    /**
     * Batch update recharge option status
     * POST /api/v1/admin/recharge/options/batch-status
     */
    batchUpdateOptionStatus: (data: BatchUpdateOptionStatusDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/recharge/options/batch-status', data),

    /**
     * Batch delete recharge options
     * POST /api/v1/admin/recharge/options/batch-delete
     */
    batchDeleteOptions: (data: BatchDeleteOptionDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/recharge/options/batch-delete', data),

    // ------------------------------------------------------------------------
    // Recharge Record Management (充值记录管理)
    // ------------------------------------------------------------------------

    /**
     * Get recharge records list with pagination
     * GET /api/v1/admin/recharge/records
     */
    getRechargeOrders: (params?: RechargeRecordQueryParams) =>
        apiClient.get<ApiResponse<RechargeRecord[]>>('/admin/recharge/records', { params }),

    /**
     * Get recharge record detail by ID
     * GET /api/v1/admin/recharge/records/:id
     */
    getRechargeRecordDetail: (id: number) =>
        apiClient.get<ApiResponse<RechargeRecord>>(`/admin/recharge/records/${id}`),

    /**
     * Refund recharge record
     * POST /api/v1/admin/recharge/records/:id/refund
     */
    refundRechargeRecord: (id: number, data: RefundRecordDto) =>
        apiClient.post<ApiResponse<{ message: string }>>(`/admin/recharge/records/${id}/refund`, data),

    // ------------------------------------------------------------------------
    // Recharge Statistics (充值统计)
    // ------------------------------------------------------------------------

    /**
     * Get recharge statistics
     * GET /api/v1/admin/recharge/stats
     */
    getRechargeStats: () =>
        apiClient.get<ApiResponse<RechargeStats>>('/admin/recharge/stats'),
};

// ============================================================================
// Re-exports
// ============================================================================

export type {
    RechargeOption,
    CreateRechargeOptionDto,
    UpdateRechargeOptionDto,
    RechargeOptionQueryParams,
    BatchUpdateOptionStatusDto,
    BatchDeleteOptionDto,
    RechargeRecord,
    RechargeRecordQueryParams,
    RefundRecordDto,
    RechargeStats,
};
