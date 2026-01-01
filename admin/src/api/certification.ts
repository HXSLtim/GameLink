import apiClient from './client';
import type { ApiResponse, Pagination } from './admin';

// ============================================================================
// Certification Types (认证类型)
// ============================================================================

/**
 * 认证类型
 */
export type CertificationType = 'identity' | 'rank';

/**
 * 认证状态
 */
export type CertificationStatus = 'pending' | 'approved' | 'rejected';

/**
 * 实名认证申请
 */
interface IdentityCertification {
    id: number;
    userId: number;
    realName: string;
    idCardNumber: string;
    idCardFrontUrl: string;
    idCardBackUrl: string;
    status: CertificationStatus;
    rejectReason?: string;
    reviewedAt?: string;
    reviewedBy?: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * 段位认证申请
 */
interface RankCertification {
    id: number;
    userId: number;
    gameType: string;
    currentRank: string;
    targetRank: string;
    screenshotUrls: string[];
    videoUrl?: string;
    additionalInfo?: string;
    status: CertificationStatus;
    rejectReason?: string;
    reviewedAt?: string;
    reviewedBy?: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * 创建实名认证申请 DTO
 */
interface CreateIdentityCertificationDto {
    realName: string;
    idCardNumber: string;
    idCardFrontUrl: string;
    idCardBackUrl: string;
}

/**
 * 创建段位认证申请 DTO
 */
interface CreateRankCertificationDto {
    gameType: string;
    currentRank: string;
    targetRank: string;
    screenshotUrls: string[];
    videoUrl?: string;
    additionalInfo?: string;
}

/**
 * 审核认证申请 DTO
 */
interface ReviewCertificationDto {
    status: 'approved' | 'rejected';
    rejectReason?: string;
}

/**
 * 认证查询参数
 */
interface CertificationQueryParams {
    page?: number;
    page_size?: number;
    userId?: number;
    status?: CertificationStatus;
    type?: CertificationType;
    startTime?: string;
    endTime?: string;
}

// ============================================================================
// Certification API
// ============================================================================

export const certificationApi = {
    // ------------------------------------------------------------------------
    // Identity Certification (实名认证)
    // ------------------------------------------------------------------------

    /**
     * 获取实名认证申请列表
     * GET /api/v1/players/certifications/identity
     */
    getIdentityCertifications: (params?: CertificationQueryParams) =>
        apiClient.get<ApiResponse<IdentityCertification[]>>('/players/certifications/identity', { params }),

    /**
     * 获取实名认证详情
     * GET /api/v1/players/certifications/identity/:id
     */
    getIdentityCertificationDetail: (id: number) =>
        apiClient.get<ApiResponse<IdentityCertification>>(`/players/certifications/identity/${id}`),

    /**
     * 创建实名认证申请
     * POST /api/v1/players/certifications/identity
     */
    createIdentityCertification: (data: CreateIdentityCertificationDto) =>
        apiClient.post<ApiResponse<IdentityCertification>>('/players/certifications/identity', data),

    /**
     * 审核实名认证申请
     * POST /api/v1/players/certifications/identity/:id/review
     */
    reviewIdentityCertification: (id: number, data: ReviewCertificationDto) =>
        apiClient.post<ApiResponse<{ message: string }>>(`/players/certifications/identity/${id}/review`, data),

    // ------------------------------------------------------------------------
    // Rank Certification (段位认证)
    // ------------------------------------------------------------------------

    /**
     * 获取段位认证申请列表
     * GET /api/v1/players/certifications/rank
     */
    getRankCertifications: (params?: CertificationQueryParams) =>
        apiClient.get<ApiResponse<RankCertification[]>>('/players/certifications/rank', { params }),

    /**
     * 获取段位认证详情
     * GET /api/v1/players/certifications/rank/:id
     */
    getRankCertificationDetail: (id: number) =>
        apiClient.get<ApiResponse<RankCertification>>(`/players/certifications/rank/${id}`),

    /**
     * 创建段位认证申请
     * POST /api/v1/players/certifications/rank
     */
    createRankCertification: (data: CreateRankCertificationDto) =>
        apiClient.post<ApiResponse<RankCertification>>('/players/certifications/rank', data),

    /**
     * 审核段位认证申请
     * POST /api/v1/players/certifications/rank/:id/review
     */
    reviewRankCertification: (id: number, data: ReviewCertificationDto) =>
        apiClient.post<ApiResponse<{ message: string }>>(`/players/certifications/rank/${id}/review`, data),

    // ------------------------------------------------------------------------
    // My Certifications (我的认证)
    // ------------------------------------------------------------------------

    /**
     * 获取我的认证状态
     * GET /api/v1/players/certifications/my-status
     */
    getMyCertificationStatus: () =>
        apiClient.get<ApiResponse<{
            identityCertified: boolean;
            rankCertified: boolean;
            identityCertification?: IdentityCertification;
            rankCertification?: RankCertification;
        }>>('/players/certifications/my-status'),
};

// ============================================================================
// Re-exports
// ============================================================================

export type {
    IdentityCertification,
    RankCertification,
    CreateIdentityCertificationDto,
    CreateRankCertificationDto,
    ReviewCertificationDto,
    CertificationQueryParams,
};
