import apiClient from './client';
import type { ApiResponse } from './admin';

// ============================================================================
// VIP Level Types
// ============================================================================

interface VIPLevel {
    id: number;
    slug: string;
    title: string;
    expRequired: number;
    orderDiscount: number;
    monthlyCouponTemplateId?: number;
    monthlyCouponCount: number;
    iconUrl: string;
    color: string;
    benefits: string;
    sortOrder: number;
    isDefault: boolean;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
}

interface CreateVIPLevelDto {
    slug: string;
    title: string;
    expRequired?: number;
    orderDiscount?: number;
    monthlyCouponTemplateId?: number;
    monthlyCouponCount?: number;
    iconUrl?: string;
    color?: string;
    benefits?: string;
    sortOrder?: number;
    isDefault?: boolean;
    isActive?: boolean;
}

interface UpdateVIPLevelDto {
    slug: string;
    title: string;
    expRequired?: number;
    orderDiscount?: number;
    monthlyCouponTemplateId?: number;
    monthlyCouponCount?: number;
    iconUrl?: string;
    color?: string;
    benefits?: string;
    sortOrder?: number;
    isDefault?: boolean;
    isActive?: boolean;
}

interface VIPLevelQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    isActive?: boolean;
}

// ============================================================================
// VIP Config Types
// ============================================================================

interface VIPConfig {
    id: number;
    configKey: string;
    configValue: string;
    description: string;
    createdAt: string;
    updatedAt: string;
}

interface CreateVIPConfigDto {
    configKey: string;
    configValue: string;
    description?: string;
}

interface UpdateVIPConfigDto {
    configValue: string;
    description?: string;
}

// ============================================================================
// VIP Config Key Constants
// ============================================================================

export const VIP_CONFIG_KEYS = {
    UNLOCK_BY_CONSUME: 'unlock_by_consume',  // 累计消费解锁门槛（分）
    UNLOCK_BY_RECHARGE: 'unlock_by_recharge', // 累计充值解锁门槛（分）
    EXPIRE_DAYS: 'expire_days',              // VIP过期天数（0=永久）
} as const;

export type VIPConfigKey = typeof VIP_CONFIG_KEYS[keyof typeof VIP_CONFIG_KEYS];

// ============================================================================
// Batch Operation Types
// ============================================================================

interface VIPBatchUpdateStatusDto {
    ids: number[];
    isActive: boolean;
}

interface VIPBatchDeleteDto {
    ids: number[];
}

// ============================================================================
// VIP API
// ============================================================================

export const vipApi = {
    // ------------------------------------------------------------------------
    // VIP Level Management
    // ------------------------------------------------------------------------

    /**
     * Get VIP levels list with pagination
     * GET /api/v1/admin/vip/levels
     */
    getVIPLevels: (params?: VIPLevelQueryParams) =>
        apiClient.get<ApiResponse<VIPLevel[]>>('/admin/vip/levels', { params }),

    /**
     * Get VIP level detail by ID
     * GET /api/v1/admin/vip/levels/:id
     */
    getVIPLevelDetail: (id: number) =>
        apiClient.get<ApiResponse<VIPLevel>>(`/admin/vip/levels/${id}`),

    /**
     * Create a new VIP level
     * POST /api/v1/admin/vip/levels
     */
    createVIPLevel: (data: CreateVIPLevelDto) =>
        apiClient.post<ApiResponse<VIPLevel>>('/admin/vip/levels', data),

    /**
     * Update VIP level by ID
     * PUT /api/v1/admin/vip/levels/:id
     */
    updateVIPLevel: (id: number, data: UpdateVIPLevelDto) =>
        apiClient.put<ApiResponse<VIPLevel>>(`/admin/vip/levels/${id}`, data),

    /**
     * Delete VIP level by ID
     * DELETE /api/v1/admin/vip/levels/:id
     */
    deleteVIPLevel: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/vip/levels/${id}`),

    /**
     * Set default VIP level
     * POST /api/v1/admin/vip/levels/:id/default
     */
    setDefaultVIPLevel: (id: number) =>
        apiClient.post<ApiResponse<void>>(`/admin/vip/levels/${id}/default`),

    /**
     * Batch update VIP level status
     * POST /api/v1/admin/vip/levels/batch-status
     */
    batchUpdateVIPLevelStatus: (data: VIPBatchUpdateStatusDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/vip/levels/batch-status', data),

    /**
     * Batch delete VIP levels
     * POST /api/v1/admin/vip/levels/batch-delete
     */
    batchDeleteVIPLevels: (data: VIPBatchDeleteDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/vip/levels/batch-delete', data),

    // ------------------------------------------------------------------------
    // VIP Config Management
    // ------------------------------------------------------------------------

    /**
     * Get all VIP configs
     * GET /api/v1/admin/vip/configs
     */
    getVIPConfigs: () =>
        apiClient.get<ApiResponse<VIPConfig[]>>('/admin/vip/configs'),

    /**
     * Get VIP config by key
     * GET /api/v1/admin/vip/configs/:key
     */
    getVIPConfig: (key: string) =>
        apiClient.get<ApiResponse<VIPConfig>>(`/admin/vip/configs/${key}`),

    /**
     * Create or update VIP config
     * POST /api/v1/admin/vip/configs
     */
    saveVIPConfig: (data: CreateVIPConfigDto) =>
        apiClient.post<ApiResponse<VIPConfig>>('/admin/vip/configs', data),

    /**
     * Update VIP config by key
     * PUT /api/v1/admin/vip/configs/:key
     */
    updateVIPConfig: (key: string, data: UpdateVIPConfigDto) =>
        apiClient.put<ApiResponse<VIPConfig>>(`/admin/vip/configs/${key}`, data),

    /**
     * Delete VIP config by key
     * DELETE /api/v1/admin/vip/configs/:key
     */
    deleteVIPConfig: (key: string) =>
        apiClient.delete<ApiResponse<void>>(`/admin/vip/configs/${key}`),
};

// ============================================================================
// Re-exports
// ============================================================================

export type {
    VIPLevel,
    CreateVIPLevelDto,
    UpdateVIPLevelDto,
    VIPLevelQueryParams,
    VIPConfig,
    CreateVIPConfigDto,
    UpdateVIPConfigDto,
    VIPBatchUpdateStatusDto,
    VIPBatchDeleteDto,
};
