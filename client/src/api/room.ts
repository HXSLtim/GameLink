/**
 * Room API
 * Handles game rooms, room creation, room management
 */

import { http } from '@/lib/http';
import type {
    Room,
    CreateRoomRequest,
    RoomListParams,
    RoomMember,
    PaginatedResponse
} from '@/types/api';

export const roomApi = {
    /**
     * Get room list
     */
    list: (params: RoomListParams) =>
        http.get<PaginatedResponse<Room>>('/room/list', { params }),

    /**
     * Get room detail
     */
    get: (id: number) =>
        http.get<Room>(`/room/${id}`),

    /**
     * Create room
     */
    create: (data: CreateRoomRequest) =>
        http.post<Room>('/room/create', data),

    /**
     * Join room
     */
    join: (roomId: number, password?: string) =>
        http.post<void>(`/room/${roomId}/join`, { password }),

    /**
     * Leave room
     */
    leave: (roomId: number) =>
        http.post<void>(`/room/${roomId}/leave`),

    /**
     * Kick member (room owner only)
     */
    kickMember: (roomId: number, memberId: number) =>
        http.post<void>(`/room/${roomId}/kick/${memberId}`),

    /**
     * Close room (room owner only)
     */
    close: (roomId: number) =>
        http.post<void>(`/room/${roomId}/close`),

    /**
     * Get room members
     */
    getMembers: (roomId: number) =>
        http.get<RoomMember[]>(`/room/${roomId}/members`),

    /**
     * Update room settings (room owner only)
     */
    updateSettings: (roomId: number, data: {
        name?: string;
        description?: string;
        maxMembers?: number;
        isPrivate?: boolean;
        password?: string;
    }) =>
        http.put<Room>(`/room/${roomId}/settings`, data),

    /**
     * Get user's rooms
     */
    getMyRooms: (params: { page: number; pageSize: number }) =>
        http.get<PaginatedResponse<Room>>('/room/my', { params }),
};
