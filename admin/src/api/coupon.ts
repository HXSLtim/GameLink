import apiClient from './client';
import type { ApiResponse } from '@/types/api';

// Re-export for backward compatibility
export type { ApiResponse };

// ============================================================================
// Type Definitions
// ============================================================================

/**
 * Coupon Types
 */
export type CouponType = 'deduct' | 'discount';

/**
 * Coupon Scope (Applicability)
 */
export type CouponScope = 'all' | 'game' | 'item';

/**
 * Coupon Source
 */
export type CouponSource = 'new_user' | 'link' | 'vip' | 'activity' | 'manual' | 'referral' | 'team';

/**
 * Coupon State (Status)
 */
export type CouponState = 'available' | 'locked' | 'used' | 'expired';

/**
 * Validity Type
 */
export type ValidityType = 'days' | 'fixed';

/**
 * Coupon Template Interface
 * Matches backend model.CouponTemplate
 */
export interface CouponTemplate {
    id: number;
    name: string;
    type: CouponType;
    source: CouponSource;
    description?: string;

    // Discount configuration
    minAmountCents: number;
    deductAmountCents: number;
    discountRate: number;
    maxDiscountCents: number;

    // Applicability
    scope: CouponScope;
    gameIds: string; // JSON array string
    itemIds: string; // JSON array string

    // Validity
    validityType: ValidityType;
    validityDays: number;
    fixedExpireAt?: string;

    // Claim configuration
    totalCount: number;
    claimedCount: number;
    perUserLimit: number;
    claimLink?: string;

    // Status
    isActive: boolean;

    // Timestamps
    createdAt: string;
    updatedAt: string;
}

/**
 * User Coupon Interface
 * Matches backend model.Coupon
 */
export interface Coupon {
    id: number;
    templateId: number;
    userId: number;
    state: CouponState;

    // Denormalized fields from template
    name: string;
    type: CouponType;
    source: CouponSource;
    minAmountCents: number;
    deductAmountCents: number;
    discountRate: number;
    maxDiscountCents: number;
    scope: CouponScope;
    gameIds: string;
    itemIds: string;

    // Time fields
    claimedAt?: string;
    expireAt: string;
    usedAt?: string;

    // Lock information
    lockedByOrderId?: number;
    lockedAt?: string;

    // Usage information
    usedOrderId?: number;
    discountCents: number;

    // Timestamps
    createdAt: string;
    updatedAt: string;
}

/**
 * Coupon with Template details
 */
export interface CouponWithTemplate extends Coupon {
    template?: CouponTemplate;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
}

/**
 * Coupon Statistics
 */
export interface CouponStats {
    totalTemplates: number;
    activeTemplates: number;
    totalCoupons: number;
    availableCoupons: number;
    usedCoupons: number;
    expiredCoupons: number;
    totalDiscountCents: number;
}

// ============================================================================
// Request DTOs
// ============================================================================

/**
 * Create Template Request
 */
export interface CreateTemplateDto {
    name: string;
    type: CouponType;
    source: CouponSource;
    description?: string;
    minAmountCents?: number;
    deductAmountCents?: number;
    discountRate?: number;
    maxDiscountCents?: number;
    scope?: CouponScope;
    gameIds?: string;
    itemIds?: string;
    validityType?: ValidityType;
    validityDays?: number;
    fixedExpireAt?: string;
    totalCount?: number;
    perUserLimit?: number;
    claimLink?: string;
    isActive?: boolean;
}

/**
 * Update Template Request
 */
export type UpdateTemplateDto = Partial<CreateTemplateDto>;

/**
 * Template Query Parameters
 */
export interface TemplateQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    type?: CouponType;
    source?: CouponSource;
    isActive?: boolean;
}

/**
 * Coupon Query Parameters
 */
export interface CouponQueryParams {
    page?: number;
    page_size?: number;
    userId?: number;
    templateId?: number;
    state?: CouponState;
    type?: CouponType;
}

/**
 * Issue Coupon Request
 */
export interface IssueCouponDto {
    userId: number;
    templateId: number;
    source?: CouponSource;
}

/**
 * Batch Update Template Status Request
 */
export interface BatchUpdateTemplateStatusDto {
    ids: number[];
    isActive: boolean;
}

