import apiClient from './client';
import type { ApiResponse } from '@/types/api';

// Re-export for backward compatibility
export type { ApiResponse };

// ============================================================================
// Type Definitions
// ============================================================================

/**
 * Referral Type
 * Matches backend model.ReferralType
 */
export type ReferralType = 'user' | 'player';

/**
 * Referral Status
 * Matches backend model.ReferralStatus
 */
export type ReferralStatus = 'pending' | 'completed' | 'cancelled';

/**
 * Reward Type
 * Matches backend model.RewardType
 */
export type RewardType = 'referrer' | 'referee';

/**
 * Referral Reward Status
 * Matches backend model.ReferralRewardStatus
 */
export type ReferralRewardStatus = 'pending' | 'issued' | 'failed';

/**
 * Referral Config
 * Configuration for referral system
 */
export interface ReferralConfig {
    id: number;
    configKey: string;
    configValue: string;
    description: string;
    createdAt: string;
    updatedAt: string;
}

/**
 * Update Referral Config Request
 */
export interface UpdateReferralConfigDto {
    value: string;
    description?: string;
}

// ============================================================================
// Referral Config Key Constants
// ============================================================================

export const REFERRAL_CONFIG_KEYS = {
    ENABLED: 'referral_enabled',              // 推荐功能开关
    REFERRER_REWARD_CENTS: 'referrer_reward_cents',  // 推荐人奖励（分）
    REFEREE_REWARD_CENTS: 'referee_reward_cents',    // 被推荐人奖励（分）
    MAX_REWARD_USES: 'max_reward_uses',       // 每个邀请码最大使用次数
    EXPIRE_DAYS: 'expire_days',               // 邀请码有效期（天）
} as const;

export type ReferralConfigKey = typeof REFERRAL_CONFIG_KEYS[keyof typeof REFERRAL_CONFIG_KEYS];

// ============================================================================
// Referral Code Types
// ============================================================================

/**
 * Referral Code Interface
 * Matches backend model.ReferralCode
 */
export interface ReferralCode {
    id: number;
    code: string;
    ownerId: number;
    type: ReferralType;
    maxUses: number;
    usedCount: number;
    expiresAt: string;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
    owner?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
}

/**
 * Create Referral Code Request
 */
export interface CreateReferralCodeDto {
    userId: number;
    type: ReferralType;
    maxUse?: number;
    expireAt?: string;
}

/**
 * Update Referral Code Request
 */
export interface UpdateReferralCodeDto {
    isActive?: boolean;
    maxUse?: number;
    expireAt?: string;
}

/**
 * Referral Code Query Parameters
 */
export interface ReferralCodeQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    userId?: number;
    type?: ReferralType;
    isActive?: boolean;
}

/**
 * Batch Update Referral Codes Status Request
 */
export interface BatchUpdateCodesStatusDto {
    ids: number[];
    isActive: boolean;
}

/**
 * Batch Delete Referral Codes Request
 */
export interface BatchDeleteCodesDto {
    ids: number[];
}

// ============================================================================
// Referral Types
// ============================================================================

/**
 * Referral Interface
 * Matches backend model.Referral
 */
export interface Referral {
    id: number;
    referrerId: number;
    refereeId: number;
    codeId: number;
    type: ReferralType;
    status: ReferralStatus;
    completedAt?: string;
    cancelledAt?: string;
    cancelReason?: string;
    createdAt: string;
    updatedAt: string;
    referrer?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    referee?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    code?: ReferralCode;
}

/**
 * Update Referral Status Request
 */
export interface UpdateReferralStatusDto {
    status: ReferralStatus;
}

/**
 * Referral Query Parameters
 */
export interface ReferralQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    referrerId?: number;
    refereeId?: number;
    type?: ReferralType;
    status?: ReferralStatus;
}

/**
 * Batch Update Referrals Status Request
 */
export interface BatchUpdateReferralsStatusDto {
    ids: number[];
    status: ReferralStatus;
}

/**
 * Batch Delete Referrals Request
 */
export interface BatchDeleteReferralsDto {
    ids: number[];
}

// ============================================================================
// Referral Reward Types
// ============================================================================

/**
 * Referral Reward Interface
 * Matches backend model.ReferralReward
 */
export interface ReferralReward {
    id: number;
    referralId: number;
    userId: number;
    type: RewardType;
    amountCents: number;
    status: ReferralRewardStatus;
    issuedAt?: string;
    failedAt?: string;
    failureReason?: string;
    createdAt: string;
    updatedAt: string;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    referral?: Referral;
}

