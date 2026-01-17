/**
 * Team API
 * Handles team creation, management, orders
 */

import { http } from '@/lib/http';
import type {
    Team,
    CreateTeamRequest,
    TeamMember,
    TeamOrder,
    PaginatedResponse
} from '@/types/api';

export const teamApi = {
    /**
     * Create team
     */
    create: (data: CreateTeamRequest) =>
        http.post<Team>('/team/create', data),

    /**
     * Get team list
     */
    list: (params: {
        page: number;
        pageSize: number;
        status?: string;
    }) =>
        http.get<PaginatedResponse<Team>>('/team/list', { params }),

    /**
     * Get team detail
     */
    get: (id: number) =>
        http.get<Team>(`/team/${id}`),

    /**
     * Join team
     */
    join: (teamId: number) =>
        http.post<void>(`/team/${teamId}/join`),

    /**
     * Leave team
     */
    leave: (teamId: number) =>
        http.post<void>(`/team/${teamId}/leave`),

    /**
     * Kick member (team leader only)
     */
    kickMember: (teamId: number, memberId: number) =>
        http.post<void>(`/team/${teamId}/kick/${memberId}`),

    /**
     * Transfer leadership (team leader only)
     */
    transferLeadership: (teamId: number, newLeaderId: number) =>
        http.post<void>(`/team/${teamId}/transfer-leadership`, { newLeaderId }),

    /**
     * Disband team (team leader only)
     */
    disband: (teamId: number) =>
        http.post<void>(`/team/${teamId}/disband`),

    /**
     * Get team members
     */
    getMembers: (teamId: number) =>
        http.get<TeamMember[]>(`/team/${teamId}/members`),

    /**
     * Get team orders
     */
    getOrders: (teamId: number, params: {
        page: number;
        pageSize: number;
    }) =>
        http.get<PaginatedResponse<TeamOrder>>(`/team/${teamId}/orders`, { params }),

    /**
     * Create team order
     */
    createOrder: (teamId: number, data: {
        serviceItemId: number;
        duration: number;
    }) =>
        http.post<TeamOrder>(`/team/${teamId}/order`, data),
};
