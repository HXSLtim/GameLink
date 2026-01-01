import apiClient from './client';
import type { ApiResponse, Pagination } from './admin';

/**
 * Settlement Company Types
 */

export type CompanyType = 'individual' | 'company';
export type CompanyStatus = 'active' | 'suspended';

export interface SettlementCompany {
    id: number;
    name: string;
    type: CompanyType;
    businessLicense?: string;
    taxNumber?: string;
    bankName?: string;
    bankAccount?: string;
    contactPerson?: string;
    contactPhone?: string;
    status: CompanyStatus;
    playerCount: number;
    createdAt: string;
    updatedAt: string;
}

export interface SettlementCompanyHistory {
    id: number;
    settlementCompanyId: number;
    fieldName: string;
    oldValue: string;
    newValue: string;
    changedBy: number;
    changedByAdmin?: {
        id: number;
        name: string;
        email: string;
    };
    changedAt: string;
}

export interface PlayerCompanyAssignment {
    id: number;
    playerId: number;
    settlementCompanyId: number;
    effectiveDate: string;
    reason: string;
    assignedBy: number;
    assignedByAdmin?: {
        id: number;
        name: string;
        email: string;
    };
    createdAt: string;
    player?: {
        id: number;
        nickname: string;
        user?: {
            name: string;
            phone: string;
        };
    };
    settlementCompany?: {
        id: number;
        name: string;
        type: CompanyType;
    };
}

export interface SettlementCompanyListParams {
    page?: number;
    pageSize?: number;
    status?: CompanyStatus;
    keyword?: string;
    sortBy?: 'name' | 'created_at' | 'player_count';
    sortOrder?: 'asc' | 'desc';
}

export interface CreateSettlementCompanyDto {
    name: string;
    type: CompanyType;
    businessLicense?: string;
    taxNumber?: string;
    bankName?: string;
    bankAccount?: string;
    contactPerson?: string;
    contactPhone?: string;
}

export interface UpdateSettlementCompanyDto {
    name?: string;
    type?: CompanyType;
    businessLicense?: string;
    taxNumber?: string;
    bankName?: string;
    bankAccount?: string;
    contactPerson?: string;
    contactPhone?: string;
}

export interface ToggleStatusPayload {
    enabled: boolean;
}

export interface AssignPlayerToCompanyDto {
    settlementCompanyId: number;
    effectiveDate: string;
    reason: string;
}

export interface BatchAssignPlayersDto {
    playerIds: number[];
    settlementCompanyId: number;
    effectiveDate: string;
    reason: string;
}

export interface BatchUpdateCompanyStatusDto {
    companyIds: number[];
    isActive: boolean;
}

export interface BatchDeleteCompaniesDto {
    companyIds: number[];
}

export interface BatchOperationResult {
    successCount: number;
    failedCount: number;
    totalCount: number;
    failedItems: Array<{ id: number; message: string }>;
    successItems: number[];
}

export interface CompanyPlayersParams {
    page?: number;
    pageSize?: number;
    keyword?: string;
    sortBy?: string;
    sortOrder?: 'asc' | 'desc';
}

/**
 * Settlement Company API
 */
export const settlementApi = {
    /**
     * Get settlement companies list
     * GET /admin/settlement-companies
     */
    getSettlementCompanies: (params?: SettlementCompanyListParams) =>
        apiClient.get<ApiResponse<SettlementCompany[]>>('/admin/settlement-companies', { params }),

    /**
     * Get settlement company detail
     * GET /admin/settlement-companies/:id
     */
    getSettlementCompanyDetail: (id: number) =>
        apiClient.get<ApiResponse<SettlementCompany>>(`/admin/settlement-companies/${id}`),

    /**
     * Create settlement company
     * POST /admin/settlement-companies
     */
    createSettlementCompany: (data: CreateSettlementCompanyDto) =>
        apiClient.post<ApiResponse<SettlementCompany>>('/admin/settlement-companies', data),

    /**
     * Update settlement company
     * PUT /admin/settlement-companies/:id
     */
    updateSettlementCompany: (id: number, data: UpdateSettlementCompanyDto) =>
        apiClient.put<ApiResponse<SettlementCompany>>(`/admin/settlement-companies/${id}`, data),

    /**
     * Toggle settlement company status
     * POST /admin/settlement-companies/:id/toggle
     */
    toggleSettlementCompanyStatus: (id: number, data: ToggleStatusPayload) =>
        apiClient.post<ApiResponse<void>>(`/admin/settlement-companies/${id}/toggle`, data),

    /**
     * Delete settlement company
     * DELETE /admin/settlement-companies/:id
     */
    deleteSettlementCompany: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/settlement-companies/${id}`),

    /**
     * Get settlement company history
     * GET /admin/settlement-companies/:id/history
     */
    getSettlementCompanyHistory: (id: number) =>
        apiClient.get<ApiResponse<SettlementCompanyHistory[]>>(`/admin/settlement-companies/${id}/history`),

    /**
     * Get company players
     * GET /admin/settlement-companies/:companyId/players
     */
    getCompanyPlayers: (companyId: number, params?: CompanyPlayersParams) =>
        apiClient.get<ApiResponse<PlayerCompanyAssignment[]>>(`/admin/settlement-companies/${companyId}/players`, { params }),

    /**
     * Assign player to company
     * POST /admin/players/:id/assign-company
     */
    assignPlayerToCompany: (playerId: number, data: AssignPlayerToCompanyDto) =>
        apiClient.post<ApiResponse<PlayerCompanyAssignment>>(`/admin/players/${playerId}/assign-company`, data),

    /**
     * Batch assign players to company
     * POST /admin/players/batch-assign-company
     */
    batchAssignPlayersToCompany: (data: BatchAssignPlayersDto) =>
        apiClient.post<ApiResponse<{ assignedCount: number; message: string }>>('/admin/players/batch-assign-company', data),

    /**
     * Get player current assignment
     * GET /admin/players/:id/current-company
     */
    getPlayerCurrentAssignment: (playerId: number) =>
        apiClient.get<ApiResponse<PlayerCompanyAssignment>>(`/admin/players/${playerId}/current-company`),

    /**
     * Get player assignment history
     * GET /admin/players/:id/company-history
     */
    getPlayerAssignmentHistory: (playerId: number, params?: { page?: number; pageSize?: number }) =>
        apiClient.get<ApiResponse<{ assignments: PlayerCompanyAssignment[]; total: number }>>(`/admin/players/${playerId}/company-history`, { params }),

    /**
     * Batch update company status
     * POST /admin/settlement-companies/batch/status
     */
    batchUpdateCompanyStatus: (data: BatchUpdateCompanyStatusDto) =>
        apiClient.post<ApiResponse<BatchOperationResult>>('/admin/settlement-companies/batch/status', data),

    /**
     * Batch delete companies
     * POST /admin/settlement-companies/batch/delete
     */
    batchDeleteCompanies: (data: BatchDeleteCompaniesDto) =>
        apiClient.post<ApiResponse<BatchOperationResult>>('/admin/settlement-companies/batch/delete', data),
};

// Re-export types for convenience
export type {
    SettlementCompany,
    SettlementCompanyHistory,
    PlayerCompanyAssignment,
    SettlementCompanyListParams,
    CreateSettlementCompanyDto,
    UpdateSettlementCompanyDto,
    ToggleStatusPayload,
    AssignPlayerToCompanyDto,
    BatchAssignPlayersDto,
    BatchUpdateCompanyStatusDto,
    BatchDeleteCompaniesDto,
    BatchOperationResult,
    CompanyPlayersParams,
};