/**
 * Issue Referral Reward Request
 */
export type IssueReferralRewardDto = Record<string, never>;

/**
 * Fail Referral Reward Request
 */
export interface FailReferralRewardDto {
    reason: string;
}

/**
 * Referral Reward Query Parameters
 */
export interface ReferralRewardQueryParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    userId?: number;
    referralId?: number;
    type?: RewardType;
    status?: ReferralRewardStatus;
}

// ============================================================================
// Referral Statistics Types
// ============================================================================

/**
 * Referral Statistics
 */
export interface ReferralStats {
    totalReferrals: number;
    completedReferrals: number;
    pendingReferrals: number;
    cancelledReferrals: number;
    totalRewardsCents: number;
    issuedRewardsCents: number;
    pendingRewardsCents: number;
    failedRewardsCents: number;
    activeCodes: number;
    totalCodes: number;
}

// ============================================================================
// API Client
// ============================================================================

/**
 * Referral API
 * Provides methods for managing referral system including configs, codes, referrals, and rewards
 */
export const referralApi = {
    // ========================================================================
    // Config Management
    // ========================================================================

    /**
     * Get all referral configs
     * GET /api/v1/admin/referrals/configs
     */
    getReferralConfigs: () =>
        apiClient.get<ApiResponse<ReferralConfig[]>>('/admin/referrals/configs'),

    /**
     * Update referral config by key
     * PUT /api/v1/admin/referrals/configs/:key
     */
    updateReferralConfig: (key: string, data: UpdateReferralConfigDto) =>
        apiClient.put<ApiResponse<ReferralConfig>>(`/admin/referrals/configs/${key}`, data),

    // ========================================================================
    // Referral Code Management
    // ========================================================================

    /**
     * Get referral codes list with pagination
     * GET /api/v1/admin/referrals/codes
     */
    getReferralCodes: (params?: ReferralCodeQueryParams) =>
        apiClient.get<ApiResponse<ReferralCode[]>>('/admin/referrals/codes', { params }),

    /**
     * Get referral code detail by ID
     * GET /api/v1/admin/referrals/codes/:id
     */
    getReferralCode: (id: number) =>
        apiClient.get<ApiResponse<ReferralCode>>(`/admin/referrals/codes/${id}`),

    /**
     * Create a new referral code
     * POST /api/v1/admin/referrals/codes
     */
    createReferralCode: (data: CreateReferralCodeDto) =>
        apiClient.post<ApiResponse<ReferralCode>>('/admin/referrals/codes', data),

    /**
     * Update referral code by ID
     * PUT /api/v1/admin/referrals/codes/:id
     */
    updateReferralCode: (id: number, data: UpdateReferralCodeDto) =>
        apiClient.put<ApiResponse<ReferralCode>>(`/admin/referrals/codes/${id}`, data),

    /**
     * Delete referral code by ID
     * DELETE /api/v1/admin/referrals/codes/:id
     */
    deleteReferralCode: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/referrals/codes/${id}`),

    /**
     * Batch update referral codes status
     * PUT /api/v1/admin/referrals/codes/batch/status
     */
    batchUpdateCodesStatus: (data: BatchUpdateCodesStatusDto) =>
        apiClient.put<ApiResponse<{ successCount: number; failedCount: number }>>('/admin/referrals/codes/batch/status', data),

    /**
     * Batch delete referral codes
     * DELETE /api/v1/admin/referrals/codes/batch
     */
    batchDeleteCodes: (data: BatchDeleteCodesDto) =>
        apiClient.delete<ApiResponse<{ successCount: number; failedCount: number }>>('/admin/referrals/codes/batch', { data }),

    // ========================================================================
    // Referral Record Management
    // ========================================================================

    /**
     * Get referrals list with pagination
     * GET /api/v1/admin/referrals
     */
    getReferrals: (params?: ReferralQueryParams) =>
        apiClient.get<ApiResponse<Referral[]>>('/admin/referrals', { params }),

    /**
     * Get referral detail by ID
     * GET /api/v1/admin/referrals/:id
     */
    getReferral: (id: number) =>
        apiClient.get<ApiResponse<Referral>>(`/admin/referrals/${id}`),

    /**
     * Update referral status by ID
     * PUT /api/v1/admin/referrals/:id/status
     */
    updateReferralStatus: (id: number, data: UpdateReferralStatusDto) =>
        apiClient.put<ApiResponse<void>>(`/admin/referrals/${id}/status`, data),

    /**
     * Batch update referrals status
     * PUT /api/v1/admin/referrals/batch/status
     */
    batchUpdateReferralsStatus: (data: BatchUpdateReferralsStatusDto) =>
        apiClient.put<ApiResponse<{ successCount: number; failedCount: number }>>('/admin/referrals/batch/status', data),

    /**
     * Batch delete referrals
     * DELETE /api/v1/admin/referrals/batch
     */
    batchDeleteReferrals: (data: BatchDeleteReferralsDto) =>
        apiClient.delete<ApiResponse<{ successCount: number; failedCount: number }>>('/admin/referrals/batch', { data }),

    // ========================================================================
    // Reward Management
    // ========================================================================

    /**
     * Get referral rewards list with pagination
     * GET /api/v1/admin/referrals/rewards
     */
    getReferralRewards: (params?: ReferralRewardQueryParams) =>
        apiClient.get<ApiResponse<ReferralReward[]>>('/admin/referrals/rewards', { params }),

    /**
     * Get referral reward detail by ID
     * GET /api/v1/admin/referrals/rewards/:id
     */
    getReferralReward: (id: number) =>
        apiClient.get<ApiResponse<ReferralReward>>(`/admin/referrals/rewards/${id}`),

    /**
     * Issue referral reward by ID
     * POST /api/v1/admin/referrals/rewards/:id/issue
     */
    issueReferralReward: (id: number) =>
        apiClient.post<ApiResponse<void>>(`/admin/referrals/rewards/${id}/issue`),

    /**
     * Mark referral reward as failed by ID
     * POST /api/v1/admin/referrals/rewards/:id/fail
     */
    failReferralReward: (id: number, data: FailReferralRewardDto) =>
        apiClient.post<ApiResponse<void>>(`/admin/referrals/rewards/${id}/fail`, data),

    // ========================================================================
    // Statistics
    // ========================================================================

    /**
     * Get referral statistics
     * GET /api/v1/admin/referrals/stats
     */
    getReferralStats: () =>
        apiClient.get<ApiResponse<ReferralStats>>('/admin/referrals/stats'),
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
 * Get referral type display label
 */
export const getReferralTypeLabel = (type: ReferralType): string => {
    const labels: Record<ReferralType, string> = {
        user: '用户推荐',
        player: '陪玩师推荐',
    };
    return labels[type] || type;
};

/**
 * Get referral status display label
 */
export const getReferralStatusLabel = (status: ReferralStatus): string => {
    const labels: Record<ReferralStatus, string> = {
        pending: '待完成',
        completed: '已完成',
        cancelled: '已取消',
    };
    return labels[status] || status;
};

/**
 * Get referral status color (for Ant Design tags)
 */
export const getReferralStatusColor = (status: ReferralStatus): string => {
    const colors: Record<ReferralStatus, string> = {
        pending: 'orange',
        completed: 'green',
        cancelled: 'red',
    };
    return colors[status] || 'default';
};

/**
 * Get reward type display label
 */
export const getRewardTypeLabel = (type: RewardType): string => {
    const labels: Record<RewardType, string> = {
        referrer: '推荐人奖励',
        referee: '被推荐人奖励',
    };
    return labels[type] || type;
};

/**
 * Get reward status display label
 */
export const getRewardStatusLabel = (status: ReferralRewardStatus): string => {
    const labels: Record<ReferralRewardStatus, string> = {
        pending: '待发放',
        issued: '已发放',
        failed: '发放失败',
    };
    return labels[status] || status;
};

/**
 * Get reward status color (for Ant Design tags)
 */
export const getRewardStatusColor = (status: ReferralRewardStatus): string => {
    const colors: Record<ReferralRewardStatus, string> = {
        pending: 'orange',
        issued: 'green',
        failed: 'red',
    };
    return colors[status] || 'default';
};

/**
 * Check if referral code is expired
 */
export const isCodeExpired = (code: ReferralCode): boolean => {
    return new Date(code.expiresAt) < new Date();
};

/**
 * Check if referral code is fully used
 */
export const isCodeFullyUsed = (code: ReferralCode): boolean => {
    return code.usedCount >= code.maxUses;
};

/**
 * Check if referral code is available for use
 */
export const isCodeAvailable = (code: ReferralCode): boolean => {
    return code.isActive && !isCodeExpired(code) && !isCodeFullyUsed(code);
};

/**
 * Calculate code usage percentage
 */
export const getCodeUsagePercent = (code: ReferralCode): number => {
    if (code.maxUses <= 0) return 0;
    return Math.min(100, Math.round((code.usedCount / code.maxUses) * 100));
};

// ============================================================================
// Default export
// ============================================================================

export default referralApi;
