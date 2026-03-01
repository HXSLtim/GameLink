import apiClient from './client';
import type { ApiResponse } from './admin';

/**
 * Reconciliation Types
 */

export type ReconciliationType = 'payment' | 'internal' | 'bank' | 'manual';
export type ReconciliationStatus = 'pending' | 'progress' | 'success' | 'failed' | 'exception';

/**
 * 对账单数据结构
 */
export interface Reconciliation {
    id: number;
    reconciliationNo: string;
    reconciliationDate: string;
    type: ReconciliationType;
    status: ReconciliationStatus;
    periodStart: string;
    periodEnd: string;
    totalRecords: number;
    matchedRecords: number;
    differenceAmount: number;
    abstract?: string;
    processedAt?: string;
    processedBy?: number;
    createdAt: string;
    updatedAt: string;
    details?: ReconciliationDetail[];
}

/**
 * 对账明细数据结构
 */
export interface ReconciliationDetail {
    id: number;
    reconciliationId: number;
    lineNo: number;
    externalType: string;
    externalNo: string;
    externalAmount: number;
    externalDate: string;
    internalType: string;
    internalNo: string;
    internalAmount: number;
    internalDate: string;
    status: string;
    differenceAmount: number;
    remark?: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * 对账单列表查询参数
 */
export interface ReconciliationListParams {
    page?: number;
    pageSize?: number;
    type?: ReconciliationType;
    status?: ReconciliationStatus;
    date_from?: string;
    date_to?: string;
}

/**
 * 创建对账单请求
 */
export interface CreateReconciliationDto {
    reconciliationNo?: string;
    reconciliationDate: string;
    type: ReconciliationType;
    periodStart: string;
    periodEnd: string;
    abstract?: string;
    details: CreateReconciliationDetailDto[];
}

/**
 * 创建对账明细请求
 */
export interface CreateReconciliationDetailDto {
    externalType: string;
    externalNo: string;
    externalAmount?: number;
    externalDate: string;
    internalType: string;
    internalNo: string;
    internalAmount?: number;
    internalDate: string;
    remark?: string;
}

/**
 * 执行对账请求
 */
export interface ExecuteReconciliationDto {
    status?: ReconciliationStatus;
}

/**
 * 对账统计信息
 */
export interface ReconciliationStats {
    total: number;
    pending: number;
    success: number;
    failed: number;
    exception: number;
}

/**
 * Reconciliation API
 */
export const reconciliationApi = {
    /**
     * Get reconciliation list
     * GET /admin/reconciliations
     */
    getReconciliations: (params?: ReconciliationListParams) =>
        apiClient.get<ApiResponse<Reconciliation[]>>('/admin/reconciliations', { params }),

    /**
     * Get reconciliation detail
     * GET /admin/reconciliations/:id
     */
    getReconciliationDetail: (id: number) =>
        apiClient.get<ApiResponse<Reconciliation>>(`/admin/reconciliations/${id}`),

    /**
     * Create reconciliation
     * POST /admin/reconciliations
     */
    createReconciliation: (data: CreateReconciliationDto) =>
        apiClient.post<ApiResponse<Reconciliation>>('/admin/reconciliations', data),

    /**
     * Execute reconciliation
     * POST /admin/reconciliations/:id/execute
     */
    executeReconciliation: (id: number, data: ExecuteReconciliationDto) =>
        apiClient.post<ApiResponse<Reconciliation>>(`/admin/reconciliations/${id}/execute`, data),
};

// Re-export types for convenience
export type {
    Reconciliation,
    ReconciliationDetail,
    ReconciliationListParams,
    CreateReconciliationDto,
    CreateReconciliationDetailDto,
    ExecuteReconciliationDto,
    ReconciliationStats,
};
