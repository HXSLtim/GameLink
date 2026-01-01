import apiClient from './client';
import type { ApiResponse, PaginatedResponse } from '@/types/api';

// ============================================================================
// Types
// ============================================================================

export type TeamStatus = 'active' | 'busy' | 'inactive';

export type TeamMemberRole = 'leader' | 'member';

export type TeamMemberStatus = 'active' | 'left' | 'kicked';

export type TeamInviteStatus = 'pending' | 'accepted' | 'rejected' | 'expired';

export interface Team {
    id: number;
    name: string;
    description?: string;
    avatarUrl?: string;
    leaderId: number;
    leader?: {
        id: number;
        nickname: string;
        avatar?: string;
        rank?: string;
    };
    status: TeamStatus;
    maxMembers: number;
    memberCount: number;
    incomeShareType: 'equal' | 'custom';
    leaderBonusRate: number;
    totalOrderCount: number;
    totalIncomeCents: number;
    currentOrderId?: number;
    createdAt: string;
    updatedAt: string;
}

export interface TeamMember {
    id: number;
    teamId: number;
    playerId: number;
    player?: {
        id: number;
        nickname: string;
        avatar?: string;
        rank?: string;
    };
    role: TeamMemberRole;
    status: TeamMemberStatus;
    sortOrder: number;
    orderCount: number;
    incomeCents: number;
    joinedAt: string;
    leftAt?: string;
}

export interface TeamInvite {
    id: number;
    teamId: number;
    playerId: number;
    player?: {
        id: number;
        nickname: string;
        avatar?: string;
    };
    inviterId: number;
    inviter?: {
        id: number;
        nickname: string;
        avatar?: string;
    };
    status: TeamInviteStatus;
    expireAt: string;
    message?: string;
    createdAt: string;
    updatedAt: string;
}

export interface TeamStats {
    totalTeams: number;
    activeTeams: number;
    busyTeams: number;
    inactiveTeams: number;
    totalMembers: number;
    totalIncomeCents: number;
}

// ============================================================================
// Request/Response Types
// ============================================================================

export interface TeamListParams {
    page?: number;
    page_size?: number;
    keyword?: string;
    status?: TeamStatus;
    leaderId?: number;
    minMember?: number;
    maxMember?: number;
}

export interface TeamListResponse extends PaginatedResponse {
    items: Team[];
}

export interface TeamCreateRequest {
    name: string;
    description?: string;
    avatarUrl?: string;
    leaderPlayerId: number;
    maxMembers?: number;
    incomeShareType?: 'equal' | 'custom';
    leaderBonusRate?: number;
}

export interface TeamUpdateRequest {
    name: string;
    description?: string;
    avatarUrl?: string;
    maxMembers?: number;
    incomeShareType?: 'equal' | 'custom';
    leaderBonusRate?: number;
}

export interface TeamStatusUpdateRequest {
    status: TeamStatus;
}

export interface MemberListParams {
    page?: number;
    page_size?: number;
    teamId?: number;
    playerId?: number;
    role?: TeamMemberRole;
    status?: TeamMemberStatus;
}

export interface MemberListResponse extends PaginatedResponse {
    items: TeamMember[];
}

export interface AddMemberRequest {
    playerId: number;
}

export interface TransferLeaderRequest {
    newLeaderPlayerId: number;
}

export interface InviteListParams {
    page?: number;
    page_size?: number;
    teamId?: number;
    playerId?: number;
    status?: TeamInviteStatus;
}

export interface InviteListResponse extends PaginatedResponse {
    items: TeamInvite[];
}

// Batch operations

export interface BatchDeleteTeamsRequest {
    team_ids: number[];
}

export interface BatchUpdateTeamsStatusRequest {
    team_ids: number[];
    status: TeamStatus;
}

export interface BatchAddTeamMembersRequest {
    team_id: number;
    player_ids: number[];
}

export interface BatchOperationError {
    id: number;
    message: string;
}

export interface BatchOperationResponse {
    successCount: number;
    failedCount: number;
    totalCount: number;
    successItems: number[];
    failedItems: BatchOperationError[];
}

// ============================================================================
// API
// ============================================================================

export const teamApi = {
    // Team Management
    getTeams: (params?: TeamListParams) =>
        apiClient.get<ApiResponse<TeamListResponse>>('/admin/teams', { params }),

    getTeamDetail: (id: number) =>
        apiClient.get<ApiResponse<Team>>(`/admin/teams/${id}`),

    createTeam: (data: TeamCreateRequest) =>
        apiClient.post<ApiResponse<Team>>('/admin/teams', data),

    updateTeam: (id: number, data: TeamUpdateRequest) =>
        apiClient.put<ApiResponse<Team>>(`/admin/teams/${id}`, data),

    deleteTeam: (id: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/teams/${id}`),

    updateTeamStatus: (id: number, data: TeamStatusUpdateRequest) =>
        apiClient.put<ApiResponse<void>>(`/admin/teams/${id}/status`, data),

    getTeamStats: () =>
        apiClient.get<ApiResponse<TeamStats>>('/admin/teams/stats'),

    // Batch Operations
    batchDeleteTeams: (data: BatchDeleteTeamsRequest) =>
        apiClient.delete<ApiResponse<BatchOperationResponse>>('/admin/teams/batch', { data }),

    batchUpdateTeamsStatus: (data: BatchUpdateTeamsStatusRequest) =>
        apiClient.put<ApiResponse<BatchOperationResponse>>('/admin/teams/batch/status', data),

    batchAddTeamMembers: (data: BatchAddTeamMembersRequest) =>
        apiClient.post<ApiResponse<BatchOperationResponse>>('/admin/teams/batch/members', data),

    // Member Management
    getTeamMembers: (teamId: number) =>
        apiClient.get<ApiResponse<TeamMember[]>>(`/admin/teams/${teamId}/members`),

    addTeamMember: (teamId: number, data: AddMemberRequest) =>
        apiClient.post<ApiResponse<void>>(`/admin/teams/${teamId}/members`, data),

    removeTeamMember: (teamId: number, playerId: number) =>
        apiClient.delete<ApiResponse<void>>(`/admin/teams/${teamId}/members/${playerId}`),

    transferCaptain: (teamId: number, data: TransferLeaderRequest) =>
        apiClient.post<ApiResponse<void>>(`/admin/teams/${teamId}/transfer-leader`, data),

    // Global Member List
    listMembers: (params?: MemberListParams) =>
        apiClient.get<ApiResponse<MemberListResponse>>('/admin/teams/members', { params }),

    // Invite Management
    listInvites: (params?: InviteListParams) =>
        apiClient.get<ApiResponse<InviteListResponse>>('/admin/teams/invites', { params }),

    getInviteDetail: (inviteId: number) =>
        apiClient.get<ApiResponse<TeamInvite>>(`/admin/teams/invites/${inviteId}`),
};

export default teamApi;