/**
 * Batch Delete Templates Request
 */
export interface BatchDeleteTemplatesDto {
    ids: number[];
}

// ============================================================================
// API Client
// ============================================================================

/**
 * Coupon API
 * Provides methods for managing coupon templates and user coupons
 */
export const couponApi = {
    // ========================================================================
    // Template Management
    // ========================================================================

    /**
     * Get coupon templates list
     * GET /admin/coupons/templates
     */
    getTemplates: (params?: TemplateQueryParams) =>
        apiClient.get<ApiResponse<CouponTemplate[]>>('/admin/coupons/templates', { params }),

    /**
     * Get coupon template detail
     * GET /admin/coupons/templates/:id
     */
    getTemplate: (id: number) =>
        apiClient.get<ApiResponse<CouponTemplate>>(`/admin/coupons/templates/${id}`),

    /**
     * Create new coupon template
     * POST /admin/coupons/templates
     */
    createTemplate: (data: CreateTemplateDto) =>
        apiClient.post<ApiResponse<CouponTemplate>>('/admin/coupons/templates', data),

    /**
     * Update coupon template
     * PUT /admin/coupons/templates/:id
     */
    updateTemplate: (id: number, data: UpdateTemplateDto) =>
        apiClient.put<ApiResponse<CouponTemplate>>(`/admin/coupons/templates/${id}`, data),

    /**
     * Delete coupon template
     * DELETE /admin/coupons/templates/:id
     */
    deleteTemplate: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/coupons/templates/${id}`),

    /**
     * Batch update template status
     * POST /admin/coupons/templates/batch-status
     */
    batchUpdateTemplateStatus: (data: BatchUpdateTemplateStatusDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/coupons/templates/batch-status', data),

    /**
     * Batch delete templates
     * POST /admin/coupons/templates/batch-delete
     */
    batchDeleteTemplates: (data: BatchDeleteTemplatesDto) =>
        apiClient.post<ApiResponse<{ affected: number }>>('/admin/coupons/templates/batch-delete', data),

    // ========================================================================
    // User Coupon Management
    // ========================================================================

    /**
     * Get user coupons list
     * GET /admin/coupons
     */
    getCoupons: (params?: CouponQueryParams) =>
        apiClient.get<ApiResponse<CouponWithTemplate[]>>('/admin/coupons', { params }),

    /**
     * Get coupon detail with template
     * GET /admin/coupons/:id
     */
    getCouponDetail: (id: number) =>
        apiClient.get<ApiResponse<CouponWithTemplate>>(`/admin/coupons/${id}`),

    /**
     * Delete coupon
     * DELETE /admin/coupons/:id
     */
    deleteCoupon: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/coupons/${id}`),

    /**
     * Issue coupon to user
     * POST /admin/coupons/issue
     */
    issueCoupon: (data: IssueCouponDto) =>
        apiClient.post<ApiResponse<Coupon>>('/admin/coupons/issue', data),

    /**
     * Get coupon statistics
     * GET /admin/coupons/stats
     */
    getCouponStats: () =>
        apiClient.get<ApiResponse<CouponStats>>('/admin/coupons/stats'),

    // ========================================================================
    // Alias methods for backward compatibility
    // ========================================================================

    // Aliases for template methods
    getCouponTemplates: (params?: TemplateQueryParams) => couponApi.getTemplates(params),
    getCouponTemplateDetail: (id: number) => couponApi.getTemplate(id),
    createCoupon: (data: CreateTemplateDto) => couponApi.createTemplate(data),
    updateCoupon: (id: number, data: UpdateTemplateDto) => couponApi.updateTemplate(id, data),

    /**
     * Toggle coupon template active status
     * Convenience method for enabling/disabling templates
     */
    toggleCouponTemplate: (id: number, isActive: boolean) =>
        couponApi.updateTemplate(id, { isActive }),

    /**
     * Get coupon usage statistics
     * Alias for getCouponStats
     */
    getCouponUsage: () => couponApi.getCouponStats(),

    /**
     * Batch update coupon template status
     * Alias for batchUpdateTemplateStatus
     */
    batchUpdateCouponStatus: (ids: number[], isActive: boolean) =>
        couponApi.batchUpdateTemplateStatus({ ids, isActive }),

    /**
     * Batch delete coupon templates
     * Alias for batchDeleteTemplates
     */
    batchDeleteCoupons: (ids: number[]) =>
        couponApi.batchDeleteTemplates({ ids }),
};

// ============================================================================
// Helper Functions
// ============================================================================

/**
 * Format cents to yuan
 */
export const centsToYuan = (cents: number): string => {
    return (cents / 100).toFixed(2);
};

/**
 * Format yuan to cents
 */
export const yuanToCents = (yuan: number): number => {
    return Math.round(yuan * 100);
};

/**
 * Parse JSON string to array
 */
export const parseJsonArray = (jsonString: string): number[] => {
    try {
        const parsed = JSON.parse(jsonString || '[]');
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return [];
    }
};

/**
 * Format coupon type display name
 */
export const getCouponTypeLabel = (type: CouponType): string => {
    const labels: Record<CouponType, string> = {
        deduct: '满减券',
        discount: '折扣券',
    };
    return labels[type] || type;
};

/**
 * Format coupon scope display name
 */
export const getCouponScopeLabel = (scope: CouponScope): string => {
    const labels: Record<CouponScope, string> = {
        all: '通用',
        game: '指定游戏',
        item: '指定服务',
    };
    return labels[scope] || scope;
};

/**
 * Format coupon source display name
 */
export const getCouponSourceLabel = (source: CouponSource): string => {
    const labels: Record<CouponSource, string> = {
        new_user: '新用户',
        link: '链接领取',
        vip: 'VIP会员',
        activity: '活动发放',
        manual: '手动发放',
        referral: '推荐奖励',
        team: '团队奖励',
    };
    return labels[source] || source;
};

/**
 * Format coupon state display name
 */
export const getCouponStateLabel = (state: CouponState): string => {
    const labels: Record<CouponState, string> = {
        available: '可用',
        locked: '已锁定',
        used: '已使用',
        expired: '已过期',
    };
    return labels[state] || state;
};

/**
 * Format coupon state color (for Ant Design tags)
 */
export const getCouponStateColor = (state: CouponState): string => {
    const colors: Record<CouponState, string> = {
        available: 'green',
        locked: 'orange',
        used: 'blue',
        expired: 'default',
    };
    return colors[state] || 'default';
};

/**
 * Calculate discount amount
 * For deduct coupon: return deductAmountCents
 * For discount coupon: calculate based on discountRate
 */
export const calculateDiscount = (
    coupon: Coupon | CouponTemplate,
    orderAmountCents: number
): number => {
    if (coupon.type === 'deduct') {
        return coupon.deductAmountCents;
    }

    if (coupon.type === 'discount') {
        const discount = Math.round(orderAmountCents * (1 - coupon.discountRate));
        if (coupon.maxDiscountCents > 0) {
            return Math.min(discount, coupon.maxDiscountCents);
        }
        return discount;
    }

    return 0;
};

/**
 * Check if coupon is applicable to order
 */
export const isCouponApplicable = (
    coupon: Coupon | CouponTemplate,
    orderAmountCents: number,
    gameId?: number,
    itemIds?: number[]
): { applicable: boolean; reason?: string } => {
    // Check minimum amount
    if (coupon.minAmountCents > 0 && orderAmountCents < coupon.minAmountCents) {
        return {
            applicable: false,
            reason: `订单金额不足，最低需要 ¥${centsToYuan(coupon.minAmountCents)}`,
        };
    }

    // Check scope
    if (coupon.scope === 'game' && gameId) {
        const applicableGameIds = parseJsonArray(coupon.gameIds);
        if (!applicableGameIds.includes(gameId)) {
            return { applicable: false, reason: '该优惠券不适用于此游戏' };
        }
    }

    if (coupon.scope === 'item' && itemIds) {
        const applicableItemIds = parseJsonArray(coupon.itemIds);
        const hasApplicableItem = itemIds.some(id => applicableItemIds.includes(id));
        if (!hasApplicableItem) {
            return { applicable: false, reason: '该优惠券不适用于此服务项目' };
        }
    }

    return { applicable: true };
};

// ============================================================================
// Default export
// ============================================================================

export default couponApi;
